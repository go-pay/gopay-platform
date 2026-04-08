package dao

import (
	"context"
	"errors"

	"gopay/app/dm"
	"gopay/app/model"

	"gorm.io/gorm"
)

// ---- 交易流水 ----

// FlowList 交易流水列表
func (d *Dao) FlowList(ctx context.Context, req *model.FlowListReq) (list []*dm.TransactionFlow, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.TransactionFlow{})
	if req.FlowNo != "" {
		db = db.Where("flow_no LIKE ?", "%"+req.FlowNo+"%")
	}
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}
	if req.ChannelType != "" {
		db = db.Where("channel_type = ?", req.ChannelType)
	}
	if req.Date != "" {
		db = db.Where("DATE(ctime) = ?", req.Date)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}

// GetFlowByID 根据 ID 查询交易流水
func (d *Dao) GetFlowByID(ctx context.Context, id int64) (*dm.TransactionFlow, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	flow := new(dm.TransactionFlow)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(flow).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return flow, err
}

// FlowStats 流水统计
func (d *Dao) FlowStats(ctx context.Context) (*model.FlowStatsResp, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	stats := &model.FlowStatsResp{}
	// 收入总额
	d.GopayDB.WithContext(ctx).Model(&dm.TransactionFlow{}).
		Where("direction = 'in' AND status = 1").
		Select("COALESCE(SUM(amount), 0)").Scan(&stats.IncomeTotal)
	// 支出总额
	d.GopayDB.WithContext(ctx).Model(&dm.TransactionFlow{}).
		Where("direction = 'out' AND status = 1").
		Select("COALESCE(SUM(amount), 0)").Scan(&stats.ExpenseTotal)
	// 总笔数
	d.GopayDB.WithContext(ctx).Model(&dm.TransactionFlow{}).Count(&stats.TotalCount)
	return stats, nil
}

// ---- 回调通知 ----

// CallbackList 回调通知列表
func (d *Dao) CallbackList(ctx context.Context, req *model.CallbackListReq) (list []*dm.CallbackRecord, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.CallbackRecord{})
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}
	if req.Status >= 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.ChannelType != "" {
		db = db.Where("channel_type = ?", req.ChannelType)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}

// GetCallbackByID 根据 ID 查询回调记录
func (d *Dao) GetCallbackByID(ctx context.Context, id int64) (*dm.CallbackRecord, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	cb := new(dm.CallbackRecord)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(cb).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return cb, err
}

// UpdateCallbackRetry 更新回调重试
func (d *Dao) UpdateCallbackRetry(ctx context.Context, id int64, status int8, httpStatus int, responseBody string) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.CallbackRecord{}).Where("id = ?", id).Updates(map[string]interface{}{
		"retry_count":   gorm.Expr("retry_count + 1"),
		"status":        status,
		"http_status":   httpStatus,
		"response_body": responseBody,
	}).Error
}
