package service

import (
	"sync"

	"gopay/app/conf"
	"gopay/app/dao"

	"github.com/go-pay/gopher/smap"
)

type Service struct {
	rwMu sync.RWMutex
	Cfg  *conf.Config
	dao  *dao.Dao

	// cache
	kvMap smap.Map[string, string] // key: k, value: v
}
