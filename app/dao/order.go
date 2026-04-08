package dao

import (
	"context"
	"errors"
	"time"

	"gopay/app/dm"
	"gopay/app/model"

	"gorm.io/gorm"
)

// ---- 支付订单 ----

// PaymentOrderList 支付订单列表
func (d *Dao) PaymentOrderList(ctx context.Context, req *model.PaymentOrderListReq) (list []*dm.PaymentOrder, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.PaymentOrder{})
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.MerchantName != "" {
		db = db.Where("merchant_name LIKE ?", "%"+req.MerchantName+"%")
	}
	if req.Status >= 0 {
		db = db.Where("status = ?", req.Status)
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

// GetPaymentOrderByID 根据 ID 查询支付订单
func (d *Dao) GetPaymentOrderByID(ctx context.Context, id int64) (*dm.PaymentOrder, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	order := new(dm.PaymentOrder)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return order, err
}

// UpdatePaymentOrderStatus 更新支付订单状态
func (d *Dao) UpdatePaymentOrderStatus(ctx context.Context, id int64, status int8) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Model(&dm.PaymentOrder{}).Where("id = ?", id).Update("status", status).Error
}

// PaymentOrderListAll 支付订单列表（不分页，用于导出）
func (d *Dao) PaymentOrderListAll(ctx context.Context, req *model.PaymentOrderListReq) (list []*dm.PaymentOrder, err error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.PaymentOrder{})
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.MerchantName != "" {
		db = db.Where("merchant_name LIKE ?", "%"+req.MerchantName+"%")
	}
	if req.Status >= 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.ChannelType != "" {
		db = db.Where("channel_type = ?", req.ChannelType)
	}
	if req.Date != "" {
		db = db.Where("DATE(ctime) = ?", req.Date)
	}
	err = db.Order("id DESC").Limit(10000).Find(&list).Error
	return
}

// ---- 退款订单 ----

// RefundOrderList 退款订单列表
func (d *Dao) RefundOrderList(ctx context.Context, req *model.RefundOrderListReq) (list []*dm.RefundOrder, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.RefundOrder{})
	if req.RefundNo != "" {
		db = db.Where("refund_no LIKE ?", "%"+req.RefundNo+"%")
	}
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
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

// GetRefundOrderByID 根据 ID 查询退款订单
func (d *Dao) GetRefundOrderByID(ctx context.Context, id int64) (*dm.RefundOrder, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	order := new(dm.RefundOrder)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return order, err
}

// CreateRefundOrder 创建退款订单
func (d *Dao) CreateRefundOrder(ctx context.Context, order *dm.RefundOrder) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Create(order).Error
}

// CountTodayRefunds 统计今日退款数量
func (d *Dao) CountTodayRefunds(ctx context.Context) (int64, error) {
	if d.GopayDB == nil {
		return 0, ErrNoDatabase
	}
	var count int64
	today := time.Now().Format("2006-01-02")
	err := d.GopayDB.WithContext(ctx).Model(&dm.RefundOrder{}).
		Where("DATE(ctime) = ?", today).Count(&count).Error
	return count, err
}

// ---- 转账订单 ----

// TransferOrderList 转账订单列表
func (d *Dao) TransferOrderList(ctx context.Context, req *model.TransferOrderListReq) (list []*dm.TransferOrder, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.TransferOrder{})
	if req.TransferNo != "" {
		db = db.Where("transfer_no LIKE ?", "%"+req.TransferNo+"%")
	}
	if req.MerchantName != "" {
		db = db.Where("merchant_name LIKE ?", "%"+req.MerchantName+"%")
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

// GetTransferOrderByID 根据 ID 查询转账订单
func (d *Dao) GetTransferOrderByID(ctx context.Context, id int64) (*dm.TransferOrder, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	order := new(dm.TransferOrder)
	err := d.GopayDB.WithContext(ctx).Where("id = ?", id).First(order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return order, err
}

// CreateTransferOrder 创建转账订单
func (d *Dao) CreateTransferOrder(ctx context.Context, order *dm.TransferOrder) error {
	if d.GopayDB == nil {
		return ErrNoDatabase
	}
	return d.GopayDB.WithContext(ctx).Create(order).Error
}

// CountTodayTransfers 统计今日转账数量
func (d *Dao) CountTodayTransfers(ctx context.Context) (int64, error) {
	if d.GopayDB == nil {
		return 0, ErrNoDatabase
	}
	var count int64
	today := time.Now().Format("2006-01-02")
	err := d.GopayDB.WithContext(ctx).Model(&dm.TransferOrder{}).
		Where("DATE(ctime) = ?", today).Count(&count).Error
	return count, err
}
