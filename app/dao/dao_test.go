package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"gopay/app/conf"
	"gopay/pkg/config"
)

var (
	d   *Dao
	ctx = context.Background()
)

func TestMain(m *testing.M) {
	os.Setenv("RUNTIME_ENV", "local")
	flag.Set("conf", "../cfg/config.yaml")
	if err := config.ParseYaml(conf.Conf); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
