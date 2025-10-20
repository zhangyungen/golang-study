package laopsdk

import (
	"fmt"
	"net/http"

	"github.com/lakala/laop-sdk-go/pkg/laopsdk/utils"
)

// Config 对应Java版本以文件路径配置证书的定义
type Config struct {
	AppID             string
	SerialNo          string
	PriKeyPath        string
	LklCertPath       string
	LklNotifyCertPath string
	SM4Key            string
	ConnectTimeoutMS  int
	ReadTimeoutMS     int
	SocketTimeoutMS   int
	ServerURL         string
	HTTPClient        *http.Client
}

// Config2 对应Java版本支持直接注入证书内容的定义
type Config2 struct {
	AppID            string
	SerialNo         string
	PriKey           string
	LklCert          string
	LklNotifyCert    string
	SM4Key           string
	ConnectTimeoutMS int
	ReadTimeoutMS    int
	SocketTimeoutMS  int
	ServerURL        string
	HTTPClient       *http.Client
}

// ToConfig2 将文件路径配置转换为直接注入内容的Config2
func (c Config) ToConfig2() (Config2, error) {
	priKey, err := utils.ReadFileString(c.PriKeyPath)
	if err != nil {
		return Config2{}, newErrorWithExtra(ErrKeystoreInit, fmt.Sprintf("priKeyPath=%s", c.PriKeyPath), err)
	}
	lklCert, err := utils.ReadFileString(c.LklCertPath)
	if err != nil {
		return Config2{}, newErrorWithExtra(ErrKeystoreInit, fmt.Sprintf("lklCertPath=%s", c.LklCertPath), err)
	}
	lklNotifyCert, err := utils.ReadFileString(c.LklNotifyCertPath)
	if err != nil {
		return Config2{}, newErrorWithExtra(ErrKeystoreInit, fmt.Sprintf("lklNotifyCertPath=%s", c.LklNotifyCertPath), err)
	}
	return Config2{
		AppID:            c.AppID,
		SerialNo:         c.SerialNo,
		PriKey:           priKey,
		LklCert:          lklCert,
		LklNotifyCert:    lklNotifyCert,
		SM4Key:           c.SM4Key,
		ConnectTimeoutMS: c.ConnectTimeoutMS,
		ReadTimeoutMS:    c.ReadTimeoutMS,
		SocketTimeoutMS:  c.SocketTimeoutMS,
		ServerURL:        c.ServerURL,
		HTTPClient:       c.HTTPClient,
	}, nil
}
