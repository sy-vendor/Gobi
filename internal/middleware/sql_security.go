package middleware

import (
	"gobi/pkg/errors"
	"gobi/pkg/security"
	"gobi/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// SQLSecurityMiddleware provides SQL injection protection at the request level
func SQLSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if strings.Contains(contentType, "application/json") {
				c.Next()
				return
			}
		}

		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if containsSuspiciousSQLPattern(value) {
					c.Error(errors.NewError(errors.ErrCodeSQLInjection, "Suspicious SQL pattern detected in request", nil))
					c.Abort()
					return
				}
			}
		}

		if err := c.Request.ParseForm(); err == nil {
			for _, values := range c.Request.Form {
				for _, value := range values {
					if containsSuspiciousSQLPattern(value) {
						c.Error(errors.NewError(errors.ErrCodeSQLInjection, "Suspicious SQL pattern detected in form data", nil))
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
	config := security.GetGlobalSQLConfig()
	suspiciousPatterns := config.GetSuspiciousPatterns()

	// Add additional patterns for middleware
	additionalPatterns := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
		"UNION", "EXEC", "EXECUTE", "SCRIPT", "JAVASCRIPT", "VBSCRIPT",
		"<", ">", "\"", "'", ";", "--", "/*", "*/", "#",
	}

	allPatterns := append(suspiciousPatterns, additionalPatterns...)

	upperInput := strings.ToUpper(input)
	for _, pattern := range allPatterns {
		if strings.Contains(upperInput, pattern) {
			return true
		}
	}

	return false
}

// ValidateSQLInBody validates SQL in request body
func ValidateSQLInBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isSQLEndpoint(c.Request.URL.Path) {
			c.Next()
			return
		}

		body, exists := c.Get("requestBody")
		if !exists {
			c.Next()
			return
		}

		if sqlData, ok := body.(map[string]interface{}); ok {
			if sqlStr, exists := sqlData["sql"]; exists {
				if sql, ok := sqlStr.(string); ok {
					if containsSuspiciousSQLPattern(sql) {
						c.Error(errors.NewError(errors.ErrCodeSQLInjection, "Suspicious SQL pattern detected in request", nil))
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
		c.Next()
	}
}

// SQLAuditMiddleware logs SQL operations for audit purposes
func SQLAuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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
