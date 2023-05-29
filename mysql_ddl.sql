CREATE DATABASE `gopay` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

-- 表随便写写，后续再优化修改

CREATE TABLE IF NOT EXISTS `account`
(
    `id`    BIGINT      NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
    `uname` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '用户名',
    `pwd`   VARCHAR(32) NOT NULL DEFAULT '' COMMENT '密码',
    `phone` VARCHAR(16) NOT NULL DEFAULT '' COMMENT '手机号',
    `ctime` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `utime` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uname` (`uname`),
    KEY `ctime` (`ctime`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '账户表';


CREATE TABLE IF NOT EXISTS `company`
(
    `id`    BIGINT      NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
    `name`  VARCHAR(32) NOT NULL DEFAULT '' COMMENT '公司名',
    `ctime` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `utime` TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',

    PRIMARY KEY (`id`),
    KEY `name` (`name`),
    KEY `ctime` (`ctime`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '公司表';


CREATE TABLE IF NOT EXISTS `payment_info`
(
    `id`                     BIGINT        NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
    `payment_type`           TINYINT       NOT NULL DEFAULT 0 COMMENT '支付类型：0-微信，1-支付宝',
    `wx_mch_id`              VARCHAR(16)   NOT NULL DEFAULT '' COMMENT '微信商户号',
    `wx_api_key`             VARCHAR(32)   NOT NULL DEFAULT '' COMMENT '微信API密钥',
    `wx_apiv3_key`           VARCHAR(32)   NOT NULL DEFAULT '' COMMENT '微信APIv3密钥',
    `wx_serial_no`           VARCHAR(64)   NOT NULL DEFAULT '' COMMENT '微信APIv3证书序列号',
    `wx_private_key`         VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '微信APIv3私钥内容',
    `wx_platform_serial_no`  VARCHAR(64)   NOT NULL DEFAULT '' COMMENT '微信平台公钥证书序列号',
    `wx_platform_public_key` VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '微信平台公钥内容',
    `ali_private_key`        VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '支付宝应用私钥内容',
    `ali_public_key`         VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '支付宝应用公钥内容',
    `ali_root_cert`          TEXT          NULL COMMENT '支付宝根证书内容',
    `ali_app_cert`           VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '支付宝APP公钥证书内容',
    `ali_alipay_public_cert` VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '支付宝公钥证书内容',
    `ctime`                  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `utime`                  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',

    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '支付信息配置表';


CREATE TABLE IF NOT EXISTS `app`
(
    `id`            BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
    `platform`      TINYINT      NOT NULL DEFAULT 0 COMMENT '应用平台：0-微信移动应用，1-微信网站应用，2-微信公众号，3-微信小程序，4-微信第三方平台，5-支付宝网页/移动应用，6-支付宝小程序，7-支付宝生活号，8-支付宝第三方平台',
    `merchant_type` TINYINT      NOT NULL DEFAULT 0 COMMENT '商户类型：0-商户，1-服务商',
    `appid`         VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '应用appid',
    `return_url`    VARCHAR(256) NOT NULL DEFAULT '' COMMENT '支付return_url',
    `notify_url`    VARCHAR(256) NOT NULL DEFAULT '' COMMENT '回调通知URL',
    `ctime`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `utime`         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',

    PRIMARY KEY (`id`),
    KEY `platform` (`platform`),
    KEY `appid` (`appid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT 'app应用表';


CREATE TABLE IF NOT EXISTS `app_payment_info`
(
    `id`              BIGINT      NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
    `app_id`          BIGINT      NOT NULL DEFAULT 0 COMMENT 'app应用表id',
    `payment_info_id` BIGINT      NOT NULL DEFAULT 0 COMMENT '支付信息配置表id',
    `appid`           VARCHAR(32) NOT NULL DEFAULT '' COMMENT '应用appid',
    `ctime`           TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `utime`           TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',

    PRIMARY KEY (`id`),
    KEY `app_id` (`app_id`),
    KEY `payment_info_id` (`payment_info_id`),
    KEY `appid` (`appid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT 'app应用支付配置关联表';
