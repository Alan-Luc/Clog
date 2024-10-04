package gContext

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HandleReqError(ctx *gin.Context, err error, statusCode int) bool {
	if err != nil {
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		clientIP := ctx.ClientIP()
		userAgent := ctx.Request.UserAgent()
		queryParams := ctx.Request.URL.RawQuery

		// Log the error with additional request context
		log.Printf(
			"An error occurred: %v | Status Code: %d | Method: %s | Path: %s | Client IP: %s | User-Agent: %s | Query: %s",
			err,
			statusCode,
			method,
			path,
			clientIP,
			userAgent,
			queryParams,
		)

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return true
	}
	return false
}
