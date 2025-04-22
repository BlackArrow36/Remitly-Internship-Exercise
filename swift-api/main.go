package main

import (
	"swift-api/config"
	"swift-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	config.InitDB()
	// Set up the Gin router
	r := gin.Default()
	// Setup routes
	routes.SetupRoutes(r)
	// Run the server
	r.Run(":8080")
}
