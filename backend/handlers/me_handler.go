package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/dto"
	"ChiQuoc/HocGolang/models"
	"ChiQuoc/HocGolang/utils"

	"github.com/gin-gonic/gin"
)

// MeHandler trả về thông tin user hiện tại dựa trên token
func MeHandler(ctx *gin.Context) {
	userID, _ := ctx.Get("user")

	var user models.User
	if err := config.DB.Preload("Role").
		Preload("Employee.Department").
		Preload("Employee.Position").
		Preload("Employee.Position.Department").
		First(&user, userID).Error; err != nil {
		utils.Unauthorized(ctx, "Tài khoản không còn tồn tại!")
		return
	}

	name := ""
	avatar := ""
	if user.Employee != nil {
		name = user.Employee.Name
		avatar = user.Employee.AvatarURL
	}

	utils.Success(ctx, dto.MeResponse{
		User: dto.AuthUserDTO{
			ID:     user.ID,
			Email:  user.Email,
			Name:   name,
			Avatar: avatar,
			Role:   user.Role.Name,
		},
		Employee: user.Employee,
	})
}
