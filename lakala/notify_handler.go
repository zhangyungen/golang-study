package lakala

//
//import (
//	"crypto"
//	"crypto/rsa"
//	"crypto/sha256"
//	"crypto/x509"
//	"encoding/base64"
//	"encoding/pem"
//	"fmt"
//	"github.com/pkg/errors"
//	"go.uber.org/zap"
//	"io"
//	"os"
//	"time"
//	"zyj.com/golang-study/tslog"
//)
//
//var (
//	// 获取当前目录
//	currentDir, _  = os.Getwd()
//	certName       = "_cert.cer"
//	privateKeyName = "_private_key.pem"
//)
//
//type NotifyValidator struct {
//	Env          string
//	CertFilePath string
//	certificate  x509.Certificate
//}
//
//func NewNotifyValidator(env string, appId string) (*NotifyValidator, error) {
//	filePath := currentDir + "/env/" + env + "/" + appId + certName
//	cert, err := LoadCertificate(filePath)
//	if err != nil {
//		tslog.Error("Failed to load certificate", zap.String("file_path", filePath), zap.Error(err))
//		return nil, err
//	}
//	return &NotifyValidator{
//		Env:          env,
//		CertFilePath: filePath,
//		certificate:  *cert,
//	}, nil
//}
//
//func (h *NotifyValidator) ValidateSignature(bodyMsg []byte, signature string) (bool, error) {
//	return VerifySignature(&h.certificate, bodyMsg, signature)
//}
//
//var notifyValidator NotifyValidator
//
//// LoadCertificate 从输入流加载X.509证书
//func LoadCertificate(filepath string) (*x509.Certificate, error) {
//	file, err := os.Open(filepath)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	data, err := io.ReadAll(file)
//	if err != nil {
//		return nil, fmt.Errorf("读取证书数据失败: %v", err)
//	}
//	// 尝试解析PEM格式（常见格式）
//	block, rest := pem.Decode(data)
//	if block != nil {
//		// 如果是PEM格式，使用PEM块中的数据
//		return parseCertificate(block.Bytes)
//	}
//
//	// 如果不是PEM格式，直接解析DER数据
//	if len(rest) == 0 {
//		return parseCertificate(data)
//	}
//
//	// 如果PEM解码后还有剩余数据，说明可能不是有效的证书格式
//	return nil, fmt.Errorf("无效的证书格式")
//}
//
//// parseCertificate 实际解析证书的函数
//func parseCertificate(certData []byte) (*x509.Certificate, error) {
//	// 解析X.509证书
//	cert, err := x509.ParseCertificate(certData)
//	if err != nil {
//		return nil, fmt.Errorf("无效的证书: %v", err)
//	}
//
//	// 获取公钥的Base64编码
//	publicKeyDER, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
//	if err != nil {
//		return nil, fmt.Errorf("无法编码公钥: %v", err)
//	}
//	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyDER)
//	fmt.Println(publicKeyBase64)
//
//	// 检查证书有效期（对应Java的checkValidity()）
//	now := time.Now()
//	if now.Before(cert.NotBefore) {
//		return nil, fmt.Errorf("证书尚未生效")
//	}
//	if now.After(cert.NotAfter) {
//		return nil, fmt.Errorf("证书已过期")
//	}
//
//	return cert, nil
//}
//
//// VerifySignature 验证RSA签名（对应Java的SHA256withRSA）
//// certificate: X509证书对象
//// message: 原始消息字节数组
//// signature: Base64编码的签名字符串
//// 返回值: 验证成功返回true，失败返回false，错误时返回error
//func VerifySignature(certificate *x509.Certificate, message []byte, signature string) (bool, error) {
//	// 从证书中获取公钥
//	publicKey, ok := certificate.PublicKey.(*rsa.PublicKey)
//	if !ok {
//		return false, errors.New("证书公钥不是RSA类型")
//	}
//
//	// 解码Base64签名
//	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
//	if err != nil {
//		return false, fmt.Errorf("签名Base64解码失败: %v", err)
//	}
//
//	// 计算消息的SHA256哈希值
//	hashed := sha256.Sum256(message)
//
//	// 使用RSA公钥验证PKCS#1 v1.5签名
//	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signatureBytes)
//	if err != nil {
//		if errors.Is(err, rsa.ErrVerification) {
//			// 签名验证不通过（不是错误，是正常的验证失败）
//			return false, nil
//		}
//		// 其他验证错误
//		return false, fmt.Errorf("签名验证过程发生错误: %v", err)
//	}
//
//	return true, nil
//}
//
//// 辅助函数：从PEM格式证书数据加载证书
//func LoadCertificateFromPEM(pemData []byte) (*x509.Certificate, error) {
//	block, _ := pem.Decode(pemData)
//	if block == nil {
//		return nil, errors.New("无法解析PEM数据")
//	}
//
//	certificate, err := x509.ParseCertificate(block.Bytes)
//	if err != nil {
//		return nil, fmt.Errorf("解析证书失败: %v", err)
//	}
//
//	return certificate, nil
//}
//
////// 读取配置文件
////func ReadToByte(filepath string) ([]byte, error) {
////	file, err := os.Open(filepath)
////	if err != nil {
////		return nil, err
////	}
////	defer file.Close()
////
////	data, err := io.ReadAll(file)
////	if err != nil {
////		logs.Error("Failed to read file", zap.String("file", filepath))
////		return data, err
////	}
////	// 解析配置数据...
////	return data, nil
////}
//
//// // 定义一个处理函数，例如用于 POST /note_sync
////
////	func (h *NotifyValidator) NotifyHandler(c *gin.Context) {
////		// 1. 记录请求IP和路径（对应Java中的System.out.println）
////		clientIP := c.ClientIP()
////		c.String(http.StatusOK, "%s:note_sync\n", clientIP) // 简单返回IP信息，通常用日志记录
////		// 2. 获取请求体
////		body, err := io.ReadAll(c.Request.Body) // 使用io.ReadAll读取整个请求体
////		if err != nil {
////			// 处理读取错误
////			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
////			return
////		}
////		// 重要：读取后，需要将Body重置，以便后续可能再次读取
////		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
////		// 3. 此时body是[]byte类型，可以直接使用或转换为字符串
////		bodyString := string(body)
////
////		h.ValidateSignature([]byte(body), bodyString)
////		// 4. 这里简单地将接收到的body内容返回，实际应用中会进行业务处理
////		c.String(http.StatusOK, "Received body: %s", bodyString)
////	}
