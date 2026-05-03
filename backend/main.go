package main

import (
	"fmt"
	"smart_transportation/config"
	"smart_transportation/database"
	"smart_transportation/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	database.InitDB()

	if config.AppConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	routes.SetupRoutes(router)

	port := config.AppConfig.Server.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server starting on port %s...\n", port)
	fmt.Printf("📊 API endpoints available:\n")
	fmt.Printf("   - Health check: http://localhost:%s/api/ping\n", port)
	fmt.Printf("   - User API: http://localhost:%s/api/user/*\n", port)
	fmt.Printf("   - Admin API: http://localhost:%s/api/admin/*\n", port)
	fmt.Printf("   - Public API: http://localhost:%s/api/public/*\n", port)
	fmt.Printf("\n📋 Default admin credentials:\n")
	fmt.Printf("   - Username: admin\n")
	fmt.Printf("   - Password: admin123\n")

	if err := router.Run(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
