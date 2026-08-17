package handler

import (
	"net/http"

	"appointment/internal/model"
	"appointment/pkg/httpx"
)

// registerPatientRoutes 注册患者相关路由。
func (s *Server) registerPatientRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/patients", s.createPatient)
	mux.HandleFunc("GET /api/patients", s.listPatients)
	mux.HandleFunc("GET /api/patients/{id}", s.getPatient)
	mux.HandleFunc("PUT /api/patients/{id}", s.updatePatient)
	mux.HandleFunc("DELETE /api/patients/{id}", s.deletePatient)
}

type patientRequest struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	IDCard string `json:"id_card"`
	Gender string `json:"gender"`
}

func (s *Server) createPatient(w http.ResponseWriter, r *http.Request) {
	var req patientRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreatePatient(model.Patient{
		Name:   req.Name,
		Phone:  req.Phone,
		IDCard: req.IDCard,
		Gender: req.Gender,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listPatients(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListPatients()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) getPatient(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.GetPatient(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) updatePatient(w http.ResponseWriter, r *http.Request) {
	var req patientRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdatePatient(r.PathValue("id"), model.Patient{
		Name:   req.Name,
		Phone:  req.Phone,
		IDCard: req.IDCard,
		Gender: req.Gender,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deletePatient(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeletePatient(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
