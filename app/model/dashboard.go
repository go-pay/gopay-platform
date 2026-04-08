package model

// DashboardStatsResp 仪表盘统计
type DashboardStatsResp struct {
	TodayAmount      int64   `json:"todayAmount"`
	TodayCount       int64   `json:"todayCount"`
	TodaySuccessRate float64 `json:"todaySuccessRate"`
	PendingApply     int64   `json:"pendingApply"`
	PendingRefund    int64   `json:"pendingRefund"`
}

// DashboardChannelDistResp 通道分布
type DashboardChannelDistResp struct {
	Alipay int64 `json:"alipay"`
	Wechat int64 `json:"wechat"`
}

// DashboardTrendResp 近7天趋势
type DashboardTrendResp struct {
	Dates   []string `json:"dates"`
	Amounts []int64  `json:"amounts"`
	Counts  []int64  `json:"counts"`
}
