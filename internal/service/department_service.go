package service

import (
	"sort"
	"time"

	"appointment/internal/model"
	"appointment/pkg/idgen"
)

// CreateDepartment 创建科室。
func (s *Service) CreateDepartment(input model.Department) (*model.Department, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	d := &model.Department{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateDepartment(d); err != nil {
		return nil, err
	}
	s.log.Infof("创建科室 %s", d.Name)
	return d, nil
}

// GetDepartment 按 ID 查询科室。
func (s *Service) GetDepartment(id string) (*model.Department, error) {
	return s.store.GetDepartment(id)
}

// ListDepartments 列出全部科室。
func (s *Service) ListDepartments() ([]*model.Department, error) {
	list := s.store.ListDepartments()
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}

// UpdateDepartment 更新科室。
func (s *Service) UpdateDepartment(id string, input model.Department) (*model.Department, error) {
	exist, err := s.store.GetDepartment(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Description = input.Description
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateDepartment(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteDepartment 删除科室。
func (s *Service) DeleteDepartment(id string) error {
	if err := s.store.DeleteDepartment(id); err != nil {
		return err
	}
	s.log.Infof("删除科室 %s", id)
	return nil
}
