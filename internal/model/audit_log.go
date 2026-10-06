package model

import "time"

type AuditAction string

const (
	ActionCreate AuditAction = "create"
	ActionUpdate AuditAction = "update"
	ActionDelete AuditAction = "delete"
)

type AuditLog struct {
	ID        uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    *uint       `gorm:"index" json:"user_id"`
	User      *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Entity    string      `gorm:"size:50;not null;index" json:"entity"`
	EntityID  uint        `gorm:"not null;index" json:"entity_id"`
	Action    AuditAction `gorm:"type:enum('create','update','delete');not null" json:"action"`
	OldData   string      `gorm:"type:json" json:"old_data,omitempty"`
	NewData   string      `gorm:"type:json" json:"new_data,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}
