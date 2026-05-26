package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Smallm1nd/Department-and-employee/internal/models"
	"github.com/Smallm1nd/Department-and-employee/internal/service"
)

type EmployeeHandler struct {
	service service.Employee
}

func NewEmployeeHandler(svc service.Employee) *EmployeeHandler {
	return &EmployeeHandler{service: svc}
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {

	deptIDStr := r.PathValue("id")
	deptID, err := strconv.Atoi(deptIDStr)
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	emp.DepartmentID = deptID

	if err := h.service.Create(&emp); err != nil {
		if err.Error() == "department not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}
