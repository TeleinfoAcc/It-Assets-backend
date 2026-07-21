package handlers

import (
	"crypto/tls"
	"log"
	"net/smtp"
)

func SendEmailNotification() {
	username := "dbaauto@teleinfomedia.co.th"
	password := "Tkyp@Auto#1188"
	smtpHost := "mailawn.thaicloudsolutions.com"
	smtpPort := "25"

	to := []string{"jirawas.tul@gmail.com"}
	fromHeader := "From: Equipment Readiness System <" + username + ">\n" // 👈 ชื่อผู้ส่งที่แสดงในอีเมล แก้ที่ตัวแปร from ด้านบน
	subject := "Subject: ⚠️ แจ้งเตือน: ตรวจสอบสรุปรายงานประจำเดือน\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<h3>เรียนผู้จัดการ</h3><p>ระบบตรวจสอบพบว่าถึงกำหนดเวลา วันที่ 10 แล้ว กรุณาเข้าตรวจสอบระบบครับ</p>"
	message := []byte(fromHeader + subject + mime + body)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // 👈 สั่งให้ Go หลับตาข้างนึง ยอมรับใบรับรองใบนี้
		ServerName:         smtpHost,
	}

	// 2. เชื่อมต่อไปยังเซิร์ฟเวอร์ SMTP
	conn, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		log.Println("❌ ไม่สามารถเชื่อมต่อ SMTP Server ได้:", err)
		return
	}
	defer conn.Close()

	// 2. ทำ StartTLS ข้ามตรวจใบรับรองเหมือนเดิม
	if err = conn.StartTLS(tlsConfig); err != nil {
		log.Println("❌ ทำ StartTLS ล้มเหลว:", err)
		return
	}

	// 3. ยืนยันตัวตนด้วย username/password
	auth := smtp.PlainAuth("", username, password, smtpHost)
	if err = conn.Auth(auth); err != nil {
		log.Println("❌ Authentication ล้มเหลว:", err)
		return
	}

	// 4. ระบุผู้ส่ง
	// 4. ระบุผู้ส่ง
	if err = conn.Mail(username); err != nil {
		log.Println("❌ ระบุผู้ส่งล้มเหลว:", err)
		return
	}

	// 5. ระบุผู้รับ
	for _, k := range to {
		if err = conn.Rcpt(k); err != nil {
			log.Println("❌ ระบุผู้รับล้มเหลว:", err)
			return
		}
	}

	// 6. ส่งเนื้อหาอีเมลตามปกติ
	w, err := conn.Data()
	if err != nil {
		log.Println("❌ เตรียมส่งข้อมูลล้มเหลว:", err)
		return
	}
	_, err = w.Write(message)
	if err != nil {
		log.Println("❌ เขียนข้อมูลเมลล้มเหลว:", err)
		return
	}
	err = w.Close()
	if err != nil {
		log.Println("❌ ปิดการส่งข้อมูลล้มเหลว:", err)
		return
	}

	err = conn.Quit() // สั่งปิดการเชื่อมต่อแบบสุภาพ (Quit)
	if err != nil {
		log.Println("⚠️ เซิร์ฟเวอร์ปิดการเชื่อมต่อแปลกๆ แต่ส่งข้อมูลไปแล้ว:", err)
		return
	}

	log.Println("🚀 [Success] ส่งอีเมลผ่านระบบบริษัทสำเร็จแล้วโดยไม่ต้องใช้รหัสผ่าน!")
}
