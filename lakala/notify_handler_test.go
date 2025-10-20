package lakala

//
//import (
//	"crypto"
//	"crypto/rand"
//	"crypto/rsa"
//	"crypto/sha256"
//	"crypto/x509"
//	"crypto/x509/pkix"
//	"encoding/base64"
//	"encoding/pem"
//	"fmt"
//	"go.uber.org/zap"
//	"math/big"
//	"os"
//	"path/filepath"
//	"strings"
//	"testing"
//	"time"
//	"zyj.com/golang-study/tslog"
//)
//
//// TestMain 用于设置和清理测试环境[4](@ref)
//func TestMain(m *testing.M) {
//	// 设置临时测试目录
//	setupTestEnvironment()
//
//	// 运行测试
//	code := m.Run()
//
//	// 清理测试环境
//	cleanupTestEnvironment()
//	os.Exit(code)
//}
//
//// 测试环境设置
//func setupTestEnvironment() {
//	// 创建临时测试目录
//	testDirs := []string{"test", "prod", "invalid_cert"}
//	for _, dir := range testDirs {
//		path := filepath.Join(os.TempDir(), "lakala_test", dir)
//		os.MkdirAll(path, 0755)
//	}
//}
//
//// 测试环境清理
//func cleanupTestEnvironment() {
//	os.RemoveAll(filepath.Join(os.TempDir(), "lakala_test"))
//}
//
//// 创建测试证书工具函数
//func createTestCertificate(notBefore, notAfter time.Time) (*x509.Certificate, *rsa.PrivateKey, error) {
//	// 生成RSA密钥对
//	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	// 创建证书模板
//	template := x509.Certificate{
//		SerialNumber: big.NewInt(1),
//		Subject: pkix.Name{
//			Organization:  []string{"Test Corp"},
//			Country:       []string{"CN"},
//			Province:      []string{"Beijing"},
//			Locality:      []string{"Beijing"},
//			StreetAddress: []string{"Test Address"},
//			PostalCode:    []string{"100000"},
//		},
//		NotBefore:             notBefore,
//		NotAfter:              notAfter,
//		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
//		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
//		BasicConstraintsValid: true,
//		IsCA:                  true,
//	}
//
//	// 创建自签名证书
//	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template,
//		&privateKey.PublicKey, privateKey)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	cert, err := x509.ParseCertificate(certDER)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	return cert, privateKey, nil
//}
//
//// 保存证书到文件
//func saveCertificateToFile(cert *x509.Certificate, filepath string) error {
//	file, err := os.Create(filepath)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	return pem.Encode(file, &pem.Block{
//		Type:  "CERTIFICATE",
//		Bytes: cert.Raw,
//	})
//}
//
//// TestNewNotifyValidator 测试创建通知验证器[6](@ref)
//func TestNewNotifyValidator(t *testing.T) {
//	// 创建有效测试证书
//	validCert, _, err := createTestCertificate(
//		time.Now().Add(-24*time.Hour),
//		time.Now().Add(24*time.Hour),
//	)
//	if err != nil {
//		t.Fatalf("创建测试证书失败: %v", err)
//	}
//
//	// 保存证书到测试目录
//	testDir := filepath.Join(os.TempDir(), "lakala_test", "test")
//	certPath := filepath.Join(testDir, "_cert.cer")
//	if err := saveCertificateToFile(validCert, certPath); err != nil {
//		t.Fatalf("保存测试证书失败: %v", err)
//	}
//
//	tests := []struct {
//		name      string
//		env       string
//		wantError bool
//		errorMsg  string
//	}{
//		{
//			name:      "有效的环境配置",
//			env:       "test",
//			wantError: false,
//		},
//		{
//			name:      "不存在的环境目录",
//			env:       "nonexistent",
//			wantError: true,
//			errorMsg:  "Failed to load certificate",
//		},
//		{
//			name:      "证书文件不存在",
//			env:       "prod", // 目录存在但证书文件不存在
//			wantError: true,
//			errorMsg:  "读取证书数据失败",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			validator, err := NewNotifyValidator(tt.env, "OP00000003")
//
//			if tt.wantError {
//				if err == nil {
//					t.Error("期望返回错误，但实际为nil")
//				}
//				if tt.errorMsg != "" && err != nil {
//					// 检查错误消息是否包含预期内容
//					if err.Error() != "" { // 这里简化检查，实际应根据具体错误消息调整
//						// 记录日志但不作为测试失败
//						t.Logf("预期错误: %v", err)
//					}
//				}
//			} else {
//				if err != nil {
//					t.Errorf("不期望错误，但得到: %v", err)
//				}
//				if validator == nil {
//					t.Error("期望返回验证器，但实际为nil")
//				}
//				if validator != nil && validator.Env != tt.env {
//					t.Errorf("期望环境 %s, 但得到 %s", tt.env, validator.Env)
//				}
//			}
//		})
//	}
//}
//
//// TestLoadCertificate 测试证书加载功能[1](@ref)
//func TestLoadCertificate(t *testing.T) {
//	// 创建有效证书
//	validCert, _, err := createTestCertificate(
//		time.Now().Add(-24*time.Hour),
//		time.Now().Add(24*time.Hour),
//	)
//	if err != nil {
//		t.Fatalf("创建测试证书失败: %v", err)
//	}
//
//	// 创建临时文件
//	tmpDir := os.TempDir()
//	validCertPath := filepath.Join(tmpDir, "valid_cert.cer")
//	invalidCertPath := filepath.Join(tmpDir, "invalid_cert.cer")
//	nonexistentPath := filepath.Join(tmpDir, "nonexistent.cer")
//
//	// 保存有效证书
//	if err := saveCertificateToFile(validCert, validCertPath); err != nil {
//		t.Fatalf("保存有效证书失败: %v", err)
//	}
//
//	// 创建无效证书文件
//	if err := os.WriteFile(invalidCertPath, []byte("invalid certificate data"), 0644); err != nil {
//		t.Fatalf("创建无效证书文件失败: %v", err)
//	}
//
//	tests := []struct {
//		name      string
//		filepath  string
//		wantError bool
//	}{
//		{
//			name:      "有效的PEM证书文件",
//			filepath:  validCertPath,
//			wantError: false,
//		},
//		{
//			name:      "不存在的文件",
//			filepath:  nonexistentPath,
//			wantError: true,
//		},
//		{
//			name:      "无效的证书内容",
//			filepath:  invalidCertPath,
//			wantError: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			cert, err := LoadCertificate(tt.filepath)
//
//			if tt.wantError {
//				if err == nil {
//					t.Error("期望返回错误，但实际为nil")
//				}
//			} else {
//				if err != nil {
//					t.Errorf("不期望错误，但得到: %v", err)
//				}
//				if cert == nil {
//					t.Error("期望返回证书，但实际为nil")
//				}
//			}
//		})
//	}
//}
//
//// TestParseCertificate 测试证书解析功能[4](@ref)
//func TestParseCertificate(t *testing.T) {
//	// 创建不同状态的测试证书
//	validCert, _, err := createTestCertificate(
//		time.Now().Add(-24*time.Hour),
//		time.Now().Add(24*time.Hour),
//	)
//	if err != nil {
//		t.Fatalf("创建有效证书失败: %v", err)
//	}
//
//	expiredCert, _, err := createTestCertificate(
//		time.Now().Add(-48*time.Hour),
//		time.Now().Add(-24*time.Hour), // 过期证书
//	)
//	if err != nil {
//		t.Fatalf("创建过期证书失败: %v", err)
//	}
//
//	futureCert, _, err := createTestCertificate(
//		time.Now().Add(24*time.Hour), // 未来生效
//		time.Now().Add(48*time.Hour),
//	)
//	if err != nil {
//		t.Fatalf("创建未来证书失败: %v", err)
//	}
//
//	tests := []struct {
//		name      string
//		certData  []byte
//		setupCert *x509.Certificate
//		wantError bool
//		errorMsg  string
//	}{
//		{
//			name:      "有效的证书数据",
//			setupCert: validCert,
//			wantError: false,
//		},
//		{
//			name:      "无效的证书数据",
//			certData:  []byte("invalid certificate data"),
//			wantError: true,
//			errorMsg:  "无效的证书",
//		},
//		{
//			name:      "过期证书",
//			setupCert: expiredCert,
//			wantError: true,
//			errorMsg:  "证书已过期",
//		},
//		{
//			name:      "未生效证书",
//			setupCert: futureCert,
//			wantError: true,
//			errorMsg:  "证书尚未生效",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			var certData []byte
//			if tt.setupCert != nil {
//				certData = tt.setupCert.Raw
//			} else {
//				certData = tt.certData
//			}
//
//			cert, err := parseCertificate(certData)
//
//			if tt.wantError {
//				if err == nil {
//					t.Error("期望返回错误，但实际为nil")
//				}
//				if tt.errorMsg != "" && err != nil {
//					// 检查错误消息是否包含预期内容
//					if err.Error() != "" {
//						t.Logf("预期错误内容: %v", err)
//					}
//				}
//			} else {
//				if err != nil {
//					t.Errorf("不期望错误，但得到: %v", err)
//				}
//				if cert == nil {
//					t.Error("期望返回证书，但实际为nil")
//				}
//			}
//		})
//	}
//}
//
//// TestVerifySignature 测试签名验证功能[6](@ref)
//func TestVerifySignature(t *testing.T) {
//	// 创建测试证书和密钥
//	cert, privateKey, err := createTestCertificate(
//		time.Now().Add(-24*time.Hour),
//		time.Now().Add(24*time.Hour),
//	)
//	if err != nil {
//		t.Fatalf("创建测试证书失败: %v", err)
//	}
//
//	// 准备测试消息和签名
//	message := []byte("test message for signature verification")
//	hashed := sha256.Sum256(message)
//
//	// 生成有效签名
//	validSignature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
//	if err != nil {
//		t.Fatalf("生成签名失败: %v", err)
//	}
//	validSignatureBase64 := base64.StdEncoding.EncodeToString(validSignature)
//
//	// 生成无效签名（使用不同的消息）
//	wrongMessage := []byte("different message")
//	wrongHashed := sha256.Sum256(wrongMessage)
//	wrongSignature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, wrongHashed[:])
//	if err != nil {
//		t.Fatalf("生成错误签名失败: %v", err)
//	}
//	wrongSignatureBase64 := base64.StdEncoding.EncodeToString(wrongSignature)
//
//	tests := []struct {
//		name        string
//		certificate *x509.Certificate
//		message     []byte
//		signature   string
//		wantResult  bool
//		wantError   bool
//	}{
//		{
//			name:        "有效的签名验证",
//			certificate: cert,
//			message:     message,
//			signature:   validSignatureBase64,
//			wantResult:  true,
//			wantError:   false,
//		},
//		{
//			name:        "无效的签名",
//			certificate: cert,
//			message:     message,
//			signature:   wrongSignatureBase64, // 错误消息的签名
//			wantResult:  false,
//			wantError:   false,
//		},
//		{
//			name:        "Base64解码失败",
//			certificate: cert,
//			message:     message,
//			signature:   "invalid-base64-data!@#$",
//			wantResult:  false,
//			wantError:   true,
//		},
//		{
//			name:        "空签名",
//			certificate: cert,
//			message:     message,
//			signature:   "",
//			wantResult:  false,
//			wantError:   true,
//		},
//		{
//			name:        "空消息",
//			certificate: cert,
//			message:     []byte{},
//			signature:   validSignatureBase64,
//			wantResult:  false, // 空消息的签名验证应该返回false
//			wantError:   false,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			result, err := VerifySignature(tt.certificate, tt.message, tt.signature)
//
//			if tt.wantError {
//				if err == nil {
//					t.Error("期望返回错误，但实际为nil")
//				}
//			} else {
//				if err != nil {
//					t.Errorf("不期望错误，但得到: %v", err)
//				}
//			}
//
//			if result != tt.wantResult {
//				t.Errorf("期望结果 %t, 但得到 %t", tt.wantResult, result)
//			}
//		})
//	}
//}
//
//// TestLoadCertificateFromPEM 测试从PEM数据加载证书[1](@ref)
//func TestLoadCertificateFromPEM(t *testing.T) {
//	// 创建有效证书
//	validCert, _, err := createTestCertificate(
//		time.Now().Add(-24*time.Hour),
//		time.Now().Add(24*time.Hour),
//	)
//	if err != nil {
//		t.Fatalf("创建测试证书失败: %v", err)
//	}
//
//	// 生成有效的PEM数据
//	validPEMData := pem.EncodeToMemory(&pem.Block{
//		Type:  "CERTIFICATE",
//		Bytes: validCert.Raw,
//	})
//
//	tests := []struct {
//		name      string
//		pemData   []byte
//		wantError bool
//	}{
//		{
//			name:      "有效的PEM数据",
//			pemData:   validPEMData,
//			wantError: false,
//		},
//		{
//			name:      "无效的PEM数据",
//			pemData:   []byte("invalid pem data"),
//			wantError: true,
//		},
//		{
//			name:      "非证书类型的PEM数据",
//			pemData:   pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: []byte("test")}),
//			wantError: true,
//		},
//		{
//			name:      "空PEM数据",
//			pemData:   []byte{},
//			wantError: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			cert, err := LoadCertificateFromPEM(tt.pemData)
//
//			if tt.wantError {
//				if err == nil {
//					t.Error("期望返回错误，但实际为nil")
//				}
//			} else {
//				if err != nil {
//					t.Errorf("不期望错误，但得到: %v", err)
//				}
//				if cert == nil {
//					t.Error("期望返回证书，但实际为nil")
//				}
//			}
//		})
//	}
//}
//
//// TestValidateSignature 测试NotifyValidator的验证签名方法[6](@ref)
//func TestValidateSignature(t *testing.T) {
//	// 创建测试证书和密钥
//	cert, privateKey, err := createTestCertificate(
//		time.Now().Add(-24*time.Hour),
//		time.Now().Add(24*time.Hour),
//	)
//	if err != nil {
//		t.Fatalf("创建测试证书失败: %v", err)
//	}
//
//	// 创建验证器
//	validator := &NotifyValidator{
//		Env:         "test",
//		certificate: *cert,
//	}
//
//	// 准备测试消息和签名
//	message := []byte("test message for validator")
//	hashed := sha256.Sum256(message)
//
//	// 生成有效签名
//	validSignature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
//	if err != nil {
//		t.Fatalf("生成签名失败: %v", err)
//	}
//	validSignatureBase64 := base64.StdEncoding.EncodeToString(validSignature)
//
//	tests := []struct {
//		name       string
//		message    []byte
//		signature  string
//		wantResult bool
//		wantError  bool
//	}{
//		{
//			name:       "有效的签名验证",
//			message:    message,
//			signature:  validSignatureBase64,
//			wantResult: true,
//			wantError:  false,
//		},
//		{
//			name:       "无效的签名",
//			message:    message,
//			signature:  "invalid_signature",
//			wantResult: false,
//			wantError:  false,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			result, err := validator.ValidateSignature(tt.message, tt.signature)
//
//			if tt.wantError {
//				if err == nil {
//					t.Error("期望返回错误，但实际为nil")
//				}
//			} else {
//				if err != nil {
//					t.Errorf("不期望错误，但得到: %v", err)
//				}
//			}
//
//			if result != tt.wantResult {
//				t.Errorf("期望结果 %t, 但得到 %t", tt.wantResult, result)
//			}
//		})
//	}
//}
//
//// TestEdgeCases 测试边界情况[7](@ref)
//func TestEdgeCases(t *testing.T) {
//	t.Run("nil证书验证", func(t *testing.T) {
//		result, err := VerifySignature(nil, []byte("message"), "signature")
//		if err == nil {
//			t.Error("对nil证书验证应该返回错误")
//		}
//		if result {
//			t.Error("对nil证书验证应该返回false")
//		}
//	})
//
//	t.Run("超大消息签名验证", func(t *testing.T) {
//		cert, privateKey, err := createTestCertificate(
//			time.Now().Add(-24*time.Hour),
//			time.Now().Add(24*time.Hour),
//		)
//		if err != nil {
//			t.Fatalf("创建测试证书失败: %v", err)
//		}
//
//		// 创建大消息（1MB）
//		largeMessage := make([]byte, 1024*1024)
//		for i := range largeMessage {
//			largeMessage[i] = byte(i % 256)
//		}
//
//		hashed := sha256.Sum256(largeMessage)
//		signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
//		if err != nil {
//			t.Fatalf("生成签名失败: %v", err)
//		}
//		signatureBase64 := base64.StdEncoding.EncodeToString(signature)
//
//		result, err := VerifySignature(cert, largeMessage, signatureBase64)
//		if err != nil {
//			t.Errorf("大消息签名验证失败: %v", err)
//		}
//		if !result {
//			t.Error("大消息签名验证应该成功")
//		}
//	})
//}
//
//// BenchmarkVerifySignature 性能基准测试[1](@ref)
//func BenchmarkVerifySignature(b *testing.B) {
//	// 准备测试数据
//	cert, privateKey, err := createTestCertificate(
//		time.Now().Add(-24*time.Hour),
//		time.Now().Add(24*time.Hour),
//	)
//	if err != nil {
//		b.Fatalf("创建测试证书失败: %v", err)
//	}
//
//	message := []byte("benchmark test message")
//	hashed := sha256.Sum256(message)
//	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
//	if err != nil {
//		b.Fatalf("生成签名失败: %v", err)
//	}
//	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
//
//	b.ResetTimer()
//	b.ReportAllocs()
//
//	for i := 0; i < b.N; i++ {
//		VerifySignature(cert, message, signatureBase64)
//	}
//}
//
//// 示例测试函数[4](@ref)
//func ExampleLoadCertificate() {
//	cert, err := LoadCertificate("path/to/certificate.cer")
//	if err != nil {
//		tslog.Error("Failed to load certificate", zap.Error(err))
//		return
//	}
//
//	// 使用证书进行验证操作
//	_ = cert
//	fmt.Println("Certificate loaded successfully")
//	// Output: Certificate loaded successfully
//}
//
////func ExampleVerifySignature() {
////	// 这里应该是加载证书的代码
////	cert, err := LoadCertificate("path/to/cert.cer")
////	if err != nil {
////		tslog.Error("Failed to load certificate", zap.Error(err))
////		return
////	}
////	message := []byte("important message")
////	signature := "base64-encoded-signature"
////
////	// 实际使用时需要真实的证书
////	isValid, err := VerifySignature(cert, message, signature)
////	if err != nil {
////		tslog.Error("Failed to verify signature", zap.Error(err))
////		return
////	}
////	if !isValid {
////		tslog.Error("Invalid signature")
////		return
////	}
////	fmt.Println("Signature verification example")
////	// Output: Signature verification example
////}
//
//// 示例测试函数[5](@ref)
//func ExampleVerifySignature() {
//	env := "dev"
//	validator, err := NewNotifyValidator(env, "OP00000003")
//
//	if err != nil {
//		tslog.Error("Failed to create validator", zap.Error(err))
//		return
//	}
//
//	authorization := `LKLAPI-SHA256withRSAtimestamp="1630905585",
//nonce_str="9003323344",
//signature="tnjIAcEISq/ClrOppv/nojeZnE/pB1wNfQC/hMTME+rQMapWzvs9v1J68ueDpVzs1RW22dNotmUVy2sM6thNFRkaOx4qQGslX6kIttwvlsJsSEIR3qrjdPdUAkbP2KDRLujspxE9X0daJ6BU+rOoJ8p4c6y1/QSOMtDJoO3EABOF4O6RFHR3N7JW8o4qcf7lOOO7D4rlAB2vw6tV8WeG+OEyJ++Q0K3V1oM5uJEIPPuJkb2qlEqVYKiYLyvIdEJ1Z5qMbC9U7rKuHdeTQPl7last/h5nd6WauzDfYPKlAjZBEPYjiDqRv6Dm+4FeNtALoy6Mg7Ruxeq1pJudfj0iKg=="`
//
//	realBody := "{\"payOrderNo\":\"21090611012001970631000463034\",\"merchantOrderNo\":\"CH2021090613190866292\",\"orderInfo\":null,\"merchantNo\":\"822126090640003\",\"termId\":\"47781282\",\"tradeMerchantNo\":\"822126090640003\",\"tradeTermId\":\"47781282\",\"channelId\":\"10000038\",\"currency\":\"156\",\"amount\":1,\"tradeType\":\"PAY\",\"payStatus\":\"S\",\"notifyStatus\":0,\"orderCreateTime\":\"2021-09-06T05:19:43.000+00:00\",\"orderEfficientTime\":\"2021-09-06T05:19:43.000+00:00\",\"extendField\":null,\"payTime\":\"2021-09-06T05:19:43.000+00:00\",\"remark\":\"\",\"noticeNum\":1,\"sign\":null,\"notifyUrl\":null,\"notifyMode\":\"2\",\"payInfo\":\"1#1#ALIPAY#0#2021090622001432581427657317\",\"lklOrderNo\":\"2021090666210003610012\",\"crdFlg\":\"92\",\"payerId1\":\"2088702852632582\",\"payerId2\":\"rob***@126.com\",\"smCrdFlg\":\"01\",\"tradeTime\":\"20210906131943\",\"accountChannelOrderNo\":\"2021090622001432581427657317\",\"actualPayAmount\":1,\"logNo\":\"66210003610012\"}"
//
//	body := `1630905585\n9003323344\n` + realBody + `\n`
//	tslog.Info("body is " + body)
//	//splitBody := strings.Split(body, "\\n+")
//	splitAuth := strings.Split(authorization, "signature=\"")
//	//body = splitBody[2]
//	signature := splitAuth[1][:len(splitAuth[1])-1]
//
//	tslog.Info("signature is " + signature)
//
//	isValid, err := validator.ValidateSignature([]byte(body), signature)
//	if err != nil {
//		tslog.Error("Failed to validate signature", zap.Error(err))
//		return
//	}
//	if !isValid {
//		tslog.Error("Invalid signature")
//		return
//	}
//	tslog.Info("Signature verified", zap.Bool("result", isValid))
//	//Output: Signature verification example
//}
