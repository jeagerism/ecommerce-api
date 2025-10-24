package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/jeagerism/ecommerce-api/util"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func RequireUserRole() echo.MiddlewareFunc {
	return AuthorizeRoles("user")
}

func RequireShopRole() echo.MiddlewareFunc {
	return AuthorizeRoles("shop")
}

func UserAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			skipPaths := map[string]bool{
				"/api/user/login":    true,
				"/api/user/register": true,
			}
			if skipPaths[c.Request().URL.Path] {
				// Skip authentication for paths that don't require it
				return next(c)
			}
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logrus.Warn("[Middleware.UserAuthMiddleware]: Missing authorization header")
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing authorization header"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logrus.Warn("[Middleware.UserAuthMiddleware]: Invalid authorization header format")
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid authorization header format"})
			}
			userSecret := []byte(os.Getenv("USER_SECRET")) // Replace with your actual secret key
			tokenString := parts[1]
			claims, err := util.ParseJWT(tokenString, userSecret)
			if err != nil || claims.Role != "user" {
				logrus.Warnf("Middleware.UserAuthMiddleware]: Invalid token: %v", err)
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
			}

			// Put all claims to context
			c.Set("userID", claims.ID)
			c.Set("userRole", claims.Role)

			return next(c)
		}
	}
}

func ShopAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			skipPaths := map[string]bool{
				"/api/shop/login":    true,
				"/api/shop/register": true,
			}
			if skipPaths[c.Request().URL.Path] {
				// Skip authentication for paths that don't require it
				return next(c)
			}
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logrus.Warn("[Middleware.ShopAuthMiddleware]: Missing authorization header")
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing authorization header"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logrus.Warn("[Middleware.ShopAuthMiddleware]: Invalid authorization header format")
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid authorization header format"})
			}
			shopSecret := []byte(os.Getenv("SHOP_SECRET")) // Replace with your actual secret key
			tokenString := parts[1]
			claims, err := util.ParseJWT(tokenString, shopSecret)
			if err != nil || claims.Role != "shop" {
				logrus.Warnf("Middleware.ShopAuthMiddleware]: Invalid token: %v", err)
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
			}

			// Put all claims to context
			c.Set("shopID", claims.ID)
			c.Set("userRole", claims.Role)

			return next(c)
		}
	}
}

// AuthorizeRoles only allows if role in allowedRoles
func AuthorizeRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			skipPaths := map[string]bool{
				"/api/user/login":    true,
				"/api/user/register": true,
				"/api/shop/login":    true,
				"/api/shop/register": true,
			}
			if skipPaths[c.Request().URL.Path] {
				// Skip authentication for paths that don't require it
				return next(c)
			}
			roleVal := c.Get("userRole")
			if roleVal == nil {
				return c.JSON(http.StatusForbidden, echo.Map{"error": "missing role in context"})
			}

			userRole := roleVal.(string)

			for _, role := range allowedRoles {
				if role == userRole {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, echo.Map{"error": "forbidden: role not allowed"})
		}
	}
}
