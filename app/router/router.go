package router

import (
	"net/http/pprof"

	"gopay/app/service"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopher/ecode"
	"github.com/go-pay/gopher/web"
)

var svc *service.Service

func StartHttpServer(s *service.Service) (g *web.GinEngine) {
	svc = s
	g = web.InitGin(s.Cfg.Http)
	g.Gin.TrustedPlatform = "x_forwarded_for"
	g.Gin.Use(g.CORS())

	ecode.Success = ecode.NewV2(0, "SUCCESS", "成功")
	initRoute(g.Gin)
	g.Start()
	return g
}

func initRoute(g *gin.Engine) {
	pp := g.Group("/debug/pprof")
	{
		pp.GET("/index", func(c *gin.Context) { pprof.Index(c.Writer, c.Request) })
		pp.GET("/cmdline", func(c *gin.Context) { pprof.Cmdline(c.Writer, c.Request) })
		pp.GET("/profile", func(c *gin.Context) { pprof.Profile(c.Writer, c.Request) })
		pp.GET("/symbol", func(c *gin.Context) { pprof.Symbol(c.Writer, c.Request) })
		pp.GET("/trace", func(c *gin.Context) { pprof.Trace(c.Writer, c.Request) })
	}

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
