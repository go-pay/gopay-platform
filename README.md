# gopay-platform

GoPay payment platform

## Terminal Run

```bash
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go && strip main

main -conf app/conf/config.yaml
```

安装总结

✅ MySQL 9.6 已安装并运行
- 内存占用: 约 335 MB（已优化配置）
- InnoDB 缓冲池: 64 MB（默认 128 MB）
- 最大连接数: 50（默认 151）

配置信息

- Root 密码: root123456
- 端口: 3306
- 字符集: utf8mb4
- 时区: +08:00
- 已创建数据库: gopay

常用命令

# 连接 MySQL
mysql -u root -proot123456

# 连接到 gopay_dev 数据库
mysql -u root -proot123456 gopay

# 启动/停止/重启服务
brew services start mysql
brew services stop mysql
brew services restart mysql

# 查看服务状态
brew services list | grep mysql

配置文件位置

- 配置文件: /usr/local/etc/my.cnf
- 数据目录: /usr/local/var/mysql
- 错误日志: /usr/local/var/mysql/mac.err
