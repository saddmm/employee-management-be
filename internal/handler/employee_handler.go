package handler

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/saddam/employee-management-be/internal/repository"
	"github.com/saddam/employee-management-be/internal/service"
)

type EmployeeHandler struct {
	empService service.EmployeeService
	validator  *validator.Validate
}

func NewEmployeeHandler(empService service.EmployeeService, validate *validator.Validate) *EmployeeHandler {
	return &EmployeeHandler{
		empService: empService,
		validator:  validate,
	}
}

// GetAll godoc
// @Summary List employees
// @Description Searches, filters, sorts and paginates employees list
// @Tags Employees
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search query by name, email, position"
// @Param department_id query int false "Department ID filter"
// @Param status query string false "Status filter (active / inactive)"
// @Param sort_by query string false "Sort field (name, email, position, status, joined_at, created_at)" default(created_at)
// @Param sort_order query string false "Sort order (asc / desc)" default(desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size limit" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} middleware.ErrorResponse
// @Router /api/employees [get]
func (h *EmployeeHandler) GetAll(c *fiber.Ctx) error {
	search := c.Query("search")
	status := c.Query("status")
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	var deptIDPtr *uint
	if deptStr := c.Query("department_id"); deptStr != "" {
		if id, err := strconv.ParseUint(deptStr, 10, 32); err == nil {
			uid := uint(id)
			deptIDPtr = &uid
		}
	}

	filter := repository.EmployeeFilter{
		Search:       search,
		DepartmentID: deptIDPtr,
		Status:       status,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		Page:         page,
		Limit:        limit,
	}

	employees, total, err := h.empService.GetAll(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	totalPages := int(total) / limit
	if limit > 0 && int(total)%limit != 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    employees,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetByID godoc
// @Summary Get employee by ID
// @Description Returns employee details by ID
// @Tags Employees
// @Produce json
// @Security BearerAuth
// @Param id path int true "Employee ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /api/employees/{id} [get]
func (h *EmployeeHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid employee ID",
		})
	}

	emp, err := h.empService.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    emp,
	})
}

// Create godoc
// @Summary Create an employee
// @Description Creates a new employee (Admin only)
// @Tags Employees
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateEmployeeRequest true "Employee creation data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Router /api/employees [post]
func (h *EmployeeHandler) Create(c *fiber.Ctx) error {
	var req service.CreateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	userID := c.Locals("user_id").(uint)
	emp, err := h.empService.Create(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Employee created successfully",
		"data":    emp,
	})
}

// Update godoc
// @Summary Update an employee
// @Description Updates an existing employee (Admin only)
// @Tags Employees
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Employee ID"
// @Param request body service.UpdateEmployeeRequest true "Employee update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /api/employees/{id} [put]
func (h *EmployeeHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid employee ID",
		})
	}

	var req service.UpdateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	userID := c.Locals("user_id").(uint)
	emp, err := h.empService.Update(userID, uint(id), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Employee updated successfully",
		"data":    emp,
	})
}

// Delete godoc
// @Summary Delete an employee
// @Description Deletes an employee by ID (Admin only)
// @Tags Employees
// @Produce json
// @Security BearerAuth
// @Param id path int true "Employee ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /api/employees/{id} [delete]
func (h *EmployeeHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid employee ID",
		})
	}

	userID := c.Locals("user_id").(uint)
	if err := h.empService.Delete(userID, uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Employee deleted successfully",
	})
}
