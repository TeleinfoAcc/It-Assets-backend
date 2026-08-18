package handlers

import (
	"checklist-backend/database"
	"checklist-backend/models"

	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	var totalCount int64
	var aioCount int64
	var aioUsageCount int64
	var aioOnStockCount int64
	var aioBrokeCount int64
	var notebookCount int64
	var notebookUsageCount int64
	var notebookOnStockCount int64
	var notebookBrokeCount int64

	var projcount []struct {
		Asset_project string `json:"asset_project"`
		Count         int64  `json:"count"`
	}

	if err := database.DB.Model(&models.AbcAssetRent{}).Count(&totalCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ?", "AIO").Count(&aioCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status = ?", "AIO", "1").Count(&aioUsageCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status = ?", "AIO", "0").Count(&aioOnStockCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status IN ?", "AIO", []string{"9", "3"}).Count(&aioBrokeCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ?", "NB").Count(&notebookCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status = ?", "NB", "1").Count(&notebookUsageCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}
	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status = ?", "NB", "0").Count(&notebookOnStockCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}
	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status IN ?", "NB", []string{"9", "3"}).Count(&notebookBrokeCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	if err := database.DB.Model(&models.AbcAssetRent{}).Select("asset_project, COUNT(*) as count").Group("asset_project").Scan(&projcount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"total_assets":            totalCount,
		"aio_count":               aioCount,
		"aio_usage_count":         aioUsageCount,
		"aio_on_stock_count":      aioOnStockCount,
		"aio_broke_count":         aioBrokeCount,
		"notebook_count":          notebookCount,
		"notebook_usage_count":    notebookUsageCount,
		"notebook_on_stock_count": notebookOnStockCount,
		"notebook_broke_count":    notebookBrokeCount,
		"project_count":           projcount,
	})
}
