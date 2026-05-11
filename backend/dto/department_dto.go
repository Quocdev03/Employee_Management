package dto

// DepartmentRequest dùng để bind dữ liệu khi tạo hoặc cập nhật phòng ban
type DepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}
