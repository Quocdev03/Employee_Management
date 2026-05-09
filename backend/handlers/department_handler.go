package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDepartments(ctx *gin.Context) {
	var departments []models.Department
	config.DB.Find(&departments)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    departments,
	})
}
