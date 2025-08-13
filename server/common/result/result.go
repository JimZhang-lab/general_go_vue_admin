package result

import (
	"net/http"
	"server/common/errors"
	"time"

	"github.com/gin-gonic/gin"
)

// 通用返回结构
type Result struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
	TraceID   string      `json:"trace_id,omitempty"`
	Success   bool        `json:"success"`
}

// 分页结果结构
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Pages    int         `json:"pages"`
}

// 返回成功
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}

	traceID := getTraceID(c)

	res := Result{
		Code:      int(ApiCode.SUCCESS),
		Message:   ApiCode.GetMessage(ApiCode.SUCCESS),
		Data:      data,
		Timestamp: time.Now().Unix(),
		TraceID:   traceID,
		Success:   true,
	}
	c.JSON(http.StatusOK, res)
}

// 返回分页成功
func SuccessWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	pageResult := PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}

	Success(c, pageResult)
}

// 返回失败
func Failed(c *gin.Context, code int, message string) {
	traceID := getTraceID(c)

	res := Result{
		Code:      code,
		Message:   message,
		Data:      gin.H{},
		Timestamp: time.Now().Unix(),
		TraceID:   traceID,
		Success:   false,
	}
	c.JSON(http.StatusOK, res)
}

// 返回应用错误
func FailedWithError(c *gin.Context, err *errors.AppError) {
	if err == nil {
		Success(c, nil)
		return
	}

	traceID := getTraceID(c)
	if err.TraceID == "" {
		err.TraceID = traceID
	}

	res := Result{
		Code:      int(err.Code),
		Message:   err.Message,
		Data:      gin.H{},
		Timestamp: time.Now().Unix(),
		TraceID:   err.TraceID,
		Success:   false,
	}

	// 根据错误级别设置HTTP状态码
	httpStatus := getHTTPStatusFromError(err)
	c.JSON(httpStatus, res)
}

// 获取追踪ID
func getTraceID(c *gin.Context) string {
	if traceID := c.GetHeader("X-Trace-ID"); traceID != "" {
		return traceID
	}
	if traceID := c.GetString("trace_id"); traceID != "" {
		return traceID
	}
	return ""
}

// 根据错误类型获取HTTP状态码
func getHTTPStatusFromError(err *errors.AppError) int {
	switch err.Code {
	case errors.ErrAuthentication:
		return http.StatusUnauthorized
	case errors.ErrAuthorization:
		return http.StatusForbidden
	case errors.ErrNotFound:
		return http.StatusNotFound
	case errors.ErrValidation:
		return http.StatusBadRequest
	case errors.ErrRateLimit:
		return http.StatusTooManyRequests
	case errors.ErrTimeout:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}
