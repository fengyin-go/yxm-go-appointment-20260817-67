package service

import (
	"sort"
	"time"

	"appointment/internal/model"
	"appointment/pkg/idgen"
)

// CreatePatient 创建患者。
func (s *Service) CreatePatient(input model.Patient) (*model.Patient, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	p := &model.Patient{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Phone:     input.Phone,
		IDCard:    input.IDCard,
		Gender:    input.Gender,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreatePatient(p); err != nil {
		return nil, err
	}
	s.log.Infof("创建患者 %s", p.Name)
	return p, nil
}

// GetPatient 按 ID 查询患者。
func (s *Service) GetPatient(id string) (*model.Patient, error) {
	return s.store.GetPatient(id)
}

// ListPatients 列出全部患者。
func (s *Service) ListPatients() ([]*model.Patient, error) {
	list := s.store.ListPatients()
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}

// UpdatePatient 更新患者资料。
func (s *Service) UpdatePatient(id string, input model.Patient) (*model.Patient, error) {
	exist, err := s.store.GetPatient(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Phone = input.Phone
	exist.IDCard = input.IDCard
	exist.Gender = input.Gender
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdatePatient(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeletePatient 删除患者。
func (s *Service) DeletePatient(id string) error {
	if err := s.store.DeletePatient(id); err != nil {
		return err
	}
	s.log.Infof("删除患者 %s", id)
	return nil
}
