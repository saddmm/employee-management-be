package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// HealthCheck godoc
// @Summary Check server health status
// @Description Returns the health status and UTC timestamp of the server
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
		Service:   "employee-management-be",
	})
}
