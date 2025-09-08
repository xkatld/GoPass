package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	LowerChars   = "abcdefghijklmnopqrstuvwxyz"
	UpperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumberChars  = "0123456789"
	SymbolChars  = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// GeneratePassword 生成随机密码
func GeneratePassword(length int, includeUpper, includeLower, includeNumbers, includeSymbols bool) (string, error) {
	if length <= 0 {
		length = 12
	}

	var charset string
	if includeLower {
		charset += LowerChars
	}
	if includeUpper {
		charset += UpperChars
	}
	if includeNumbers {
		charset += NumberChars
	}
	if includeSymbols {
		charset += SymbolChars
	}

	// 如果没有选择任何字符集，默认使用字母和数字
	if charset == "" {
		charset = LowerChars + UpperChars + NumberChars
	}

	var password strings.Builder
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		password.WriteByte(charset[randomIndex.Int64()])
	}

	return password.String(), nil
}

// CheckPasswordStrength 检查密码强度
func CheckPasswordStrength(password string) map[string]interface{} {
	result := map[string]interface{}{
		"length":      len(password),
		"has_lower":   false,
		"has_upper":   false,
		"has_number":  false,
		"has_symbol":  false,
		"strength":    "weak",
		"score":       0,
	}

	score := 0

	// 检查长度
	if len(password) >= 8 {
		score += 1
	}
	if len(password) >= 12 {
		score += 1
	}

	// 检查字符类型
	for _, char := range password {
		switch {
		case strings.ContainsRune(LowerChars, char):
			if !result["has_lower"].(bool) {
				result["has_lower"] = true
				score += 1
			}
		case strings.ContainsRune(UpperChars, char):
			if !result["has_upper"].(bool) {
				result["has_upper"] = true
				score += 1
			}
		case strings.ContainsRune(NumberChars, char):
			if !result["has_number"].(bool) {
				result["has_number"] = true
				score += 1
			}
		case strings.ContainsRune(SymbolChars, char):
			if !result["has_symbol"].(bool) {
				result["has_symbol"] = true
				score += 1
			}
		}
	}

	result["score"] = score

	// 评估强度
	switch {
	case score >= 5:
		result["strength"] = "strong"
	case score >= 3:
		result["strength"] = "medium"
	default:
		result["strength"] = "weak"
	}

	return result
}
