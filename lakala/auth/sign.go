package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
)

type SignatureResult struct {
	Sign                    string
	CertificateSerialNumber string
}

// Signer 签名器接口
type Signer interface {
	Sign(message []byte) (*SignatureResult, error)
}

// PrivateKeySigner 私钥签名器
type PrivateKeySigner struct {
	certificateSerialNumber string
	privateKey              *rsa.PrivateKey
}

func NewPrivateKeySigner(serialNumber string, privateKey *rsa.PrivateKey) *PrivateKeySigner {
	return &PrivateKeySigner{
		certificateSerialNumber: serialNumber,
		privateKey:              privateKey,
	}
}

func (s *PrivateKeySigner) Sign(message []byte) (*SignatureResult, error) {
	hasher := crypto.SHA256.New()
	hasher.Write(message)
	hashed := hasher.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}

	return &SignatureResult{
		Sign:                    base64.StdEncoding.EncodeToString(signature),
		CertificateSerialNumber: s.certificateSerialNumber,
	}, nil
}
