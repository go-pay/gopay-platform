package router

import (
	"gopay/app/model"
	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

// ---- 商户管理 ----

func merchantList(c *gin.Context) {
	req := new(model.MerchantListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("merchantList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.MerchantList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func merchantAdd(c *gin.Context) {
	req := new(model.MerchantAddReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("merchantAdd ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	id, err := svc.MerchantAdd(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, gin.H{"id": id}, nil)
}

func merchantUpdate(c *gin.Context) {
	req := new(model.MerchantUpdateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("merchantUpdate ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.MerchantUpdate(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func merchantToggleStatus(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("merchantToggleStatus ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.MerchantToggleStatus(c, req.ID); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func merchantOptions(c *gin.Context) {
	opts, err := svc.MerchantOptions(c)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, opts, nil)
}

// ---- 商户应用 ----

func merchantAppList(c *gin.Context) {
	req := new(model.MerchantAppListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("merchantAppList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.MerchantAppList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func merchantAppAdd(c *gin.Context) {
	req := new(model.MerchantAppAddReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("merchantAppAdd ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.MerchantAppAdd(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func merchantAppUpdate(c *gin.Context) {
	req := new(model.MerchantAppUpdateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("merchantAppUpdate ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.MerchantAppUpdate(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}
