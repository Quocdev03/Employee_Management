package models

import "time"

type Department struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"uniqueIndex;not null;size:100"`
	Positions []Position `json:"positions,omitempty" gorm:"foreignKey:DepartmentID"`
	CreatedAt time.Time  `json:"created_at"`
}
