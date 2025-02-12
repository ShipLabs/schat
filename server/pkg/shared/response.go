package shared

import (
	"github.com/gin-gonic/gin"
)

func ErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{"success": false, "error": message})
}

func SuccessResponse(ctx *gin.Context, statusCode int, message string, data any) {
	ctx.JSON(statusCode, gin.H{"success": true, "message": message, "data": data})
}
