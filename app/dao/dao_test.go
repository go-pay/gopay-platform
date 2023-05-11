package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"gopay/app/cfg"

	"github.com/go-pay/gopher/conf"
)

var (
	d   *Dao
	ctx = context.Background()
)

func TestMain(m *testing.M) {
	os.Setenv("RUNTIME_ENV", "local")
	flag.Set("conf", "../cfg/config.yaml")
	if err := conf.ParseYaml(cfg.Conf); err != nil {
		panic(err)
	}
	d = New(cfg.Conf)
	os.Exit(m.Run())
}
