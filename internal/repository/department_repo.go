package repository

import (
	"github.com/saddam/employee-management-be/internal/model"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	FindAll() ([]model.Department, error)
	FindByID(id uint) (*model.Department, error)
	FindByName(name string) (*model.Department, error)
	Create(dept *model.Department) error
	Update(dept *model.Department) error
	Delete(id uint) error
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) FindAll() ([]model.Department, error) {
	var depts []model.Department
	if err := r.db.Order("name ASC").Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *departmentRepository) FindByID(id uint) (*model.Department, error) {
	var dept model.Department
	if err := r.db.First(&dept, id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepository) FindByName(name string) (*model.Department, error) {
	var dept model.Department
	if err := r.db.Where("name = ?", name).First(&dept).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepository) Create(dept *model.Department) error {
	return r.db.Create(dept).Error
}

func (r *departmentRepository) Update(dept *model.Department) error {
	return r.db.Save(dept).Error
}

func (r *departmentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Department{}, id).Error
}
