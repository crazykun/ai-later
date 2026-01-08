package middleware

import (
	"ai-navigator/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AddGlobalContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		copyright := config.AppConfig.Copyright
		if copyright == "" {
			copyright = "AI导航 © 2024"
		}
		c.Set("Copyright", copyright)
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
