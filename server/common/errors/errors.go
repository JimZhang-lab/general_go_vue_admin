/*
 * @Author: JimZhang
 * @Date: 2025-07-27
 * @Description: 增强的错误处理系统，支持错误链式传递和上下文信息
 */
package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// ErrorCode 错误码类型
type ErrorCode int

// 预定义错误码
const (
	// 系统级错误 (1000-1999)
	ErrSystem ErrorCode = 1000 + iota
	ErrDatabase
	ErrRedis
	ErrConfig
	ErrInternal
	ErrTimeout
	ErrRateLimit
	ErrCircuitBreaker

	// 业务级错误 (2000-2999)
	ErrBusiness ErrorCode = 2000 + iota
	ErrValidation
	ErrAuthentication
	ErrAuthorization
	ErrNotFound
	ErrAlreadyExists
	ErrInvalidOperation
	ErrResourceLocked

	// 网络级错误 (3000-3999)
	ErrNetwork ErrorCode = 3000 + iota
	ErrHTTPRequest
	ErrHTTPResponse
	ErrConnection
	ErrProtocol
)

// ErrorLevel 错误级别
type ErrorLevel int

const (
	LevelDebug ErrorLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// AppError 应用错误结构
type AppError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Level     ErrorLevel             `json:"level"`
	Timestamp time.Time              `json:"timestamp"`
	Stack     []StackFrame           `json:"stack,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Cause     error                  `json:"cause,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
}

// StackFrame 堆栈帧信息
type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持 Go 1.13+ 的错误链
func (e *AppError) Unwrap() error {
	return e.Cause
}

// Is 支持 errors.Is
func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Code == t.Code
	}
	return false
}

// WithContext 添加上下文信息
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithTraceID 添加追踪ID
func (e *AppError) WithTraceID(traceID string) *AppError {
	e.TraceID = traceID
	return e
}

// GetContext 获取上下文信息
func (e *AppError) GetContext(key string) (interface{}, bool) {
	if e.Context == nil {
		return nil, false
	}
	value, exists := e.Context[key]
	return value, exists
}

// JSON 序列化为JSON
func (e *AppError) JSON() string {
	data, _ := json.Marshal(e)
	return string(data)
}

// New 创建新的应用错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Level:     LevelError,
		Timestamp: time.Now(),
		Stack:     captureStack(2),
		Context:   make(map[string]interface{}),
	}
}

// Wrap 包装现有错误
func Wrap(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}

	return &AppError{
		Code:      code,
		Message:   message,
		Level:     LevelError,
		Timestamp: time.Now(),
		Stack:     captureStack(2),
		Context:   make(map[string]interface{}),
		Cause:     err,
	}
}

// Wrapf 格式化包装错误
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *AppError {
	return Wrap(err, code, fmt.Sprintf(format, args...))
}

// captureStack 捕获调用栈
func captureStack(skip int) []StackFrame {
	var frames []StackFrame
	for i := skip; i < skip+10; i++ { // 最多捕获10层
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		
		// 简化文件路径
		if idx := strings.LastIndex(file, "/"); idx >= 0 {
			file = file[idx+1:]
		}
		
		frames = append(frames, StackFrame{
			Function: fn.Name(),
			File:     file,
			Line:     line,
		})
	}
	return frames
}

// 预定义错误创建函数
func SystemError(message string) *AppError {
	return New(ErrSystem, message)
}

func DatabaseError(err error) *AppError {
	return Wrap(err, ErrDatabase, "数据库操作失败")
}

func RedisError(err error) *AppError {
	return Wrap(err, ErrRedis, "Redis操作失败")
}

func ValidationError(message string) *AppError {
	return New(ErrValidation, message)
}

func AuthenticationError(message string) *AppError {
	return New(ErrAuthentication, message)
}

func AuthorizationError(message string) *AppError {
	return New(ErrAuthorization, message)
}

func NotFoundError(resource string) *AppError {
	return New(ErrNotFound, fmt.Sprintf("%s 不存在", resource))
}

func AlreadyExistsError(resource string) *AppError {
	return New(ErrAlreadyExists, fmt.Sprintf("%s 已存在", resource))
}

func TimeoutError(operation string) *AppError {
	return New(ErrTimeout, fmt.Sprintf("%s 操作超时", operation))
}

func RateLimitError() *AppError {
	return New(ErrRateLimit, "请求频率过高，请稍后重试")
}

// ErrorChain 错误链，用于收集多个错误
type ErrorChain struct {
	errors []*AppError
}

// NewErrorChain 创建错误链
func NewErrorChain() *ErrorChain {
	return &ErrorChain{
		errors: make([]*AppError, 0),
	}
}

// Add 添加错误到链中
func (ec *ErrorChain) Add(err *AppError) *ErrorChain {
	if err != nil {
		ec.errors = append(ec.errors, err)
	}
	return ec
}

// HasErrors 检查是否有错误
func (ec *ErrorChain) HasErrors() bool {
	return len(ec.errors) > 0
}

// Errors 获取所有错误
func (ec *ErrorChain) Errors() []*AppError {
	return ec.errors
}

// First 获取第一个错误
func (ec *ErrorChain) First() *AppError {
	if len(ec.errors) > 0 {
		return ec.errors[0]
	}
	return nil
}

// Error 实现 error 接口
func (ec *ErrorChain) Error() string {
	if len(ec.errors) == 0 {
		return ""
	}
	
	var messages []string
	for _, err := range ec.errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}
