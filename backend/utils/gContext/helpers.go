package gContext

import (
	"github.com/Alan-Luc/VertiLog/backend/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleReqError(ctx *gin.Context, err error, statusCode int) bool {
	if err != nil {
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		clientIP := ctx.ClientIP()
		userAgent := ctx.Request.UserAgent()
		queryParams := ctx.Request.URL.RawQuery

		// Log the error with additional request context
		logger.Logger.Info(
			"An error occured:",
			zap.Error(err),
			zap.Int("status_code", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.String("query_params", queryParams),
		)

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return true
	}
	return false
}
