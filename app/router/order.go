package router

import (
	"gopay/app/model"
	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

// ---- 支付订单 ----

func orderPaymentList(c *gin.Context) {
	req := new(model.PaymentOrderListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("orderPaymentList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.OrderPaymentList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func orderPaymentDetail(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("orderPaymentDetail ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	detail, err := svc.OrderPaymentDetail(c, req.ID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, detail, nil)
}

func orderPaymentClose(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("orderPaymentClose ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.OrderPaymentClose(c, req.ID); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func orderPaymentRefund(c *gin.Context) {
	req := new(model.PaymentOrderRefundReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("orderPaymentRefund ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	username, _ := c.Get(CtxKeyUsername)
	operator, _ := username.(string)
	if err := svc.OrderPaymentRefund(c, req, operator); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func orderPaymentExport(c *gin.Context) {
	req := new(model.PaymentOrderListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("orderPaymentExport ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	list, err := svc.OrderPaymentExport(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, list, nil)
}

// ---- 退款订单 ----

func orderRefundList(c *gin.Context) {
	req := new(model.RefundOrderListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("orderRefundList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.OrderRefundList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func orderRefundDetail(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("orderRefundDetail ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	detail, err := svc.OrderRefundDetail(c, req.ID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, detail, nil)
}

// ---- 转账订单 ----

func orderTransferList(c *gin.Context) {
	req := new(model.TransferOrderListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("orderTransferList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.OrderTransferList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func orderTransferAdd(c *gin.Context) {
	req := new(model.TransferOrderAddReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("orderTransferAdd ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.OrderTransferAdd(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func orderTransferDetail(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("orderTransferDetail ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	detail, err := svc.OrderTransferDetail(c, req.ID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, detail, nil)
}
