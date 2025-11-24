package haozpay

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/go-resty/resty/v2"
)

// signatureMiddleware 请求签名中间件
// 在每个请求发送前自动添加签名字段
//
// 皓臻支付签名算法:
//  1. 收集请求参数(排除sign字段)
//  2. 按参数名ASCII码升序排序
//  3. 按"key=value"格式用&拼接成字符串
//  4. 用SHA256算法生成摘要
//  5. 用商户私钥对摘要进行RSA加密
//
// 参数:
//   - privateKeyPEM: 商户私钥(PEM格式)
//
// 返回:
//   - resty.RequestMiddleware: resty 请求中间件函数
func signatureMiddleware(privateKeyPEM string) resty.RequestMiddleware {
	return func(c *resty.Client, r *resty.Request) error {
		if r.Body == nil {
			return nil
		}

		haozReq, ok := r.Body.(*HaozPayRequest)
		if !ok {
			return nil
		}

		paramsMap := make(map[string]interface{})

		// 展开 bizBody JSON 到 paramsMap
		if haozReq.BizBody != "" {
			var bizBodyMap map[string]interface{}
			if err := json.Unmarshal([]byte(haozReq.BizBody), &bizBodyMap); err != nil {
				return fmt.Errorf("failed to unmarshal bizBody: %w", err)
			}
			// 将 bizBody 中的所有字段添加到 paramsMap
			for k, v := range bizBodyMap {
				paramsMap[k] = v
			}
		}

		// 添加 merchantNo 和 timestamp（使用数字类型，不是字符串）
		paramsMap["merchantNo"] = haozReq.MerchantNo
		paramsMap["timestamp"] = haozReq.Timestamp

		sign, err := GenerateSign(paramsMap, privateKeyPEM)
		if err != nil {
			return fmt.Errorf("failed to generate signature: %w", err)
		}

		haozReq.Sign = sign
		r.SetBody(haozReq)

		return nil
	}
}

// verifyHaozPaySignature 验证皓臻支付回调签名
// 验签算法流程:
//  1. 构建签名字符串(按参数名ASCII升序排序)
//  2. 计算SHA256摘要
//  3. 使用平台公钥解密签名
//  4. 比较解密后的摘要与计算的摘要是否一致
//
// 参数:
//   - publicKeyPEM: 平台公钥(PEM格式)
//   - params: 回调参数(不含sign字段)
//   - signature: Base64编码的签名字符串
//
// 返回:
//   - error: 验签失败时返回错误
func verifyHaozPaySignature(publicKeyPEM string, params map[string]string, signature string) error {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for i, key := range keys {
		value := params[key]
		if value != "" {
			if i > 0 {
				sb.WriteString("&")
			}
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(value)
		}
	}

	paramsStr := sb.String()
	hash := sha256.Sum256([]byte(paramsStr))
	hashHex := fmt.Sprintf("%x", hash)

	publicKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	decrypted, err := decryptWithPublicKey(publicKey, sigBytes)
	if err != nil {
		return fmt.Errorf("failed to decrypt with public key: %w", err)
	}

	if string(decrypted) != hashHex {
		return fmt.Errorf("signature verification failed: hash mismatch")
	}

	return nil
}

// decryptWithPublicKey 使用公钥解密数据
// 这是非标准的RSA用法，但与Java的Hutool库行为一致
// Java的Hutool库实际上是用公钥做"验签"操作（textbook RSA）
func decryptWithPublicKey(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	c := new(big.Int).SetBytes(data)
	if c.Cmp(publicKey.N) >= 0 {
		return nil, fmt.Errorf("message too long")
	}

	// 使用公钥的 E 和 N 进行模幂运算: m = c^e mod n
	m := new(big.Int).Exp(c, big.NewInt(int64(publicKey.E)), publicKey.N)

	// 去除前导零，返回原始数据
	return m.Bytes(), nil
}

// parsePublicKey 解析PEM格式的公钥
// 支持两种格式:
//  1. 完整的 PEM 格式(带 -----BEGIN/END----- 标志)
//  2. 纯 Base64 编码的密钥字符串(不带标志)
func parsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	var keyBytes []byte

	// 尝试 PEM 解码
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block != nil {
		// PEM 格式
		keyBytes = block.Bytes
	} else {
		// 可能是纯 Base64 格式，尝试直接解码
		decoded, err := base64.StdEncoding.DecodeString(publicKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to decode public key: not valid PEM or Base64 format")
		}
		keyBytes = decoded
	}

	pubInterface, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return pubKey, nil
}

// errorHandlerMiddleware 错误处理中间件
// 在接收到响应后检查 HTTP 状态码，如果是错误状态则解析错误信息
//
// 处理逻辑:
//  1. 检查 HTTP 状态码是否 >= 400
//  2. 如果是错误状态，尝试解析响应体中的错误信息
//  3. 将错误信息包装为 SDKError 类型返回
//
// 返回:
//   - resty.ResponseMiddleware: resty 响应中间件函数
func errorHandlerMiddleware() resty.ResponseMiddleware {
	return func(c *resty.Client, r *resty.Response) error {
		// 检查是否为错误状态码
		if r.StatusCode() >= 400 {
			var errResp Response

			// 尝试解析错误响应
			if err := json.Unmarshal(r.Body(), &errResp); err != nil {
				// 解析失败时返回通用错误
				return NewSDKError(
					0,
					"failed to parse error response",
					r.StatusCode(),
				)
			}

			// 返回包含详细信息的 SDK 错误
			return NewSDKErrorWithRequestID(
				errResp.Code,
				errResp.Message,
				r.StatusCode(),
				errResp.RequestID,
			)
		}
		return nil
	}
}

// requestLogMiddleware 请求日志中间件
// 在调试模式下打印请求详情
//
// 打印内容:
//   - 请求方法和 URL
//   - 请求体内容(格式化的 JSON)
//
// 参数:
//   - debug: 是否开启调试模式
//
// 返回:
//   - resty.RequestMiddleware: resty 请求中间件函数
func requestLogMiddleware(debug bool) resty.RequestMiddleware {
	return func(c *resty.Client, r *resty.Request) error {
		if debug {
			// 打印请求行
			fmt.Printf("[SDK Request] %s %s\n", r.Method, r.URL)

			// 打印请求体
			if r.Body != nil {
				bodyBytes, _ := json.MarshalIndent(r.Body, "", "  ")
				fmt.Printf("[SDK Request Body] %s\n", string(bodyBytes))
			}
		}
		return nil
	}
}

// responseLogMiddleware 响应日志中间件
// 在调试模式下打印响应详情
//
// 打印内容:
//   - HTTP 状态码
//   - 请求耗时
//   - 响应体内容
//
// 参数:
//   - debug: 是否开启调试模式
//
// 返回:
//   - resty.ResponseMiddleware: resty 响应中间件函数
func responseLogMiddleware(debug bool) resty.ResponseMiddleware {
	return func(c *resty.Client, r *resty.Response) error {
		if debug {
			// 打印响应状态和耗时
			fmt.Printf("[SDK Response] Status: %d, Time: %v\n",
				r.StatusCode(), r.Time())

			// 打印响应体
			fmt.Printf("[SDK Response Body] %s\n", string(r.Body()))
		}
		return nil
	}
}
