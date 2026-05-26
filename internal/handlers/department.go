package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Smallm1nd/Department-and-employee/internal/models"
	"github.com/Smallm1nd/Department-and-employee/internal/service"
)

type DepartmentHandler struct {
	service service.Department
}

func NewDepartmentHandler(svc service.Department) *DepartmentHandler {
	return &DepartmentHandler{service: svc}
}

func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dept models.Department
	if err := json.NewDecoder(r.Body).Decode(&dept); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.service.Create(&dept); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	depts, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(depts)
}

func (h *DepartmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	depth, _ := strconv.Atoi(r.URL.Query().Get("depth"))
	if depth == 0 {
		depth = 1
	}

	includeEmpStr := r.URL.Query().Get("include_employees")
	includeEmp := true // По умолчанию true
	if includeEmpStr == "false" {
		includeEmp = false
	}

	dept, err := h.service.GetTree(id, depth, includeEmp)
	if err != nil {
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	var req struct {
		Name     *string `json:"name"`
		ParentID *int    `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dept, err := h.service.Update(id, req.Name, req.ParentID)
	if err != nil {
		if err.Error() == "circular dependency detected" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "cascade"
	}

	var reassignTo *int
	reassignStr := r.URL.Query().Get("reassign_to_department_id")
	if reassignStr != "" {
		val, _ := strconv.Atoi(reassignStr)
		reassignTo = &val
	}

	if err := h.service.Delete(id, mode, reassignTo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
