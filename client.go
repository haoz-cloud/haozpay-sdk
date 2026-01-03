package haozpay

import (
	"github.com/go-resty/resty/v2"
)

const (
	// SDKVersion SDK 版本号
	SDKVersion = "1.0.0"
	// UserAgent HTTP 请求的 User-Agent 标识
	UserAgent = "haozPay/" + SDKVersion
)

// Client SDK 客户端，提供皓臻支付业务服务的访问入口
// 通过 NewClient 函数创建实例
type Client struct {
	// config SDK 配置信息
	config *Config
	// restyClient 底层 HTTP 客户端
	restyClient *resty.Client

	// Payment 支付服务，提供皓臻支付相关的 API 操作
	// 包含统一下单、订单取消、退款、退款查询、账户提现等功能
	Payment *PaymentService
}

// NewClient 创建并初始化一个新的 SDK 客户端
//
// 参数:
//   - cfg: 客户端配置，包含 API 地址、商户编号、密钥、超时等设置
//
// 返回:
//   - *Client: 初始化完成的客户端实例
//   - error: 配置验证失败时返回错误
//
// 功能说明:
//   - 验证配置的有效性
//   - 创建并配置底层 HTTP 客户端
//   - 注册请求签名和日志中间件
//   - 初始化支付服务
//
// 示例:
//
//	config := sdk.DefaultConfig().
//	    WithBaseURL("https://gate.haozpay.com").
//	    WithMerchantNo("HZ1971294971928846336").
//	    WithPrivateKey(privateKeyPEM).
//	    WithPublicKey(platformPublicKeyPEM)
//
//	client, err := sdk.NewClient(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 使用客户端调用支付接口
//	order, err := client.Payment.CreateOrder(ctx, req)
func NewClient(cfg *Config) (*Client, error) {
	// 验证配置的有效性
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// 创建并配置底层 HTTP 客户端
	restyClient := resty.New().
		SetBaseURL(cfg.BaseURL).                      // 设置 API 基础地址
		SetTimeout(cfg.Timeout).                      // 设置请求超时时间
		SetDebug(cfg.Debug).                          // 设置调试模式
		SetRetryCount(cfg.RetryCount).                // 设置重试次数
		SetRetryWaitTime(cfg.RetryWaitTime).          // 设置重试等待时间
		SetRetryMaxWaitTime(cfg.RetryMaxWait).        // 设置最大重试等待时间
		SetHeader("User-Agent", UserAgent).           // 设置 User-Agent
		SetHeader("Content-Type", "application/json") // 设置内容类型

	// 如果配置了代理，则设置代理
	if cfg.Proxy != "" {
		restyClient.SetProxy(cfg.Proxy)
	}

	// 如果配置了 TLS，则应用 TLS 配置
	if cfg.TLSConfig != nil {
		restyClient.SetTLSClientConfig(cfg.TLSConfig)
	}

	// 注册请求和响应中间件
	restyClient.OnBeforeRequest(requestLogMiddleware(cfg.Debug))    // 请求日志中间件（调试模式时打印请求详情）
	restyClient.OnBeforeRequest(signatureMiddleware(cfg.PrivateKey)) // 请求签名中间件（使用RSA私钥自动签名）
	restyClient.OnAfterResponse(responseLogMiddleware(cfg.Debug))   // 响应日志中间件（调试模式时打印响应详情）
	restyClient.OnAfterResponse(errorHandlerMiddleware())           // 错误处理中间件（统一处理错误响应）

	// 创建客户端实例
	client := &Client{
		config:      cfg,
		restyClient: restyClient,
	}

	// 初始化支付服务
	// PaymentService 提供以下功能：
	//   - CreateOrder: 统一下单
	//   - CancelOrder: 订单取消
	//   - CreateRefund: 退款
	//   - QueryRefund: 退款查询
	//   - CreateWithdraw: 账户提现
	client.Payment = NewPaymentService(client.restyClient, cfg)

	return client, nil
}

// GetConfig 获取客户端的配置信息
//
// 返回:
//   - *Config: 当前客户端使用的配置对象
//
// 示例:
//
//	config := client.GetConfig()
//	fmt.Println("MerchantNo:", config.MerchantNo)
func (c *Client) GetConfig() *Config {
	return c.config
}

// GetRestyClient 获取底层的 resty HTTP 客户端
// 高级用户可以使用此方法获取底层客户端进行自定义操作
//
// 返回:
//   - *resty.Client: resty HTTP 客户端实例
//
// 注意:
//   - 此方法供高级用户使用，一般情况下不需要直接操作底层客户端
//   - 直接使用底层客户端可能会绕过SDK的签名和错误处理机制
//
// 示例:
//
//	restyClient := client.GetRestyClient()
//	resp, err := restyClient.R().Get("/custom/endpoint")
func (c *Client) GetRestyClient() *resty.Client {
	return c.restyClient
}

// VerifyCallback 验证皓臻支付平台回调签名
// 用于在接收到支付/退款回调通知时验证签名的真实性
//
// 参数:
//   - params: 回调参数（不包含 sign 字段），例如从 HTTP 请求中获取的所有参数
//   - signature: 回调中的签名字符串（Base64 编码）
//
// 返回:
//   - error: 验签失败时返回错误，成功返回 nil
//
// 使用场景:
//   - 支付成功回调验签
//   - 退款成功回调验签
//   - 其他异步通知回调验签
//
// 示例:
//
//	// 从 HTTP 请求中获取回调参数
//	params := map[string]string{
//	    "merchantNo":      r.FormValue("merchantNo"),
//	    "orderNo":         r.FormValue("orderNo"),
//	    "payStatus":       r.FormValue("payStatus"),
//	    "payAmount":       r.FormValue("payAmount"),
//	    // ... 其他回调参数，不包括 sign
//	}
//	signature := r.FormValue("sign")
//
//	// 验证签名
//	if err := client.VerifyCallback(params, signature); err != nil {
//	    log.Printf("回调签名验证失败: %v", err)
//	    http.Error(w, "signature verification failed", http.StatusBadRequest)
//	    return
//	}
//
//	// 签名验证通过，处理业务逻辑
//	log.Println("回调签名验证成功")
func (c *Client) VerifyCallback(params map[string]string, signature string) error {
	return verifyHaozPaySignature(c.config.PlatFormPublicKey, params, signature)
}