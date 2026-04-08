package service

import (
	"context"

	"gopay/app/dm"
	"gopay/app/model"
	ec "gopay/errcode"

	"github.com/go-pay/xlog"
)

// ---- 对账单 ----

// ReconBillList 对账单列表
func (s *Service) ReconBillList(ctx context.Context, req *model.ReconBillListReq) (*model.PageResp, error) {
	list, total, err := s.dao.ReconBillList(ctx, req)
	if err != nil {
		xlog.Errorf("ReconBillList, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ReconBillDetail 对账单详情
func (s *Service) ReconBillDetail(ctx context.Context, id int64) (*dm.ReconciliationBill, error) {
	bill, err := s.dao.GetReconBillByID(ctx, id)
	if err != nil {
		xlog.Errorf("ReconBillDetail GetReconBillByID(%d), err:%v", id, err)
		return nil, ec.ServerErr
	}
	if bill == nil {
		return nil, ec.NotFound
	}
	return bill, nil
}

// ReconBillGenerate 生成对账单
func (s *Service) ReconBillGenerate(ctx context.Context, req *model.ReconBillGenerateReq) error {
	// 检查是否已存在
	existing, err := s.dao.GetReconBillByDateAndChannel(ctx, req.Date, req.ChannelType)
	if err != nil {
		xlog.Errorf("ReconBillGenerate GetReconBillByDateAndChannel, err:%v", err)
		return ec.ServerErr
	}
	if existing != nil {
		return ec.Conflict
	}
	// 统计平台订单数据
	count, amount, err := s.dao.GetPlatformStats(ctx, req.Date, req.ChannelType)
	if err != nil {
		xlog.Errorf("ReconBillGenerate GetPlatformStats, err:%v", err)
		return ec.ServerErr
	}
	bill := &dm.ReconciliationBill{
		BillDate:       req.Date,
		ChannelType:    req.ChannelType,
		PlatformCount:  int(count),
		PlatformAmount: amount,
		Status:         0, // 待对账
	}
	if err = s.dao.CreateReconBill(ctx, bill); err != nil {
		xlog.Errorf("ReconBillGenerate CreateReconBill, err:%v", err)
		return ec.ServerErr
	}
	return nil
}

// ReconBillReconcile 执行对账
func (s *Service) ReconBillReconcile(ctx context.Context, id int64) error {
	bill, err := s.dao.GetReconBillByID(ctx, id)
	if err != nil {
		xlog.Errorf("ReconBillReconcile GetReconBillByID(%d), err:%v", id, err)
		return ec.ServerErr
	}
	if bill == nil {
		return ec.NotFound
	}
	if bill.Status != 0 {
		return ec.RequestErr
	}
	// TODO: 实际对账逻辑应下载通道对账文件并逐笔比对
	// 暂时模拟: 假设通道数据与平台一致，无差异
	updates := map[string]interface{}{
		"channel_count":  bill.PlatformCount,
		"channel_amount": bill.PlatformAmount,
		"diff_count":     0,
		"diff_amount":    0,
		"status":         1, // 已对账
	}
	if err = s.dao.UpdateReconBill(ctx, id, updates); err != nil {
		xlog.Errorf("ReconBillReconcile UpdateReconBill(%d), err:%v", id, err)
		return ec.ServerErr
	}
	return nil
}

// ---- 对账差异 ----

// ReconDiffList 对账差异列表
func (s *Service) ReconDiffList(ctx context.Context, req *model.ReconDiffListReq) (*model.PageResp, error) {
	list, total, err := s.dao.ReconDiffList(ctx, req)
	if err != nil {
		xlog.Errorf("ReconDiffList, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ReconDiffDetail 对账差异详情
func (s *Service) ReconDiffDetail(ctx context.Context, id int64) (*dm.ReconciliationDiff, error) {
	diff, err := s.dao.GetReconDiffByID(ctx, id)
	if err != nil {
		xlog.Errorf("ReconDiffDetail GetReconDiffByID(%d), err:%v", id, err)
		return nil, ec.ServerErr
	}
	if diff == nil {
		return nil, ec.NotFound
	}
	return diff, nil
}

// ReconDiffHandle 处理对账差异
func (s *Service) ReconDiffHandle(ctx context.Context, req *model.ReconDiffHandleReq, handler string) error {
	diff, err := s.dao.GetReconDiffByID(ctx, req.ID)
	if err != nil {
		xlog.Errorf("ReconDiffHandle GetReconDiffByID(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	if diff == nil {
		return ec.NotFound
	}
	if diff.HandleStatus != 0 {
		return ec.RequestErr
	}
	var status int8
	switch req.Action {
	case "resolve":
		status = 1
	case "ignore":
		status = 2
	default:
		return ec.RequestErr
	}
	updates := map[string]interface{}{
		"handle_status": status,
		"handle_remark": req.Remark,
		"handler":       handler,
	}
	if err = s.dao.UpdateReconDiff(ctx, req.ID, updates); err != nil {
		xlog.Errorf("ReconDiffHandle UpdateReconDiff(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	return nil
}

// ReconDiffExport 导出对账差异
func (s *Service) ReconDiffExport(ctx context.Context, req *model.ReconDiffListReq) ([]*dm.ReconciliationDiff, error) {
	list, err := s.dao.ReconDiffListAll(ctx, req)
	if err != nil {
		xlog.Errorf("ReconDiffExport, err:%v", err)
		return nil, ec.ServerErr
	}
	return list, nil
}
