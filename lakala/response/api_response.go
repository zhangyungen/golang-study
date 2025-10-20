package response

import (
	"encoding/json"
)

/**
* API返回结果V3 结构体
 */
type ApiResponseV3[T any] struct {
	Msg      string `json:"msg"`
	RespTime string `json:"resp_time"`
	Code     string `json:"code"`
	RespData T      `json:"resp_data"`
}

/**
* API返回结果V2 结构体
 */
type ApiResponseV2[T any] struct {
	RetCode  string `json:"retCode"`
	RetMsg   string `json:"retMsg"`
	RespData T      `json:"respData"`
}

/**
* API返回结果V3 聚合主扫-聚合交易退款返回结构体
 */
type LabsTransOrderRefundResponse struct {
	OutTradeNo       string `json:"out_trade_no"`
	TradeNo          string `json:"trade_no"`
	LogNo            string `json:"log_no"`
	AccTradeNo       string `json:"acc_trade_no"`
	AccountType      string `json:"account_type"`
	TotalAmount      string `json:"total_amount"`
	RefundAmount     string `json:"refund_amount"`
	PayerAmount      string `json:"payer_amount"`
	TradeTime        string `json:"trade_time"`
	OriginTradeNo    string `json:"origin_trade_no"`
	OriginOutTradeNo string `json:"origin_out_trade_no"`
	UpIssAddnData    string `json:"up_iss_addn_data"`
	UpCouponInfo     string `json:"up_coupon_info"`
	TradeInfo        string `json:"trade_info"`
}

/*
*
统一退货接口
*/
type TransOrderRefundResponse struct {
	TradeState       string `json:"trade_state"`
	RefundType       string `json:"refund_type"`
	MerchantNo       string `json:"merchant_no"`
	OutTradeNo       string `json:"out_trade_no"`
	TradeNo          string `json:"trade_no"`
	LogNo            string `json:"log_no"`
	AccountType      string `json:"account_type"`
	TotalAmount      string `json:"total_amount"`
	RefundAmount     string `json:"refund_amount"`
	PayerAmount      string `json:"payer_amount"`
	TradeTime        string `json:"trade_time"`
	OriginTradeNo    string `json:"origin_trade_no"`
	OriginOutTradeNo string `json:"origin_out_trade_no"`
	OriginLogNo      string `json:"origin_log_no"`
	ChannelRetDesc   string `json:"channel_ret_desc"`
}

/**
* API返回结果V3 聚合主扫-聚合扫码-交易查询返回结构体
 */
type LabsTransOrderQueryResponse struct {
	OutTradeNo string `json:"out_trade_no"`
	TradeNo    string `json:"trade_no"`
	LogNo      string `json:"log_no"`
	SplitAttr  string `json:"split_attr"`
	SplitInfo  []struct {
		SubTradeNo    string `json:"sub_trade_no"`
		SubLogNo      string `json:"sub_log_no"`
		OutSubTradeNo string `json:"out_sub_trade_no"`
		MerchantNo    string `json:"merchant_no"`
		TermNo        string `json:"term_no"`
		Amount        string `json:"amount"`
	} `json:"split_info"`
	AccTradeNo         string `json:"acc_trade_no"`
	AccountType        string `json:"account_type"`
	SettleMerchantNo   string `json:"settle_merchant_no"`
	SettleTermNo       string `json:"settle_term_no"`
	TradeState         string `json:"trade_state"`
	TradeStateDesc     string `json:"trade_state_desc"`
	TotalAmount        string `json:"total_amount"`
	PayerAmount        string `json:"payer_amount"`
	AccSettleAmount    string `json:"acc_settle_amount"`
	AccMdiscountAmount string `json:"acc_mdiscount_amount"`
	AccDiscountAmount  string `json:"acc_discount_amount"`
	TradeTime          string `json:"trade_time"`
	UserId1            string `json:"user_id1"`
	UserId2            string `json:"user_id2"`
	BankType           string `json:"bank_type"`
	AccActivityId      string `json:"acc_activity_id"`
	UpCouponInfo       string `json:"up_coupon_info"`
	TradeInfo          string `json:"trade_info"`
}

/**
* API返回结果V3 聚合主扫-交易返回结构体
 */
type LabsTransPreorderResponse struct {
	MerchantNo       string `json:"merchant_no"`
	OutTradeNo       string `json:"out_trade_no"`
	TradeNo          string `json:"trade_no"`
	LogNo            string `json:"log_no"`
	SettleMerchantNo string `json:"settle_merchant_no"`
	SettleTermNo     string `json:"settle_term_no"`
	AccRespFields    struct {
		Code        string `json:"code"`
		CodeImage   string `json:"code_image"`
		PrepayId    string `json:"prepay_id"`
		AppId       string `json:"app_id"`
		PaySign     string `json:"pay_sign"`
		TimeStamp   string `json:"time_stamp"`
		NonceStr    string `json:"nonce_str"`
		Package     string `json:"package"`
		SignType    string `json:"sign_type"`
		FormData    string `json:"form_data"`
		RedirectUrl string `json:"redirect_url"`
		BestPayInfo string `json:"best_pay_info"`
	} `json:"acc_resp_fields"`
}

/**
* API返回结果V2 商户分账-商户信息查询
 */
type LedgerMerResponse struct {
	MerInnerNo       string  `json:"merInnerNo"`
	MerCupNo         string  `json:"merCupNo"`
	SplitLowestRatio float64 `json:"splitLowestRatio"`
	OrgId            string  `json:"orgId"`
	OrgName          string  `json:"orgName"`
	SplitStatus      string  `json:"splitStatus"`
	SplitStatusText  string  `json:"splitStatusText"`
	SplitRange       string  `json:"splitRange"`
	SepFundSource    string  `json:"sepFundSource"`
	BindRelations    []struct {
		MerInnerNo   string `json:"merInnerNo"`
		MerCupNo     string `json:"merCupNo"`
		ReceiverNo   string `json:"receiverNo"`
		ReceiverName string `json:"receiverName"`
	} `json:"bindRelations"`
}

/**
* API返回结果V2 分账信息接收方创建变更返回结构体·
 */
type LedgerReceiverBindResponse struct {
	Version    string `json:"version"`
	OrderNo    string `json:"orderNo"`
	OrgCode    string `json:"orgCode"`
	OrgId      string `json:"orgId"`
	OrgName    string `json:"orgName"`
	ReceiverNo string `json:"receiverNo"`
}

/**
* API返回结果V3 可分账金额查询返回结构体
 */
type V3SacsQueryAmtResponse struct {
	MerchantNo       string `json:"merchant_no"`
	TotalSeparateAmt string `json:"total_separate_amt"`
	CanSeparateAmt   string `json:"can_separate_amt"`
	LogDate          string `json:"log_date"`
	LogNo            string `json:"log_no"`
}

type EwalletBalanceResponse struct {
	PayNo        string `json:"payNo"`
	PayType      string `json:"payType"`
	AcctSt       string `json:"acctSt"`
	ForceBalance string `json:"forceBalance"`
	HisBalance   string `json:"hisBalance"`
	ReBalance    string `json:"reBalance"`
	CurBalance   string `json:"curBalance"`
}

/**
* API返回结果V3 分账结果查询
 */
type SacsBalanceSeparateQueryResponse struct {
	TotalAmt    string `json:"total_amt"`
	FinishDate  string `json:"finish_date"`
	LogDate     string `json:"log_date"`
	CalType     string `json:"cal_type"`
	DetailDatas []struct {
		RecvNo string `json:"recv_no"`
		Amt    int    `json:"amt"`
	} `json:"detail_datas"`
	SeparateNo    string `json:"separate_no"`
	LogNo         string `json:"log_no"`
	SeparateType  string `json:"separate_type"`
	OutSeparateNo string `json:"out_separate_no"`
	CmdType       string `json:"cmd_type"`
	SeparateDate  string `json:"separate_date"`
	Status        string `json:"status"`
}

/**
* API返回结果V3 分账撤销
 */
type SacsSeparateCancelResponse struct {
	TotalAmt         string `json:"total_amt"`
	SeparateNo       string `json:"separate_no"`
	OriginSeparateNo string `json:"origin_separate_no"`
	OutSeparateNo    string `json:"out_separate_no"`
	Status           string `json:"status"`
}

/**
* API返回结果V3 分账回退
 */
type SacsSeparateFallBackResponse struct {
	TotalAmt         string `json:"total_amt"`
	SeparateNo       string `json:"separate_no"`
	OriginSeparateNo string `json:"origin_separate_no"`
	OutSeparateNo    string `json:"out_separate_no"`
	Status           string `json:"status"`
}

/**
* API返回结果V3 分账返回参数
 */
type SacsBalanceSeparateResponse struct {
	TotalAmt      string `json:"total_amt"`
	LogDate       string `json:"log_date"`
	SeparateNo    string `json:"separate_no"`
	LogNo         string `json:"log_no"`
	OutSeparateNo string `json:"out_separate_no"`
	Status        string `json:"status"`
}

/**
* API返回结果V2 分账关系绑定和解除绑定,变更返回结构体
 */
type LedgerApplyBindResponse struct {
	Version string `json:"version"`
	OrderNo string `json:"orderNo"`
	OrgCode string `json:"orgCode"`
	ApplyId int64  `json:"applyId"`
}

/*
* API返回结果V2 分账关系变更返回结构体
 */
type ModifyLedgerMerResponse struct {
	Version string `json:"version"`
	OrderNo string `json:"orderNo"`
	OrgCode string `json:"orgCode"`
	ApplyId int64  `json:"applyId"`
}

/**
* API返回结果V2 分账接受者详情返回结构体
 */
type LedgerReceiverDetailResponse struct {
	RowStatus                  interface{} `json:"rowStatus"`
	RowSno                     string      `json:"rowSno"`
	RowCreateUser              interface{} `json:"rowCreateUser"`
	RowCreateUserName          string      `json:"rowCreateUserName"`
	RowCreateTm                string      `json:"rowCreateTm"`
	RowModifyUser              interface{} `json:"rowModifyUser"`
	RowModifyUserName          interface{} `json:"rowModifyUserName"`
	RowModifyTm                string      `json:"rowModifyTm"`
	RowVerNo                   interface{} `json:"rowVerNo"`
	Id                         int         `json:"id"`
	ReceiverNo                 string      `json:"receiverNo"`
	ReceiverName               string      `json:"receiverName"`
	ContactMobile              string      `json:"contactMobile"`
	LicenseNo                  interface{} `json:"licenseNo"`
	LicenseName                string      `json:"licenseName"`
	LegalPersonName            string      `json:"legalPersonName"`
	LegalPersonCertificateType string      `json:"legalPersonCertificateType"`
	LegalPersonCertificateNo   interface{} `json:"legalPersonCertificateNo"`
	AcctNo                     string      `json:"acctNo"`
	AcctName                   string      `json:"acctName"`
	AcctTypeCode               string      `json:"acctTypeCode"`
	AcctCertificateType        string      `json:"acctCertificateType"`
	AcctCertificateNo          string      `json:"acctCertificateNo"`
	AcctOpenBankCode           string      `json:"acctOpenBankCode"`
	AcctOpenBankName           string      `json:"acctOpenBankName"`
	AcctClearBankCode          string      `json:"acctClearBankCode"`
	SettlePeriod               interface{} `json:"settlePeriod"`
	SettleModel                interface{} `json:"settleModel"`
	ClearDt                    interface{} `json:"clearDt"`
	OrgId                      string      `json:"orgId"`
	OrgName                    string      `json:"orgName"`
	OrgPath                    string      `json:"orgPath"`
	ReceiverStatus             interface{} `json:"receiverStatus"`
	WalletId                   string      `json:"walletId"`
	Remark                     interface{} `json:"remark"`
}

// 附件返回
type AttachFileResponse struct {
	CmdRetCode string `json:"cmdRetCode"`
	ReqId      string `json:"reqId"`
	RespData   struct {
		AttType   string `json:"attType"`
		OrderNo   string `json:"orderNo"`
		OrgCode   string `json:"orgCode"`
		AttFileId string `json:"attFileId"`
	} `json:"respData"`
	RetCode   string `json:"retCode"`
	RetMsg    string `json:"retMsg"`
	Timestamp int64  `json:"timestamp"`
	Ver       string `json:"ver"`
}

/* *
* 解析V2返回结果
 */
func ParseV2ResponseToStruct[T any](str string) (*ApiResponseV2[T], error) {
	var response ApiResponseV2[T]
	err := json.Unmarshal([]byte(str), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

/* *
* 解析V2返回结果
 */
func ParseV3ResponseToStruct[T any](str string) (*ApiResponseV3[T], error) {
	var response ApiResponseV3[T]
	err := json.Unmarshal([]byte(str), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
