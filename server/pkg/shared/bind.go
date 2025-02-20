package shared

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseBody(ctx *gin.Context, body any) bool {
	if err := ctx.ShouldBindBodyWithJSON(body); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}
