// main.go
package main

import (
	"ai-navigator/config"
	"ai-navigator/handlers"
	"ai-navigator/middleware"
	"log"

	"ai-navigator/global"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	initConfig()

	// Create a new Gin router with default middleware
	r := gin.Default()

	// Setup session
	store := cookie.NewStore([]byte("secret_key"))
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

	// Get port from config
	port := getPort()

	// Start server on configured port
	log.Println("Server starting on http://localhost:" + port)
	r.Run(":" + port)
}

func initConfig() {
	// Load configuration from file
	viper.SetConfigName("config")
	viper.AddConfigPath(".") // Look for config in the current directory

	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在, 使用默认配置
		log.Println("No configuration file found, using default values")
		global.ConfigData = &config.Config{
			Port: "8080",
		}
		return
	}

	if err := viper.Unmarshal(&global.ConfigData); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}

func getPort() string {
	// Try to get port from the new nested structure first
	if global.ConfigData.Server.Port != "" {
		return global.ConfigData.Server.Port
	}
	// Fall back to the old flat structure
	if global.ConfigData.Port != "" {
		return global.ConfigData.Port
	}
	// Default port
	return "8080"
}
