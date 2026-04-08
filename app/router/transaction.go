package router

import (
	"gopay/app/model"
	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

// ---- 交易流水 ----

func transactionFlowList(c *gin.Context) {
	req := new(model.FlowListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("transactionFlowList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.FlowList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func transactionFlowDetail(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("transactionFlowDetail ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	detail, err := svc.FlowDetail(c, req.ID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, detail, nil)
}

func transactionFlowStats(c *gin.Context) {
	stats, err := svc.FlowStats(c)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, stats, nil)
}

// ---- 回调通知 ----

func callbackList(c *gin.Context) {
	req := new(model.CallbackListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("callbackList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.CallbackList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func callbackDetail(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("callbackDetail ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	detail, err := svc.CallbackDetail(c, req.ID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, detail, nil)
}

func callbackRetry(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("callbackRetry ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.CallbackRetry(c, req.ID); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}
