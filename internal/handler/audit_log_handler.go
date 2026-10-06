package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/saddam/employee-management-be/internal/service"
)

type AuditLogHandler struct {
	auditService service.AuditLogService
}

func NewAuditLogHandler(auditService service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{
		auditService: auditService,
	}
}

func (h *AuditLogHandler) GetAll(c *fiber.Ctx) error {
	entity := c.Query("entity")
	action := c.Query("action")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	logs, total, err := h.auditService.GetAll(entity, action, page, limit)
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
		"data":    logs,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
