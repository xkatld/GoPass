package handlers

import (
	"net/http"
	"strconv"
	"time"

	"gopass/internal/crypto"
	"gopass/internal/database"
	"gopass/internal/models"
	"gopass/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreatePassword 创建密码条目
func CreatePassword(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var req models.PasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	// 验证输入
	if valid, msg := utils.ValidatePasswordTitle(req.Title); !valid {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: msg,
		})
		return
	}

	if req.Website != "" && !utils.ValidateURL(req.Website) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "网站URL格式不正确",
		})
		return
	}

	if valid, msg := utils.ValidateCategory(req.Category); !valid {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: msg,
		})
		return
	}

	// 清理输入
	req.Title = utils.SanitizeInput(req.Title)
	req.Website = utils.SanitizeInput(req.Website)
	req.Username = utils.SanitizeInput(req.Username)
	req.Category = utils.SanitizeInput(req.Category)
	req.Notes = utils.SanitizeInput(req.Notes)

	// 加密密码
	encryptionKey := crypto.GenerateKey("user-master-key-" + strconv.Itoa(userID))
	encryptedPassword, err := crypto.Encrypt(req.Password, encryptionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to encrypt password",
		})
		return
	}

	// 插入密码条目
	result, err := database.DB.Exec(`
		INSERT INTO passwords (user_id, title, website, username, password, category, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.Title, req.Website, req.Username, encryptedPassword, req.Category, req.Notes, time.Now(), time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create password entry",
		})
		return
	}

	passwordID, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Password entry created successfully",
		Data: map[string]interface{}{
			"id": passwordID,
		},
	})
}

// GetPasswords 获取用户的所有密码条目
func GetPasswords(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	category := c.Query("category")
	search := c.Query("search")

	query := `SELECT id, title, website, username, password, category, notes, created_at, updated_at 
			  FROM passwords WHERE user_id = ?`
	args := []interface{}{userID}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	if search != "" {
		query += " AND (title LIKE ? OR website LIKE ? OR username LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	query += " ORDER BY created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}
	defer rows.Close()

	var passwords []models.Password
	encryptionKey := crypto.GenerateKey("user-master-key-" + strconv.Itoa(userID))

	for rows.Next() {
		var p models.Password
		var encryptedPassword string
		err := rows.Scan(&p.ID, &p.Title, &p.Website, &p.Username, &encryptedPassword, &p.Category, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			continue
		}

		// 解密密码
		decryptedPassword, err := crypto.Decrypt(encryptedPassword, encryptionKey)
		if err != nil {
			continue
		}
		p.Password = decryptedPassword
		p.UserID = userID

		passwords = append(passwords, p)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Passwords retrieved successfully",
		Data:    passwords,
	})
}

// GetPassword 获取单个密码条目
func GetPassword(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	passwordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid password ID",
		})
		return
	}

	var p models.Password
	var encryptedPassword string
	err = database.DB.QueryRow(`
		SELECT id, title, website, username, password, category, notes, created_at, updated_at
		FROM passwords WHERE id = ? AND user_id = ?`,
		passwordID, userID,
	).Scan(&p.ID, &p.Title, &p.Website, &p.Username, &encryptedPassword, &p.Category, &p.Notes, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Password entry not found",
		})
		return
	}

	// 解密密码
	encryptionKey := crypto.GenerateKey("user-master-key-" + strconv.Itoa(userID))
	decryptedPassword, err := crypto.Decrypt(encryptedPassword, encryptionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to decrypt password",
		})
		return
	}
	p.Password = decryptedPassword
	p.UserID = userID

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password retrieved successfully",
		Data:    p,
	})
}

// GeneratePassword 生成随机密码
func GeneratePassword(c *gin.Context) {
	var req models.GeneratePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	// 设置默认值
	if req.Length == 0 {
		req.Length = 12
	}

	password, err := utils.GeneratePassword(req.Length, req.IncludeUpper, req.IncludeLower, req.IncludeNumbers, req.IncludeSymbols)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate password",
		})
		return
	}

	strength := utils.CheckPasswordStrength(password)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password generated successfully",
		Data: map[string]interface{}{
			"password": password,
			"strength": strength,
		},
	})
}

// UpdatePassword 更新密码条目
func UpdatePassword(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	passwordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid password ID",
		})
		return
	}

	var req models.PasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	// 加密密码
	encryptionKey := crypto.GenerateKey("user-master-key-" + strconv.Itoa(userID))
	encryptedPassword, err := crypto.Encrypt(req.Password, encryptionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to encrypt password",
		})
		return
	}

	// 更新密码条目
	result, err := database.DB.Exec(`
		UPDATE passwords SET title = ?, website = ?, username = ?, password = ?, category = ?, notes = ?, updated_at = ?
		WHERE id = ? AND user_id = ?`,
		req.Title, req.Website, req.Username, encryptedPassword, req.Category, req.Notes, time.Now(), passwordID, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update password entry",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Password entry not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password entry updated successfully",
	})
}

// DeletePassword 删除密码条目
func DeletePassword(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	passwordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid password ID",
		})
		return
	}

	result, err := database.DB.Exec("DELETE FROM passwords WHERE id = ? AND user_id = ?", passwordID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete password entry",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Password entry not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password entry deleted successfully",
	})
}

// getUserID 从上下文中获取用户ID
func getUserID(c *gin.Context) int {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return 0
	}
	return userID.(int)
}
