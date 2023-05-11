package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopay/app/cfg"
	"gopay/app/router"
	"gopay/app/service"

	"github.com/go-pay/gopher/conf"
	"github.com/go-pay/gopher/xlog"
)

func main() {
	xlog.Warnf("welcome to gopay platform")
	// Parse Config
	err := conf.ParseYaml(cfg.Conf)
	if err != nil {
		panic(err)
	}
	// New Service
	svc := service.New(cfg.Conf)
	// Start Web Server
	httpServer := router.StartHttpServer(svc)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			xlog.Warnf("get a signal %s, stop the process", si.String())
			httpServer.Close()
			// wait for a second
			time.Sleep(time.Second)
			svc.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
