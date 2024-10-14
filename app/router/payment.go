package router

import (
	"gopay/app/model"
	"gopay/ecode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

func alipayGetPaymentQrcode(c *gin.Context) {
	req := new(model.AlipayGetPaymentQrcodeReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		xlog.Errorf("c.ShouldBindJSON(%v), err:%v", req, err)
		web.JSON(c, nil, ecode.RequestErr)
		return
	}
	rsp, err := svc.AlipayGetPaymentQrCode(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func alipayPagePayUrl(c *gin.Context) {
	req := new(model.AlipayPagePayUrlReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		xlog.Errorf("c.ShouldBindJSON(%v), err:%v", req, err)
		web.JSON(c, nil, ecode.RequestErr)
		return
	}
	rsp, err := svc.AlipayPagePayUrl(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}
