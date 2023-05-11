package service

import (
	"context"
	"sync"

	"gopay/app/cfg"
	"gopay/app/dao"

	"github.com/go-pay/gopher/smap"
)

type Service struct {
	rwMu sync.RWMutex
	Cfg  *cfg.Config
	dao  *dao.Dao

	// cache
	kvMap smap.Map[string, string] // key: k, value: v
}

var (
	srv *Service
	ctx = context.Background()
)

func New(c *cfg.Config) (s *Service) {
	srv = &Service{
		Cfg: c,
		dao: dao.New(c),
	}

	// loop job
	srv.initLoop()
	return srv
}

// Close 关闭相关资源
func (s *Service) Close() {
	if s.dao != nil {
		s.dao.Close()
	}
}

// 初始化 loop
func (s *Service) initLoop() {
	// 初始化 loop
}
