package dm

import (
	"github.com/go-pay/xtime"
)

// Merchant 商户表
type Merchant struct {
	ID      int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name    string     `gorm:"column:name" json:"name"`
	Contact string     `gorm:"column:contact" json:"contact"`
	Phone   string     `gorm:"column:phone" json:"phone"`
	Email   string     `gorm:"column:email" json:"email"`
	Status  int8       `gorm:"column:status" json:"status"`
	Remark  string     `gorm:"column:remark" json:"remark"`
	Ctime   xtime.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime   xtime.Time `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (Merchant) TableName() string { return "merchant" }

// MerchantApp 商户应用表
type MerchantApp struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string     `gorm:"column:name" json:"name"`
	Appid        string     `gorm:"column:appid" json:"appid"`
	MerchantID   int64      `gorm:"column:merchant_id" json:"merchantId"`
	MerchantName string     `gorm:"-" json:"merchantName"`
	PlatformType int8       `gorm:"column:platform_type" json:"platformType"`
	MerchantType int8       `gorm:"column:merchant_type" json:"merchantType"`
	NotifyUrl    string     `gorm:"column:notify_url" json:"notifyUrl"`
	ReturnUrl    string     `gorm:"column:return_url" json:"returnUrl"`
	Status       int8       `gorm:"column:status" json:"status"`
	Ctime        xtime.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime        xtime.Time `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (MerchantApp) TableName() string { return "merchant_app" }
