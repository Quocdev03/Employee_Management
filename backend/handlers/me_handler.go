package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MeHandler trả về thông tin user hiện tại dựa trên token
func MeHandler(ctx *gin.Context) {
	userID, _ := ctx.Get("user")

	var user models.User
	if err := config.DB.Preload("Role").Preload("Employee.Department").First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Tài khoản không còn tồn tại!",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"user": gin.H{
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
		},
	})
}