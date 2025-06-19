package middleware

import (
	"gobi/pkg/errors"
	"gobi/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			path := c.Request.URL.Path
			method := c.Request.Method
			userID, _ := c.Get("userID")
			role, _ := c.Get("role")
			if customErr, ok := err.(*errors.CustomError); ok {
				utils.Logger.WithFields(map[string]interface{}{
					"path":    path,
					"method":  method,
					"userID":  userID,
					"role":    role,
					"code":    customErr.Code,
					"message": customErr.Message,
					"error":   customErr.Error(),
				}).Error("API error")
				c.JSON(customErr.Code, errors.ErrorResponse{
					Code:    customErr.Code,
					Message: customErr.Message,
					Error:   customErr.Error(),
				})
			} else {
				utils.Logger.WithFields(map[string]interface{}{
					"path":   path,
					"method": method,
					"userID": userID,
					"role":   role,
					"error":  err.Error(),
				}).Error("Internal server error")
				c.JSON(errors.ErrInternalServer.Code, errors.ErrorResponse{
					Code:    errors.ErrInternalServer.Code,
					Message: errors.ErrInternalServer.Message,
					Error:   err.Error(),
				})
			}
			c.Abort()
		}
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				path := c.Request.URL.Path
				method := c.Request.Method
				userID, _ := c.Get("userID")
				role, _ := c.Get("role")
				utils.Logger.WithFields(map[string]interface{}{
					"path":   path,
					"method": method,
					"userID": userID,
					"role":   role,
					"panic":  err,
				}).Error("Panic recovered")
				c.JSON(errors.ErrInternalServer.Code, errors.ErrorResponse{
					Code:    errors.ErrInternalServer.Code,
					Message: errors.ErrInternalServer.Message,
					Error:   "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
