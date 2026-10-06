package model

import "time"

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleViewer Role = "viewer"
	RoleUser   Role = "user"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:150;not null" json:"name" validate:"required,min=2,max=150"`
	Email     string    `gorm:"size:150;not null;uniqueIndex" json:"email" validate:"required,email,max=150"`
	Password  string    `gorm:"size:255;not null" json:"-" validate:"required,min=6"`
	Role      Role      `gorm:"type:enum('admin','viewer','user');default:'user';not null" json:"role" validate:"omitempty,oneof=admin viewer user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
