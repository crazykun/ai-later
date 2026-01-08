// main.go
package main

import (
	"ai-navigator/config"
	"ai-navigator/handlers"
	"ai-navigator/middleware"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create a new Gin router with default middleware
	r := gin.Default()

	// Setup session
	store := cookie.NewStore([]byte(config.AppConfig.Session.Secret))
	r.Use(sessions.Sessions("admin_session", store))

	// Load HTML templates - explicitly list files to avoid directory issues
	r.LoadHTMLFiles(
		"templates/index.html",
		"templates/error.html",
		"templates/layout.html",
		"templates/admin/admin-login.html",
		"templates/admin/admin-index.html",
		"templates/admin/admin-sites.html",
		"templates/admin/admin-add-site.html",
		"templates/admin/admin-edit-site.html",
	)

	// Serve static files
	r.Static("/static", "./static")

	// 在router中使用中间件
	r.Use(middleware.AddGlobalContext())

	// Frontend routes
	r.GET("/", handlers.HomeHandler)
	r.GET("/search", handlers.SearchHandler)

	// Admin routes
	admin := r.Group("/admin")
	{
		admin.GET("/login", handlers.AdminLoginHandler)
		admin.POST("/login", handlers.AdminLoginPostHandler)
		admin.GET("/logout", handlers.AdminLogoutHandler)
		admin.GET("/captcha", handlers.CaptchaHandler)

		// Protected admin routes
		adminAuth := admin.Group("/")
		adminAuth.Use(middleware.AdminAuthMiddleware())
		{
			adminAuth.GET("/", handlers.AdminIndexHandler)
			adminAuth.GET("/sites", handlers.AdminSitesHandler)
			adminAuth.GET("/sites/add", handlers.AdminAddSiteHandler)
			adminAuth.POST("/sites/add", handlers.AdminAddSitePostHandler)
			adminAuth.GET("/sites/edit/:id", handlers.AdminEditSiteHandler)
			adminAuth.POST("/sites/edit/:id", handlers.AdminEditSitePostHandler)
			adminAuth.GET("/sites/delete/:id", handlers.AdminDeleteSiteHandler)
		}
	}

	// Start server on configured port
	port := config.AppConfig.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on http://localhost:%s", port)
	r.Run(":" + port)
}
