package handlers

import (
	"checklist-backend/database"
	"checklist-backend/models"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret-key")

func Login(c *gin.Context) {
	var body struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	// หา user จาก DB
	var user models.IamAgents
	if err := database.DB.Where("login = ?", body.Login).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "login  ไม่ถูกต้อง"})
		return
	}

	// เช็ค password
	if user.Password != body.Password {
		c.JSON(401, gin.H{"error": "password ไม่ถูกต้อง"})
		return
	}

	// สร้าง JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"agent_id":      user.Agent_id,
		"login":         user.Login,
		"permission_id": user.Permission_id, // เพิ่มตรงนี้
		"exp":           time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecret)

	c.JSON(200, gin.H{
		"token": tokenString,
		"agent": gin.H{
			"permission_id": user.Permission_id,
			"title_name_th": user.Title_name_th,
			"first_name_th": user.First_name_th,
			"last_name_th":  user.Last_name_th,
			"login":         user.Login,
		},
	})
}
