package model

import (
	"strings"
	"time"
)

// Department 科室。
type Department struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Validate 规范化并校验科室字段。
func (d *Department) Validate() error {
	d.Name = strings.TrimSpace(d.Name)
	d.Description = strings.TrimSpace(d.Description)
	if d.Name == "" {
		return NewValidationError("name", "科室名称不能为空")
	}
	return nil
}
