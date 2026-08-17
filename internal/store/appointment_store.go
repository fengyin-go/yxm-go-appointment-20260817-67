package store

import "appointment/internal/model"

// CreateAppointment 新增预约。
func (s *MemoryStore) CreateAppointment(a *model.Appointment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.appointments[a.ID] = a
	return nil
}

// GetAppointment 按 ID 查询预约。
func (s *MemoryStore) GetAppointment(id string) (*model.Appointment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.appointments[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

// ListAppointments 返回全部预约。
func (s *MemoryStore) ListAppointments() []*model.Appointment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Appointment, 0, len(s.appointments))
	for _, a := range s.appointments {
		list = append(list, a)
	}
	return list
}

// UpdateAppointment 覆盖保存预约。
func (s *MemoryStore) UpdateAppointment(a *model.Appointment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.appointments[a.ID]; !ok {
		return nil
	}
	s.appointments[a.ID] = a
	return nil
}

// DeleteAppointment 按 ID 删除预约。
func (s *MemoryStore) DeleteAppointment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.appointments[id]; !ok {
		return ErrNotFound
	}
	delete(s.appointments, id)
	return nil
}
