package request

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/lakala/laop-sdk-go/pkg/laopsdk/enums"
	"github.com/lakala/laop-sdk-go/pkg/laopsdk/utils"
)

// LklRequest 对应Java接口
type LklRequest interface {
	FunctionCode() enums.FunctionCode
	ToBody() (string, error)
	LklAppID() string
}

// BaseRequest 提供通用appId维护
type BaseRequest struct {
	appID string
}

func (b *BaseRequest) SetLklAppID(appID string) {
	b.appID = appID
}

func (b *BaseRequest) LklAppID() string {
	return b.appID
}

// V2Common 封装V2公共报文
type V2Common struct {
	BaseRequest
}

func BuildV2Body(data interface{}) (string, error) {
	envelope := map[string]interface{}{
		"reqId":   randomHex(32),
		"reqTime": time.Now().Format("20060102150405"),
		"version": "2.0",
		"reqData": data,
	}
	return utils.ToJSONString(envelope)
}

func (c *V2Common) BuildBody(data interface{}) (string, error) {
	return BuildV2Body(data)
}

// V3Common 封装V3公共报文
type V3Common struct {
	BaseRequest
}

func BuildV3Body(data interface{}) (string, error) {
	envelope := map[string]interface{}{
		"req_time": time.Now().Format("20060102150405"),
		"version":  "3.0",
		"req_data": data,
	}
	return utils.ToJSONString(envelope)
}

func (c *V3Common) BuildBody(data interface{}) (string, error) {
	return BuildV3Body(data)
}

func randomHex(length int) string {
	bytesLen := length / 2
	if length%2 != 0 {
		bytesLen++
	}
	buf := make([]byte, bytesLen)
	if _, err := rand.Read(buf); err != nil {
		return "00000000000000000000000000000000"
	}
	hexStr := hex.EncodeToString(buf)
	if len(hexStr) > length {
		hexStr = hexStr[:length]
	}
	return hexStr
}
