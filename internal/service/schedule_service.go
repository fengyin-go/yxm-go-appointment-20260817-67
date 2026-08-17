package service

import (
	"sort"
	"time"

	"appointment/internal/model"
	"appointment/pkg/idgen"
)

// CreateSchedule 创建排班（号源）。
func (s *Service) CreateSchedule(input model.Schedule) (*model.Schedule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetDoctor(input.DoctorID); err != nil {
		return nil, err
	}
	sc := &model.Schedule{
		ID:          idgen.Hex(),
		DoctorID:    input.DoctorID,
		Date:        input.Date,
		TimeSlot:    input.TimeSlot,
		TotalSlots:  input.TotalSlots,
		BookedSlots: 0,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateSchedule(sc); err != nil {
		return nil, err
	}
	s.log.Infof("创建排班 doctor=%s date=%s slot=%s", input.DoctorID, input.Date, input.TimeSlot)
	return sc, nil
}

// GetSchedule 按 ID 查询排班。
func (s *Service) GetSchedule(id string) (*model.Schedule, error) {
	return s.store.GetSchedule(id)
}

// ListSchedules 列出排班，支持按医生、日期筛选与分页。
func (s *Service) ListSchedules(filter model.ScheduleFilter, page, size int) ([]*model.Schedule, int, error) {
	all := s.store.ListSchedules()
	matched := make([]*model.Schedule, 0, len(all))
	for _, sc := range all {
		if filter.Match(sc) {
			matched = append(matched, sc)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].Date != matched[j].Date {
			return matched[i].Date < matched[j].Date
		}
		return matched[i].TimeSlot < matched[j].TimeSlot
	})
	total := len(matched)
	if page < 1 {
		page = 1
	}
	// page 为 1 基：第一页对应偏移 0。
	start := (page - 1) * size
	if start >= total {
		return []*model.Schedule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// DeleteSchedule 删除排班。
func (s *Service) DeleteSchedule(id string) error {
	if err := s.store.DeleteSchedule(id); err != nil {
		return err
	}
	s.log.Infof("删除排班 %s", id)
	return nil
}
