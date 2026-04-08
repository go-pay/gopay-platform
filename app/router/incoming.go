package router

import (
	"gopay/app/model"
	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

func incomingApplyList(c *gin.Context) {
	req := new(model.IncomingApplyListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("incomingApplyList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.IncomingApplyList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func incomingApplyAdd(c *gin.Context) {
	req := new(model.IncomingApplyAddReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("incomingApplyAdd ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.IncomingApplyAdd(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func incomingApplySubmit(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("incomingApplySubmit ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.IncomingApplySubmit(c, req.ID); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func incomingApplyReview(c *gin.Context) {
	req := new(model.IncomingApplyReviewReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("incomingApplyReview ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	username, _ := c.Get(CtxKeyUsername)
	reviewer, _ := username.(string)
	if err := svc.IncomingApplyReview(c, req, reviewer); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func incomingRecordList(c *gin.Context) {
	req := new(model.IncomingRecordListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("incomingRecordList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.IncomingRecordList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func incomingRecordDetail(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("incomingRecordDetail ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	detail, err := svc.IncomingRecordDetail(c, req.ID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, detail, nil)
}
