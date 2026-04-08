package dm

import "github.com/go-pay/xtime"

// IncomingApply 进件申请表
type IncomingApply struct {
	ID           int64       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApplyNo      string      `gorm:"column:apply_no" json:"applyNo"`
	MerchantID   int64       `gorm:"column:merchant_id" json:"merchantId"`
	MerchantName string      `gorm:"-" json:"merchantName"`
	ChannelType  string      `gorm:"column:channel_type" json:"channelType"`
	MerchantNo   string      `gorm:"column:merchant_no" json:"merchantNo"`
	LicenseNo    string      `gorm:"column:license_no" json:"licenseNo"`
	LicenseImg   string      `gorm:"column:license_img" json:"licenseImg"`
	LegalPerson  string      `gorm:"column:legal_person" json:"legalPerson"`
	IdCardFront  string      `gorm:"column:id_card_front" json:"idCardFront"`
	IdCardBack   string      `gorm:"column:id_card_back" json:"idCardBack"`
	Phone        string      `gorm:"column:phone" json:"phone"`
	Status       int8        `gorm:"column:status" json:"status"`
	Remark       string      `gorm:"column:remark" json:"remark"`
	Reviewer     string      `gorm:"column:reviewer" json:"reviewer"`
	ReviewRemark string      `gorm:"column:review_remark" json:"reviewRemark"`
	ReviewTime   *xtime.Time `gorm:"column:review_time" json:"reviewTime"`
	Ctime        xtime.Time  `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime        xtime.Time  `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (IncomingApply) TableName() string { return "incoming_apply" }
