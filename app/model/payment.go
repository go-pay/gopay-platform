package model

type AlipayGetPaymentQrcodeReq struct {
	Subject string `json:"subject"`
	Money   int64  `json:"money"` // 分
}

type AlipayGetPaymentQrcodeRsp struct {
	OutTradeNo string `json:"out_trade_no"`
	QrCode     string `json:"qr_code"`
}

type AlipayPagePayUrlReq struct {
	Subject string `json:"subject"`
	Money   int64  `json:"money"` // 分
}

type AlipayPagePayUrlRsp struct {
	OutTradeNo string `json:"out_trade_no"`
	PagePayUrl string `json:"page_pay_url"`
}

type WxGetPaymentQrCodeReq struct {
	GoodsId int `json:"goods_id"`
}

type WxGetPaymentQrcodeRsp struct {
	OutTradeNo string `json:"out_trade_no"`
	QrCode     string `json:"qr_code"`
}
