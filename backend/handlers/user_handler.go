package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/dto"
	"ChiQuoc/HocGolang/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(ctx *gin.Context) {
	var users []models.User
	if err := config.DB.Preload("Role").Preload("Employee").Find(&users).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

func GetUserByID(ctx *gin.Context) {
	var user models.User
	if err := config.DB.Preload("Role").Preload("Employee").First(&user, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Không tìm thấy user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

func CreateUser(ctx *gin.Context) {
	var input dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Mọi user nên khuyến khích có employee liên kết (nhưng admin có thể chưa cần ngay)
	// Bỏ check bắt buộc nếu cần, hoặc giữ lại tuỳ logic.
	// User bình thường bắt buộc phải có employee.
	if input.RoleID != 1 && input.EmployeeID == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Tài khoản nhân viên phải được liên kết với hồ sơ nhân viên",
		})
		return
	}

	// Kiểm tra emp có tồn tại không
	if input.EmployeeID != nil {
		var emp models.Employee
		if err := config.DB.First(&emp, *input.EmployeeID).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Employee không tồn tại",
			})
			return
		}

		// check employee đã có user chưa
		var count int64
		config.DB.Model(&models.User{}).Where("employee_id = ?", *input.EmployeeID).Count(&count)
		if count > 0 {
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "Employee đã có tài khoản",
			})
			return
		}
	}

	// hash pass
	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	user := models.User{
		Email:        input.Email,
		PasswordHash: string(hash),
		RoleID:       input.RoleID,
		EmployeeID:   input.EmployeeID,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"success": false, "error": "Email đã tồn tại"})
		return
	}
	config.DB.Preload("Role").Preload("Employee").First(&user, user.ID)

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "data": user})

}

func UpdateUser(ctx *gin.Context) {

	// Tìm user trong db
	var user models.User
	if err := config.DB.First(&user, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Không tìm thấy user"})
		return
	}

	// Cho phép tự sửa bản thân (nhưng không được đổi Role của chính mình để tránh tự xoá quyền Admin)
	currentUser, _ := ctx.Get("user")
	isSelf := currentUser.(uint) == user.ID

	// Quản trị viên này không thể đổi thông tin quản trị viên khác
	if user.RoleID == 1 && !isSelf {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Không có quyền chỉnh sửa tài khoản Quản trị viên khác",
		})
		return
	}

	// Lấy data từ body (tất cả field đều tuỳ chọn, nhập 1 hoặc nhiều field)
	var input dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
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
		// Ngăn chặn tự đổi Role của chính mình (để bảo vệ quyền Admin)
		if isSelf && input.RoleID != user.RoleID {
			ctx.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Không thể tự thay đổi quyền hạn của chính mình"})
			return
		}
		updates["role_id"] = input.RoleID
	}

	if input.EmployeeID != nil {
		// Nếu thay đổi employee_id, kiểm tra xem nó đã được dùng bởi user khác chưa
		if *input.EmployeeID != 0 && (user.EmployeeID == nil || *input.EmployeeID != *user.EmployeeID) {
			var count int64
			config.DB.Model(&models.User{}).Where("employee_id = ? AND id != ?", *input.EmployeeID, user.ID).Count(&count)
			if count > 0 {
				ctx.JSON(http.StatusConflict, gin.H{
					"success": false,
					"error":   "Nhân viên này đã được liên kết với một tài khoản khác",
				})
				return
			}
		}
		updates["employee_id"] = *input.EmployeeID
	}

	// Phải có ít nhất 1 field để update
	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Vui lòng nhập ít nhất 1 trường cần cập nhật (email, password, role_id)"})
		return
	}

	config.DB.Model(&user).Updates(updates)
	config.DB.Preload("Role").Preload("Employee").First(&user, user.ID)

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": user})

}

func DeleteUser(ctx *gin.Context) {
	// Tìm user trong db
	var user models.User

	if err := config.DB.First(&user, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Không tìm thấy user",
			"error":   err.Error(),
		})
		return
	}

	currentUser, _ := ctx.Get("user")

	if currentUser.(uint) == user.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Không thể xoá chính mình"})
		return
	}

	// Quản trị viên này không thể xoá quản trị viên khác
	if user.RoleID == 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Không có quyền xoá tài khoản Quản trị viên khác",
		})
		return
	}

	config.DB.Delete(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Xoá user thành công",
	})
}
