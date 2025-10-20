package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"strings"
	"sync"
	"time"
)

// CertificatesVerifier 管理平台证书集合
type CertificatesVerifier struct {
	certs map[string]*x509.Certificate
	mu    sync.RWMutex
}

func NewCertificatesVerifier(certificates []*x509.Certificate) *CertificatesVerifier {
	result := &CertificatesVerifier{certs: make(map[string]*x509.Certificate)}
	for _, cert := range certificates {
		result.certs[strings.ToUpper(cert.SerialNumber.Text(16))] = cert
	}
	return result
}

func (v *CertificatesVerifier) Verify(serial string, message []byte, signature string) (bool, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if serial == "" {
		return false, errors.New("缺少证书序列号")
	}
	cert, ok := v.certs[strings.ToUpper(serial)]
	if !ok {
		return false, errors.New("未找到匹配的证书")
	}
	signBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}
	hash := sha256.Sum256(message)
	pub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return false, ErrInvalidPublicKey
	}
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], signBytes)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (v *CertificatesVerifier) ValidCertificate() (*x509.Certificate, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	for _, cert := range v.certs {
		if isCertCurrentlyValid(cert) {
			return cert, nil
		}
	}
	return nil, errors.New("没有有效的平台证书")
}

func isCertCurrentlyValid(cert *x509.Certificate) bool {
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return false
	}
	if now.After(cert.NotAfter) {
		return false
	}
	return true
}
