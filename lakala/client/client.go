package client

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"
	"zyj.com/golang-study/lakala/auth"
	"zyj.com/golang-study/lakala/util"
)

// LKLApiClient LKL API客户端
type LKLApiClient struct {
	AppId               string
	HttpClient          *http.Client
	NotificationHandler *auth.NotificationHandler
	Sm4Util             *util.SM4Util
	SdkServerUrl        string
}

func NewLKLApiClient(appId string, httpClient *http.Client, notificationHandler *auth.NotificationHandler,
	sm4Util *util.SM4Util, sdkServerUrl string) *LKLApiClient {
	return &LKLApiClient{
		AppId:               appId,
		HttpClient:          httpClient,
		NotificationHandler: notificationHandler,
		Sm4Util:             sm4Util,
		SdkServerUrl:        sdkServerUrl,
	}
}

// LklHttpClientBuilder HTTP客户端构建器
type LklHttpClientBuilder struct {
	credentials *auth.LklApiCredentials
	validator   auth.Validator
	timeout     time.Duration
}

// Credentials 凭证接口
type Credentials interface {
	GetSchema() string
	GetToken(request *http.Request) (string, error)
	GetOpAppId() string
}

// Validator 验证器接口
type Validator interface {
	Validate(response *http.Response) error
	ValidateRequest(request *http.Request) error
	ValidateNotification(body, authorization string) error
}

func NewLklHttpClientBuilder() *LklHttpClientBuilder {
	return &LklHttpClientBuilder{
		timeout: 30 * time.Second,
	}
}

func (b *LklHttpClientBuilder) WithMerchant(appId, serialNo string, privateKey *rsa.PrivateKey) *LklHttpClientBuilder {
	signer := auth.NewPrivateKeySigner(serialNo, privateKey)
	b.credentials = auth.NewLklApiCredentials(appId, serialNo, signer)
	return b
}

func (b *LklHttpClientBuilder) WithCredentials(credentials *auth.LklApiCredentials) *LklHttpClientBuilder {
	b.credentials = credentials
	return b
}

func (b *LklHttpClientBuilder) WithLklpay(certificates []*x509.Certificate) *LklHttpClientBuilder {
	verifier := auth.NewCertificatesVerifier(certificates)
	b.validator = auth.NewLklApiValidator(verifier)
	return b
}

func (b *LklHttpClientBuilder) WithValidator(validator auth.Validator) *LklHttpClientBuilder {
	b.validator = validator
	return b
}

func (b *LklHttpClientBuilder) WithTimeout(timeout time.Duration) *LklHttpClientBuilder {
	b.timeout = timeout
	return b
}

func (b *LklHttpClientBuilder) Build() *http.Client {
	if b.credentials == nil {
		panic("缺少身份认证信息")
	}
	if b.validator == nil {
		panic("缺少签名验证信息")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	return &http.Client{
		Transport: b.createSignatureTransport(transport, b.credentials, b.validator),
		Timeout:   b.timeout,
	}
}

func (b *LklHttpClientBuilder) createSignatureTransport(transport http.RoundTripper,
	credentials *auth.LklApiCredentials, validator auth.Validator) http.RoundTripper {
	return &SignatureTransport{
		Transport:   transport,
		Credentials: credentials,
		Validator:   validator,
	}
}

// SignatureTransport 签名传输层
type SignatureTransport struct {
	Transport   http.RoundTripper
	Credentials *auth.LklApiCredentials
	Validator   auth.Validator
}

func (t *SignatureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 添加认证头
	token, err := t.Credentials.GetToken(req)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", t.Credentials.GetSchema()+" "+token)
	req.Header.Set("lkl-op-sdk", "lkl-go-sdk-1.0.0")
	req.Header.Set("lkl-op-flowgroup", "NORMAL")
	req.Header.Set("lkl-op-appid", t.Credentials.GetOpAppId())

	// 执行请求
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// 验证响应
	if resp.StatusCode == 200 {
		if err := t.Validator.Validate(resp); err != nil {
			return nil, err
		}
	}

	return resp, nil
}
