package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors 格式化验证错误
func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   strings.ToLower(e.Field()),
				Message: getErrorMessage(e),
			})
		}
	}

	return errors
}

// getErrorMessage 获取错误消息
func getErrorMessage(e validator.FieldError) string {
	field := e.Field()

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s 是必填项", field)
	case "email":
		return fmt.Sprintf("%s 必须是有效的邮箱地址", field)
	case "min":
		return fmt.Sprintf("%s 最小长度为 %s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s 最大长度为 %s", field, e.Param())
	case "len":
		return fmt.Sprintf("%s 长度必须为 %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s 必须大于 %s", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s 必须大于等于 %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s 必须小于 %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s 必须小于等于 %s", field, e.Param())
	case "alphanum":
		return fmt.Sprintf("%s 只能包含字母和数字", field)
	case "alpha":
		return fmt.Sprintf("%s 只能包含字母", field)
	case "numeric":
		return fmt.Sprintf("%s 只能包含数字", field)
	case "url":
		return fmt.Sprintf("%s 必须是有效的 URL", field)
	case "uri":
		return fmt.Sprintf("%s 必须是有效的 URI", field)
	case "oneof":
		return fmt.Sprintf("%s 必须是以下值之一: %s", field, e.Param())
	default:
		return fmt.Sprintf("%s 验证失败", field)
	}
}

// GetValidationErrorMessage 获取验证错误的第一条消息
func GetValidationErrorMessage(err error) string {
	errors := FormatValidationErrors(err)
	if len(errors) > 0 {
		return errors[0].Message
	}
	return "参数验证失败"
}
