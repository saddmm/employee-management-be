package database

import (
	"fmt"
	"log"
	"time"

	"github.com/saddam/employee-management-be/internal/config"
	"github.com/saddam/employee-management-be/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	logLevel := logger.Info
	if cfg.AppEnv == "production" {
		logLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Println("Database connection established successfully")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	err := db.AutoMigrate(
		&model.Department{},
		&model.User{},
		&model.Employee{},
		&model.AuditLog{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("Database migration completed")
	return nil
}

func Seed(db *gorm.DB) error {
	log.Println("Seeding initial data...")

	// Seed Admin User
	var count int64
	db.Model(&model.User{}).Where("email = ?", "admin@example.com").Count(&count)
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash seed admin password: %w", err)
		}

		admin := model.User{
			Name:     "Administrator",
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			Role:     model.RoleAdmin,
		}
		if err := db.Create(&admin).Error; err != nil {
			return fmt.Errorf("failed to seed admin user: %w", err)
		}
		log.Println("Admin user seeded (admin@example.com / admin123)")
	}

	// Seed Viewer User
	db.Model(&model.User{}).Where("email = ?", "viewer@example.com").Count(&count)
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("viewer123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash seed viewer password: %w", err)
		}

		viewer := model.User{
			Name:     "Viewer Demo",
			Email:    "viewer@example.com",
			Password: string(hashedPassword),
			Role:     model.RoleViewer,
		}
		if err := db.Create(&viewer).Error; err != nil {
			return fmt.Errorf("failed to seed viewer user: %w", err)
		}
		log.Println("Viewer user seeded (viewer@example.com / viewer123)")
	}

	// Seed Departments
	departments := []model.Department{
		{Name: "Engineering", Description: "Software engineering and technical infrastructure"},
		{Name: "Human Resources", Description: "Talent acquisition and people management"},
		{Name: "Product Management", Description: "Product strategy and user research"},
		{Name: "Marketing", Description: "Brand management and growth marketing"},
		{Name: "Finance", Description: "Financial planning and accounting"},
	}

	for _, dept := range departments {
		var deptCount int64
		db.Model(&model.Department{}).Where("name = ?", dept.Name).Count(&deptCount)
		if deptCount == 0 {
			if err := db.Create(&dept).Error; err != nil {
				return fmt.Errorf("failed to seed department %s: %w", dept.Name, err)
			}
		}
	}
	log.Println("Departments seeded")

	// Get department IDs for employee relationships
	var engDept, hrDept, prodDept, mktDept, finDept model.Department
	db.Where("name = ?", "Engineering").First(&engDept)
	db.Where("name = ?", "Human Resources").First(&hrDept)
	db.Where("name = ?", "Product Management").First(&prodDept)
	db.Where("name = ?", "Marketing").First(&mktDept)
	db.Where("name = ?", "Finance").First(&finDept)

	parseDate := func(d string) *time.Time {
		t, _ := time.Parse("2006-01-02", d)
		return &t
	}

	// Seed 10 Dummy Employees
	dummyEmployees := []model.Employee{
		{
			Name:         "Sarah Connor",
			Email:        "sarah.connor@example.com",
			Phone:        "+1 555-0101",
			Position:     "Principal Software Engineer",
			DepartmentID: &engDept.ID,
			Status:       model.StatusActive,
			JoinedAt:     parseDate("2023-01-15"),
		},
		{
			Name:         "Michael Chang",
			Email:        "michael.chang@example.com",
			Phone:        "+1 555-0102",
			Position:     "Senior Frontend Developer",
			DepartmentID: &engDept.ID,
			Status:       model.StatusActive,
			JoinedAt:     parseDate("2023-03-20"),
		},
		{
			Name:         "Emma Watson",
			Email:        "emma.watson@example.com",
			Phone:        "+1 555-0103",
			Position:     "DevOps & Cloud Specialist",
			DepartmentID: &engDept.ID,
			Status:       model.StatusActive,
			JoinedAt:     parseDate("2023-06-01"),
		},
		{
			Name:         "David Miller",
			Email:        "david.miller@example.com",
			Phone:        "+1 555-0104",
			Position:     "Head of People & Talent",
			DepartmentID: &hrDept.ID,
			Status:       model.StatusActive,
			JoinedAt:     parseDate("2022-11-10"),
		},
		{
			Name:         "Jessica Alba",
			Email:        "jessica.alba@example.com",
			Phone:        "+1 555-0105",
			Position:     "HR Operations Specialist",
			DepartmentID: &hrDept.ID,
			Status:       model.StatusActive,
			JoinedAt:     parseDate("2024-02-01"),
		},
		{
			Name:         "Robert Downey",
			Email:        "robert.downey@example.com",
			Phone:        "+1 555-0106",
			Position:     "Group Product Manager",
			DepartmentID: &prodDept.ID,
			Status:       model.StatusActive,
			JoinedAt:     parseDate("2022-08-15"),
		},
		{
			Name:         "Sophia Loren",
			Email:        "sophia.loren@example.com",
			Phone:        "+1 555-0107",
			Position:     "Lead UX/UI Researcher",
			DepartmentID: &prodDept.ID,
			Status:       model.StatusInactive,
			JoinedAt:     parseDate("2023-09-01"),
		},
		{
			Name:         "Alexander Wright",
			Email:        "alexander.wright@example.com",
			Phone:        "+1 555-0108",
			Position:     "Growth Marketing Director",
			DepartmentID: &mktDept.ID,
			Status:       model.StatusActive,
			JoinedAt:     parseDate("2023-04-12"),
		},
		{
			Name:         "Olivia Wilde",
			Email:        "olivia.wilde@example.com",
			Phone:        "+1 555-0109",
			Position:     "Content & Brand Strategist",
			DepartmentID: &mktDept.ID,
			Status:       model.StatusActive,
			JoinedAt:     parseDate("2024-01-10"),
		},
		{
			Name:         "James Anderson",
			Email:        "james.anderson@example.com",
			Phone:        "+1 555-0110",
			Position:     "Senior Financial Analyst",
			DepartmentID: &finDept.ID,
			Status:       model.StatusActive,
			JoinedAt:     parseDate("2023-05-18"),
		},
	}

	for _, emp := range dummyEmployees {
		var empCount int64
		db.Model(&model.Employee{}).Where("email = ?", emp.Email).Count(&empCount)
		if empCount == 0 {
			if err := db.Create(&emp).Error; err != nil {
				return fmt.Errorf("failed to seed employee %s: %w", emp.Name, err)
			}
		}
	}
	log.Println("10 Dummy employees seeded successfully")

	return nil
}
