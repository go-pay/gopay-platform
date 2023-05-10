package router

import (
	"github.com/gin-gonic/gin"
	"gopay/app/service"

	"github.com/go-pay/gopher/ecode"
	"github.com/go-pay/gopher/web"
)

var srv *service.Service

func StartHttpServer(s *service.Service) (g *web.GinEngine) {
	srv = s
	g = web.InitGin(s.Cfg.Http)
	g.Gin.TrustedPlatform = "x_forwarded_for"
	g.Gin.Use(g.CORS())

	initRoute(g.Gin)
	ecode.Success = ecode.NewV2(1, "SUCCESS", "success")
	g.Start()
	return g
}

func initRoute(g *gin.Engine) {

}
