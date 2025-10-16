package auth

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"time"
)

// Verifier 验证器接口
type Verifier interface {
	Verify(serialNumber string, message []byte, signature string) bool
	GetValidCertificate() *x509.Certificate
}

// CertificatesVerifier 证书验证器
type CertificatesVerifier struct {
	certificates map[string]*x509.Certificate
}

func NewCertificatesVerifier(certificates []*x509.Certificate) *CertificatesVerifier {
	certMap := make(map[string]*x509.Certificate)
	for _, cert := range certificates {
		certMap[cert.SerialNumber.String()] = cert
	}
	return &CertificatesVerifier{certificates: certMap}
}

func (v *CertificatesVerifier) Verify(serialNumber string, message []byte, signature string) bool {
	cert, exists := v.certificates[serialNumber]
	if !exists {
		return false
	}

	return v.verify(cert, message, signature)
}

func (v *CertificatesVerifier) verify(cert *x509.Certificate, message []byte, signature string) bool {
	hasher := crypto.SHA256.New()
	hasher.Write(message)
	hashed := hasher.Sum(nil)

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	err = cert.CheckSignature(x509.SHA256WithRSA, hashed, sigBytes)
	return err == nil
}

func (v *CertificatesVerifier) GetValidCertificate() *x509.Certificate {
	now := time.Now()
	for _, cert := range v.certificates {
		if now.After(cert.NotBefore) && now.Before(cert.NotAfter) {
			return cert
		}
	}
	return nil
}
