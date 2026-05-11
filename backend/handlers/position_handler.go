package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/dto"
	"ChiQuoc/HocGolang/models"
	"ChiQuoc/HocGolang/utils"

	"github.com/gin-gonic/gin"
)

// GetPositionsByDepartment: Lấy danh sách chức vụ theo phòng ban
func GetPositionsByDepartment(ctx *gin.Context) {
	var positions []models.Position
	if err := config.DB.Where("department_id = ?", ctx.Param("id")).
		Order("id asc").Find(&positions).Error; err != nil {
		utils.InternalError(ctx, "Không thể lấy danh sách chức vụ")
		return
	}
	utils.Success(ctx, positions)
}

// CreatePosition: Thêm chức vụ mới vào phòng ban
func CreatePosition(ctx *gin.Context) {
	var input dto.CreatePositionRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(ctx, "Tên chức vụ và phòng ban không được để trống")
		return
	}

	// Kiểm tra trùng tên trong cùng phòng ban
	var count int64
	config.DB.Model(&models.Position{}).
		Where("name = ? AND department_id = ?", input.Name, input.DepartmentID).
		Count(&count)
	if count > 0 {
		utils.Conflict(ctx, "Chức vụ này đã tồn tại trong phòng ban")
		return
	}

	pos := models.Position{Name: input.Name, DepartmentID: input.DepartmentID}
	if err := config.DB.Create(&pos).Error; err != nil {
		utils.InternalError(ctx, "Thêm chức vụ thất bại")
		return
	}
	utils.Create(ctx, pos)
}

// UpdatePosition: Cập nhật tên chức vụ
func UpdatePosition(ctx *gin.Context) {
	var pos models.Position
	if err := config.DB.First(&pos, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy chức vụ")
		return
	}

	var input dto.UpdatePositionRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(ctx, "Tên chức vụ không được để trống")
		return
	}

	if err := config.DB.Model(&pos).Update("name", input.Name).Error; err != nil {
		utils.InternalError(ctx, "Cập nhật chức vụ thất bại")
		return
	}
	utils.Success(ctx, pos)
}

// DeletePosition: Xoá chức vụ (chỉ khi không có nhân viên đang giữ)
func DeletePosition(ctx *gin.Context) {
	var pos models.Position
	if err := config.DB.First(&pos, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy chức vụ")
		return
	}

	var empCount int64
	config.DB.Model(&models.Employee{}).Where("position_id = ?", pos.ID).Count(&empCount)
	if empCount > 0 {
		utils.BadRequest(ctx, "Không thể xoá chức vụ đang có nhân viên đảm nhiệm")
		return
	}

	if err := config.DB.Delete(&pos).Error; err != nil {
		utils.InternalError(ctx, "Xoá chức vụ thất bại")
		return
	}
	utils.Success(ctx, "Xoá chức vụ thành công")
}
