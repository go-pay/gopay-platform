package router

import (
	"gopay/app/model"
	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

func paymentChannelList(c *gin.Context) {
	req := new(model.ChannelListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("paymentChannelList ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.ChannelList(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func paymentChannelAdd(c *gin.Context) {
	req := new(model.ChannelAddReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("paymentChannelAdd ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.ChannelAdd(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func paymentChannelUpdate(c *gin.Context) {
	req := new(model.ChannelUpdateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("paymentChannelUpdate ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.ChannelUpdate(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func paymentChannelToggleStatus(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("paymentChannelToggleStatus ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.ChannelToggleStatus(c, req.ID); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func paymentChannelDetail(c *gin.Context) {
	req := new(model.IDReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("paymentChannelDetail ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	detail, err := svc.ChannelDetail(c, req.ID)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, detail, nil)
}

func paymentChannelConfig(c *gin.Context) {
	req := new(model.ChannelConfigReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("paymentChannelConfig ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	if err := svc.ChannelConfig(c, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}
