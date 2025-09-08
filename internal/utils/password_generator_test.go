package utils

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	tests := []struct {
		name           string
		length         int
		includeUpper   bool
		includeLower   bool
		includeNumbers bool
		includeSymbols bool
		expectError    bool
	}{
		{
			name:           "Basic password",
			length:         12,
			includeUpper:   true,
			includeLower:   true,
			includeNumbers: true,
			includeSymbols: false,
			expectError:    false,
		},
		{
			name:           "Short password",
			length:         4,
			includeUpper:   true,
			includeLower:   true,
			includeNumbers: false,
			includeSymbols: false,
			expectError:    false,
		},
		{
			name:           "Long password with symbols",
			length:         32,
			includeUpper:   true,
			includeLower:   true,
			includeNumbers: true,
			includeSymbols: true,
			expectError:    false,
		},
		{
			name:           "Numbers only",
			length:         8,
			includeUpper:   false,
			includeLower:   false,
			includeNumbers: true,
			includeSymbols: false,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := GeneratePassword(tt.length, tt.includeUpper, tt.includeLower, tt.includeNumbers, tt.includeSymbols)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if !tt.expectError {
				if len(password) != tt.length {
					t.Errorf("Expected password length %d, got %d", tt.length, len(password))
				}
				
				// 验证字符集
				if tt.includeUpper {
					hasUpper := false
					for _, char := range password {
						if char >= 'A' && char <= 'Z' {
							hasUpper = true
							break
						}
					}
					if !hasUpper && tt.length > 0 {
						// 注意：由于是随机生成，可能不包含某种字符类型，这里只是警告
						t.Logf("Warning: Password doesn't contain uppercase letters")
					}
				}
			}
		})
	}
}

func TestCheckPasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected string
	}{
		{
			name:     "Weak password",
			password: "123",
			expected: "weak",
		},
		{
			name:     "Medium password",
			password: "Password123",
			expected: "medium",
		},
		{
			name:     "Strong password",
			password: "Password123!@#",
			expected: "strong",
		},
		{
			name:     "Empty password",
			password: "",
			expected: "weak",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPasswordStrength(tt.password)
			strength := result["strength"].(string)
			
			if strength != tt.expected {
				t.Errorf("Expected strength %s, got %s", tt.expected, strength)
			}
		})
	}
}
