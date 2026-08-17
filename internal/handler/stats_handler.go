package handler

import (
	"net/http"

	"appointment/pkg/httpx"
)

// registerStatsRoutes 注册统计相关路由。
func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.overview)
	mux.HandleFunc("GET /api/stats/appointments", s.appointmentStats)
	mux.HandleFunc("GET /api/stats/doctors", s.doctorSummaries)
}

func (s *Server) overview(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.Overview()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) appointmentStats(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.AppointmentStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) doctorSummaries(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.DoctorScheduleSummaries()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
