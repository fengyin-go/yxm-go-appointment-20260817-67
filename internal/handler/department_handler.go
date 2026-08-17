package handler

import (
	"net/http"

	"appointment/internal/model"
	"appointment/pkg/httpx"
)

// registerDepartmentRoutes 注册科室相关路由。
func (s *Server) registerDepartmentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/departments", s.createDepartment)
	mux.HandleFunc("GET /api/departments", s.listDepartments)
	mux.HandleFunc("GET /api/departments/{id}", s.getDepartment)
	mux.HandleFunc("PUT /api/departments/{id}", s.updateDepartment)
	mux.HandleFunc("DELETE /api/departments/{id}", s.deleteDepartment)
}

type departmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) createDepartment(w http.ResponseWriter, r *http.Request) {
	var req departmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDepartment(model.Department{Name: req.Name, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDepartments(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListDepartments()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) getDepartment(w http.ResponseWriter, r *http.Request) {
	d, err := s.svc.GetDepartment(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) updateDepartment(w http.ResponseWriter, r *http.Request) {
	var req departmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.UpdateDepartment(r.PathValue("id"), model.Department{Name: req.Name, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDepartment(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteDepartment(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
