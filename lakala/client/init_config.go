package client

// Config 配置文件格式配置证书
type Config struct {
	AppId            string `json:"appId"`            // 拉卡拉appId
	SerialNo         string `json:"serialNo"`         // 商户证书序列号
	PriKeyPath       string `json:"priKeyPath"`       // 商户私钥信息地址
	LklCerPath       string `json:"lklCerPath"`       // 拉卡拉支付平台证书地址
	LklNotifyCerPath string `json:"lklNotifyCerPath"` // 拉卡拉通知推送验签证书地址
	Sm4Key           string `json:"sm4Key"`           // 拉卡拉报文加密对称性密钥
	ConnectTimeout   int    `json:"connectTimeout"`   // 连接超时时间
	ReadTimeout      int    `json:"readTimeout"`      // 建立连接后超时时间
	SocketTimeout    int    `json:"socketTimeout"`    // 读取数据超时时间
	ServerUrl        string `json:"serverUrl"`        // API服务地址
}

// Config2 支持证书以字符串格式传入
type Config2 struct {
	AppId          string `json:"appId"`          // 拉卡拉appId
	SerialNo       string `json:"serialNo"`       // 商户证书序列号
	PriKey         string `json:"priKey"`         // 商户私钥信息
	LklCer         string `json:"lklCer"`         // 拉卡拉支付平台证书
	LklNotifyCer   string `json:"lklNotifyCer"`   // 拉卡拉通知推送验签证书
	Sm4Key         string `json:"sm4Key"`         // 拉卡拉报文加密对称性密钥
	ConnectTimeout int    `json:"connectTimeout"` // 连接超时时间
	ReadTimeout    int    `json:"readTimeout"`    // 建立连接后超时时间
	SocketTimeout  int    `json:"socketTimeout"`  // 读取数据超时时间
	ServerUrl      string `json:"serverUrl"`      // API服务地址
}

// NewConfig 创建Config
func NewConfig(appId, serialNo, priKeyPath, lklCerPath, lklNotifyCerPath, serverUrl string) *Config {
	return &Config{
		AppId:            appId,
		SerialNo:         serialNo,
		PriKeyPath:       priKeyPath,
		LklCerPath:       lklCerPath,
		LklNotifyCerPath: lklNotifyCerPath,
		ServerUrl:        serverUrl,
	}
}

// NewConfig2 创建Config2
func NewConfig2(appId, serialNo, priKey, lklCer, lklNotifyCer, serverUrl string) *Config2 {
	return &Config2{
		AppId:        appId,
		SerialNo:     serialNo,
		PriKey:       priKey,
		LklCer:       lklCer,
		LklNotifyCer: lklNotifyCer,
		ServerUrl:    serverUrl,
	}
}
