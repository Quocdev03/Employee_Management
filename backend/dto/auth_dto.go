package dto

import "ChiQuoc/HocGolang/models"

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  AuthUserDTO `json:"user"`
}

type AuthUserDTO struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
}

type MeResponse struct {
	User     AuthUserDTO      `json:"user"`
	Employee *models.Employee `json:"employee"`
}
