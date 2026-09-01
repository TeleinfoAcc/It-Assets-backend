package handlers

import (
	"checklist-backend/database"
	"checklist-backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ---------- helpers ----------

// parseDateRange อ่าน query param start_date / end_date (รูปแบบ YYYY-MM-DD)
// ถ้าไม่ส่งมา จะ default ย้อนหลัง defaultDays วันจนถึงตอนนี้
// end จะถูกบวก 1 วันเพื่อให้ครอบคลุมทั้งวันสุดท้าย (ใช้เป็น timestamp < end)
func parseDateRange(c *gin.Context, defaultDays int) (start time.Time, end time.Time) {
	now := time.Now()
	end = now
	start = now.AddDate(0, 0, -defaultDays)

	if s := c.Query("start_date"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			start = t
		}
	}
	if e := c.Query("end_date"); e != "" {
		if t, err := time.ParseInLocation("2006-01-02", e, time.Local); err == nil {
			end = t.AddDate(0, 0, 1)
		}
	}
	return start, end
}

// =====================================================================
//  1. รายงานสรุปครุภัณฑ์ (Inventory Summary)
//     GET /api/report/inventory
//     query: com_type, curr_room_code, asset_project, include_retired=true, detail=true
//     ใช้ตอนตรวจนับทรัพย์สิน / ส่งฝ่ายบัญชี (frontend เอาไป export excel ได้)
//
// =====================================================================
func GetReportInventory(c *gin.Context) {
	where := "1=1"
	args := []interface{}{}

	if v := c.Query("com_type"); v != "" {
		where += " AND r.com_type = ?"
		args = append(args, v)
	}
	if v := c.Query("curr_room_code"); v != "" {
		where += " AND r.curr_room_code = ?"
		args = append(args, v)
	}
	if v := c.Query("asset_project"); v != "" {
		where += " AND r.asset_project = ?"
		args = append(args, v)
	}
	if c.Query("include_retired") != "true" {
		where += " AND r.asset_status::text <> '5'"
	}

	type kv struct {
		Key   string `json:"key"`
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}

	run := func(sql string) []kv {
		var rows []kv
		database.DB.Raw(sql, args...).Scan(&rows)
		return rows
	}

	var total int64
	database.DB.Raw(`SELECT COUNT(*) FROM abcinv.abc_asset_rent r WHERE `+where, args...).Scan(&total)

	byStatus := run(`
		SELECT r.asset_status::text AS key,
		       COALESCE(s.asset_status_name, '-') AS name,
		       COUNT(*) AS count
		FROM abcinv.abc_asset_rent r
		LEFT JOIN abcinv.abc_asset_status s ON s.asset_status::text = r.asset_status::text
		WHERE ` + where + `
		GROUP BY r.asset_status::text, s.asset_status_name
		ORDER BY count DESC`)

	byType := run(`
		SELECT COALESCE(NULLIF(r.com_type, ''), '-') AS key, '' AS name, COUNT(*) AS count
		FROM abcinv.abc_asset_rent r WHERE ` + where + `
		GROUP BY r.com_type ORDER BY count DESC`)

	byBrand := run(`
		SELECT COALESCE(NULLIF(r.com_brand, ''), '-') AS key, '' AS name, COUNT(*) AS count
		FROM abcinv.abc_asset_rent r WHERE ` + where + `
		GROUP BY r.com_brand ORDER BY count DESC`)

	byProject := run(`
		SELECT COALESCE(NULLIF(r.asset_project, ''), '-') AS key, '' AS name, COUNT(*) AS count
		FROM abcinv.abc_asset_rent r WHERE ` + where + `
		GROUP BY r.asset_project ORDER BY count DESC`)

	byRoom := run(`
		SELECT COALESCE(NULLIF(r.curr_room_code, ''), '-') AS key, '' AS name, COUNT(*) AS count
		FROM abcinv.abc_asset_rent r WHERE ` + where + `
		GROUP BY r.curr_room_code ORDER BY count DESC`)

	resp := gin.H{
		"total":      total,
		"by_status":  byStatus,
		"by_type":    byType,
		"by_brand":   byBrand,
		"by_project": byProject,
		"by_room":    byRoom,
	}

	if c.Query("detail") == "true" {
		var items []models.AbcAssetRent
		q := database.DB.Preload("Asset_status_name")
		if v := c.Query("com_type"); v != "" {
			q = q.Where("com_type = ?", v)
		}
		if v := c.Query("curr_room_code"); v != "" {
			q = q.Where("curr_room_code = ?", v)
		}
		if v := c.Query("asset_project"); v != "" {
			q = q.Where("asset_project = ?", v)
		}
		if c.Query("include_retired") != "true" {
			q = q.Where("asset_status::text <> ?", "5")
		}
		q.Order("com_name asc").Find(&items)
		resp["items"] = items
	}

	c.JSON(200, resp)
}

// =====================================================================
//  2. รายงานการเคลื่อนไหวสถานะตามช่วงเวลา (Status Movement)
//     GET /api/report/status-movement?start_date=2026-08-01&end_date=2026-08-31
//     สรุปจาก abc_history ว่าในช่วงเวลานั้นมีการเปลี่ยนสถานะไปเป็นอะไรบ้าง กี่ครั้ง
//     ใครเป็นคนทำ และแนวโน้มรายวัน
//
// =====================================================================
func GetReportStatusMovement(c *gin.Context) {
	start, end := parseDateRange(c, 30)

	type statusRow struct {
		AssetStatus     uint   `json:"asset_status"`
		AssetStatusName string `json:"asset_status_name"`
		Count           int64  `json:"count"`
	}
	var byStatus []statusRow
	database.DB.Raw(`
		SELECT h.asset_status,
		       COALESCE(s.asset_status_name, '-') AS asset_status_name,
		       COUNT(*) AS count
		FROM abcinv.abc_history h
		LEFT JOIN abcinv.abc_asset_status s ON s.asset_status = h.asset_status
		WHERE h.timestamp >= ? AND h.timestamp < ?
		GROUP BY h.asset_status, s.asset_status_name
		ORDER BY count DESC`, start, end).Scan(&byStatus)

	type dayRow struct {
		Day   time.Time `json:"day"`
		Count int64     `json:"count"`
	}
	var byDay []dayRow
	database.DB.Raw(`
		SELECT date_trunc('day', h.timestamp) AS day, COUNT(*) AS count
		FROM abcinv.abc_history h
		WHERE h.timestamp >= ? AND h.timestamp < ?
		GROUP BY 1 ORDER BY 1`, start, end).Scan(&byDay)

	type agentRow struct {
		AgentID     uint   `json:"agent_id"`
		FirstNameTH string `json:"first_name_th"`
		LastNameTH  string `json:"last_name_th"`
		Count       int64  `json:"count"`
	}
	var byAgent []agentRow
	database.DB.Raw(`
		SELECT h.agent_id, a.first_name_th, a.last_name_th, COUNT(*) AS count
		FROM abcinv.abc_history h
		LEFT JOIN qamon.iam_agents a ON a.agent_id = h.agent_id
		WHERE h.timestamp >= ? AND h.timestamp < ?
		GROUP BY h.agent_id, a.first_name_th, a.last_name_th
		ORDER BY count DESC`, start, end).Scan(&byAgent)

	var total int64
	database.DB.Raw(`SELECT COUNT(*) FROM abcinv.abc_history h
		WHERE h.timestamp >= ? AND h.timestamp < ?`, start, end).Scan(&total)

	c.JSON(200, gin.H{
		"start_date":   start.Format("2006-01-02"),
		"end_date":     end.AddDate(0, 0, -1).Format("2006-01-02"),
		"total_events": total,
		"by_status":    byStatus,
		"by_day":       byDay,
		"by_agent":     byAgent,
	})
}

// =====================================================================
//  3. รายงานเครื่องเสีย / Downtime (Repair report)
//     GET /api/report/repair?start_date=..&end_date=..&min_times=2
//     - นับจำนวนครั้งที่เครื่องถูกเปลี่ยนสถานะเป็น "เสีย" (asset_status 3 หรือ 9)
//     - MTTR = เวลาเฉลี่ยที่อยู่ในสถานะเสียก่อนถูกเปลี่ยนสถานะครั้งถัดไป
//     - list เครื่องที่เสียตั้งแต่ min_times ครั้งขึ้นไป (เครื่องที่ควรพิจารณาเปลี่ยน)
//
// =====================================================================
func GetReportRepair(c *gin.Context) {
	start, end := parseDateRange(c, 90)

	minTimes := 2
	if v := c.Query("min_times"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			minTimes = n
		}
	}

	// CTE ร่วม: หา timestamp ของ event ถัดไปของแต่ละเครื่อง แล้วคัดเฉพาะ event ที่เป็น "เสีย"
	const episodeCTE = `
		WITH ordered AS (
			SELECT h.it_asset_id, h.asset_status, h.timestamp,
			       LEAD(h.timestamp) OVER (PARTITION BY h.it_asset_id ORDER BY h.timestamp) AS next_ts
			FROM abcinv.abc_history h
			WHERE h.timestamp >= ? AND h.timestamp < ?
		),
		episodes AS (
			SELECT it_asset_id,
			       timestamp AS broke_at,
			       CASE WHEN next_ts IS NOT NULL
			            THEN EXTRACT(EPOCH FROM (next_ts - timestamp)) / 3600.0
			       END AS repair_hours
			FROM ordered
			WHERE asset_status::text IN ('3', '9')
		)`

	// สรุปภาพรวม
	var summary struct {
		TotalRepairEvents int64    `json:"total_repair_events"`
		AvgMttrHours      *float64 `json:"avg_mttr_hours"`
		AssetsAffected    int64    `json:"assets_affected"`
	}
	database.DB.Raw(episodeCTE+`
		SELECT COUNT(*)                                   AS total_repair_events,
		       ROUND(AVG(repair_hours)::numeric, 2)       AS avg_mttr_hours,
		       COUNT(DISTINCT it_asset_id)                AS assets_affected
		FROM episodes`, start, end).Scan(&summary)

	// รายเครื่องที่เสียบ่อย
	type repairRow struct {
		ItAssetID      uint       `json:"it_asset_id"`
		ComName        string     `json:"com_name"`
		Serialnumber   string     `json:"serialnumber"`
		ComType        string     `json:"com_type"`
		RepairCount    int64      `json:"repair_count"`
		AvgRepairHours *float64   `json:"avg_repair_hours"`
		LastRepairAt   *time.Time `json:"last_repair_at"`
	}
	var rows []repairRow
	database.DB.Raw(episodeCTE+`,
		agg AS (
			SELECT it_asset_id,
			       COUNT(*)                             AS repair_count,
			       ROUND(AVG(repair_hours)::numeric, 2) AS avg_repair_hours,
			       MAX(broke_at)                        AS last_repair_at
			FROM episodes
			GROUP BY it_asset_id
			HAVING COUNT(*) >= ?
		)
		SELECT g.it_asset_id,
		       COALESCE(a.com_name, r.com_name)         AS com_name,
		       COALESCE(a.serialnumber, r.serialnumber) AS serialnumber,
		       COALESCE(a.com_type, r.com_type)         AS com_type,
		       g.repair_count,
		       g.avg_repair_hours,
		       g.last_repair_at
		FROM agg g
		LEFT JOIN abcinv.abc_asset      a ON a.it_asset_id = g.it_asset_id
		LEFT JOIN abcinv.abc_asset_rent r ON r.it_asset_id = g.it_asset_id
		ORDER BY g.repair_count DESC, g.avg_repair_hours DESC NULLS LAST`,
		start, end, minTimes).Scan(&rows)

	c.JSON(200, gin.H{
		"start_date":      start.Format("2006-01-02"),
		"end_date":        end.AddDate(0, 0, -1).Format("2006-01-02"),
		"min_times":       minTimes,
		"summary":         summary,
		"repeat_offender": rows,
	})
}
