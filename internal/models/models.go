package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Email        string    `json:"email" db:"email"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Password 密码条目模型
type Password struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Website     string    `json:"website" db:"website"`
	Username    string    `json:"username" db:"username"`
	Password    string    `json:"password" db:"password"`
	Category    string    `json:"category" db:"category"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Category 分类模型
type Category struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// PasswordRequest 密码条目请求
type PasswordRequest struct {
	Title    string `json:"title" binding:"required"`
	Website  string `json:"website"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
	Category string `json:"category"`
	Notes    string `json:"notes"`
}

// GeneratePasswordRequest 生成密码请求
type GeneratePasswordRequest struct {
	Length         int  `json:"length" binding:"min=4,max=128"`
	IncludeUpper   bool `json:"include_upper"`
	IncludeLower   bool `json:"include_lower"`
	IncludeNumbers bool `json:"include_numbers"`
	IncludeSymbols bool `json:"include_symbols"`
}

// APIResponse 通用API响应
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
