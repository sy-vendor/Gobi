package middleware

import (
	"gobi/pkg/errors"
	"gobi/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// SQLSecurityMiddleware provides SQL injection protection at the request level
func SQLSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for SQL in request body
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if strings.Contains(contentType, "application/json") {
				// For JSON requests, we'll validate SQL in the handler
				// This middleware serves as a first line of defense
				c.Next()
				return
			}
		}

		// Check URL parameters for suspicious SQL patterns
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if containsSuspiciousSQLPattern(value) {
					c.Error(errors.NewError(400, "Suspicious SQL pattern detected in request", nil))
					c.Abort()
					return
				}
			}
		}

		// Check form data for suspicious SQL patterns
		if err := c.Request.ParseForm(); err == nil {
			for _, values := range c.Request.Form {
				for _, value := range values {
					if containsSuspiciousSQLPattern(value) {
						c.Error(errors.NewError(400, "Suspicious SQL pattern detected in form data", nil))
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

// containsSuspiciousSQLPattern checks if a string contains suspicious SQL patterns
func containsSuspiciousSQLPattern(input string) bool {
	suspiciousPatterns := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
		"UNION", "EXEC", "EXECUTE", "SCRIPT", "JAVASCRIPT", "VBSCRIPT",
		"<", ">", "\"", "'", ";", "--", "/*", "*/", "#",
		"1=1", "TRUE", "FALSE", "OR 1", "AND 1",
		"INFORMATION_SCHEMA", "SYSTEM_TABLES", "DUAL",
	}

	upperInput := strings.ToUpper(input)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(upperInput, pattern) {
			return true
		}
	}

	return false
}

// ValidateSQLInBody validates SQL in request body
func ValidateSQLInBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to specific endpoints that handle SQL
		if !isSQLEndpoint(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Get the request body
		body, exists := c.Get("requestBody")
		if !exists {
			c.Next()
			return
		}

		// Check if body contains SQL field
		if sqlData, ok := body.(map[string]interface{}); ok {
			if sqlStr, exists := sqlData["sql"]; exists {
				if sql, ok := sqlStr.(string); ok {
					// Validate SQL
					if err := utils.ValidateSQL(sql); err != nil {
						c.Error(errors.WrapError(err, "SQL validation failed"))
						c.Abort()
						return
					}

					// Ensure query is read-only
					if !utils.IsReadOnlyQuery(sql) {
						c.Error(errors.NewError(403, "Only SELECT queries are allowed", nil))
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

// isSQLEndpoint checks if the endpoint handles SQL
func isSQLEndpoint(path string) bool {
	sqlEndpoints := []string{
		"/api/queries",
		"/api/queries/",
		"/api/charts",
		"/api/charts/",
		"/api/reports",
		"/api/reports/",
	}

	for _, endpoint := range sqlEndpoints {
		if strings.HasPrefix(path, endpoint) {
			return true
		}
	}

	return false
}

// SQLRateLimitMiddleware provides rate limiting for SQL operations
func SQLRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This is a placeholder for rate limiting implementation
		// In production, you would use a proper rate limiter like Redis

		// For now, we'll just pass through
		c.Next()
	}
}

// SQLAuditMiddleware logs SQL operations for audit purposes
func SQLAuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request for audit purposes
		userID, _ := c.Get("userID")
		role, _ := c.Get("role")

		utils.Logger.WithFields(map[string]interface{}{
			"path":      c.Request.URL.Path,
			"method":    c.Request.Method,
			"userID":    userID,
			"role":      role,
			"ip":        c.ClientIP(),
			"userAgent": c.Request.UserAgent(),
		}).Info("SQL operation requested")

		c.Next()
	}
}
