package main

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/database"
	"ChiQuoc/HocGolang/middleware"
	"ChiQuoc/HocGolang/routes"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	
	config.LoadEnv()
	config.ConnectDB()

	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	// CORS: cho phép frontend gọi API
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
