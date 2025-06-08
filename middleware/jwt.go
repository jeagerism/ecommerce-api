package middleware

import (
	"net/http"
	"strings"

	"github.com/jeagerism/ecommerce-api/util"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func RoleAuthMiddleware(secret []byte, expectedRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logrus.Warn("[RoleAuthMiddleware] Missing authorization header")
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing authorization header"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logrus.Warn("[RoleAuthMiddleware] Invalid authorization header format")
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid authorization header format"})
			}

			tokenString := parts[1]
			claims, err := util.ParseJWT(tokenString, secret) // claims is *JWTClaims
			if err != nil {
				logrus.Warnf("[RoleAuthMiddleware] Invalid token: %v", err)
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
			}

			// Check role from claims.Role
			if claims.Role != expectedRole {
				logrus.Warnf("[RoleAuthMiddleware] Forbidden: expected role %s, got %s", expectedRole, claims.Role)
				return c.JSON(http.StatusForbidden, echo.Map{"error": "forbidden"})
			}

			// Extract ID from claims based on role
			if expectedRole == "shop" {
				c.Set("shopID", claims.ID)
			} else if expectedRole == "user" {
				c.Set("userID", claims.ID)
			} else {
				logrus.Warnf("[RoleAuthMiddleware] Unknown role: %s", claims.Role)
				return c.JSON(http.StatusForbidden, echo.Map{"error": "unknown role"})
			}

			return next(c)
		}
	}
}
