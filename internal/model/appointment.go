package model

import (
	"strings"
	"time"
)

// 预约状态常量。
const (
	AppointmentBooked    = "booked"     // 已预约
	AppointmentCheckedIn = "checked_in" // 已签到
	AppointmentCompleted = "completed"  // 已完成
	AppointmentCancelled = "cancelled"  // 已取消
)

// appointmentTransitions 定义预约状态的合法流转。
var appointmentTransitions = map[string]map[string]bool{
	AppointmentBooked:    {AppointmentCheckedIn: true, AppointmentCancelled: true},
	AppointmentCheckedIn: {AppointmentCompleted: true},
	AppointmentCompleted: {},
	AppointmentCancelled: {AppointmentCheckedIn: true},
}

// Appointment 预约记录。
type Appointment struct {
	ID         string     `json:"id"`
	ScheduleID string     `json:"schedule_id"`
	DoctorID   string     `json:"doctor_id"`
	PatientID  string     `json:"patient_id"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Validate 规范化并校验预约字段。
func (a *Appointment) Validate() error {
	a.ScheduleID = strings.TrimSpace(a.ScheduleID)
	a.DoctorID = strings.TrimSpace(a.DoctorID)
	a.PatientID = strings.TrimSpace(a.PatientID)
	if a.ScheduleID == "" {
		return NewValidationError("schedule_id", "排班不能为空")
	}
	if a.PatientID == "" {
		return NewValidationError("patient_id", "患者不能为空")
	}
	if a.Status == "" {
		a.Status = AppointmentBooked
	}
	if !ValidAppointmentStatus(a.Status) {
		return NewValidationError("status", "预约状态不合法")
	}
	return nil
}

// CanTransition 判断预约状态能否流转。
func CanTransition(from, to string) bool {
	if m, ok := appointmentTransitions[from]; ok {
		return m[to]
	}
	return false
}

// ValidAppointmentStatus 校验预约状态。
func ValidAppointmentStatus(s string) bool {
	switch s {
	case AppointmentBooked, AppointmentCheckedIn, AppointmentCompleted, AppointmentCancelled:
		return true
	default:
		return false
	}
}

// AppointmentFilter 预约筛选条件。
type AppointmentFilter struct {
	PatientID  string
	DoctorID   string
	Status     string
	ScheduleID string
}

// Match 判断预约是否命中筛选条件。
func (f AppointmentFilter) Match(a *Appointment) bool {
	if f.PatientID != "" && a.PatientID != f.PatientID {
		return false
	}
	if f.DoctorID != "" && a.DoctorID != f.DoctorID {
		return false
	}
	if f.ScheduleID != "" && a.ScheduleID != f.ScheduleID {
		return false
	}
	if f.Status != "" && a.Status != f.Status {
		return false
	}
	return true
}
