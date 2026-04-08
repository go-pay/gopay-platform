package router

import (
	"strings"

	"gopay/app/model"
	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

func userGetInfo(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	tokenStr := strings.TrimPrefix(auth, "Bearer ")
	if tokenStr == "" || tokenStr == auth {
		web.JSON(c, nil, errcode.TokenInvalid)
		return
	}
	rsp, err := svc.GetUserInfo(c, tokenStr)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func userChangePwd(c *gin.Context) {
	req := new(model.ChangePwdReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("userChangePwd ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	uid, _ := c.Get(CtxKeyUID)
	uidVal, _ := uid.(int64)
	if uidVal == 0 {
		web.JSON(c, nil, errcode.TokenInvalid)
		return
	}
	if err := svc.ChangePwd(c, uidVal, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}

func userProfile(c *gin.Context) {
	req := new(model.ProfileReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("userProfile ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	uid, _ := c.Get(CtxKeyUID)
	uidVal, _ := uid.(int64)
	if uidVal == 0 {
		web.JSON(c, nil, errcode.TokenInvalid)
		return
	}
	if err := svc.UpdateProfile(c, uidVal, req); err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, nil, nil)
}
