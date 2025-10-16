package util

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/tjfoc/gmsm/sm4"
	"zyj.com/golang-study/lakala"
)

// SM4Util SM4工具类
type SM4Util struct {
	sm4Key []byte
}

const (
	ALGORITHM_NAME     = "SM4"
	ALGORITHM_NAME_ECB = "SM4/ECB/PKCS5Padding"
	DEFAULT_KEY_SIZE   = 128
	ENCODING           = "UTF-8"
)

// NewSM4Util 创建SM4工具实例
func NewSM4Util(key []byte) *SM4Util {
	return &SM4Util{
		sm4Key: key,
	}
}

// NewSM4UtilFromBase64 从Base64字符串创建SM4工具实例
func NewSM4UtilFromBase64(keyBase64 string) (*SM4Util, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, lakala.NewSDKException("Base64解码失败", err)
	}
	return NewSM4Util(key), nil
}

// GenerateKeyToBase64 生成Base64编码的密钥
func GenerateKeyToBase64() (string, error) {
	key, err := GenerateKey()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// GenerateKey 生成密钥
func GenerateKey() ([]byte, error) {
	return GenerateKeyWithSize(DEFAULT_KEY_SIZE / 8) // 128位 = 16字节
}

// GenerateKeyToBase64WithSize 生成指定大小的Base64编码密钥
func GenerateKeyToBase64WithSize(keySize int) (string, error) {
	key, err := GenerateKeyWithSize(keySize)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// GenerateKeyWithSize 生成指定大小的密钥
func GenerateKeyWithSize(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, lakala.NewSDKException("生成随机密钥失败", err)
	}
	return key, nil
}

// Encrypt ECB_PKCS5Padding 加密，返回Base64编码后的密文
func (s *SM4Util) Encrypt(data string) (string, error) {
	if s.sm4Key == nil {
		return "", lakala.NewSDKException("SM4密钥未初始化", nil)
	}

	// SM4 ECB加密
	cipherText, err := sm4.Sm4Ecb(s.sm4Key, []byte(data), true)
	if err != nil {
		return "", lakala.NewSDKException("SM4加密失败", err)
	}

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt ECB_PKCS5Padding 解密，返回明文数据
func (s *SM4Util) Decrypt(cipherText string) (string, error) {
	if s.sm4Key == nil {
		return "", lakala.NewSDKException("SM4密钥未初始化", nil)
	}

	// Base64解码
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", lakala.NewSDKException("Base64解码失败", err)
	}

	// SM4 ECB解密
	plainText, err := sm4.Sm4Ecb(s.sm4Key, cipherBytes, false)
	if err != nil {
		return "", lakala.NewSDKException("SM4解密失败", err)
	}

	return string(plainText), nil
}

// VerifyKey 验证密钥格式
func VerifyKey(key string) bool {
	if key == "" {
		return false
	}

	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return false
	}

	// SM4密钥长度必须是16字节(128位)
	return len(keyBytes) == DEFAULT_KEY_SIZE/8
}

// GetKey 获取密钥
func (s *SM4Util) GetKey() []byte {
	return s.sm4Key
}

// GetKeyBase64 获取Base64编码的密钥
func (s *SM4Util) GetKeyBase64() string {
	return base64.StdEncoding.EncodeToString(s.sm4Key)
}
