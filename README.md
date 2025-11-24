# haozPay SDK for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/WeiZzz-D/haozpay-sdk)](https://goreportcard.com/report/github.com/WeiZzz-D/haozpay-sdk)
[![GoDoc](https://godoc.org/github.com/WeiZzz-D/haozpay-sdk?status.svg)](https://godoc.org/github.com/WeiZzz-D/haozpay-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

皓臻支付 Go SDK，提供简洁易用的接口集成皓臻支付平台服务。

## ✨ 特性

- 🔐 **安全可靠**: RSA SHA256WithRSA 签名算法，确保请求安全
- 🚀 **简单易用**: 链式配置，简洁的 API 设计
- 📦 **功能完整**: 支持统一下单、订单取消、退款、退款查询、账户提现
- 🛠 **生产就绪**: 内置重试机制、超时控制、调试模式
- 📝 **文档完善**: 详细的代码注释和使用示例

## 📋 支持的接口

| 接口 | 方法 | 说明 |
|------|------|------|
| 统一下单 | `CreateOrder` | 创建支付订单 |
| 订单取消 | `CancelOrder` | 取消未支付订单 |
| 退款 | `CreateRefund` | 发起退款请求 |
| 退款查询 | `QueryRefund` | 查询退款状态 |
| 账户提现 | `CreateWithdraw` | 商户账户提现 |

## 📦 安装

```bash
go get github.com/WeiZzz-D/haozpay-sdk
```

## 🚀 快速开始

### 1. 初始化客户端

```go
package main

import (
    "context"
    "log"
    
    haozpay "github.com/WeiZzz-D/haozpay-sdk"
)

func main() {
    // 配置客户端
    config := haozpay.DefaultConfig().
        WithBaseURL("https://gate.haozpay.com").
        WithMerchantNo("HZ1971294971928846336").
        WithPrivateKey(privateKeyPEM)  // 商户RSA私钥
    
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
    OrderNo:     "ORDER123456",
    RefundSeqId: "REFUND20251020001",
}

refundStatus, err := client.Payment.QueryRefund(ctx, queryReq)
if err != nil {
    log.Fatal(err)
}

log.Printf("退款状态: %s (代码: %d)", 
    refundStatus.RefundStatusDesc, 
    refundStatus.RefundStatus)
```

### 6. 账户提现

```go
withdrawReq := &haozpay.CreateWithdrawRequest{
    PayChannel:     "HFDG",
    WithdrawAmount: 100.00,
    ReqSeqId:       "TX20251116001",
    Remark:         "商户提现",
    NotifyUrl:      "https://yourdomain.com/withdraw-callback",
}

err := client.Payment.CreateWithdraw(ctx, withdrawReq)
if err != nil {
    log.Fatal(err)
}

log.Println("提现申请成功")
```

## 🔐 密钥配置

### 生成 RSA 密钥对

```bash
# 1. 生成 RSA 私钥 (PKCS#1 格式)
openssl genrsa -out rsa_private_key.pem 2048

# 2. 生成 RSA 公钥
openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem

# 3. 转换为 PKCS#8 格式 (推荐)
openssl pkcs8 -topk8 -in rsa_private_key.pem -out pkcs8_private_key.pem -nocrypt
```

### 配置密钥

1. **商户私钥**: 将生成的私钥通过 `WithPrivateKey()` 配置，用于请求签名
2. **商户公钥**: 将生成的公钥上传到皓臻支付平台控台

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

完整的 API 文档请访问: [GoDoc](https://godoc.org/github.com/WeiZzz-D/haozpay-sdk)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

## 🔗 相关链接

- [皓臻支付官网](https://gate.haozpay.com)
- [皓臻支付文档](https://gate.haozpay.com/docs)
- [问题反馈](https://github.com/WeiZzz-D/haozpay-sdk/issues)

## ⚠️ 注意事项

1. **生产环境请关闭调试模式**，避免泄露敏感信息
2. **妥善保管商户私钥**，不要提交到代码仓库
3. **建议使用环境变量**存储敏感配置信息
4. **异步回调请验证签名**，防止伪造请求

## 📮 联系方式

如有问题，请提交 [Issue](https://github.com/WeiZzz-D/haozpay-sdk/issues)。