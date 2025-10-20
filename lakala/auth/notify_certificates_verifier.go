package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

// NotifyCertificatesVerifier 专用于通知验签
type NotifyCertificatesVerifier struct {
	cert *x509.Certificate
}

func NewNotifyCertificatesVerifier(cert *x509.Certificate) *NotifyCertificatesVerifier {
	return &NotifyCertificatesVerifier{cert: cert}
}

func (v *NotifyCertificatesVerifier) Verify(_ string, message []byte, signature string) (bool, error) {
	signBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}
	hash := sha256.Sum256(message)
	pub, ok := v.cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return false, ErrInvalidPublicKey
	}
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], signBytes); err != nil {
		return false, err
	}
	return true, nil
}

func (v *NotifyCertificatesVerifier) ValidCertificate() (*x509.Certificate, error) {
	if !isCertCurrentlyValid(v.cert) {
		return nil, errors.New("证书无效")
	}
	return v.cert, nil
}

var ErrInvalidPublicKey = errors.New("证书公钥类型异常")
