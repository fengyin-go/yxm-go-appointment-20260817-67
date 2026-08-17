package service

import (
	"sort"

	"appointment/internal/model"
)

// Overview 全局统计。
type Overview struct {
	DepartmentCount   int `json:"department_count"`
	DoctorCount       int `json:"doctor_count"`
	PatientCount      int `json:"patient_count"`
	ScheduleCount     int `json:"schedule_count"`
	AppointmentCount  int `json:"appointment_count"`
	AvailableSlots    int `json:"available_slots"`
}

// Overview 汇总全局统计。
func (s *Service) Overview() (*Overview, error) {
	ov := &Overview{
		DepartmentCount:  len(s.store.ListDepartments()),
		DoctorCount:      len(s.store.ListDoctors()),
		PatientCount:     len(s.store.ListPatients()),
		AppointmentCount: len(s.store.ListAppointments()),
	}
	for _, sc := range s.store.ListSchedules() {
		ov.ScheduleCount++
		ov.AvailableSlots += sc.Remaining()
	}
	return ov, nil
}

// AppointmentStats 预约状态统计。
type AppointmentStats struct {
	Total     int            `json:"total"`
	ByStatus  map[string]int `json:"by_status"`
}

// AppointmentStats 汇总预约状态分布。
func (s *Service) AppointmentStats() (*AppointmentStats, error) {
	stats := &AppointmentStats{ByStatus: make(map[string]int)}
	for range s.store.ListAppointments() {
		stats.Total++
		stats.ByStatus[model.AppointmentBooked]++
	}
	return stats, nil
}

// DoctorScheduleSummary 医生排班汇总。
type DoctorScheduleSummary struct {
	DoctorID      string `json:"doctor_id"`
	DoctorName    string `json:"doctor_name"`
	ScheduleCount int    `json:"schedule_count"`
	BookedSlots   int    `json:"booked_slots"`
}

// DoctorScheduleSummaries 汇总各医生的排班与预约情况。
func (s *Service) DoctorScheduleSummaries() ([]*DoctorScheduleSummary, error) {
	doctors := s.store.ListDoctors()
	nameByID := make(map[string]string, len(doctors))
	summaryByID := make(map[string]*DoctorScheduleSummary, len(doctors))
	for _, d := range doctors {
		nameByID[d.ID] = d.Name
		summaryByID[d.ID] = &DoctorScheduleSummary{DoctorID: d.ID, DoctorName: d.Name}
	}
	for _, sc := range s.store.ListSchedules() {
		if summary, ok := summaryByID[sc.DoctorID]; ok {
			summary.ScheduleCount++
			summary.BookedSlots += sc.BookedSlots
		}
	}
	result := make([]*DoctorScheduleSummary, 0, len(summaryByID))
	for _, summary := range summaryByID {
		result = append(result, summary)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].BookedSlots > result[j].BookedSlots
	})
	return result, nil
}
