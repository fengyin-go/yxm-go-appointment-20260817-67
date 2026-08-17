package service

import (
	"testing"

	"appointment/internal/model"
)

func TestRegression_AppointmentPagination(t *testing.T) {
	s := newTestService()
	_, _, sid, _ := s.seed(t)
	sc, err := s.GetSchedule(sid)
	if err != nil {
		t.Fatalf("查询排班失败: %v", err)
	}
	sc.TotalSlots = 5
	if err := s.store.UpdateSchedule(sc); err != nil {
		t.Fatalf("更新排班失败: %v", err)
	}

	for i := 0; i < 3; i++ {
		p, err := s.CreatePatient(model.Patient{Name: "患者" + string(rune('A'+i)), Phone: "1390000000" + string(rune('0'+i))})
		if err != nil {
			t.Fatalf("创建患者失败: %v", err)
		}
		if _, err := s.BookAppointment(sid, p.ID); err != nil {
			t.Fatalf("预约失败: %v", err)
		}
	}

	page1, total, err := s.ListAppointments(model.AppointmentFilter{}, 1, 2)
	if err != nil {
		t.Fatalf("查询第一页失败: %v", err)
	}
	if total != 3 || len(page1) != 2 {
		t.Fatalf("第一页应返回 2 条，得到 total=%d len=%d", total, len(page1))
	}
	page2, _, err := s.ListAppointments(model.AppointmentFilter{}, 2, 2)
	if err != nil {
		t.Fatalf("查询第二页失败: %v", err)
	}
	if len(page2) != 1 {
		t.Fatalf("第二页应返回 1 条，得到 %d", len(page2))
	}
}

func TestRegression_SchedulePagination(t *testing.T) {
	s := newTestService()
	dep, err := s.CreateDepartment(model.Department{Name: "外科"})
	if err != nil {
		t.Fatalf("创建科室失败: %v", err)
	}
	doc, err := s.CreateDoctor(model.Doctor{Name: "王医生", DepartmentID: dep.ID})
	if err != nil {
		t.Fatalf("创建医生失败: %v", err)
	}
	for i := 0; i < 3; i++ {
		_, err := s.CreateSchedule(model.Schedule{
			DoctorID:   doc.ID,
			Date:       "2026-08-16",
			TimeSlot:   "0" + string(rune('9'+i)) + ":00-10:00",
			TotalSlots: 2,
		})
		if err != nil {
			t.Fatalf("创建排班失败: %v", err)
		}
	}

	page1, total, err := s.ListSchedules(model.ScheduleFilter{}, 1, 2)
	if err != nil {
		t.Fatalf("查询第一页失败: %v", err)
	}
	if total != 3 || len(page1) != 2 {
		t.Fatalf("第一页应返回 2 条，得到 total=%d len=%d", total, len(page1))
	}
	page2, _, err := s.ListSchedules(model.ScheduleFilter{}, 2, 2)
	if err != nil {
		t.Fatalf("查询第二页失败: %v", err)
	}
	if len(page2) != 1 {
		t.Fatalf("第二页应返回 1 条，得到 %d", len(page2))
	}
}

func TestRegression_ValidationErrorIsDetectable(t *testing.T) {
	s := newTestService()
	_, _, _, _ = s.seed(t)

	_, err := s.CreateSchedule(model.Schedule{
		DoctorID:   "missing-doctor",
		Date:       "bad-date",
		TimeSlot:   "09:00-10:00",
		TotalSlots: 1,
	})
	if !model.IsValidationError(err) {
		t.Fatalf("非法排班应识别为校验错误，得到 %T %v", err, err)
	}
}

func TestRegression_CreateScheduleDoesNotPanic(t *testing.T) {
	s := newTestService()
	dep, err := s.CreateDepartment(model.Department{Name: "儿科"})
	if err != nil {
		t.Fatalf("创建科室失败: %v", err)
	}
	doc, err := s.CreateDoctor(model.Doctor{Name: "刘医生", DepartmentID: dep.ID})
	if err != nil {
		t.Fatalf("创建医生失败: %v", err)
	}
	if _, err := s.CreateSchedule(model.Schedule{
		DoctorID:   doc.ID,
		Date:       "2026-08-16",
		TimeSlot:   "09:00-10:00",
		TotalSlots: 1,
	}); err != nil {
		t.Fatalf("创建排班失败: %v", err)
	}
}

func TestRegression_CancelledAppointmentCannotCheckIn(t *testing.T) {
	s := newTestService()
	_, _, sid, pid := s.seed(t)

	a, err := s.BookAppointment(sid, pid)
	if err != nil {
		t.Fatalf("预约失败: %v", err)
	}
	if _, err := s.CancelAppointment(a.ID); err != nil {
		t.Fatalf("取消失败: %v", err)
	}
	if _, err := s.CheckInAppointment(a.ID); !model.IsValidationError(err) {
		t.Fatalf("已取消预约不应再签到，得到 %v", err)
	}
}

func TestRegression_CancelReleasesOneSlot(t *testing.T) {
	s := newTestService()
	_, _, sid, pid := s.seed(t)

	a, err := s.BookAppointment(sid, pid)
	if err != nil {
		t.Fatalf("预约失败: %v", err)
	}
	if _, err := s.CancelAppointment(a.ID); err != nil {
		t.Fatalf("取消失败: %v", err)
	}
	sc, err := s.GetSchedule(sid)
	if err != nil {
		t.Fatalf("查询排班失败: %v", err)
	}
	if sc.BookedSlots != 0 || sc.Remaining() != 2 {
		t.Fatalf("取消后号源应恢复，得到 %+v", sc)
	}
}
