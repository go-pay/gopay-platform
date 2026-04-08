package service

import (
	"context"

	"gopay/app/dm"
	"gopay/app/model"
	ec "gopay/errcode"

	"github.com/go-pay/xlog"
)

// ---- 交易流水 ----

// FlowList 交易流水列表
func (s *Service) FlowList(ctx context.Context, req *model.FlowListReq) (*model.PageResp, error) {
	list, total, err := s.dao.FlowList(ctx, req)
	if err != nil {
		xlog.Errorf("FlowList, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// FlowDetail 交易流水详情
func (s *Service) FlowDetail(ctx context.Context, id int64) (*dm.TransactionFlow, error) {
	flow, err := s.dao.GetFlowByID(ctx, id)
	if err != nil {
		xlog.Errorf("FlowDetail GetFlowByID(%d), err:%v", id, err)
		return nil, ec.ServerErr
	}
	if flow == nil {
		return nil, ec.NotFound
	}
	return flow, nil
}

// FlowStats 交易流水统计
func (s *Service) FlowStats(ctx context.Context) (*model.FlowStatsResp, error) {
	stats, err := s.dao.FlowStats(ctx)
	if err != nil {
		xlog.Errorf("FlowStats, err:%v", err)
		return nil, ec.ServerErr
	}
	return stats, nil
}

// ---- 回调通知 ----

// CallbackList 回调通知列表
func (s *Service) CallbackList(ctx context.Context, req *model.CallbackListReq) (*model.PageResp, error) {
	list, total, err := s.dao.CallbackList(ctx, req)
	if err != nil {
		xlog.Errorf("CallbackList, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// CallbackDetail 回调通知详情
func (s *Service) CallbackDetail(ctx context.Context, id int64) (*dm.CallbackRecord, error) {
	cb, err := s.dao.GetCallbackByID(ctx, id)
	if err != nil {
		xlog.Errorf("CallbackDetail GetCallbackByID(%d), err:%v", id, err)
		return nil, ec.ServerErr
	}
	if cb == nil {
		return nil, ec.NotFound
	}
	return cb, nil
}

// CallbackRetry 手动重试回调
func (s *Service) CallbackRetry(ctx context.Context, id int64) error {
	cb, err := s.dao.GetCallbackByID(ctx, id)
	if err != nil {
		xlog.Errorf("CallbackRetry GetCallbackByID(%d), err:%v", id, err)
		return ec.ServerErr
	}
	if cb == nil {
		return ec.NotFound
	}
	if cb.Status == 1 {
		return ec.RequestErr
	}
	// TODO: 实际发送 HTTP 通知到 notifyUrl
	// 暂时模拟重试：更新状态为待重试
	if err = s.dao.UpdateCallbackRetry(ctx, id, 2, 0, ""); err != nil {
		xlog.Errorf("CallbackRetry UpdateCallbackRetry(%d), err:%v", id, err)
		return ec.ServerErr
	}
	return nil
}
