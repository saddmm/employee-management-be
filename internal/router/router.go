package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "github.com/saddam/employee-management-be/docs"
	"github.com/saddam/employee-management-be/internal/config"
	"github.com/saddam/employee-management-be/internal/handler"
	"github.com/saddam/employee-management-be/internal/middleware"
	"github.com/saddam/employee-management-be/internal/model"
)

type RouterConfig struct {
	Config       *config.Config
	AuthHandler  *handler.AuthHandler
	DeptHandler  *handler.DepartmentHandler
	EmpHandler   *handler.EmployeeHandler
	AuditHandler *handler.AuditLogHandler
}

func SetupRoutes(app *fiber.App, cfg *RouterConfig) {
	// Global Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(middleware.RateLimiter(cfg.Config))

	// Swagger documentation route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Get("/api-documentation", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html", fiber.StatusMovedPermanently)
	})
	app.Get("/api-documentation/*", fiberSwagger.WrapHandler)

	// Health check (Public)
	app.Get("/health", handler.HealthCheck)

	// API Group
	api := app.Group("/api")

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/register", cfg.AuthHandler.Register)
	auth.Post("/login", cfg.AuthHandler.Login)
	auth.Get("/me", middleware.JWTProtected(cfg.Config), cfg.AuthHandler.GetMe)

	// Protected Group for authenticated users
	jwtAuth := middleware.JWTProtected(cfg.Config)

	// Departments
	depts := api.Group("/departments", jwtAuth)
	depts.Get("/", cfg.DeptHandler.GetAll)
	depts.Get("/:id", cfg.DeptHandler.GetByID)
	depts.Post("/", middleware.RequireRole(model.RoleAdmin), cfg.DeptHandler.Create)
	depts.Put("/:id", middleware.RequireRole(model.RoleAdmin), cfg.DeptHandler.Update)
	depts.Delete("/:id", middleware.RequireRole(model.RoleAdmin), cfg.DeptHandler.Delete)

	// Employees
	employees := api.Group("/employees", jwtAuth)
	employees.Get("/export/csv", cfg.EmpHandler.ExportCSV)
	employees.Get("/", cfg.EmpHandler.GetAll)
	employees.Get("/:id", cfg.EmpHandler.GetByID)
	employees.Post("/", middleware.RequireRole(model.RoleAdmin), cfg.EmpHandler.Create)
	employees.Put("/:id", middleware.RequireRole(model.RoleAdmin), cfg.EmpHandler.Update)
	employees.Delete("/:id", middleware.RequireRole(model.RoleAdmin), cfg.EmpHandler.Delete)

	// Audit Logs (Admin only)
	auditLogs := api.Group("/audit-logs", jwtAuth, middleware.RequireRole(model.RoleAdmin))
	auditLogs.Get("/", cfg.AuditHandler.GetAll)
}
