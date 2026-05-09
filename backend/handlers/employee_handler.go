package handlers

import (
	"ChiQuoc/HocGolang/config"
	"ChiQuoc/HocGolang/dto"
	"ChiQuoc/HocGolang/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func parseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	// Thử format YYYY-MM-DD
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		// Thử format RFC3339
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
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "4"))
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
		Preload("User")

	if search != "" {
		words := strings.Fields(strings.TrimSpace(search))
		for _, w := range words {
			like := "%" + w + "%"
			query = query.Where("name LIKE ? OR phone LIKE ?", like, like)
		}
	}

	if deptID != "" {
		query = query.Where("department_id = ?", deptID)
	}

	var total int64
	query.Count(&total)

	var employees []models.Employee
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&employees)

	// Admin thấy đầy đủ, nhân viên thường ẩn salary
	role, _ := ctx.Get("role")
	if role == "admin" {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    employees,
			"total":   total,
			"page":    page,
			"limit":   limit,
		})
		return
	}

	public := make([]dto.EmployeePublicResponse, len(employees))
	for i, e := range employees {
		public[i] = toPublicResponse(e)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    public,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func GetEmployeeID(ctx *gin.Context) {
	var employee models.Employee

	if err := config.DB.
		Preload("Department").
		Preload("User").
		First(&employee, ctx.Param("id")).Error; err != nil {

		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Không tìm thấy nhân viên!",
		})
		return
	}

	// Admin thấy đầy đủ, nhân viên thường ẩn salary
	role, _ := ctx.Get("role")
	if role == "admin" {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "data": employee})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": toPublicResponse(employee)})
}

func CreateEmployee(ctx *gin.Context) {
	var input dto.CreateEmployeeInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	phone := strings.TrimSpace(input.Phone)
	if phone != "" {
		var exist models.Employee
		if err := config.DB.Where("phone = ?", phone).First(&exist).Error; err == nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "Số điện thoại đã tồn tại",
			})
			return
		}
	}

	if input.DepartmentID != nil {
		var dept models.Department
		if err := config.DB.First(&dept, *input.DepartmentID).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Phòng ban không tồn tại",
			})
			return
		}
	}

	hireDate, err := parseDate(input.HireDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Định dạng ngày vào làm không hợp lệ"})
		return
	}

	dob, err := parseDate(input.DateOfBirth)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Định dạng ngày sinh không hợp lệ"})
		return
	}

	employee := models.Employee{
		Name:         input.Name,
		Gender:       input.Gender,
		Phone:        phone,
		DepartmentID: input.DepartmentID,
		Position:     input.Position,
		Salary:       input.Salary,
		HireDate:     hireDate,
		DateOfBirth:  dob,
		Status:       models.StatusActive,
	}

	if err := config.DB.Create(&employee).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Tạo nhân viên thất bại",
		})
		return
	}

	config.DB.Preload("Department").First(&employee, employee.ID)

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    employee,
	})
}

func UpdateEmployee(ctx *gin.Context) {
	var employee models.Employee
	if err := config.DB.First(&employee, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Không tìm thấy nhân viên!"})
		return
	}

	var input dto.UpdateEmployeeInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Gender != nil {
		updates["gender"] = *input.Gender
	}
	if input.Phone != nil {
		phone := strings.TrimSpace(*input.Phone)
		if phone != "" {
			var exist models.Employee
			if err := config.DB.
				Where("phone = ? AND id != ?", phone, employee.ID).
				First(&exist).Error; err == nil {
				ctx.JSON(http.StatusConflict, gin.H{"success": false, "error": "Số điện thoại đã tồn tại"})
				return
			}
		}
		updates["phone"] = phone
	}
	if input.DepartmentID != nil {
		updates["department_id"] = input.DepartmentID
	}
	if input.Position != nil {
		updates["position"] = *input.Position
	}
	if input.Salary != nil {
		updates["salary"] = *input.Salary
	}
	if input.HireDate != nil {
		hireDate, err := parseDate(*input.HireDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Định dạng ngày vào làm không hợp lệ"})
			return
		}
		updates["hire_date"] = hireDate
	}
	if input.DateOfBirth != nil {
		dob, err := parseDate(*input.DateOfBirth)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Định dạng ngày sinh không hợp lệ"})
			return
		}
		updates["date_of_birth"] = dob
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	if err := config.DB.Model(&employee).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Cập nhật thất bại"})
		return
	}

	config.DB.Preload("Department").First(&employee, employee.ID)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": employee})
}

func DeleteEmployee(ctx *gin.Context) {
	var employee models.Employee
	if err := config.DB.First(&employee, ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Không tìm thấy nhân viên!"})
		return
	}

	// Xoá tài khoản associated nếu có
	config.DB.Where("employee_id = ?", employee.ID).Delete(&models.User{})

	// Cập nhật trạng thái thành nghỉ việc trước khi xoá mềm
	config.DB.Model(&employee).Update("status", models.StatusInactive)

	config.DB.Delete(&employee)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Xoá nhân viên thành công!",
	})
}
