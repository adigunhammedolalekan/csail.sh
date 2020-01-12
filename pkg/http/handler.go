package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SuccessResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

func BadRequestResponse(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, &ErrorResponse{Error: true, Message: message})
}

func ForbiddenRequestResponse(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusForbidden, &ErrorResponse{Error: true, Message: message})
}

func InternalServerErrorResponse(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: true, Message: message})
}

