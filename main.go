package main

import (
	"context"

	"gopay/app/conf"
	"gopay/app/router"
	"gopay/app/service"
	"gopay/pkg/config"
)

func main() {
	// Parse Config
	err := config.ParseYaml(conf.Conf)
	if err != nil {
		panic(err)
	}
	config.SetLogLevel(conf.Conf.Cfg.LogLevel)

	// New Service
	svc := service.New(conf.Conf)
	// Start Web Server
	router.NewHttpServer(svc).AddExitHook(func(c context.Context) {
		svc.Close()
	}).Start()
}
