package middleware

import (
	"strings"

	"sm-client-backend/auth"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// Protected verifies the JWT and validates the layer 2 security (matching slug)
func Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		tokenString := parts[1]

		// Layer 1: Verify Token
		token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
			return auth.SecretKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(*auth.JWTClaim)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Layer 2: Match Slug
		// The route is expected to be like /api/vehicles/:slug
		requestSlug := c.Params("slug")
		if requestSlug != "" && requestSlug != claims.Slug {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: You do not have permission to access this resource",
			})
		}

		// Pass the claims to the next handler if needed
		c.Locals("user_id", claims.UserID)
		c.Locals("slug", claims.Slug)

		return c.Next()
	}
}
