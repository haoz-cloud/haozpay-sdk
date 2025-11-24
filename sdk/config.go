package sdk

import (
	"crypto/tls"
	"time"
)

// Config SDK 客户端配置
// 包含 API 连接、认证、超时、重试等所有配置项
type Config struct {
	// BaseURL API 服务的基础地址，例如: https://gate.haozpay.com
	BaseURL string
	// MerchantNo 商户编号，由皓臻支付平台分配
	MerchantNo string
	// PrivateKey 商户RSA私钥(PEM格式)，用于请求签名
	// 需要妥善保管，不可泄露
	PrivateKey string
	// Timeout 单个请求的超时时间，默认 30 秒
	Timeout time.Duration
	// RetryCount 请求失败时的重试次数，默认 3 次
	RetryCount int
	// RetryWaitTime 重试之间的等待时间，默认 1 秒
	RetryWaitTime time.Duration
	// RetryMaxWait 重试的最大等待时间，默认 5 秒
	RetryMaxWait time.Duration
	// Debug 是否开启调试模式，开启后会打印请求和响应详情
	Debug bool
	// Proxy 代理服务器地址，例如: http://proxy.example.com:8080
	Proxy string
	// TLSConfig 自定义 TLS 配置，用于 HTTPS 连接
	TLSConfig *tls.Config
}

// DefaultConfig 创建一个具有默认值的配置对象
//
// 默认值:
//   - Timeout: 30秒
//   - RetryCount: 3次
//   - RetryWaitTime: 1秒
//   - RetryMaxWait: 5秒
//   - Debug: false
//
// 返回:
//   - *Config: 包含默认值的配置对象
//
// 示例:
//
//	config := sdk.DefaultConfig().
//	    WithBaseURL("https://gate.haozpay.com").
//	    WithMerchantNo("HZ1971294971928846336").
//	    WithPrivateKey(privateKeyPEM)
func DefaultConfig() *Config {
	return &Config{
		Timeout:       30 * time.Second,
		RetryCount:    3,
		RetryWaitTime: 1 * time.Second,
		RetryMaxWait:  5 * time.Second,
		Debug:         false,
	}
}

// WithBaseURL 设置 API 基础地址
// 支持链式调用
//
// 参数:
//   - baseURL: API 服务的基础地址，例如 "https://gate.haozpay.com"
//
// 返回:
//   - *Config: 返回自身以支持链式调用
func (c *Config) WithBaseURL(baseURL string) *Config {
	c.BaseURL = baseURL
	return c
}

// WithMerchantNo 设置商户编号
// 支持链式调用
//
// 参数:
//   - merchantNo: 商户编号，由皓臻支付平台分配
//
// 返回:
//   - *Config: 返回自身以支持链式调用
//
// 示例:
//
//	config.WithMerchantNo("HZ1971294971928846336")
func (c *Config) WithMerchantNo(merchantNo string) *Config {
	c.MerchantNo = merchantNo
	return c
}

// WithPrivateKey 设置商户RSA私钥
// 支持链式调用
//
// 参数:
//   - privateKey: 商户RSA私钥(PEM格式)，用于请求签名
//
// 返回:
//   - *Config: 返回自身以支持链式调用
//
// 注意:
//   - 私钥必须是PEM格式
//   - 支持PKCS#1和PKCS#8两种格式
//   - 请妥善保管私钥，不可泄露
//
// 示例:
//
//	privateKeyPEM := `-----BEGIN RSA PRIVATE KEY-----
//	MIIEpAIBAAKCAQEA...
//	-----END RSA PRIVATE KEY-----`
//	config.WithPrivateKey(privateKeyPEM)
func (c *Config) WithPrivateKey(privateKey string) *Config {
	c.PrivateKey = privateKey
	return c
}

// WithTimeout 设置请求超时时间
// 支持链式调用
//
// 参数:
//   - timeout: 超时时间，例如 30*time.Second
//
// 返回:
//   - *Config: 返回自身以支持链式调用
//
// 示例:
//
//	config.WithTimeout(60 * time.Second)
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	return c
}

// WithRetry 设置重试策略
// 支持链式调用
//
// 参数:
//   - count: 重试次数
//   - waitTime: 重试等待时间
//   - maxWait: 最大重试等待时间
//
// 返回:
//   - *Config: 返回自身以支持链式调用
//
// 示例:
//
//	config.WithRetry(5, 2*time.Second, 10*time.Second)
func (c *Config) WithRetry(count int, waitTime, maxWait time.Duration) *Config {
	c.RetryCount = count
	c.RetryWaitTime = waitTime
	c.RetryMaxWait = maxWait
	return c
}

// WithDebug 设置调试模式
// 开启后会在控制台打印详细的请求和响应信息
// 支持链式调用
//
// 参数:
//   - debug: true 开启调试，false 关闭调试
//
// 返回:
//   - *Config: 返回自身以支持链式调用
//
// 注意:
//   - 调试模式会打印敏感信息，生产环境请关闭
//
// 示例:
//
//	config.WithDebug(true)
func (c *Config) WithDebug(debug bool) *Config {
	c.Debug = debug
	return c
}

// WithProxy 设置代理服务器
// 支持链式调用
//
// 参数:
//   - proxy: 代理服务器地址，例如 "http://proxy.example.com:8080"
//
// 返回:
//   - *Config: 返回自身以支持链式调用
//
// 示例:
//
//	config.WithProxy("http://127.0.0.1:8888")
func (c *Config) WithProxy(proxy string) *Config {
	c.Proxy = proxy
	return c
}

// WithTLSConfig 设置自定义 TLS 配置
// 支持链式调用
//
// 参数:
//   - tlsConfig: TLS 配置对象
//
// 返回:
//   - *Config: 返回自身以支持链式调用
//
// 示例:
//
//	tlsConfig := &tls.Config{InsecureSkipVerify: true}
//	config.WithTLSConfig(tlsConfig)
func (c *Config) WithTLSConfig(tlsConfig *tls.Config) *Config {
	c.TLSConfig = tlsConfig
	return c
}

// Validate 验证配置的有效性
// 检查必填字段是否已设置
//
// 返回:
//   - error: 如果配置无效则返回错误，否则返回 nil
//
// 必填字段:
//   - BaseURL: API 基础地址
//   - MerchantNo: 商户编号
//   - PrivateKey: 商户RSA私钥
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return ErrInvalidConfig("BaseURL is required")
	}
	if c.MerchantNo == "" {
		return ErrInvalidConfig("MerchantNo is required")
	}
	if c.PrivateKey == "" {
		return ErrInvalidConfig("PrivateKey is required")
	}
	return nil
}