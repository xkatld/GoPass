package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopass/internal/crypto"
	"gopass/internal/database"
	"gopass/internal/models"

	"github.com/gin-gonic/gin"
)

// ExportData 导出用户数据为CSV
func ExportData(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	// 查询用户的所有密码条目
	rows, err := database.DB.Query(`
		SELECT title, website, username, password, category, notes, created_at
		FROM passwords WHERE user_id = ? ORDER BY created_at DESC`,
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

	// 设置响应头
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=gopass_export_%s.csv", time.Now().Format("20060102_150405")))

	// 创建CSV写入器
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// 写入CSV头部
	writer.Write([]string{"Title", "Website", "Username", "Password", "Category", "Notes", "Created At"})

	// 解密密钥
	encryptionKey := crypto.GenerateKey("user-master-key-" + strconv.Itoa(userID))

	// 写入数据行
	for rows.Next() {
		var title, website, username, encryptedPassword, category, notes string
		var createdAt time.Time

		err := rows.Scan(&title, &website, &username, &encryptedPassword, &category, &notes, &createdAt)
		if err != nil {
			continue
		}

		// 解密密码
		decryptedPassword, err := crypto.Decrypt(encryptedPassword, encryptionKey)
		if err != nil {
			decryptedPassword = "[DECRYPTION_ERROR]"
		}

		writer.Write([]string{
			title,
			website,
			username,
			decryptedPassword,
			category,
			notes,
			createdAt.Format("2006-01-02 15:04:05"),
		})
	}
}

// ImportData 从CSV导入数据
func ImportData(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	// 获取上传的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "No file uploaded",
		})
		return
	}
	defer file.Close()

	// 读取CSV文件
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Failed to parse CSV file",
		})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "CSV file is empty or has no data rows",
		})
		return
	}

	// 验证CSV头部
	header := records[0]
	if len(header) < 4 || !strings.EqualFold(header[0], "Title") || !strings.EqualFold(header[3], "Password") {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid CSV format. Expected headers: Title, Website, Username, Password, Category, Notes, Created At",
		})
		return
	}

	// 加密密钥
	encryptionKey := crypto.GenerateKey("user-master-key-" + strconv.Itoa(userID))

	// 导入数据
	imported := 0
	failed := 0

	for _, record := range records[1:] {
		if len(record) < 4 {
			failed++
			continue
		}

		title := record[0]
		website := ""
		username := ""
		password := record[3]
		category := ""
		notes := ""

		if len(record) > 1 {
			website = record[1]
		}
		if len(record) > 2 {
			username = record[2]
		}
		if len(record) > 4 {
			category = record[4]
		}
		if len(record) > 5 {
			notes = record[5]
		}

		// 验证必填字段
		if title == "" || password == "" {
			failed++
			continue
		}

		// 加密密码
		encryptedPassword, err := crypto.Encrypt(password, encryptionKey)
		if err != nil {
			failed++
			continue
		}

		// 插入数据库
		_, err = database.DB.Exec(`
			INSERT INTO passwords (user_id, title, website, username, password, category, notes, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			userID, title, website, username, encryptedPassword, category, notes, time.Now(), time.Now(),
		)
		if err != nil {
			failed++
			continue
		}

		imported++
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: fmt.Sprintf("Import completed. %d entries imported, %d failed", imported, failed),
		Data: map[string]interface{}{
			"imported": imported,
			"failed":   failed,
		},
	})
}
