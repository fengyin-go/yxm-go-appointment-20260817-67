package model

import (
	"strings"
	"time"
)

// Schedule 医生的排班，某日某个时段的号源。
type Schedule struct {
	ID          string    `json:"id"`
	DoctorID    string    `json:"doctor_id"`
	Date        string    `json:"date"`       // 排班日期，格式 "2006-01-02"
	TimeSlot    string    `json:"time_slot"`  // 时段，如 "09:00-10:00"
	TotalSlots  int       `json:"total_slots"` // 总号源数
	BookedSlots int       `json:"booked_slots"` // 已预约数
	CreatedAt   time.Time `json:"created_at"`
}

// Validate 规范化并校验排班字段。
func (s *Schedule) Validate() error {
	s.DoctorID = strings.TrimSpace(s.DoctorID)
	s.Date = strings.TrimSpace(s.Date)
	s.TimeSlot = strings.TrimSpace(s.TimeSlot)
	if s.DoctorID == "" {
		return NewValidationError("doctor_id", "医生不能为空")
	}
	if s.Date == "" {
		return NewValidationError("date", "排班日期不能为空")
	}
	if _, err := time.Parse("2006-01-02", s.Date); err != nil {
		return NewValidationError("date", "日期格式需为 2006-01-02")
	}
	if s.TimeSlot == "" {
		return NewValidationError("time_slot", "时段不能为空")
	}
	if s.TotalSlots <= 0 {
		return NewValidationError("total_slots", "号源数必须为正数")
	}
	if s.BookedSlots < 0 {
		return NewValidationError("booked_slots", "已预约数不能为负数")
	}
	return nil
}

// Remaining 返回剩余号源数。
func (s *Schedule) Remaining() int {
	if s.TotalSlots < s.BookedSlots {
		return 0
	}
	return s.TotalSlots - s.BookedSlots
}

// IsFull 号源是否已约满。
func (s *Schedule) IsFull() bool {
	return s.BookedSlots >= s.TotalSlots
}

// ScheduleFilter 排班筛选条件。
type ScheduleFilter struct {
	DoctorID string
	Date     string
}

// Match 判断排班是否命中筛选条件。
func (f ScheduleFilter) Match(s *Schedule) bool {
	if f.DoctorID != "" && s.DoctorID != f.DoctorID {
		return false
	}
	if f.Date != "" && s.Date != f.Date {
		return false
	}
	return true
}
