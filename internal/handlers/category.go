package handlers

import (
	"net/http"
	"time"

	"gopass/internal/database"
	"gopass/internal/models"

	"github.com/gin-gonic/gin"
)

// GetCategories 获取用户的所有分类
func GetCategories(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	rows, err := database.DB.Query(
		"SELECT id, name, description, created_at FROM categories WHERE user_id = ? ORDER BY name",
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt)
		if err != nil {
			continue
		}
		cat.UserID = userID
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Categories retrieved successfully",
		Data:    categories,
	})
}

// CreateCategory 创建新分类
func CreateCategory(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	// 检查分类名是否已存在
	var exists bool
	err := database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM categories WHERE user_id = ? AND name = ?)",
		userID, req.Name,
	).Scan(&exists)
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
			Message: "Category already exists",
		})
		return
	}

	// 插入新分类
	result, err := database.DB.Exec(
		"INSERT INTO categories (user_id, name, description, created_at) VALUES (?, ?, ?, ?)",
		userID, req.Name, req.Description, time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create category",
		})
		return
	}

	categoryID, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Category created successfully",
		Data: map[string]interface{}{
			"id": categoryID,
		},
	})
}
