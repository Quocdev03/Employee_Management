package dto

import (
	"ChiQuoc/HocGolang/models"
	"time"
)

// CreateEmployeeInput: Dữ liệu đầu vào khi tạo nhân viên
type CreateEmployeeInput struct {
	Name         string  `json:"name" binding:"required"`
	Gender       string  `json:"gender"`
	DateOfBirth  string  `json:"date_of_birth"`
	Phone        string  `json:"phone" binding:"required"`
	DepartmentID *uint   `json:"department_id"`
	PositionID   *uint   `json:"position_id"`
	Salary       float64 `json:"salary" binding:"required,gt=0"`
	HireDate     string  `json:"hire_date"`
	AvatarURL    string  `json:"avatar_url"`
}

// UpdateEmployeeInput: Dữ liệu đầu vào khi cập nhật nhân viên
type UpdateEmployeeInput struct {
	Name         *string  `json:"name"`
	Gender       *string  `json:"gender"`
	DateOfBirth  *string  `json:"date_of_birth"`
	Phone        *string  `json:"phone"`
	DepartmentID *uint    `json:"department_id"`
	PositionID   *uint    `json:"position_id"`
	Salary       *float64 `json:"salary"`
	HireDate     *string  `json:"hire_date"`
	Status       *string  `json:"status"`
	AvatarURL    *string  `json:"avatar_url"`
}

// EmployeeBaseResponse: Chứa các trường chung của nhân viên
type EmployeeBaseResponse struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Gender      string             `json:"gender"`
	Email       string             `json:"email"`
	DateOfBirth *time.Time         `json:"date_of_birth"`
	Phone       string             `json:"phone"`
	AvatarURL   string             `json:"avatar_url"`
	Status      models.Status      `json:"status"`
	HireDate    *time.Time         `json:"hire_date"`
	Department  *models.Department `json:"department"`
	Position    *models.Position   `json:"position"`
}

// EmployeeResponse: Trả về cho Admin (có lương)
type EmployeeResponse struct {
	EmployeeBaseResponse
	Salary float64      `json:"salary"`
	User   *models.User `json:"user,omitempty"`
}

// EmployeePublicResponse: Trả về cho Nhân viên (không có lương)
type EmployeePublicResponse struct {
	EmployeeBaseResponse
	User *models.User `json:"user,omitempty"`
}
