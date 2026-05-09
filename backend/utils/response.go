package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ctx = request context chứa request và response
// data interface: nhận bất kì kiểu dữ liệu nào

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func Create(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

func BadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   message,
	})
}

func NotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error":   message,
	})
}

func Unauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error":   message,
	})
}

func Forbidden(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusForbidden, gin.H{
		"success": false,
		"error":   message,
	})
}

func InternalError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   message,
	})
}
