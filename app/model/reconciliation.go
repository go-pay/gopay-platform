package model

// ---- 对账单 ----

type ReconBillListReq struct {
	PageReq
	Date        string `json:"date"`
	ChannelType string `json:"channelType"` // ""=全部
	Status      int8   `json:"status"`      // -1=全部
}

type ReconBillGenerateReq struct {
	Date        string `json:"date" binding:"required"`
	ChannelType string `json:"channelType" binding:"required"`
}

// ---- 对账差异 ----

type ReconDiffListReq struct {
	PageReq
	BillDate     string `json:"billDate"`
	ChannelType  string `json:"channelType"`  // ""=全部
	DiffType     string `json:"diffType"`     // ""=全部
	HandleStatus int8   `json:"handleStatus"` // -1=全部
}

type ReconDiffHandleReq struct {
	ID     int64  `json:"id" binding:"required"`
	Action string `json:"action" binding:"required"` // "resolve" or "ignore"
	Remark string `json:"remark"`
}
