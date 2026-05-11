package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/dto"
	"ChiQuoc/HocGolang/models"
	"ChiQuoc/HocGolang/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// --- Helpers ---

func toUserResponse(u models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		RoleID:     u.RoleID,
		Role:       &u.Role,
		EmployeeID: u.EmployeeID,
		Employee:   u.Employee,
		IsActive:   u.IsActive,
	}
}

// --- Handlers ---

func GetUsers(ctx *gin.Context) {
	var users []models.User
	if err := config.DB.Preload("Role").Preload("Employee").Find(&users).Error; err != nil {
		utils.InternalError(ctx, "Không thể lấy danh sách người dùng")
		return
	}

	res := make([]dto.UserResponse, len(users))
	for i := range users {
		res[i] = toUserResponse(users[i])
	}
	utils.Success(ctx, res)
}

func GetUserByID(ctx *gin.Context) {
	var user models.User
	if err := config.DB.Preload("Role").Preload("Employee").First(&user, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy người dùng")
		return
	}
	utils.Success(ctx, toUserResponse(user))
}

func CreateUser(ctx *gin.Context) {
	var input dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	// Logic kiểm tra Employee đã có account chưa
	if input.EmployeeID != nil {
		var count int64
		config.DB.Model(&models.User{}).Where("employee_id = ?", *input.EmployeeID).Count(&count)
		if count > 0 {
			utils.Conflict(ctx, "Nhân viên này đã có tài khoản")
			return
		}
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	user := models.User{
		Email:        input.Email,
		PasswordHash: string(hash),
		RoleID:       input.RoleID,
		EmployeeID:   input.EmployeeID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.Conflict(ctx, "Email đã tồn tại hoặc có lỗi khi tạo tài khoản")
		return
	}

	config.DB.Preload("Role").Preload("Employee").First(&user, user.ID)
	utils.Create(ctx, toUserResponse(user))
}

func UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := config.DB.First(&user, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy người dùng")
		return
	}

	currentUser, _ := ctx.Get("user")
	isSelf := currentUser.(uint) == user.ID

	if user.RoleID == 1 && !isSelf {
		utils.Forbidden(ctx, "Không có quyền chỉnh sửa Quản trị viên khác")
		return
	}

	var input dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if input.Email != "" {
		updates["email"] = input.Email
	}
	if input.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
		updates["password_hash"] = string(hash)
	}
	if input.RoleID != 0 {
		if isSelf && input.RoleID != user.RoleID {
			utils.Forbidden(ctx, "Không thể tự thay đổi quyền hạn của chính mình")
			return
		}
		updates["role_id"] = input.RoleID
	}
	if input.EmployeeID != nil {
		updates["employee_id"] = input.EmployeeID
	}

	if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
		utils.InternalError(ctx, "Cập nhật thất bại")
		return
	}

	config.DB.Preload("Role").Preload("Employee").First(&user, user.ID)
	utils.Success(ctx, toUserResponse(user))
}

func DeleteUser(ctx *gin.Context) {
	var user models.User
	if err := config.DB.First(&user, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy người dùng")
		return
	}

	currentUser, _ := ctx.Get("user")
	if currentUser.(uint) == user.ID {
		utils.Forbidden(ctx, "Không thể tự xoá chính mình")
		return
	}

	if user.RoleID == 1 {
		utils.Forbidden(ctx, "Không thể xoá Quản trị viên khác")
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		utils.InternalError(ctx, "Xoá thất bại")
		return
	}

	utils.Success(ctx, "Xoá thành công")
}
