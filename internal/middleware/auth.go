package middleware

import (
	"net/http"
	"strings"

	"backend-padel-go/internal/models"
	"backend-padel-go/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Authorization header required", "UNAUTHORIZED"))
			c.Abort()
			return
		}

		// Verificar formato "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid authorization header format", "UNAUTHORIZED"))
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validar token
		claims, err := services.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid token", "UNAUTHORIZED"))
			c.Abort()
			return
		}

		// Agregar información del usuario al contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User role not found", "UNAUTHORIZED"))
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Invalid user role type", "INTERNAL_ERROR"))
			c.Abort()
			return
		}

		// Verificar si el rol del usuario está en la lista de roles permitidos
		allowed := false
		for _, role := range roles {
			if roleStr == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, models.NewErrorResponse("Insufficient permissions", "FORBIDDEN"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func OwnerOrAdminRequired() gin.HandlerFunc {
	return RoleRequired("owner", "admin")
}

func AdminRequired() gin.HandlerFunc {
	return RoleRequired("admin")
}
