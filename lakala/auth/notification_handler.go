package auth

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"zyj.com/golang-study/lakala"
	"zyj.com/golang-study/tslog"
)

// NotificationHandler 通知处理器
type NotificationHandler struct {
	validator LklApiValidator
}

// NewNotificationHandler 创建通知处理器
func NewNotificationHandler(certificate *x509.Certificate) *NotificationHandler {
	verifier := NewNotifyCertificatesVerifier(certificate)
	validator := NewLklApiValidator(verifier)

	return &NotificationHandler{
		validator: *validator,
	}
}

// Parse 解析HTTP请求并验证签名
func (h *NotificationHandler) Parse(request *http.Request) (string, error) {
	body, err := h.validator.ValidateRequest(request)
	if err != nil {
		return "", lakala.NewSDKException("推送通知的签名验证失败", err)
	}
	return body, nil
}

// Validate 验证通知内容和授权信息
func (h *NotificationHandler) Validate(body, authorization string) error {
	if !h.validator.ValidateNotification(body, authorization) {
		return lakala.NewSDKException("推送通知的签名验证失败", nil)
	}
	return nil
}

// LklApiValidator LKL API验证器
type LklApiValidator struct {
	verifier Verifier
}

// NewLklApiValidator 创建LKL API验证器
func NewLklApiValidator(verifier Verifier) *LklApiValidator {
	return &LklApiValidator{
		verifier: verifier,
	}
}

// Validate 验证HTTP响应
func (v *LklApiValidator) Validate(response *http.Response) error {
	if err := v.validateParameters(response); err != nil {
		return err
	}

	message, err := v.buildMessage(response)
	if err != nil {
		return err
	}

	serial := response.Header.Get("Lklapi-Serial")
	signature := response.Header.Get("Lklapi-Signature")

	if !v.verifier.Verify(serial, []byte(message), signature) {
		return fmt.Errorf("签名验证失败: serial=[%s] message=[%s] sign=[%s]",
			serial, message, signature)
	}

	return nil
}

// ValidateNotification 验证通知内容
func (v *LklApiValidator) ValidateNotification(body, authorization string) bool {
	authorizationMap := v.getAuthorizationMap(authorization)
	message := v.buildNotificationMessage(body, authorizationMap)
	signature := authorizationMap["signature"]

	return v.verifier.Verify("", []byte(message), signature)
}

// ValidateRequest 验证HTTP请求
func (v *LklApiValidator) ValidateRequest(request *http.Request) (string, error) {
	body, err := v.getRequestBody(request)
	if err != nil {
		return "", err
	}

	authorization := request.Header.Get("Authorization")
	authorizationMap := v.getAuthorizationMap(authorization)
	message := v.buildNotificationMessage(body, authorizationMap)
	signature := authorizationMap["signature"]

	if !v.verifier.Verify("", []byte(message), signature) {
		return "", fmt.Errorf("请求签名验证失败: message=[%s] sign=[%s]", message, signature)
	}

	return body, nil
}

// validateParameters 验证响应参数
func (v *LklApiValidator) validateParameters(response *http.Response) error {
	if response.Header.Get("Lklapi-Serial") == "" {
		return fmt.Errorf("缺少Lklapi-Serial头")
	}
	if response.Header.Get("Lklapi-Signature") == "" {
		return fmt.Errorf("缺少Lklapi-Signature头")
	}
	if response.Header.Get("Lklapi-Timestamp") == "" {
		return fmt.Errorf("缺少Lklapi-Timestamp头")
	}
	if response.Header.Get("Lklapi-Nonce") == "" {
		return fmt.Errorf("缺少Lklapi-Nonce头")
	}
	return nil
}

// buildMessage 构建响应验证消息
func (v *LklApiValidator) buildMessage(response *http.Response) (string, error) {
	timestamp := response.Header.Get("Lklapi-Timestamp")
	nonce := response.Header.Get("Lklapi-Nonce")
	serialNo := response.Header.Get("Lklapi-Serial")
	appid := response.Header.Get("Lklapi-Appid")

	body, err := v.getResponseBody(response)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", appid, serialNo, timestamp, nonce, body), nil
}

// buildNotificationMessage 构建通知验证消息
func (v *LklApiValidator) buildNotificationMessage(body string, authorizationMap map[string]string) string {
	timestamp := authorizationMap["timestamp"]
	nonce := authorizationMap["nonce_str"]
	return fmt.Sprintf("%s\n%s\n%s\n", timestamp, nonce, body)
}

// getResponseBody 获取响应体
func (v *LklApiValidator) getResponseBody(response *http.Response) (string, error) {
	if response.Body == nil {
		return "", nil
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// 由于响应体已经被读取，需要重新设置以便后续使用
	response.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	return string(bodyBytes), nil
}

// getRequestBody 获取请求体
func (v *LklApiValidator) getRequestBody(request *http.Request) (string, error) {
	if request.Body == nil {
		return "", nil
	}

	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return "", err
	}

	// 重新设置请求体以便后续使用
	request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	return string(bodyBytes), nil
}

// getAuthorizationMap 解析授权头信息
func (v *LklApiValidator) getAuthorizationMap(authorization string) map[string]string {
	result := make(map[string]string)

	authorization = strings.TrimSpace(authorization)
	if authorization == "" {
		return result
	}

	// 解析授权类型
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) < 2 {
		return result
	}

	authType := parts[0]
	signInfo := parts[1]

	// 解析授权类型中的算法信息
	typeParts := strings.Split(authType, "-")
	if len(typeParts) > 1 {
		result["signSystemCode"] = typeParts[0]
		result["signAlgorithm"] = typeParts[1]
	}

	// 解析签名信息
	pairs := strings.Split(signInfo, ",")
	for _, pair := range pairs {
		if strings.Contains(pair, "=") {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				// 去除值的引号
				if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
					value = value[1 : len(value)-1]
				}
				result[key] = value
			}
		}
	}

	return result
}

// NotifyCertificatesVerifier 通知证书验证器
type NotifyCertificatesVerifier struct {
	notifyCertificate *x509.Certificate
}

// NewNotifyCertificatesVerifier 创建通知证书验证器
func NewNotifyCertificatesVerifier(certificate *x509.Certificate) *NotifyCertificatesVerifier {
	return &NotifyCertificatesVerifier{
		notifyCertificate: certificate,
	}
}

// Verify 验证签名
func (v *NotifyCertificatesVerifier) Verify(serialNumber string, message []byte, signature string) bool {
	return v.verify(message, signature)
}

// verify 内部验证方法
func (v *NotifyCertificatesVerifier) verify(message []byte, signature string) bool {
	hasher := crypto.SHA256.New()
	hasher.Write(message)
	hashed := hasher.Sum(nil)

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	err = v.notifyCertificate.CheckSignature(x509.SHA256WithRSA, hashed, sigBytes)
	if err != nil {
		tslog.Error("推送通知签名验证失败", zap.Error(err))
		return false
	}
	return true
}

// GetValidCertificate 获取有效证书
func (v *NotifyCertificatesVerifier) GetValidCertificate() *x509.Certificate {
	return v.notifyCertificate
}

// Validator 验证器接口
type Validator interface {
	// Validate 验证HTTP响应
	Validate(response *http.Response) error

	// ValidateRequest 验证HTTP请求并返回请求体
	ValidateRequest(request *http.Request) (string, error)

	// ValidateNotification 验证通知内容
	ValidateNotification(body, authorization string) bool
}
