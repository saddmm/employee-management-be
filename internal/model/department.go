package model

import "time"

type Department struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"size:100;not null;uniqueIndex" json:"name" validate:"required,min=2,max=100"`
	Description string     `gorm:"type:text" json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Employees   []Employee `gorm:"foreignKey:DepartmentID;constraint:OnDelete:SET NULL" json:"employees,omitempty"`
}
