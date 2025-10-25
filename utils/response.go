package utils

import (
	"github.com/gin-gonic/gin"
)

func JSONResponse(c *gin.Context, status int, msg string, data any) {
	c.JSON(status, gin.H{
		"message": msg,
		"data":    data,
	})
}

func JSONError(c *gin.Context, status int, code string, msg string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": msg,
		},
	})
}
