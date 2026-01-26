package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mock auth until we have a real auth system
		const mockUserID = "a1b2c3d4-e5f6-7890-1234-567890abcdef"
		c.Set("userID", mockUserID)
		c.Next()
	}
}
