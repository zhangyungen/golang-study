package laopsdk

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lakala/laop-sdk-go/pkg/laopsdk/auth"
	"github.com/lakala/laop-sdk-go/pkg/laopsdk/notification"
	"github.com/lakala/laop-sdk-go/pkg/laopsdk/request"
	"github.com/lakala/laop-sdk-go/pkg/laopsdk/utils"
)

const (
	LklOpSDK       = "lkl-go-sdk-1.0.7"
	LklOpFlowGroup = "NORMAL"

	fallbackConnectTimeout = 50000
	fallbackReadTimeout    = 100000
	fallbackSocketTimeout  = 50000
)

type clientEntry struct {
	appID      string
	serverURL  string
	httpClient *http.Client
	sm4        *utils.SM4
	notifier   *notification.Handler
}

var (
	clientsMu  sync.RWMutex
	clients    = make(map[string]*clientEntry)
	defaultApp string
)

func parseInt(value string, fallback int) int {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	v, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return v
}

func durationOrFallback(ms, fallback int) time.Duration {
	if ms <= 0 {
		ms = fallback
	}
	return time.Duration(ms) * time.Millisecond
}

// InitFromProperties 读取properties配置并初始化
func InitFromProperties(path string) error {
	props, err := utils.ReadProperties(path)
	if err != nil {
		return newErrorWithExtra(ErrFileRead, path, err)
	}
	cfg := Config{
		AppID:             props["appId"],
		SerialNo:          props["serialNo"],
		PriKeyPath:        props["priKeyPath"],
		LklCertPath:       props["lklCerPath"],
		LklNotifyCertPath: props["lklNotifyCerPath"],
		SM4Key:            props["sm4Key"],
		ServerURL:         props["serverUrl"],
		ConnectTimeoutMS:  parseInt(props["connectTimeout"], 5000),
		ReadTimeoutMS:     parseInt(props["readTimeout"], 10000),
		SocketTimeoutMS:   parseInt(props["socketTimeout"], 5000),
	}
	return Init(cfg)
}

// Init 通过文件路径配置初始化
func Init(cfg Config) error {
	cfg2, err := cfg.ToConfig2()
	if err != nil {
		return err
	}
	return InitConfig(cfg2)
}

// InitConfig 通过字符串形式的配置初始化SDK
func InitConfig(cfg Config2) error {
	if strings.TrimSpace(cfg.AppID) == "" {
		return newErrorWithExtra(ErrCheckFail, "appId", nil)
	}
	if strings.TrimSpace(cfg.SerialNo) == "" {
		return newErrorWithExtra(ErrCheckFail, "serialNo", nil)
	}
	if strings.TrimSpace(cfg.PriKey) == "" {
		return newErrorWithExtra(ErrCheckFail, "priKey", nil)
	}
	if strings.TrimSpace(cfg.LklCert) == "" {
		return newErrorWithExtra(ErrCheckFail, "lklCert", nil)
	}
	if strings.TrimSpace(cfg.LklNotifyCert) == "" {
		return newErrorWithExtra(ErrCheckFail, "lklNotifyCert", nil)
	}
	if strings.TrimSpace(cfg.ServerURL) == "" {
		return newErrorWithExtra(ErrCheckFail, "serverUrl", nil)
	}
	var sm4Util *utils.SM4
	if strings.TrimSpace(cfg.SM4Key) != "" {
		if !utils.VerifyKey(cfg.SM4Key) {
			return newErrorWithExtra(ErrCheckFail, "sm4Key", nil)
		}
		keyBytes, err := base64.StdEncoding.DecodeString(cfg.SM4Key)
		if err != nil {
			return newErrorWithExtra(ErrCheckFail, "sm4Key", err)
		}
		sm4Util, err = utils.NewSM4(keyBytes)
		if err != nil {
			return newErrorWithExtra(ErrCheckFail, "sm4Key", err)
		}
	}

	privateKey, err := utils.LoadPrivateKey(cfg.PriKey)
	if err != nil {
		return newErrorWithExtra(ErrKeystoreInit, "priKey", err)
	}
	platformCerts, err := utils.LoadCertificates(cfg.LklCert)
	if err != nil {
		return newErrorWithExtra(ErrKeystoreInit, "lklCert", err)
	}
	notifyCert, err := utils.LoadCertificate(cfg.LklNotifyCert)
	if err != nil {
		return newErrorWithExtra(ErrKeystoreInit, "lklNotifyCert", err)
	}

	signer := auth.NewPrivateKeySigner(cfg.SerialNo, privateKey)
	credentials := auth.NewLklApiCredentials(cfg.AppID, cfg.SerialNo, signer)
	verifier := auth.NewCertificatesVerifier(platformCerts)
	validator := auth.NewLklApiValidator(verifier)

	httpClient, err := buildHTTPClient(cfg, credentials, validator)
	if err != nil {
		return err
	}

	notifyVerifier := auth.NewNotifyCertificatesVerifier(notifyCert)
	handler := notification.NewHandler(notifyVerifier)

	entry := &clientEntry{
		appID:      cfg.AppID,
		serverURL:  strings.TrimRight(cfg.ServerURL, "/"),
		httpClient: httpClient,
		sm4:        sm4Util,
		notifier:   handler,
	}

	clientsMu.Lock()
	clients[cfg.AppID] = entry
	if defaultApp == "" {
		defaultApp = cfg.AppID
	}
	clientsMu.Unlock()
	return nil
}

type signatureRoundTripper struct {
	base        http.RoundTripper
	credentials auth.Credentials
	validator   *auth.LklApiValidator
}

func (s *signatureRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if s.base == nil {
		s.base = http.DefaultTransport
	}
	var bodyBytes []byte
	if req.Body != nil {
		data, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, newError(ErrBadRequest, err)
		}
		bodyBytes = data
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		req.ContentLength = int64(len(bodyBytes))
	}
	token, err := s.credentials.Token(bodyBytes)
	if err != nil {
		return nil, newError(ErrBadRequest, err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", s.credentials.Schema(), token))
	req.Header.Set("lkl-op-sdk", LklOpSDK)
	req.Header.Set("lkl-op-flowgroup", LklOpFlowGroup)
	req.Header.Set("lkl-op-appid", s.credentials.AppID())
	if len(bodyBytes) > 0 && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}

	resp, err := s.base.RoundTrip(req)
	if err != nil {
		return nil, newError(ErrPostError, err)
	}
	if resp.StatusCode == http.StatusOK {
		if err := s.validator.ValidateResponse(resp); err != nil {
			resp.Body.Close()
			return nil, newError(ErrBadRequest, err)
		}
		return resp, nil
	}
	respBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	message := strings.TrimSpace(string(respBody))
	if message == "" {
		message = resp.Status
	}
	return nil, newErrorWithExtra(ErrBadRequest, message, nil)
}

func buildHTTPClient(cfg Config2, credentials auth.Credentials, validator *auth.LklApiValidator) (*http.Client, error) {
	var baseTransport http.RoundTripper
	timeout := durationOrFallback(cfg.SocketTimeoutMS, fallbackSocketTimeout)
	var checkRedirect func(*http.Request, []*http.Request) error
	var jar http.CookieJar
	if cfg.HTTPClient != nil {
		if cfg.HTTPClient.Transport != nil {
			baseTransport = cfg.HTTPClient.Transport
		}
		if cfg.HTTPClient.Timeout > 0 {
			timeout = cfg.HTTPClient.Timeout
		}
		checkRedirect = cfg.HTTPClient.CheckRedirect
		jar = cfg.HTTPClient.Jar
	}
	if baseTransport == nil {
		dialer := &net.Dialer{
			Timeout:   durationOrFallback(cfg.ConnectTimeoutMS, fallbackConnectTimeout),
			KeepAlive: 30 * time.Second,
		}
		baseTransport = &http.Transport{
			DialContext:           dialer.DialContext,
			TLSHandshakeTimeout:   durationOrFallback(cfg.ConnectTimeoutMS, fallbackConnectTimeout),
			ResponseHeaderTimeout: durationOrFallback(cfg.ReadTimeoutMS, fallbackReadTimeout),
			ExpectContinueTimeout: 1 * time.Second,
			IdleConnTimeout:       90 * time.Second,
		}
	}
	client := &http.Client{
		Transport:     &signatureRoundTripper{base: baseTransport, credentials: credentials, validator: validator},
		Timeout:       timeout,
		CheckRedirect: checkRedirect,
		Jar:           jar,
	}
	return client, nil
}

func getClient(appID string) (*clientEntry, error) {
	clientsMu.RLock()
	id := strings.TrimSpace(appID)
	if id == "" {
		id = defaultApp
	}
	entry, ok := clients[id]
	clientsMu.RUnlock()
	if !ok || entry == nil {
		if id == "" {
			return nil, newError(ErrSDKNotInit, nil)
		}
		return nil, newErrorWithExtra(ErrAppIDNotInit, id, nil)
	}
	return entry, nil
}

// HTTPPost 默认不加密
func HTTPPost(req request.LklRequest) (string, error) {
	return HTTPPostWithCrypt(req, false, false)
}

// HTTPPostWithCrypt 支持请求加密与响应解密
func HTTPPostWithCrypt(req request.LklRequest, reqEncrypt, respDecrypt bool) (string, error) {
	entry, err := getClient(req.LklAppID())
	if err != nil {
		return "", err
	}

	needCrypt := reqEncrypt || respDecrypt
	if needCrypt && entry.sm4 == nil {
		return "", newErrorWithExtra(ErrSM4InitFail, fmt.Sprintf("appId=%s", entry.appID), nil)
	}

	body, err := req.ToBody()
	if err != nil {
		return "", err
	}

	payload := body
	if reqEncrypt {
		cipherText, err := entry.sm4.EncryptString(body)
		if err != nil {
			return "", err
		}
		payload = cipherText
	}

	url := entry.serverURL + req.FunctionCode().URL()
	log.Printf("POST %s\n %s", url, payload)
	resp, err := doPost(entry.httpClient, url, payload)
	if err != nil {
		return "", err
	}
	if respDecrypt {
		plain, err := entry.sm4.DecryptString(resp)
		if err != nil {
			return "", err
		}
		return plain, nil
	}
	return resp, nil
}

func doPost(client *http.Client, url, body string) (string, error) {
	if client == nil {
		return "", newError(ErrSDKNotInit, nil)
	}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return "", newError(ErrBadRequest, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", newError(ErrPostError, err)
	}
	return string(data), nil
}

// HandleNotificationRequest 处理回调请求
func HandleNotificationRequest(r *http.Request) (string, error) {
	return HandleNotificationRequestWithApp(r, "")
}

// HandleNotificationRequestWithApp 指定appId处理回调请求
func HandleNotificationRequestWithApp(r *http.Request, appID string) (string, error) {
	entry, err := getClient(appID)
	if err != nil {
		return "", err
	}
	if entry.notifier == nil {
		return "", newError(ErrUnknown, fmt.Errorf("未初始化通知证书"))
	}
	return entry.notifier.ParseRequest(r)
}

// HandleNotification 验证通知报文
func HandleNotification(body, authorization string) error {
	return HandleNotificationWithApp(body, authorization, "")
}

// HandleNotificationWithApp 指定appId验证通知报文
func HandleNotificationWithApp(body, authorization, appID string) error {
	entry, err := getClient(appID)
	if err != nil {
		return err
	}
	if entry.notifier == nil {
		return newError(ErrUnknown, fmt.Errorf("未初始化通知证书"))
	}
	return entry.notifier.ValidateBody(body, authorization)
}

// SM4Encrypt 使用默认appId加密
func SM4Encrypt(body string) (string, error) {
	return SM4EncryptWithApp(body, "")
}

// SM4EncryptWithApp 指定appId加密
func SM4EncryptWithApp(body, appID string) (string, error) {
	entry, err := getClient(appID)
	if err != nil {
		return "", err
	}
	if entry.sm4 == nil {
		return "", newErrorWithExtra(ErrSM4InitFail, fmt.Sprintf("appId=%s", entry.appID), nil)
	}
	if strings.TrimSpace(body) == "" {
		return "", newErrorWithExtra(ErrCheckFail, "body为空", nil)
	}
	return entry.sm4.EncryptString(body)
}

// SM4Decrypt 使用默认appId解密
func SM4Decrypt(body string) (string, error) {
	return SM4DecryptWithApp(body, "")
}

// SM4DecryptWithApp 指定appId解密
func SM4DecryptWithApp(body, appID string) (string, error) {
	entry, err := getClient(appID)
	if err != nil {
		return "", err
	}
	if entry.sm4 == nil {
		return "", newErrorWithExtra(ErrSM4InitFail, fmt.Sprintf("appId=%s", entry.appID), nil)
	}
	if strings.TrimSpace(body) == "" {
		return "", newErrorWithExtra(ErrCheckFail, "body为空", nil)
	}
	return entry.sm4.DecryptString(body)
}
