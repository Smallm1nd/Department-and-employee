package repository

import (
	"github.com/Smallm1nd/Department-and-employee/internal/models"
	"gorm.io/gorm"
)

type DepartmentRepo struct {
	conn *gorm.DB
}

func NewDepartmentRepo(conn *gorm.DB) *DepartmentRepo {
	return &DepartmentRepo{conn: conn}
}

func (dr *DepartmentRepo) Create(dept *models.Department) error {
	return dr.conn.Create(dept).Error
}

func (dr *DepartmentRepo) GetAll() ([]models.Department, error) {
	var depts []models.Department
	err := dr.conn.Find(&depts).Error
	return depts, err
}

func (dr *DepartmentRepo) GetByID(id int) (*models.Department, error) {
	var dept models.Department
	err := dr.conn.First(&dept, id).Error
	return &dept, err
}

// GetTree - динамически собирает дерево нужной глубины
func (dr *DepartmentRepo) GetTree(id int, depth int, includeEmp bool) (*models.Department, error) {
	var dept models.Department
	query := dr.conn.Model(&models.Department{})

	orderEmp := func(db *gorm.DB) *gorm.DB { return db.Order("full_name ASC") }

	if includeEmp {
		query = query.Preload("Employees", orderEmp)
	}

	var preloadStr string
	for i := 1; i < depth; i++ {
		if preloadStr == "" {
			preloadStr = "Children"
		} else {
			preloadStr += ".Children"
		}
		query = query.Preload(preloadStr)
		if includeEmp {
			query = query.Preload(preloadStr+".Employees", orderEmp)
		}
	}

	err := query.First(&dept, id).Error
	return &dept, err
}

func (dr *DepartmentRepo) Update(dept *models.Department) error {
	return dr.conn.Save(dept).Error
}

func (dr *DepartmentRepo) Delete(id int) error {

	return dr.conn.Delete(&models.Department{}, id).Error
}

func (dr *DepartmentRepo) ReassignEmployees(oldDeptID, newDeptID int) error {
	return dr.conn.Model(&models.Employee{}).Where("department_id = ?", oldDeptID).Update("department_id", newDeptID).Error
}
