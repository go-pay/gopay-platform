package dao

import (
	"context"
	"errors"
	"time"

	"gopay/app/dm"
	"gopay/app/model"

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

// GetAccountByID 根据 ID 查询账户
func (d *Dao) GetAccountByID(ctx context.Context, id int64) (account *dm.Account, err error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	account = new(dm.Account)
	err = d.GopayDB.WithContext(ctx).Where("id = ?", id).First(account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

// UpdateAccountLastLogin 更新最后登录时间
func (d *Dao) UpdateAccountLastLogin(ctx context.Context, id int64, t time.Time) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.Account{}).Where("id = ?", id).Update("last_login", t).Error
}

// UpdateAccountPwd 更新密码
func (d *Dao) UpdateAccountPwd(ctx context.Context, id int64, pwd string) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.Account{}).Where("id = ?", id).Update("pwd", pwd).Error
}

// UpdateAccountProfile 更新个人资料
func (d *Dao) UpdateAccountProfile(ctx context.Context, id int64, realName, phone, email string) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.Account{}).Where("id = ?", id).
		Updates(map[string]interface{}{"real_name": realName, "phone": phone, "email": email}).Error
}

// AccountList 用户列表
func (d *Dao) AccountList(ctx context.Context, req *model.UserListReq) (list []*dm.Account, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.Account{})
	if req.Username != "" {
		db = db.Where("uname LIKE ?", "%"+req.Username+"%")
	}
	if req.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+req.Phone+"%")
	}
	if req.Status >= 0 {
		db = db.Where("status = ?", req.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}

// CreateAccount 创建用户
func (d *Dao) CreateAccount(ctx context.Context, account *dm.Account) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Create(account).Error
}

// UpdateAccount 更新用户信息
func (d *Dao) UpdateAccount(ctx context.Context, id int64, updates map[string]interface{}) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.Account{}).Where("id = ?", id).Updates(updates).Error
}

// ToggleAccountStatus 切换用户状态
func (d *Dao) ToggleAccountStatus(ctx context.Context, id int64) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Exec("UPDATE account SET status = IF(status=1,0,1) WHERE id = ?", id).Error
}
