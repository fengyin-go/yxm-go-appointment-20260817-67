// Package model 定义预约挂号系统的领域模型与校验逻辑。
package model

import "errors"

// ValidationError 表示字段校验失败。
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}

// Is 使 ValidationError 匹配 ErrValidation 哨兵，
// 这样即使被 fmt.Errorf("%w", err) 包裹，errors.Is(err, ErrValidation) 仍可识别。
func (e *ValidationError) Is(target error) bool {
	return target == ErrValidation
}

// NewValidationError 构造字段校验错误。
func NewValidationError(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}

var ErrValidation = errors.New("validation error")

// IsValidationError 判断错误是否为字段校验错误。
func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidation)
}
