package auth

// Signer 定义签名行为
type Signer interface {
	Sign(message []byte) (SignatureResult, error)
}

// SignatureResult 对应Java中的签名结果
type SignatureResult struct {
	Sign                 string
	CertificateSerialNum string
}
