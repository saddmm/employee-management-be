package handler

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/saddam/employee-management-be/internal/service"
)

type DepartmentHandler struct {
	deptService service.DepartmentService
	validator   *validator.Validate
}

func NewDepartmentHandler(deptService service.DepartmentService, validate *validator.Validate) *DepartmentHandler {
	return &DepartmentHandler{
		deptService: deptService,
		validator:   validate,
	}
}

// GetAll godoc
// @Summary List all departments
// @Description Returns list of all departments
// @Tags Departments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} middleware.ErrorResponse
// @Router /api/departments [get]
func (h *DepartmentHandler) GetAll(c *fiber.Ctx) error {
	departments, err := h.deptService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    departments,
	})
}

// GetByID godoc
// @Summary Get department by ID
// @Description Returns single department by ID
// @Tags Departments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /api/departments/{id} [get]
func (h *DepartmentHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid department ID",
		})
	}

	dept, err := h.deptService.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dept,
	})
}

// Create godoc
// @Summary Create a department
// @Description Creates a new department (Admin only)
// @Tags Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateDepartmentRequest true "Department info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Router /api/departments [post]
func (h *DepartmentHandler) Create(c *fiber.Ctx) error {
	var req service.CreateDepartmentRequest
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
	dept, err := h.deptService.Create(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Department created successfully",
		"data":    dept,
	})
}

// Update godoc
// @Summary Update department
// @Description Updates an existing department (Admin only)
// @Tags Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Param request body service.UpdateDepartmentRequest true "Department update info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /api/departments/{id} [put]
func (h *DepartmentHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid department ID",
		})
	}

	var req service.UpdateDepartmentRequest
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
	dept, err := h.deptService.Update(userID, uint(id), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Department updated successfully",
		"data":    dept,
	})
}

// Delete godoc
// @Summary Delete department
// @Description Deletes a department by ID (Admin only)
// @Tags Departments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /api/departments/{id} [delete]
func (h *DepartmentHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid department ID",
		})
	}

	userID := c.Locals("user_id").(uint)
	if err := h.deptService.Delete(userID, uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Department deleted successfully",
	})
}
