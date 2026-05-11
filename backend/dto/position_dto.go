package dto

// CreatePositionRequest dùng để bind dữ liệu khi thêm chức vụ mới
type CreatePositionRequest struct {
	Name         string `json:"name" binding:"required"`
	DepartmentID uint   `json:"department_id" binding:"required"`
}

// UpdatePositionRequest dùng để bind dữ liệu khi cập nhật tên chức vụ
type UpdatePositionRequest struct {
	Name string `json:"name" binding:"required"`
}
