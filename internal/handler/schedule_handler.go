package handler

import (
	"net/http"

	"appointment/internal/model"
	"appointment/pkg/httpx"
)

// registerScheduleRoutes 注册排班相关路由。
func (s *Server) registerScheduleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/schedules", s.createSchedule)
	mux.HandleFunc("GET /api/schedules", s.listSchedules)
	mux.HandleFunc("GET /api/schedules/{id}", s.getSchedule)
	mux.HandleFunc("DELETE /api/schedules/{id}", s.deleteSchedule)
}

type createScheduleRequest struct {
	DoctorID   string `json:"doctor_id"`
	Date       string `json:"date"`       // "2006-01-02"
	TimeSlot   string `json:"time_slot"`  // "09:00-10:00"
	TotalSlots int    `json:"total_slots"`
}

func (s *Server) createSchedule(w http.ResponseWriter, r *http.Request) {
	var req createScheduleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sc, err := s.svc.CreateSchedule(model.Schedule{
		DoctorID:   req.DoctorID,
		Date:       req.Date,
		TimeSlot:   req.TimeSlot,
		TotalSlots: req.TotalSlots,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sc)
}

// listSchedules 排班列表：GET /api/schedules?doctor_id=&date=&page=&size=
func (s *Server) listSchedules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ScheduleFilter{
		DoctorID: r.URL.Query().Get("doctor_id"),
		Date:     r.URL.Query().Get("date"),
	}
	items, total, err := s.svc.ListSchedules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSchedule(w http.ResponseWriter, r *http.Request) {
	sc, err := s.svc.GetSchedule(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sc)
}

func (s *Server) deleteSchedule(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSchedule(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
