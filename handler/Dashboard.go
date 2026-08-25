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

	var AIOcountbyProj []struct {
		Asset_project string `json:"asset_project"`
		Count         int64  `json:"count"`
	}

	var NBcountbuProj []struct {
		Asset_project string `json:"asset_project"`
		Count         int64  `json:"count"`
	}

	// Count total assets
	if err := database.DB.Model(&models.AbcAssetRent{}).Where("asset_status != ?", "5").Count(&totalCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	// Count AIO assets
	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status != ?", "AIO", "5").Count(&aioCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	// Count AIO assets in use
	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status = ?", "AIO", "1").Count(&aioUsageCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	// Count AIO assets in stock
	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status = ?", "AIO", "0").Count(&aioOnStockCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	// Count AIO assets that are broken
	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status IN ?", "AIO", []string{"9", "3"}).Count(&aioBrokeCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	// Count Notebook assets
	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status != ?", "NB", "5").Count(&notebookCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	// Count Notebook assets in use
	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status = ?", "NB", "1").Count(&notebookUsageCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	// Count Notebook assets in stock
	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status = ?", "NB", "0").Count(&notebookOnStockCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	// Count Notebook assets that are broken
	if err := database.DB.Model(&models.AbcAssetRent{}).Where("com_type = ? AND asset_status IN ?", "NB", []string{"9", "3"}).Count(&notebookBrokeCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	//AIO count by project
	if err := database.DB.Model(&models.AbcAssetRent{}).Select("asset_project, COUNT(*) as count").Where("com_type = ? AND asset_status != ?", "AIO", "5").Group("asset_project").Scan(&AIOcountbyProj).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลได้", "details": err.Error()})
		return
	}

	//Notebook count by project
	if err := database.DB.Model(&models.AbcAssetRent{}).Select("asset_project, COUNT(*) as count").Where("com_type = ? AND asset_status != ?", "NB", "5").Group("asset_project").Scan(&NBcountbuProj).Error; err != nil {
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
		"aio_project_count":       AIOcountbyProj,
		"notebook_project_count":  NBcountbuProj,
	})
}
