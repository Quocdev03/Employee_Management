package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Load .env
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Không load được .env (có thể đang chạy Docker)")
	}
}

func ConnectDB() {
	// Kết nối đến database
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Dùng GORM mở kết nối MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Hiển thị log SQL lỗi và warning thôi để đỡ rối mắt
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic("Không thể kết nối MySQL: " + err.Error())
	}

	// // Lấy kết nối thật từ GORM
	// sqlDB, err := db.DB()
	// if err != nil {
	// 	panic(err)
	// }

	// // Tối đa 10 kết nối cùng 1 lúc
	// sqlDB.SetMaxOpenConns(10)
	// // Giữ sẳn 5 connections rảnh
	// sqlDB.SetMaxIdleConns(5)
	// // Tự động đóng kết nối sau 5 phút không dùng
	// sqlDB.SetConnMaxLifetime(time.Minute * 5)

	DB = db
	fmt.Println("Kết nối MySQL thành công!")
}
