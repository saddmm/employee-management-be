package repository

import (
	"fmt"
	"strings"

	"github.com/saddam/employee-management-be/internal/model"
	"gorm.io/gorm"
)

type EmployeeFilter struct {
	Search       string
	DepartmentID *uint
	Status       string
	SortBy       string
	SortOrder    string
	Page         int
	Limit        int
}

type EmployeeRepository interface {
	FindAll(filter EmployeeFilter) ([]model.Employee, int64, error)
	FindAllForExport(filter EmployeeFilter) ([]model.Employee, error)
	FindByID(id uint) (*model.Employee, error)
	FindByEmail(email string) (*model.Employee, error)
	Create(emp *model.Employee) error
	Update(emp *model.Employee) error
	Delete(id uint) error
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) applyFilter(db *gorm.DB, filter EmployeeFilter) *gorm.DB {
	query := db

	if filter.Search != "" {
		searchPattern := "%" + strings.TrimSpace(filter.Search) + "%"
		query = query.Where("(name LIKE ? OR email LIKE ? OR position LIKE ?)", searchPattern, searchPattern, searchPattern)
	}

	if filter.DepartmentID != nil && *filter.DepartmentID > 0 {
		query = query.Where("department_id = ?", *filter.DepartmentID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	return query
}

func (r *employeeRepository) FindAll(filter EmployeeFilter) ([]model.Employee, int64, error) {
	var employees []model.Employee
	var total int64

	baseQuery := r.applyFilter(r.db.Model(&model.Employee{}), filter)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Validate sort column to avoid SQL injection
	sortBy := "created_at"
	allowedSorts := map[string]bool{
		"name":       true,
		"email":      true,
		"position":   true,
		"status":     true,
		"joined_at":  true,
		"created_at": true,
	}
	if allowedSorts[filter.SortBy] {
		sortBy = filter.SortBy
	}

	sortOrder := "DESC"
	if strings.ToLower(filter.SortOrder) == "asc" {
		sortOrder = "ASC"
	}

	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	dataQuery := r.applyFilter(r.db.Model(&model.Employee{}), filter).
		Preload("Department").
		Order(orderClause).
		Offset(offset).
		Limit(limit)

	if err := dataQuery.Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

func (r *employeeRepository) FindAllForExport(filter EmployeeFilter) ([]model.Employee, error) {
	var employees []model.Employee

	sortBy := "name"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "ASC"
	if strings.ToLower(filter.SortOrder) == "desc" {
		sortOrder = "DESC"
	}

	query := r.applyFilter(r.db.Model(&model.Employee{}), filter).
		Preload("Department").
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	if err := query.Find(&employees).Error; err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *employeeRepository) FindByID(id uint) (*model.Employee, error) {
	var emp model.Employee
	if err := r.db.Preload("Department").First(&emp, id).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *employeeRepository) FindByEmail(email string) (*model.Employee, error) {
	var emp model.Employee
	if err := r.db.Where("email = ?", email).First(&emp).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *employeeRepository) Create(emp *model.Employee) error {
	return r.db.Create(emp).Error
}

func (r *employeeRepository) Update(emp *model.Employee) error {
	return r.db.Save(emp).Error
}

func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&model.Employee{}, id).Error
}
