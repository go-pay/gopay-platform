package service

import (
	"context"

	"gopay/app/model"
	ec "gopay/errcode"

	"github.com/go-pay/xlog"
)

// DashboardStats 仪表盘统计
func (s *Service) DashboardStats(ctx context.Context) (*model.DashboardStatsResp, error) {
	total, successCount, totalAmount, err := s.dao.DashboardTodayStats(ctx)
	if err != nil {
		xlog.Errorf("DashboardStats DashboardTodayStats, err:%v", err)
		return nil, ec.ServerErr
	}
	pendingApply, _ := s.dao.DashboardPendingApply(ctx)
	pendingRefund, _ := s.dao.DashboardPendingRefund(ctx)

	successRate := float64(0)
	if total > 0 {
		successRate = float64(successCount) / float64(total) * 100
		// 保留一位小数
		successRate = float64(int(successRate*10)) / 10
	}
	return &model.DashboardStatsResp{
		TodayAmount:      totalAmount,
		TodayCount:       total,
		TodaySuccessRate: successRate,
		PendingApply:     pendingApply,
		PendingRefund:    pendingRefund,
	}, nil
}

// DashboardRecentOrders 最近订单
func (s *Service) DashboardRecentOrders(ctx context.Context, req *model.PageReq) (*model.PageResp, error) {
	list, total, err := s.dao.DashboardRecentOrders(ctx, req)
	if err != nil {
		xlog.Errorf("DashboardRecentOrders, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DashboardChannelDistribution 通道分布
func (s *Service) DashboardChannelDistribution(ctx context.Context) (*model.DashboardChannelDistResp, error) {
	alipay, wechat, err := s.dao.DashboardChannelDistribution(ctx)
	if err != nil {
		xlog.Errorf("DashboardChannelDistribution, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.DashboardChannelDistResp{
		Alipay: alipay,
		Wechat: wechat,
	}, nil
}

// DashboardTrend 近7天趋势
func (s *Service) DashboardTrend(ctx context.Context) (*model.DashboardTrendResp, error) {
	resp, err := s.dao.DashboardTrend(ctx)
	if err != nil {
		xlog.Errorf("DashboardTrend, err:%v", err)
		return nil, ec.ServerErr
	}
	return resp, nil
}
