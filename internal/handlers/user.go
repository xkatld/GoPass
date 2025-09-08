package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"gopass/internal/auth"
	"gopass/internal/crypto"
	"gopass/internal/database"
	"gopass/internal/models"
	"gopass/internal/utils"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	// 验证输入
	if valid, msg := utils.ValidateUsername(req.Username); !valid {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: msg,
		})
		return
	}

	if !utils.ValidateEmail(req.Email) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "邮箱格式不正确",
		})
		return
	}

	if valid, msg := utils.ValidatePassword(req.Password, 8, true); !valid {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: msg,
		})
		return
	}

	// 清理输入
	req.Username = utils.SanitizeInput(req.Username)
	req.Email = utils.SanitizeInput(req.Email)

	// 检查用户名是否已存在
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", req.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "Username already exists",
		})
		return
	}

	// 哈希密码
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to hash password",
		})
		return
	}

	// 插入用户
	result, err := database.DB.Exec(
		"INSERT INTO users (username, password_hash, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		req.Username, hashedPassword, req.Email, time.Now(), time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create user",
		})
		return
	}

	userID, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "User created successfully",
		Data: map[string]interface{}{
			"user_id": userID,
		},
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	// 查询用户
	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, password_hash, email FROM users WHERE username = ?",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	// 验证密码
	if !crypto.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// 生成JWT令牌
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})
}
