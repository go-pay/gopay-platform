package conf

import (
	"github.com/go-pay/orm"
	"github.com/go-pay/web"
	"github.com/go-pay/xtime"
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
	Appid          string `yaml:"appid"`
	MchId          string `yaml:"mch_id"`
	SerialNo       string `yaml:"serial_no"`
	ApiV3Key       string `yaml:"api_v3_key"`
	KeyFileContent string `yaml:"key_file_content"`
}

type Alipay struct {
	Appid                   string `yaml:"appid"`
	PrivateKey              string `yaml:"private_key"`
	AppPublicCertContent    string `yaml:"app_public_cert_content"`
	AlipayRootCertContent   string `yaml:"alipay_root_cert_content"`
	AlipayPublicCertContent string `yaml:"alipay_public_cert_content"`
}
