package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/saddam/employee-management-be/internal/model"
	"github.com/saddam/employee-management-be/internal/repository"
	"gorm.io/gorm"
)

type CreateEmployeeRequest struct {
	DepartmentID *uint                `json:"department_id"`
	Name         string               `json:"name" validate:"required,min=2,max=150"`
	Email        string               `json:"email" validate:"required,email,max=150"`
	Phone        string               `json:"phone" validate:"omitempty,max=20"`
	Position     string               `json:"position" validate:"required,max=100"`
	Status       model.EmployeeStatus `json:"status" validate:"omitempty,oneof=active inactive"`
	JoinedAt     *string              `json:"joined_at"`
}

type UpdateEmployeeRequest struct {
	DepartmentID *uint                `json:"department_id"`
	Name         string               `json:"name" validate:"required,min=2,max=150"`
	Email        string               `json:"email" validate:"required,email,max=150"`
	Phone        string               `json:"phone" validate:"omitempty,max=20"`
	Position     string               `json:"position" validate:"required,max=100"`
	Status       model.EmployeeStatus `json:"status" validate:"omitempty,oneof=active inactive"`
	JoinedAt     *string              `json:"joined_at"`
}

type EmployeeService interface {
	GetAll(filter repository.EmployeeFilter) ([]model.Employee, int64, error)
	GetAllForExport(filter repository.EmployeeFilter) ([]model.Employee, error)
	GetByID(id uint) (*model.Employee, error)
	Create(userID uint, req *CreateEmployeeRequest) (*model.Employee, error)
	Update(userID uint, id uint, req *UpdateEmployeeRequest) (*model.Employee, error)
	Delete(userID uint, id uint) error
}

type employeeService struct {
	empRepo  repository.EmployeeRepository
	deptRepo repository.DepartmentRepository
	auditSvc AuditLogger
}

func NewEmployeeService(empRepo repository.EmployeeRepository, deptRepo repository.DepartmentRepository, auditSvc AuditLogger) EmployeeService {
	return &employeeService{
		empRepo:  empRepo,
		deptRepo: deptRepo,
		auditSvc: auditSvc,
	}
}

func (s *employeeService) GetAll(filter repository.EmployeeFilter) ([]model.Employee, int64, error) {
	return s.empRepo.FindAll(filter)
}

func (s *employeeService) GetAllForExport(filter repository.EmployeeFilter) ([]model.Employee, error) {
	return s.empRepo.FindAllForExport(filter)
}

func (s *employeeService) GetByID(id uint) (*model.Employee, error) {
	emp, err := s.empRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}
	return emp, nil
}

func (s *employeeService) Create(userID uint, req *CreateEmployeeRequest) (*model.Employee, error) {
	existing, err := s.empRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("employee with this email already exists")
	}

	if req.DepartmentID != nil && *req.DepartmentID > 0 {
		if _, err := s.deptRepo.FindByID(*req.DepartmentID); err != nil {
			return nil, errors.New("specified department does not exist")
		}
	}

	status := model.StatusActive
	if req.Status != "" {
		status = req.Status
	}

	var joinedTime *time.Time
	if req.JoinedAt != nil && *req.JoinedAt != "" {
		parsed, err := time.Parse("2006-01-02", *req.JoinedAt)
		if err != nil {
			return nil, errors.New("invalid joined_at date format. Expected YYYY-MM-DD")
		}
		joinedTime = &parsed
	}

	emp := &model.Employee{
		DepartmentID: req.DepartmentID,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Position:     req.Position,
		Status:       status,
		JoinedAt:     joinedTime,
	}

	if err := s.empRepo.Create(emp); err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	// Reload with relations for audit and response
	savedEmp, _ := s.empRepo.FindByID(emp.ID)
	if savedEmp != nil {
		emp = savedEmp
	}

	if s.auditSvc != nil {
		_ = s.auditSvc.Log(&userID, "employee", emp.ID, model.ActionCreate, nil, emp)
	}

	return emp, nil
}

func (s *employeeService) Update(userID uint, id uint, req *UpdateEmployeeRequest) (*model.Employee, error) {
	emp, err := s.empRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}

	if req.Email != emp.Email {
		existing, err := s.empRepo.FindByEmail(req.Email)
		if err == nil && existing != nil && existing.ID != emp.ID {
			return nil, errors.New("employee with this email already exists")
		}
	}

	if req.DepartmentID != nil && *req.DepartmentID > 0 {
		if _, err := s.deptRepo.FindByID(*req.DepartmentID); err != nil {
			return nil, errors.New("specified department does not exist")
		}
	}

	var joinedTime *time.Time
	if req.JoinedAt != nil && *req.JoinedAt != "" {
		parsed, err := time.Parse("2006-01-02", *req.JoinedAt)
		if err != nil {
			return nil, errors.New("invalid joined_at date format. Expected YYYY-MM-DD")
		}
		joinedTime = &parsed
	}

	oldCopy := *emp

	emp.DepartmentID = req.DepartmentID
	emp.Name = req.Name
	emp.Email = req.Email
	emp.Phone = req.Phone
	emp.Position = req.Position
	if req.Status != "" {
		emp.Status = req.Status
	}
	emp.JoinedAt = joinedTime

	if err := s.empRepo.Update(emp); err != nil {
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	updatedEmp, _ := s.empRepo.FindByID(emp.ID)
	if updatedEmp != nil {
		emp = updatedEmp
	}

	if s.auditSvc != nil {
		_ = s.auditSvc.Log(&userID, "employee", emp.ID, model.ActionUpdate, oldCopy, emp)
	}

	return emp, nil
}

func (s *employeeService) Delete(userID uint, id uint) error {
	emp, err := s.empRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("employee not found")
		}
		return err
	}

	if err := s.empRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	if s.auditSvc != nil {
		_ = s.auditSvc.Log(&userID, "employee", emp.ID, model.ActionDelete, emp, nil)
	}

	return nil
}
