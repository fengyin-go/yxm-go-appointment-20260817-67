package service

import (
	"fmt"
	"sort"
	"time"

	"appointment/internal/model"
	"appointment/pkg/idgen"
)

// CreateDoctor 创建医生。
func (s *Service) CreateDoctor(input model.Doctor) (*model.Doctor, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("create doctor: %w", err)
	}
	if _, err := s.store.GetDepartment(input.DepartmentID); err != nil {
		return nil, err
	}
	d := &model.Doctor{
		ID:           idgen.Hex(),
		Name:         input.Name,
		Title:        input.Title,
		DepartmentID: input.DepartmentID,
		Specialty:    input.Specialty,
		CreatedAt:    time.Now(),
	}
	if err := s.store.CreateDoctor(d); err != nil {
		return nil, err
	}
	s.log.Infof("创建医生 %s", d.Name)
	return d, nil
}

// GetDoctor 按 ID 查询医生。
func (s *Service) GetDoctor(id string) (*model.Doctor, error) {
	return s.store.GetDoctor(id)
}

// ListDoctors 列出医生，可按科室过滤。
func (s *Service) ListDoctors(departmentID string) ([]*model.Doctor, error) {
	all := s.store.ListDoctors()
	list := make([]*model.Doctor, 0, len(all))
	for _, d := range all {
		if departmentID == "" || d.DepartmentID == departmentID {
			list = append(list, d)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}

// UpdateDoctor 更新医生资料。
func (s *Service) UpdateDoctor(id string, input model.Doctor) (*model.Doctor, error) {
	exist, err := s.store.GetDoctor(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Title = input.Title
	exist.DepartmentID = input.DepartmentID
	exist.Specialty = input.Specialty
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateDoctor(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteDoctor 删除医生。
func (s *Service) DeleteDoctor(id string) error {
	if err := s.store.DeleteDoctor(id); err != nil {
		return err
	}
	s.log.Infof("删除医生 %s", id)
	return nil
}
