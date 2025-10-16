package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"zyj.com/golang-study/lakala"
)

// pemUtil PEM工具类
type PemUtil struct{}

// LoadPrivateKeyFromString 从字符串加载私钥
func (p *PemUtil) LoadPrivateKeyFromString(privateKeyStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, lakala.NewSDKException("PEM解码失败: 无法找到PEM块", nil)
	}

	// 尝试PKCS1格式
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privateKey, nil
	}

	// 尝试PKCS8格式
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, lakala.NewSDKException("解析私钥失败", err)
	}

	return key.(*rsa.PrivateKey), nil
}

// LoadPrivateKeyFromFile 从文件加载私钥
func (p *PemUtil) LoadPrivateKeyFromFile(filePath string) (*rsa.PrivateKey, error) {
	content, err := ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return p.LoadPrivateKeyFromString(content)
}

// LoadCertificateFromString 从字符串加载证书
func (p *PemUtil) LoadCertificateFromString(certStr string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certStr))
	if block == nil {
		return nil, lakala.NewSDKException("PEM解码失败: 无法找到PEM块", nil)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, lakala.NewSDKException("解析证书失败", err)
	}

	return cert, nil
}

// LoadCertificateFromFile 从文件加载证书
func (p *PemUtil) LoadCertificateFromFile(filePath string) (*x509.Certificate, error) {
	content, err := ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return p.LoadCertificateFromString(content)
}

// 全局工具实例
var pemUtil = &PemUtil{}

// LoadPrivateKeyFromString 全局函数：从字符串加载私钥
func LoadPrivateKeyFromString(privateKeyStr string) (*rsa.PrivateKey, error) {
	return pemUtil.LoadPrivateKeyFromString(privateKeyStr)
}

// LoadPrivateKeyFromFile 全局函数：从文件加载私钥
func LoadPrivateKeyFromFile(filePath string) (interface{}, error) {
	return pemUtil.LoadPrivateKeyFromFile(filePath)
}

// LoadCertificateFromString 全局函数：从字符串加载证书
func LoadCertificateFromString(certStr string) (*x509.Certificate, error) {
	return pemUtil.LoadCertificateFromString(certStr)
}

// LoadCertificateFromFile 全局函数：从文件加载证书
func LoadCertificateFromFile(filePath string) (*x509.Certificate, error) {
	return pemUtil.LoadCertificateFromFile(filePath)
}
