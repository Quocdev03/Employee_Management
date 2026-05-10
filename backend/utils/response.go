package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success: Trả về thành công 200 kèm data
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// SuccessWithMeta: Trả về thành công kèm data và thông tin bổ sung (phân trang, ...)
func SuccessWithMeta(ctx *gin.Context, data interface{}, meta gin.H) {
	res := gin.H{
		"success": true,
		"data":    data,
	}
	for k, v := range meta {
		res[k] = v
	}
	ctx.JSON(http.StatusOK, res)
}

// Create: Trả về thành công 201 khi tạo mới dữ liệu
func Create(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

// BadRequest: Lỗi dữ liệu đầu vào 400
func BadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   message,
	})
}

// Unauthorized: Lỗi xác thực 401
func Unauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error":   message,
	})
}

// Forbidden: Lỗi phân quyền 403
func Forbidden(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusForbidden, gin.H{
		"success": false,
		"error":   message,
	})
}

// NotFound: Lỗi không tìm thấy tài nguyên 404
func NotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error":   message,
	})
}

// Conflict: Lỗi xung đột dữ liệu 409
func Conflict(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusConflict, gin.H{
		"success": false,
		"error":   message,
	})
}

// InternalError: Lỗi hệ thống 500
func InternalError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   message,
	})
}

// Error: Hàm tổng quát để trả về bất kỳ mã lỗi nào
func Error(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{
		"success": false,
		"error":   message,
	})
}
