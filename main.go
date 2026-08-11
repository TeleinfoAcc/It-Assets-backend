package main

import (
	"checklist-backend/cron"
	"checklist-backend/database"
	handlers "checklist-backend/handler"
	"checklist-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cron.StartEmailCron()

	database.Connect()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:4200", "http://172.21.142.211:8005", "http://10.238.99.152:8005"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	// Public routes (ไม่ต้อง login)
	r.POST("/api/auth/login", handlers.Login)

	// Protected routes (ต้อง login)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{

		protected.GET("/getAssets", handlers.GetAssets)
		protected.GET("/getAssetsRent", handlers.GetAssetsRent)
		protected.GET("/getAssetsRent/:id", handlers.GetAssetsRentById)
		protected.GET("/getAssetStatus", handlers.GetAssetStatus)
		protected.GET("/getSites", handlers.GetSites)
		protected.GET("/getRooms", handlers.GetRooms)
		protected.GET("/getAssets/:id", handlers.GetAssetById)
		protected.POST("/updateAssets", handlers.UpdateAsset)
		protected.POST("/updateAssetsRent", handlers.UpdateAssetRent)

		// protected.DELETE("/users/:id", handlers.DeleteUser)
	}

	r.Run(":8080")
}
