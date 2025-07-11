package middleware

import (
	"gobi/config"
	"gobi/pkg/errors"
	"strings"

	"gobi/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config, userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errors.NewError(errors.ErrCodeUnauthorized, "Authorization header is required", nil))
			c.Abort()
			return
		}

		if strings.HasPrefix(authHeader, "ApiKey ") {
			plainKey := strings.TrimPrefix(authHeader, "ApiKey ")
			if len(plainKey) < 12 {
				c.Error(errors.NewError(errors.ErrCodeInvalidAPIKey, "Invalid API key format", nil))
				c.Abort()
				return
			}
			prefix := plainKey[:12]
			apiKey, err := userService.GetAPIKeyByPrefix(prefix)
			if err != nil || !userService.ValidateAPIKey(apiKey, plainKey) {
				c.Error(errors.NewError(errors.ErrCodeInvalidAPIKey, "Invalid or expired API key", nil))
				c.Abort()
				return
			}
			c.Set("userID", apiKey.UserID)
			c.Set("role", "service")
			c.Next()
			return
		}

		// JWT authentication (default)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(errors.NewError(errors.ErrCodeUnauthorized, "Invalid authorization header format", nil))
			c.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			c.Error(errors.NewError(errors.ErrCodeInvalidToken, "Invalid token", err))
			c.Abort()
			return
		}

		if !token.Valid {
			c.Error(errors.ErrInvalidToken)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Error(errors.NewError(errors.ErrCodeInvalidToken, "Invalid token claims", nil))
			c.Abort()
			return
		}

		// 检查必要的claims是否存在
		userID, exists := claims["user_id"]
		if !exists {
			c.Error(errors.ErrTokenMissingClaims)
			c.Abort()
			return
		}

		role, exists := claims["role"]
		if !exists {
			c.Error(errors.ErrTokenMissingClaims)
			c.Abort()
			return
		}

		// 类型断言
		userIDFloat, ok := userID.(float64)
		if !ok {
			c.Error(errors.NewError(errors.ErrCodeInvalidToken, "Invalid user_id type in token", nil))
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.Error(errors.NewError(errors.ErrCodeInvalidToken, "Invalid role type in token", nil))
			c.Abort()
			return
		}

		c.Set("userID", uint(userIDFloat))
		c.Set("role", roleStr)
		c.Next()
	}
}
