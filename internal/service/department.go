package service

import (
	"errors"
	"strings"

	"github.com/Smallm1nd/Department-and-employee/internal/models"
	"github.com/Smallm1nd/Department-and-employee/internal/repository"
)

type Department interface {
	Create(dept *models.Department) error
	GetAll() ([]models.Department, error)
	GetTree(id int, depth int, includeEmp bool) (*models.Department, error)
	Update(id int, name *string, parentID *int) (*models.Department, error)
	Delete(id int, mode string, reassignTo *int) error
}

type departmentService struct {
	repo *repository.DepartmentRepo
}

func NewDepartmentService(repo *repository.DepartmentRepo) Department {
	return &departmentService{repo: repo}
}

func (ds *departmentService) Create(dept *models.Department) error {
	dept.Name = strings.TrimSpace(dept.Name)
	if dept.Name == "" {
		return errors.New("department name cannot be empty")
	}
	return ds.repo.Create(dept)
}

func (ds *departmentService) GetAll() ([]models.Department, error) {
	return ds.repo.GetAll()
}

func (ds *departmentService) GetTree(id int, depth int, includeEmp bool) (*models.Department, error) {
	if depth < 1 {
		depth = 1
	}
	if depth > 5 {
		depth = 5
	}
	return ds.repo.GetTree(id, depth, includeEmp)
}

func (ds *departmentService) Update(id int, name *string, parentID *int) (*models.Department, error) {
	dept, err := ds.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("department not found")
	}

	if name != nil {
		trimmedName := strings.TrimSpace(*name)
		if trimmedName == "" {
			return nil, errors.New("name cannot be empty")
		}
		dept.Name = trimmedName
	}

	if parentID != nil {
		if *parentID == id {
			return nil, errors.New("cannot set department as its own parent")
		}
		if ds.isCyclic(id, *parentID) {
			return nil, errors.New("circular dependency detected")
		}
		dept.ParentID = parentID
	}

	if err := ds.repo.Update(dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (ds *departmentService) Delete(id int, mode string, reassignTo *int) error {
	if mode == "reassign" {
		if reassignTo == nil {
			return errors.New("reassign_to_department_id is required for reassign mode")
		}
		if *reassignTo == id {
			return errors.New("cannot reassign to the department being deleted")
		}
		if _, err := ds.repo.GetByID(*reassignTo); err != nil {
			return errors.New("target reassign department not found")
		}
		if err := ds.repo.ReassignEmployees(id, *reassignTo); err != nil {
			return err
		}
	}

	return ds.repo.Delete(id)
}

func (ds *departmentService) isCyclic(deptID int, targetParentID int) bool {
	currentID := targetParentID
	for currentID != 0 {
		if currentID == deptID {
			return true
		}
		parent, err := ds.repo.GetByID(currentID)
		if err != nil || parent.ParentID == nil {
			break
		}
		currentID = *parent.ParentID
	}
	return false
}
