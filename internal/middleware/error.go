package middleware

import (
	"gobi/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if customErr, ok := err.(*errors.CustomError); ok {
				c.JSON(customErr.Code, errors.ErrorResponse{
					Code:    customErr.Code,
					Message: customErr.Message,
					Error:   customErr.Error(),
				})
			} else {
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
