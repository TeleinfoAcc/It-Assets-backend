package handlers

import (
	"checklist-backend/database"
	"checklist-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

// func AddTools(c *gin.Context) {
// 	var body struct {
// 		Tool_name string `json:"tool_name"`
// 		Quantity  int    `json:"quantity"`
// 		Type      string `json:"type"`
// 		Remark    string `json:"remark"`
// 	}

// 	if err := c.BindJSON(&body); err != nil {
// 		c.JSON(400, gin.H{"error": "ข้อมูลไม่ถูกต้อง", "error details": err.Error()})
// 		return
// 	}

// 	// validate ว่าไม่ว่าง
// 	if body.Tool_name == "" {
// 		c.JSON(400, gin.H{"error": "กรุณาระบุ Tool Name"})
// 		return
// 	}

// 	// เช็คว่ามีอยู่แล้วไหม โดยดูจาก error โดยตรง
// 	var existing models.Tools
// 	err := database.DB.Where("tool_name = ?", body.Tool_name).First(&existing).Error

// 	if err == nil {
// 		// err == nil หมายถึงเจอ record → มีอยู่แล้ว
// 		c.JSON(409, gin.H{"error": "มี Tool นี้อยู่แล้ว"})
// 		return
// 	}

// 	if !errors.Is(err, gorm.ErrRecordNotFound) {
// 		// error อื่นที่ไม่ใช่ "ไม่เจอ record" เช่น DB ล่ม
// 		c.JSON(500, gin.H{"error": "เกิดข้อผิดพลาด", "detailswwss": err.Error()})
// 		return
// 	}

// 	// ไม่เจอ record → สร้างใหม่ได้
// 	newTool := models.Tools{
// 		Tool_name: body.Tool_name,
// 		Is_active: 1,
// 		Quantity:  body.Quantity,
// 		Type:      body.Type,
// 		Remark:    body.Remark,
// 	}

// 	if err := database.DB.Create(&newTool).Error; err != nil {
// 		c.JSON(500, gin.H{"error": "ไม่สามารถเพิ่ม Tool ได้", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(201, gin.H{"message": "เพิ่มอุปกรณ์สำเร็จ"})
// }

func GetAssetStatus(c *gin.Context) {

	var AssetStatuses []models.AbcAssetStatus
	if err := database.DB.Order("asset_status_name asc").Find(&AssetStatuses).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูล Asset Status ได้", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"asset_statuses": AssetStatuses})
}

func GetSites(c *gin.Context) {

	var Sites []models.AbcSite
	if err := database.DB.Order("site_name asc").Find(&Sites).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูล Sites ได้", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"sites": Sites})
}

func GetRooms(c *gin.Context) {

	var Rooms []models.AbcSiteRoom
	if err := database.DB.Order("room_name asc").Find(&Rooms).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูล Rooms ได้", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"rooms": Rooms})
}

func GetAssets(c *gin.Context) {

	var Assets []models.AbcAsset
	if err := database.DB.Preload("Asset_status_name").Order("com_name asc").Find(&Assets).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูล Assets ได้", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"tools": Assets})
}

func GetAssetsRent(c *gin.Context) {

	var AssetsRent []models.AbcAssetRent
	if err := database.DB.Preload("Asset_status_name").Order("com_name asc").Find(&AssetsRent).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูล Assets Rent ได้", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"assets_rent": AssetsRent})

}

func GetAssetsRentById(c *gin.Context) {

	var AssetsRent models.AbcAssetRent
	id := c.Param("id")
	if err := database.DB.Order("com_name asc").Where("it_asset_id = ?", id).Find(&AssetsRent).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูล Assets Rent ได้", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"assets_rent": AssetsRent})

}

func GetAssetById(c *gin.Context) {

	var Asset models.AbcAsset
	id := c.Param("id")
	if err := database.DB.Where("it_asset_id = ?", id).First(&Asset).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูล Asset ได้", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"tools": Asset})
}

func UpdateAsset(c *gin.Context) {

	agent_id := uint(c.GetFloat64("agent_id"))

	var body struct {
		It_asset_id    uint   `json:"it_asset_id"`
		Gl_asset_code  string `json:"gl_asset_code"`
		Serialnumber   string `json:"serialnumber"`
		Com_name       string `json:"com_name"`
		Com_local_ip   string `json:"com_local_ip"`
		Com_join_ip    string `json:"com_join_ip"`
		Curr_room_code string `json:"curr_room_code"`
		Com_brand      string `json:"com_brand"`
		Com_model      string `json:"com_model"`
		Com_type       string `json:"com_type"`
		Mdf_date       string `json:"mdf_date"`
		Asset_status   uint   `json:"asset_status"`
		Cap_date       string `json:"cap_date"`
		Loc_type       string `json:"loc_type"`
		Location       string `json:"location"`
		Com_desc1      string `json:"com_desc1"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "ข้อมูลไม่ถูกต้อง", "details": err.Error()})
		return
	}

	tx := database.DB.Begin()

	if err := tx.Model(&models.AbcAsset{}).Where("it_asset_id = ?", body.It_asset_id).Updates(models.AbcAsset{
		Com_name:       body.Com_name,
		Com_type:       body.Com_type,
		Gl_asset_code:  body.Gl_asset_code,
		Com_model:      body.Com_model,
		Cap_date:       body.Cap_date,
		Serialnumber:   body.Serialnumber,
		Com_local_ip:   body.Com_local_ip,
		Com_join_ip:    body.Com_join_ip,
		Curr_room_code: body.Curr_room_code,
		Com_brand:      body.Com_brand,
		Mdf_date:       body.Mdf_date,
		Asset_status:   body.Asset_status,
		Loc_type:       body.Loc_type,
		Location:       body.Location,
		Com_desc1:      body.Com_desc1,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "ไม่สามารถอัปเดต Asset ได้", "details": err.Error()})
		return
	}

	if err := tx.Model(&models.AbcHistory{}).Create(&models.AbcHistory{
		Agent_id:     agent_id,
		Asset_status: body.Asset_status,
		Com_desc1:    body.Com_desc1,
		It_asset_id:  body.It_asset_id,
		Timestamp:    time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "ไม่สามารถสร้างประวัติ Asset ได้", "details": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"message": "อัปเดต Asset สำเร็จ"})
}

func UpdateAssetRent(c *gin.Context) {

	agent_id := uint(c.GetFloat64("agent_id"))

	var body struct {
		It_asset_id       uint   `json:"it_asset_id" gorm:"primaryKey"`
		Serialnumber      string `json:"serialnumber"`
		Emp_id            string `json:"emp_id"`
		Com_name          string `json:"com_name"`
		Com_local_ip      string `json:"com_local_ip"`
		Com_join_ip       string `json:"com_join_ip"`
		Curr_room_code    string `json:"curr_room_code"`
		Com_brand         string `json:"com_brand"`
		Com_model         string `json:"com_model"`
		Com_type          string `json:"com_type"`
		Com_desc1         string `json:"com_desc1"`
		Com_desc2         string `json:"com_desc2"`
		Com_desc3         string `json:"com_desc3"`
		Gl_asset_code     string `json:"gl_asset_code"`
		Loc_type          string `json:"loc_type"`
		Create_date       string `json:"create_date"`
		Mdf_date          string `json:"mdf_date"`
		Asset_status      uint   `json:"asset_status"`
		Loc_seat          string `json:"loc_seat"`
		Location          string `json:"location"`
		Com_status        string `json:"com_status"`
		Cap_date          string `json:"cap_date"`
		Mdf_agent_id      uint   `json:"mdf_agent_id"`
		Iss_date          string `json:"iss_date"`
		Return_date       string `json:"return_date"`
		Asset_type        string `json:"asset_type"`
		Asset_project     string `json:"asset_project"`
		Com_hdd           string `json:"com_hdd"`
		Com_wifi_mac      string `json:"com_wifi_mac"`
		Com_lan_mac       string `json:"com_lan_mac"`
		Com_adapt_sn      string `json:"com_adapt_sn"`
		Com_mouse_sn      string `json:"com_mouse_sn"`
		Com_ssd           string `json:"com_ssd"`
		Com_usb_sn        string `json:"com_usb_sn"`
		Asset_status_name string `json:"asset_status_name"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "ข้อมูลไม่ถูกต้อง", "details": err.Error()})
		return
	}

	tx := database.DB.Begin()

	if err := tx.Model(&models.AbcAssetRent{}).Where("it_asset_id = ?", body.It_asset_id).Updates(models.AbcAssetRent{
		Com_name:       body.Com_name,
		Com_type:       body.Com_type,
		Gl_asset_code:  body.Gl_asset_code,
		Com_model:      body.Com_model,
		Cap_date:       body.Cap_date,
		Serialnumber:   body.Serialnumber,
		Com_local_ip:   body.Com_local_ip,
		Com_join_ip:    body.Com_join_ip,
		Curr_room_code: body.Curr_room_code,
		Com_brand:      body.Com_brand,
		Mdf_date:       body.Mdf_date,
		Asset_status:   body.Asset_status,
		Loc_type:       body.Loc_type,
		Location:       body.Location,
		Com_desc1:      body.Com_desc1,
		Com_desc2:      body.Com_desc2,
		Com_desc3:      body.Com_desc3,
		Loc_seat:       body.Loc_seat,
		Com_status:     body.Com_status,
		Mdf_agent_id:   body.Mdf_agent_id,
		Iss_date:       body.Iss_date,
		Return_date:    body.Return_date,
		Asset_type:     body.Asset_type,
		Asset_project:  body.Asset_project,
		Com_hdd:        body.Com_hdd,
		Com_wifi_mac:   body.Com_wifi_mac,
		Com_lan_mac:    body.Com_lan_mac,
		Com_adapt_sn:   body.Com_adapt_sn,
		Com_mouse_sn:   body.Com_mouse_sn,
		Com_ssd:        body.Com_ssd,
		Com_usb_sn:     body.Com_usb_sn,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "ไม่สามารถอัปเดต Asset ได้", "details": err.Error()})
		return
	}

	if err := tx.Model(&models.AbcHistory{}).Create(&models.AbcHistory{
		Agent_id:     agent_id,
		Asset_status: body.Asset_status,
		Com_desc1:    body.Com_desc1,
		It_asset_id:  body.It_asset_id,
		Timestamp:    time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "ไม่สามารถสร้างประวัติ Asset ได้", "details": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"message": "อัปเดต Asset สำเร็จ"})
}

func GetHistory(c *gin.Context) {

	type HistoryResponse struct {
		HisID           uint      `json:"his_id"`
		AgentID         uint      `json:"agent_id"`
		AssetStatus     uint      `json:"asset_status"`
		ComDesc1        string    `json:"com_desc1"`
		ItAssetID       uint      `json:"it_asset_id"`
		ComName         string    `json:"com_name"`
		Serialnumber    string    `json:"serialnumber"`
		FirstNameTH     string    `json:"first_name_th"`
		LastNameTH      string    `json:"last_name_th"`
		AssetStatusName string    `json:"asset_status_name"`
		Timestamp       time.Time `json:"timestamp"`
	}

	var history []HistoryResponse

	err := database.DB.
		Table("abcinv.abc_history AS h").
		Select(`
            h.his_id,
            h.agent_id,
            h.asset_status,
            h.com_desc1,
            h.it_asset_id,
            COALESCE(asset.com_name, asset_rent.com_name) AS com_name,
            COALESCE(asset.serialnumber, asset_rent.serialnumber) AS serialnumber,
            agent.first_name_th,
            agent.last_name_th,
			abc_asset_status.asset_status_name,
			h.timestamp
        `).
		Joins(`
            LEFT JOIN qamon.iam_agents AS agent
                ON agent.agent_id = h.agent_id
        `).
		Joins(`
            LEFT JOIN abcinv.abc_asset AS asset
                ON asset.it_asset_id = h.it_asset_id
        `).
		Joins(`
            LEFT JOIN abcinv.abc_asset_rent AS asset_rent
                ON asset_rent.it_asset_id = h.it_asset_id
        `).
		Joins(`
            LEFT JOIN abcinv.abc_asset_status
                ON abc_asset_status.asset_status = h.asset_status
        `).
		Find(&history).Error

	if err != nil {
		c.JSON(500, gin.H{
			"error":   "ไม่สามารถดึงประวัติ Asset ได้",
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"history": history})
}

func GetHistoryById(c *gin.Context) {

	itAssetID := c.Param("it_asset_id")

	type HistoryResponse struct {
		HisID        uint      `json:"his_id"`
		AgentID      uint      `json:"agent_id"`
		AssetStatus  uint      `json:"asset_status"`
		ComDesc1     string    `json:"com_desc1"`
		ItAssetID    uint      `json:"it_asset_id"`
		ComName      string    `json:"com_name"`
		Serialnumber string    `json:"serialnumber"`
		FirstNameTH  string    `json:"first_name_th"`
		LastNameTH   string    `json:"last_name_th"`
		Timestamp    time.Time `json:"timestamp"`
	}

	var history []HistoryResponse

	err := database.DB.
		Table("abcinv.abc_history AS h").
		Select(`
            h.his_id,
            h.agent_id,
            h.asset_status,
            h.com_desc1,
            h.it_asset_id,
            COALESCE(asset.com_name, asset_rent.com_name) AS com_name,
            COALESCE(asset.serialnumber, asset_rent.serialnumber) AS serialnumber,
            agent.first_name_th,
            agent.last_name_th,
			h.timestamp
        `).
		Joins(`
            LEFT JOIN qamon.iam_agents AS agent
                ON agent.agent_id = h.agent_id
        `).
		Joins(`
            LEFT JOIN abcinv.abc_asset AS asset
                ON asset.it_asset_id = h.it_asset_id
        `).
		Joins(`
            LEFT JOIN abcinv.abc_asset_rent AS asset_rent
                ON asset_rent.it_asset_id = h.it_asset_id
        `).
		Where("h.it_asset_id = ?", itAssetID).
		Find(&history).Error

	if err != nil {
		c.JSON(500, gin.H{
			"error":   "ไม่สามารถดึงประวัติ Asset ได้",
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"history": history})
}

func GetProjects(c *gin.Context) {

	var Projects []models.IamProject
	if err := database.DB.Where("is_active = ?", 1).Order("proj_name asc").Find(&Projects).Error; err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูล Projects ได้", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"projects": Projects})
}
