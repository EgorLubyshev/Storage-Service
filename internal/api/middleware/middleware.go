package middleware

import "github.com/gin-gonic/gin"

func (m *AuthMiddleware) JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Parse token
		// 2. Validate
		// 3. c.Set("user_id", userID)
		c.Next()
	}
}
