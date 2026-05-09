package middleware

import (
	"net/http"
	"os"
	"strings"

	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddlewareJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Thiếu token xác thực!",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token không hợp lệ hoặc đã hết hạn!",
			})
			return
		}

		// Kiểm tra user còn tồn tại trong DB không (phòng trường hợp reset DB hoặc bị xóa)
		var user models.User
		if err := config.DB.First(&user, claims.UserID).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Tài khoản không còn tồn tại!",
			})
			return
		}

		// Lưu thông tin user vào context
		ctx.Set("user", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
