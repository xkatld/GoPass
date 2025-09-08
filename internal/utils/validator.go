package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) (bool, string) {
	if len(username) < 3 {
		return false, "用户名长度至少需要3个字符"
	}
	if len(username) > 50 {
		return false, "用户名长度不能超过50个字符"
	}
	
	// 只允许字母、数字、下划线和连字符
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return false, "用户名只能包含字母、数字、下划线和连字符"
	}
	
	return true, ""
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string, minLength int, requireStrong bool) (bool, string) {
	if len(password) < minLength {
		return false, fmt.Sprintf("密码长度至少需要%d个字符", minLength)
	}
	
	if !requireStrong {
		return true, ""
	}
	
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		return false, "密码必须包含至少一个大写字母"
	}
	if !hasLower {
		return false, "密码必须包含至少一个小写字母"
	}
	if !hasNumber {
		return false, "密码必须包含至少一个数字"
	}
	if !hasSpecial {
		return false, "密码必须包含至少一个特殊字符"
	}
	
	return true, ""
}

// SanitizeInput 清理输入字符串
func SanitizeInput(input string) string {
	// 移除前后空白字符
	input = strings.TrimSpace(input)
	
	// 移除潜在的危险字符
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#39;")
	input = strings.ReplaceAll(input, "&", "&amp;")
	
	return input
}

// ValidateURL 验证URL格式
func ValidateURL(url string) bool {
	if url == "" {
		return true // 允许空URL
	}
	
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(url)
}

// ValidatePasswordTitle 验证密码条目标题
func ValidatePasswordTitle(title string) (bool, string) {
	title = strings.TrimSpace(title)
	if title == "" {
		return false, "标题不能为空"
	}
	if len(title) > 100 {
		return false, "标题长度不能超过100个字符"
	}
	return true, ""
}

// ValidateCategory 验证分类名称
func ValidateCategory(category string) (bool, string) {
	if category == "" {
		return true, "" // 允许空分类
	}
	
	category = strings.TrimSpace(category)
	if len(category) > 50 {
		return false, "分类名称长度不能超过50个字符"
	}
	
	return true, ""
}
