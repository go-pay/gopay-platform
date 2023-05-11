package router

import (
	"net/http/pprof"

	"gopay/app/service"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopher/web"
)

var srv *service.Service

func StartHttpServer(s *service.Service) (g *web.GinEngine) {
	srv = s
	g = web.InitGin(s.Cfg.Http)
	g.Gin.TrustedPlatform = "x_forwarded_for"
	g.Gin.Use(g.CORS())

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
}
