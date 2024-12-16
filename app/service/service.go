package service

import (
	"context"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay/cert"
	"sync"

	"gopay/app/conf"
	"gopay/app/dao"

	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/go-pay/smap"
)

type Service struct {
	rwMu   sync.RWMutex
	Config *conf.Config
	dao    *dao.Dao
	alipay *alipay.Client
	wxpay  *wechat.ClientV3

	// cache
	kvMap smap.Map[string, string] // key: k, value: v
}

var (
	srv *Service
	ctx = context.Background()
)

func New(c *conf.Config) (s *Service) {
	srv = &Service{
		Config: c,
		//dao:    dao.New(c),
	}
	// 初始化支付宝 client
	if c.PayPlatform.Alipay.Appid != "" && c.PayPlatform.Alipay.PrivateKey != "" {
		alipayCli, err := alipay.NewClient(c.PayPlatform.Alipay.Appid, c.PayPlatform.Alipay.PrivateKey, false)
		if err != nil {
			panic(err)
		}
		// Debug开关，输出/关闭日志
		alipayCli.DebugSwitch = gopay.DebugOff
		// 配置公共参数
		alipayCli.SetCharset(alipay.UTF8).
			SetSignType(alipay.RSA2).
			// SetAppAuthToken("")
			SetReturnUrl("https://baidu.com").
			SetNotifyUrl("https://baidu.com")

		// 自动同步验签（只支持证书模式）
		// 传入 支付宝公钥证书 alipayPublicCert.crt 内容
		alipayCli.AutoVerifySign(cert.AlipayPublicContentRSA2)
		// 传入证书内容
		if err = alipayCli.SetCertSnByContent([]byte(c.PayPlatform.Alipay.AppPublicCertContent),
			[]byte(c.PayPlatform.Alipay.AlipayRootCertContent),
			[]byte(c.PayPlatform.Alipay.AlipayPublicCertContent)); err != nil {
			panic(err)
		}
		srv.alipay = alipayCli
	}

	// 初始化微信v3 client
	wxCli, err := wechat.NewClientV3(c.PayPlatform.Wechat.MchId, c.PayPlatform.Wechat.SerialNo, c.PayPlatform.Wechat.ApiV3Key, c.PayPlatform.Wechat.KeyFileContent)
	if err != nil {
		panic(err)
	}
	wxCli.DebugSwitch = gopay.DebugOn
	srv.wxpay = wxCli

	// loop job
	srv.initLoop()
	return srv
}

// Close 关闭相关资源
func (s *Service) Close() {
	if s.dao != nil {
		s.dao.Close()
	}
}

// 初始化 loop
func (s *Service) initLoop() {
	// 初始化 loop
}
