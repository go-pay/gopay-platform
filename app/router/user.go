package router

import (
	"strings"

	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
)

func userGetInfo(c *gin.Context) {
	// 从 Authorization header 提取 token
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
