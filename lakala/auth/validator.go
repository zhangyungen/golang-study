package auth

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// LklApiValidator 负责HTTP响应与回调验签
type LklApiValidator struct {
	verifier Verifier
}

func NewLklApiValidator(verifier Verifier) *LklApiValidator {
	return &LklApiValidator{verifier: verifier}
}

// ValidateResponse 校验响应头签名并回填body
func (v *LklApiValidator) ValidateResponse(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("响应为空")
	}
	headers := []string{"Lklapi-Appid", "Lklapi-Serial", "Lklapi-Timestamp", "Lklapi-Nonce", "Lklapi-Signature"}
	for _, h := range headers {
		if resp.Header.Get(h) == "" {
			return fmt.Errorf("缺少必要响应头: %s", h)
		}
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	message := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
		resp.Header.Get("Lklapi-Appid"),
		resp.Header.Get("Lklapi-Serial"),
		resp.Header.Get("Lklapi-Timestamp"),
		resp.Header.Get("Lklapi-Nonce"),
		string(bodyBytes))

	ok, err := v.verifier.Verify(resp.Header.Get("Lklapi-Serial"), []byte(message), resp.Header.Get("Lklapi-Signature"))
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("应答验签失败")
	}
	return nil
}

// ValidateNotification 校验通知签名
func (v *LklApiValidator) ValidateNotification(body []byte, authorization string) error {
	authMap, err := parseAuthorization(authorization)
	if err != nil {
		return err
	}
	timestamp, nonce, signature := authMap["timestamp"], authMap["nonce_str"], authMap["signature"]
	if timestamp == "" || nonce == "" || signature == "" {
		return fmt.Errorf("Authorization缺少必要字段")
	}
	message := fmt.Sprintf("%s\n%s\n%s\n", timestamp, nonce, string(body))
	ok, err := v.verifier.Verify("", []byte(message), signature)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("通知验签失败")
	}
	return nil
}

// ValidateRequest 读取请求体并完成验签
func (v *LklApiValidator) ValidateRequest(r *http.Request) ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("请求为空")
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return nil, fmt.Errorf("缺少Authorization头")
	}
	if err := v.ValidateNotification(bodyBytes, authorization); err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func parseAuthorization(header string) (map[string]string, error) {
	result := make(map[string]string)
	header = strings.TrimSpace(header)
	if header == "" {
		return nil, fmt.Errorf("Authorization为空")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("Authorization格式错误")
	}
	result["schema"] = parts[0]
	kvs := strings.Split(parts[1], ",")
	for _, kv := range kvs {
		kv = strings.TrimSpace(kv)
		if kv == "" {
			continue
		}
		idx := strings.Index(kv, "=")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(kv[:idx])
		val := strings.Trim(kv[idx+1:], "\"")
		result[key] = val
	}
	return result, nil
}
