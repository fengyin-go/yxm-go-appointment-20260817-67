package model

import (
	"strings"
	"time"
)

// Patient 患者。
type Patient struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	IDCard    string    `json:"id_card"`
	Gender    string    `json:"gender"` // male / female / other
	CreatedAt time.Time `json:"created_at"`
}

// Validate 规范化并校验患者字段。
func (p *Patient) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	p.Phone = strings.TrimSpace(p.Phone)
	p.IDCard = strings.TrimSpace(p.IDCard)
	p.Gender = strings.TrimSpace(p.Gender)
	if p.Name == "" {
		return NewValidationError("name", "患者姓名不能为空")
	}
	if p.Phone == "" {
		return NewValidationError("phone", "联系电话不能为空")
	}
	if p.Gender == "" {
		p.Gender = "other"
	}
	if p.Gender != "male" && p.Gender != "female" && p.Gender != "other" {
		return NewValidationError("gender", "性别不合法")
	}
	return nil
}
