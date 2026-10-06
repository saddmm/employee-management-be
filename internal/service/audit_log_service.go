package service

import (
	"encoding/json"

	"github.com/saddam/employee-management-be/internal/model"
	"github.com/saddam/employee-management-be/internal/repository"
)

type AuditLogService interface {
	Log(userID *uint, entity string, entityID uint, action model.AuditAction, oldData interface{}, newData interface{}) error
	GetAll(entity, action string, page, limit int) ([]model.AuditLog, int64, error)
}

type auditLogService struct {
	auditRepo repository.AuditLogRepository
}

func NewAuditLogService(auditRepo repository.AuditLogRepository) AuditLogService {
	return &auditLogService{
		auditRepo: auditRepo,
	}
}

func (s *auditLogService) Log(userID *uint, entity string, entityID uint, action model.AuditAction, oldData interface{}, newData interface{}) error {
	var oldJSON, newJSON string

	if oldData != nil {
		if bytes, err := json.Marshal(oldData); err == nil {
			oldJSON = string(bytes)
		}
	}

	if newData != nil {
		if bytes, err := json.Marshal(newData); err == nil {
			newJSON = string(bytes)
		}
	}

	logEntry := &model.AuditLog{
		UserID:   userID,
		Entity:   entity,
		EntityID: entityID,
		Action:   action,
		OldData:  oldJSON,
		NewData:  newJSON,
	}

	return s.auditRepo.Create(logEntry)
}

func (s *auditLogService) GetAll(entity, action string, page, limit int) ([]model.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.auditRepo.FindAll(entity, action, page, limit)
}
