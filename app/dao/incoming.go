package dao

import (
	"context"
	"errors"
	"time"

	"gopay/app/dm"
	"gopay/app/model"

	"gorm.io/gorm"
)

// IncomingApplyList 进件申请列表
func (d *Dao) IncomingApplyList(ctx context.Context, req *model.IncomingApplyListReq) (list []*dm.IncomingApply, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.IncomingApply{})
	if req.ChannelType != "" {
		db = db.Where("channel_type = ?", req.ChannelType)
	}
	if req.Status >= 0 {
		db = db.Where("status = ?", req.Status)
	}
	// merchantName 需要 join 或者子查询
	if req.MerchantName != "" {
		db = db.Where("merchant_id IN (SELECT id FROM merchant WHERE name LIKE ?)", "%"+req.MerchantName+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}

// GetIncomingApplyByID 根据 ID 查询进件申请
func (d *Dao) GetIncomingApplyByID(ctx context.Context, id int64) (*dm.IncomingApply, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	apply := new(dm.IncomingApply)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(apply).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return apply, err
}

// CreateIncomingApply 创建进件申请
func (d *Dao) CreateIncomingApply(ctx context.Context, apply *dm.IncomingApply) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Create(apply).Error
}

// UpdateIncomingApplyStatus 更新进件状态
func (d *Dao) UpdateIncomingApplyStatus(ctx context.Context, id int64, status int8) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.IncomingApply{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateIncomingApplyReview 审核进件
func (d *Dao) UpdateIncomingApplyReview(ctx context.Context, id int64, status int8, reviewer, remark string) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	now := time.Now()
	return d.GopayDB.WithContext(ctx).Model(&dm.IncomingApply{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"reviewer":      reviewer,
		"review_remark": remark,
		"review_time":   &now,
	}).Error
}

// CountTodayApplies 统计今日进件数量（用于生成编号）
func (d *Dao) CountTodayApplies(ctx context.Context) (int64, error) {
	if d.GopayDB == nil {
		return 0, ErrNoDatabase
	}
	var count int64
	today := time.Now().Format("2006-01-02")
	err := d.GopayDB.WithContext(ctx).Model(&dm.IncomingApply{}).
		Where("DATE(ctime) = ?", today).Count(&count).Error
	return count, err
}

// IncomingRecordList 进件记录列表（仅 status=2 或 3）
func (d *Dao) IncomingRecordList(ctx context.Context, req *model.IncomingRecordListReq) (list []*dm.IncomingApply, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.IncomingApply{}).Where("status IN (2, 3)")
	if req.ChannelType != "" {
		db = db.Where("channel_type = ?", req.ChannelType)
	}
	if req.Status >= 2 {
		db = db.Where("status = ?", req.Status)
	}
	if req.MerchantName != "" {
		db = db.Where("merchant_id IN (SELECT id FROM merchant WHERE name LIKE ?)", "%"+req.MerchantName+"%")
	}
	if req.ReviewDate != "" {
		db = db.Where("DATE(review_time) = ?", req.ReviewDate)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}
