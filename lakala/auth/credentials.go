package auth

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var symbols = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// LklApiCredentials LKL API凭证
type LklApiCredentials struct {
	appId    string
	signer   Signer
	serialNo string
}

func NewLklApiCredentials(appId, serialNo string, signer Signer) *LklApiCredentials {
	return &LklApiCredentials{
		appId:    appId,
		signer:   signer,
		serialNo: serialNo,
	}
}

func (c *LklApiCredentials) GetSchema() string {
	return "LKLAPI-SHA256withRSA"
}

func (c *LklApiCredentials) GetToken(request *http.Request) (string, error) {
	nonceStr := c.generateNonceStr()
	timestamp := c.generateTimestamp()

	message, err := c.buildMessage(c.appId, c.serialNo, nonceStr, timestamp, request)
	if err != nil {
		return "", err
	}

	signature, err := c.signer.Sign([]byte(message))
	if err != nil {
		return "", err
	}

	token := fmt.Sprintf("appid=\"%s\",nonce_str=\"%s\",timestamp=\"%s\",serial_no=\"%s\",signature=\"%s\"",
		c.appId, nonceStr, strconv.FormatInt(timestamp, 10), signature.CertificateSerialNumber, signature.Sign)

	return token, nil
}

func (c *LklApiCredentials) GetOpAppId() string {
	return c.appId
}

func (c *LklApiCredentials) generateTimestamp() int64 {
	return time.Now().Unix()
}

func (c *LklApiCredentials) generateNonceStr() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为后备
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	for i, b := range bytes {
		bytes[i] = symbols[b%byte(len(symbols))]
	}
	return string(bytes)
}

func (c *LklApiCredentials) buildMessage(appid, serialNo, nonce string, timestamp int64, request *http.Request) (string, error) {
	var body string
	if request.Body != nil {
		bodyBytes, err := io.ReadAll(request.Body)
		if err != nil {
			return "", err
		}
		body = string(bodyBytes)
		// 重置request.Body，因为已经被读取
		request.Body = io.NopCloser(strings.NewReader(body))
	}

	return fmt.Sprintf("%s\n%s\n%d\n%s\n%s\n", appid, serialNo, timestamp, nonce, body), nil
}
