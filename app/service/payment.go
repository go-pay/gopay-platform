package service

import (
	"context"
	"strings"

	"gopay/app/model"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopher/xlog"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AlipayGetPaymentQrCode 支付码获取(用户扫码支付)
func (s *Service) AlipayGetPaymentQrCode(c context.Context, req *model.AlipayGetPaymentQrcodeReq) (rsp *model.AlipayGetPaymentQrcodeRsp, err error) {
	// 生成单号
	tradeNo := strings.ReplaceAll(uuid.New().String(), "-", "")
	amount := decimal.NewFromInt(req.Money).DivRound(decimal.NewFromInt(100), 2).String()
	xlog.Infof("tradeNo: %s, amount: %s", tradeNo, amount)
	// 构造参数
	bm := make(gopay.BodyMap)
	bm.Set("subject", req.Subject).
		Set("out_trade_no", tradeNo).
		Set("total_amount", amount)
	// 发起支付
	aliRsp, err := s.alipay.TradePrecreate(c, bm)
	if err != nil {
		if bizError, ok := alipay.IsBizError(err); ok {
			xlog.Errorf("s.alipay.TradePrecreate(%v), bizError:%v", bm, bizError)
			return nil, err
		}
		xlog.Errorf("s.alipay.TradePrecreate(%v), err:%v", bm, err)
		return nil, err
	}
	// return
	rsp = &model.AlipayGetPaymentQrcodeRsp{
		OutTradeNo: aliRsp.Response.OutTradeNo,
		QrCode:     aliRsp.Response.QrCode,
	}
	return rsp, nil
}

// AlipayPagePayUrl 支付宝网页支付链接地址获取
func (s *Service) AlipayPagePayUrl(c context.Context, req *model.AlipayPagePayUrlReq) (rsp *model.AlipayPagePayUrlRsp, err error) {
	// 生成单号
	tradeNo := strings.ReplaceAll(uuid.New().String(), "-", "")
	amount := decimal.NewFromInt(req.Money).DivRound(decimal.NewFromInt(100), 2).String()
	xlog.Infof("tradeNo: %s, amount: %s", tradeNo, amount)
	// 构造参数
	bm := make(gopay.BodyMap)
	bm.Set("subject", req.Subject).
		Set("out_trade_no", tradeNo).
		Set("total_amount", amount)
	// 发起支付
	pagePayUrl, err := s.alipay.TradePagePay(c, bm)
	if err != nil {
		xlog.Errorf("s.alipay.TradePagePay(%v), err:%v", bm, err)
		return nil, err
	}
	// return
	rsp = &model.AlipayPagePayUrlRsp{
		OutTradeNo: tradeNo,
		PagePayUrl: pagePayUrl,
	}
	return rsp, nil
}
