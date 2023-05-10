package conf

import "github.com/go-pay/gopher/web"

var Conf = &Config{}

type Config struct {
	Http *web.Config `yaml:"http"`
}
