package dao

import (
	"gopay/app/cfg"

	"github.com/go-pay/gopher/orm"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dao struct {
	cfg      *cfg.Config
	GopayDB  *gorm.DB
	GopayRds *redis.Client
}

func New(c *cfg.Config) (d *Dao) {
	db := orm.InitGorm(c.MySQL.Gopay)
	rds := orm.InitRedis(c.Redis.Gopay)

	d = &Dao{
		cfg:      c,
		GopayDB:  db,
		GopayRds: rds,
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
