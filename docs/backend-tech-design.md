# GoPay Platform 后端技术设计文档

> 版本: v1.0 | 日期: 2026-04-08 | 对齐前端: gopay-platform-web

---

## 一、系统概览

### 1.1 架构

```
前端 (Vue3 + Vuetify3, :3000)
    ↓ Vite proxy /gopay/v1/ → :2233
后端 (Gin + GORM, :2233)
    ├── router/   → HTTP 处理器，参数校验，响应封装
    ├── service/  → 业务逻辑，支付编排
    ├── dao/      → 数据访问 (MySQL + Redis)
    ├── model/    → 请求/响应 DTO
    └── dm/       → 数据库实体 (GORM Model)
```

### 1.2 统一响应格式

```json
{
  "code": 0,
  "msg": "success",
  "data": { ... }
}
```

错误码规范 (errcode/ecode.go):
| 码值 | 名称 | 说明 |
|------|------|------|
| 0 | Success | 成功 |
| 10400 | RequestErr | 参数错误 |
| 10401 | HeaderVerifyFailed | Header 校验失败 |
| 10402 | LoginFailed | 用户名或密码错误 |
| 10403 | TokenInvalid | Token 无效/过期 |
| 10404 | NotFound | 资源不存在 |
| 10405 | Forbidden | 无权限 |
| 10409 | Conflict | 数据冲突(如重名) |
| 10500 | ServerErr | 服务器错误 |
| 10501 | UnAvailableErr | 服务不可用 |

### 1.3 分页约定

请求参数:
```json
{ "page": 1, "pageSize": 20, "...其他筛选条件" }
```

响应:
```json
{
  "code": 0,
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "pageSize": 20
  }
}
```

---

## 二、数据库设计

### 2.1 现有表修改

#### account (扩展为系统用户表)

```sql
ALTER TABLE `account`
  ADD COLUMN `real_name`   VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '真实姓名' AFTER `pwd`,
  ADD COLUMN `email`       VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '邮箱' AFTER `phone`,
  ADD COLUMN `role`        VARCHAR(16)  NOT NULL DEFAULT 'viewer' COMMENT '角色: admin/operator/finance/viewer' AFTER `email`,
  ADD COLUMN `status`      TINYINT      NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用 1-正常' AFTER `role`,
  ADD COLUMN `last_login`  DATETIME     NULL COMMENT '最后登录时间' AFTER `status`;
```

#### payment_order (扩展)

```sql
ALTER TABLE `payment_order`
  ADD COLUMN `order_no`      VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '平台订单号(展示用)' AFTER `id`,
  ADD COLUMN `out_trade_no`  VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商户订单号' AFTER `trade_no`,
  ADD COLUMN `merchant_id`   BIGINT       NOT NULL DEFAULT 0 COMMENT '商户ID' AFTER `user_id`,
  ADD COLUMN `merchant_name` VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商户名称(冗余)' AFTER `merchant_id`,
  ADD COLUMN `channel_type`  VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat' AFTER `payment_type`,
  ADD COLUMN `pay_method`    VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '支付方式: qrcode/page/wap/app/jsapi/miniapp' AFTER `channel_type`,
  ADD COLUMN `subject`       VARCHAR(128) NOT NULL DEFAULT '' COMMENT '商品描述' AFTER `pay_method`,
  ADD COLUMN `client_ip`     VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '客户端IP' AFTER `subject`,
  ADD COLUMN `notified`      TINYINT      NOT NULL DEFAULT 0 COMMENT '是否已回调通知: 0-否 1-是' AFTER `notify_body`,
  ADD KEY `idx_order_no` (`order_no`),
  ADD KEY `idx_merchant_id` (`merchant_id`);
```

### 2.2 新增表

#### merchant (商户表)

```sql
CREATE TABLE IF NOT EXISTS `merchant` (
  `id`      BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `name`    VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商户名称',
  `contact` VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '联系人',
  `phone`   VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '联系电话',
  `email`   VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '邮箱',
  `status`  TINYINT      NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用 1-正常',
  `remark`  VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
  `ctime`   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商户表';
```

#### merchant_app (商户应用表)

```sql
CREATE TABLE IF NOT EXISTS `merchant_app` (
  `id`            BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `name`          VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '应用名称',
  `appid`         VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '应用AppID',
  `merchant_id`   BIGINT       NOT NULL DEFAULT 0 COMMENT '所属商户ID',
  `platform_type` TINYINT      NOT NULL DEFAULT 0 COMMENT '平台类型: 0-微信移动 1-微信网站 2-微信公众号 3-微信小程序 5-支付宝网页/移动 6-支付宝小程序 7-支付宝生活号',
  `merchant_type` TINYINT      NOT NULL DEFAULT 0 COMMENT '商户类型: 0-商户 1-服务商',
  `notify_url`    VARCHAR(256) NOT NULL DEFAULT '' COMMENT '回调通知URL',
  `return_url`    VARCHAR(256) NOT NULL DEFAULT '' COMMENT '支付Return URL',
  `status`        TINYINT      NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用 1-正常',
  `ctime`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  KEY `idx_merchant_id` (`merchant_id`),
  KEY `idx_appid` (`appid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商户应用表';
```

#### payment_channel (支付通道表)

```sql
CREATE TABLE IF NOT EXISTS `payment_channel` (
  `id`            BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `name`          VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '通道名称',
  `code`          VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '通道编码',
  `type`          VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `merchant_id`   BIGINT       NOT NULL DEFAULT 0 COMMENT '所属商户ID',
  `pay_methods`   VARCHAR(128) NOT NULL DEFAULT '' COMMENT '支持的支付方式(逗号分隔): qrcode,page,wap,app,jsapi,miniapp',
  `fee_rate`      DECIMAL(5,2) NOT NULL DEFAULT 0.00 COMMENT '费率(%)',
  `status`        TINYINT      NOT NULL DEFAULT 1 COMMENT '状态: 0-停用 1-启用',
  `remark`        VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
  `ctime`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_merchant_id` (`merchant_id`),
  KEY `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付通道表';
```

#### payment_channel_config (支付通道参数配置表)

```sql
CREATE TABLE IF NOT EXISTS `payment_channel_config` (
  `id`           BIGINT        NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `channel_id`   BIGINT        NOT NULL DEFAULT 0 COMMENT '支付通道ID',
  `app_id`       VARCHAR(32)   NOT NULL DEFAULT '' COMMENT 'AppID',
  `mch_id`       VARCHAR(32)   NOT NULL DEFAULT '' COMMENT '商户号(微信)',
  `private_key`  TEXT          NULL COMMENT '私钥',
  `public_key`   TEXT          NULL COMMENT '公钥',
  `api_key`      VARCHAR(256)  NOT NULL DEFAULT '' COMMENT 'API密钥(微信V3)',
  `serial_no`    VARCHAR(64)   NOT NULL DEFAULT '' COMMENT '证书序列号(微信)',
  `notify_url`   VARCHAR(256)  NOT NULL DEFAULT '' COMMENT '回调通知URL',
  `sign_type`    VARCHAR(8)    NOT NULL DEFAULT 'RSA2' COMMENT '签名方式: RSA2/RSA',
  `sandbox`      TINYINT       NOT NULL DEFAULT 0 COMMENT '是否沙箱: 0-生产 1-沙箱',
  `ctime`        TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`        TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_channel_id` (`channel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付通道参数配置表';
```

#### incoming_apply (进件申请表)

```sql
CREATE TABLE IF NOT EXISTS `incoming_apply` (
  `id`           BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `apply_no`     VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '申请编号',
  `merchant_id`  BIGINT       NOT NULL DEFAULT 0 COMMENT '商户ID',
  `channel_type` VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `merchant_no`  VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '通道方商户号',
  `license_no`   VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '营业执照号',
  `license_img`  VARCHAR(256) NOT NULL DEFAULT '' COMMENT '营业执照照片URL',
  `legal_person` VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '法人姓名',
  `id_card_front`VARCHAR(256) NOT NULL DEFAULT '' COMMENT '身份证正面URL',
  `id_card_back` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '身份证背面URL',
  `phone`        VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '联系电话',
  `status`       TINYINT      NOT NULL DEFAULT 0 COMMENT '状态: 0-待提交 1-审核中 2-已通过 3-已驳回',
  `remark`       VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
  `reviewer`     VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '审核人',
  `review_remark`VARCHAR(256) NOT NULL DEFAULT '' COMMENT '审核意见',
  `review_time`  DATETIME     NULL COMMENT '审核时间',
  `ctime`        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_apply_no` (`apply_no`),
  KEY `idx_merchant_id` (`merchant_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='进件申请表';
```

#### refund_order (退款订单表)

```sql
CREATE TABLE IF NOT EXISTS `refund_order` (
  `id`              BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `refund_no`       VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '退款单号',
  `order_no`        VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '原支付订单号',
  `trade_refund_no` VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '通道退款单号',
  `merchant_id`     BIGINT       NOT NULL DEFAULT 0 COMMENT '商户ID',
  `merchant_name`   VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商户名称(冗余)',
  `refund_amount`   BIGINT       NOT NULL DEFAULT 0 COMMENT '退款金额(分)',
  `order_amount`    BIGINT       NOT NULL DEFAULT 0 COMMENT '原订单金额(分)',
  `channel_type`    VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `status`          TINYINT      NOT NULL DEFAULT 0 COMMENT '状态: 0-退款中 1-退款成功 2-退款失败',
  `reason`          VARCHAR(256) NOT NULL DEFAULT '' COMMENT '退款原因',
  `operator`        VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '操作人',
  `finish_time`     DATETIME     NULL COMMENT '完成时间',
  `ctime`           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_refund_no` (`refund_no`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_merchant_id` (`merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款订单表';
```

#### transfer_order (转账订单表)

```sql
CREATE TABLE IF NOT EXISTS `transfer_order` (
  `id`                BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `transfer_no`       VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '转账单号',
  `trade_transfer_no` VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '通道转账单号',
  `merchant_id`       BIGINT       NOT NULL DEFAULT 0 COMMENT '商户ID',
  `merchant_name`     VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商户名称(冗余)',
  `amount`            BIGINT       NOT NULL DEFAULT 0 COMMENT '转账金额(分)',
  `channel_type`      VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `payee_type`        VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '收款方式: account/openid/phone',
  `payee_account`     VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '收款账号',
  `payee_name`        VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '收款人姓名',
  `status`            TINYINT      NOT NULL DEFAULT 0 COMMENT '状态: 0-处理中 1-成功 2-失败',
  `remark`            VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
  `finish_time`       DATETIME     NULL COMMENT '完成时间',
  `ctime`             TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`             TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_transfer_no` (`transfer_no`),
  KEY `idx_merchant_id` (`merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='转账订单表';
```

#### transaction_flow (交易流水表)

```sql
CREATE TABLE IF NOT EXISTS `transaction_flow` (
  `id`              BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `flow_no`         VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '流水号',
  `order_no`        VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '关联订单号',
  `type`            VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '交易类型: pay/refund/transfer',
  `merchant_id`     BIGINT       NOT NULL DEFAULT 0 COMMENT '商户ID',
  `merchant_name`   VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商户名称(冗余)',
  `amount`          BIGINT       NOT NULL DEFAULT 0 COMMENT '交易金额(分)',
  `channel_type`    VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `channel_flow_no` VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '通道流水号',
  `direction`       VARCHAR(4)   NOT NULL DEFAULT '' COMMENT '资金方向: in/out',
  `status`          TINYINT      NOT NULL DEFAULT 0 COMMENT '状态: 0-处理中 1-已完成',
  `remark`          VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
  `ctime`           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_flow_no` (`flow_no`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_merchant_id` (`merchant_id`),
  KEY `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易流水表';
```

#### callback_record (回调通知记录表)

```sql
CREATE TABLE IF NOT EXISTS `callback_record` (
  `id`            BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `order_no`      VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '关联订单号',
  `type`          VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通知类型: pay/refund/transfer',
  `channel_type`  VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `direction`     VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通知方向: upstream(上游→平台)/downstream(平台→商户)',
  `notify_url`    VARCHAR(256) NOT NULL DEFAULT '' COMMENT '通知URL',
  `status`        TINYINT      NOT NULL DEFAULT 0 COMMENT '状态: 0-失败 1-成功 2-待重试',
  `http_status`   INT          NOT NULL DEFAULT 0 COMMENT 'HTTP状态码',
  `retry_count`   INT          NOT NULL DEFAULT 0 COMMENT '已重试次数',
  `max_retry`     INT          NOT NULL DEFAULT 5 COMMENT '最大重试次数',
  `request_body`  TEXT         NULL COMMENT '请求内容',
  `response_body` TEXT         NULL COMMENT '响应内容',
  `ctime`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='回调通知记录表';
```

#### reconciliation_bill (对账单表)

```sql
CREATE TABLE IF NOT EXISTS `reconciliation_bill` (
  `id`               BIGINT      NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `bill_date`        DATE        NOT NULL COMMENT '对账日期',
  `channel_type`     VARCHAR(16) NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `platform_count`   INT         NOT NULL DEFAULT 0 COMMENT '平台笔数',
  `platform_amount`  BIGINT      NOT NULL DEFAULT 0 COMMENT '平台金额(分)',
  `channel_count`    INT         NOT NULL DEFAULT 0 COMMENT '通道笔数',
  `channel_amount`   BIGINT      NOT NULL DEFAULT 0 COMMENT '通道金额(分)',
  `diff_count`       INT         NOT NULL DEFAULT 0 COMMENT '差异笔数',
  `diff_amount`      BIGINT      NOT NULL DEFAULT 0 COMMENT '差异金额(分)',
  `status`           TINYINT     NOT NULL DEFAULT 0 COMMENT '状态: 0-待对账 1-已对账 2-有差异',
  `ctime`            TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`            TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_date_channel` (`bill_date`, `channel_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='对账单表';
```

#### reconciliation_diff (对账差异表)

```sql
CREATE TABLE IF NOT EXISTS `reconciliation_diff` (
  `id`              BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `bill_date`       DATE         NOT NULL COMMENT '对账日期',
  `order_no`        VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '订单号',
  `channel_type`    VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `diff_type`       VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '差异类型: platform_only/channel_only/amount_mismatch/status_mismatch',
  `platform_amount` BIGINT       NULL COMMENT '平台金额(分)',
  `channel_amount`  BIGINT       NULL COMMENT '通道金额(分)',
  `diff_amount`     BIGINT       NOT NULL DEFAULT 0 COMMENT '差异金额(分)',
  `handle_status`   TINYINT      NOT NULL DEFAULT 0 COMMENT '处理状态: 0-待处理 1-已处理 2-已忽略',
  `handle_remark`   VARCHAR(256) NOT NULL DEFAULT '' COMMENT '处理备注',
  `handler`         VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '处理人',
  `ctime`           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  KEY `idx_bill_date` (`bill_date`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_handle_status` (`handle_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='对账差异表';
```

#### sys_role (角色表)

```sql
CREATE TABLE IF NOT EXISTS `sys_role` (
  `id`          BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `code`        VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '角色编码',
  `name`        VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '角色名称',
  `description` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '角色描述',
  `built_in`    TINYINT      NOT NULL DEFAULT 0 COMMENT '是否内置: 0-否 1-是',
  `status`      TINYINT      NOT NULL DEFAULT 1 COMMENT '状态: 0-停用 1-启用',
  `ctime`       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 初始角色数据
INSERT IGNORE INTO `sys_role` (`code`, `name`, `description`, `built_in`) VALUES
('admin', '管理员', '系统管理员，拥有全部权限', 1),
('operator', '运营', '运营人员，管理商户和订单', 1),
('finance', '财务', '财务人员，查看订单和对账', 1),
('viewer', '只读', '只读权限，仅查看数据', 1);
```

#### sys_role_perm (角色权限表)

```sql
CREATE TABLE IF NOT EXISTS `sys_role_perm` (
  `id`      BIGINT      NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `role_id` BIGINT      NOT NULL DEFAULT 0 COMMENT '角色ID',
  `perm`    VARCHAR(64) NOT NULL DEFAULT '' COMMENT '权限标识: merchant:list, order:payment 等',
  `ctime`   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_perm` (`role_id`, `perm`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限表';
```

#### operation_log (操作日志表)

```sql
CREATE TABLE IF NOT EXISTS `operation_log` (
  `id`           BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `operator`     VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '操作人',
  `module`       VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '操作模块: auth/merchant/incoming/payment/order/system',
  `action`       VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '操作类型: login/create/update/delete/export',
  `description`  VARCHAR(256) NOT NULL DEFAULT '' COMMENT '操作描述',
  `ip`           VARCHAR(64)  NOT NULL DEFAULT '' COMMENT 'IP地址',
  `user_agent`   VARCHAR(256) NOT NULL DEFAULT '' COMMENT 'UserAgent',
  `success`      TINYINT      NOT NULL DEFAULT 1 COMMENT '是否成功: 0-失败 1-成功',
  `duration`     INT          NOT NULL DEFAULT 0 COMMENT '耗时(ms)',
  `request_data` TEXT         NULL COMMENT '请求参数(JSON)',
  `ctime`        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

  PRIMARY KEY (`id`),
  KEY `idx_operator` (`operator`),
  KEY `idx_module` (`module`),
  KEY `idx_ctime` (`ctime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';
```

### 2.3 ER 关系

```
merchant (1) ──< merchant_app (N)
merchant (1) ──< payment_channel (N)
payment_channel (1) ──(1) payment_channel_config
merchant (1) ──< incoming_apply (N)
merchant (1) ──< payment_order (N)
merchant (1) ──< refund_order (N)
merchant (1) ──< transfer_order (N)
payment_order (1) ──< refund_order (N)   [通过 order_no 关联]
payment_order/refund_order/transfer_order ──< transaction_flow [通过 order_no]
payment_order/refund_order/transfer_order ──< callback_record  [通过 order_no]
sys_role (1) ──< sys_role_perm (N)
account.role ──> sys_role.code
```

---

## 三、接口设计

### 3.1 认证模块 (已实现，需扩展)

#### POST /gopay/v1/sso/login (已实现)
- 请求: `{ "username": "admin", "password": "123456" }`
- 响应: `{ "token": "jwt...", "userInfo": { "id", "username", "realName", "phone", "email", "role" } }`
- 变更: userInfo 新增 realName, email 字段；登录成功更新 last_login；记录操作日志

#### POST /gopay/v1/user/getInfo (已实现)
- 请求: Header `Authorization: Bearer <token>`
- 响应: `{ "id", "username", "realName", "phone", "email", "role", "lastLogin" }`
- 变更: 响应新增 realName, email, lastLogin 字段

#### POST /gopay/v1/user/changePwd (新增)
- 请求: `{ "oldPassword": "", "newPassword": "" }`
- 响应: `{ "code": 0 }`

#### POST /gopay/v1/user/profile (新增)
- 请求: `{ "realName": "", "phone": "", "email": "" }`
- 响应: `{ "code": 0 }`

### 3.2 仪表盘模块 (新增)

#### POST /gopay/v1/dashboard/stats
- 响应:
```json
{
  "todayAmount": 12845600,
  "todayCount": 1234,
  "todaySuccessRate": 98.6,
  "pendingApply": 3,
  "pendingRefund": 5
}
```

#### POST /gopay/v1/dashboard/recentOrders
- 请求: `{ "page": 1, "pageSize": 5 }`
- 响应: 分页的 payment_order 列表 (最近订单)

#### POST /gopay/v1/dashboard/channelDistribution
- 响应: `{ "alipay": 83496, "wechat": 44960 }` (今日各通道交易额)

#### POST /gopay/v1/dashboard/trend
- 响应: `{ "dates": ["04-01",...], "amounts": [820000,...], "counts": [680,...] }` (近7天趋势)

### 3.3 商户管理模块 (新增)

#### POST /gopay/v1/merchant/list
- 请求: `{ "page": 1, "pageSize": 20, "name": "", "contact": "", "status": -1 }`
- 响应: 分页 merchant 列表

#### POST /gopay/v1/merchant/add
- 请求: `{ "name": "", "contact": "", "phone": "", "email": "", "remark": "" }`
- 响应: `{ "id": 1 }`

#### POST /gopay/v1/merchant/update
- 请求: `{ "id": 1, "name": "", "contact": "", "phone": "", "email": "", "remark": "" }`

#### POST /gopay/v1/merchant/toggleStatus
- 请求: `{ "id": 1 }`
- 切换 status 0↔1

#### POST /gopay/v1/merchant/app/list
- 请求: `{ "page": 1, "pageSize": 20, "name": "", "appid": "", "platformType": -1, "merchantId": 0 }`
- 响应: 分页 merchant_app 列表 (携带 merchantName)

#### POST /gopay/v1/merchant/app/add
- 请求: `{ "name", "appid", "merchantId", "platformType", "merchantType", "notifyUrl", "returnUrl" }`
- 响应: `{ "id": 1 }`

#### POST /gopay/v1/merchant/app/update
- 请求: 同 add + id

#### POST /gopay/v1/merchant/options
- 响应: `[{ "id": 1, "name": "星辰科技" }, ...]` (状态正常的商户下拉选项)

### 3.4 进件管理模块 (新增)

#### POST /gopay/v1/incoming/apply/list
- 请求: `{ "page", "pageSize", "merchantName", "status", "channelType" }`
- 响应: 分页 incoming_apply 列表 (携带 merchantName)

#### POST /gopay/v1/incoming/apply/add
- 请求: `{ "merchantId", "channelType", "merchantNo", "licenseNo", "licenseImg", "legalPerson", "idCardFront", "idCardBack", "phone", "remark" }`
- status 默认 0(待提交)

#### POST /gopay/v1/incoming/apply/submit
- 请求: `{ "id": 1 }`
- status 0→1(审核中)

#### POST /gopay/v1/incoming/apply/review
- 请求: `{ "id": 1, "action": "pass|reject", "remark": "" }`
- pass: status 1→2, reject: status 1→3

#### POST /gopay/v1/incoming/record/list
- 请求: `{ "page", "pageSize", "merchantName", "channelType", "status", "reviewDate" }`
- 响应: 仅查询 status=2或3 的进件记录

#### POST /gopay/v1/incoming/record/detail
- 请求: `{ "id": 1 }`
- 响应: 进件详情

#### POST /gopay/v1/upload/image (通用图片上传)
- 请求: multipart/form-data, file 字段
- 响应: `{ "url": "https://..." }`

### 3.5 支付配置模块 (新增)

#### POST /gopay/v1/payment/channel/list
- 请求: `{ "page", "pageSize", "name", "code", "type", "status" }`
- 响应: 分页 payment_channel 列表 (携带 merchantName)

#### POST /gopay/v1/payment/channel/add
- 请求: `{ "name", "code", "type", "merchantId", "payMethods":["qrcode","page"], "feeRate": 0.6, "remark" }`

#### POST /gopay/v1/payment/channel/update
- 请求: 同 add + id

#### POST /gopay/v1/payment/channel/toggleStatus
- 请求: `{ "id": 1 }`

#### POST /gopay/v1/payment/channel/detail
- 请求: `{ "id": 1 }`
- 响应: channel + config 完整信息 (敏感字段脱敏)

#### POST /gopay/v1/payment/channel/config
- 请求: `{ "channelId": 1, "appId", "privateKey", "publicKey", "mchId", "apiKey", "serialNo", "notifyUrl", "signType", "sandbox" }`
- 保存/更新通道参数配置

### 3.6 订单中心模块 (新增)

#### POST /gopay/v1/order/payment/list
- 请求: `{ "page", "pageSize", "orderNo", "merchantName", "status", "channelType", "date" }`
- 响应: 分页 payment_order 列表

#### POST /gopay/v1/order/payment/detail
- 请求: `{ "id": 1 }`

#### POST /gopay/v1/order/payment/close
- 请求: `{ "id": 1 }`
- status 0→3(关闭)

#### POST /gopay/v1/order/payment/refund
- 请求: `{ "id": 1, "amount": 29900, "reason": "" }`
- 创建 refund_order，调用退款接口

#### POST /gopay/v1/order/payment/export
- 请求: 同 list 筛选条件
- 响应: CSV/Excel 文件下载

#### POST /gopay/v1/order/refund/list
- 请求: `{ "page", "pageSize", "refundNo", "orderNo", "status", "channelType" }`

#### POST /gopay/v1/order/refund/detail
- 请求: `{ "id": 1 }`

#### POST /gopay/v1/order/transfer/list
- 请求: `{ "page", "pageSize", "transferNo", "merchantName", "status", "channelType" }`

#### POST /gopay/v1/order/transfer/add
- 请求: `{ "merchantId", "channelType", "amount", "payeeType", "payeeAccount", "payeeName", "remark" }`
- 创建 transfer_order，调用转账接口

#### POST /gopay/v1/order/transfer/detail
- 请求: `{ "id": 1 }`

### 3.7 交易记录模块 (新增)

#### POST /gopay/v1/transaction/flow/list
- 请求: `{ "page", "pageSize", "flowNo", "orderNo", "type", "channelType", "date" }`
- 响应: 分页 transaction_flow 列表

#### POST /gopay/v1/transaction/flow/detail
- 请求: `{ "id": 1 }`

#### POST /gopay/v1/transaction/flow/stats
- 响应: `{ "incomeTotal": 100000, "expenseTotal": 50000, "totalCount": 120 }`

#### POST /gopay/v1/transaction/callback/list
- 请求: `{ "page", "pageSize", "orderNo", "type", "status", "channelType" }`

#### POST /gopay/v1/transaction/callback/detail
- 请求: `{ "id": 1 }`

#### POST /gopay/v1/transaction/callback/retry
- 请求: `{ "id": 1 }`
- 手动重试回调通知

### 3.8 对账管理模块 (新增)

#### POST /gopay/v1/recon/bill/list
- 请求: `{ "page", "pageSize", "date", "channelType", "status" }`

#### POST /gopay/v1/recon/bill/detail
- 请求: `{ "id": 1 }`

#### POST /gopay/v1/recon/bill/generate
- 请求: `{ "date": "2026-04-01", "channelType": "alipay" }`
- 生成对账单

#### POST /gopay/v1/recon/bill/reconcile
- 请求: `{ "id": 1 }`
- 执行对账

#### POST /gopay/v1/recon/diff/list
- 请求: `{ "page", "pageSize", "orderNo", "diffType", "handleStatus", "date" }`

#### POST /gopay/v1/recon/diff/detail
- 请求: `{ "id": 1 }`

#### POST /gopay/v1/recon/diff/handle
- 请求: `{ "id": 1, "action": "resolve|ignore", "remark": "" }`

#### POST /gopay/v1/recon/diff/export
- 请求: 同 list 筛选条件
- 响应: 文件下载

### 3.9 系统管理模块 (新增)

#### POST /gopay/v1/system/user/list
- 请求: `{ "page", "pageSize", "username", "phone", "status" }`

#### POST /gopay/v1/system/user/add
- 请求: `{ "username", "password", "realName", "phone", "email", "role" }`

#### POST /gopay/v1/system/user/update
- 请求: `{ "id", "realName", "phone", "email", "role" }`

#### POST /gopay/v1/system/user/toggleStatus
- 请求: `{ "id": 1 }`

#### POST /gopay/v1/system/user/resetPwd
- 请求: `{ "id": 1 }`
- 重置密码为默认值

#### POST /gopay/v1/system/role/list
- 请求: `{ "page", "pageSize" }`

#### POST /gopay/v1/system/role/add
- 请求: `{ "code", "name", "description" }`

#### POST /gopay/v1/system/role/update
- 请求: `{ "id", "name", "description" }`

#### POST /gopay/v1/system/role/toggleStatus
- 请求: `{ "id": 1 }`
- 不允许停用 built_in 角色

#### POST /gopay/v1/system/role/perms/update
- 请求: `{ "roleId": 1, "perms": ["merchant:list", "order:payment", ...] }`
- 全量更新角色权限

#### POST /gopay/v1/system/role/perms/list
- 请求: `{ "roleId": 1 }`
- 响应: `["merchant:list", ...]`

#### POST /gopay/v1/system/log/list
- 请求: `{ "page", "pageSize", "operator", "module", "action", "date" }`

#### POST /gopay/v1/system/log/detail
- 请求: `{ "id": 1 }`

#### POST /gopay/v1/system/log/export
- 请求: 同 list 筛选条件

---

## 四、权限标识定义

与前端 role.vue 中的权限组对齐:

| 权限组 | 权限标识 | 说明 |
|--------|---------|------|
| 商户管理 | merchant:list | 商户列表 |
| | merchant:app | 应用管理 |
| | merchant:edit | 商户编辑 |
| 进件管理 | incoming:apply | 进件申请 |
| | incoming:record | 进件记录 |
| | incoming:review | 进件审核 |
| 支付配置 | payment:channel | 通道管理 |
| | payment:config | 通道配置 |
| 订单中心 | order:payment | 支付订单 |
| | order:refund | 退款订单 |
| | order:transfer | 转账订单 |
| 交易记录 | transaction:flow | 交易流水 |
| | transaction:callback | 回调记录 |
| 对账管理 | recon:bill | 对账单 |
| | recon:diff | 对账差异 |
| 系统管理 | system:user | 用户管理 |
| | system:role | 角色管理 |
| | system:log | 操作日志 |

---

## 五、中间件设计

### 5.1 JWT 认证中间件 (已有，需完善)

除 `/sso/login` 和 `/monitor/ping` 外，所有接口需 Bearer Token 认证。

### 5.2 操作日志中间件 (新增)

在关键写操作接口上记录操作日志:
- 拦截请求参数、操作人、IP、UserAgent
- 异步写入 operation_log 表
- 记录耗时

### 5.3 权限校验中间件 (新增，可选)

根据 account.role → sys_role → sys_role_perm 判断接口权限。
初期可简化为: admin 拥有全部权限，其他角色按配置校验。

---

## 六、实现计划

### Phase 1: 基础设施 (预计工作量最大)

| 步骤 | 任务 | 涉及文件 |
|------|------|---------|
| 1.1 | 执行数据库 DDL (新建表 + ALTER 现有表) | mysql_ddl.sql |
| 1.2 | 新增 dm 层实体定义 (所有新表的 GORM Model) | app/dm/ |
| 1.3 | 新增错误码 (NotFound, Forbidden, Conflict) | errcode/ecode.go |
| 1.4 | 完善 JWT 认证中间件 (从 service 提取到 middleware) | app/router/middleware.go |
| 1.5 | 实现操作日志中间件 | app/router/middleware.go |
| 1.6 | 实现统一分页工具函数 | app/model/page.go |

### Phase 2: 认证与系统管理

| 步骤 | 任务 | 前端页面 |
|------|------|---------|
| 2.1 | 扩展 Login/GetInfo 返回新字段 | login, profile |
| 2.2 | 实现 changePwd, profile 接口 | profile |
| 2.3 | 实现系统用户 CRUD | system/user |
| 2.4 | 实现角色 CRUD + 权限管理 | system/role |
| 2.5 | 实现操作日志查询 | system/log |

### Phase 3: 商户与进件

| 步骤 | 任务 | 前端页面 |
|------|------|---------|
| 3.1 | 实现商户 CRUD + 状态切换 | merchant/list |
| 3.2 | 实现商户应用 CRUD | merchant/app |
| 3.3 | 实现进件申请 (新建/提交/审核) | incoming/apply |
| 3.4 | 实现进件记录查询 | incoming/record |
| 3.5 | 实现图片上传接口 | (通用) |

### Phase 4: 支付配置与通道

| 步骤 | 任务 | 前端页面 |
|------|------|---------|
| 4.1 | 实现支付通道 CRUD + 状态切换 | payment/channel |
| 4.2 | 实现通道参数配置 (Alipay/WeChat 分别处理) | payment/channel |
| 4.3 | 通道详情接口 (敏感字段脱敏) | payment/channel |

### Phase 5: 订单中心

| 步骤 | 任务 | 前端页面 |
|------|------|---------|
| 5.1 | 实现支付订单列表/详情/关闭 | order/payment |
| 5.2 | 实现退款功能 (创建退款单 + 调用支付宝/微信退款) | order/refund |
| 5.3 | 实现转账功能 (创建转账单 + 调用转账接口) | order/transfer |
| 5.4 | 实现订单导出 | order/payment |

### Phase 6: 交易记录

| 步骤 | 任务 | 前端页面 |
|------|------|---------|
| 6.1 | 实现交易流水查询 + 统计 | transaction/flow |
| 6.2 | 实现回调通知记录查询 | transaction/callback |
| 6.3 | 实现回调手动重试 | transaction/callback |

### Phase 7: 对账管理

| 步骤 | 任务 | 前端页面 |
|------|------|---------|
| 7.1 | 实现对账单列表/详情/生成/执行 | reconciliation/bill |
| 7.2 | 实现对账差异查询/处理/导出 | reconciliation/diff |

### Phase 8: 仪表盘

| 步骤 | 任务 | 前端页面 |
|------|------|---------|
| 8.1 | 实现统计数据聚合接口 | dashboard |
| 8.2 | 实现近7天趋势/通道分布接口 | dashboard |
| 8.3 | 实现最近订单接口 | dashboard |

---

## 七、文件新增/修改清单

### 新增文件

```
app/dm/
  merchant.go           — Merchant, MerchantApp 实体
  payment_channel.go    — PaymentChannel, PaymentChannelConfig 实体
  incoming.go           — IncomingApply 实体
  refund_order.go       — RefundOrder 实体
  transfer_order.go     — TransferOrder 实体
  transaction.go        — TransactionFlow, CallbackRecord 实体
  reconciliation.go     — ReconciliationBill, ReconciliationDiff 实体
  sys.go                — SysRole, SysRolePerm 实体
  operation_log.go      — OperationLog 实体

app/model/
  page.go               — PageReq, PageResp 分页通用结构
  merchant.go           — 商户相关请求/响应 DTO
  merchant_app.go       — 商户应用 DTO
  payment_channel.go    — 支付通道 DTO
  incoming.go           — 进件 DTO
  refund.go             — 退款 DTO
  transfer.go           — 转账 DTO
  transaction.go        — 交易流水/回调 DTO
  reconciliation.go     — 对账 DTO
  system.go             — 系统管理 DTO
  dashboard.go          — 仪表盘 DTO

app/dao/
  merchant.go           — 商户 DAO
  merchant_app.go       — 商户应用 DAO
  payment_channel.go    — 支付通道 DAO
  incoming.go           — 进件 DAO
  payment_order.go      — 支付订单 DAO
  refund_order.go       — 退款订单 DAO
  transfer_order.go     — 转账订单 DAO
  transaction.go        — 交易流水 DAO
  callback.go           — 回调记录 DAO
  reconciliation.go     — 对账 DAO
  sys_role.go           — 角色 DAO
  operation_log.go      — 操作日志 DAO

app/service/
  merchant.go           — 商户管理业务逻辑
  merchant_app.go       — 商户应用业务逻辑
  payment_channel.go    — 支付通道业务逻辑
  incoming.go           — 进件业务逻辑
  order.go              — 订单管理业务逻辑 (支付/退款/转账)
  transaction.go        — 交易记录业务逻辑
  reconciliation.go     — 对账业务逻辑
  system.go             — 系统管理业务逻辑 (用户/角色/日志)
  dashboard.go          — 仪表盘业务逻辑

app/router/
  middleware.go         — JWT 中间件、操作日志中间件
  merchant.go           — 商户路由处理器
  incoming.go           — 进件路由处理器
  payment_channel.go    — 支付通道路由处理器
  order.go              — 订单路由处理器
  transaction.go        — 交易记录路由处理器
  reconciliation.go     — 对账路由处理器
  system.go             — 系统管理路由处理器
  dashboard.go          — 仪表盘路由处理器
```

### 修改文件

```
mysql_ddl.sql           — 新增表 DDL + ALTER 语句
errcode/ecode.go        — 新增错误码
app/dm/db_gopay.go      — Account 实体新增字段
app/model/model.go      — LoginRsp/UserInfo 新增字段
app/service/account.go  — Login/GetUserInfo 扩展
app/service/service.go  — Service 结构无需改动 (DAO 已共享)
app/router/router.go    — 注册新路由
app/router/sso.go       — Login 响应扩展
app/router/user.go      — GetInfo 响应扩展 + 新增 changePwd/profile
```

---

## 八、接口总览 (按前端页面对齐)

| 前端页面 | 接口数量 | 关键接口 |
|---------|---------|---------|
| 登录 | 2 (已有) | login, getInfo |
| 个人中心 | 2 | changePwd, profile |
| 仪表盘 | 4 | stats, recentOrders, channelDistribution, trend |
| 商户列表 | 4 | list, add, update, toggleStatus |
| 商户应用 | 4 | list, add, update, options |
| 进件申请 | 4 | list, add, submit, review |
| 进件记录 | 2 | list, detail |
| 支付通道 | 5 | list, add, update, toggleStatus, config, detail |
| 支付订单 | 4 | list, detail, close, refund, export |
| 退款订单 | 2 | list, detail |
| 转账订单 | 3 | list, add, detail |
| 交易流水 | 3 | list, detail, stats |
| 回调记录 | 3 | list, detail, retry |
| 对账单 | 4 | list, detail, generate, reconcile |
| 对账差异 | 4 | list, detail, handle, export |
| 用户管理 | 4 | list, add, update, toggleStatus, resetPwd |
| 角色管理 | 5 | list, add, update, toggleStatus, perms/update, perms/list |
| 操作日志 | 3 | list, detail, export |
| 图片上传 | 1 | upload/image |
| **合计** | **~57** | |
