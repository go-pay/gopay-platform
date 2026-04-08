package dm

import "github.com/go-pay/xtime"

// ReconciliationBill 对账单表
type ReconciliationBill struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BillDate       string     `gorm:"column:bill_date" json:"billDate"`
	ChannelType    string     `gorm:"column:channel_type" json:"channelType"`
	PlatformCount  int        `gorm:"column:platform_count" json:"platformCount"`
	PlatformAmount int64      `gorm:"column:platform_amount" json:"platformAmount"`
	ChannelCount   int        `gorm:"column:channel_count" json:"channelCount"`
	ChannelAmount  int64      `gorm:"column:channel_amount" json:"channelAmount"`
	DiffCount      int        `gorm:"column:diff_count" json:"diffCount"`
	DiffAmount     int64      `gorm:"column:diff_amount" json:"diffAmount"`
	Status         int8       `gorm:"column:status" json:"status"`
	Ctime          xtime.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime          xtime.Time `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (ReconciliationBill) TableName() string { return "reconciliation_bill" }

// ReconciliationDiff 对账差异表
type ReconciliationDiff struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BillDate       string     `gorm:"column:bill_date" json:"billDate"`
	OrderNo        string     `gorm:"column:order_no" json:"orderNo"`
	ChannelType    string     `gorm:"column:channel_type" json:"channelType"`
	DiffType       string     `gorm:"column:diff_type" json:"diffType"`
	PlatformAmount *int64     `gorm:"column:platform_amount" json:"platformAmount"`
	ChannelAmount  *int64     `gorm:"column:channel_amount" json:"channelAmount"`
	DiffAmount     int64      `gorm:"column:diff_amount" json:"diffAmount"`
	HandleStatus   int8       `gorm:"column:handle_status" json:"handleStatus"`
	HandleRemark   string     `gorm:"column:handle_remark" json:"handleRemark"`
	Handler        string     `gorm:"column:handler" json:"handler"`
	Ctime          xtime.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime          xtime.Time `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (ReconciliationDiff) TableName() string { return "reconciliation_diff" }
