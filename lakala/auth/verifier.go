package auth

import "crypto/x509"

// Verifier 定义验签行为
type Verifier interface {
	Verify(serial string, message []byte, signature string) (bool, error)
	ValidCertificate() (*x509.Certificate, error)
}
