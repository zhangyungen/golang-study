package lakala

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"hash"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"zyj.com/golang-study/util/httpclient"
)

// 接入方唯一编号（appid）	OP00000003	拉卡拉分配
// 证书序列号（serial_no）	00dfba8194c41b84cf	接入方生成的cer证书序列号
// 加签证书	点击下载	接入方生成
// 商户号（merchant_no）	82229007392000A	商户进件产生
// 终端号（term_no）	D9296400	商户进件产生
// 异步通知验签证书	点击下载	拉卡拉分配
// 国密4（SM4Key）	uIj6CPg1GZAY10dXFfsEAQ==	拉卡拉分配
var (
	APPLY_PATH = "/sit/api/v2/mms/openApi/ledger/applyLedgerMer"
	//https://test.wsmsd.cn/sit/api/v2/mms/openApi/ledger/applyLedgerMer
	// APIBaseURL API基础URL
	APIBaseURL     = "https://s2.lakala.com"
	APIBaseURLTEST = "https://test.wsmsd.cn"
	appId          = "OP00000003"
	serialNo       = "00dfba8194c41b84cf"
	merchantNo     = "82229007392000A"
	termNo         = "D9296400"
	SM4Key         = "LHo55AjrT4aDhAIBZhb5KQ=="
)

// Client lakala
type Client struct {
	ctx        context.Context // 上下文
	IsProd     bool            // 是否生产环境
	hc         httpclient.HTTPClient
	sha256Hash hash.Hash
	mu         sync.Mutex
}

func (c *Client) Sign(sParam *SignParam, reqOpts *httpclient.RequestOptions) string {
	str := sParam.ToSignStr()
	c.mu.Lock()
	defer func() {
		c.sha256Hash.Reset()
		c.mu.Unlock()
	}()
	c.sha256Hash.Write([]byte(str))
	sign := strings.ToLower(hex.EncodeToString(c.sha256Hash.Sum(nil)))
	reqOpts.Headers["Authorization"] = sign
	return sign
}

// NewClient 初始化lakala户端
// isProd: 是否生产环境
// baseURL: 基础URL
// s2.lakala.com	主域名	https	拉卡拉开放平台的主域名，优先使用该域名
// openapi.lakala.com	备用	https	拉卡拉开放平台的备用域名，主域名故障时，可采用此域名。
// openapi.lklbiz.com	备用	https	拉卡拉开放平台的备用域名，主域名故障时，可采用此域名。
func NewClient(isProd bool, baseURL string) (client *Client, err error) {
	withTransport := httpclient.WithTransport(defaultClientTransport())
	timeout := httpclient.WithTimeout(0 * time.Second)
	client = &Client{
		ctx:        context.Background(),
		IsProd:     isProd,
		hc:         *httpclient.NewHTTPClient(baseURL, withTransport, timeout),
		sha256Hash: sha256.New(),
	}
	return client, nil
}

func (c *Client) Get(path string, opts *httpclient.RequestOptions) (*httpclient.Response, error) {
	return c.hc.Get(path, opts)
}
func (c *Client) Post(path string, body interface{}, opts *httpclient.RequestOptions) (*httpclient.Response, error) {
	return c.hc.Post(path, body, opts)
}
func (c *Client) Put(path string, body interface{}, opts *httpclient.RequestOptions) (*httpclient.Response, error) {
	return c.hc.Put(path, body, opts)
}
func (c *Client) Delete(path string, opts *httpclient.RequestOptions) (*httpclient.Response, error) {
	return c.hc.Delete(path, opts)
}
func (c *Client) Patch(path string, body interface{}, opts *httpclient.RequestOptions) (*httpclient.Response, error) {
	return c.hc.Patch(path, body, opts)
}

func (c *Client) Head(path string, opts *httpclient.RequestOptions) (*httpclient.Response, error) {
	return c.hc.Head(path, opts)
}
func (c *Client) Options(path string, opts *httpclient.RequestOptions) (*httpclient.Response, error) {
	return c.hc.Options(path, opts)
}
func (c *Client) DownloadFile(url, filepath string, opts *httpclient.RequestOptions) error {
	return c.hc.DownloadFile(url, filepath, opts)
}

/**
 * PostJSON POST JSON 数据
 */
func (c *Client) PostJSON(path string, strJson string, param SignParam) (*httpclient.Response, error) {
	options := &httpclient.RequestOptions{
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "application/json",
			"Authorization": param.ToSignStr(),
		},
	}
	param.body = strJson
	return c.hc.PostJSON(path, strJson, options)
}

func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
func defaultClientTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: defaultTransportDialContext(&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost:   1000,
		MaxConnsPerHost:       3000,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     true,
		ForceAttemptHTTP2:     true,
	}
}

func generateRandomString(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// generateTimestamp13 生成13位时间戳(毫秒级)
func generateTimestamp13() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
}

type SignParam struct {
	Appid     string //应用id
	SerialNo  string //证书序列号
	timeStamp string // 时间戳 13位 毫秒
	nonceStr  string
	body      string
}

func (s *SignParam) ToSignStr() string {
	s.nonceStr = generateRandomString(16)
	s.timeStamp = fmt.Sprintf("%d", time.Now().Unix())
	body := s.body
	sign := s.Appid + "\n" + s.SerialNo + "\n" + s.timeStamp + "\n" + s.nonceStr + "\n" + body + "\n"
	return sign
}

func BuildRequest() *httpclient.RequestOptions {
	options := httpclient.RequestOptions{}
	return &options

}
