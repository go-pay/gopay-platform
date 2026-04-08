# GoPay Platform 后端实现进度

> 最后更新: 2026-04-08 | 总接口: 67 个 | 已完成: 2 个 | 进度: 3%

---

## 快速上下文 (新会话必读)

**项目是什么**: Go 支付平台后端，前端 18 个页面已全部完成(mock 数据)，后端需实现 67 个接口对接。

**当前状态**: Phase 1 (基础设施) 尚未开始。已有的基础功能仅登录、获取用户信息和支付宝支付。

**下一步该做什么**: 从 Phase 1 开始 — 执行数据库 DDL、创建 dm 实体、完善中间件。

**开发前准备**:
1. 确保 MySQL 已启动，执行 `mysql_ddl.sql` 建库建表
2. 读 `docs/api-detail.md` 了解接口的请求/响应格式
3. 读下方的代码模式参考，了解本项目的编码风格
4. 按 Phase 顺序开发，每完成一个接口更新本文档状态

**决策记录**:
| 日期 | 决策 | 原因 |
|------|------|------|
| 2026-04-08 | 全部接口用 POST 方法 | 与前端现有调用方式对齐，简化 CORS 配置 |
| 2026-04-08 | 金额统一用分(int64) | 避免浮点精度问题 |
| 2026-04-08 | Phase 顺序: 基础设施→系统管理→商户→通道→订单→交易→对账→仪表盘 | 仪表盘依赖订单和交易数据，放最后；系统管理是其他模块的权限基础 |

---

## 进度总览

| 阶段 | 模块 | 接口数 | 状态 | 说明 |
|------|------|-------|------|------|
| Phase 1 | 基础设施 | - | ⏳ 待开始 | DDL、dm 实体、中间件、分页工具 |
| Phase 2 | 认证与系统管理 | 17 | ⏳ 待开始 | 登录扩展、用户/角色/日志 CRUD |
| Phase 3 | 商户与进件 | 12 | ⏳ 待开始 | 商户/应用/进件 CRUD |
| Phase 4 | 支付通道配置 | 6 | ⏳ 待开始 | 通道 CRUD + 参数配置 |
| Phase 5 | 订单中心 | 10 | ⏳ 待开始 | 支付/退款/转账订单 |
| Phase 6 | 交易记录 | 6 | ⏳ 待开始 | 交易流水 + 回调通知 |
| Phase 7 | 对账管理 | 8 | ⏳ 待开始 | 对账单 + 对账差异 |
| Phase 8 | 仪表盘 | 4 | ⏳ 待开始 | 统计数据聚合 |

状态标识: ✅ 已完成 | 🔧 开发中 | ⏳ 待开始 | 🚫 阻塞

---

## Phase 1: 基础设施

| # | 任务 | 涉及文件 | 状态 | 备注 |
|---|------|---------|------|------|
| 1.1 | 执行数据库 DDL - 新建 11 张表 | mysql_ddl.sql | ⏳ | merchant, merchant_app, payment_channel, payment_channel_config, incoming_apply, refund_order, transfer_order, transaction_flow, callback_record, reconciliation_bill, reconciliation_diff, sys_role, sys_role_perm, operation_log |
| 1.2 | ALTER 现有表 - account 加字段 | mysql_ddl.sql | ⏳ | real_name, email, role, status, last_login |
| 1.3 | ALTER 现有表 - payment_order 加字段 | mysql_ddl.sql | ⏳ | order_no, out_trade_no, merchant_id, channel_type, pay_method, subject, client_ip, notified |
| 1.4 | 新增 dm 层实体 | app/dm/*.go | ⏳ | 所有新表的 GORM Model |
| 1.5 | 扩展 Account dm 实体 | app/dm/db_gopay.go | ⏳ | 新增字段 |
| 1.6 | 新增错误码 | errcode/ecode.go | ⏳ | NotFound(10404), Forbidden(10405), Conflict(10409) |
| 1.7 | JWT 认证中间件 (从 service 提取) | app/router/middleware.go | ⏳ | |
| 1.8 | 操作日志中间件 | app/router/middleware.go | ⏳ | |
| 1.9 | 统一分页 DTO | app/model/page.go | ⏳ | PageReq, PageResp |

---

## Phase 2: 认证与系统管理

### 2A - 认证扩展 (前端: login, profile)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 2A.1 | 登录(扩展) | POST /gopay/v1/sso/login | 🔧 | 已有基础实现，需加 realName/email/lastLogin |
| 2A.2 | 获取用户信息(扩展) | POST /gopay/v1/user/getInfo | 🔧 | 同上 |
| 2A.3 | 修改密码 | POST /gopay/v1/user/changePwd | ⏳ | |
| 2A.4 | 更新个人资料 | POST /gopay/v1/user/profile | ⏳ | |

### 2B - 系统用户管理 (前端: system/user.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 2B.1 | 用户列表 | POST /gopay/v1/system/user/list | ⏳ | |
| 2B.2 | 新增用户 | POST /gopay/v1/system/user/add | ⏳ | |
| 2B.3 | 编辑用户 | POST /gopay/v1/system/user/update | ⏳ | |
| 2B.4 | 切换用户状态 | POST /gopay/v1/system/user/toggleStatus | ⏳ | |
| 2B.5 | 重置密码 | POST /gopay/v1/system/user/resetPwd | ⏳ | |

### 2C - 角色管理 (前端: system/role.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 2C.1 | 角色列表 | POST /gopay/v1/system/role/list | ⏳ | |
| 2C.2 | 新增角色 | POST /gopay/v1/system/role/add | ⏳ | |
| 2C.3 | 编辑角色 | POST /gopay/v1/system/role/update | ⏳ | |
| 2C.4 | 切换角色状态 | POST /gopay/v1/system/role/toggleStatus | ⏳ | |
| 2C.5 | 更新角色权限 | POST /gopay/v1/system/role/perms/update | ⏳ | |
| 2C.6 | 查询角色权限 | POST /gopay/v1/system/role/perms/list | ⏳ | |

### 2D - 操作日志 (前端: system/log.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 2D.1 | 日志列表 | POST /gopay/v1/system/log/list | ⏳ | |
| 2D.2 | 日志详情 | POST /gopay/v1/system/log/detail | ⏳ | |
| 2D.3 | 导出日志 | POST /gopay/v1/system/log/export | ⏳ | |

---

## Phase 3: 商户与进件

### 3A - 商户管理 (前端: merchant/list.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 3A.1 | 商户列表 | POST /gopay/v1/merchant/list | ⏳ | |
| 3A.2 | 新增商户 | POST /gopay/v1/merchant/add | ⏳ | |
| 3A.3 | 编辑商户 | POST /gopay/v1/merchant/update | ⏳ | |
| 3A.4 | 切换商户状态 | POST /gopay/v1/merchant/toggleStatus | ⏳ | |
| 3A.5 | 商户下拉选项 | POST /gopay/v1/merchant/options | ⏳ | |

### 3B - 商户应用 (前端: merchant/app.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 3B.1 | 应用列表 | POST /gopay/v1/merchant/app/list | ⏳ | |
| 3B.2 | 新增应用 | POST /gopay/v1/merchant/app/add | ⏳ | |
| 3B.3 | 编辑应用 | POST /gopay/v1/merchant/app/update | ⏳ | |

### 3C - 进件管理 (前端: incoming/apply.vue, record.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 3C.1 | 进件申请列表 | POST /gopay/v1/incoming/apply/list | ⏳ | |
| 3C.2 | 新建进件 | POST /gopay/v1/incoming/apply/add | ⏳ | |
| 3C.3 | 提交审核 | POST /gopay/v1/incoming/apply/submit | ⏳ | |
| 3C.4 | 审核 | POST /gopay/v1/incoming/apply/review | ⏳ | |
| 3C.5 | 进件记录列表 | POST /gopay/v1/incoming/record/list | ⏳ | |
| 3C.6 | 进件记录详情 | POST /gopay/v1/incoming/record/detail | ⏳ | |
| 3C.7 | 图片上传 | POST /gopay/v1/upload/image | ⏳ | |

---

## Phase 4: 支付通道配置

### 4A - 支付通道 (前端: payment/channel.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 4A.1 | 通道列表 | POST /gopay/v1/payment/channel/list | ⏳ | |
| 4A.2 | 新增通道 | POST /gopay/v1/payment/channel/add | ⏳ | |
| 4A.3 | 编辑通道 | POST /gopay/v1/payment/channel/update | ⏳ | |
| 4A.4 | 切换通道状态 | POST /gopay/v1/payment/channel/toggleStatus | ⏳ | |
| 4A.5 | 通道详情 | POST /gopay/v1/payment/channel/detail | ⏳ | |
| 4A.6 | 参数配置 | POST /gopay/v1/payment/channel/config | ⏳ | |

---

## Phase 5: 订单中心

### 5A - 支付订单 (前端: order/payment.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 5A.1 | 支付订单列表 | POST /gopay/v1/order/payment/list | ⏳ | |
| 5A.2 | 支付订单详情 | POST /gopay/v1/order/payment/detail | ⏳ | |
| 5A.3 | 关闭订单 | POST /gopay/v1/order/payment/close | ⏳ | |
| 5A.4 | 发起退款 | POST /gopay/v1/order/payment/refund | ⏳ | |
| 5A.5 | 导出支付订单 | POST /gopay/v1/order/payment/export | ⏳ | |

### 5B - 退款订单 (前端: order/refund.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 5B.1 | 退款订单列表 | POST /gopay/v1/order/refund/list | ⏳ | |
| 5B.2 | 退款详情 | POST /gopay/v1/order/refund/detail | ⏳ | |

### 5C - 转账订单 (前端: order/transfer.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 5C.1 | 转账订单列表 | POST /gopay/v1/order/transfer/list | ⏳ | |
| 5C.2 | 发起转账 | POST /gopay/v1/order/transfer/add | ⏳ | |
| 5C.3 | 转账详情 | POST /gopay/v1/order/transfer/detail | ⏳ | |

---

## Phase 6: 交易记录

### 6A - 交易流水 (前端: transaction/flow.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 6A.1 | 流水列表 | POST /gopay/v1/transaction/flow/list | ⏳ | |
| 6A.2 | 流水详情 | POST /gopay/v1/transaction/flow/detail | ⏳ | |
| 6A.3 | 流水统计 | POST /gopay/v1/transaction/flow/stats | ⏳ | |

### 6B - 回调通知 (前端: transaction/callback.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 6B.1 | 回调列表 | POST /gopay/v1/transaction/callback/list | ⏳ | |
| 6B.2 | 回调详情 | POST /gopay/v1/transaction/callback/detail | ⏳ | |
| 6B.3 | 手动重试 | POST /gopay/v1/transaction/callback/retry | ⏳ | |

---

## Phase 7: 对账管理

### 7A - 对账单 (前端: reconciliation/bill.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 7A.1 | 对账单列表 | POST /gopay/v1/recon/bill/list | ⏳ | |
| 7A.2 | 对账单详情 | POST /gopay/v1/recon/bill/detail | ⏳ | |
| 7A.3 | 生成对账单 | POST /gopay/v1/recon/bill/generate | ⏳ | |
| 7A.4 | 执行对账 | POST /gopay/v1/recon/bill/reconcile | ⏳ | |

### 7B - 对账差异 (前端: reconciliation/diff.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 7B.1 | 差异列表 | POST /gopay/v1/recon/diff/list | ⏳ | |
| 7B.2 | 差异详情 | POST /gopay/v1/recon/diff/detail | ⏳ | |
| 7B.3 | 处理差异 | POST /gopay/v1/recon/diff/handle | ⏳ | |
| 7B.4 | 导出差异 | POST /gopay/v1/recon/diff/export | ⏳ | |

---

## Phase 8: 仪表盘

### 8A - Dashboard (前端: dashboard/index.vue)

| # | 接口 | 路由 | 状态 | 备注 |
|---|------|------|------|------|
| 8A.1 | 统计数据 | POST /gopay/v1/dashboard/stats | ⏳ | |
| 8A.2 | 最近订单 | POST /gopay/v1/dashboard/recentOrders | ⏳ | |
| 8A.3 | 通道分布 | POST /gopay/v1/dashboard/channelDistribution | ⏳ | |
| 8A.4 | 近7天趋势 | POST /gopay/v1/dashboard/trend | ⏳ | |

---

## 已完成的功能 (项目初始状态)

| 功能 | 文件 | 说明 |
|------|------|------|
| ✅ 登录接口 (基础) | app/router/sso.go, app/service/account.go | JWT 登录，缺少新字段 |
| ✅ 获取用户信息 (基础) | app/router/user.go, app/service/account.go | 缺少新字段 |
| ✅ 支付宝扫码支付 | app/router/payment.go, app/service/payment.go | TradePrecreate |
| ✅ 支付宝网页支付 | app/router/payment.go, app/service/payment.go | TradePagePay |
| ✅ 健康检查 | app/router/router.go | GET /monitor/ping |
| ✅ 数据库连接 | app/dao/dao.go | GORM MySQL |
| ✅ 配置加载 | app/conf/ | YAML 配置 |
| ✅ 错误码框架 | errcode/ecode.go | 基础错误码 |
| ✅ CORS 中间件 | app/router/router.go | 跨域支持 |

---

## 变更日志

| 日期 | 变更 |
|------|------|
| 2026-04-08 | 初始化文档，分析前后端状态，制定 8 阶段实施计划 |
