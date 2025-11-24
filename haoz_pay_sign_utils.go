package haozpay

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
)

// BuildSignString 构建签名字符串
// 参数按字典序升序排列，如果参数值为空字符串则略过
//
// params: 参数Map
// 返回: 签名字符串，格式为: key1=value1&key2=value2
func BuildSignString(params map[string]interface{}) string {
	if params == nil || len(params) == 0 {
		return ""
	}

	// 提取key并排序（字典序）
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 构建签名字符串
	var sb strings.Builder
	for _, key := range keys {
		value := params[key]
		// 跳过sign字段和nil值以及空字符串
		if key == "sign" || value == nil {
			continue
		}
		valueStr := fmt.Sprintf("%v", value)
		if strings.TrimSpace(valueStr) == "" {
			continue
		}

		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(valueStr)
		sb.WriteString("&")
	}

	// 删除最后一个&
	result := sb.String()
	if len(result) > 0 {
		result = result[:len(result)-1]
	}

	return result
}

// GenerateSign 生成签名
// 步骤：
// 1. 构建签名字符串（字典序排序，空值跳过）
// 2. SHA256摘要（转为HEX字符串）
// 3. 使用私钥对SHA256摘要进行RSA加密（兼容Java Hutool的encryptBase64）
// 4. Base64编码
//
// params: 参数Map
// privateKeyStr: 私钥字符串（支持纯私钥字符串或完整PEM格式）
func GenerateSign(params map[string]interface{}, privateKeyStr string) (string, error) {
	// 1. 构建签名字符串
	signString := BuildSignString(params)

	// 2. SHA256摘要，转为HEX字符串
	hash := sha256.Sum256([]byte(signString))
	sha256Hash := fmt.Sprintf("%x", hash)

	// 3. 解析私钥
	privateKey, err := parsePrivateKey(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %w", err)
	}

	// 4. 使用私钥进行RSA加密（兼容Java Hutool的非标准加密方式）
	// Java的Hutool库使用私钥加密时，实际是用私钥进行模运算
	encryptedData, err := encryptWithPrivateKey(privateKey, []byte(sha256Hash))
	if err != nil {
		return "", fmt.Errorf("RSA私钥加密失败: %w", err)
	}

	// 5. Base64编码
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// encryptWithPrivateKey 使用私钥加密（模拟Java Hutool的encryptBase64行为）
// 注意：这是非标准的RSA操作，仅用于兼容Java Hutool库
func encryptWithPrivateKey(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	// 将数据转换为大整数
	m := new(big.Int).SetBytes(data)

	// 检查数据是否超过模数
	if m.Cmp(privateKey.N) >= 0 {
		return nil, errors.New("数据过长，超过RSA模数")
	}

	// 使用私钥的指数进行模运算: c = m^d mod n
	c := new(big.Int).Exp(m, privateKey.D, privateKey.N)

	// 将结果转换为字节数组，填充到密钥长度
	k := (privateKey.N.BitLen() + 7) / 8
	encrypted := make([]byte, k)
	cBytes := c.Bytes()

	// 在前面补零，确保长度一致
	copy(encrypted[k-len(cBytes):], cBytes)

	return encrypted, nil
}

// parsePrivateKey 解析私钥（支持PKCS1和PKCS8格式，自动兼容纯私钥字符串和PEM格式）
func parsePrivateKey(keyStr string) (*rsa.PrivateKey, error) {
	// 去除首尾空白字符
	keyStr = strings.TrimSpace(keyStr)

	// 智能检测并补全PEM格式标志
	keyStr = normalizePEMFormat(keyStr)

	// 解析PEM格式
	block, _ := pem.Decode([]byte(keyStr))
	if block == nil {
		return nil, errors.New("私钥PEM格式解析失败")
	}

	// 尝试PKCS1格式
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// 尝试PKCS8格式
		keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("不支持的私钥格式: %w", err)
		}
		var ok bool
		privateKey, ok = keyInterface.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("不是RSA私钥")
		}
	}

	return privateKey, nil
}

// normalizePEMFormat 标准化PEM格式，自动补全头尾标志
func normalizePEMFormat(keyStr string) string {
	keyStr = strings.TrimSpace(keyStr)

	// 检查是否已经包含PEM格式标志
	hasPKCS1Header := strings.Contains(keyStr, "-----BEGIN RSA PRIVATE KEY-----")
	hasPKCS8Header := strings.Contains(keyStr, "-----BEGIN PRIVATE KEY-----")
	hasPKCS1Footer := strings.Contains(keyStr, "-----END RSA PRIVATE KEY-----")
	hasPKCS8Footer := strings.Contains(keyStr, "-----END PRIVATE KEY-----")

	// 如果已经是完整的PEM格式，直接返回
	if (hasPKCS1Header && hasPKCS1Footer) || (hasPKCS8Header && hasPKCS8Footer) {
		return keyStr
	}

	// 如果只有头部或尾部，先移除它们
	keyStr = strings.ReplaceAll(keyStr, "-----BEGIN RSA PRIVATE KEY-----", "")
	keyStr = strings.ReplaceAll(keyStr, "-----END RSA PRIVATE KEY-----", "")
	keyStr = strings.ReplaceAll(keyStr, "-----BEGIN PRIVATE KEY-----", "")
	keyStr = strings.ReplaceAll(keyStr, "-----END PRIVATE KEY-----", "")
	keyStr = strings.TrimSpace(keyStr)

	// 移除可能存在的换行符，然后重新格式化（每64个字符一行）
	keyStr = strings.ReplaceAll(keyStr, "\n", "")
	keyStr = strings.ReplaceAll(keyStr, "\r", "")
	keyStr = strings.ReplaceAll(keyStr, " ", "")

	// 每64个字符插入换行符（PEM标准格式）
	var formatted strings.Builder
	for i := 0; i < len(keyStr); i += 64 {
		end := i + 64
		if end > len(keyStr) {
			end = len(keyStr)
		}
		formatted.WriteString(keyStr[i:end])
		formatted.WriteString("\n")
	}

	// 默认使用PKCS1格式头尾（兼容性更好）
	return "-----BEGIN RSA PRIVATE KEY-----\n" + formatted.String() + "-----END RSA PRIVATE KEY-----"
}
