/*
 * @Author: JimZhang
 * @Date: 2025-07-27
 * @Description: 增强的请求绑定中间件，支持多种数据格式
 */
package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"server/common/errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ContentType 常量
const (
	ContentTypeJSON          = "application/json"
	ContentTypeForm          = "application/x-www-form-urlencoded"
	ContentTypeMultipartForm = "multipart/form-data"
	ContentTypeXML           = "application/xml"
	ContentTypeYAML          = "application/yaml"
	ContentTypePlain         = "text/plain"
)

// BindingConfig 绑定配置
type BindingConfig struct {
	MaxMemory        int64                     // 最大内存限制 (multipart form)
	AllowedTypes     []string                  // 允许的内容类型
	RequiredHeaders  map[string]string         // 必需的请求头
	CustomValidators map[string]validator.Func // 自定义验证器
	EnableValidation bool                      // 是否启用验证
	ValidationMode   string                    // 验证模式: strict, loose
}

// DefaultBindingConfig 默认配置
var DefaultBindingConfig = &BindingConfig{
	MaxMemory: 32 << 20, // 32MB
	AllowedTypes: []string{
		ContentTypeJSON,
		ContentTypeForm,
		ContentTypeMultipartForm,
		ContentTypeXML,
		ContentTypeYAML,
	},
	EnableValidation: true,
	ValidationMode:   "strict",
}

// SmartBinder 智能绑定器
type SmartBinder struct {
	config    *BindingConfig
	validator *validator.Validate
}

// NewSmartBinder 创建智能绑定器
func NewSmartBinder(config *BindingConfig) *SmartBinder {
	if config == nil {
		config = DefaultBindingConfig
	}

	v := validator.New()

	// 注册自定义验证器
	if config.CustomValidators != nil {
		for tag, fn := range config.CustomValidators {
			v.RegisterValidation(tag, fn)
		}
	}

	return &SmartBinder{
		config:    config,
		validator: v,
	}
}

// BindRequest 智能绑定请求数据
func (sb *SmartBinder) BindRequest(c *gin.Context, obj interface{}) *errors.AppError {
	// 检查必需的请求头
	if err := sb.checkRequiredHeaders(c); err != nil {
		return err
	}

	contentType := c.GetHeader("Content-Type")
	if contentType == "" {
		contentType = ContentTypeJSON // 默认为JSON
	}

	// 移除charset等参数
	if idx := strings.Index(contentType, ";"); idx > 0 {
		contentType = strings.TrimSpace(contentType[:idx])
	}

	// 检查内容类型是否被允许
	if !sb.isAllowedContentType(contentType) {
		return errors.ValidationError(fmt.Sprintf("不支持的内容类型: %s", contentType))
	}

	// 根据内容类型进行绑定
	var err error
	switch contentType {
	case ContentTypeJSON:
		err = sb.bindJSON(c, obj)
	case ContentTypeForm:
		err = sb.bindForm(c, obj)
	case ContentTypeMultipartForm:
		err = sb.bindMultipartForm(c, obj)
	case ContentTypeXML:
		err = sb.bindXML(c, obj)
	case ContentTypeYAML:
		err = sb.bindYAML(c, obj)
	default:
		// 尝试自动检测
		err = sb.bindAuto(c, obj)
	}

	if err != nil {
		return errors.Wrap(err, errors.ErrValidation, "请求数据绑定失败")
	}

	// 验证数据
	if sb.config.EnableValidation {
		if err := sb.validateStruct(obj); err != nil {
			return err
		}
	}

	return nil
}

// bindJSON 绑定JSON数据
func (sb *SmartBinder) bindJSON(c *gin.Context, obj interface{}) error {
	if c.Request.Body == nil {
		return fmt.Errorf("请求体为空")
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("读取请求体失败: %v", err)
	}

	if len(body) == 0 {
		return fmt.Errorf("请求体为空")
	}

	if err := json.Unmarshal(body, obj); err != nil {
		return fmt.Errorf("JSON解析失败: %v", err)
	}

	return nil
}

// bindForm 绑定表单数据
func (sb *SmartBinder) bindForm(c *gin.Context, obj interface{}) error {
	if err := c.Request.ParseForm(); err != nil {
		return fmt.Errorf("表单解析失败: %v", err)
	}

	return sb.mapFormToStruct(c.Request.Form, obj)
}

// bindMultipartForm 绑定多部分表单数据
func (sb *SmartBinder) bindMultipartForm(c *gin.Context, obj interface{}) error {
	if err := c.Request.ParseMultipartForm(sb.config.MaxMemory); err != nil {
		return fmt.Errorf("多部分表单解析失败: %v", err)
	}

	// 绑定表单字段
	if err := sb.mapFormToStruct(c.Request.MultipartForm.Value, obj); err != nil {
		return err
	}

	// 处理文件上传
	return sb.handleFileUploads(c.Request.MultipartForm.File, obj)
}

// bindXML 绑定XML数据
func (sb *SmartBinder) bindXML(c *gin.Context, obj interface{}) error {
	return c.ShouldBindXML(obj)
}

// bindYAML 绑定YAML数据
func (sb *SmartBinder) bindYAML(c *gin.Context, obj interface{}) error {
	return c.ShouldBindYAML(obj)
}

// bindAuto 自动检测并绑定
func (sb *SmartBinder) bindAuto(c *gin.Context, obj interface{}) error {
	// 首先尝试从URL参数绑定
	if err := c.ShouldBindQuery(obj); err == nil {
		return nil
	}

	// 然后尝试JSON
	if err := sb.bindJSON(c, obj); err == nil {
		return nil
	}

	// 最后尝试表单
	return sb.bindForm(c, obj)
}

// mapFormToStruct 将表单数据映射到结构体
func (sb *SmartBinder) mapFormToStruct(form map[string][]string, obj interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("目标对象必须是结构体指针")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		// 获取字段标签
		tag := fieldType.Tag.Get("form")
		if tag == "" {
			tag = fieldType.Tag.Get("json")
		}
		if tag == "" {
			tag = strings.ToLower(fieldType.Name)
		}

		// 跳过忽略的字段
		if tag == "-" {
			continue
		}

		// 获取表单值
		values, exists := form[tag]
		if !exists || len(values) == 0 {
			continue
		}

		// 设置字段值
		if err := sb.setFieldValue(field, values[0]); err != nil {
			return fmt.Errorf("设置字段 %s 失败: %v", fieldType.Name, err)
		}
	}

	return nil
}

// setFieldValue 设置字段值
func (sb *SmartBinder) setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(intVal)
		} else {
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintVal, err := strconv.ParseUint(value, 10, 64); err == nil {
			field.SetUint(uintVal)
		} else {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			field.SetFloat(floatVal)
		} else {
			return err
		}
	case reflect.Bool:
		if boolVal, err := strconv.ParseBool(value); err == nil {
			field.SetBool(boolVal)
		} else {
			return err
		}
	default:
		return fmt.Errorf("不支持的字段类型: %s", field.Kind())
	}

	return nil
}

// handleFileUploads 处理文件上传
func (sb *SmartBinder) handleFileUploads(files map[string][]*multipart.FileHeader, obj interface{}) error {
	// 这里可以根据需要实现文件处理逻辑
	// 例如将文件信息设置到结构体的相应字段中
	return nil
}

// validateStruct 验证结构体
func (sb *SmartBinder) validateStruct(obj interface{}) *errors.AppError {
	if err := sb.validator.Struct(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, e := range validationErrors {
				errorMessages = append(errorMessages, sb.formatValidationError(e))
			}
			return errors.ValidationError(strings.Join(errorMessages, "; "))
		}
		return errors.ValidationError(err.Error())
	}
	return nil
}

// formatValidationError 格式化验证错误
func (sb *SmartBinder) formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("字段 %s 是必需的", err.Field())
	case "min":
		return fmt.Sprintf("字段 %s 的最小值为 %s", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("字段 %s 的最大值为 %s", err.Field(), err.Param())
	case "email":
		return fmt.Sprintf("字段 %s 必须是有效的邮箱地址", err.Field())
	default:
		return fmt.Sprintf("字段 %s 验证失败: %s", err.Field(), err.Tag())
	}
}

// checkRequiredHeaders 检查必需的请求头
func (sb *SmartBinder) checkRequiredHeaders(c *gin.Context) *errors.AppError {
	if sb.config.RequiredHeaders == nil {
		return nil
	}

	for header, expectedValue := range sb.config.RequiredHeaders {
		actualValue := c.GetHeader(header)
		if actualValue != expectedValue {
			return errors.ValidationError(fmt.Sprintf("缺少或错误的请求头: %s", header))
		}
	}

	return nil
}

// isAllowedContentType 检查内容类型是否被允许
func (sb *SmartBinder) isAllowedContentType(contentType string) bool {
	if sb.config.AllowedTypes == nil {
		return true
	}

	for _, allowed := range sb.config.AllowedTypes {
		if contentType == allowed {
			return true
		}
	}

	return false
}

// SmartBindMiddleware 智能绑定中间件
func SmartBindMiddleware(config *BindingConfig) gin.HandlerFunc {
	binder := NewSmartBinder(config)

	return func(c *gin.Context) {
		// 将绑定器存储到上下文中，供后续使用
		c.Set("smart_binder", binder)
		c.Next()
	}
}

// GetSmartBinder 从上下文获取智能绑定器
func GetSmartBinder(c *gin.Context) *SmartBinder {
	if binder, exists := c.Get("smart_binder"); exists {
		if sb, ok := binder.(*SmartBinder); ok {
			return sb
		}
	}
	return NewSmartBinder(nil)
}
