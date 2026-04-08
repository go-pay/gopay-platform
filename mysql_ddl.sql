CREATE DATABASE `gopay` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;


-- 账户表
CREATE TABLE IF NOT EXISTS `account` (
  `id`         BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `uname`      VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '用户名',
  `pwd`        VARCHAR(255) NOT NULL DEFAULT '' COMMENT '密码(bcrypt hash)',
  `real_name`  VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '真实姓名',
  `phone`      VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '手机号',
  `email`      VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '邮箱',
  `role`       VARCHAR(16)  NOT NULL DEFAULT 'viewer' COMMENT '角色: admin/operator/finance/viewer',
  `status`     TINYINT      NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用 1-正常',
  `last_login` DATETIME     NULL COMMENT '最后登录时间',
  `ctime`      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_uname` (`uname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='账户表';


-- 支付订单表
CREATE TABLE IF NOT EXISTS `payment_order` (
  `id`             BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `order_no`       VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '平台订单号',
  `user_id`        BIGINT       NOT NULL DEFAULT 0 COMMENT '绑定用户ID',
  `merchant_id`    BIGINT       NOT NULL DEFAULT 0 COMMENT '商户ID',
  `merchant_name`  VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商户名称',
  `qrcode`         VARCHAR(128) NOT NULL DEFAULT '' COMMENT '支付二维码',
  `trade_no`       VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '订单号',
  `out_trade_no`   VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商户订单号',
  `transaction_id` VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '支付交易流水号',
  `payment_type`   TINYINT      NOT NULL DEFAULT 0 COMMENT '支付类型: 0-微信 1-支付宝',
  `channel_type`   VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `pay_method`     VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '支付方式: qrcode/page/wap/app/jsapi/miniapp',
  `subject`        VARCHAR(128) NOT NULL DEFAULT '' COMMENT '商品描述',
  `client_ip`      VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '客户端IP',
  `pay_money`      BIGINT       NOT NULL DEFAULT 0 COMMENT '支付金额(分)',
  `status`         TINYINT      NOT NULL DEFAULT 0 COMMENT '订单状态: 0-待支付 1-支付成功 3-订单关闭',
  `pay_time`       DATETIME     NULL COMMENT '支付时间',
  `remark`         VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
  `notify_body`    TEXT         NULL COMMENT '回调参数信息',
  `notified`       TINYINT      NOT NULL DEFAULT 0 COMMENT '是否已回调通知: 0-否 1-是',
  `ctime`          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_trade_no` (`trade_no`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_merchant_id` (`merchant_id`),
  KEY `idx_transaction_id` (`transaction_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付订单表';


-- 商户表
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


-- 商户应用表
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


-- 支付通道表
CREATE TABLE IF NOT EXISTS `payment_channel` (
  `id`          BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `name`        VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '通道名称',
  `code`        VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '通道编码',
  `type`        VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `merchant_id` BIGINT       NOT NULL DEFAULT 0 COMMENT '所属商户ID',
  `pay_methods` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '支持的支付方式(逗号分隔)',
  `fee_rate`    DECIMAL(5,2) NOT NULL DEFAULT 0.00 COMMENT '费率(%)',
  `status`      TINYINT      NOT NULL DEFAULT 1 COMMENT '状态: 0-停用 1-启用',
  `remark`      VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
  `ctime`       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_merchant_id` (`merchant_id`),
  KEY `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付通道表';


-- 支付通道参数配置表
CREATE TABLE IF NOT EXISTS `payment_channel_config` (
  `id`         BIGINT        NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `channel_id` BIGINT        NOT NULL DEFAULT 0 COMMENT '支付通道ID',
  `app_id`     VARCHAR(32)   NOT NULL DEFAULT '' COMMENT 'AppID',
  `mch_id`     VARCHAR(32)   NOT NULL DEFAULT '' COMMENT '商户号(微信)',
  `private_key`TEXT          NULL COMMENT '私钥',
  `public_key` TEXT          NULL COMMENT '公钥',
  `api_key`    VARCHAR(256)  NOT NULL DEFAULT '' COMMENT 'API密钥(微信V3)',
  `serial_no`  VARCHAR(64)   NOT NULL DEFAULT '' COMMENT '证书序列号(微信)',
  `notify_url` VARCHAR(256)  NOT NULL DEFAULT '' COMMENT '回调通知URL',
  `sign_type`  VARCHAR(8)    NOT NULL DEFAULT 'RSA2' COMMENT '签名方式: RSA2/RSA',
  `sandbox`    TINYINT       NOT NULL DEFAULT 0 COMMENT '是否沙箱: 0-生产 1-沙箱',
  `ctime`      TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`      TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_channel_id` (`channel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付通道参数配置表';


-- 进件申请表
CREATE TABLE IF NOT EXISTS `incoming_apply` (
  `id`            BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `apply_no`      VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '申请编号',
  `merchant_id`   BIGINT       NOT NULL DEFAULT 0 COMMENT '商户ID',
  `channel_type`  VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `merchant_no`   VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '通道方商户号',
  `license_no`    VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '营业执照号',
  `license_img`   VARCHAR(256) NOT NULL DEFAULT '' COMMENT '营业执照照片URL',
  `legal_person`  VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '法人姓名',
  `id_card_front` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '身份证正面URL',
  `id_card_back`  VARCHAR(256) NOT NULL DEFAULT '' COMMENT '身份证背面URL',
  `phone`         VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '联系电话',
  `status`        TINYINT      NOT NULL DEFAULT 0 COMMENT '状态: 0-待提交 1-审核中 2-已通过 3-已驳回',
  `remark`        VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
  `reviewer`      VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '审核人',
  `review_remark` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '审核意见',
  `review_time`   DATETIME     NULL COMMENT '审核时间',
  `ctime`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_apply_no` (`apply_no`),
  KEY `idx_merchant_id` (`merchant_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='进件申请表';


-- 退款订单表
CREATE TABLE IF NOT EXISTS `refund_order` (
  `id`              BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `refund_no`       VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '退款单号',
  `order_no`        VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '原支付订单号',
  `trade_refund_no` VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '通道退款单号',
  `merchant_id`     BIGINT       NOT NULL DEFAULT 0 COMMENT '商户ID',
  `merchant_name`   VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商户名称',
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


-- 转账订单表
CREATE TABLE IF NOT EXISTS `transfer_order` (
  `id`                BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `transfer_no`       VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '转账单号',
  `trade_transfer_no` VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '通道转账单号',
  `merchant_id`       BIGINT       NOT NULL DEFAULT 0 COMMENT '商户ID',
  `merchant_name`     VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商户名称',
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


-- 交易流水表
CREATE TABLE IF NOT EXISTS `transaction_flow` (
  `id`              BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `flow_no`         VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '流水号',
  `order_no`        VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '关联订单号',
  `type`            VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '交易类型: pay/refund/transfer',
  `merchant_id`     BIGINT       NOT NULL DEFAULT 0 COMMENT '商户ID',
  `merchant_name`   VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商户名称',
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


-- 回调通知记录表
CREATE TABLE IF NOT EXISTS `callback_record` (
  `id`            BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `order_no`      VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '关联订单号',
  `type`          VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通知类型: pay/refund/transfer',
  `channel_type`  VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `direction`     VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '通知方向: upstream/downstream',
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


-- 对账单表
CREATE TABLE IF NOT EXISTS `reconciliation_bill` (
  `id`              BIGINT      NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `bill_date`       DATE        NOT NULL COMMENT '对账日期',
  `channel_type`    VARCHAR(16) NOT NULL DEFAULT '' COMMENT '通道类型: alipay/wechat',
  `platform_count`  INT         NOT NULL DEFAULT 0 COMMENT '平台笔数',
  `platform_amount` BIGINT      NOT NULL DEFAULT 0 COMMENT '平台金额(分)',
  `channel_count`   INT         NOT NULL DEFAULT 0 COMMENT '通道笔数',
  `channel_amount`  BIGINT      NOT NULL DEFAULT 0 COMMENT '通道金额(分)',
  `diff_count`      INT         NOT NULL DEFAULT 0 COMMENT '差异笔数',
  `diff_amount`     BIGINT      NOT NULL DEFAULT 0 COMMENT '差异金额(分)',
  `status`          TINYINT     NOT NULL DEFAULT 0 COMMENT '状态: 0-待对账 1-已对账 2-有差异',
  `ctime`           TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `utime`           TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_date_channel` (`bill_date`, `channel_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='对账单表';


-- 对账差异表
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


-- 角色表
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


-- 角色权限表
CREATE TABLE IF NOT EXISTS `sys_role_perm` (
  `id`      BIGINT      NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `role_id` BIGINT      NOT NULL DEFAULT 0 COMMENT '角色ID',
  `perm`    VARCHAR(64) NOT NULL DEFAULT '' COMMENT '权限标识',
  `ctime`   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_perm` (`role_id`, `perm`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限表';


-- 操作日志表
CREATE TABLE IF NOT EXISTS `operation_log` (
  `id`           BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `operator`     VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '操作人',
  `module`       VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '操作模块',
  `action`       VARCHAR(16)  NOT NULL DEFAULT '' COMMENT '操作类型',
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


-- 初始角色数据
INSERT IGNORE INTO `sys_role` (`code`, `name`, `description`, `built_in`) VALUES
('admin', '管理员', '系统管理员，拥有全部权限', 1),
('operator', '运营', '运营人员，管理商户和订单', 1),
('finance', '财务', '财务人员，查看订单和对账', 1),
('viewer', '只读', '只读权限，仅查看数据', 1);

-- 初始管理员账号（密码: admin，明文存储，首次登录后建议改为 bcrypt hash）
INSERT IGNORE INTO `account` (`uname`, `pwd`, `real_name`, `phone`, `role`, `status`) VALUES ('admin', 'admin', '超级管理员', '13800000001', 'admin', 1);
