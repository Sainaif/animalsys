package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/internal/core/entities"
)

// AuditLogger interface for logging audit trails
type AuditLogger interface {
	Log(auditLog *entities.AuditLog) error
}

// AuditConfig configures audit logging
type AuditConfig struct {
	Logger             AuditLogger
	SensitiveFields    []string // Fields to redact from logs
	EnableBodyLogging  bool     // Log request/response bodies
	SkipPaths         []string // Paths to skip audit logging
}

// Audit creates an audit trail middleware
func Audit(config AuditConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip audit logging for certain paths
		for _, skipPath := range config.SkipPaths {
			if c.Request.URL.Path == skipPath {
				c.Next()
				return
			}
		}

		// Only audit state-changing operations and sensitive reads
		method := c.Request.Method
		if !shouldAudit(method) {
			c.Next()
			return
		}

		// Get user info from context (set by auth middleware)
		userID, _ := c.Get("user_id")
		userEmail, _ := c.Get("user_email")
		userRole, _ := c.Get("user_role")

		// Determine action based on HTTP method
		action := getAuditAction(method)

		// Extract resource info from path
		resourceType, resourceID := extractResourceInfo(c.Request.URL.Path)

		// Create audit log
		auditLog := entities.NewAuditLog(
			toString(userID),
			toString(userEmail),
			toString(userRole),
			action,
			resourceType,
			resourceID,
			generateDescription(method, resourceType),
		)

		// Set IP and user agent
		auditLog.SetIPAndUserAgent(c.ClientIP(), c.Request.UserAgent())

		// Capture request body if enabled
		var requestBody map[string]interface{}
		if config.EnableBodyLogging && shouldCaptureBody(method) {
			requestBody = captureRequestBody(c)
			if requestBody != nil {
				redactSensitiveFields(requestBody, config.SensitiveFields)
				auditLog.SetBefore(requestBody)
			}
		}

		// Capture response
		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Capture response body if enabled and request was successful
		if config.EnableBodyLogging && responseWriter.Status() < 400 {
			var responseBody map[string]interface{}
			if err := json.Unmarshal(responseWriter.body.Bytes(), &responseBody); err == nil {
				redactSensitiveFields(responseBody, config.SensitiveFields)
				auditLog.SetAfter(responseBody)
			}
		}

		// Mark as failed if status code indicates error
		if responseWriter.Status() >= 400 {
			errorMsg := "Request failed"
			if len(c.Errors) > 0 {
				errorMsg = c.Errors.String()
			}
			auditLog.SetError(errorMsg)
		}

		// Log audit trail asynchronously
		go func() {
			if err := config.Logger.Log(auditLog); err != nil {
				// Log error but don't fail the request
				// In production, you might want to send this to a monitoring service
				println("Failed to log audit trail:", err.Error())
			}
		}()
	}
}

// shouldAudit determines if the request should be audited
func shouldAudit(method string) bool {
	// Audit all state-changing operations
	return method == http.MethodPost ||
		method == http.MethodPut ||
		method == http.MethodPatch ||
		method == http.MethodDelete
}

// getAuditAction maps HTTP method to audit action
func getAuditAction(method string) entities.AuditAction {
	switch method {
	case http.MethodPost:
		return entities.ActionCreate
	case http.MethodPut, http.MethodPatch:
		return entities.ActionUpdate
	case http.MethodDelete:
		return entities.ActionDelete
	default:
		return entities.ActionRead
	}
}

// extractResourceInfo extracts resource type and ID from path
func extractResourceInfo(path string) (string, string) {
	// Remove leading slash
	path = strings.TrimPrefix(path, "/")

	// Split path into segments
	segments := strings.Split(path, "/")

	if len(segments) < 2 {
		return segments[0], ""
	}

	// Assume format: /api/v1/resourceType/resourceID
	resourceType := segments[len(segments)-2]
	resourceID := ""

	// Check if last segment looks like an ID
	lastSegment := segments[len(segments)-1]
	if !strings.Contains(lastSegment, "?") && len(lastSegment) > 0 {
		// Could be an ID if it's not an action keyword
		if !isActionKeyword(lastSegment) {
			resourceID = lastSegment
		}
	}

	return resourceType, resourceID
}

// isActionKeyword checks if a segment is an action keyword
func isActionKeyword(segment string) bool {
	keywords := []string{"create", "update", "delete", "list", "search", "export", "import"}
	for _, keyword := range keywords {
		if segment == keyword {
			return true
		}
	}
	return false
}

// generateDescription generates a human-readable description
func generateDescription(method, resourceType string) string {
	action := ""
	switch method {
	case http.MethodPost:
		action = "created"
	case http.MethodPut, http.MethodPatch:
		action = "updated"
	case http.MethodDelete:
		action = "deleted"
	default:
		action = "accessed"
	}

	return strings.Title(action) + " " + resourceType
}

// shouldCaptureBody determines if request body should be captured
func shouldCaptureBody(method string) bool {
	return method == http.MethodPost ||
		method == http.MethodPut ||
		method == http.MethodPatch
}

// captureRequestBody captures and returns request body
func captureRequestBody(c *gin.Context) map[string]interface{} {
	if c.Request.Body == nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}

	// Restore body for further processing
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var body map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		return nil
	}

	return body
}

// redactSensitiveFields redacts sensitive fields from data
func redactSensitiveFields(data map[string]interface{}, sensitiveFields []string) {
	for key := range data {
		for _, sensitiveField := range sensitiveFields {
			if strings.EqualFold(key, sensitiveField) {
				data[key] = "***REDACTED***"
				break
			}
		}

		// Recursively redact nested objects
		if nested, ok := data[key].(map[string]interface{}); ok {
			redactSensitiveFields(nested, sensitiveFields)
		}
	}
}

// responseBodyWriter wraps ResponseWriter to capture response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// toString safely converts interface{} to string
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
