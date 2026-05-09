package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/dto"
	"ChiQuoc/HocGolang/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(ctx *gin.Context) {
	var stats dto.DashboardResponse
	config.DB.Model(&models.Employee{}).Count(&stats.TotalEmployees)

	// Lấy nv trạng thái active
	config.DB.Model(&models.Employee{}).Where("status = ?", models.StatusActive).Count(&stats.ActiveEmployees)

	// Lấy nv trạng thái inactive
	config.DB.Model(&models.Employee{}).Where("status = ?", models.StatusInactive).Count(&stats.InactiveEmployees)

	// Lấy user có role admin
	config.DB.Model(&models.User{}).Where("role_id = 1").Count(&stats.TotalAdminRole)

	// Lấy tổng các phòng ban
	config.DB.Model(&models.Department{}).Count(&stats.TotalDepartments)

	// Tổng user
	config.DB.Model(&models.User{}).Count(&stats.TotalUsers)

	// Thống kê nhân viên theo phòng ban (Lấy tên phòng ban và số lượng)
	config.DB.Table("employees").
		Select("departments.name as department_name, count(*) as count").
		Joins("left join departments on departments.id = employees.department_id").
		Group("departments.name").
		Scan(&stats.EmployeesByDept)

	ctx.JSON(http.StatusOK, stats)

}
