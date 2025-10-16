package main

import (
	"zyj.com/golang-study/lakala/auth"
	"zyj.com/golang-study/lakala/util"
	"zyj.com/golang-study/tslog"
)

func main() {
	// 示例1：处理HTTP通知回调
	//http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
	//	// 加载通知证书
	//	certStr := `-----BEGIN CERTIFICATE-----
	//    ... 证书内容 ...
	//    -----END CERTIFICATE-----`
	//
	//	privateKey, err := lakala.LoadCertificateFromString(certStr)
	//	if err != nil {
	//		http.Error(w, "证书加载失败", http.StatusInternalServerError)
	//		return
	//	}
	//
	//	// 创建通知处理器
	//	handler := auth.NewNotificationHandler(privateKey)
	//
	//	// 解析并验证通知
	//	body, err := handler.Parse(r)
	//	if err != nil {
	//		http.Error(w, "签名验证失败", http.StatusBadRequest)
	//		return
	//	}
	//
	//	fmt.Printf("收到通知: %s\n", body)
	//	w.WriteHeader(http.StatusOK)
	//	w.Write([]byte("SUCCESS"))
	//})
	//
	//
	//log.Println("启动服务器在 :8080")
	//log.Fatal(http.ListenAndServe(":8080", nil))

	// 验证方法3 加签
	privateKey, err := util.LoadPrivateKeyFromString(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDvDBZyHUDndAGx
rIcsCV2njhNO3vCEZotTaWYSYwtDvkcAb1EjsBFabXZaKigpqFXk5XXNI3NIHP9M
8XKzIgGvc65NpLAfRjVql8JiTvLyYd1gIUcOXMInabu+oX7dQSI1mS8XzqaoVRhD
ZQWhXcJW9bxMulgnzvk0Ggw07AjGF7si+hP/Va8SJmN7EJwfQq6TpSxR+WdIHpbW
dhZ+NHwitnQwAJTLBFvfk28INM39G7XOsXdVLfsooFdglVTOHpNuRiQAj9gShCCN
rpGsNQxDiJIxE43qRsNsRwigyo6DPJk/klgDJa417E2wgP8VrwiXparO4FMzOGK1
5quuoD7DAgMBAAECggEBANhmWOt1EAx3OBFf3f4/fEjylQgRSiqRqg8Ymw6KGuh4
mE4Md6eW/B6geUOmZjVP7nIIR1wte28M0REWgn8nid8LGf+v1sB5DmIwgAf+8G/7
qCwd8/VMg3aqgQtRp0ckb5OV2Mv0h2pbnltkWHR8LDIMwymyh5uCApbn/aTrCAZK
NXcPOyAn9tM8Bu3FHk3Pf24Er3SN+bnGxgpzDrFjsDSHjDFT9UMIc2WdA3tuMv9X
3DDn0bRCsHnsIw3WrwY6HQ8mumdbURk+2Ey3eRFfMYxyS96kOgBC2hqZOlDwVPAK
TPtS4hoq+cQ0sRaJQ4T0UALJrBVHa+EESgRaTvrXqAECgYEA+WKmy9hcvp6IWZlk
9Q1JZ+dgIVxrO65zylK2FnD1/vcTx2JMn73WKtQb6vdvTuk+Ruv9hY9PEsf7S8gH
STTmzHOUgo5x0F8yCxXFnfji2juoUnDdpkjtQK5KySDcpQb5kcCJWEVi9v+zObM0
Zr1Nu5/NreE8EqUl3+7MtHOu1TMCgYEA9WM9P6m4frHPW7h4gs/GISA9LuOdtjLv
AtgCK4cW2mhtGNAMttD8zOBQrRuafcbFAyU9de6nhGwetOhkW9YSV+xRNa7HWTeI
RgXJuJBrluq5e1QGTIwZU/GujpNaR4Qiu0B8TodM/FME7htsyxjmCwEfT6SDYlke
MzTbMa9Q0DECgYBqsR/2+dvD2YMwAgZFKKgNAdoIq8dcwyfamUQ5mZ5EtGQL2yw4
8zibHh/LiIxgUD1Kjk/qQgNsX45NP4iOc0mCkrgomtRqdy+rumbPTNmQ0BEVJCBP
scd+8pIgNiTvnWpMRvj7gMP0NDTzLI3wnnCRIq8WAtR2jZ0Ejt+ZHBziLQKBgQDi
bEe/zqNmhDuJrpXEXmO7fTv3YB/OVwEj5p1Z/LSho2nHU3Hn3r7lbLYEhUvwctCn
Ll2fzC7Wic1rsGOqOcWDS5NDrZpUQGGF+yE/JEOiZcPwgH+vcjaMtp0TAfRzuQEz
NzV8YGwxB4mtC7E/ViIuVULHAk4ZGZI8PbFkDxjKgQKBgG8jEuLTI1tsP3kyaF3j
Aylnw7SkBc4gfe9knsYlw44YlrDSKr8AOp/zSgwvMYvqT+fygaJ3yf9uIBdrIilq
CHKXccZ9uA/bT5JfIi6jbg3EoE9YhB0+1aGAS1O2dBvUiD8tJ+BjAT4OB0UDpmM6
QsFLQgFyXgvDnzr/o+hQJelW
-----END PRIVATE KEY-----
`)
	if err != nil {
		tslog.Error("证书加载失败" + err.Error())
		return
	}

	signer := auth.NewPrivateKeySigner("OP00000003", privateKey)
	sign, err := signer.Sign([]byte("hello，rsa，privates"))
	tslog.Info("签名结果: " + sign.Sign)

	// 解签名
	body3 := "{\"payOrderNo\":\"21090611012001970631000463034\",\"merchantOrderNo\":\"CH2021090613190866292\",\"orderInfo\":null,\"merchantNo\":\"822126090640003\",\"termId\":\"47781282\",\"tradeMerchantNo\":\"822126090640003\",\"tradeTermId\":\"47781282\",\"channelId\":\"10000038\",\"currency\":\"156\",\"amount\":1,\"tradeType\":\"PAY\",\"payStatus\":\"S\",\"notifyStatus\":0,\"orderCreateTime\":\"2021-09-06T05:19:43.000+00:00\",\"orderEfficientTime\":\"2021-09-06T05:19:43.000+00:00\",\"extendField\":null,\"payTime\":\"2021-09-06T05:19:43.000+00:00\",\"remark\":\"\",\"noticeNum\":1,\"sign\":null,\"notifyUrl\":null,\"notifyMode\":\"2\",\"payInfo\":\"1#1#ALIPAY#0#2021090622001432581427657317\",\"lklOrderNo\":\"2021090666210003610012\",\"crdFlg\":\"92\",\"payerId1\":\"2088702852632582\",\"payerId2\":\"rob***@126.com\",\"smCrdFlg\":\"01\",\"tradeTime\":\"20210906131943\",\"accountChannelOrderNo\":\"2021090622001432581427657317\",\"actualPayAmount\":1,\"logNo\":\"66210003610012\"}"

	signature := "LKLAPI-SHA256withRSA timestamp=\"1630905585\",\n" +
		"nonce_str=\"9003323344\",\n" +
		"signature=\"tnjIAcEISq/ClrOppv/nojeZnE/pB1wNfQC/hMTME+rQMapWzvs9v1J68ueDpVzs1RW22dNotmUVy2sM6thNFRkaOx4qQGslX6kIttwvlsJsSEIR3qrjdPdUAkbP2KDRLujspxE9X0daJ6BU+rOoJ8p4c6y1/QSOMtDJoO3EABOF4O6RFHR3N7JW8o4qcf7lOOO7D4rlAB2vw6tV8WeG+OEyJ++Q0K3V1oM5uJEIPPuJkb2qlEqVYKiYLyvIdEJ1Z5qMbC9U7rKuHdeTQPl7last/h5nd6WauzDfYPKlAjZBEPYjiDqRv6Dm+4FeNtALoy6Mg7Ruxeq1pJudfj0iKg==\"\n"

	lkNotifycer := "-----BEGIN CERTIFICATE-----\nMIIEMTCCAxmgAwIBAgIGAXRTgcMnMA0GCSqGSIb3DQEBCwUAMHYxCzAJBgNVBAYT\nAkNOMRAwDgYDVQQIDAdCZWlKaW5nMRAwDgYDVQQHDAdCZWlKaW5nMRcwFQYDVQQK\nDA5MYWthbGEgQ28uLEx0ZDEqMCgGA1UEAwwhTGFrYWxhIE9yZ2FuaXphdGlvbiBW\nYWxpZGF0aW9uIENBMB4XDTIwMTAxMDA1MjQxNFoXDTMwMTAwODA1MjQxNFowZTEL\nMAkGA1UEBhMCQ04xEDAOBgNVBAgMB0JlaUppbmcxEDAOBgNVBAcMB0JlaUppbmcx\nFzAVBgNVBAoMDkxha2FsYSBDby4sTHRkMRkwFwYDVQQDDBBBUElHVy5MQUtBTEEu\nQ09NMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt1zHL54HiI8d2sLJ\nlwoQji3/ln0nsvfZ/XVpOjuB+1YR6/0LdxEDMC/hxI6iH2Rm5MjwWz3dmN/6BZeI\ngwGeTOWJUZFARo8UduKrlhC6gWMRpAiiGC8wA8stikc5gYB+UeFVZi/aJ0WN0cpP\nJYCvPBhxhMvhVDnd4hNohnR1L7k0ypuWg0YwGjC25FaNAEFBYP9EYUyCJjE//9Z7\nsMzHR9SJYCqqo6r9bOH9G6sWKuEp+osuAh+kJIxJMHfipw7w3tEcWG0hce9u/el4\ncYJtg8/PPMVoccKmeCzMvarr7jdKP4lenJbtwlgyfs+JgNu60KMUJH8RS72wC9NY\nuFz09wIDAQABo4HVMIHSMIGSBgNVHSMEgYowgYeAFCnH4DkZPR6CZxRn/kIqVsMo\ndJHpoWekZTBjMQswCQYDVQQGEwJDTjEQMA4GA1UECAwHQmVpSmluZzEQMA4GA1UE\nBwwHQmVpSmluZzEXMBUGA1UECgwOTGFrYWxhIENvLixMdGQxFzAVBgNVBAMMDkxh\na2FsYSBSb290IENBggYBaiUALIowHQYDVR0OBBYEFJ2Kx9YZfmWpkKFnC33C0r5D\nK3rFMAwGA1UdEwEB/wQCMAAwDgYDVR0PAQH/BAQDAgeAMA0GCSqGSIb3DQEBCwUA\nA4IBAQBZoeU0XyH9O0LGF9R+JyGwfU/O5amoB97VeM+5n9v2z8OCiIJ8eXVGKN9L\ntl9QkpTEanYwK30KkpHcJP1xfVkhPi/cCMgfTWQ5eKYC7Zm16zk7n4CP6IIgZIqm\nTVGsIGKk8RzWseyWPB3lfqMDR52V1tdA1S8lJ7a2Xnpt5M2jkDXoArl3SVSwCb4D\nAmThYhak48M++fUJNYII9JBGRdRGbfJ2GSFdPXgesUL2CwlReQwbW4GZkYGOg9LK\nCNPK6XShlNdvgPv0CCR08KCYRwC3HZ0y1F0NjaKzYdGNPrvOq9lA495ONZCvzYDo\ngmsu/kd6eqxTs/JwdaIYr4sCMg8Z\n-----END CERTIFICATE-----"
	certificate2, err := util.LoadCertificateFromString(lkNotifycer)
	handler := auth.NewNotificationHandler(certificate2)
	err = handler.Validate(body3, signature)
	if err != nil {
		tslog.Error("验证失败" + err.Error())
		return
	}

	//	解签
	//	V2MmsOpenApiLedgerQueryLedgerMerRequestDemo
	//V2MmsOpenApiLedgerModifyLedgerMerRequestDemo

}

// 示例2：验证字符串格式的通知
func validateNotification(body, authorization, certStr string) error {
	certificate, err := util.LoadCertificateFromString(certStr)
	if err != nil {
		return err
	}
	handler := auth.NewNotificationHandler(certificate)
	return handler.Validate(body, authorization)
}
