package dto

type CreateUserInput struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	RoleID     uint   `json:"role_id" binding:"required,oneof=1 2"`
	EmployeeID *uint  `json:"employee_id"`
}

type UpdateUserInput struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RoleID     uint   `json:"role_id"`
	EmployeeID *uint  `json:"employee_id"`
}
