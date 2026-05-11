package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/dto"
	"ChiQuoc/HocGolang/models"
	"ChiQuoc/HocGolang/utils"

	"github.com/gin-gonic/gin"
)

// GetDepartments: Lấy toàn bộ phòng ban kèm theo danh sách chức vụ
func GetDepartments(ctx *gin.Context) {
	var departments []models.Department
	if err := config.DB.Preload("Positions").Order("id asc").Find(&departments).Error; err != nil {
		utils.InternalError(ctx, "Không thể lấy danh sách phòng ban")
		return
	}
	utils.Success(ctx, departments)
}

// GetDepartmentByID: Lấy chi tiết một phòng ban kèm chức vụ
func GetDepartmentByID(ctx *gin.Context) {
	var dept models.Department
	if err := config.DB.Preload("Positions").First(&dept, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy phòng ban")
		return
	}
	utils.Success(ctx, dept)
}

// CreateDepartment: Tạo phòng ban mới
func CreateDepartment(ctx *gin.Context) {
	var input dto.DepartmentRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(ctx, "Tên phòng ban không được để trống")
		return
	}

	dept := models.Department{Name: input.Name}
	if err := config.DB.Create(&dept).Error; err != nil {
		utils.Conflict(ctx, "Tên phòng ban đã tồn tại")
		return
	}

	config.DB.Preload("Positions").First(&dept, dept.ID)
	utils.Create(ctx, dept)
}

// UpdateDepartment: Cập nhật tên phòng ban
func UpdateDepartment(ctx *gin.Context) {
	var dept models.Department
	if err := config.DB.First(&dept, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy phòng ban")
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(ctx, "Tên phòng ban không được để trống")
		return
	}

	if err := config.DB.Model(&dept).Update("name", input.Name).Error; err != nil {
		utils.Conflict(ctx, "Tên phòng ban đã tồn tại")
		return
	}

	config.DB.Preload("Positions").First(&dept, dept.ID)
	utils.Success(ctx, dept)
}

// DeleteDepartment: Xoá phòng ban (chỉ khi không có nhân viên)
func DeleteDepartment(ctx *gin.Context) {
	var dept models.Department
	if err := config.DB.First(&dept, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy phòng ban")
		return
	}

	var empCount int64
	config.DB.Model(&models.Employee{}).Where("department_id = ?", dept.ID).Count(&empCount)
	if empCount > 0 {
		utils.BadRequest(ctx, "Không thể xoá phòng ban đang có nhân viên")
		return
	}

	// Xoá tất cả chức vụ liên quan trước
	config.DB.Where("department_id = ?", dept.ID).Delete(&models.Position{})

	if err := config.DB.Delete(&dept).Error; err != nil {
		utils.InternalError(ctx, "Xoá phòng ban thất bại")
		return
	}

	utils.Success(ctx, "Xoá phòng ban thành công")
}
