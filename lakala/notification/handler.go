package notification

import (
	"net/http"

	"github.com/lakala/laop-sdk-go/pkg/laopsdk/auth"
)

// Handler 负责通知验签解析
type Handler struct {
	validator *auth.LklApiValidator
}

func NewHandler(notifyVerifier auth.Verifier) *Handler {
	return &Handler{validator: auth.NewLklApiValidator(notifyVerifier)}
}

// ParseRequest 验签并返回正文
func (h *Handler) ParseRequest(r *http.Request) (string, error) {
	body, err := h.validator.ValidateRequest(r)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// ValidateBody 直接验证报文与Authorization
func (h *Handler) ValidateBody(body, authorization string) error {
	return h.validator.ValidateNotification([]byte(body), authorization)
}
