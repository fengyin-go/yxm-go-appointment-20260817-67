package service

import (
	"testing"

	"appointment/internal/config"
	"appointment/internal/model"
	"appointment/internal/store"
	"appointment/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

// seed 准备科室、医生、排班与患者，返回各 ID。
func (s *Service) seed(t *testing.T) (departmentID, doctorID, scheduleID, patientID string) {
	t.Helper()
	dep, err := s.CreateDepartment(model.Department{Name: "内科"})
	if err != nil {
		t.Fatalf("创建科室失败: %v", err)
	}
	doc, err := s.CreateDoctor(model.Doctor{Name: "张医生", DepartmentID: dep.ID})
	if err != nil {
		t.Fatalf("创建医生失败: %v", err)
	}
	sc, err := s.CreateSchedule(model.Schedule{
		DoctorID:   doc.ID,
		Date:       "2026-08-16",
		TimeSlot:   "09:00-10:00",
		TotalSlots: 2,
	})
	if err != nil {
		t.Fatalf("创建排班失败: %v", err)
	}
	p, err := s.CreatePatient(model.Patient{Name: "李四", Phone: "13800000000"})
	if err != nil {
		t.Fatalf("创建患者失败: %v", err)
	}
	return dep.ID, doc.ID, sc.ID, p.ID
}

func TestService_BookAppointment(t *testing.T) {
	s := newTestService()
	_, _, sid, pid := s.seed(t)

	a, err := s.BookAppointment(sid, pid)
	if err != nil {
		t.Fatalf("预约失败: %v", err)
	}
	if a.Status != model.AppointmentBooked || a.DoctorID == "" {
		t.Fatalf("预约记录异常: %+v", a)
	}

	sc, _ := s.GetSchedule(sid)
	if sc.BookedSlots != 1 || sc.Remaining() != 1 {
		t.Fatalf("号源占用异常: %+v", sc)
	}
}

func TestService_BookFullSchedule(t *testing.T) {
	s := newTestService()
	_, _, sid, pid := s.seed(t)

	// TotalSlots=2，预约两次后应约满
	if _, err := s.BookAppointment(sid, pid); err != nil {
		t.Fatalf("首次预约失败: %v", err)
	}
	p2, _ := s.CreatePatient(model.Patient{Name: "王五", Phone: "13900000000"})
	if _, err := s.BookAppointment(sid, p2.ID); err != nil {
		t.Fatalf("第二次预约失败: %v", err)
	}
	p3, _ := s.CreatePatient(model.Patient{Name: "赵六", Phone: "13700000000"})
	if _, err := s.BookAppointment(sid, p3.ID); !model.IsValidationError(err) {
		t.Fatalf("期望号源已满校验错误，得到 %v", err)
	}
}

func TestService_CancelReleasesSlot(t *testing.T) {
	s := newTestService()
	_, _, sid, pid := s.seed(t)

	a, err := s.BookAppointment(sid, pid)
	if err != nil {
		t.Fatalf("预约失败: %v", err)
	}
	if _, err := s.CancelAppointment(a.ID); err != nil {
		t.Fatalf("取消失败: %v", err)
	}

	sc, _ := s.GetSchedule(sid)
	if sc.BookedSlots != 0 {
		t.Fatalf("取消后应释放号源，得到 booked=%d", sc.BookedSlots)
	}
}

func TestService_AppointmentLifecycle(t *testing.T) {
	s := newTestService()
	_, _, sid, pid := s.seed(t)

	a, _ := s.BookAppointment(sid, pid)

	// booked -> checked_in
	a, err := s.CheckInAppointment(a.ID)
	if err != nil {
		t.Fatalf("签到失败: %v", err)
	}
	if a.Status != model.AppointmentCheckedIn {
		t.Fatalf("期望 checked_in，得到 %s", a.Status)
	}

	// checked_in -> completed
	a, err = s.CompleteAppointment(a.ID)
	if err != nil {
		t.Fatalf("完成失败: %v", err)
	}
	if a.Status != model.AppointmentCompleted {
		t.Fatalf("期望 completed，得到 %s", a.Status)
	}

	// completed -> cancelled 非法
	if _, err := s.TransitionAppointment(a.ID, model.AppointmentCancelled); !model.IsValidationError(err) {
		t.Fatalf("期望校验错误，得到 %v", err)
	}
}

func TestService_Stats(t *testing.T) {
	s := newTestService()
	_, _, sid, pid := s.seed(t)
	if _, err := s.BookAppointment(sid, pid); err != nil {
		t.Fatalf("预约失败: %v", err)
	}

	ov, err := s.Overview()
	if err != nil {
		t.Fatalf("总览失败: %v", err)
	}
	if ov.DepartmentCount != 1 || ov.DoctorCount != 1 || ov.AppointmentCount != 1 {
		t.Fatalf("总览异常: %+v", ov)
	}
	if ov.AvailableSlots != 1 {
		t.Fatalf("剩余号源应为 1，得到 %d", ov.AvailableSlots)
	}

	stats, _ := s.AppointmentStats()
	if stats.ByStatus[model.AppointmentBooked] != 1 {
		t.Fatalf("预约统计异常: %+v", stats)
	}
}
