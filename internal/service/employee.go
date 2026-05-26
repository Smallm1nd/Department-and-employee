package service

import (
	"errors"
	"strings"

	"github.com/Smallm1nd/Department-and-employee/internal/models"
	"github.com/Smallm1nd/Department-and-employee/internal/repository"
)

type Employee interface {
	Create(emp *models.Employee) error
}

type employeeService struct {
	empRepo  *repository.EmployeeRepo
	deptRepo *repository.DepartmentRepo
}

func NewEmployeeService(empRepo *repository.EmployeeRepo, deptRepo *repository.DepartmentRepo) Employee {
	return &employeeService{empRepo: empRepo, deptRepo: deptRepo}
}

func (es *employeeService) Create(emp *models.Employee) error {
	emp.FullName = strings.TrimSpace(emp.FullName)
	emp.Position = strings.TrimSpace(emp.Position)

	if emp.FullName == "" || emp.Position == "" {
		return errors.New("full_name and position are required")
	}

	_, err := es.deptRepo.GetByID(emp.DepartmentID)
	if err != nil {
		return errors.New("department not found")
	}

	return es.empRepo.Create(emp)
}
