package model

import "time"

type EmployeeStatus string

const (
	StatusActive   EmployeeStatus = "active"
	StatusInactive EmployeeStatus = "inactive"
)

type Employee struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	DepartmentID *uint          `gorm:"index" json:"department_id"`
	Department   *Department    `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Name         string         `gorm:"size:150;not null" json:"name" validate:"required,min=2,max=150"`
	Email        string         `gorm:"size:150;not null;uniqueIndex" json:"email" validate:"required,email,max=150"`
	Phone        string         `gorm:"size:20" json:"phone" validate:"omitempty,max=20"`
	Position     string         `gorm:"size:100;not null" json:"position" validate:"required,max=100"`
	Status       EmployeeStatus `gorm:"type:enum('active','inactive');default:'active';not null" json:"status" validate:"omitempty,oneof=active inactive"`
	JoinedAt     *time.Time     `gorm:"type:date" json:"joined_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
