package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/saddam/employee-management-be/internal/repository"
)

func (h *EmployeeHandler) ExportCSV(c *fiber.Ctx) error {
	search := c.Query("search")
	status := c.Query("status")
	sortBy := c.Query("sort_by", "name")
	sortOrder := c.Query("sort_order", "asc")

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
	}

	employees, err := h.empService.GetAllForExport(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{"ID", "Name", "Email", "Phone", "Position", "Department", "Status", "Joined At", "Created At"}
	if err := writer.Write(headers); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate CSV header",
		})
	}

	for _, emp := range employees {
		deptName := "-"
		if emp.Department != nil {
			deptName = emp.Department.Name
		}

		joinedAtStr := "-"
		if emp.JoinedAt != nil {
			joinedAtStr = emp.JoinedAt.Format("2006-01-02")
		}

		record := []string{
			strconv.FormatUint(uint64(emp.ID), 10),
			emp.Name,
			emp.Email,
			emp.Phone,
			emp.Position,
			deptName,
			string(emp.Status),
			joinedAtStr,
			emp.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(record); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to write CSV record",
			})
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to finalize CSV output",
		})
	}

	fileName := fmt.Sprintf("employees_%s.csv", time.Now().Format("20060102_150405"))
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	return c.Send(buf.Bytes())
}
