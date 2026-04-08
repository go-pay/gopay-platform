# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

gopay-platform 是一个 Go 支付处理平台，集成支付宝和微信支付。模块名 `gopay`，Go 1.22。

## 常用命令

```bash
# 构建
go build -o main main.go

# 运行
./main -conf app/conf/config.yaml

# 测试
go test ./... -v
go test ./app/dao -v -run TestName    # 单个测试
go test ./... -cover                   # 覆盖率

# Docker 构建
docker build -t gopay-platform:latest .
```

构建环境：CGO_ENABLED=0, GOOS=linux, GOARCH=amd64。

## 架构

分层架构，依赖方向从上到下：

```
main.go (入口：解析配置 → 创建 Service → 启动 HTTP 服务器 → 优雅关闭)
    ↓
router/  (HTTP 处理器，Gin 框架，路由前缀 /gopay/v1/)
    ↓
service/ (业务逻辑，支付编排，全局单例 srv)
    ↓
dao/     (数据访问，GORM + Redis)
```

### 关键目录

- `app/conf/` — 配置结构体定义，YAML 格式，通过 `-conf` 标志加载
- `app/router/` — HTTP 路由与处理器，使用 `go-pay/web` 封装 Gin
- `app/service/` — 业务逻辑层，初始化支付宝/微信客户端和 DAO
- `app/dao/` — 数据库连接和查询（GORM MySQL + Redis）
- `app/model/` — 请求/响应 DTO
- `app/dm/` — 数据映射层（占位）
- `pkg/config/` — 配置解析工具（支持 YAML/JSON/TOML）
- `ecode/` — 自定义错误码（RequestErr 10400, ServerErr 10500 等）

### 核心依赖

| 用途 | 库 |
|------|-----|
| Web 框架 | gin-gonic/gin |
| ORM | gorm.io/gorm (MySQL) |
| 缓存 | redis/go-redis/v9 |
| 支付 | go-pay/gopay（支付宝 & 微信） |
| Web 封装 | go-pay/web（中间件、限流、优雅关闭） |
| 日志 | go-pay/xlog |
| JSON | bytedance/sonic |
| 精确计算 | shopspring/decimal |

### 配置结构

配置文件 `app/conf/config.yaml`，结构体在 `app/conf/conf.go`：
- `cfg` — 日志级别、重载间隔
- `http` — 服务地址（默认 :2233）、超时、限流
- `redis` / `mysql` — 连接配置
- `pay_platform` — 支付宝/微信凭证和证书

### 支付流程

1. 客户端请求支付接口（传入 subject、money 分为单位）
2. Service 生成 UUID trade_no
3. 调用支付宝/微信 API
4. 返回二维码或支付链接

### 数据库

DDL 在 `mysql_ddl.sql`，主要表：account、company、payment_cfg、app、app_payment_cfg、payment_order。

## 开发模式

- 手动依赖注入（无 Wire 等代码生成工具）
- Service 为全局单例，通过构造函数初始化
- 错误处理使用 `go-pay/ecode` 自定义错误码
- 路由处理器通过包级变量 `svc` 访问 Service
- 微信支付 V3 部分实现，部分端点为 stub

## 关联前端项目

本项目的前端管理平台为 **gopay-platform-web**，位于 `/Users/my40138ml/workspace/jerry/go-pay/gopay-platform-web`。

- **技术栈**: Vue 3 + TypeScript + Vite + Vuetify 3 + Pinia + ECharts
- **开发端口**: 前端 localhost:3000，后端 localhost:2233
- **API 对接**: 前端通过 Vite proxy 将 `/gopay/v1/` 请求转发到后端 `:2233`
- **前端状态**: 全部 18 个页面 UI 已完成（使用 mock 数据），等待后端接口对接

## 后端开发必读文档

**开始开发前务必先阅读以下 3 份文档**，它们包含了完整的前后端对齐信息：

| 文档 | 路径 | 内容 |
|------|------|------|
| 技术设计 | `docs/backend-tech-design.md` | 数据库设计（14 张新表 + 2 张扩展表的完整 DDL）、67 个接口设计、中间件设计、8 阶段实施计划、文件新增/修改清单 |
| 接口详细设计 | `docs/api-detail.md` | 每个接口的完整请求/响应 JSON 示例、前端 TypeScript 接口定义、字段映射关系、枚举值映射、前端页面与接口对应关系表 |
| 实现进度 | `docs/IMPLEMENTATION_PROGRESS.md` | 逐接口的完成状态追踪（✅/🔧/⏳），每完成一个接口需更新此文档 |

### 开发流程

1. 读 `docs/IMPLEMENTATION_PROGRESS.md` 确认当前进度，找到下一个待实现的 Phase/接口
2. 读 `docs/api-detail.md` 中对应接口的详细设计（请求/响应示例、字段映射）
3. 按 `docs/backend-tech-design.md` 中的文件清单和分层架构实现: dm → dao → service → router
4. 实现完成后更新 `docs/IMPLEMENTATION_PROGRESS.md` 中对应条目状态为 ✅

### 关键约定

- **统一响应**: `{ "code": 0, "msg": "success", "data": {...} }`
- **金额**: 数据库和接口统一用分(int64)，前端展示时除以 100 为元
- **时间格式**: `YYYY-MM-DD HH:mm:ss`
- **分页**: 请求 `{ "page": 1, "pageSize": 20 }`，响应 `{ "list": [], "total": N, "page": 1, "pageSize": 20 }`
- **认证**: 除 login 和 ping 外，所有接口需 Bearer Token（JWT HS256, 24h）
- **错误码**: 在 `errcode/ecode.go` 中定义，详见技术设计文档

## 环境搭建

```bash
# 1. MySQL (本地默认无密码)
mysql -u root -e "source mysql_ddl.sql"

# 2. 配置文件已就绪，默认连接:
#    MySQL: root@127.0.0.1:3306/gopay (无密码)
#    Redis: localhost:6379, password="password" (当前已注释掉，未使用)
#    HTTP:  :2233

# 3. 启动后端
go build -o main main.go && ./main -conf app/conf/config.yaml

# 4. 前端 (另一个终端)
cd /Users/my40138ml/workspace/jerry/go-pay/gopay-platform-web
npm run dev   # localhost:3000, 自动代理 /gopay/v1/ → :2233
```

## 代码模式参考

新增接口时严格遵循以下模式，保持项目代码风格一致。

### Router 层 (app/router/*.go)

```go
// 包级变量 svc 访问 Service (已在 router.go 定义)
// var svc *service.Service

func merchantList(c *gin.Context) {
    req := new(model.MerchantListReq)
    if err := c.ShouldBindJSON(req); err != nil {
        xlog.Errorf("merchantList ShouldBindJSON(%v), err:%v", req, err)
        web.JSON(c, nil, errcode.RequestErr)
        return
    }
    rsp, err := svc.MerchantList(c, req)
    if err != nil {
        web.JSON(c, nil, err)
        return
    }
    web.JSON(c, rsp, nil)
}
```

**要点**: Bind → 调 svc → web.JSON 返回。不在 router 里写业务逻辑。

### 路由注册 (app/router/router.go)

```go
// 在 initRoute 函数的 v1 group 内追加
merchant := v1.Group("/merchant")
{
    merchant.POST("/list", merchantList)
    merchant.POST("/add", merchantAdd)
}
```

### Service 层 (app/service/*.go)

```go
func (s *Service) MerchantList(ctx context.Context, req *model.MerchantListReq) (*model.PageResp, error) {
    list, total, err := s.dao.MerchantList(ctx, req)
    if err != nil {
        xlog.Errorf("MerchantList dao err:%v", err)
        return nil, ec.ServerErr
    }
    return &model.PageResp{
        List:     list,
        Total:    total,
        Page:     req.Page,
        PageSize: req.PageSize,
    }, nil
}
```

**要点**: 调 dao → 组装响应。错误统一返回 errcode，用 xlog 记日志。

### DAO 层 (app/dao/*.go)

```go
func (d *Dao) MerchantList(ctx context.Context, req *model.MerchantListReq) (list []*dm.Merchant, total int64, err error) {
    if d.GopayDB == nil {
        return nil, 0, ErrNoDatabase
    }
    db := d.GopayDB.WithContext(ctx).Model(&dm.Merchant{})
    // 条件筛选
    if req.Name != "" {
        db = db.Where("name LIKE ?", "%"+req.Name+"%")
    }
    if req.Status >= 0 {
        db = db.Where("status = ?", req.Status)
    }
    // 先 count 再分页查
    if err = db.Count(&total).Error; err != nil {
        return
    }
    err = db.Order("id DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
    return
}
```

**要点**: 先判 GopayDB != nil → 链式条件 → Count → 分页 Find。

### DM 实体层 (app/dm/*.go)

```go
type Merchant struct {
    ID      int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    Name    string    `gorm:"column:name" json:"name"`
    Contact string    `gorm:"column:contact" json:"contact"`
    Phone   string    `gorm:"column:phone" json:"phone"`
    Email   string    `gorm:"column:email" json:"email"`
    Status  int8      `gorm:"column:status" json:"status"`
    Remark  string    `gorm:"column:remark" json:"remark"`
    Ctime   time.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
    Utime   time.Time `gorm:"column:utime;autoUpdateTime" json:"utime"`
}

func (Merchant) TableName() string { return "merchant" }
```

**要点**: gorm tag 指定 column，json tag 驼峰，Pwd 字段用 `json:"-"` 隐藏。时间用 autoCreateTime/autoUpdateTime。

### Model DTO 层 (app/model/*.go)

```go
// 分页通用
type PageReq struct {
    Page     int `json:"page" binding:"required,min=1"`
    PageSize int `json:"pageSize" binding:"required,min=1,max=100"`
}

type PageResp struct {
    List     interface{} `json:"list"`
    Total    int64       `json:"total"`
    Page     int         `json:"page"`
    PageSize int         `json:"pageSize"`
}

// 业务请求 (内嵌 PageReq)
type MerchantListReq struct {
    PageReq
    Name    string `json:"name"`
    Contact string `json:"contact"`
    Status  int8   `json:"status"` // -1=全部, 0=禁用, 1=正常
}
```

**要点**: 列表请求内嵌 PageReq。json tag 使用 camelCase 与前端对齐。必填字段加 `binding:"required"`。
