package dao

import (
	"context"
	"errors"

	"gopay/app/dm"

	"gorm.io/gorm"
)

var ErrNoDatabase = errors.New("database not initialized")

// GetAccountByUname 根据用户名查询账户，用户不存在时返回 nil, nil
func (d *Dao) GetAccountByUname(ctx context.Context, uname string) (account *dm.Account, err error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	account = new(dm.Account)
	err = d.GopayDB.WithContext(ctx).Where("uname = ?", uname).First(account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}
