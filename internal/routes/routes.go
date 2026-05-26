package routes

import (
	"net/http"

	"github.com/Smallm1nd/Department-and-employee/internal/handlers"
)

func NewRouter(deptHandler *handlers.DepartmentHandler, empHandler *handlers.EmployeeHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /departments", deptHandler.Create)
	mux.HandleFunc("GET /departments", deptHandler.GetAll)
	mux.HandleFunc("GET /departments/{id}", deptHandler.GetByID)
	mux.HandleFunc("PATCH /departments/{id}", deptHandler.Update)
	mux.HandleFunc("DELETE /departments/{id}", deptHandler.Delete)

	mux.HandleFunc("POST /departments/{id}/employees", empHandler.Create)

	return mux
}
