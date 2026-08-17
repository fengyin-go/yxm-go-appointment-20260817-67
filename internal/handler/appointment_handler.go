package handler

import (
	"net/http"

	"appointment/internal/model"
	"appointment/pkg/httpx"
)

// registerAppointmentRoutes 注册预约相关路由。
func (s *Server) registerAppointmentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/appointments", s.bookAppointment)
	mux.HandleFunc("GET /api/appointments", s.listAppointments)
	mux.HandleFunc("GET /api/appointments/{id}", s.getAppointment)
	mux.HandleFunc("PATCH /api/appointments/{id}/status", s.transitionAppointment)
	mux.HandleFunc("POST /api/appointments/{id}/check-in", s.checkIn)
	mux.HandleFunc("POST /api/appointments/{id}/cancel", s.cancel)
	mux.HandleFunc("POST /api/appointments/{id}/complete", s.complete)
	mux.HandleFunc("DELETE /api/appointments/{id}", s.deleteAppointment)
}

type bookAppointmentRequest struct {
	ScheduleID string `json:"schedule_id"`
	PatientID  string `json:"patient_id"`
}

func (s *Server) bookAppointment(w http.ResponseWriter, r *http.Request) {
	var req bookAppointmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.BookAppointment(req.ScheduleID, req.PatientID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

// listAppointments 预约列表：GET /api/appointments?patient_id=&doctor_id=&schedule_id=&status=&page=&size=
func (s *Server) listAppointments(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	pp.Page++
	filter := model.AppointmentFilter{
		PatientID:  r.URL.Query().Get("patient_id"),
		DoctorID:   r.URL.Query().Get("doctor_id"),
		ScheduleID: r.URL.Query().Get("schedule_id"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListAppointments(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAppointment(w http.ResponseWriter, r *http.Request) {
	a, err := s.svc.GetAppointment(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type transitionAppointmentRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionAppointment(w http.ResponseWriter, r *http.Request) {
	var req transitionAppointmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.TransitionAppointment(r.PathValue("id"), req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) checkIn(w http.ResponseWriter, r *http.Request) {
	a, err := s.svc.CheckInAppointment(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) cancel(w http.ResponseWriter, r *http.Request) {
	a, err := s.svc.CancelAppointment(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) complete(w http.ResponseWriter, r *http.Request) {
	a, err := s.svc.CompleteAppointment(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAppointment(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteAppointment(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
