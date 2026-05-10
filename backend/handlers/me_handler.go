package handlers

import (
	"ChiQuoc/HocGolang/config"
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

	utils.Success(ctx, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name": func() string {
			if user.Employee != nil {
				return user.Employee.Name
			}
			return ""
		}(),
		"avatar": func() string {
			if user.Employee != nil {
				return user.Employee.AvatarURL
			}
			return ""
		}(),
		"role":     user.Role.Name,
		"employee": user.Employee,
	})
}
