package middleware

import (
	"net/http"
	"strings"

	"backend-padel-go/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateRequest[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		
		// Bind JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request body", err.Error()))
			c.Abort()
			return
		}

		// Validar estructura
		if err := validate.Struct(req); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, getValidationError(err))
			}
			
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Validation failed", strings.Join(errors, "; ")))
			c.Abort()
			return
		}

		// Agregar request validado al contexto
		c.Set("validated_request", req)
		c.Next()
	}
}

func ValidateQuery[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		
		// Bind query parameters
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid query parameters", err.Error()))
			c.Abort()
			return
		}

		// Validar estructura
		if err := validate.Struct(req); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, getValidationError(err))
			}
			
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Validation failed", strings.Join(errors, "; ")))
			c.Abort()
			return
		}

		// Agregar query validado al contexto
		c.Set("validated_query", req)
		c.Next()
	}
}

func getValidationError(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return field + " is required"
	case "min":
		return field + " must be at least " + param
	case "max":
		return field + " must be at most " + param
	case "email":
		return field + " must be a valid email address"
	case "oneof":
		return field + " must be one of: " + param
	case "len":
		return field + " must be exactly " + param + " characters"
	default:
		return field + " is invalid"
	}
}
