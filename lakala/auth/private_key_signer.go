package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// PrivateKeySigner 使用RSA私钥实现SHA256withRSA
type PrivateKeySigner struct {
	serialNumber string
	privateKey   *rsa.PrivateKey
}

func NewPrivateKeySigner(serial string, key *rsa.PrivateKey) *PrivateKeySigner {
	return &PrivateKeySigner{serialNumber: serial, privateKey: key}
}

func (s *PrivateKeySigner) Sign(message []byte) (SignatureResult, error) {
	hash := sha256.Sum256(message)
	sign, err := rsa.SignPKCS1v15(nil, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return SignatureResult{}, fmt.Errorf("签名计算失败: %w", err)
	}
	return SignatureResult{
		Sign:                 base64.StdEncoding.EncodeToString(sign),
		CertificateSerialNum: s.serialNumber,
	}, nil
}
