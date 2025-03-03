package middleware

import (
	"ai-navigator/global"

	"github.com/gin-gonic/gin"
)

func AddGlobalContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Copyright", global.ConfigData.Copyright)
		c.Next()
	}
}
