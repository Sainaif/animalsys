package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateRequest validates the request body against a struct
func ValidateRequest(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&obj); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
				"message": err.Error(),
			})
			return
		}

		if err := validate.Struct(obj); err != nil {
			validationErrors := formatValidationErrors(err)
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Validation failed",
				"errors": validationErrors,
			})
			return
		}

		c.Next()
	}
}

// formatValidationErrors formats validation errors for response
func formatValidationErrors(err error) []map[string]string {
	var errors []map[string]string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errors = append(errors, map[string]string{
				"field":   fieldError.Field(),
				"message": getErrorMessage(fieldError),
				"tag":     fieldError.Tag(),
			})
		}
	}

	return errors
}

// getErrorMessage returns a human-readable error message
func getErrorMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, fe.Param())
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "numeric":
		return fmt.Sprintf("%s must be a number", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// SanitizeInput sanitizes user input to prevent injection attacks
func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get raw body
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Next()
			return
		}

		// Sanitize all string values
		sanitizeMap(body)

		// Set sanitized body back to context
		c.Set("sanitized_body", body)
		c.Next()
	}
}

// sanitizeMap recursively sanitizes a map
func sanitizeMap(m map[string]interface{}) {
	for key, value := range m {
		switch v := value.(type) {
		case string:
			m[key] = sanitizeString(v)
		case map[string]interface{}:
			sanitizeMap(v)
		case []interface{}:
			sanitizeSlice(v)
		}
	}
}

// sanitizeSlice recursively sanitizes a slice
func sanitizeSlice(s []interface{}) {
	for i, value := range s {
		switch v := value.(type) {
		case string:
			s[i] = sanitizeString(v)
		case map[string]interface{}:
			sanitizeMap(v)
		case []interface{}:
			sanitizeSlice(v)
		}
	}
}

// sanitizeString removes potentially dangerous characters
func sanitizeString(s string) string {
	// Remove null bytes
	s = strings.ReplaceAll(s, "\x00", "")

	// Trim whitespace
	s = strings.TrimSpace(s)

	// Remove MongoDB operators (basic protection)
	dangerousPatterns := []string{"$where", "$ne", "$gt", "$lt", "$gte", "$lte", "$in", "$nin", "$regex"}
	for _, pattern := range dangerousPatterns {
		s = strings.ReplaceAll(s, pattern, "")
	}

	return s
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}
