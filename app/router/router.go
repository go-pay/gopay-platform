package router

import (
	"gopay/app/service"
	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/ecode"
	"github.com/go-pay/web"
	"github.com/go-pay/web/middleware"
)

var svc *service.Service

func NewHttpServer(s *service.Service) (g *web.GinEngine) {
	svc = s
	g = web.InitGin(s.Cfg.Http)
	g.Gin.TrustedPlatform = "x_forwarded_for"
	g.Gin.Use(middleware.CORS())

	ecode.Success = errcode.Success
	initRoute(g.Gin)
	return g
}

func initRoute(g *gin.Engine) {
	// 静态文件服务 - 上传文件访问
	g.Static("/uploads", "./uploads")

	monitor := g.Group("/gopay/v1/monitor")
	{
		monitor.GET("/ping", func(c *gin.Context) { web.JSON(c, "PingOK: "+c.ClientIP(), nil) })
	}

	v1 := g.Group("/gopay/v1")
	{
		// 无需认证的接口
		sso := v1.Group("/sso")
		{
			sso.POST("/login", ssoLogin)
		}
	}

	// 需要 JWT 认证的接口
	auth := g.Group("/gopay/v1", JWTAuth())
	{
		// 用户相关
		user := auth.Group("/user")
		{
			user.POST("/getInfo", userGetInfo)
			user.POST("/changePwd", userChangePwd)
			user.POST("/profile", userProfile)
		}

		// 仪表盘
		dashboard := auth.Group("/dashboard")
		{
			dashboard.POST("/stats", dashboardStats)
			dashboard.POST("/recentOrders", dashboardRecentOrders)
			dashboard.POST("/channelDistribution", dashboardChannelDistribution)
			dashboard.POST("/trend", dashboardTrend)
		}

		// 商户管理
		merchant := auth.Group("/merchant")
		{
			merchant.POST("/list", merchantList)
			merchant.POST("/add", OperationLogger("merchant", "create", "新增商户"), merchantAdd)
			merchant.POST("/update", OperationLogger("merchant", "update", "编辑商户"), merchantUpdate)
			merchant.POST("/toggleStatus", OperationLogger("merchant", "update", "切换商户状态"), merchantToggleStatus)
			merchant.POST("/options", merchantOptions)

			app := merchant.Group("/app")
			{
				app.POST("/list", merchantAppList)
				app.POST("/add", OperationLogger("merchant", "create", "新增应用"), merchantAppAdd)
				app.POST("/update", OperationLogger("merchant", "update", "编辑应用"), merchantAppUpdate)
			}
		}

		// 进件管理
		incoming := auth.Group("/incoming")
		{
			apply := incoming.Group("/apply")
			{
				apply.POST("/list", incomingApplyList)
				apply.POST("/add", OperationLogger("incoming", "create", "新建进件"), incomingApplyAdd)
				apply.POST("/submit", OperationLogger("incoming", "update", "提交审核"), incomingApplySubmit)
				apply.POST("/review", OperationLogger("incoming", "update", "审核进件"), incomingApplyReview)
			}
			record := incoming.Group("/record")
			{
				record.POST("/list", incomingRecordList)
				record.POST("/detail", incomingRecordDetail)
			}
		}

		// 文件上传
		auth.POST("/upload/image", uploadImage)

		// 支付通道配置
		channel := auth.Group("/payment/channel")
		{
			channel.POST("/list", paymentChannelList)
			channel.POST("/add", OperationLogger("payment", "create", "新增通道"), paymentChannelAdd)
			channel.POST("/update", OperationLogger("payment", "update", "编辑通道"), paymentChannelUpdate)
			channel.POST("/toggleStatus", OperationLogger("payment", "update", "切换通道状态"), paymentChannelToggleStatus)
			channel.POST("/detail", paymentChannelDetail)
			channel.POST("/config", OperationLogger("payment", "update", "配置通道参数"), paymentChannelConfig)
		}

		// 支付宝支付 (保留已有接口)
		ali := auth.Group("/payment/alipay")
		{
			ali.POST("/getPaymentQrcode", alipayGetPaymentQrcode)
			ali.POST("/getPagePayUrl", alipayPagePayUrl)
		}

		// 订单中心
		order := auth.Group("/order")
		{
			payOrder := order.Group("/payment")
			{
				payOrder.POST("/list", orderPaymentList)
				payOrder.POST("/detail", orderPaymentDetail)
				payOrder.POST("/close", OperationLogger("order", "update", "关闭订单"), orderPaymentClose)
				payOrder.POST("/refund", OperationLogger("order", "create", "发起退款"), orderPaymentRefund)
				payOrder.POST("/export", OperationLogger("order", "export", "导出支付订单"), orderPaymentExport)
			}
			refund := order.Group("/refund")
			{
				refund.POST("/list", orderRefundList)
				refund.POST("/detail", orderRefundDetail)
			}
			transfer := order.Group("/transfer")
			{
				transfer.POST("/list", orderTransferList)
				transfer.POST("/add", OperationLogger("order", "create", "发起转账"), orderTransferAdd)
				transfer.POST("/detail", orderTransferDetail)
			}
		}

		// 交易记录
		transaction := auth.Group("/transaction")
		{
			flow := transaction.Group("/flow")
			{
				flow.POST("/list", transactionFlowList)
				flow.POST("/detail", transactionFlowDetail)
				flow.POST("/stats", transactionFlowStats)
			}
			callback := transaction.Group("/callback")
			{
				callback.POST("/list", callbackList)
				callback.POST("/detail", callbackDetail)
				callback.POST("/retry", OperationLogger("transaction", "update", "手动重试回调"), callbackRetry)
			}
		}

		// 对账管理
		recon := auth.Group("/recon")
		{
			bill := recon.Group("/bill")
			{
				bill.POST("/list", reconBillList)
				bill.POST("/detail", reconBillDetail)
				bill.POST("/generate", OperationLogger("recon", "create", "生成对账单"), reconBillGenerate)
				bill.POST("/reconcile", OperationLogger("recon", "update", "执行对账"), reconBillReconcile)
			}
			diff := recon.Group("/diff")
			{
				diff.POST("/list", reconDiffList)
				diff.POST("/detail", reconDiffDetail)
				diff.POST("/handle", OperationLogger("recon", "update", "处理对账差异"), reconDiffHandle)
				diff.POST("/export", OperationLogger("recon", "export", "导出对账差异"), reconDiffExport)
			}
		}

		// 系统管理
		system := auth.Group("/system")
		{
			sysUser := system.Group("/user")
			{
				sysUser.POST("/list", systemUserList)
				sysUser.POST("/add", OperationLogger("system", "create", "新增用户"), systemUserAdd)
				sysUser.POST("/update", OperationLogger("system", "update", "编辑用户"), systemUserUpdate)
				sysUser.POST("/toggleStatus", OperationLogger("system", "update", "切换用户状态"), systemUserToggleStatus)
				sysUser.POST("/resetPwd", OperationLogger("system", "update", "重置密码"), systemUserResetPwd)
			}
			role := system.Group("/role")
			{
				role.POST("/list", systemRoleList)
				role.POST("/add", OperationLogger("system", "create", "新增角色"), systemRoleAdd)
				role.POST("/update", OperationLogger("system", "update", "编辑角色"), systemRoleUpdate)
				role.POST("/toggleStatus", OperationLogger("system", "update", "切换角色状态"), systemRoleToggleStatus)
				perms := role.Group("/perms")
				{
					perms.POST("/update", OperationLogger("system", "update", "更新角色权限"), systemRolePermsUpdate)
					perms.POST("/list", systemRolePermsList)
				}
			}
			sysLog := system.Group("/log")
			{
				sysLog.POST("/list", systemLogList)
				sysLog.POST("/detail", systemLogDetail)
				sysLog.POST("/export", OperationLogger("system", "export", "导出操作日志"), systemLogExport)
			}
		}
	}
}
