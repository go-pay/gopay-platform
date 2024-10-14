package router

import (
	"gopay/app/service"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/ecode"
	"github.com/go-pay/web"
	"github.com/go-pay/web/middleware"
)

var svc *service.Service

func NewHttpServer(s *service.Service) (g *web.GinEngine) {
	svc = s
	g = web.InitGin(s.Cfg.Http)
	g.Gin.TrustedPlatform = "x_forwarded_for"
	g.Gin.Use(middleware.CORS())

	ecode.Success = ecode.New(0, "SUCCESS", "success")
	initRoute(g.Gin)
	return g
}

func initRoute(g *gin.Engine) {
	monitor := g.Group("/gopay/v1/monitor")
	{
		monitor.GET("/ping", func(c *gin.Context) { web.JSON(c, "PingOK: "+c.ClientIP(), nil) })
	}
	v1 := g.Group("/gopay/v1")
	{
		// sso
		sso := v1.Group("/sso")
		{
			sso.POST("/login", ssoLogin) // 登录
		}
		// 用户相关
		user := v1.Group("/user")
		{
			user.POST("/getInfo", userGetInfo) // 获取用户信息
		}
		// 支付相关
		pay := v1.Group("/payment")
		{
			manage := pay.Group("/manage")
			{
				manage.POST("/getPaymentInfoList") // 获取支付信息列表
				manage.POST("/addPaymentInfo")     // 添加支付配置信息
			}
			// 支付宝支付
			ali := pay.Group("/alipay")
			{
				ali.POST("/getPaymentQrcode", alipayGetPaymentQrcode) // 获取支付宝支付二维码
				ali.POST("/getPagePayUrl", alipayPagePayUrl)          // 获取支付宝网页支付链接
			}
			// 微信支付
			wx := pay.Group("/wechat")
			{
				wx.POST("/getPaymentQrcode") // 获取微信支付二维码
			}
		}
		// 订单相关
		order := v1.Group("/order")
		{
			order.POST("/getOrderList") // 获取支付订单列表
		}
	}
}
