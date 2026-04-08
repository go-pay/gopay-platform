package dm

import "time"

// Account 账户表
type Account struct {
	ID    int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Uname string    `gorm:"column:uname" json:"uname"`
	Pwd   string    `gorm:"column:pwd" json:"-"`
	Phone string    `gorm:"column:phone" json:"phone"`
	Ctime time.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime time.Time `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (Account) TableName() string {
	return "account"
}
