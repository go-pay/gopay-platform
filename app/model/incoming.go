package model

// ---- 进件管理 ----

type IncomingApplyListReq struct {
	PageReq
	MerchantName string `json:"merchantName"`
	Status       int8   `json:"status"`      // -1=全部
	ChannelType  string `json:"channelType"` // ""=全部
}

type IncomingApplyAddReq struct {
	MerchantID  int64  `json:"merchantId" binding:"required"`
	ChannelType string `json:"channelType" binding:"required"`
	MerchantNo  string `json:"merchantNo"`
	LicenseNo   string `json:"licenseNo"`
	LicenseImg  string `json:"licenseImg"`
	LegalPerson string `json:"legalPerson"`
	IdCardFront string `json:"idCardFront"`
	IdCardBack  string `json:"idCardBack"`
	Phone       string `json:"phone"`
	Remark      string `json:"remark"`
	Submit      bool   `json:"submit"` // true=直接提交审核
}

type IncomingApplyReviewReq struct {
	ID     int64  `json:"id" binding:"required"`
	Action string `json:"action" binding:"required"` // "pass" or "reject"
	Remark string `json:"remark"`
}

type IncomingRecordListReq struct {
	PageReq
	MerchantName string `json:"merchantName"`
	ChannelType  string `json:"channelType"`
	Status       int8   `json:"status"` // -1=全部, 仅 2或3
	ReviewDate   string `json:"reviewDate"`
}
