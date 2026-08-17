package store

import "appointment/internal/model"

// CreateDoctor 新增医生。
func (s *MemoryStore) CreateDoctor(d *model.Doctor) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.doctors[d.ID] = d
	return nil
}

// GetDoctor 按 ID 查询医生。
func (s *MemoryStore) GetDoctor(id string) (*model.Doctor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.doctors[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

// ListDoctors 返回全部医生。
func (s *MemoryStore) ListDoctors() []*model.Doctor {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Doctor, 0, len(s.doctors))
	for _, d := range s.doctors {
		list = append(list, d)
	}
	return list
}

// UpdateDoctor 覆盖保存医生。
func (s *MemoryStore) UpdateDoctor(d *model.Doctor) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.doctors[d.ID]; !ok {
		return ErrNotFound
	}
	s.doctors[d.ID] = d
	return nil
}

// DeleteDoctor 按 ID 删除医生。
func (s *MemoryStore) DeleteDoctor(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.doctors[id]; !ok {
		return ErrNotFound
	}
	delete(s.doctors, id)
	return nil
}
