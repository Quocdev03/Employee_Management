package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/dto"
	"ChiQuoc/HocGolang/middleware"
	"ChiQuoc/HocGolang/models"
	"ChiQuoc/HocGolang/utils"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	var user models.User
	if err := config.DB.Preload("Role").Preload("Employee").Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.Unauthorized(ctx, "Tài khoản không tồn tại!")
		return
	}

	if !user.IsActive {
		utils.Forbidden(ctx, "Tài khoản bị vô hiệu hoá!")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		utils.Unauthorized(ctx, "Sai mật khẩu!")
		return
	}

	// Token claims
	claims := &middleware.Claims{
		UserID: user.ID,
		Role:   user.Role.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	// Map to AuthUserDTO
	name := ""
	avatar := ""
	if user.Employee != nil {
		name = user.Employee.Name
		avatar = user.Employee.AvatarURL
	}

	res := dto.LoginResponse{
		Token: tokenString,
		User: dto.AuthUserDTO{
			ID:     user.ID,
			Email:  user.Email,
			Name:   name,
			Avatar: avatar,
			Role:   user.Role.Name,
		},
	}

	utils.Success(ctx, res)
}
