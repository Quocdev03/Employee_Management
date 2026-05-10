package models

import "time"

// Position là chức vụ thuộc một phòng ban cụ thể
type Position struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Name         string     `json:"name" gorm:"not null;size:150;index"`
	DepartmentID uint       `json:"department_id" gorm:"not null;index"`
	Department   *Department `json:"department,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    time.Time  `json:"created_at"`
}
