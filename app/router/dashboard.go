package router

import (
	"gopay/app/model"
	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

func dashboardStats(c *gin.Context) {
	stats, err := svc.DashboardStats(c)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, stats, nil)
}

func dashboardRecentOrders(c *gin.Context) {
	req := new(model.PageReq)
	if err := c.ShouldBindJSON(req); err != nil {
		xlog.Errorf("dashboardRecentOrders ShouldBindJSON, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	rsp, err := svc.DashboardRecentOrders(c, req)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}

func dashboardChannelDistribution(c *gin.Context) {
	dist, err := svc.DashboardChannelDistribution(c)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, dist, nil)
}

func dashboardTrend(c *gin.Context) {
	trend, err := svc.DashboardTrend(c)
	if err != nil {
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, trend, nil)
}
