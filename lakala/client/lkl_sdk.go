package client

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
	"zyj.com/golang-study/lakala/util"
)

var (
	lklApiClientMap = make(map[string]*LKLApiClient)
	defaultAppId    string
	mutex           sync.RWMutex
)

const (
	LKL_OP_SDK       = "lkl-go-sdk-1.0.0"
	LKL_OP_FLOWGROUP = "NORMAL"

	SDK_CONNECT_TIMEOUT = 50000
	SDK_READ_TIMEOUT    = 100000
	SDK_SOCKET_TIMEOUT  = 50000
)

// InitWithConfigPath 通过配置文件路径初始化SDK
func InitWithConfigPath(configPath string) error {
	config := Config{}
	properties, err := util.GetProperties(configPath)
	if err != nil {
		return err
	}

	config.AppId = properties["appId"]
	config.SerialNo = properties["serialNo"]
	config.PriKeyPath = properties["priKeyPath"]
	config.LklCerPath = properties["lklCerPath"]
	config.LklNotifyCerPath = properties["lklNotifyCerPath"]

	// 解析超时配置
	if connectTimeout, exists := properties["connectTimeout"]; exists && connectTimeout != "" {
		fmt.Sscanf(connectTimeout, "%d", &config.ConnectTimeout)
	} else {
		config.ConnectTimeout = 5000
	}

	if readTimeout, exists := properties["readTimeout"]; exists && readTimeout != "" {
		fmt.Sscanf(readTimeout, "%d", &config.ReadTimeout)
	} else {
		config.ReadTimeout = 10000
	}

	if socketTimeout, exists := properties["socketTimeout"]; exists && socketTimeout != "" {
		fmt.Sscanf(socketTimeout, "%d", &config.SocketTimeout)
	} else {
		config.SocketTimeout = 5000
	}

	config.ServerUrl = properties["serverUrl"]
	config.Sm4Key = properties["sm4Key"]

	return InitWithConfig(config)
}

// InitWithConfig 通过Config初始化SDK
func InitWithConfig(config *Config) error {

	priKey, err := util.ReadFile(config.PriKeyPath)
	cert, err := util.ReadFile(config.LklCerPath)
	notifyCert, err := util.ReadFile(config.LklNotifyCerPath)
	if err != nil {
		return err
	}
	config2 := &Config2{
		AppId:          config.AppId,
		SerialNo:       config.SerialNo,
		PriKey:         priKey,
		LklCer:         cert,
		LklNotifyCer:   notifyCert,
		ConnectTimeout: config.ConnectTimeout,
		ReadTimeout:    config.ReadTimeout,
		SocketTimeout:  config.SocketTimeout,
		ServerUrl:      config.ServerUrl,
		Sm4Key:         config.Sm4Key,
	}
	return InitWithConfig2(config2)
}

// InitWithConfig2 通过Config2初始化SDK
func InitWithConfig2(config *Config2) error {
	mutex.Lock()
	defer mutex.Unlock()

	// 参数验证
	if err := validateConfig(config); err != nil {
		return err
	}

	// 创建HTTP客户端
	httpClient, err := createHttpClient(config)
	if err != nil {
		return fmt.Errorf("创建HTTP客户端失败: %v", err)
	}

	// 创建SM4工具
	sm4Util, err := createSM4Util(config.Sm4Key)
	if err != nil {
		return fmt.Errorf("创建SM4工具失败: %v", err)
	}

	// 创建通知处理器
	notificationHandler, err := createNotificationHandler(config.LklNotifyCer)
	if err != nil {
		return fmt.Errorf("创建通知处理器失败: %v", err)
	}

	// 创建API客户端
	apiClient := client.NewLKLApiClient(
		config.AppId,
		httpClient,
		notificationHandler,
		sm4Util,
		config.ServerUrl,
	)

	// 保存客户端
	lklApiClientMap[config.AppId] = apiClient

	// 设置默认AppId
	if defaultAppId == "" {
		defaultAppId = config.AppId
	}

	fmt.Printf("SDK初始化成功...appId=%s\n", config.AppId)
	return nil
}

// HttpPost 发送POST请求 - 默认无需加解密
func HttpPost(request *model.LklRequest) (string, error) {
	return HttpPostWithEncryption(request, false, false)
}

// HttpPostWithEncryption 发送带加密的POST请求
func HttpPostWithEncryption(request *model.LklRequest, reqEncrypt, respDecrypt bool) (string, error) {
	appId := request.GetLklAppId()
	if appId == "" {
		appId = defaultAppId
	}

	apiClient, err := getLklApiClient(appId)
	if err != nil {
		return "", err
	}

	needCrypt := reqEncrypt || respDecrypt
	if needCrypt && apiClient.Sm4Util == nil {
		return "", exception.NewSDKExceptionFromEnumsWithInfo(exception.SM4_INIT_FAIL, "appId="+apiClient.AppId)
	}

	url := request.GetFunctionCode().GetUrl()
	var body string

	if reqEncrypt {
		encryptedBody, err := apiClient.Sm4Util.Encrypt(request.ToBody())
		if err != nil {
			return "", err
		}
		body = encryptedBody
	} else {
		body = request.ToBody()
	}

	fullUrl := apiClient.SdkServerUrl + url
	response, err := doHttpPost(fullUrl, body, apiClient)
	if err != nil {
		return "", err
	}

	if respDecrypt {
		return apiClient.Sm4Util.Decrypt(response)
	}

	return response, nil
}

// HttpPostWithUrl 发送POST请求到指定URL
func HttpPostWithUrl(fullUrl, body string) (string, error) {
	apiClient, err := getLklApiClient(defaultAppId)
	if err != nil {
		return "", err
	}
	return doHttpPost(fullUrl, body, apiClient)
}

// HttpPostWithAppId 指定AppId发送POST请求
func HttpPostWithAppId(fullUrl, body, appId string) (string, error) {
	apiClient, err := getLklApiClient(appId)
	if err != nil {
		return "", err
	}
	return doHttpPost(fullUrl, body, apiClient)
}

// NotificationHandle 处理通知回调
func NotificationHandle(request *http.Request) (string, error) {
	return NotificationHandleWithAppId(request, defaultAppId)
}

// NotificationHandleWithBody 处理字符串格式的通知
func NotificationHandleWithBody(body, authorization string) error {
	return NotificationHandleWithBodyAndAppId(body, authorization, defaultAppId)
}

// NotificationHandleWithAppId 指定AppId处理通知回调
func NotificationHandleWithAppId(request *http.Request, appId string) (string, error) {
	apiClient, err := getLklApiClient(appId)
	if err != nil {
		return "", err
	}
	return apiClient.NotificationHandler.Parse(request)
}

// NotificationHandleWithBodyAndAppId 指定AppId处理字符串格式的通知
func NotificationHandleWithBodyAndAppId(body, authorization, appId string) error {
	apiClient, err := getLklApiClient(appId)
	if err != nil {
		return err
	}
	return apiClient.NotificationHandler.Validate(body, authorization)
}

// SM4Encrypt SM4加密
func SM4Encrypt(body string) (string, error) {
	return SM4EncryptWithAppId(body, defaultAppId)
}

// SM4EncryptWithAppId 指定AppId进行SM4加密
func SM4EncryptWithAppId(body, appId string) (string, error) {
	apiClient, err := getLklApiClient(appId)
	if err != nil {
		return "", err
	}
	if apiClient.Sm4Util == nil {
		return "", exception.NewSDKExceptionFromEnumsWithInfo(exception.SM4_INIT_FAIL, "appId="+apiClient.AppId)
	}
	if strings.TrimSpace(body) == "" {
		return "", exception.NewSDKExceptionFromEnums(exception.CHECK_FAIL)
	}
	return apiClient.Sm4Util.Encrypt(body)
}

// SM4Decrypt SM4解密
func SM4Decrypt(body string) (string, error) {
	return SM4DecryptWithAppId(body, defaultAppId)
}

// SM4DecryptWithAppId 指定AppId进行SM4解密
func SM4DecryptWithAppId(body, appId string) (string, error) {
	apiClient, err := getLklApiClient(appId)
	if err != nil {
		return "", err
	}
	if apiClient.Sm4Util == nil {
		return "", exception.NewSDKExceptionFromEnumsWithInfo(exception.SM4_INIT_FAIL, "appId="+apiClient.AppId)
	}
	if strings.TrimSpace(body) == "" {
		return "", exception.NewSDKExceptionFromEnums(exception.CHECK_FAIL)
	}
	return apiClient.Sm4Util.Decrypt(body)
}

// GetDefaultAppId 获取默认AppId
func GetDefaultAppId() string {
	return defaultAppId
}

// SetDefaultAppId 设置默认AppId
func SetDefaultAppId(appId string) {
	mutex.Lock()
	defer mutex.Unlock()
	defaultAppId = appId
}

// 内部工具函数
func validateConfig(config *Config2) error {
	if strings.TrimSpace(config.AppId) == "" {
		return exception.NewSDKExceptionFromEnumsWithInfo(exception.CHECK_FAIL, "拉卡拉appId")
	}
	if strings.TrimSpace(config.SerialNo) == "" {
		return exception.NewSDKExceptionFromEnumsWithInfo(exception.CHECK_FAIL, "商户证书序列号")
	}
	if strings.TrimSpace(config.PriKey) == "" {
		return exception.NewSDKExceptionFromEnumsWithInfo(exception.CHECK_FAIL, "商户私钥信息")
	}
	if strings.TrimSpace(config.LklCer) == "" {
		return exception.NewSDKExceptionFromEnumsWithInfo(exception.CHECK_FAIL, "拉卡拉平台证书")
	}
	if strings.TrimSpace(config.LklNotifyCer) == "" {
		return exception.NewSDKExceptionFromEnumsWithInfo(exception.CHECK_FAIL, "拉卡拉平台通知验签证书")
	}
	if strings.TrimSpace(config.ServerUrl) == "" {
		return exception.NewSDKExceptionFromEnumsWithInfo(exception.CHECK_FAIL, "拉卡拉serverUrl")
	}
	if config.Sm4Key != "" {
		if !util.VerifySM4Key(config.Sm4Key) {
			return exception.NewSDKExceptionFromEnumsWithInfo(exception.CHECK_FAIL, "拉卡拉报文加密对称性密钥")
		}
	}
	return nil
}

func createHttpClient(config *client.Config2) (*http.Client, error) {
	// 加载私钥
	privateKey, err := util.LoadPrivateKeyFromString(config.PriKey)
	if err != nil {
		return nil, err
	}

	// 加载证书
	certificate, err := util.LoadCertificateFromString(config.LklCer)
	if err != nil {
		return nil, err
	}

	builder := client.NewLklHttpClientBuilder().
		WithMerchant(config.AppId, config.SerialNo, privateKey).
		WithLklpay([]*x509.Certificate{certificate})

	// 设置超时
	timeout := time.Duration(SDK_CONNECT_TIMEOUT) * time.Millisecond
	if config.ConnectTimeout > 0 {
		timeout = time.Duration(config.ConnectTimeout) * time.Millisecond
	}
	builder.WithTimeout(timeout)

	return builder.Build(), nil
}

func createSM4Util(sm4Key string) (*util.SM4Util, error) {
	if sm4Key == "" {
		return nil, nil
	}
	key, err := base64.StdEncoding.DecodeString(sm4Key)
	if err != nil {
		return nil, err
	}
	return util.NewSM4Util(key), nil
}

func createNotificationHandler(certStr string) (*notification.NotificationHandler, error) {
	certificate, err := util.LoadCertificateFromString(certStr)
	if err != nil {
		return nil, err
	}
	return notification.NewNotificationHandler(certificate), nil
}

func getLklApiClient(appId string) (*client.LKLApiClient, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	if appId == "" {
		appId = defaultAppId
	}

	if appId == "" {
		return nil, exception.NewSDKExceptionFromEnums(exception.SDK_NOT_INIT)
	}

	apiClient, exists := lklApiClientMap[appId]
	if !exists {
		return nil, exception.NewSDKExceptionFromEnumsWithInfo(exception.SDK_APPID_NOT_INIT, "appId="+appId)
	}

	return apiClient, nil
}

func doHttpPost(url, body string, apiClient *client.LKLApiClient) (string, error) {
	fmt.Printf("httpPost url:%s\n", url)

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return "", exception.NewSDKException("创建请求失败", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	resp, err := apiClient.HttpClient.Do(req)
	if err != nil {
		return "", exception.NewSDKException("请求执行失败", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", exception.NewSDKExceptionWithCode("HTTP_ERROR",
			fmt.Sprintf("请求响应失败：%s", string(respBody)))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", exception.NewSDKException("读取响应失败", err)
	}

	return string(respBody), nil
}
