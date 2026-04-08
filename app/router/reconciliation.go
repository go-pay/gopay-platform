package router

import (
	"gopay/app/model"
	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

// ---- 对账单 ----

func reconBillList(c *gin.Context) {
	req := new(model.ReconBillListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("reconBillList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.ReconBillList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func reconBillDetail(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("reconBillDetail ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	detail, err := svc.ReconBillDetail(c, req.ID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, detail, nil)
}

func reconBillGenerate(c *gin.Context) {
	req := new(model.ReconBillGenerateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("reconBillGenerate ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.ReconBillGenerate(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func reconBillReconcile(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("reconBillReconcile ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.ReconBillReconcile(c, req.ID); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

// ---- 对账差异 ----

func reconDiffList(c *gin.Context) {
	req := new(model.ReconDiffListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("reconDiffList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.ReconDiffList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func reconDiffDetail(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("reconDiffDetail ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	detail, err := svc.ReconDiffDetail(c, req.ID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, detail, nil)
}

func reconDiffHandle(c *gin.Context) {
	req := new(model.ReconDiffHandleReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("reconDiffHandle ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	username, _ := c.Get(CtxKeyUsername)
	handler, _ := username.(string)
	if err := svc.ReconDiffHandle(c, req, handler); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func reconDiffExport(c *gin.Context) {
	req := new(model.ReconDiffListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("reconDiffExport ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	list, err := svc.ReconDiffExport(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, list, nil)
}
