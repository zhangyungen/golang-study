package client

// UrlMaping 函数代码结构体
type UrlMaping struct {
	Code string
	Name string
}

// GetCode 获取代码
func (f *UrlMaping) GetCode() string {
	return f.Code
}

// GetName 获取名称
func (f *UrlMaping) GetName() string {
	return f.Name
}

// GetUrl 获取URL
func (f *UrlMaping) GetUrl() string {
	return f.Code
}

// 定义所有函数代码常量
var (
	// V3接口
	API_V3_LABS_QUERY_TRADEQUERY                   = &UrlMaping{"/api/v3/labs/query/tradequery", "聚合扫码-交易查询"}
	API_V3_LABS_TRANS_PREORDER                     = &UrlMaping{"/api/v3/labs/trans/preorder", "聚合扫码-聚合主扫"}
	API_V3_LABS_TRANS_MICROPAY                     = &UrlMaping{"/api/v3/labs/trans/micropay", "聚合扫码-聚合被扫"}
	API_V3_LABS_RELATION_REFUND                    = &UrlMaping{"/api/v3/labs/relation/refund", "聚合扫码-退款交易"}
	API_V3_SACS_SEPARATE                           = &UrlMaping{"/api/v3/sacs/separate", "分账指令V3接口-订单分账"}
	API_V3_SACS_QUERYAMT                           = &UrlMaping{"/api/v3/sacs/queryAmt", "分账指令V3接口-可分账金额查询"}
	API_V3_SACS_QUERY                              = &UrlMaping{"/api/v3/sacs/query", "分账指令V3接口-分账结果查询"}
	API_V3_LABS_QUERY_IDMREFUNDQUERY               = &UrlMaping{"/api/v3/labs/query/idmrefundquery", "商户订单退款查询"}
	API_V3_LABS_RELATION_IDMREFUND                 = &UrlMaping{"/api/v3/labs/relation/idmrefund", "商户订单退款交易"}
	API_V3_LABS_RELATION_CLOSE                     = &UrlMaping{"/api/v3/labs/relation/close", "聚合扫码-关单"}
	API_V3_LABS_RELATION_REVOKED                   = &UrlMaping{"/api/v3/labs/relation/revoked", "扫码-撤销"}
	API_V3_CCSS_COUNTER_ORDER_SPECIAL_CREATE       = &UrlMaping{"/api/v3/ccss/counter/order/special_create", "收银台订单创建"}
	API_V3_CCSS_COUNTER_ORDER_QUERY                = &UrlMaping{"/api/v3/ccss/counter/order/query", "收银台订单查询"}
	API_V3_SACS_CANCEL                             = &UrlMaping{"/api/v3/sacs/cancel", "订单分账撤销"}
	API_V3_SACS_FALLBACK                           = &UrlMaping{"/api/v3/sacs/fallback", "订单分账回退"}
	API_V3_SACS_BALANCESEPARATE                    = &UrlMaping{"/api/v3/sacs/balanceSeparate", "余额分账"}
	API_V3_SACS_BALANCECANCEL                      = &UrlMaping{"/api/v3/sacs/balanceCancel", "余额分账撤销"}
	API_V3_SACS_BALANCEFALLBACK                    = &UrlMaping{"/api/v3/sacs/balanceFallback", "余额分账回退"}
	API_V3_SACS_BALANCESEPARATEQUERY               = &UrlMaping{"/api/v3/sacs/balanceSeparateQuery", "余额分账结果查询"}
	API_V3_LABS_TRANS_SHARE_CODE                   = &UrlMaping{"/api/v3/labs/trans/share_code", "支付类接口-申请分享码(支付宝吱口令)"}
	API_V3_LABS_TRANS_PREORDER_ENCRY               = &UrlMaping{"/api/v3/labs/trans/preorder_encry", "支付类接口-主扫交易接口（全报文加密）"}
	API_V3_LABS_TRANS_MICROPAY_ENCRY               = &UrlMaping{"/api/v3/labs/trans/micropay_encry", "支付类接口-被扫接口（全报文加密）"}
	API_V3_CCSS_COUNTER_ORDER_SPECIAL_CREATE_ENCRY = &UrlMaping{"/api/v3/ccss/counter/order/special_create_encry", "收银台服务系统-收银台订单创建 （全报文加密）"}
	API_V3_LABS_QUERY_GETFACEAUTHINFO              = &UrlMaping{"/api/v3/labs/query/getfaceauthinfo", "支付类接口-微信刷脸授权信息获取"}
	API_V3_LAMS_MERCHANT_QUERY_TRANS               = &UrlMaping{"/api/v3/lams/merchant/query_trans", "商户服务接口V3.0-商户计费查询"}
	API_V3_LAMS_TRADE_TRADE_REFUND                 = &UrlMaping{"/api/v3/lams/trade/trade_refund", "商户服务接口V3.0-统一退货"}
	API_V3_LAMS_TRADE_TRADE_REFUND_QUERY           = &UrlMaping{"/api/v3/lams/trade/trade_refund_query", "商户服务接口V3.0-退货查询"}

	API_V3_IPSQP_TRANS_QUICKSIGNAPPLY    = &UrlMaping{"/api/v3/ipsqp/trans/quickSignApply", "快捷签约申请(发送短信)"}
	API_V3_IPSQP_TRANS_QUICKSIGNCONFIRM  = &UrlMaping{"/api/v3/ipsqp/trans/quickSignConfirm", "快捷签约确认"}
	API_V3_IPSQP_TRANS_QUICKSIGNCANCEL   = &UrlMaping{"/api/v3/ipsqp/trans/quickSignCancel", "快捷解约"}
	API_V3_IPSQP_TRANS_APPLYQUICKPAY     = &UrlMaping{"/api/v3/ipsqp/trans/applyQuickPay", "快捷支付申请（已有协议号）"}
	API_V3_IPSQP_TRANS_CONFIRMQUICKPAY   = &UrlMaping{"/api/v3/ipsqp/trans/confirmQuickPay", "快捷支付确认接口（已有协议号）"}
	API_V3_IPSQP_TRANS_APPLYSIGNANDPAY   = &UrlMaping{"/api/v3/ipsqp/trans/applySignAndPay", "快捷签约并支付申请（没有协议号）"}
	API_V3_IPSQP_TRANS_CONFIRMSIGNANDPAY = &UrlMaping{"/api/v3/ipsqp/trans/confirmSignAndPay", "快捷支付确认接口（没有协议号）"}
	API_V3_IPSQP_QUERY_TRADEQUERY        = &UrlMaping{"/api/v3/ipsqp/query/tradeQuery", "代收查询"}
	API_V3_IPSQP_TRANS_GATEWAYSIGN       = &UrlMaping{"/api/v3/ipsqp/trans/gatewaySign", "网关签约（一键绑卡）"}
	API_V3_IPSQP_TRANS_DIRECTLYQUICKPAY  = &UrlMaping{"/api/v3/ipsqp/trans/directlyQuickPay", "快捷直接支付"}
	API_V3_IPSQP_QUERY_QUERYSIGNBYPAN    = &UrlMaping{"/api/v3/ipsqp/query/querySignByPan", "根据卡号查询协议号"}

	API_V3_IPSDK_TRANS_COLLECTIONSIGNAPPLY   = &UrlMaping{"/api/v3/ipsdk/trans/collectionSignApply", "代收签约申请(发送短信)"}
	API_V3_IPSDK_TRANS_COLLECTIONSIGNCONFIRM = &UrlMaping{"/api/v3/ipsdk/trans/collectionSignConfirm", "代收签约确认"}
	API_V3_IPSDK_TRANS_COLLECTIONSIGNCANCEL  = &UrlMaping{"/api/v3/ipsdk/trans/collectionSignCancel", "代收解约"}
	API_V3_IPSDK_TRANS_BATCHCOLLECTION       = &UrlMaping{"/api/v3/ipsdk/trans/batchCollection", "批量代收"}
	API_V3_IPSDK_TRANS_BATCHQUERY            = &UrlMaping{"/api/v3/ipsdk/trans/batchQuery", "批量代收查询"}
	API_V3_IPSDK_TRANS_REALTIMECOLLECTION    = &UrlMaping{"/api/v3/ipsdk/trans/realTimeCollection", "实时代收"}
	API_V3_IPSDK_RELATION_REFUND             = &UrlMaping{"/api/v3/ipsdk/relation/refund", "代收退款"}

	API_V3_IPSDF_PAID_PAY        = &UrlMaping{"/api/v3/ipsdf/paid/pay", "实时代付接口"}
	API_V3_IPSDF_PAID_QUERY      = &UrlMaping{"/api/v3/ipsdf/paid/query", "实时代付查询接口"}
	API_V3_IPSDF_PAID_BATCHPAY   = &UrlMaping{"/api/v3/ipsdf/paid/batchPay", "批量代付接口"}
	API_V3_IPSDF_PAID_BATCHQUERY = &UrlMaping{"/api/v3/ipsdf/paid/batchQuery", "批量代付查询接口"}

	API_V3_CCSS_COUNTER_ORDER_CLOSE = &UrlMaping{"/api/v3/ccss/counter/order/close", "收银台关单"}

	API_V3_RFD_REFUND_FRONT_REFUND       = &UrlMaping{"/api/v3/rfd/refund_front/refund", "退货前置-退货(同步)"}
	API_V3_RFD_REFUND_FRONT_REFUND_QUERY = &UrlMaping{"/api/v3/rfd/refund_front/refund_query", "退货前置-退货查询"}
	API_V3_RFD_REFUND_FRONT_REFUND_FEE   = &UrlMaping{"/api/v3/rfd/refund_front/refund_fee", "退货前置-退货手续费试算"}
	API_V3_RFD_REFUND_FRONT_MERGE_REFUND = &UrlMaping{"/api/v3/rfd/refund_front/merge_refund", "退货前置-合单退货"}

	// V2接口
	API_V2_MMS_OPENAPI_LEDGER_APPLYLEDGERMER       = &UrlMaping{"/api/v2/mms/openApi/ledger/applyLedgerMer", "商户分账业务开通申请"}
	API_V2_MMS_OPENAPI_LEDGER_MODIFYLEDGERMER      = &UrlMaping{"/api/v2/mms/openApi/ledger/modifyLedgerMer", "商户分账信息变更申请"}
	API_V2_MMS_OPENAPI_LEDGER_QUERYLEDGERMER       = &UrlMaping{"/api/v2/mms/openApi/ledger/queryLedgerMer", "分账商户信息查询"}
	API_V2_MMS_OPENAPI_LEDGER_APPLYLEDGERRECEIVER  = &UrlMaping{"/api/v2/mms/openApi/ledger/applyLedgerReceiver", "分账接收方开通申请"}
	API_V2_MMS_OPENAPI_LEDGER_MODIFYLEDGERRECEIVER = &UrlMaping{"/api/v2/mms/openApi/ledger/modifyLedgerReceiver", "分账接收方信息变更申请"}
	API_V2_MMS_OPENAPI_LEDGER_APPLYBIND            = &UrlMaping{"/api/v2/mms/openApi/ledger/applyBind", "分账关系绑定申请"}
	API_V2_MMS_OPENAPI_LEDGER_APPLYUNBIND          = &UrlMaping{"/api/v2/mms/openApi/ledger/applyUnBind", "分账关系解绑申请"}
	API_V2_MMS_OPENAPI_ACTIVESETTLE_APPLY          = &UrlMaping{"/api/v2/mms/openApi/activeSettle/apply", "商户开通主动结算业务"}

	API_V2_LAEP_INDUSTRY_EWALLETWITHDRAW       = &UrlMaping{"/api/v2/laep/industry/ewalletWithdraw", "账管家V2.0-账管家提现"}
	API_V2_LAEP_INDUSTRY_EWALLETWITHDRAWQUERY  = &UrlMaping{"/api/v2/laep/industry/ewalletWithdrawQuery", "账管家V2.0-账管家提现结果查询"}
	API_V2_LAEP_INDUSTRY_EWALLET_SETTLEPROFILE = &UrlMaping{"/api/v2/laep/industry/ewallet/settleProfile", "账管家V2.0-账管家提款模式设置"}
	API_V2_LAEP_INDUSTRY_EWALLET_SETTLEQUERY   = &UrlMaping{"/api/v2/laep/industry/ewallet/settleQuery", "账管家V2.0-账管家提款模式查询"}
	API_V2_LAEP_INDUSTRY_EWALLET_EWALLETFEE    = &UrlMaping{"/api/v2/laep/industry/ewallet/ewalletfee", "账管家V2.0-账管家提款手续费设置"}
	API_V2_LAEP_INDUSTRY_EWALLET_QUERYFEE      = &UrlMaping{"/api/v2/laep/industry/ewallet/queryfee", "账管家V2.0-账管家提款手续费查询"}
	API_V2_LAEP_INDUSTRY_EWALLETWITHDRAWD1     = &UrlMaping{"/api/v2/laep/industry/ewalletWithdrawD1", "账管家V2.0-账管家提现D1"}
	API_V2_LAEP_INDUSTRY_TRANSFER_DIRECT       = &UrlMaping{"/api/v2/laep/industry/transfer/direct", "账管家V2.0-可用余额定向转账"}
	API_V2_LAEP_INDUSTRY_TRANSFER_QUERY        = &UrlMaping{"/api/v2/laep/industry/transfer/query", "账管家V2.0-转账订单查询"}
	API_V2_LAEP_INDUSTRY_EWALLETBALANCEQUERY   = &UrlMaping{"/api/v2/laep/industry/ewalletBalanceQuery", "账管家V2.0-账管家余额查询"}
	API_V2_LAEP_EWALLETACCOUNT_QUERYACCTINFO   = &UrlMaping{"/api/v2/laep/ewalletAccount/queryAcctInfo", "账管家V2.0-收单账户信息查询"}
	API_V2_MRSSQUERY_QUERYACCTINFO             = &UrlMaping{"/api/v2/mrssQuery/queryAcctInfo", "收单清结算接口V2-账户信息查询"}
	API_V2_LAEP_INDUSTRY_BANKCARD_BIND         = &UrlMaping{"/api/v2/laep/industry/bankcard/bind", "账管家V2.0-绑定银行卡"}
	API_V2_LAEP_INDUSTRY_BANKCARD_LIST         = &UrlMaping{"/api/v2/laep/industry/bankcard/list", "账管家V2.0-查询银行卡列表"}
	API_V2_LAEP_INDUSTRY_BANKCARD_UNBIND       = &UrlMaping{"/api/v2/laep/industry/bankcard/unbind", "账管家V2.0-解绑银行卡"}
	API_V2_LAEP_EWALLETACCOUNT_BILLSQRY        = &UrlMaping{"/api/v2/laep/ewalletAccount/billsQry", "账管家V2.0-账单分页查询"}
	API_V2_LAEP_CREATEELERECEIPT               = &UrlMaping{"/api/v2/laep/createEleReceipt", "账管家V2.0-生成电子回单"}

	API_V2_MCQS_MERCHANT_QUERYMERCHANTDETAIL = &UrlMaping{"/api/v2/mcqs/merchant/queryMerchantDetail", "商户信息查询"}
	API_V2_MCQS_LIMIT_QUERYLIMIT             = &UrlMaping{"/api/v2/mcqs/limit/queryLimit", "商户限额查询"}
	API_V2_MCQS_MERCHANT_QUERYMERACCOUNT     = &UrlMaping{"/api/v2/mcqs/merchant/queryMerAccount", "商户结算账户查询"}
	API_V2_MCQS_MERCHANT_QUERYTERMACCOUNT    = &UrlMaping{"/api/v2/mcqs/merchant/queryTermAccount", "终端结算账户查询"}
	API_V2_MCQS_BUSI_QUERYTERMTRANTYPE       = &UrlMaping{"/api/v2/mcqs/busi/queryTermTranType", "终端交易权限查询"}
	API_V2_MCQS_MERCHANT_QUERYSHOPLISTBYNUM  = &UrlMaping{"/api/v2/mcqs/merchant/queryShopListByNum", "商户网点列表查询"}
	API_V2_MCQS_BUSI_QUERYBUSILISTBYNUM      = &UrlMaping{"/api/v2/mcqs/busi/queryBusiListByNum", "商户终端列表查询"}
	API_V2_MCQS_BUSI_QUERYTERMFEE            = &UrlMaping{"/api/v2/mcqs/busi/queryTermFee", "终端费率查询"}
	API_V2_MMS_OPENAPI_UPLOADFILE            = &UrlMaping{"/api/v2/mms/openApi/uploadFile", "附件上传"}
)

// GetUrlName 根据URL获取函数代码
func GetUrlName(url string) *UrlMaping {
	// 这里可以通过反射或其他方式实现URL到FunctionCode的映射
	// 简化实现，实际使用时可以根据需要完善
	switch url {
	case "/api/v3/labs/query/tradequery":
		return API_V3_LABS_QUERY_TRADEQUERY
	case "/api/v3/labs/trans/preorder":
		return API_V3_LABS_TRANS_PREORDER
	// 添加其他URL映射...
	default:
		return &UrlMaping{Code: url, Name: "未知接口"}
	}
}
