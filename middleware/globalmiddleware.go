package middleware

import (
	"ai-navigator/global"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AddGlobalContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Copyright", global.ConfigData.Copyright)
		c.Next()
	}
}

// AdminAuthMiddleware checks if user is authenticated for admin access
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("admin_logged_in") == true {
			c.Next()
			return
		}
		c.Redirect(http.StatusFound, "/admin/login")
	}
}
