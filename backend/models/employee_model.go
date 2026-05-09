package models

import (
	"time"

	"gorm.io/gorm"
)

type Status int

const (
	StatusInactive Status = 0
	StatusActive   Status = 1
)

type Employee struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Name         string         `json:"name" gorm:"not null;index"`
	Gender       string         `json:"gender" gorm:"type:enum('male','female','other');default:null"`
	DateOfBirth  *time.Time     `json:"date_of_birth"`
	Phone        string         `json:"phone" gorm:"unique;not null,index"`
	DepartmentID *uint          `json:"department_id"`
	Department   *Department    `json:"department" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Position     string         `json:"position"`
	Salary       float64        `json:"salary" gorm:"type:decimal(15,2);default:0"`
	HireDate     *time.Time     `json:"hire_date"`
	Status       Status         `json:"status" gorm:"default:1"`
	AvatarURL    string         `json:"avatar_url"`
	User         *User          `json:"user" gorm:"foreignKey:EmployeeID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}
