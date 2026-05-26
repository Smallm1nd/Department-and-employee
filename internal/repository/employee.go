package repository

import (
	"github.com/Smallm1nd/Department-and-employee/internal/models"
	"gorm.io/gorm"
)

type EmployeeRepo struct {
	conn *gorm.DB
}

func NewEmployeeRepo(conn *gorm.DB) *EmployeeRepo {
	return &EmployeeRepo{conn: conn}
}

func (er *EmployeeRepo) Create(emp *models.Employee) error {
	return er.conn.Create(emp).Error
}
