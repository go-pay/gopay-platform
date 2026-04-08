package dm

import "time"

// PaymentChannel 支付通道表
type PaymentChannel struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"column:name" json:"name"`
	Code         string    `gorm:"column:code" json:"code"`
	Type         string    `gorm:"column:type" json:"type"`
	MerchantID   int64     `gorm:"column:merchant_id" json:"merchantId"`
	MerchantName string    `gorm:"-" json:"merchantName"`
	PayMethods   string    `gorm:"column:pay_methods" json:"-"`
	FeeRate      float64   `gorm:"column:fee_rate" json:"feeRate"`
	Status       int8      `gorm:"column:status" json:"status"`
	Remark       string    `gorm:"column:remark" json:"remark"`
	Ctime        time.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime        time.Time `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (PaymentChannel) TableName() string { return "payment_channel" }

// PaymentChannelConfig 支付通道参数配置表
type PaymentChannelConfig struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ChannelID  int64     `gorm:"column:channel_id" json:"channelId"`
	AppID      string    `gorm:"column:app_id" json:"appId"`
	MchID      string    `gorm:"column:mch_id" json:"mchId"`
	PrivateKey string    `gorm:"column:private_key" json:"-"`
	PublicKey  string    `gorm:"column:public_key" json:"-"`
	ApiKey     string    `gorm:"column:api_key" json:"-"`
	SerialNo   string    `gorm:"column:serial_no" json:"serialNo"`
	NotifyUrl  string    `gorm:"column:notify_url" json:"notifyUrl"`
	SignType   string    `gorm:"column:sign_type" json:"signType"`
	Sandbox    int8      `gorm:"column:sandbox" json:"sandbox"`
	Ctime      time.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime      time.Time `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (PaymentChannelConfig) TableName() string { return "payment_channel_config" }
