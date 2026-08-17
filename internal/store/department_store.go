package store

import "appointment/internal/model"

// CreateDepartment 新增科室，名称重复时返回 ErrConflict。
func (s *MemoryStore) CreateDepartment(d *model.Department) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.departments {
		if exist.Name == d.Name {
			return ErrConflict
		}
	}
	s.departments[d.ID] = d
	return nil
}

// GetDepartment 按 ID 查询科室。
func (s *MemoryStore) GetDepartment(id string) (*model.Department, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.departments[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

// ListDepartments 返回全部科室。
func (s *MemoryStore) ListDepartments() []*model.Department {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Department, 0, len(s.departments))
	for _, d := range s.departments {
		list = append(list, d)
	}
	return list
}

// UpdateDepartment 覆盖保存科室。
func (s *MemoryStore) UpdateDepartment(d *model.Department) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.departments[d.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.departments {
		if exist.ID != d.ID && exist.Name == d.Name {
			return ErrConflict
		}
	}
	s.departments[d.ID] = d
	return nil
}

// DeleteDepartment 按 ID 删除科室。
func (s *MemoryStore) DeleteDepartment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.departments[id]; !ok {
		return ErrNotFound
	}
	delete(s.departments, id)
	return nil
}
