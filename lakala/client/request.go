package client

import (
	"time"
	"zyj.com/golang-study/lakala/util"
	"zyj.com/golang-study/util/strutil"
)

// LklRequest 请求接口
type LklRequest interface {
	// GetFunctionCode 获取功能代码，对应接口URL
	GetFunctionCode() *UrlMaping

	// ToBody 转换为请求报文内容
	ToBody() string

	// GetLklAppId 获取lkl-app-id
	GetLklAppId() string
}

// BaseRequest 基础请求结构体
type BaseRequest struct {
	LklAppId string `json:"-"`
}

// GetLklAppId 获取LKL AppId
func (b *BaseRequest) GetLklAppId() string {
	return b.LklAppId
}

// SetLklAppId 设置LKL AppId
func (b *BaseRequest) SetLklAppId(appId string) {
	b.LklAppId = appId
}

// V2CommRequest V2接口通用公共参数
type V2CommRequest struct {
	BaseRequest
}

// V2RequestData V2请求数据结构
type V2RequestData struct {
	ReqId   string      `json:"reqId"`
	ReqTime string      `json:"reqTime"`
	Version string      `json:"version"`
	ReqData interface{} `json:"reqData"`
}

// NewV2CommRequest 创建V2通用请求
func NewV2CommRequest() *V2CommRequest {
	return &V2CommRequest{}
}

// ToBody 转换为请求报文内容
func (v *V2CommRequest) ToBody() string {
	// 生成UUID作为请求ID
	reqId := generateUUID()

	// 格式化当前时间
	reqTime := time.Now().Format("20060102150405")

	// 构建请求数据
	requestData := V2RequestData{
		ReqId:   reqId,
		ReqTime: reqTime,
		Version: "2.0",
		ReqData: v, // 注意：这里需要子类重写GetFunctionCode方法
	}

	return util.ToJSONString(requestData)
}

// GetFunctionCode 获取功能代码 - 需要子类重写
func (v *V2CommRequest) GetFunctionCode() *UrlMaping {
	// 默认实现，子类应该重写这个方法
	return nil
}

// generateUUID 生成UUID
func generateUUID() string {
	// 简单的UUID生成实现
	// 在实际项目中可以使用github.com/google/uuid等库'
	return strutil.GenerateUUIDV4()
}
