package models

import "time"

type Department struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;size:200"`
	ParentID  *int      `json:"parent_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	Employees []Employee    `json:"employees,omitempty" gorm:"foreignKey:DepartmentID;constraint:OnDelete:CASCADE;"`
	Children  []*Department `json:"children,omitempty" gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE;"`
}

type Employee struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	DepartmentID int        `json:"department_id" gorm:"not null"`
	FullName     string     `json:"full_name" gorm:"not null;size:200"`
	Position     string     `json:"position" gorm:"not null;size:200"`
	HiredAt      *time.Time `json:"hired_at"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
}
