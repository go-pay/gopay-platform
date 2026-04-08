package dm

import "time"

// PaymentOrder 支付订单表 (扩展后)
type PaymentOrder struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo       string     `gorm:"column:order_no" json:"orderNo"`
	UserID        int64      `gorm:"column:user_id" json:"userId"`
	MerchantID    int64      `gorm:"column:merchant_id" json:"merchantId"`
	MerchantName  string     `gorm:"column:merchant_name" json:"merchantName"`
	Qrcode        string     `gorm:"column:qrcode" json:"qrcode"`
	TradeNo       string     `gorm:"column:trade_no" json:"tradeNo"`
	OutTradeNo    string     `gorm:"column:out_trade_no" json:"outTradeNo"`
	TransactionID string     `gorm:"column:transaction_id" json:"transactionId"`
	PaymentType   int8       `gorm:"column:payment_type" json:"paymentType"`
	ChannelType   string     `gorm:"column:channel_type" json:"channelType"`
	PayMethod     string     `gorm:"column:pay_method" json:"payMethod"`
	Subject       string     `gorm:"column:subject" json:"subject"`
	ClientIP      string     `gorm:"column:client_ip" json:"clientIp"`
	PayMoney      int64      `gorm:"column:pay_money" json:"amount"`
	Status        int8       `gorm:"column:status" json:"status"`
	PayTime       *time.Time `gorm:"column:pay_time" json:"payTime"`
	Remark        string     `gorm:"column:remark" json:"remark"`
	NotifyBody    string     `gorm:"column:notify_body" json:"-"`
	Notified      int8       `gorm:"column:notified" json:"notified"`
	Ctime         time.Time  `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime         time.Time  `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (PaymentOrder) TableName() string { return "payment_order" }

// RefundOrder 退款订单表
type RefundOrder struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RefundNo      string     `gorm:"column:refund_no" json:"refundNo"`
	OrderNo       string     `gorm:"column:order_no" json:"orderNo"`
	TradeRefundNo string     `gorm:"column:trade_refund_no" json:"tradeRefundNo"`
	MerchantID    int64      `gorm:"column:merchant_id" json:"merchantId"`
	MerchantName  string     `gorm:"column:merchant_name" json:"merchantName"`
	RefundAmount  int64      `gorm:"column:refund_amount" json:"refundAmount"`
	OrderAmount   int64      `gorm:"column:order_amount" json:"orderAmount"`
	ChannelType   string     `gorm:"column:channel_type" json:"channelType"`
	Status        int8       `gorm:"column:status" json:"status"`
	Reason        string     `gorm:"column:reason" json:"reason"`
	Operator      string     `gorm:"column:operator" json:"operator"`
	FinishTime    *time.Time `gorm:"column:finish_time" json:"finishTime"`
	Ctime         time.Time  `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime         time.Time  `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (RefundOrder) TableName() string { return "refund_order" }

// TransferOrder 转账订单表
type TransferOrder struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TransferNo      string     `gorm:"column:transfer_no" json:"transferNo"`
	TradeTransferNo string     `gorm:"column:trade_transfer_no" json:"tradeTransferNo"`
	MerchantID      int64      `gorm:"column:merchant_id" json:"merchantId"`
	MerchantName    string     `gorm:"column:merchant_name" json:"merchantName"`
	Amount          int64      `gorm:"column:amount" json:"amount"`
	ChannelType     string     `gorm:"column:channel_type" json:"channelType"`
	PayeeType       string     `gorm:"column:payee_type" json:"payeeType"`
	PayeeAccount    string     `gorm:"column:payee_account" json:"payeeAccount"`
	PayeeName       string     `gorm:"column:payee_name" json:"payeeName"`
	Status          int8       `gorm:"column:status" json:"status"`
	Remark          string     `gorm:"column:remark" json:"remark"`
	FinishTime      *time.Time `gorm:"column:finish_time" json:"finishTime"`
	Ctime           time.Time  `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime           time.Time  `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (TransferOrder) TableName() string { return "transfer_order" }
