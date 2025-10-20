package notification

type NotificationRequest struct {
	OutTradeNo      string `json:"out_trade_no"`
	TradeNo         string `json:"trade_no"`
	LogNo           string `json:"log_no"`
	AccTradeNo      string `json:"acc_trade_no"`
	TradeStatus     string `json:"trade_status"`
	TradeState      string `json:"trade_state"`
	TotalAmount     string `json:"total_amount"`
	PayerAmount     string `json:"payer_amount"`
	AccSettleAmount string `json:"acc_settle_amount"`
	TradeTime       string `json:"trade_time"`
	UserId1         string `json:"user_id1"`
	UserId2         string `json:"user_id2"`
	NotifyUrl       string `json:"notify_url"`
	AccountType     string `json:"account_type"`
	CardType        string `json:"card_type"`
}
