package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	RoleID       uint      `json:"role_id" gorm:"not null;default:2"`
	Role         Role      `json:"role" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	EmployeeID   *uint     `json:"employee_id" gorm:"uniqueIndex"`
	Employee     *Employee `json:"employee"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}
