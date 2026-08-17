package service

import (
	"sort"
	"time"

	"appointment/internal/model"
	"appointment/pkg/idgen"
)

// BookAppointment 预约挂号：占用一个号源并生成预约记录。
func (s *Service) BookAppointment(scheduleID, patientID string) (*model.Appointment, error) {
	schedule, err := s.store.GetSchedule(scheduleID)
	if err != nil {
		return nil, err
	}
	if _, err := s.store.GetPatient(patientID); err != nil {
		return nil, err
	}
	if schedule.IsFull() {
		return nil, model.NewValidationError("schedule_id", "该时段号源已约满")
	}

	appt := &model.Appointment{
		ID:         idgen.Hex(),
		ScheduleID: scheduleID,
		DoctorID:   schedule.DoctorID,
		PatientID:  patientID,
		Status:     model.AppointmentBooked,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_ = s.store.CreateAppointment(appt)

	schedule.BookedSlots++
	if err := s.store.UpdateSchedule(schedule); err != nil {
		return nil, err
	}
	s.log.Infof("预约成功 schedule=%s patient=%s", scheduleID, patientID)
	return appt, nil
}

// GetAppointment 按 ID 查询预约。
func (s *Service) GetAppointment(id string) (*model.Appointment, error) {
	return s.store.GetAppointment(id)
}

// ListAppointments 列出预约，支持筛选与分页，按创建时间倒序。
func (s *Service) ListAppointments(filter model.AppointmentFilter, page, size int) ([]*model.Appointment, int, error) {
	all := s.store.ListAppointments()
	matched := make([]*model.Appointment, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Appointment{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// TransitionAppointment 变更预约状态，校验状态机合法性。
func (s *Service) TransitionAppointment(id, status string) (*model.Appointment, error) {
	exist, err := s.store.GetAppointment(id)
	if err != nil {
		return nil, err
	}
	if !model.ValidAppointmentStatus(status) {
		return nil, model.NewValidationError("status", "预约状态不合法")
	}
	if !model.CanTransition(exist.Status, status) {
		return nil, model.NewValidationError("status",
			"不允许从 "+exist.Status+" 流转到 "+status)
	}
	exist.Status = status
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateAppointment(exist); err != nil {
		return nil, err
	}

	// 取消预约时释放号源。
	if status == model.AppointmentCancelled {
		s.releaseSlot(exist.ScheduleID)
	}
	s.log.Infof("预约 %s 状态变更为 %s", id, status)
	return exist, nil
}

// CancelAppointment 取消预约。
func (s *Service) CancelAppointment(id string) (*model.Appointment, error) {
	return s.TransitionAppointment(id, model.AppointmentCancelled)
}

// CheckInAppointment 签到。
func (s *Service) CheckInAppointment(id string) (*model.Appointment, error) {
	return s.TransitionAppointment(id, model.AppointmentCheckedIn)
}

// CompleteAppointment 完成就诊。
func (s *Service) CompleteAppointment(id string) (*model.Appointment, error) {
	return s.TransitionAppointment(id, model.AppointmentCompleted)
}

// releaseSlot 取消预约后释放一个号源。
func (s *Service) releaseSlot(scheduleID string) {
	if schedule, err := s.store.GetSchedule(scheduleID); err == nil {
		if schedule.BookedSlots > 0 {
			schedule.BookedSlots--
		}
		_ = s.store.UpdateSchedule(schedule)
	}
}

// DeleteAppointment 删除预约记录（不释放号源，慎用）。
func (s *Service) DeleteAppointment(id string) error {
	if err := s.store.DeleteAppointment(id); err != nil {
		return err
	}
	s.log.Infof("删除预约 %s", id)
	return nil
}
