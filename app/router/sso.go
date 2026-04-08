package router

import (
	"gopay/app/model"
	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

func ssoLogin(c *gin.Context) {
	req := new(model.LoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("ssoLogin ShouldBindJSON(%v), err:%v", req, err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.Login(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}
