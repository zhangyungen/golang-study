package auth

import (
	"crypto/rand"
	"fmt"
	"time"
)

// Credentials 定义生成请求头的接口
type Credentials interface {
	Schema() string
	Token(body []byte) (string, error)
	AppID() string
}

// LklApiCredentials 实现Authorization生成逻辑
type LklApiCredentials struct {
	appID  string
	serial string
	signer Signer
}

func NewLklApiCredentials(appID, serial string, signer Signer) *LklApiCredentials {
	return &LklApiCredentials{appID: appID, serial: serial, signer: signer}
}

func (c *LklApiCredentials) AppID() string {
	return c.appID
}

func (c *LklApiCredentials) Schema() string {
	return "LKLAPI-SHA256withRSA"
}

func (c *LklApiCredentials) Token(body []byte) (string, error) {
	nonce, err := generateNonce(32)
	if err != nil {
		return "", err
	}
	timestamp := time.Now().Unix()
	message := fmt.Sprintf("%s\n%s\n%d\n%s\n%s\n", c.appID, c.serial, timestamp, nonce, string(body))
	result, err := c.signer.Sign([]byte(message))
	if err != nil {
		return "", err
	}
	token := fmt.Sprintf("appid=\"%s\",nonce_str=\"%s\",timestamp=\"%d\",serial_no=\"%s\",signature=\"%s\"",
		c.appID, nonce, timestamp, result.CertificateSerialNum, result.Sign)
	return token, nil
}

func generateNonce(length int) (string, error) {
	const symbols = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	buf := make([]byte, length)
	random := make([]byte, length)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	for i := range buf {
		buf[i] = symbols[int(random[i])%len(symbols)]
	}
	return string(buf), nil
}
