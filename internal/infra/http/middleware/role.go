package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}

		if role.(string) != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
