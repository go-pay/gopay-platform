package model

// ---- 支付订单 ----

type PaymentOrderListReq struct {
	PageReq
	OrderNo      string `json:"orderNo"`
	MerchantName string `json:"merchantName"`
	Status       int8   `json:"status"`      // -1=全部
	ChannelType  string `json:"channelType"` // ""=全部
	Date         string `json:"date"`
}

type PaymentOrderRefundReq struct {
	ID     int64  `json:"id" binding:"required"`
	Amount int64  `json:"amount" binding:"required,min=1"`
	Reason string `json:"reason" binding:"required"`
}

// ---- 退款订单 ----

type RefundOrderListReq struct {
	PageReq
	RefundNo    string `json:"refundNo"`
	OrderNo     string `json:"orderNo"`
	Status      int8   `json:"status"`      // -1=全部
	ChannelType string `json:"channelType"` // ""=全部
}

// ---- 转账订单 ----

type TransferOrderListReq struct {
	PageReq
	TransferNo   string `json:"transferNo"`
	MerchantName string `json:"merchantName"`
	Status       int8   `json:"status"`      // -1=全部
	ChannelType  string `json:"channelType"` // ""=全部
}

type TransferOrderAddReq struct {
	MerchantID   int64  `json:"merchantId" binding:"required"`
	ChannelType  string `json:"channelType" binding:"required"`
	Amount       int64  `json:"amount" binding:"required,min=1"`
	PayeeType    string `json:"payeeType" binding:"required"`
	PayeeAccount string `json:"payeeAccount" binding:"required"`
	PayeeName    string `json:"payeeName" binding:"required"`
	Remark       string `json:"remark"`
}
