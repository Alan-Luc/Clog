package gContext

import (
	"github.com/gin-gonic/gin"
)

func HandleReqError(ctx *gin.Context, err error, statusCode int) bool {
	if err != nil {
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return true
	}
	return false
}
