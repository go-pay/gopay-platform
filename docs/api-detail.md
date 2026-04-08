# GoPay Platform 接口详细设计 - 前后端字段映射

> 本文档包含每个前端页面的完整字段定义、Mock 数据结构、API 请求/响应示例。
> 任何人拿到此文档即可直接开发后端接口，无需再分析前端代码。

---

## 通用约定

### 统一响应格式
```json
// 成功
{ "code": 0, "msg": "success", "data": { ... } }

// 失败
{ "code": 10400, "msg": "参数错误", "data": null }
```

### 统一分页请求
```json
{ "page": 1, "pageSize": 20, "...筛选字段" }
```

### 统一分页响应
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

### 金额处理
- 数据库和接口传输统一用 **分(int64)**
- 前端展示时除以 100 显示为元，保留两位小数
- 前端 mock 数据中的 `amount: 299.00` 实际传输应为 `29900`

### 时间格式
- 数据库: TIMESTAMP / DATETIME
- 接口响应: `"2026-04-01 10:30:00"` (格式 YYYY-MM-DD HH:mm:ss)

### 认证
- 除 `/sso/login` 和 `/monitor/ping` 外，所有接口需 Header: `Authorization: Bearer <jwt_token>`
- JWT 有效期 24 小时，HS256 签名

---

## 一、登录与认证 (已实现，需扩展)

### 1.1 POST /gopay/v1/sso/login

**前端文件**: `src/store/user.ts` → `login()` action
**前端调用**: `POST /gopay/v1/sso/login`, body: `{ username, password }`

请求:
```json
{
  "username": "admin",
  "password": "admin"
}
```

响应:
```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "userInfo": {
      "id": 1,
      "username": "admin",
      "realName": "超级管理员",
      "phone": "13800000001",
      "email": "admin@gopay.com",
      "role": "admin",
      "lastLogin": "2026-04-07 09:30:00"
    }
  }
}
```

**与现有实现的差异**: 需新增 `realName`, `email`, `lastLogin` 字段。

### 1.2 POST /gopay/v1/user/getInfo

**前端文件**: `src/store/user.ts` → `fetchUserInfo()` action
**前端调用**: 携带 Bearer Token

响应:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "username": "admin",
    "realName": "超级管理员",
    "phone": "13800000001",
    "email": "admin@gopay.com",
    "role": "admin",
    "lastLogin": "2026-04-07 09:30:00"
  }
}
```

### 1.3 POST /gopay/v1/user/changePwd

**前端文件**: `src/views/profile/index.vue`

请求:
```json
{
  "oldPassword": "admin",
  "newPassword": "newpass123",
  "confirmPassword": "newpass123"
}
```

响应:
```json
{ "code": 0, "msg": "success" }
```

### 1.4 POST /gopay/v1/user/profile

**前端文件**: `src/views/profile/index.vue`

请求:
```json
{
  "realName": "超级管理员",
  "phone": "13800000001",
  "email": "admin@gopay.com"
}
```

响应:
```json
{ "code": 0, "msg": "success" }
```

---

## 二、仪表盘 Dashboard

**前端文件**: `src/views/dashboard/index.vue`

### 2.1 POST /gopay/v1/dashboard/stats

响应:
```json
{
  "code": 0,
  "data": {
    "todayAmount": 12845600,
    "todayCount": 1234,
    "todaySuccessRate": 98.6,
    "pendingApply": 3,
    "pendingRefund": 5
  }
}
```

前端展示映射:
| 前端展示 | 字段 | 格式化 |
|---------|------|--------|
| 今日交易额 ¥128,456.00 | todayAmount | 分→元, 千分位 |
| 今日交易笔数 1,234 | todayCount | 千分位 |
| 今日成功率 98.6% | todaySuccessRate | 直接拼接% |
| 待审核进件 3 | pendingApply | 直接展示 |
| 待处理退款 5 | pendingRefund | 直接展示 |

### 2.2 POST /gopay/v1/dashboard/recentOrders

请求:
```json
{ "page": 1, "pageSize": 5 }
```

响应:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "orderNo": "PAY20260407001234",
        "merchantName": "星辰科技有限公司",
        "amount": 258000,
        "channelType": "alipay",
        "status": 1,
        "ctime": "2026-04-07 14:23:15"
      }
    ],
    "total": 100
  }
}
```

前端展示映射:
| 前端列 | 后端字段 | 格式化 |
|--------|---------|--------|
| 订单号 | orderNo | 原样 (monospace) |
| 商户名称 | merchantName | 原样 |
| 金额 | amount | 分→元 (¥2,580.00) |
| 支付通道 | channelType | alipay→支付宝, wechat→微信 |
| 状态 | status | 0→待支付(amber), 1→已支付(green), 3→已关闭(grey) |
| 创建时间 | ctime | 原样 |

### 2.3 POST /gopay/v1/dashboard/channelDistribution

响应:
```json
{
  "code": 0,
  "data": {
    "alipay": 8349600,
    "wechat": 4496000
  }
}
```

### 2.4 POST /gopay/v1/dashboard/trend

响应:
```json
{
  "code": 0,
  "data": {
    "dates": ["04-02", "04-03", "04-04", "04-05", "04-06", "04-07", "04-08"],
    "amounts": [820000, 950000, 1130000, 1080000, 1250000, 1310000, 1280000],
    "counts": [680, 720, 910, 850, 1030, 1120, 1230]
  }
}
```

---

## 三、商户管理

### 3.1 商户列表 - POST /gopay/v1/merchant/list

**前端文件**: `src/views/merchant/list.vue`

**前端 TypeScript 接口:**
```typescript
interface Merchant {
  id: number
  name: string
  contact: string
  phone: string
  email: string
  status: number    // 0-禁用 1-正常
  remark: string
  ctime: string
}
```

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "name": "",
  "contact": "",
  "status": -1
}
```
> status 传 -1 表示全部

响应:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "星辰科技有限公司",
        "contact": "张三",
        "phone": "13800138001",
        "email": "zhangsan@xingchen.com",
        "status": 1,
        "remark": "",
        "ctime": "2026-03-15 10:30:00"
      }
    ],
    "total": 6,
    "page": 1,
    "pageSize": 20
  }
}
```

前端表格列:
| 列名 | 宽度 | 字段 | 格式 |
|------|------|------|------|
| ID | 60px | id | 数字 |
| 商户名称 | 180px | name | 文本 |
| 联系人 | 120px | contact | 文本 |
| 联系电话 | 140px | phone | 文本 |
| 邮箱 | auto | email | 文本 |
| 状态 | 80px | status | 1→正常(green chip), 0→禁用(grey chip) |
| 创建时间 | 160px | ctime | 时间字符串 |
| 操作 | 120px | - | 编辑/查看应用/禁用\|启用 |

### 3.2 新增商户 - POST /gopay/v1/merchant/add

请求:
```json
{
  "name": "星辰科技有限公司",
  "contact": "张三",
  "phone": "13800138001",
  "email": "zhangsan@xingchen.com",
  "remark": ""
}
```
> name, contact 必填

响应:
```json
{ "code": 0, "data": { "id": 7 } }
```

### 3.3 编辑商户 - POST /gopay/v1/merchant/update

请求:
```json
{
  "id": 1,
  "name": "星辰科技有限公司",
  "contact": "张三丰",
  "phone": "13800138001",
  "email": "zhangsan@xingchen.com",
  "remark": "VIP商户"
}
```

### 3.4 切换商户状态 - POST /gopay/v1/merchant/toggleStatus

请求:
```json
{ "id": 1 }
```

### 3.5 商户下拉选项 - POST /gopay/v1/merchant/options

响应:
```json
{
  "code": 0,
  "data": [
    { "id": 1, "name": "星辰科技有限公司" },
    { "id": 2, "name": "云海数字传媒" },
    { "id": 3, "name": "极光电子商务" },
    { "id": 5, "name": "蓝鲸网络科技" },
    { "id": 6, "name": "九州在线商贸" }
  ]
}
```
> 仅返回 status=1 的商户

---

## 四、商户应用

### 4.1 应用列表 - POST /gopay/v1/merchant/app/list

**前端文件**: `src/views/merchant/app.vue`

**前端 TypeScript 接口:**
```typescript
interface MerchantApp {
  id: number
  name: string
  appid: string
  merchantId: number
  merchantName: string
  platformType: number
  merchantType: number  // 0-商户, 1-服务商
  status: number
  notifyUrl: string
  returnUrl: string
  ctime: string
}
```

**platformType 枚举:**
| 值 | 标签 | Chip 颜色 |
|---|------|----------|
| 0 | 微信移动应用 | green |
| 1 | 微信网站应用 | green |
| 2 | 微信公众号 | green |
| 3 | 微信小程序 | green |
| 5 | 支付宝网页/移动 | blue |
| 6 | 支付宝小程序 | blue |
| 7 | 支付宝生活号 | blue |

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "name": "",
  "appid": "",
  "platformType": -1,
  "merchantId": 0
}
```

响应:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "星辰小程序",
        "appid": "wx1234567890abcdef",
        "merchantId": 1,
        "merchantName": "星辰科技有限公司",
        "platformType": 3,
        "merchantType": 0,
        "status": 1,
        "notifyUrl": "https://example.com/notify",
        "returnUrl": "https://example.com/return",
        "ctime": "2026-03-16 10:00:00"
      }
    ],
    "total": 5
  }
}
```

### 4.2 新增应用 - POST /gopay/v1/merchant/app/add

请求:
```json
{
  "name": "星辰小程序",
  "appid": "wx1234567890abcdef",
  "merchantId": 1,
  "platformType": 3,
  "merchantType": 0,
  "notifyUrl": "https://example.com/notify",
  "returnUrl": "https://example.com/return"
}
```
> name, appid, merchantId, platformType 必填

### 4.3 编辑应用 - POST /gopay/v1/merchant/app/update

请求: 同 add + id (appid 不可修改)

---

## 五、进件管理

### 5.1 进件申请列表 - POST /gopay/v1/incoming/apply/list

**前端文件**: `src/views/incoming/apply.vue`

**前端 TypeScript 接口:**
```typescript
interface ApplyItem {
  id: number
  applyNo: string          // 格式: INC20260401001
  merchantId: number
  merchantName: string
  channelType: 'alipay' | 'wechat'
  merchantNo: string       // 通道方商户号
  licenseNo: string        // 营业执照号
  legalPerson: string      // 法人姓名
  phone: string
  status: number           // 0-待提交 1-审核中 2-已通过 3-已驳回
  remark: string
  ctime: string
}
```

**status 状态映射:**
| 值 | 标签 | Chip 颜色 | 可用操作 |
|---|------|----------|---------|
| 0 | 待提交 | grey | 编辑, 提交审核 |
| 1 | 审核中 | amber | 审核 |
| 2 | 已通过 | green | 查看 |
| 3 | 已驳回 | red | 查看 |

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "merchantName": "",
  "status": -1,
  "channelType": ""
}
```

响应:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "applyNo": "INC20260401001",
        "merchantId": 1,
        "merchantName": "星辰科技有限公司",
        "channelType": "alipay",
        "merchantNo": "2088441234567890",
        "licenseNo": "91110108MA12345678",
        "legalPerson": "张三",
        "phone": "13800138001",
        "status": 2,
        "remark": "",
        "ctime": "2026-04-01 10:00:00"
      }
    ],
    "total": 6
  }
}
```

### 5.2 新建进件申请 - POST /gopay/v1/incoming/apply/add

请求:
```json
{
  "merchantId": 1,
  "channelType": "alipay",
  "merchantNo": "",
  "licenseNo": "91110108MA12345678",
  "licenseImg": "https://cdn.example.com/license.jpg",
  "legalPerson": "张三",
  "idCardFront": "https://cdn.example.com/front.jpg",
  "idCardBack": "https://cdn.example.com/back.jpg",
  "phone": "13800138001",
  "remark": "",
  "submit": false
}
```
> submit=false 保存草稿(status=0), submit=true 直接提交审核(status=1)
> applyNo 由后端自动生成: INC + YYYYMMDD + 3位序号

### 5.3 提交审核 - POST /gopay/v1/incoming/apply/submit

请求:
```json
{ "id": 5 }
```
> 将 status 从 0→1

### 5.4 审核 - POST /gopay/v1/incoming/apply/review

请求:
```json
{
  "id": 3,
  "action": "pass",
  "remark": "资料齐全，审核通过"
}
```
> action: "pass"(status 1→2) 或 "reject"(status 1→3)
> 自动记录 reviewer, review_time

### 5.5 进件记录列表 - POST /gopay/v1/incoming/record/list

**前端文件**: `src/views/incoming/record.vue`

**前端 TypeScript 接口 (比 apply 多几个审核字段):**
```typescript
interface RecordItem {
  id: number
  applyNo: string
  merchantName: string
  channelType: 'alipay' | 'wechat'
  merchantNo: string
  licenseNo: string
  legalPerson: string
  phone: string
  status: number           // 仅 2(已通过) 和 3(已驳回)
  reviewer: string         // 审核人
  reviewRemark: string     // 审核意见
  ctime: string
  reviewTime: string       // 审核时间
}
```

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "merchantName": "",
  "channelType": "",
  "status": -1,
  "reviewDate": ""
}
```
> 此接口仅返回 status=2 或 status=3 的记录

### 5.6 进件详情 - POST /gopay/v1/incoming/record/detail

请求: `{ "id": 1 }`
响应: 同 RecordItem 完整字段

### 5.7 图片上传 - POST /gopay/v1/upload/image

请求: multipart/form-data
- `file`: 图片文件 (max 5MB, image/*)

响应:
```json
{
  "code": 0,
  "data": { "url": "https://cdn.example.com/uploads/2026/04/abc123.jpg" }
}
```

---

## 六、支付通道配置

### 6.1 通道列表 - POST /gopay/v1/payment/channel/list

**前端文件**: `src/views/payment/channel.vue`

**前端 TypeScript 接口:**
```typescript
interface ChannelItem {
  id: number
  name: string
  code: string
  type: 'alipay' | 'wechat'
  merchantId: number
  merchantName: string
  payMethods: string[]       // ["qrcode", "page", "wap"]
  feeRate: number            // 0.6 (百分比)
  status: number
  remark: string
  ctime: string
}
```

**payMethods 枚举:**
| 值 | 标签 |
|---|------|
| qrcode | 扫码支付 |
| page | 网页支付 |
| wap | WAP支付 |
| app | APP支付 |
| jsapi | JSAPI支付 |
| miniapp | 小程序支付 |

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "name": "",
  "code": "",
  "type": "",
  "status": -1
}
```

响应:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "支付宝当面付",
        "code": "alipay_face",
        "type": "alipay",
        "merchantId": 1,
        "merchantName": "星辰科技有限公司",
        "payMethods": ["qrcode", "page"],
        "feeRate": 0.6,
        "status": 1,
        "remark": "",
        "ctime": "2026-03-16 10:00:00"
      }
    ],
    "total": 6
  }
}
```

> 注意: 数据库存储 payMethods 为逗号分隔字符串 "qrcode,page"，响应时转为数组

### 6.2 新增通道 - POST /gopay/v1/payment/channel/add

请求:
```json
{
  "name": "支付宝当面付",
  "code": "alipay_face",
  "type": "alipay",
  "merchantId": 1,
  "payMethods": ["qrcode", "page"],
  "feeRate": 0.6,
  "remark": ""
}
```
> name, code, type, merchantId, payMethods 必填
> code 全局唯一

### 6.3 编辑通道 - POST /gopay/v1/payment/channel/update

请求: 同 add + id (code 不可修改)

### 6.4 切换通道状态 - POST /gopay/v1/payment/channel/toggleStatus

请求: `{ "id": 1 }`

### 6.5 通道参数配置 - POST /gopay/v1/payment/channel/config

**支付宝通道参数:**
```json
{
  "channelId": 1,
  "appId": "2021000122672388",
  "privateKey": "MIIEvQIBADANBg...",
  "publicKey": "MIIBIjANBgkqhk...",
  "notifyUrl": "https://example.com/alipay/notify",
  "signType": "RSA2",
  "sandbox": true
}
```

**微信通道参数:**
```json
{
  "channelId": 3,
  "appId": "wx1234567890",
  "mchId": "1600123456",
  "apiKey": "your-apiv3-key",
  "serialNo": "cert-serial-no",
  "privateKey": "-----BEGIN RSA PRIVATE KEY-----...",
  "notifyUrl": "https://example.com/wechat/notify"
}
```

### 6.6 通道详情 - POST /gopay/v1/payment/channel/detail

请求: `{ "id": 1 }`

响应:
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "支付宝当面付",
    "code": "alipay_face",
    "type": "alipay",
    "merchantId": 1,
    "merchantName": "星辰科技有限公司",
    "payMethods": ["qrcode", "page"],
    "feeRate": 0.6,
    "status": 1,
    "remark": "",
    "ctime": "2026-03-16 10:00:00",
    "config": {
      "appId": "2021000122672388",
      "privateKey": "******",
      "publicKey": "******",
      "notifyUrl": "https://example.com/alipay/notify",
      "signType": "RSA2",
      "sandbox": true
    }
  }
}
```
> 敏感字段 privateKey, publicKey, apiKey 脱敏为 "******"

---

## 七、订单中心 - 支付订单

### 7.1 支付订单列表 - POST /gopay/v1/order/payment/list

**前端文件**: `src/views/order/payment.vue`

**前端 TypeScript 接口:**
```typescript
interface PaymentOrder {
  id: number
  orderNo: string           // 平台订单号 PAY20260401100001
  outTradeNo: string        // 商户订单号 (UUID去连字符)
  tradeNo: string           // 通道交易号
  merchantId: number
  merchantName: string
  amount: number            // 金额(分)
  channelType: 'alipay' | 'wechat'
  payMethod: string         // qrcode/page/wap/app/jsapi/miniapp
  status: number            // 0-待支付 1-支付成功 2-支付失败 3-已关闭
  subject: string           // 商品描述
  clientIp: string
  notified: boolean         // 是否已回调通知
  remark: string
  ctime: string
  payTime: string           // 支付时间 (可能为空)
}
```

**payMethod 标签映射:**
| 值 | 标签 |
|---|------|
| qrcode | 扫码支付 |
| page | 网页支付 |
| wap | WAP支付 |
| app | APP支付 |
| jsapi | JSAPI |
| miniapp | 小程序 |

**status 映射:**
| 值 | 标签 | Chip 颜色 | 可用操作 |
|---|------|----------|---------|
| 0 | 待支付 | amber | 查看详情, 关闭订单 |
| 1 | 支付成功 | green | 查看详情, 发起退款 |
| 2 | 支付失败 | red | 查看详情 |
| 3 | 已关闭 | grey | 查看详情 |

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "orderNo": "",
  "merchantName": "",
  "status": -1,
  "channelType": "",
  "date": ""
}
```

响应:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "orderNo": "PAY20260401100001",
        "outTradeNo": "a1b2c3d4e5f6",
        "tradeNo": "2026040122001435481234",
        "merchantId": 1,
        "merchantName": "星辰科技有限公司",
        "amount": 29900,
        "channelType": "alipay",
        "payMethod": "qrcode",
        "status": 1,
        "subject": "星辰科技-商品购买",
        "clientIp": "192.168.1.100",
        "notified": true,
        "remark": "",
        "ctime": "2026-04-01 10:30:00",
        "payTime": "2026-04-01 10:31:05"
      }
    ],
    "total": 8
  }
}
```

### 7.2 支付订单详情 - POST /gopay/v1/order/payment/detail
请求: `{ "id": 1 }`
响应: 同列表单条数据

### 7.3 关闭订单 - POST /gopay/v1/order/payment/close
请求: `{ "id": 1 }`
> 仅 status=0(待支付) 可关闭 → status=3

### 7.4 发起退款 - POST /gopay/v1/order/payment/refund
请求:
```json
{
  "id": 1,
  "amount": 29900,
  "reason": "用户申请退款"
}
```
> 仅 status=1(支付成功) 可退款
> 创建 refund_order 记录，调用支付宝/微信退款 API

### 7.5 导出 - POST /gopay/v1/order/payment/export
请求: 同 list 筛选条件 (不含分页)
响应: Content-Type: application/octet-stream, CSV 文件

---

## 八、订单中心 - 退款订单

### 8.1 退款订单列表 - POST /gopay/v1/order/refund/list

**前端文件**: `src/views/order/refund.vue`

**前端 TypeScript 接口:**
```typescript
interface RefundOrder {
  id: number
  refundNo: string          // REF20260402001
  orderNo: string           // 原支付订单号
  tradeRefundNo: string     // 通道退款单号
  merchantName: string
  refundAmount: number      // 退款金额(分)
  orderAmount: number       // 原订单金额(分)
  channelType: 'alipay' | 'wechat'
  status: number            // 0-退款中 1-退款成功 2-退款失败
  reason: string
  operator: string
  ctime: string
  finishTime: string        // 完成时间 (可能为空)
}
```

**status 映射:**
| 值 | 标签 | Chip 颜色 |
|---|------|----------|
| 0 | 退款中 | amber |
| 1 | 退款成功 | green |
| 2 | 退款失败 | red |

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "refundNo": "",
  "orderNo": "",
  "status": -1,
  "channelType": ""
}
```

响应:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "refundNo": "REF20260402001",
        "orderNo": "PAY20260401100001",
        "tradeRefundNo": "2026040200001234",
        "merchantName": "星辰科技有限公司",
        "refundAmount": 29900,
        "orderAmount": 29900,
        "channelType": "alipay",
        "status": 1,
        "reason": "用户申请退款",
        "operator": "admin",
        "ctime": "2026-04-02 09:00:00",
        "finishTime": "2026-04-02 09:05:00"
      }
    ],
    "total": 5
  }
}
```

### 8.2 退款详情 - POST /gopay/v1/order/refund/detail
请求: `{ "id": 1 }`

---

## 九、订单中心 - 转账订单

### 9.1 转账订单列表 - POST /gopay/v1/order/transfer/list

**前端文件**: `src/views/order/transfer.vue`

**前端 TypeScript 接口:**
```typescript
interface TransferOrder {
  id: number
  transferNo: string        // TRF20260401001
  tradeTransferNo: string   // 通道转账单号
  merchantId: number
  merchantName: string
  amount: number            // 转账金额(分)
  channelType: 'alipay' | 'wechat'
  payeeType: string         // account/openid/phone
  payeeAccount: string      // 收款账号
  payeeName: string         // 收款人
  status: number            // 0-处理中 1-成功 2-失败
  remark: string
  ctime: string
  finishTime: string
}
```

**payeeType 映射:**
| 值 | 标签 |
|---|------|
| account | 账户 |
| openid | OpenID |
| phone | 手机号 |

**status 映射:**
| 值 | 标签 | Chip 颜色 |
|---|------|----------|
| 0 | 处理中 | amber |
| 1 | 转账成功 | green |
| 2 | 转账失败 | red |

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "transferNo": "",
  "merchantName": "",
  "status": -1,
  "channelType": ""
}
```

### 9.2 发起转账 - POST /gopay/v1/order/transfer/add

请求:
```json
{
  "merchantId": 1,
  "channelType": "alipay",
  "amount": 500000,
  "payeeType": "account",
  "payeeAccount": "alipay_account@example.com",
  "payeeName": "张三",
  "remark": "商户分润转账"
}
```
> 所有字段必填 (remark 除外)
> transferNo 由后端生成: TRF + YYYYMMDD + 3位序号

### 9.3 转账详情 - POST /gopay/v1/order/transfer/detail
请求: `{ "id": 1 }`

---

## 十、交易记录 - 交易流水

### 10.1 流水列表 - POST /gopay/v1/transaction/flow/list

**前端文件**: `src/views/transaction/flow.vue`

**前端 TypeScript 接口:**
```typescript
interface FlowItem {
  id: number
  flowNo: string            // FLW20260401001
  orderNo: string           // 关联订单号
  type: 'pay' | 'refund' | 'transfer'
  merchantName: string
  amount: number            // 交易金额(分)
  channelType: 'alipay' | 'wechat'
  channelFlowNo: string     // 通道流水号
  direction: 'in' | 'out'   // 资金方向
  status: number            // 0-处理中 1-已完成
  remark: string
  ctime: string
}
```

**type 映射:**
| 值 | 标签 | Chip 颜色 |
|---|------|----------|
| pay | 支付 | blue |
| refund | 退款 | red |
| transfer | 转账 | amber |

**direction 映射:**
| 值 | 标签 | 颜色 | 金额前缀 |
|---|------|------|---------|
| in | 收入 | green | + |
| out | 支出 | red | - |

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "flowNo": "",
  "orderNo": "",
  "type": "",
  "channelType": "",
  "date": ""
}
```

### 10.2 流水统计 - POST /gopay/v1/transaction/flow/stats

响应:
```json
{
  "code": 0,
  "data": {
    "incomeTotal": 113490,
    "expenseTotal": 2591700,
    "totalCount": 12
  }
}
```
> incomeTotal: direction=in 且 status=1 的金额总和
> expenseTotal: direction=out 且 status=1 的金额总和

### 10.3 流水详情 - POST /gopay/v1/transaction/flow/detail
请求: `{ "id": 1 }`

---

## 十一、交易记录 - 回调通知

### 11.1 回调列表 - POST /gopay/v1/transaction/callback/list

**前端文件**: `src/views/transaction/callback.vue`

**前端 TypeScript 接口:**
```typescript
interface CallbackItem {
  id: number
  orderNo: string
  type: 'pay' | 'refund' | 'transfer'
  channelType: 'alipay' | 'wechat'
  direction: 'upstream' | 'downstream'
  notifyUrl: string
  status: number            // 0-失败 1-成功 2-待重试
  httpStatus: number
  retryCount: number
  maxRetry: number          // 固定为 5
  requestBody: string       // JSON 字符串
  responseBody: string      // JSON 字符串
  ctime: string
}
```

**type 映射:**
| 值 | 标签 | Chip 颜色 |
|---|------|----------|
| pay | 支付通知 | blue |
| refund | 退款通知 | red |
| transfer | 转账通知 | amber |

**direction 映射:**
| 值 | 标签 | Chip 颜色 |
|---|------|----------|
| upstream | 上游→平台 | purple |
| downstream | 平台→商户 | teal |

**status 映射:**
| 值 | 标签 | Chip 颜色 | 可操作 |
|---|------|----------|--------|
| 0 | 失败 | red | 手动重试 |
| 1 | 成功 | green | - |
| 2 | 待重试 | amber | 手动重试 |

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "orderNo": "",
  "type": "",
  "status": -1,
  "channelType": ""
}
```

### 11.2 回调详情 - POST /gopay/v1/transaction/callback/detail
请求: `{ "id": 1 }`
响应: 包含完整的 requestBody 和 responseBody

### 11.3 手动重试 - POST /gopay/v1/transaction/callback/retry
请求: `{ "id": 1 }`
> 仅 status != 1 时可操作
> retryCount++, 重新发送通知

---

## 十二、对账管理 - 对账单

### 12.1 对账单列表 - POST /gopay/v1/recon/bill/list

**前端文件**: `src/views/reconciliation/bill.vue`

**前端 TypeScript 接口:**
```typescript
interface BillItem {
  id: number
  billDate: string          // "2026-04-01"
  channelType: 'alipay' | 'wechat'
  platformCount: number     // 平台笔数
  platformAmount: number    // 平台金额(分)
  channelCount: number      // 通道笔数
  channelAmount: number     // 通道金额(分)
  diffCount: number         // 差异笔数
  diffAmount: number        // 差异金额(分)
  status: number            // 0-待对账 1-已对账 2-有差异
  ctime: string
}
```

**status 映射:**
| 值 | 标签 | Chip 颜色 | 可操作 |
|---|------|----------|--------|
| 0 | 待对账 | grey | 执行对账 |
| 1 | 已对账 | green | 查看/下载 |
| 2 | 有差异 | red | 查看/下载 |

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "date": "",
  "channelType": "",
  "status": -1
}
```

### 12.2 生成对账单 - POST /gopay/v1/recon/bill/generate
请求:
```json
{ "date": "2026-04-01", "channelType": "alipay" }
```

### 12.3 执行对账 - POST /gopay/v1/recon/bill/reconcile
请求: `{ "id": 1 }`
> 仅 status=0(待对账) 可执行

### 12.4 对账单详情 - POST /gopay/v1/recon/bill/detail
请求: `{ "id": 1 }`

---

## 十三、对账管理 - 对账差异

### 13.1 差异列表 - POST /gopay/v1/recon/diff/list

**前端文件**: `src/views/reconciliation/diff.vue`

**前端 TypeScript 接口:**
```typescript
interface DiffItem {
  id: number
  billDate: string
  orderNo: string
  channelType: 'alipay' | 'wechat'
  diffType: string
  platformAmount: number | null   // 平台金额(分)
  channelAmount: number | null    // 通道金额(分)
  diffAmount: number              // 差异金额(分)
  handleStatus: number            // 0-待处理 1-已处理 2-已忽略
  handleRemark: string
  handler: string
}
```

**diffType 映射:**
| 值 | 标签 | Chip 颜色 |
|---|------|----------|
| platform_only | 平台多单 | amber |
| channel_only | 通道多单 | purple |
| amount_mismatch | 金额不一致 | red |
| status_mismatch | 状态不一致 | teal |

**handleStatus 映射:**
| 值 | 标签 | Chip 颜色 | 可操作 |
|---|------|----------|--------|
| 0 | 待处理 | amber | 处理 |
| 1 | 已处理 | green | - |
| 2 | 已忽略 | grey | - |

### 13.2 处理差异 - POST /gopay/v1/recon/diff/handle
请求:
```json
{
  "id": 1,
  "action": "resolve",
  "remark": "已与通道方确认，以平台数据为准"
}
```
> action: "resolve"(→status=1) 或 "ignore"(→status=2)

### 13.3 差异详情 - POST /gopay/v1/recon/diff/detail
请求: `{ "id": 1 }`

### 13.4 导出差异 - POST /gopay/v1/recon/diff/export
请求: 同 list 筛选条件

---

## 十四、系统管理 - 用户管理

### 14.1 用户列表 - POST /gopay/v1/system/user/list

**前端文件**: `src/views/system/user.vue`

**前端 TypeScript 接口:**
```typescript
interface UserItem {
  id: number
  username: string
  realName: string
  phone: string
  email: string
  role: string              // admin/operator/finance/viewer
  status: number            // 0-禁用 1-正常
  ctime: string
  lastLogin: string
}
```

**role 标签映射:**
| 值 | 标签 |
|---|------|
| admin | 管理员 |
| operator | 运营 |
| finance | 财务 |
| viewer | 只读 |

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "username": "",
  "phone": "",
  "status": -1
}
```

响应:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin",
        "realName": "超级管理员",
        "phone": "13800000001",
        "email": "admin@gopay.com",
        "role": "admin",
        "status": 1,
        "ctime": "2026-03-01 00:00:00",
        "lastLogin": "2026-04-08 09:30:00"
      }
    ],
    "total": 5
  }
}
```

### 14.2 新增用户 - POST /gopay/v1/system/user/add
请求:
```json
{
  "username": "zhangsan",
  "password": "123456",
  "realName": "张三",
  "phone": "13800138001",
  "email": "zhangsan@gopay.com",
  "role": "operator"
}
```
> username, password 必填
> password 存储时 bcrypt 加密

### 14.3 编辑用户 - POST /gopay/v1/system/user/update
请求:
```json
{
  "id": 2,
  "realName": "张三",
  "phone": "13800138001",
  "email": "zhangsan@gopay.com",
  "role": "operator"
}
```
> username 不可修改

### 14.4 切换用户状态 - POST /gopay/v1/system/user/toggleStatus
请求: `{ "id": 2 }`

### 14.5 重置密码 - POST /gopay/v1/system/user/resetPwd
请求: `{ "id": 2 }`
> 重置为默认密码 (如 "123456") 并 bcrypt 加密

---

## 十五、系统管理 - 角色管理

### 15.1 角色列表 - POST /gopay/v1/system/role/list

**前端文件**: `src/views/system/role.vue`

**前端 TypeScript 接口:**
```typescript
interface RoleItem {
  id: number
  code: string
  name: string
  description: string
  userCount: number         // 使用该角色的用户数
  status: number            // 0-停用 1-启用
  builtIn: boolean          // 内置角色不可停用
  perms: string[]           // 权限标识列表
  ctime: string
}
```

**权限组定义 (前端写死):**
```typescript
const permGroups = [
  { key: 'merchant', label: '商户管理', items: [
    { value: 'merchant:list', label: '商户列表' },
    { value: 'merchant:app', label: '应用管理' },
    { value: 'merchant:edit', label: '商户编辑' },
  ]},
  { key: 'incoming', label: '进件管理', items: [
    { value: 'incoming:apply', label: '进件申请' },
    { value: 'incoming:record', label: '进件记录' },
    { value: 'incoming:review', label: '进件审核' },
  ]},
  { key: 'payment', label: '支付配置', items: [
    { value: 'payment:channel', label: '通道管理' },
    { value: 'payment:config', label: '通道配置' },
  ]},
  { key: 'order', label: '订单中心', items: [
    { value: 'order:payment', label: '支付订单' },
    { value: 'order:refund', label: '退款订单' },
    { value: 'order:transfer', label: '转账订单' },
  ]},
  { key: 'transaction', label: '交易记录', items: [
    { value: 'transaction:flow', label: '交易流水' },
    { value: 'transaction:callback', label: '回调记录' },
  ]},
  { key: 'recon', label: '对账管理', items: [
    { value: 'recon:bill', label: '对账单' },
    { value: 'recon:diff', label: '对账差异' },
  ]},
  { key: 'system', label: '系统管理', items: [
    { value: 'system:user', label: '用户管理' },
    { value: 'system:role', label: '角色管理' },
    { value: 'system:log', label: '操作日志' },
  ]},
]
```

响应:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "code": "admin",
        "name": "管理员",
        "description": "系统管理员，拥有全部权限",
        "userCount": 1,
        "status": 1,
        "builtIn": true,
        "perms": ["merchant:list", "merchant:app", "merchant:edit", "...所有权限"],
        "ctime": "2026-03-01 00:00:00"
      }
    ],
    "total": 4
  }
}
```

### 15.2 新增角色 - POST /gopay/v1/system/role/add
请求:
```json
{
  "code": "custom_role",
  "name": "自定义角色",
  "description": "测试角色"
}
```
> code 全局唯一

### 15.3 编辑角色 - POST /gopay/v1/system/role/update
请求: `{ "id": 5, "name": "...", "description": "..." }`
> code 不可修改

### 15.4 切换角色状态 - POST /gopay/v1/system/role/toggleStatus
请求: `{ "id": 5 }`
> builtIn=true 的角色不可停用，返回错误

### 15.5 更新角色权限 - POST /gopay/v1/system/role/perms/update
请求:
```json
{
  "roleId": 2,
  "perms": ["merchant:list", "merchant:app", "order:payment", "order:refund"]
}
```
> 全量替换: 先删除该角色所有权限，再批量插入新权限

### 15.6 查询角色权限 - POST /gopay/v1/system/role/perms/list
请求: `{ "roleId": 2 }`
响应:
```json
{
  "code": 0,
  "data": ["merchant:list", "merchant:app", "order:payment", "order:refund"]
}
```

---

## 十六、系统管理 - 操作日志

### 16.1 日志列表 - POST /gopay/v1/system/log/list

**前端文件**: `src/views/system/log.vue`

**前端 TypeScript 接口:**
```typescript
interface LogItem {
  id: number
  operator: string
  module: string
  action: string
  description: string
  ip: string
  userAgent: string
  success: boolean
  duration: number          // 耗时(ms)
  requestData: string       // JSON 字符串
  ctime: string
}
```

**module 映射:**
| 值 | 标签 |
|---|------|
| auth | 登录认证 |
| merchant | 商户管理 |
| incoming | 进件管理 |
| payment | 支付配置 |
| order | 订单中心 |
| system | 系统管理 |

**action 映射 & Chip 颜色:**
| 值 | 标签 | 颜色 |
|---|------|------|
| login | 登录 | blue |
| create | 新增 | green |
| update | 修改 | amber |
| delete | 删除 | red |
| export | 导出 | purple |

请求:
```json
{
  "page": 1,
  "pageSize": 20,
  "operator": "",
  "module": "",
  "action": "",
  "date": ""
}
```

响应:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "operator": "admin",
        "module": "auth",
        "action": "login",
        "description": "管理员登录系统",
        "ip": "192.168.1.100",
        "userAgent": "Mozilla/5.0 ...",
        "success": true,
        "duration": 45,
        "requestData": "{\"username\":\"admin\"}",
        "ctime": "2026-04-08 09:30:00"
      }
    ],
    "total": 10
  }
}
```

### 16.2 日志详情 - POST /gopay/v1/system/log/detail
请求: `{ "id": 1 }`

### 16.3 导出日志 - POST /gopay/v1/system/log/export
请求: 同 list 筛选条件

---

## 附录 A: 全部接口汇总表

| # | 方法 | 路径 | 模块 | 说明 | 状态 |
|---|------|------|------|------|------|
| 1 | POST | /gopay/v1/sso/login | 认证 | 登录 | 已有(需扩展) |
| 2 | POST | /gopay/v1/user/getInfo | 认证 | 获取用户信息 | 已有(需扩展) |
| 3 | POST | /gopay/v1/user/changePwd | 认证 | 修改密码 | 新增 |
| 4 | POST | /gopay/v1/user/profile | 认证 | 更新个人资料 | 新增 |
| 5 | POST | /gopay/v1/dashboard/stats | 仪表盘 | 统计数据 | 新增 |
| 6 | POST | /gopay/v1/dashboard/recentOrders | 仪表盘 | 最近订单 | 新增 |
| 7 | POST | /gopay/v1/dashboard/channelDistribution | 仪表盘 | 通道分布 | 新增 |
| 8 | POST | /gopay/v1/dashboard/trend | 仪表盘 | 近7天趋势 | 新增 |
| 9 | POST | /gopay/v1/merchant/list | 商户 | 商户列表 | 新增 |
| 10 | POST | /gopay/v1/merchant/add | 商户 | 新增商户 | 新增 |
| 11 | POST | /gopay/v1/merchant/update | 商户 | 编辑商户 | 新增 |
| 12 | POST | /gopay/v1/merchant/toggleStatus | 商户 | 切换状态 | 新增 |
| 13 | POST | /gopay/v1/merchant/options | 商户 | 下拉选项 | 新增 |
| 14 | POST | /gopay/v1/merchant/app/list | 商户应用 | 应用列表 | 新增 |
| 15 | POST | /gopay/v1/merchant/app/add | 商户应用 | 新增应用 | 新增 |
| 16 | POST | /gopay/v1/merchant/app/update | 商户应用 | 编辑应用 | 新增 |
| 17 | POST | /gopay/v1/incoming/apply/list | 进件 | 申请列表 | 新增 |
| 18 | POST | /gopay/v1/incoming/apply/add | 进件 | 新建申请 | 新增 |
| 19 | POST | /gopay/v1/incoming/apply/submit | 进件 | 提交审核 | 新增 |
| 20 | POST | /gopay/v1/incoming/apply/review | 进件 | 审核 | 新增 |
| 21 | POST | /gopay/v1/incoming/record/list | 进件 | 记录列表 | 新增 |
| 22 | POST | /gopay/v1/incoming/record/detail | 进件 | 记录详情 | 新增 |
| 23 | POST | /gopay/v1/upload/image | 通用 | 图片上传 | 新增 |
| 24 | POST | /gopay/v1/payment/channel/list | 支付通道 | 通道列表 | 新增 |
| 25 | POST | /gopay/v1/payment/channel/add | 支付通道 | 新增通道 | 新增 |
| 26 | POST | /gopay/v1/payment/channel/update | 支付通道 | 编辑通道 | 新增 |
| 27 | POST | /gopay/v1/payment/channel/toggleStatus | 支付通道 | 切换状态 | 新增 |
| 28 | POST | /gopay/v1/payment/channel/detail | 支付通道 | 通道详情 | 新增 |
| 29 | POST | /gopay/v1/payment/channel/config | 支付通道 | 参数配置 | 新增 |
| 30 | POST | /gopay/v1/order/payment/list | 订单 | 支付订单列表 | 新增 |
| 31 | POST | /gopay/v1/order/payment/detail | 订单 | 支付订单详情 | 新增 |
| 32 | POST | /gopay/v1/order/payment/close | 订单 | 关闭订单 | 新增 |
| 33 | POST | /gopay/v1/order/payment/refund | 订单 | 发起退款 | 新增 |
| 34 | POST | /gopay/v1/order/payment/export | 订单 | 导出支付订单 | 新增 |
| 35 | POST | /gopay/v1/order/refund/list | 退款 | 退款订单列表 | 新增 |
| 36 | POST | /gopay/v1/order/refund/detail | 退款 | 退款详情 | 新增 |
| 37 | POST | /gopay/v1/order/transfer/list | 转账 | 转账订单列表 | 新增 |
| 38 | POST | /gopay/v1/order/transfer/add | 转账 | 发起转账 | 新增 |
| 39 | POST | /gopay/v1/order/transfer/detail | 转账 | 转账详情 | 新增 |
| 40 | POST | /gopay/v1/transaction/flow/list | 交易 | 流水列表 | 新增 |
| 41 | POST | /gopay/v1/transaction/flow/detail | 交易 | 流水详情 | 新增 |
| 42 | POST | /gopay/v1/transaction/flow/stats | 交易 | 流水统计 | 新增 |
| 43 | POST | /gopay/v1/transaction/callback/list | 回调 | 回调列表 | 新增 |
| 44 | POST | /gopay/v1/transaction/callback/detail | 回调 | 回调详情 | 新增 |
| 45 | POST | /gopay/v1/transaction/callback/retry | 回调 | 手动重试 | 新增 |
| 46 | POST | /gopay/v1/recon/bill/list | 对账 | 对账单列表 | 新增 |
| 47 | POST | /gopay/v1/recon/bill/detail | 对账 | 对账单详情 | 新增 |
| 48 | POST | /gopay/v1/recon/bill/generate | 对账 | 生成对账单 | 新增 |
| 49 | POST | /gopay/v1/recon/bill/reconcile | 对账 | 执行对账 | 新增 |
| 50 | POST | /gopay/v1/recon/diff/list | 对账差异 | 差异列表 | 新增 |
| 51 | POST | /gopay/v1/recon/diff/detail | 对账差异 | 差异详情 | 新增 |
| 52 | POST | /gopay/v1/recon/diff/handle | 对账差异 | 处理差异 | 新增 |
| 53 | POST | /gopay/v1/recon/diff/export | 对账差异 | 导出差异 | 新增 |
| 54 | POST | /gopay/v1/system/user/list | 系统 | 用户列表 | 新增 |
| 55 | POST | /gopay/v1/system/user/add | 系统 | 新增用户 | 新增 |
| 56 | POST | /gopay/v1/system/user/update | 系统 | 编辑用户 | 新增 |
| 57 | POST | /gopay/v1/system/user/toggleStatus | 系统 | 切换用户状态 | 新增 |
| 58 | POST | /gopay/v1/system/user/resetPwd | 系统 | 重置密码 | 新增 |
| 59 | POST | /gopay/v1/system/role/list | 系统 | 角色列表 | 新增 |
| 60 | POST | /gopay/v1/system/role/add | 系统 | 新增角色 | 新增 |
| 61 | POST | /gopay/v1/system/role/update | 系统 | 编辑角色 | 新增 |
| 62 | POST | /gopay/v1/system/role/toggleStatus | 系统 | 切换角色状态 | 新增 |
| 63 | POST | /gopay/v1/system/role/perms/update | 系统 | 更新权限 | 新增 |
| 64 | POST | /gopay/v1/system/role/perms/list | 系统 | 查询权限 | 新增 |
| 65 | POST | /gopay/v1/system/log/list | 系统 | 操作日志列表 | 新增 |
| 66 | POST | /gopay/v1/system/log/detail | 系统 | 日志详情 | 新增 |
| 67 | POST | /gopay/v1/system/log/export | 系统 | 导出日志 | 新增 |

**合计: 67 个接口 (2 个已有需扩展 + 65 个新增)**

---

## 附录 B: 前端页面与接口对应关系

| 前端页面路径 | 前端文件 | 依赖接口编号 |
|-------------|---------|-------------|
| /login | src/views/login/index.vue | #1 |
| /dashboard | src/views/dashboard/index.vue | #2, #5, #6, #7, #8 |
| /merchant/list | src/views/merchant/list.vue | #9, #10, #11, #12 |
| /merchant/app | src/views/merchant/app.vue | #13, #14, #15, #16 |
| /incoming/apply | src/views/incoming/apply.vue | #13, #17, #18, #19, #20, #23 |
| /incoming/record | src/views/incoming/record.vue | #21, #22 |
| /payment/channel | src/views/payment/channel.vue | #13, #24, #25, #26, #27, #28, #29 |
| /order/payment | src/views/order/payment.vue | #30, #31, #32, #33, #34 |
| /order/refund | src/views/order/refund.vue | #35, #36 |
| /order/transfer | src/views/order/transfer.vue | #13, #37, #38, #39 |
| /transaction/flow | src/views/transaction/flow.vue | #40, #41, #42 |
| /transaction/callback | src/views/transaction/callback.vue | #43, #44, #45 |
| /reconciliation/bill | src/views/reconciliation/bill.vue | #46, #47, #48, #49 |
| /reconciliation/diff | src/views/reconciliation/diff.vue | #50, #51, #52, #53 |
| /system/user | src/views/system/user.vue | #54, #55, #56, #57, #58 |
| /system/role | src/views/system/role.vue | #59, #60, #61, #62, #63, #64 |
| /system/log | src/views/system/log.vue | #65, #66, #67 |
| /profile | src/views/profile/index.vue | #2, #3, #4 |
