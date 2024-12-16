package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
)

func goodsGetGoodsList(c *gin.Context) {
	rsp, err := svc.GetGoodsList(c)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}
