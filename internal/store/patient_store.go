package store

import "appointment/internal/model"

// CreatePatient 新增患者，手机号重复时返回 ErrConflict。
func (s *MemoryStore) CreatePatient(p *model.Patient) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.patients {
		if exist.Phone == p.Phone {
			return ErrConflict
		}
	}
	s.patients[p.ID] = p
	return nil
}

// GetPatient 按 ID 查询患者。
func (s *MemoryStore) GetPatient(id string) (*model.Patient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.patients[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

// ListPatients 返回全部患者。
func (s *MemoryStore) ListPatients() []*model.Patient {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Patient, 0, len(s.patients))
	for _, p := range s.patients {
		list = append(list, p)
	}
	return list
}

// UpdatePatient 覆盖保存患者。
func (s *MemoryStore) UpdatePatient(p *model.Patient) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.patients[p.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.patients {
		if exist.ID != p.ID && exist.Phone == p.Phone {
			return ErrConflict
		}
	}
	s.patients[p.ID] = p
	return nil
}

// DeletePatient 按 ID 删除患者。
func (s *MemoryStore) DeletePatient(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.patients[id]; !ok {
		return ErrNotFound
	}
	delete(s.patients, id)
	return nil
}
