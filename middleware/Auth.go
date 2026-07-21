package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret-key")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// รับ Authorization header จาก Angular's authInterceptor
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "ไม่มี token"})
			c.Abort() // หยุดไม่ให้ไปต่อ
			return
		}

		// ตัด "Bearer " ออก
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// ตรวจ token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "token ไม่ถูกต้อง"})
			c.Abort()
			return
		}

		// ดึง agent_id ใส่ไว้ใน context
		claims := token.Claims.(jwt.MapClaims)
		c.Set("agent_id", claims["agent_id"])
		c.Set("permission_id", claims["permission_id"])

		c.Next() // ไปต่อ
	}
}
