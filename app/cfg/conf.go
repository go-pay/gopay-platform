package cfg

import (
	"github.com/go-pay/gopher/orm"
	"github.com/go-pay/gopher/web"
	"github.com/go-pay/gopher/xtime"
)

var Conf = &Config{}

type Config struct {
	Cfg         *Cfg        `yaml:"cfg"`
	Http        *web.Config `yaml:"http"`
	Redis       *Redis      `yaml:"redis"`
	MySQL       *MySQL      `yaml:"mysql"`
	PayPlatform *Payment    `yaml:"pay_platform"`
}

type Cfg struct {
	LogLevel           string         `yaml:"log_level"`
	ReloadInterval     xtime.Duration `yaml:"reload_interval"`
	ReloadLongInterval xtime.Duration `yaml:"reload_long_interval"`
}

type Redis struct {
	Gopay *orm.RedisConfig `yaml:"gopay"`
}

type MySQL struct {
	Gopay *orm.MySQLConfig `yaml:"gopay"`
}

type Payment struct {
	Wechat *WechatPay `yaml:"wechat"`
	Alipay *Alipay    `yaml:"alipay"`
}

type WechatPay struct {
	Appid             string `yaml:"appid"`
	MchId             string `yaml:"mch_id"`
	ApiKey            string `yaml:"api_key"`
	ApiV3Key          string `yaml:"api_v3_key"`
	CertFileContent   string `yaml:"cert_file_content"`
	KeyFileContent    string `yaml:"key_file_content"`
	Pkcs12FileContent string `yaml:"pkcs12_file_content"`
}

type Alipay struct {
	Appid             string `yaml:"appid"`
	PrivateKey        string `yaml:"private_key"`
	AppCertContent    string `yaml:"app_cert_content"`
	RootCertContent   string `yaml:"root_cert_content"`
	PublicCertContent string `yaml:"public_cert_content"`
}