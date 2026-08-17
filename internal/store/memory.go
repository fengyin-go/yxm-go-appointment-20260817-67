package store

import (
	"sync"

	"appointment/internal/model"
)

// MemoryStore 基于内存的 Store 实现。
type MemoryStore struct {
	mu           sync.RWMutex
	departments  map[string]*model.Department
	doctors      map[string]*model.Doctor
	schedules    map[string]*model.Schedule
	patients     map[string]*model.Patient
	appointments map[string]*model.Appointment
}

// NewMemoryStore 创建空的内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		departments:  make(map[string]*model.Department),
		doctors:      make(map[string]*model.Doctor),
		schedules:    make(map[string]*model.Schedule),
		patients:     make(map[string]*model.Patient),
	}
}

var _ Store = (*MemoryStore)(nil)
