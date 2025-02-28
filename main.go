// main.go
package main

import (
    "github.com/gin-gonic/gin"
    "log"
    "ai-navigator/handlers"
)

func main() {
    // Create a new Gin router with default middleware
    r := gin.Default()

    // Load HTML templates
    r.LoadHTMLGlob("templates/*")

    // Serve static files
    r.Static("/static", "./static")

    // Routes
    r.GET("/", handlers.HomeHandler)
    r.GET("/search", handlers.SearchHandler)
    
    // Start server on port 8080
    log.Println("Server starting on http://localhost:8080")
    r.Run(":8080")
}