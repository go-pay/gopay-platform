package dm

import "time"

// TransactionFlow 交易流水表
type TransactionFlow struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	FlowNo        string    `gorm:"column:flow_no" json:"flowNo"`
	OrderNo       string    `gorm:"column:order_no" json:"orderNo"`
	Type          string    `gorm:"column:type" json:"type"`
	MerchantID    int64     `gorm:"column:merchant_id" json:"merchantId"`
	MerchantName  string    `gorm:"column:merchant_name" json:"merchantName"`
	Amount        int64     `gorm:"column:amount" json:"amount"`
	ChannelType   string    `gorm:"column:channel_type" json:"channelType"`
	ChannelFlowNo string    `gorm:"column:channel_flow_no" json:"channelFlowNo"`
	Direction     string    `gorm:"column:direction" json:"direction"`
	Status        int8      `gorm:"column:status" json:"status"`
	Remark        string    `gorm:"column:remark" json:"remark"`
	Ctime         time.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime         time.Time `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (TransactionFlow) TableName() string { return "transaction_flow" }

// CallbackRecord 回调通知记录表
type CallbackRecord struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo      string    `gorm:"column:order_no" json:"orderNo"`
	Type         string    `gorm:"column:type" json:"type"`
	ChannelType  string    `gorm:"column:channel_type" json:"channelType"`
	Direction    string    `gorm:"column:direction" json:"direction"`
	NotifyUrl    string    `gorm:"column:notify_url" json:"notifyUrl"`
	Status       int8      `gorm:"column:status" json:"status"`
	HttpStatus   int       `gorm:"column:http_status" json:"httpStatus"`
	RetryCount   int       `gorm:"column:retry_count" json:"retryCount"`
	MaxRetry     int       `gorm:"column:max_retry" json:"maxRetry"`
	RequestBody  string    `gorm:"column:request_body" json:"requestBody"`
	ResponseBody string    `gorm:"column:response_body" json:"responseBody"`
	Ctime        time.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime        time.Time `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (CallbackRecord) TableName() string { return "callback_record" }
