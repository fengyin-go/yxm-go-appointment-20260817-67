// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"appointment/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法。
type Store interface {
	// 科室
	CreateDepartment(d *model.Department) error
	GetDepartment(id string) (*model.Department, error)
	ListDepartments() []*model.Department
	UpdateDepartment(d *model.Department) error
	DeleteDepartment(id string) error

	// 医生
	CreateDoctor(d *model.Doctor) error
	GetDoctor(id string) (*model.Doctor, error)
	ListDoctors() []*model.Doctor
	UpdateDoctor(d *model.Doctor) error
	DeleteDoctor(id string) error

	// 排班
	CreateSchedule(s *model.Schedule) error
	GetSchedule(id string) (*model.Schedule, error)
	ListSchedules() []*model.Schedule
	UpdateSchedule(s *model.Schedule) error
	DeleteSchedule(id string) error

	// 患者
	CreatePatient(p *model.Patient) error
	GetPatient(id string) (*model.Patient, error)
	ListPatients() []*model.Patient
	UpdatePatient(p *model.Patient) error
	DeletePatient(id string) error

	// 预约
	CreateAppointment(a *model.Appointment) error
	GetAppointment(id string) (*model.Appointment, error)
	ListAppointments() []*model.Appointment
	UpdateAppointment(a *model.Appointment) error
	DeleteAppointment(id string) error
}
