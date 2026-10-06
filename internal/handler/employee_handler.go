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
	if int(total)%limit != 0 {
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
