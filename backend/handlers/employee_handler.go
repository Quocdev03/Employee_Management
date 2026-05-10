package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/dto"
	"ChiQuoc/HocGolang/models"
	"ChiQuoc/HocGolang/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func parseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		t, err = time.Parse(time.RFC3339, dateStr)
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func toPublicResponse(e models.Employee) dto.EmployeePublicResponse {
	return dto.EmployeePublicResponse{
		ID:          e.ID,
		Name:        e.Name,
		Gender:      e.Gender,
		DateOfBirth: e.DateOfBirth,
		Phone:       e.Phone,
		Position:    e.Position,
		AvatarURL:   e.AvatarURL,
		Status:      e.Status,
		Department:  e.Department,
		HireDate:    e.HireDate,
	}
}

func GetEmployeeList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")
	deptID := ctx.Query("department_id")

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := config.DB.Model(&models.Employee{}).
		Preload("Department").
		Preload("Position").
		Preload("Position.Department").
		Preload("User")

	if search != "" {
		words := strings.Fields(strings.TrimSpace(search))
		for _, w := range words {
			like := "%" + w + "%"
			query = query.Where("employees.name LIKE ? OR employees.phone LIKE ?", like, like)
		}
	}

	if deptID != "" {
		query = query.Where("department_id = ?", deptID)
	}

	var total int64
	query.Count(&total)

	var employees []models.Employee
	if err := query.Offset(offset).Limit(limit).Order("employees.created_at DESC").Find(&employees).Error; err != nil {
		utils.InternalError(ctx, "Không thể lấy danh sách nhân viên")
		return
	}

	role, _ := ctx.Get("role")
	var responseData interface{}

	if role == "admin" {
		responseData = employees
	} else {
		public := make([]dto.EmployeePublicResponse, len(employees))
		for i, e := range employees {
			public[i] = toPublicResponse(e)
		}
		responseData = public
	}

	utils.SuccessWithMeta(ctx, responseData, gin.H{
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func GetEmployeeID(ctx *gin.Context) {
	var employee models.Employee
	if err := config.DB.Preload("Department").Preload("Position").Preload("Position.Department").Preload("User").First(&employee, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy nhân viên!")
		return
	}

	role, _ := ctx.Get("role")
	if role == "admin" {
		utils.Success(ctx, employee)
		return
	}

	utils.Success(ctx, toPublicResponse(employee))
}

func CreateEmployee(ctx *gin.Context) {
	var input dto.CreateEmployeeInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	phone := strings.TrimSpace(input.Phone)
	if phone != "" {
		var exist models.Employee
		if err := config.DB.Where("phone = ?", phone).First(&exist).Error; err == nil {
			utils.Conflict(ctx, "Số điện thoại đã tồn tại")
			return
		}
	}

	if input.DepartmentID != nil {
		var dept models.Department
		if err := config.DB.First(&dept, *input.DepartmentID).Error; err != nil {
			utils.BadRequest(ctx, "Phòng ban không tồn tại")
			return
		}
	}

	// Kiểm tra chức vụ có thuộc đúng phòng ban không
	if input.PositionID != nil {
		var pos models.Position
		if err := config.DB.First(&pos, *input.PositionID).Error; err != nil {
			utils.BadRequest(ctx, "Chức vụ không tồn tại")
			return
		}
		if input.DepartmentID != nil && pos.DepartmentID != *input.DepartmentID {
			utils.BadRequest(ctx, "Chức vụ không thuộc phòng ban đã chọn")
			return
		}
	}

	hireDate, err := parseDate(input.HireDate)
	if err != nil {
		utils.BadRequest(ctx, "Định dạng ngày vào làm không hợp lệ")
		return
	}

	dob, err := parseDate(input.DateOfBirth)
	if err != nil {
		utils.BadRequest(ctx, "Định dạng ngày sinh không hợp lệ")
		return
	}

	employee := models.Employee{
		Name:         input.Name,
		Gender:       input.Gender,
		Phone:        phone,
		DepartmentID: input.DepartmentID,
		PositionID:   input.PositionID,
		Salary:       input.Salary,
		HireDate:     hireDate,
		DateOfBirth:  dob,
		Status:       models.StatusActive,
		AvatarURL:    input.AvatarURL,
	}

	if err := config.DB.Create(&employee).Error; err != nil {
		utils.InternalError(ctx, "Tạo nhân viên thất bại")
		return
	}

	config.DB.Preload("Department").Preload("Position").Preload("Position.Department").First(&employee, employee.ID)
	utils.Create(ctx, employee)
}

func UpdateEmployee(ctx *gin.Context) {
	var employee models.Employee
	if err := config.DB.First(&employee, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy nhân viên!")
		return
	}

	var input dto.UpdateEmployeeInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if input.Name != nil { updates["name"] = *input.Name }
	if input.Gender != nil { updates["gender"] = *input.Gender }
	if input.Phone != nil {
		phone := strings.TrimSpace(*input.Phone)
		if phone != "" {
			var exist models.Employee
			if err := config.DB.Where("phone = ? AND id != ?", phone, employee.ID).First(&exist).Error; err == nil {
				utils.Conflict(ctx, "Số điện thoại đã tồn tại")
				return
			}
		}
		updates["phone"] = phone
	}
	if input.DepartmentID != nil { updates["department_id"] = input.DepartmentID }
	if input.PositionID != nil {
		// Kiểm tra chức vụ có thuộc đúng phòng ban không
		var pos models.Position
		if err := config.DB.First(&pos, *input.PositionID).Error; err != nil {
			utils.BadRequest(ctx, "Chức vụ không tồn tại")
			return
		}
		deptID := employee.DepartmentID
		if input.DepartmentID != nil {
			deptID = input.DepartmentID
		}
		if deptID != nil && pos.DepartmentID != *deptID {
			utils.BadRequest(ctx, "Chức vụ không thuộc phòng ban đã chọn")
			return
		}
		updates["position_id"] = *input.PositionID
	}
	if input.Salary != nil { updates["salary"] = *input.Salary }
	if input.HireDate != nil {
		hireDate, err := parseDate(*input.HireDate)
		if err != nil {
			utils.BadRequest(ctx, "Định dạng ngày vào làm không hợp lệ")
			return
		}
		updates["hire_date"] = hireDate
	}
	if input.DateOfBirth != nil {
		dob, err := parseDate(*input.DateOfBirth)
		if err != nil {
			utils.BadRequest(ctx, "Định dạng ngày sinh không hợp lệ")
			return
		}
		updates["date_of_birth"] = dob
	}
	if input.Status != nil { updates["status"] = *input.Status }
	if input.AvatarURL != nil { updates["avatar_url"] = *input.AvatarURL }

	if err := config.DB.Model(&employee).Updates(updates).Error; err != nil {
		utils.InternalError(ctx, "Cập nhật thất bại")
		return
	}

	config.DB.Preload("Department").Preload("Position").Preload("Position.Department").First(&employee, employee.ID)
	utils.Success(ctx, employee)
}

func DeleteEmployee(ctx *gin.Context) {
	var employee models.Employee
	if err := config.DB.First(&employee, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy nhân viên!")
		return
	}

	var associatedUser models.User
	if err := config.DB.Where("employee_id = ?", employee.ID).First(&associatedUser).Error; err == nil {
		currentUserID, _ := ctx.Get("user")
		if currentUserID.(uint) == associatedUser.ID {
			utils.Forbidden(ctx, "Không thể xoá nhân viên liên kết với tài khoản của chính mình")
			return
		}
		if associatedUser.RoleID == 1 {
			utils.Forbidden(ctx, "Không thể xoá nhân viên là Quản trị viên")
			return
		}
		config.DB.Delete(&associatedUser)
	}

	config.DB.Model(&employee).Update("status", models.StatusInactive)
	if err := config.DB.Delete(&employee).Error; err != nil {
		utils.InternalError(ctx, "Xoá nhân viên thất bại")
		return
	}

	utils.Success(ctx, "Xoá nhân viên thành công")
}
