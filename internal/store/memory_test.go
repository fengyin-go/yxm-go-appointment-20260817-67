package store

import (
	"testing"
	"time"

	"appointment/internal/model"
)

func TestMemoryStore_Department(t *testing.T) {
	s := NewMemoryStore()
	d := &model.Department{ID: "d1", Name: "内科", CreatedAt: time.Now()}
	if err := s.CreateDepartment(d); err != nil {
		t.Fatalf("创建科室失败: %v", err)
	}
	if err := s.CreateDepartment(&model.Department{ID: "d2", Name: "内科"}); err != ErrConflict {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}
	got, err := s.GetDepartment("d1")
	if err != nil || got.Name != "内科" {
		t.Fatalf("查询科室失败: %v, %v", got, err)
	}
}

func TestMemoryStore_Doctor(t *testing.T) {
	s := NewMemoryStore()
	doc := &model.Doctor{ID: "doc1", Name: "张医生", DepartmentID: "d1", CreatedAt: time.Now()}
	if err := s.CreateDoctor(doc); err != nil {
		t.Fatalf("创建医生失败: %v", err)
	}
	if len(s.ListDoctors()) != 1 {
		t.Fatalf("期望 1 个医生")
	}
}

func TestMemoryStore_Schedule(t *testing.T) {
	s := NewMemoryStore()
	sc := &model.Schedule{ID: "sc1", DoctorID: "doc1", Date: "2026-08-16", TimeSlot: "09:00-10:00", TotalSlots: 10, CreatedAt: time.Now()}
	if err := s.CreateSchedule(sc); err != nil {
		t.Fatalf("创建排班失败: %v", err)
	}
	sc.BookedSlots = 5
	if err := s.UpdateSchedule(sc); err != nil {
		t.Fatalf("更新排班失败: %v", err)
	}
	got, _ := s.GetSchedule("sc1")
	if got.BookedSlots != 5 || got.Remaining() != 5 {
		t.Fatalf("号源计算异常: %+v", got)
	}
}

func TestMemoryStore_PatientAndAppointment(t *testing.T) {
	s := NewMemoryStore()
	p := &model.Patient{ID: "p1", Name: "李四", Phone: "13800000000", CreatedAt: time.Now()}
	if err := s.CreatePatient(p); err != nil {
		t.Fatalf("创建患者失败: %v", err)
	}
	if err := s.CreatePatient(&model.Patient{ID: "p2", Phone: "13800000000"}); err != ErrConflict {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}

	a := &model.Appointment{ID: "a1", ScheduleID: "sc1", DoctorID: "doc1", PatientID: "p1", Status: model.AppointmentBooked, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateAppointment(a); err != nil {
		t.Fatalf("创建预约失败: %v", err)
	}
	if len(s.ListAppointments()) != 1 {
		t.Fatalf("期望 1 条预约")
	}
}
