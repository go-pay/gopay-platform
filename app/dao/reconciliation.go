package dao

import (
	"context"
	"errors"

	"gopay/app/dm"
	"gopay/app/model"

	"gorm.io/gorm"
)

// ---- 对账单 ----

// ReconBillList 对账单列表
func (d *Dao) ReconBillList(ctx context.Context, req *model.ReconBillListReq) (list []*dm.ReconciliationBill, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.ReconciliationBill{})
	if req.Date != "" {
		db = db.Where("bill_date = ?", req.Date)
	}
	if req.ChannelType != "" {
		db = db.Where("channel_type = ?", req.ChannelType)
	}
	if req.Status >= 0 {
		db = db.Where("status = ?", req.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("bill_date DESC, id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}

// GetReconBillByID 根据 ID 查询对账单
func (d *Dao) GetReconBillByID(ctx context.Context, id int64) (*dm.ReconciliationBill, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	bill := new(dm.ReconciliationBill)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(bill).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return bill, err
}

// GetReconBillByDateAndChannel 按日期和通道查询
func (d *Dao) GetReconBillByDateAndChannel(ctx context.Context, date, channelType string) (*dm.ReconciliationBill, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	bill := new(dm.ReconciliationBill)
	err := d.GopayDB.WithContext(ctx).Where("bill_date = ? AND channel_type = ?", date, channelType).First(bill).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return bill, err
}

// CreateReconBill 创建对账单
func (d *Dao) CreateReconBill(ctx context.Context, bill *dm.ReconciliationBill) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Create(bill).Error
}

// UpdateReconBill 更新对账单
func (d *Dao) UpdateReconBill(ctx context.Context, id int64, updates map[string]interface{}) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.ReconciliationBill{}).Where("id = ?", id).Updates(updates).Error
}

// GetPlatformStats 获取平台指定日期和通道的订单统计
func (d *Dao) GetPlatformStats(ctx context.Context, date, channelType string) (count int64, amount int64, err error) {
	if d.GopayDB == nil {
		return 0, 0, ErrNoDatabase
	}
	type result struct {
		Count  int64
		Amount int64
	}
	var r result
	err = d.GopayDB.WithContext(ctx).Model(&dm.PaymentOrder{}).
		Where("DATE(ctime) = ? AND channel_type = ? AND status = 1", date, channelType).
		Select("COUNT(*) as count, COALESCE(SUM(pay_money), 0) as amount").
		Scan(&r).Error
	return r.Count, r.Amount, err
}

// ---- 对账差异 ----

// ReconDiffList 对账差异列表
func (d *Dao) ReconDiffList(ctx context.Context, req *model.ReconDiffListReq) (list []*dm.ReconciliationDiff, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.ReconciliationDiff{})
	if req.BillDate != "" {
		db = db.Where("bill_date = ?", req.BillDate)
	}
	if req.ChannelType != "" {
		db = db.Where("channel_type = ?", req.ChannelType)
	}
	if req.DiffType != "" {
		db = db.Where("diff_type = ?", req.DiffType)
	}
	if req.HandleStatus >= 0 {
		db = db.Where("handle_status = ?", req.HandleStatus)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}

// GetReconDiffByID 根据 ID 查询对账差异
func (d *Dao) GetReconDiffByID(ctx context.Context, id int64) (*dm.ReconciliationDiff, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	diff := new(dm.ReconciliationDiff)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(diff).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return diff, err
}

// UpdateReconDiff 更新对账差异
func (d *Dao) UpdateReconDiff(ctx context.Context, id int64, updates map[string]interface{}) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.ReconciliationDiff{}).Where("id = ?", id).Updates(updates).Error
}

// ReconDiffListAll 对账差异列表（不分页，用于导出）
func (d *Dao) ReconDiffListAll(ctx context.Context, req *model.ReconDiffListReq) (list []*dm.ReconciliationDiff, err error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.ReconciliationDiff{})
	if req.BillDate != "" {
		db = db.Where("bill_date = ?", req.BillDate)
	}
	if req.ChannelType != "" {
		db = db.Where("channel_type = ?", req.ChannelType)
	}
	if req.DiffType != "" {
		db = db.Where("diff_type = ?", req.DiffType)
	}
	if req.HandleStatus >= 0 {
		db = db.Where("handle_status = ?", req.HandleStatus)
	}
	err = db.Order("id DESC").Limit(10000).Find(&list).Error
	return
}
