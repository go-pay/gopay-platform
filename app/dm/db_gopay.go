package dm

import "time"

// Account 账户表
type Account struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Uname     string     `gorm:"column:uname" json:"username"`
	Pwd       string     `gorm:"column:pwd" json:"-"`
	RealName  string     `gorm:"column:real_name" json:"realName"`
	Phone     string     `gorm:"column:phone" json:"phone"`
	Email     string     `gorm:"column:email" json:"email"`
	Role      string     `gorm:"column:role" json:"role"`
	Status    int8       `gorm:"column:status" json:"status"`
	LastLogin *time.Time `gorm:"column:last_login" json:"lastLogin"`
	Ctime     time.Time  `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	Utime     time.Time  `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (Account) TableName() string {
	return "account"
}
