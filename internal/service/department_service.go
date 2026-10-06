package service

import (
	"errors"
	"fmt"

	"github.com/saddam/employee-management-be/internal/model"
	"github.com/saddam/employee-management-be/internal/repository"
	"gorm.io/gorm"
)

type AuditLogger interface {
	Log(userID *uint, entity string, entityID uint, action model.AuditAction, oldData interface{}, newData interface{}) error
}

type DepartmentService interface {
	GetAll() ([]model.Department, error)
	GetByID(id uint) (*model.Department, error)
	Create(userID uint, req *CreateDepartmentRequest) (*model.Department, error)
	Update(userID uint, id uint, req *UpdateDepartmentRequest) (*model.Department, error)
	Delete(userID uint, id uint) error
}

type departmentService struct {
	deptRepo repository.DepartmentRepository
	auditSvc AuditLogger
}

type CreateDepartmentRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateDepartmentRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
}

func NewDepartmentService(deptRepo repository.DepartmentRepository, auditSvc AuditLogger) DepartmentService {
	return &departmentService{
		deptRepo: deptRepo,
		auditSvc: auditSvc,
	}
}

func (s *departmentService) GetAll() ([]model.Department, error) {
	return s.deptRepo.FindAll()
}

func (s *departmentService) GetByID(id uint) (*model.Department, error) {
	dept, err := s.deptRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, err
	}
	return dept, nil
}

func (s *departmentService) Create(userID uint, req *CreateDepartmentRequest) (*model.Department, error) {
	existing, err := s.deptRepo.FindByName(req.Name)
	if err == nil && existing != nil {
		return nil, errors.New("department with this name already exists")
	}

	dept := &model.Department{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.deptRepo.Create(dept); err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	if s.auditSvc != nil {
		_ = s.auditSvc.Log(&userID, "department", dept.ID, model.ActionCreate, nil, dept)
	}

	return dept, nil
}

func (s *departmentService) Update(userID uint, id uint, req *UpdateDepartmentRequest) (*model.Department, error) {
	dept, err := s.deptRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, err
	}

	if req.Name != dept.Name {
		existing, err := s.deptRepo.FindByName(req.Name)
		if err == nil && existing != nil && existing.ID != dept.ID {
			return nil, errors.New("department with this name already exists")
		}
	}

	oldCopy := *dept
	dept.Name = req.Name
	dept.Description = req.Description

	if err := s.deptRepo.Update(dept); err != nil {
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	if s.auditSvc != nil {
		_ = s.auditSvc.Log(&userID, "department", dept.ID, model.ActionUpdate, oldCopy, dept)
	}

	return dept, nil
}

func (s *departmentService) Delete(userID uint, id uint) error {
	dept, err := s.deptRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("department not found")
		}
		return err
	}

	if err := s.deptRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	if s.auditSvc != nil {
		_ = s.auditSvc.Log(&userID, "department", dept.ID, model.ActionDelete, dept, nil)
	}

	return nil
}
