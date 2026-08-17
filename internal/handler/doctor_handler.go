package handler

import (
	"net/http"

	"appointment/internal/model"
	"appointment/pkg/httpx"
)

// registerDoctorRoutes 注册医生相关路由。
func (s *Server) registerDoctorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/doctors", s.createDoctor)
	mux.HandleFunc("GET /api/doctors", s.listDoctors)
	mux.HandleFunc("GET /api/doctors/{id}", s.getDoctor)
	mux.HandleFunc("PUT /api/doctors/{id}", s.updateDoctor)
	mux.HandleFunc("DELETE /api/doctors/{id}", s.deleteDoctor)
}

type doctorRequest struct {
	Name         string `json:"name"`
	Title        string `json:"title"`
	DepartmentID string `json:"department_id"`
	Specialty    string `json:"specialty"`
}

func (s *Server) createDoctor(w http.ResponseWriter, r *http.Request) {
	var req doctorRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDoctor(model.Doctor{
		Name:         req.Name,
		Title:        req.Title,
		DepartmentID: req.DepartmentID,
		Specialty:    req.Specialty,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

// listDoctors 医生列表：GET /api/doctors?department_id=
func (s *Server) listDoctors(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListDoctors(r.URL.Query().Get("department_id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) getDoctor(w http.ResponseWriter, r *http.Request) {
	d, err := s.svc.GetDoctor(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) updateDoctor(w http.ResponseWriter, r *http.Request) {
	var req doctorRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.UpdateDoctor(r.PathValue("id"), model.Doctor{
		Name:         req.Name,
		Title:        req.Title,
		DepartmentID: req.DepartmentID,
		Specialty:    req.Specialty,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDoctor(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteDoctor(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
