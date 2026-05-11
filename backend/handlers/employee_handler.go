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

// --- Helpers ---

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

func ensureUserLoaded(e *models.Employee) {
	if e.User == nil {
		var user models.User
		if err := config.DB.Where("employee_id = ?", e.ID).First(&user).Error; err == nil {
			e.User = &user
		}
	}
}

func mapEmployeeBase(e models.Employee) dto.EmployeeBaseResponse {
	ensureUserLoaded(&e)
	email := ""
	if e.User != nil {
		email = e.User.Email
	}

	return dto.EmployeeBaseResponse{
		ID:          e.ID,
		Name:        e.Name,
		Gender:      e.Gender,
		Email:       email,
		DateOfBirth: e.DateOfBirth,
		Phone:       e.Phone,
		AvatarURL:   e.AvatarURL,
		Status:      e.Status,
		HireDate:    e.HireDate,
		Department:  e.Department,
		Position:    e.Position,
	}
}

func toEmployeeResponse(e models.Employee) dto.EmployeeResponse {
	return dto.EmployeeResponse{
		EmployeeBaseResponse: mapEmployeeBase(e),
		Salary:               e.Salary,
		User:                 e.User,
	}
}

func toEmployeePublicResponse(e models.Employee) dto.EmployeePublicResponse {
	return dto.EmployeePublicResponse{
		EmployeeBaseResponse: mapEmployeeBase(e),
		User:                 e.User,
	}
}

// --- Handlers ---

func GetEmployeeList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")
	deptID := ctx.Query("department_id")

	if page < 1 { page = 1 }
	if limit <= 0 { limit = 10 }

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
	
	if role == "admin" {
		res := make([]dto.EmployeeResponse, len(employees))
		for i := range employees {
			res[i] = toEmployeeResponse(employees[i])
		}
		utils.SuccessWithMeta(ctx, res, gin.H{"total": total, "page": page, "limit": limit})
	} else {
		res := make([]dto.EmployeePublicResponse, len(employees))
		for i := range employees {
			res[i] = toEmployeePublicResponse(employees[i])
		}
		utils.SuccessWithMeta(ctx, res, gin.H{"total": total, "page": page, "limit": limit})
	}
}

func GetEmployeeID(ctx *gin.Context) {
	var employee models.Employee
	if err := config.DB.Preload("Department").Preload("Position").Preload("Position.Department").Preload("User").First(&employee, ctx.Param("id")).Error; err != nil {
		utils.NotFound(ctx, "Không tìm thấy nhân viên!")
		return
	}

	role, _ := ctx.Get("role")
	if role == "admin" {
		utils.Success(ctx, toEmployeeResponse(employee))
	} else {
		utils.Success(ctx, toEmployeePublicResponse(employee))
	}
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

	hireDate, _ := parseDate(input.HireDate)
	dob, _ := parseDate(input.DateOfBirth)

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
	utils.Create(ctx, toEmployeeResponse(employee))
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
	if input.PositionID != nil { updates["position_id"] = *input.PositionID }
	if input.Salary != nil { updates["salary"] = *input.Salary }
	if input.HireDate != nil {
		hireDate, _ := parseDate(*input.HireDate)
		updates["hire_date"] = hireDate
	}
	if input.DateOfBirth != nil {
		dob, _ := parseDate(*input.DateOfBirth)
		updates["date_of_birth"] = dob
	}
	if input.Status != nil { updates["status"] = *input.Status }
	if input.AvatarURL != nil { updates["avatar_url"] = *input.AvatarURL }

	if err := config.DB.Model(&employee).Updates(updates).Error; err != nil {
		utils.InternalError(ctx, "Cập nhật thất bại")
		return
	}

	config.DB.Preload("Department").Preload("Position").Preload("Position.Department").Preload("User").First(&employee, employee.ID)
	utils.Success(ctx, toEmployeeResponse(employee))
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

	if err := config.DB.Delete(&employee).Error; err != nil {
		utils.InternalError(ctx, "Xoá nhân viên thất bại")
		return
	}

	utils.Success(ctx, "Xoá nhân viên thành công")
}
