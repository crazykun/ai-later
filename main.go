// main.go
package main

import (
	"ai-navigator/config"
	"ai-navigator/handlers"
	"ai-navigator/middleware"
	"log"

	"ai-navigator/global"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	initConfig()

	// Create a new Gin router with default middleware
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Serve static files
	r.Static("/static", "./static")

	// 在router中使用中间件
	r.Use(middleware.AddGlobalContext())

	// Routes
	r.GET("/", handlers.HomeHandler)
	r.GET("/search", handlers.SearchHandler)

	// Start server on port 8080
	log.Println("Server starting on http://localhost:" + global.ConfigData.Port)
	r.Run(":" + global.ConfigData.Port)
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
	}

	if err := viper.Unmarshal(&global.ConfigData); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
