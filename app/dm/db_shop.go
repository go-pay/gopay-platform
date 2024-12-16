package dm

import (
	"github.com/go-pay/xtime"
)

// 店铺商品表
type ShopGoods struct {
	ID        int        `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`               // 自增主键
	SkuID     string     `gorm:"column:sku_id;NOT NULL" json:"sku_id"`                         // 商品sku_id
	GoodsName string     `gorm:"column:goods_name;NOT NULL" json:"goods_name"`                 // 商品名称
	GoodsDesc string     `gorm:"column:goods_desc;NOT NULL" json:"goods_desc"`                 // 商品描述
	UnitPrice int        `gorm:"column:unit_price;default:0;NOT NULL" json:"unit_price"`       // 商品单价(分)
	Ctime     xtime.Time `gorm:"column:ctime;default:CURRENT_TIMESTAMP;NOT NULL" json:"ctime"` // 创建时间
	Utime     xtime.Time `gorm:"column:utime;default:CURRENT_TIMESTAMP;NOT NULL" json:"utime"` // 最后修改时间
}

func (m *ShopGoods) TableName() string {
	return "shop_goods"
}

// 支付订单表
type PaymentOrder struct {
	ID            int        `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`               // 自增长ID
	UserID        int        `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`             // 绑定用户id
	Qrcode        string     `gorm:"column:qrcode;NOT NULL" json:"qrcode"`                         // 支付二维码
	TradeNo       string     `gorm:"column:trade_no;NOT NULL" json:"trade_no"`                     // 订单号
	TransactionID string     `gorm:"column:transaction_id;NOT NULL" json:"transaction_id"`         // 支付交易流水号
	PaymentType   int        `gorm:"column:payment_type;default:0;NOT NULL" json:"payment_type"`   // 支付类型：0-微信，1-支付宝
	PayMoney      int        `gorm:"column:pay_money;default:0;NOT NULL" json:"pay_money"`         // 支付金额(分)
	Status        int        `gorm:"column:status;default:0;NOT NULL" json:"status"`               // 订单状态：0-待支付，1-支付成功，3-订单关闭
	PayTime       xtime.Time `gorm:"column:pay_time" json:"pay_time"`                              // 支付时间
	Remark        string     `gorm:"column:remark;NOT NULL" json:"remark"`                         // 备注
	NotifyBody    string     `gorm:"column:notify_body" json:"notify_body"`                        // 回调参数信息
	Ctime         xtime.Time `gorm:"column:ctime;default:CURRENT_TIMESTAMP;NOT NULL" json:"ctime"` // 创建时间
	Utime         xtime.Time `gorm:"column:utime;default:CURRENT_TIMESTAMP;NOT NULL" json:"utime"` // 最后更新时间
}

func (m *PaymentOrder) TableName() string {
	return "payment_order"
}
