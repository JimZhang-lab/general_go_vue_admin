/*
 * @Author: JimZhang
 * @Date: 2025-07-27
 * @Description: 增强的控制器基类，支持多种数据格式绑定和错误处理
 */
package controller

import (
	"server/common/errors"
	"server/common/result"
	"server/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BaseController 基础控制器
type BaseController struct{}

// BindRequest 智能绑定请求数据
func (bc *BaseController) BindRequest(c *gin.Context, obj interface{}) *errors.AppError {
	binder := middleware.GetSmartBinder(c)
	return binder.BindRequest(c, obj)
}

// BindJSON 绑定JSON数据
func (bc *BaseController) BindJSON(c *gin.Context, obj interface{}) *errors.AppError {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.Wrap(err, errors.ErrValidation, "JSON数据绑定失败")
	}
	return nil
}

// BindQuery 绑定查询参数
func (bc *BaseController) BindQuery(c *gin.Context, obj interface{}) *errors.AppError {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.Wrap(err, errors.ErrValidation, "查询参数绑定失败")
	}
	return nil
}

// BindForm 绑定表单数据
func (bc *BaseController) BindForm(c *gin.Context, obj interface{}) *errors.AppError {
	if err := c.ShouldBind(obj); err != nil {
		return errors.Wrap(err, errors.ErrValidation, "表单数据绑定失败")
	}
	return nil
}

// BindURI 绑定URI参数
func (bc *BaseController) BindURI(c *gin.Context, obj interface{}) *errors.AppError {
	if err := c.ShouldBindUri(obj); err != nil {
		return errors.Wrap(err, errors.ErrValidation, "URI参数绑定失败")
	}
	return nil
}

// GetIntParam 获取整数参数
func (bc *BaseController) GetIntParam(c *gin.Context, key string) (int, *errors.AppError) {
	value := c.Param(key)
	if value == "" {
		return 0, errors.ValidationError("参数 " + key + " 不能为空")
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.Wrap(err, errors.ErrValidation, "参数 "+key+" 必须是整数")
	}
	
	return intValue, nil
}

// GetIntQuery 获取整数查询参数
func (bc *BaseController) GetIntQuery(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return intValue
}

// GetStringParam 获取字符串参数
func (bc *BaseController) GetStringParam(c *gin.Context, key string) (string, *errors.AppError) {
	value := c.Param(key)
	if value == "" {
		return "", errors.ValidationError("参数 " + key + " 不能为空")
	}
	return value, nil
}

// GetStringQuery 获取字符串查询参数
func (bc *BaseController) GetStringQuery(c *gin.Context, key string, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetPageParams 获取分页参数
func (bc *BaseController) GetPageParams(c *gin.Context) (page, pageSize int) {
	page = bc.GetIntQuery(c, "page", 1)
	pageSize = bc.GetIntQuery(c, "page_size", 10)
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}
	
	return page, pageSize
}

// Success 返回成功响应
func (bc *BaseController) Success(c *gin.Context, data interface{}) {
	result.Success(c, data)
}

// SuccessWithPage 返回分页成功响应
func (bc *BaseController) SuccessWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	result.SuccessWithPage(c, list, total, page, pageSize)
}

// Failed 返回失败响应
func (bc *BaseController) Failed(c *gin.Context, code int, message string) {
	result.Failed(c, code, message)
}

// FailedWithError 返回错误响应
func (bc *BaseController) FailedWithError(c *gin.Context, err *errors.AppError) {
	result.FailedWithError(c, err)
}

// HandleError 统一错误处理
func (bc *BaseController) HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	
	if appErr, ok := err.(*errors.AppError); ok {
		bc.FailedWithError(c, appErr)
	} else {
		// 包装普通错误
		appErr := errors.Wrap(err, errors.ErrInternal, "内部服务器错误")
		bc.FailedWithError(c, appErr)
	}
}

// ValidateRequired 验证必需字段
func (bc *BaseController) ValidateRequired(fields map[string]interface{}) *errors.AppError {
	for fieldName, fieldValue := range fields {
		if fieldValue == nil {
			return errors.ValidationError("字段 " + fieldName + " 是必需的")
		}
		
		switch v := fieldValue.(type) {
		case string:
			if v == "" {
				return errors.ValidationError("字段 " + fieldName + " 不能为空")
			}
		case int, int32, int64:
			if v == 0 {
				return errors.ValidationError("字段 " + fieldName + " 不能为0")
			}
		}
	}
	return nil
}

// GetUserID 从上下文获取用户ID
func (bc *BaseController) GetUserID(c *gin.Context) (int, *errors.AppError) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.AuthenticationError("用户未登录")
	}
	
	if id, ok := userID.(int); ok {
		return id, nil
	}
	
	return 0, errors.AuthenticationError("无效的用户ID")
}

// GetUserInfo 从上下文获取用户信息
func (bc *BaseController) GetUserInfo(c *gin.Context) (map[string]interface{}, *errors.AppError) {
	userInfo, exists := c.Get("user_info")
	if !exists {
		return nil, errors.AuthenticationError("用户信息不存在")
	}
	
	if info, ok := userInfo.(map[string]interface{}); ok {
		return info, nil
	}
	
	return nil, errors.AuthenticationError("无效的用户信息")
}

// SetTraceID 设置追踪ID
func (bc *BaseController) SetTraceID(c *gin.Context, traceID string) {
	c.Set("trace_id", traceID)
	c.Header("X-Trace-ID", traceID)
}

// GetTraceID 获取追踪ID
func (bc *BaseController) GetTraceID(c *gin.Context) string {
	if traceID := c.GetString("trace_id"); traceID != "" {
		return traceID
	}
	return c.GetHeader("X-Trace-ID")
}

// LogRequest 记录请求日志
func (bc *BaseController) LogRequest(c *gin.Context, action string, params interface{}) {
	// 这里可以集成日志系统
	// log.Info("Request", "action", action, "params", params, "trace_id", bc.GetTraceID(c))
}

// LogResponse 记录响应日志
func (bc *BaseController) LogResponse(c *gin.Context, action string, result interface{}) {
	// 这里可以集成日志系统
	// log.Info("Response", "action", action, "result", result, "trace_id", bc.GetTraceID(c))
}

// CacheKey 生成缓存键
func (bc *BaseController) CacheKey(prefix string, params ...interface{}) string {
	key := prefix
	for _, param := range params {
		key += "_" + toString(param)
	}
	return key
}

// toString 将任意类型转换为字符串
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	default:
		return ""
	}
}
