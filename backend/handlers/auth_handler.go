package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/dto"
	"ChiQuoc/HocGolang/middleware"
	"ChiQuoc/HocGolang/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Gọi tới models user
	var user models.User
	// load thêm bảng role Preload - Where chạy sau First()
	/*
		SELECT * FROM roles WHERE id = user.role_id;
		SELECT * FROM users WHERE email = '...' LIMIT 1;

		Preload → config load relation
		Where   → thêm điều kiện
		First   → chạy query chính
		Preload → chạy query phụ (roles)
		Error   → trả lỗi
	*/
	if err := // Preload cả Role lẫn Employee (bao gồm cả Department)
		config.DB.Preload("Role").Preload("Employee.Department").Where("email = ?", input.Email).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Tài khoản không tồn tại!",
		})
		return
	}

	// Kiểm tra người dùng active
	if !user.IsActive {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Tài khoản bị vô hiệu hoá!",
		})
		return
	}

	//  Kiểm tra pass
	// So sánh pass input với pass đã được hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Sai mật khẩu!",
		})
		return

	}

	// Đăng ký token
	claims := &middleware.Claims{
		UserID: user.ID,
		Role:   user.Role.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
		},
	}

	// Tạo token jwt
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	// Ký bằng secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": "Lỗi server!",
		})
		return
	}

	// Trả về res
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   tokenString,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name": func() string {
				if user.Employee != nil {
					return user.Employee.Name
				}
				return ""
			}(),
			"avatar_url": func() string {
				if user.Employee != nil {
					return user.Employee.AvatarURL
				}
				return ""
			}(),
			"role":     user.Role.Name,
		},
	})
}


