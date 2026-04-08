package dao

import (
	"context"
	"errors"

	"gopay/app/dm"
	"gopay/app/model"

	"gorm.io/gorm"
)

// ChannelList 支付通道列表
func (d *Dao) ChannelList(ctx context.Context, req *model.ChannelListReq) (list []*dm.PaymentChannel, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.PaymentChannel{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
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

// GetChannelByID 根据 ID 查询通道
func (d *Dao) GetChannelByID(ctx context.Context, id int64) (*dm.PaymentChannel, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	ch := new(dm.PaymentChannel)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(ch).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return ch, err
}

// GetChannelByCode 根据 code 查询通道
func (d *Dao) GetChannelByCode(ctx context.Context, code string) (*dm.PaymentChannel, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	ch := new(dm.PaymentChannel)
	err := d.GopayDB.WithContext(ctx).Where("code = ?", code).First(ch).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return ch, err
}

// CreateChannel 创建通道
func (d *Dao) CreateChannel(ctx context.Context, ch *dm.PaymentChannel) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Create(ch).Error
}

// UpdateChannel 更新通道
func (d *Dao) UpdateChannel(ctx context.Context, id int64, updates map[string]interface{}) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.PaymentChannel{}).Where("id = ?", id).Updates(updates).Error
}

// ToggleChannelStatus 切换通道状态
func (d *Dao) ToggleChannelStatus(ctx context.Context, id int64) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Exec("UPDATE payment_channel SET status = IF(status=1,0,1) WHERE id = ?", id).Error
}

// GetChannelConfig 获取通道配置
func (d *Dao) GetChannelConfig(ctx context.Context, channelID int64) (*dm.PaymentChannelConfig, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	cfg := new(dm.PaymentChannelConfig)
	err := d.GopayDB.WithContext(ctx).Where("channel_id = ?", channelID).First(cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return cfg, err
}

// SaveChannelConfig 保存通道配置（upsert）
func (d *Dao) SaveChannelConfig(ctx context.Context, cfg *dm.PaymentChannelConfig) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	// 查询是否已存在
	existing := new(dm.PaymentChannelConfig)
	err := d.GopayDB.WithContext(ctx).Where("channel_id = ?", cfg.ChannelID).First(existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return d.GopayDB.WithContext(ctx).Create(cfg).Error
	}
	if err != nil {
		return err
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.PaymentChannelConfig{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
		"app_id":      cfg.AppID,
		"mch_id":      cfg.MchID,
		"private_key": cfg.PrivateKey,
		"public_key":  cfg.PublicKey,
		"api_key":     cfg.ApiKey,
		"serial_no":   cfg.SerialNo,
		"notify_url":  cfg.NotifyUrl,
		"sign_type":   cfg.SignType,
		"sandbox":     cfg.Sandbox,
	}).Error
}
