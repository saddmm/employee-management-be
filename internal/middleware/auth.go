package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/saddam/employee-management-be/internal/config"
	"github.com/saddam/employee-management-be/internal/model"
)

type JWTClaims struct {
	UserID uint       `json:"user_id"`
	Email  string     `json:"email"`
	Role   model.Role `json:"role"`
	jwt.RegisteredClaims
}

func JWTProtected(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
				Success: false,
				Message: "Authorization header is missing",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
				Success: false,
				Message: "Invalid authorization header format. Expected 'Bearer <token>'",
			})
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
				Success: false,
				Message: "Invalid or expired authorization token",
			})
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
				Success: false,
				Message: "Malformed token claims",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)
		c.Locals("claims", claims)

		return c.Next()
	}
}

func RequireRole(roles ...model.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleVal := c.Locals("user_role")
		if userRoleVal == nil {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Success: false,
				Message: "Access forbidden: no role assigned",
			})
		}

		userRole, ok := userRoleVal.(model.Role)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Success: false,
				Message: "Access forbidden: invalid role format",
			})
		}

		for _, r := range roles {
			if userRole == r {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
			Success: false,
			Message: "Access forbidden: insufficient permissions",
		})
	}
}
