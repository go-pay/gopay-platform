package dao

import (
	"context"
	"time"

	"gopay/app/dm"
	"gopay/app/model"
)

// DashboardTodayStats 今日交易统计
func (d *Dao) DashboardTodayStats(ctx context.Context) (total int64, successCount int64, totalAmount int64, err error) {
	if d.GopayDB == nil {
		return 0, 0, 0, ErrNoDatabase
	}
	today := time.Now().Format("2006-01-02")
	// 今日总订单数
	d.GopayDB.WithContext(ctx).Model(&dm.PaymentOrder{}).
		Where("DATE(ctime) = ?", today).Count(&total)
	// 今日成功笔数和金额
	type result struct {
		Count  int64
		Amount int64
	}
	var r result
	d.GopayDB.WithContext(ctx).Model(&dm.PaymentOrder{}).
		Where("DATE(ctime) = ? AND status = 1", today).
		Select("COUNT(*) as count, COALESCE(SUM(pay_money), 0) as amount").
		Scan(&r)
	return total, r.Count, r.Amount, nil
}

// DashboardPendingApply 待审核进件数
func (d *Dao) DashboardPendingApply(ctx context.Context) (int64, error) {
	if d.GopayDB == nil {
		return 0, ErrNoDatabase
	}
	var count int64
	err := d.GopayDB.WithContext(ctx).Model(&dm.IncomingApply{}).Where("status = 1").Count(&count).Error
	return count, err
}

// DashboardPendingRefund 待处理退款数
func (d *Dao) DashboardPendingRefund(ctx context.Context) (int64, error) {
	if d.GopayDB == nil {
		return 0, ErrNoDatabase
	}
	var count int64
	err := d.GopayDB.WithContext(ctx).Model(&dm.RefundOrder{}).Where("status = 0").Count(&count).Error
	return count, err
}

// DashboardRecentOrders 最近订单
func (d *Dao) DashboardRecentOrders(ctx context.Context, req *model.PageReq) (list []*dm.PaymentOrder, total int64, err error) {
	if d.GopayDB == nil {
		return nil, 0, ErrNoDatabase
	}
	db := d.GopayDB.WithContext(ctx).Model(&dm.PaymentOrder{})
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}

// DashboardChannelDistribution 通道金额分布（今日已支付订单）
func (d *Dao) DashboardChannelDistribution(ctx context.Context) (alipay, wechat int64, err error) {
	if d.GopayDB == nil {
		return 0, 0, ErrNoDatabase
	}
	today := time.Now().Format("2006-01-02")
	type result struct {
		Amount int64
	}
	var r result
	d.GopayDB.WithContext(ctx).Model(&dm.PaymentOrder{}).
		Where("DATE(ctime) = ? AND status = 1 AND channel_type = 'alipay'", today).
		Select("COALESCE(SUM(pay_money), 0) as amount").Scan(&r)
	alipay = r.Amount
	d.GopayDB.WithContext(ctx).Model(&dm.PaymentOrder{}).
		Where("DATE(ctime) = ? AND status = 1 AND channel_type = 'wechat'", today).
		Select("COALESCE(SUM(pay_money), 0) as amount").Scan(&r)
	wechat = r.Amount
	return
}

// DashboardTrend 近7天趋势
func (d *Dao) DashboardTrend(ctx context.Context) (*model.DashboardTrendResp, error) {
	if d.GopayDB == nil {
		return nil, ErrNoDatabase
	}
	resp := &model.DashboardTrendResp{
		Dates:   make([]string, 7),
		Amounts: make([]int64, 7),
		Counts:  make([]int64, 7),
	}
	now := time.Now()
	for i := 6; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		dateStr := day.Format("2006-01-02")
		resp.Dates[6-i] = day.Format("01-02")

		type result struct {
			Count  int64
			Amount int64
		}
		var r result
		d.GopayDB.WithContext(ctx).Model(&dm.PaymentOrder{}).
			Where("DATE(ctime) = ? AND status = 1", dateStr).
			Select("COUNT(*) as count, COALESCE(SUM(pay_money), 0) as amount").
			Scan(&r)
		resp.Amounts[6-i] = r.Amount
		resp.Counts[6-i] = r.Count
	}
	return resp, nil
}
