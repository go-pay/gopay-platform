package model

// ---- 交易流水 ----

type FlowListReq struct {
	PageReq
	FlowNo      string `json:"flowNo"`
	OrderNo     string `json:"orderNo"`
	Type        string `json:"type"`        // ""=全部, pay/refund/transfer
	ChannelType string `json:"channelType"` // ""=全部
	Date        string `json:"date"`
}

type FlowStatsResp struct {
	IncomeTotal  int64 `json:"incomeTotal"`
	ExpenseTotal int64 `json:"expenseTotal"`
	TotalCount   int64 `json:"totalCount"`
}

// ---- 回调通知 ----

type CallbackListReq struct {
	PageReq
	OrderNo     string `json:"orderNo"`
	Type        string `json:"type"`        // ""=全部
	Status      int8   `json:"status"`      // -1=全部
	ChannelType string `json:"channelType"` // ""=全部
}
