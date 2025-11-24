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
func BuildSignString(params map[string]interface{}) string {
	if params == nil || len(params) == 0 {
		return ""
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, key := range keys {
		value := params[key]
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

	result := sb.String()
	if len(result) > 0 {
		result = result[:len(result)-1]
	}

	return result
}

// GenerateSign 生成签名
func GenerateSign(params map[string]interface{}, privateKeyStr string) (string, error) {
	signString := BuildSignString(params)

	hash := sha256.Sum256([]byte(signString))
	sha256Hash := fmt.Sprintf("%x", hash)

	privateKey, err := parsePrivateKey(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %w", err)
	}

	signBytes, err := privateKeyEncryptRaw(privateKey, []byte(sha256Hash))
	if err != nil {
		return "", fmt.Errorf("RSA私钥加密失败: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signBytes), nil
}

// privateKeyEncryptRaw 使用私钥进行"加密"（实际是签名操作）
func privateKeyEncryptRaw(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	k := privateKey.Size()

	if len(data) > k-11 {
		return nil, errors.New("数据过长，超过RSA限制")
	}

	em := make([]byte, k)
	em[0] = 0x00
	em[1] = 0x01

	psLen := k - 3 - len(data)
	for i := 2; i < 2+psLen; i++ {
		em[i] = 0xFF
	}

	em[2+psLen] = 0x00
	copy(em[3+psLen:], data)

	m := new(big.Int).SetBytes(em)
	c := new(big.Int).Exp(m, privateKey.D, privateKey.N)

	encrypted := make([]byte, k)
	cBytes := c.Bytes()
	copy(encrypted[k-len(cBytes):], cBytes)

	return encrypted, nil
}

// parsePrivateKey 解析私钥
func parsePrivateKey(keyStr string) (*rsa.PrivateKey, error) {
	keyStr = strings.TrimSpace(keyStr)
	keyStr = normalizePEMFormat(keyStr)

	block, _ := pem.Decode([]byte(keyStr))
	if block == nil {
		return nil, errors.New("私钥PEM格式解析失败")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
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

// normalizePEMFormat 标准化PEM格式
func normalizePEMFormat(keyStr string) string {
	keyStr = strings.TrimSpace(keyStr)

	hasPKCS1Header := strings.Contains(keyStr, "-----BEGIN RSA PRIVATE KEY-----")
	hasPKCS8Header := strings.Contains(keyStr, "-----BEGIN PRIVATE KEY-----")
	hasPKCS1Footer := strings.Contains(keyStr, "-----END RSA PRIVATE KEY-----")
	hasPKCS8Footer := strings.Contains(keyStr, "-----END PRIVATE KEY-----")

	if (hasPKCS1Header && hasPKCS1Footer) || (hasPKCS8Header && hasPKCS8Footer) {
		return keyStr
	}

	keyStr = strings.ReplaceAll(keyStr, "-----BEGIN RSA PRIVATE KEY-----", "")
	keyStr = strings.ReplaceAll(keyStr, "-----END RSA PRIVATE KEY-----", "")
	keyStr = strings.ReplaceAll(keyStr, "-----BEGIN PRIVATE KEY-----", "")
	keyStr = strings.ReplaceAll(keyStr, "-----END PRIVATE KEY-----", "")
	keyStr = strings.TrimSpace(keyStr)

	keyStr = strings.ReplaceAll(keyStr, "\n", "")
	keyStr = strings.ReplaceAll(keyStr, "\r", "")
	keyStr = strings.ReplaceAll(keyStr, " ", "")

	var formatted strings.Builder
	for i := 0; i < len(keyStr); i += 64 {
		end := i + 64
		if end > len(keyStr) {
			end = len(keyStr)
		}
		formatted.WriteString(keyStr[i:end])
		formatted.WriteString("\n")
	}

	return "-----BEGIN RSA PRIVATE KEY-----\n" + formatted.String() + "-----END RSA PRIVATE KEY-----"
}
