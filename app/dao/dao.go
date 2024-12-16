package dao

import (
	"gopay/app/conf"

	"github.com/go-pay/orm"
	"gorm.io/gorm"
)

type Dao struct {
	cfg   *conf.Config
	BizDB *gorm.DB
	PayDB *gorm.DB
}

func New(c *conf.Config) (d *Dao) {
	bizdb := orm.InitGorm(c.MySQL.BizDB)
	paydb := orm.InitGorm(c.MySQL.PayDB)
	//rds := orm.InitRedis(c.Redis.Gopay)

	d = &Dao{
		cfg:   c,
		BizDB: bizdb,
		PayDB: paydb,
		//GopayRds: rds,
	}
	return
}

func (d *Dao) Close() {
	if d.BizDB != nil {
		db, _ := d.BizDB.DB()
		if db != nil {
			db.Close()
		}
	}
	if d.PayDB != nil {
		db, _ := d.PayDB.DB()
		if db != nil {
			db.Close()
		}
	}
}
