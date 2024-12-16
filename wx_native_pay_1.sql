CREATE TABLE IF NOT EXISTS shop_goods
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
    sku_id     VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '商品sku_id',
    goods_name VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '商品名称',
    goods_desc VARCHAR(512) NOT NULL DEFAULT '' COMMENT '商品描述',
    unit_price INT          NOT NULL DEFAULT 0 COMMENT '商品单价(分)',
    ctime      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    utime      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',

    INDEX idx_sku_id (sku_id),
    INDEX idx_goods_name (goods_name),
    INDEX idx_unit_price (unit_price)
) ENGINE = InnoDB
  CHARSET = utf8mb4 COMMENT '店铺商品表';

INSERT INTO shop_goods (sku_id, goods_name, goods_desc, unit_price)
VALUES ('GOOD0001', '商品1', '我是商品1', 1),
       ('GOOD0002', '商品2', '我是商品2', 2),
       ('GOOD0003', '商品3', '我是商品3', 3);

-- ================================================================================================================

CREATE TABLE IF NOT EXISTS payment_order
(
    `id`             BIGINT       NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
    `user_id`        BIGINT       NOT NULL DEFAULT 0 COMMENT '绑定用户id',
    `qrcode`         VARCHAR(256) NOT NULL DEFAULT '' COMMENT '支付二维码',
    `trade_no`       VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '订单号',
    `transaction_id` VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '支付交易流水号',
    `payment_type`   TINYINT      NOT NULL DEFAULT 0 COMMENT '支付类型：0-微信，1-支付宝',
    `pay_money`      BIGINT       NOT NULL DEFAULT 0 COMMENT '支付金额(分)',
    `status`         TINYINT      NOT NULL DEFAULT 0 COMMENT '订单状态：0-待支付，1-支付成功，3-订单关闭',
    `pay_time`       DATETIME     NULL COMMENT '支付时间',
    `remark`         VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
    `notify_body`    TEXT         NULL COMMENT '回调参数信息',
    `ctime`          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `utime`          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_trade_no` (`trade_no`),
    KEY `idx_transaction_id` (`transaction_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '支付订单表';