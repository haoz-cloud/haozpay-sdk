# haozPay SDK for Go

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

皓臻支付 Go SDK，提供简洁易用的接口集成皓臻支付平台服务。

## ✨ 特性

- 🔐 **安全可靠**: RSA SHA256WithRSA 签名算法，确保请求安全
- 🚀 **简单易用**: 链式配置，简洁的 API 设计
- 📦 **功能完整**: 支持统一下单、订单取消、退款、退款查询
- 🛠 **生产就绪**: 内置重试机制、超时控制、调试模式
- 📝 **文档完善**: 详细的代码注释和使用示例

## 📋 支持的接口

| 接口 | 方法 | 说明 |
|------|------|------|
| 统一下单 | `CreateOrder` | 创建支付订单 |
| 订单取消 | `CancelOrder` | 取消未支付订单 |
| 退款 | `CreateRefund` | 发起退款请求 |
| 退款查询 | `QueryRefund` | 查询退款状态 |
| 回调验证 | `VerifyCallback` | 验证支付/退款回调签名 |

## 📦 安装

### 使用 go get 安装

```bash
go get github.com/haoz-cloud/haozpay-sdk@v1.0.0
```

### 或在 go.mod 中添加依赖

```go
require github.com/haoz-cloud/haozpay-sdk v1.0.0
```

然后执行：

```bash
go mod tidy
```

## 🚀 快速开始

### 1. 初始化客户端

```go
package main

import (
    "context"
    "log"
    
    haozpay "github.com/haoz-cloud/haozpay-sdk"
)

func main() {
    // 配置客户端
    config := haozpay.DefaultConfig().
        WithBaseURL("https://gate.haozpay.com").
        WithMerchantNo("HZ1971294971928846336").
        WithPrivateKey(privateKeyPEM).              // 商户RSA私钥
        WithPlatFormPublicKey(platformPublicKeyPEM) // 平台RSA公钥

    // 创建客户端
    client, err := haozpay.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // 调用支付接口...
}
```

### 2. 统一下单

```go
// 创建支付订单
orderReq := &haozpay.CreatePaymentOrderRequest{
    OrderTitle:        "测试订单",
    OrderAmount:       0.02,
    PayType:           1,                // 1: 微信, 0: 支付宝
    UseHaozPayCashier: true,
    NotifyUrl:         "https://yourdomain.com/callback",
}

order, err := client.Payment.CreateOrder(ctx, orderReq)
if err != nil {
    log.Fatal(err)
}

log.Printf("订单创建成功: %s", order.MerchantOrderNo)
log.Printf("支付信息: %s", order.PayInfo)
```

### 3. 订单取消

```go
cancelReq := &haozpay.CancelPaymentOrderRequest{
    OrderNo:      "ORDER123456",
    CancelReason: "用户取消",
}

err := client.Payment.CancelOrder(ctx, cancelReq)
if err != nil {
    log.Fatal(err)
}

log.Println("订单取消成功")
```

### 4. 退款

```go
refundReq := &haozpay.CreateRefundRequest{
    OrderNo:      "ORDER123456",
    RefundAmount: 0.02,
    RefundReason: "商品问题",
    Remark:       "用户申请退款",
    NotifyUrl:    "https://yourdomain.com/refund-callback",
}

refund, err := client.Payment.CreateRefund(ctx, refundReq)
if err != nil {
    log.Fatal(err)
}

log.Printf("退款申请成功，退款状态: %d", refund.RefundStatus)
```

### 5. 退款查询

```go
queryReq := &haozpay.QueryRefundRequest{
    OrderNo: "ORDER123456",
}

refundStatus, err := client.Payment.QueryRefund(ctx, queryReq)
if err != nil {
    log.Fatal(err)
}

log.Printf("退款状态: %s (代码: %d)",
    refundStatus.RefundStatusDesc,
    refundStatus.RefundStatus)
```

### 6. 回调签名验证

```go
// 处理支付回调
func handlePaymentCallback(w http.ResponseWriter, r *http.Request) {
    // 从 HTTP 请求中获取所有回调参数（除了 sign）
    params := map[string]string{
        "merchantNo": r.FormValue("merchantNo"),
        "orderNo":    r.FormValue("orderNo"),
        "payStatus":  r.FormValue("payStatus"),
        "payAmount":  r.FormValue("payAmount"),
        "timestamp":  r.FormValue("timestamp"),
        // ... 其他回调参数
    }

    // 获取签名
    signature := r.FormValue("sign")

    // 验证回调签名
    if err := client.VerifyCallback(params, signature); err != nil {
        log.Printf("回调签名验证失败: %v", err)
        http.Error(w, "fail", http.StatusBadRequest)
        return
    }

    // 签名验证通过，处理业务逻辑
    log.Println("回调签名验证成功")

    // 更新订单状态等业务逻辑
    // ...

    // 返回成功响应给皓臻支付平台
    w.Write([]byte("success"))
}
```

## 🔐 密钥配置

### 配置密钥

SDK 需要配置以下密钥信息：

1. **商户私钥 (PrivateKey)**: 必填，将生成的私钥通过 `WithPrivateKey()` 配置，用于请求签名
2. **平台公钥 (PlatFormPublicKey)**: 必填，将皓臻支付平台提供的公钥通过 `WithPlatFormPublicKey()` 配置，用于验证回调签名
3. **商户公钥**: 将生成的商户公钥上传到皓臻支付平台管理控制台

### 密钥说明

- **商户私钥**: 用于SDK发起请求时进行签名，确保请求来源可信
- **平台公钥**: 用于验证皓臻支付平台的回调通知签名，防止伪造回调
- **妥善保管**: 商户私钥必须妥善保管，不可泄露

## ⚙️ 高级配置

### 调试模式

```go
config := haozpay.DefaultConfig().
    WithBaseURL("https://gate.haozpay.com").
    WithMerchantNo("HZ1971294971928846336").
    WithPrivateKey(privateKeyPEM).
    WithDebug(true)  // 开启调试模式，打印请求和响应详情
```

### 自定义超时和重试

```go
config := haozpay.DefaultConfig().
    WithBaseURL("https://gate.haozpay.com").
    WithMerchantNo("HZ1971294971928846336").
    WithPrivateKey(privateKeyPEM).
    WithTimeout(60 * time.Second).                           // 60秒超时
    WithRetry(5, 2*time.Second, 10*time.Second)             // 重试5次，等待2-10秒
```

### 代理配置

```go
config := haozpay.DefaultConfig().
    WithBaseURL("https://gate.haozpay.com").
    WithMerchantNo("HZ1971294971928846336").
    WithPrivateKey(privateKeyPEM).
    WithProxy("http://127.0.0.1:8888")  // 设置HTTP代理
```

## 🔧 错误处理

```go
order, err := client.Payment.CreateOrder(ctx, orderReq)
if err != nil {
    // 判断是否为 SDK 错误
    if sdkErr, ok := err.(*haozpay.SDKError); ok {
        log.Printf("错误码: %d", sdkErr.Code)
        log.Printf("错误信息: %s", sdkErr.Message)
        log.Printf("请求ID: %s", sdkErr.RequestID)
        log.Printf("HTTP状态码: %d", sdkErr.StatusCode)
    } else {
        log.Printf("其他错误: %v", err)
    }
    return
}
```

## 📖 API 文档

### 1. 统一下单 (CreateOrder)

#### 请求参数 (CreatePaymentOrderRequest)

| 字段名 | 类型 | 必填 | 说明                             |
|--------|------|--|--------------------------------|
| `OrderTitle` | `string` | ✅ | 订单商品描述                      |
| `OrderAmount` | `float64` | ✅ | 订单金额（单位：元）                     |
| `PayType` | `int` | ✅ | 支付类型：`0` = 支付宝正扫，`2` = JSAPI支付 |
| `UseHaozPayCashier` | `bool` | ✅ | 是否使用皓臻支付收银台交易，必传true           |
| `NotifyUrl` | `string` | ✅ | 支付结果异步通知地址                     |
| `redirectUrl` | `string` |  ❌ | 支付结果异步通知地址                     |

#### 返回参数 (PaymentOrderResponse)

| 字段名 | 类型 | 说明                             |
|--------|------|--------------------------------|
| `MerchantNo` | `string` | 商户号                            |
| `ChannelType` | `string` | 支付渠道类型                         |
| `SeqId` | `string` | 平台订单流水号                        |
| `PayType` | `int` | 支付类型：`0` = 支付宝正扫，`2` = 微信JSAPI |
| `OrderTitle` | `string` | 订单商品描述                         |
| `OrderAmount` | `float64` | 订单金额（单位：元）                     |
| `PayInfo` | `string` | 支付信息（支付宝为表单 HTML，微信为二维码链接）     |
| `MerchantOrderNo` | `string` | 商户订单号                          |

---

### 2. 订单取消 (CancelOrder)

#### 请求参数 (CancelPaymentOrderRequest)

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `OrderNo` | `string` | ✅ | 商户订单号 |
| `CancelReason` | `string` | ❌ | 取消原因 |

#### 返回参数

无返回数据，仅返回成功/失败状态（通过 `error`）

---

### 3. 退款 (CreateRefund)

#### 请求参数 (CreateRefundRequest)

| 字段名            | 类型 | 必填 | 说明 |
|----------------|------|-----|------|
| `OrderNo`      | `string` | ⚠️ | 商户订单号（与 `ReqSeqId` 二选一，不能同时为空） |
| `ReqSeqId`     | `string` | ⚠️ | 原订单请求流水号（与 `OrderNo` 二选一，不能同时为空） |
| `RefundAmount` | `float64` | ✅ | 退款金额（单位：元） |
| `RefundReason` | `string` | ❌ | 退款原因 |
| `Remark`       | `string` | ❌ | 备注信息 |
| `NotifyUrl`    | `string` | ❌ | 退款结果异步通知地址 |

#### 返回参数 (RefundResponse)

| 字段名 | 类型 | 说明                                   |
|--------|------|--------------------------------------|
| `MerchantNo` | `string` | 商户号                                  |
| `OrderNo` | `string` | 商户订单号                                |
| `SeqId` | `string` | 退款流水号                                |
| `ReqDate` | `string` | 请求日期                                 |
| `PaySeqId` | `string` | 原支付流水号                               |
| `PayReqDate` | `string` | 原支付请求日期                              |
| `PayUniqueId` | `string` | 支付唯一标识                               |
| `RefundStartDate` | `string` | 退款开始日期                               |
| `RefundStartTime` | `time.Time` | 退款开始时间                               |
| `RefundFinishTime` | `time.Time` | 退款完成时间                               |
| `RefundStatus` | `int` | 退款状态：`1` = 退款中，`2` = 退款成功，`3` = 退款失败 |
| `RefundAmount` | `float64` | 申请退款金额（单位：元）                         |
| `RealRefundAmount` | `float64` | 实际退款金额（单位：元）                         |
| `TotalRefAmount` | `string` | 原交易累计退款金额（单位：元）                      |
| `TotalRefFeeAmount` | `string` | 原交易累计退款手续费（单位：元）                     |
| `RefCount` | `string` | 累计退款次数                               |

---

### 4. 退款查询 (QueryRefund)

#### 请求参数 (QueryRefundRequest)

| 字段名           | 类型 | 必填 | 说明 |
|---------------|------|-----|------|
| `OrderNo`     | `string` | ✅ | 商户订单号 |
| `RefundSeqId` | `string` | ❌ | 退款请求流水号 |

#### 返回参数 (QueryRefundResponse)

| 字段名 | 类型 | 说明                                        |
|--------|------|-------------------------------------------|
| `MerchantNo` | `string` | 商户号                                       |
| `OrderNo` | `string` | 商户订单号                                     |
| `RefundSeqId` | `string` | 退款请求流水号                                   |
| `PaySeqId` | `string` | 支付请求流水号                                   |
| `PayReqDate` | `string` | 支付请求日期                                    |
| `RefundAmount` | `float64` | 申请退款金额（单位：元）                              |
| `ActualRefundAmount` | `float64` | 实际退款金额（单位：元）                              |
| `RefundStatus` | `int` | 退款状态码：`0` = 初始，`1` = 处理中，`2` = 成功，`3` = 失败 |
| `RefundStatusDesc` | `string` | 退款状态描述                                    |
| `TransFinishTime` | `string` | 交易完成时间，格式：yyyyMMddHHmmss                  |
| `FeeAmount` | `float64` | 手续费金额（单位：元）                               |
| `AcctSplitBunch` | `string` | 分账对象（JSON字符串）                             |
| `UnconfirmAmount` | `float64` | 待确认总金额（单位：元）                              |
| `ConfirmedAmount` | `float64` | 已确认总金额（单位：元）                              |
| `PayChannel` | `string` | 支付渠道：`A` = 支付宝，`T` = 微信，`U` = 银联二维码，`D` = 数字货币                                     |
| `Remark` | `string` | 备注                                        |

---

完整的 API 文档请查看源码注释。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

## 🔗 相关链接

- [皓臻支付文档](https://gate.haozpay.com/docs)
- [GitHub 仓库](https://github.com/haoz-cloud/haozpay-sdk)
- [问题反馈](https://github.com/haoz-cloud/haozpay-sdk/issues)

## ⚠️ 注意事项

1. **生产环境请关闭调试模式**，避免泄露敏感信息
2. **妥善保管商户私钥**，不要提交到代码仓库
3. **建议使用环境变量**存储敏感配置信息
4. **异步回调请验证签名**，防止伪造请求

## 📮 联系方式

如有问题，请提交 [Issue](https://github.com/haoz-cloud/haozpay-sdk/issues) 或联系技术支持。