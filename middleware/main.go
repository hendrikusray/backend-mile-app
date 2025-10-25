package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ah := c.GetHeader("Authorization")
		if !strings.HasPrefix(ah, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthorized", "message": "missing bearer token"}})
			return
		}
		token := strings.TrimPrefix(ah, "Bearer ")
		if token != "mock-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthorized", "message": "invalid token"}})
			return
		}

		c.Set("user_id", "U001")
		c.Next()
	}
}
