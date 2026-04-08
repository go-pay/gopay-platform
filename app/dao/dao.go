package dao

import (
	"gopay/app/conf"

	"github.com/go-pay/orm"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dao struct {
	cfg      *conf.Config
	GopayDB  *gorm.DB
	GopayRds *redis.Client
}

func New(c *conf.Config) (d *Dao) {
	db := orm.InitGorm(c.MySQL.Gopay)
	//rds := orm.InitRedis(c.Redis.Gopay)

	d = &Dao{
		cfg:     c,
		GopayDB: db,
		//GopayRds: rds,
	}
	return
}

func (d *Dao) Close() {
	if d.GopayDB != nil {
		db, _ := d.GopayDB.DB()
		if db != nil {
			db.Close()
		}
	}
	if d.GopayRds != nil {
		d.GopayRds.Close()
	}
}
