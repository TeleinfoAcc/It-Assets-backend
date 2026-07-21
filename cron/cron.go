package cron

import (
	handlers "checklist-backend/handler"
	"log"

	"github.com/robfig/cron/v3"
)

func StartEmailCron() {
	// jakartaLocation, err := time.LoadLocation("Asia/Bangkok")
	// if err != nil {
	// 	log.Println("ไม่สามารถโหลด Timezone ได้:", err)
	// 	return
	// }

	// c := cron.New(cron.WithLocation(jakartaLocation))

	c := cron.New()

	// สั่งรันทุกวันที่ 10 เวลา 9 โมงเช้า
	_, err := c.AddFunc("00 4 21 * *", func() {
		log.Println("⏰ [UTC Cron] ถึงเวลาทำงานแล้ว! เริ่มกระบวนการส่งเมล...")
		handlers.SendEmailNotification()
	})

	if err != nil {
		log.Println("ไม่สามารถสร้าง Cron Job ได้:", err)
		return
	}

	c.Start()
	log.Println("🟢 ระบบเปิดใช้งาน Cron Job (โหมด UTC) เรียบร้อยแล้ว...")
}
