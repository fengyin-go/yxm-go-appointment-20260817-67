package model

import (
	"strings"
	"time"
)

// Doctor 医生。
type Doctor struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Title        string    `json:"title"`        // 职称，如 主任医师
	DepartmentID string    `json:"department_id"` // 所属科室
	Specialty    string    `json:"specialty"`     // 擅长领域
	CreatedAt    time.Time `json:"created_at"`
}

// Validate 规范化并校验医生字段。
func (d *Doctor) Validate() error {
	d.Name = strings.TrimSpace(d.Name)
	d.Title = strings.TrimSpace(d.Title)
	d.DepartmentID = strings.TrimSpace(d.DepartmentID)
	d.Specialty = strings.TrimSpace(d.Specialty)
	if d.Name == "" {
		return NewValidationError("name", "医生姓名不能为空")
	}
	if d.DepartmentID == "" {
		return NewValidationError("department_id", "所属科室不能为空")
	}
	return nil
}
