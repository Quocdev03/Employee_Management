package main

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/database"
	"ChiQuoc/HocGolang/middleware"
	"ChiQuoc/HocGolang/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	config.LoadEnv()
	config.ConnectDB()

	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	// CORS: cho phép frontend gọi API
	middleware.SetupCORS(router)

	routes.Setup(router)

	// Chỉ chạy khi cần
	if os.Getenv("SEEDDATA") == "true" {
		// đồng bộ schema DB với struct
		database.MigrateModels()
		// khởi tạo dữ liệu mẫu
		database.Seed()
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
