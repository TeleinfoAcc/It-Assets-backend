package handlers

import (
	"checklist-backend/database"
	"checklist-backend/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddTools(c *gin.Context) {
	var body struct {
		Tool_name string `json:"tool_name"`
		Quantity  int    `json:"quantity"`
		Type      string `json:"type"`
		Remark    string `json:"remark"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "ข้อมูลไม่ถูกต้อง", "error details": err.Error()})
		return
	}

	// validate ว่าไม่ว่าง
	if body.Tool_name == "" {
		c.JSON(400, gin.H{"error": "กรุณาระบุ Tool Name"})
		return
	}

	// เช็คว่ามีอยู่แล้วไหม โดยดูจาก error โดยตรง
	var existing models.Tools
	err := database.DB.Where("tool_name = ?", body.Tool_name).First(&existing).Error

	if err == nil {
		// err == nil หมายถึงเจอ record → มีอยู่แล้ว
		c.JSON(409, gin.H{"error": "มี Tool นี้อยู่แล้ว"})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// error อื่นที่ไม่ใช่ "ไม่เจอ record" เช่น DB ล่ม
		c.JSON(500, gin.H{"error": "เกิดข้อผิดพลาด", "detailswwss": err.Error()})
		return
	}

	// ไม่เจอ record → สร้างใหม่ได้
	newTool := models.Tools{
		Tool_name: body.Tool_name,
		Is_active: 1,
		Quantity:  body.Quantity,
		Type:      body.Type,
		Remark:    body.Remark,
	}

	if err := database.DB.Create(&newTool).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถเพิ่ม Tool ได้", "details": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "เพิ่มอุปกรณ์สำเร็จ"})
}

func GetTools(c *gin.Context) {

	var tools []models.Tools
	if err := database.DB.Where("is_active = 1 ").Order("tool_name asc").Find(&tools).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูล Tools ได้", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"tools": tools})
}

func UpdateTools(c *gin.Context) {
	var body struct {
		Tool_id   uint   `json:"tool_id"`
		Tool_name string `json:"tool_name"`
		Type      string `json:"type"`
		Quantity  int    `json:"quantity"`
		Remark    string `json:"remark"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	if err := database.DB.Model(&models.Tools{}).Where("tool_id = ?", body.Tool_id).Updates(models.Tools{
		Tool_name: body.Tool_name,
		Type:      body.Type,
		Quantity:  body.Quantity,
		Remark:    body.Remark,
	}).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถอัปเดต Tool ได้", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "อัปเดต Tool สำเร็จ"})
}
