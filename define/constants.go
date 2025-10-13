package define

const (
	RESPONSE_ERROR_CODE = "code"
	RESPONSE_DATA       = "data"

	APP_VERSION_KEY = "App-Version"
	AppPlatform     = "Platform"

	HEADER_TRACE_ID_KEY = "X-Request-Id"

	App_ID_KEY = "app_id"

	LOG_DEBUG_MSG             = "debug_msg"
	RESPONSE_ERROR_DETAIL_MSG = "error"
	LOG_PANIC_STACK           = "panic_stack"
	RESPONSE_ERROR_STACK      = "err_stack"

	RESPONSE_META = "meta"

	RESPONSE_HEADER_ERROR_CODE = "Code"

	RESPONSE_ERROR_MESSAGE = "message"

	En = "en"

	CheeseIOSUpdateAppType     = 107
	CheeseAndroidUpdateAppType = 108
	StableRelease              = 1
	BetaRelease                = 2 //官网
	ForceUpdate                = 3

	CacheTokenPrefix = "token:"
	ASASourceID      = 1329
	Ios              = "ios"
	Android          = "android"
)

// Response Code
const (
	SUCCESS = 0

	//TOKEN ERROR 101-200
	TOKEN_EMPTY          = 101
	TOKEN_PARSE_ERROR    = 102
	TOKEN_NOT_VALID      = 103
	TOKEN_INFO_NOT_VALID = 104
	TOKEN_USER_ERROR     = 105

	//AUTH ERROR 201-300
	PROJECT_NONE    = 201
	PROJECT_NO_AUTH = 202

	//other
	PASSWORD_ERROR  = 1001
	NO_ENOUGH_SPACE = 1002
	RETHOUCH_ERROR  = 1100
)

const (
	ERROR_PANIC = 44
)

const (
	DaySecond = 86400
)

const (
	ReportTypeName        = "1_联机拍摄问题"
	ReportTypeDefaultName = "默认问题"
)

const (
	DefaultTemplateID   = 1000001
	DefaultTemplateName = "基础白色"

	DefaultSPBottomInfo = "已经到底部啦"
)

const (
	AppBannerPublished = 1
)

const (
	TokenErrorMessage = "用户信息过期，请重新登录"
)

const (
	EntryClick            = 0
	EntryCountdown        = 1
	DefaultEntryText      = "进入相册直播"
	DefaultEntryCountdown = 3
)

// Language constants for multilingual project settings
const (
	LanguageSimplifiedChinese  = 1 // 简体中文
	LanguageTraditionalChinese = 2 // 繁体中文
	LanguageEnglish            = 3 // English
)

// Language switch constants
const (
	LanguageSwitchDisabled = 0 // 关闭
	LanguageSwitchEnabled  = 1 // 开启
)

type CheeseRefundReason int

const (
	CheeseRefundReasonFreeChange CheeseRefundReason = 1

	CheeseRefundReasonFreeChangeStr = "因多位选片人同时付款，或您修改了套餐设置，为避免超发免费张数，系统已发起自动退款。请引导选片人刷新小程序后重新下单"
)

const (
	CheeseTradeStateWaitPay        = 0
	CheeseTradeStateSuccess        = 1 //支付成功
	CheeseTradeStateRefunding      = 4 //退款中
	CheeseTradeStateRefunded       = 5 //已退款
	CheeseTradeStateRefundAbnormal = 6 //退款异常
	CheeseTradeStateRefundClosed   = 7 //退款异常超时
)

const (
	CheeseTradeStateWaitPayStr        = "待支付"
	CheeseTradeStateSuccessStr        = "支付成功"   //支付成功
	CheeseTradeStateRefundingStr      = "自动退款中"  //退款中
	CheeseTradeStateRefundedStr       = "已退款"    //已退款
	CheeseTradeStateRefundAbnormalStr = "退款异常"   //退款异常
	CheeseTradeStateRefundClosedStr   = "退款异常超时" //退款异常超时
)

var CheeseTradeStateMap = map[int]string{
	CheeseTradeStateWaitPay:        CheeseTradeStateWaitPayStr,
	CheeseTradeStateSuccess:        CheeseTradeStateSuccessStr,
	CheeseTradeStateRefunding:      CheeseTradeStateRefundingStr,
	CheeseTradeStateRefunded:       CheeseTradeStateRefundedStr,
	CheeseTradeStateRefundAbnormal: CheeseTradeStateRefundAbnormalStr,
	CheeseTradeStateRefundClosed:   CheeseTradeStateRefundClosedStr,
}

// SMS Reminder Types 短信提醒类型
type SmsReminderType int

const (
	SmsReminderTypeEnterpriseApp SmsReminderType = 1 // 企业申请进件通知
	SendEnterpriseTodoReminder   SmsReminderType = 2 // 项目商家管理员提醒

)

const (
	SmsReminderTypeProjectAdminStr  = "项目商家管理员提醒"
	SmsReminderTypeEnterpriseAppStr = "企业申请进件通知"
)

var SmsReminderTypeMap = map[SmsReminderType]string{
	SendEnterpriseTodoReminder:   SmsReminderTypeProjectAdminStr,
	SmsReminderTypeEnterpriseApp: SmsReminderTypeEnterpriseAppStr,
}

// SMS提醒缓存key前缀
const (
	SmsReminderCacheKeyPrefix = "sms_reminder:"
)
