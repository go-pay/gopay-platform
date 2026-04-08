package service

import (
	"context"
	"fmt"
	"time"

	"gopay/app/dm"
	"gopay/app/model"
	ec "gopay/errcode"

	"github.com/go-pay/xlog"
)

// ---- 支付订单 ----

// OrderPaymentList 支付订单列表
func (s *Service) OrderPaymentList(ctx context.Context, req *model.PaymentOrderListReq) (*model.PageResp, error) {
	list, total, err := s.dao.PaymentOrderList(ctx, req)
	if err != nil {
		xlog.Errorf("OrderPaymentList, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// OrderPaymentDetail 支付订单详情
func (s *Service) OrderPaymentDetail(ctx context.Context, id int64) (*dm.PaymentOrder, error) {
	order, err := s.dao.GetPaymentOrderByID(ctx, id)
	if err != nil {
		xlog.Errorf("OrderPaymentDetail GetPaymentOrderByID(%d), err:%v", id, err)
		return nil, ec.ServerErr
	}
	if order == nil {
		return nil, ec.NotFound
	}
	return order, nil
}

// OrderPaymentClose 关闭订单 (status 0→3)
func (s *Service) OrderPaymentClose(ctx context.Context, id int64) error {
	order, err := s.dao.GetPaymentOrderByID(ctx, id)
	if err != nil {
		xlog.Errorf("OrderPaymentClose GetPaymentOrderByID(%d), err:%v", id, err)
		return ec.ServerErr
	}
	if order == nil {
		return ec.NotFound
	}
	if order.Status != 0 {
		return ec.RequestErr
	}
	return s.dao.UpdatePaymentOrderStatus(ctx, id, 3)
}

// OrderPaymentRefund 发起退款
func (s *Service) OrderPaymentRefund(ctx context.Context, req *model.PaymentOrderRefundReq, operator string) error {
	order, err := s.dao.GetPaymentOrderByID(ctx, req.ID)
	if err != nil {
		xlog.Errorf("OrderPaymentRefund GetPaymentOrderByID(%d), err:%v", req.ID, err)
		return ec.ServerErr
	}
	if order == nil {
		return ec.NotFound
	}
	if order.Status != 1 {
		return ec.RequestErr
	}
	if req.Amount > order.PayMoney {
		return ec.RequestErr
	}
	// 生成退款编号
	count, _ := s.dao.CountTodayRefunds(ctx)
	refundNo := fmt.Sprintf("REF%s%03d", time.Now().Format("20060102"), count+1)

	refund := &dm.RefundOrder{
		RefundNo:     refundNo,
		OrderNo:      order.OrderNo,
		MerchantID:   order.MerchantID,
		MerchantName: order.MerchantName,
		RefundAmount: req.Amount,
		OrderAmount:  order.PayMoney,
		ChannelType:  order.ChannelType,
		Status:       0, // 退款中
		Reason:       req.Reason,
		Operator:     operator,
	}
	if err = s.dao.CreateRefundOrder(ctx, refund); err != nil {
		xlog.Errorf("OrderPaymentRefund CreateRefundOrder, err:%v", err)
		return ec.ServerErr
	}
	// TODO: 调用支付宝/微信退款 API
	return nil
}

// OrderPaymentExport 导出支付订单
func (s *Service) OrderPaymentExport(ctx context.Context, req *model.PaymentOrderListReq) ([]*dm.PaymentOrder, error) {
	list, err := s.dao.PaymentOrderListAll(ctx, req)
	if err != nil {
		xlog.Errorf("OrderPaymentExport, err:%v", err)
		return nil, ec.ServerErr
	}
	return list, nil
}

// ---- 退款订单 ----

// OrderRefundList 退款订单列表
func (s *Service) OrderRefundList(ctx context.Context, req *model.RefundOrderListReq) (*model.PageResp, error) {
	list, total, err := s.dao.RefundOrderList(ctx, req)
	if err != nil {
		xlog.Errorf("OrderRefundList, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// OrderRefundDetail 退款订单详情
func (s *Service) OrderRefundDetail(ctx context.Context, id int64) (*dm.RefundOrder, error) {
	order, err := s.dao.GetRefundOrderByID(ctx, id)
	if err != nil {
		xlog.Errorf("OrderRefundDetail GetRefundOrderByID(%d), err:%v", id, err)
		return nil, ec.ServerErr
	}
	if order == nil {
		return nil, ec.NotFound
	}
	return order, nil
}

// ---- 转账订单 ----

// OrderTransferList 转账订单列表
func (s *Service) OrderTransferList(ctx context.Context, req *model.TransferOrderListReq) (*model.PageResp, error) {
	list, total, err := s.dao.TransferOrderList(ctx, req)
	if err != nil {
		xlog.Errorf("OrderTransferList, err:%v", err)
		return nil, ec.ServerErr
	}
	return &model.PageResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// OrderTransferAdd 发起转账
func (s *Service) OrderTransferAdd(ctx context.Context, req *model.TransferOrderAddReq) error {
	// 验证商户
	m, err := s.dao.GetMerchantByID(ctx, req.MerchantID)
	if err != nil {
		xlog.Errorf("OrderTransferAdd GetMerchantByID(%d), err:%v", req.MerchantID, err)
		return ec.ServerErr
	}
	if m == nil {
		return ec.NotFound
	}
	count, _ := s.dao.CountTodayTransfers(ctx)
	transferNo := fmt.Sprintf("TRF%s%03d", time.Now().Format("20060102"), count+1)

	transfer := &dm.TransferOrder{
		TransferNo:   transferNo,
		MerchantID:   req.MerchantID,
		MerchantName: m.Name,
		Amount:       req.Amount,
		ChannelType:  req.ChannelType,
		PayeeType:    req.PayeeType,
		PayeeAccount: req.PayeeAccount,
		PayeeName:    req.PayeeName,
		Status:       0, // 处理中
		Remark:       req.Remark,
	}
	if err = s.dao.CreateTransferOrder(ctx, transfer); err != nil {
		xlog.Errorf("OrderTransferAdd CreateTransferOrder, err:%v", err)
		return ec.ServerErr
	}
	// TODO: 调用支付宝/微信转账 API
	return nil
}

// OrderTransferDetail 转账订单详情
func (s *Service) OrderTransferDetail(ctx context.Context, id int64) (*dm.TransferOrder, error) {
	order, err := s.dao.GetTransferOrderByID(ctx, id)
	if err != nil {
		xlog.Errorf("OrderTransferDetail GetTransferOrderByID(%d), err:%v", id, err)
		return nil, ec.ServerErr
	}
	if order == nil {
		return nil, ec.NotFound
	}
	return order, nil
}
