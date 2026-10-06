package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/saddam/employee-management-be/internal/config"
	"github.com/saddam/employee-management-be/internal/database"
	"github.com/saddam/employee-management-be/internal/handler"
	"github.com/saddam/employee-management-be/internal/middleware"
	"github.com/saddam/employee-management-be/internal/repository"
	"github.com/saddam/employee-management-be/internal/router"
	"github.com/saddam/employee-management-be/internal/service"
)

// @title Employee Management System API
// @version 1.0
// @description RESTful API for Employee Management System built with Go Fiber, GORM, and MySQL.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.LoadConfig()

	// Connect to MySQL
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Auto-migrate schema
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Database migration error: %v", err)
	}

	// Seed default data
	if err := database.Seed(db); err != nil {
		log.Printf("Database seed warning: %v", err)
	}

	validate := validator.New()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	deptRepo := repository.NewDepartmentRepository(db)
	empRepo := repository.NewEmployeeRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)

	// Services
	auditService := service.NewAuditLogService(auditRepo)
	authService := service.NewAuthService(userRepo, cfg)
	deptService := service.NewDepartmentService(deptRepo, auditService)
	empService := service.NewEmployeeService(empRepo, deptRepo, auditService)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, validate)
	deptHandler := handler.NewDepartmentHandler(deptService, validate)
	empHandler := handler.NewEmployeeHandler(empService, validate)
	auditHandler := handler.NewAuditLogHandler(auditService)

	// Fiber App
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		AppName:      "Employee Management System API",
	})

	// Setup Routes
	router.SetupRoutes(app, &router.RouterConfig{
		Config:       cfg,
		AuthHandler:  authHandler,
		DeptHandler:  deptHandler,
		EmpHandler:   empHandler,
		AuditHandler: auditHandler,
	})

	addr := ":" + cfg.AppPort
	log.Printf("Server listening on http://localhost%s", addr)
	log.Printf("Swagger API documentation available at http://localhost%s/api-documentation (or /swagger/)", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
