package employee

import "github.com/google/uuid"

type ListFilters struct {
	Name         *string
	CPF          *string
	RG           *string
	DepartmentID *uuid.UUID
	Page         int
	PageSize     int
}

type EmployeeWithManager struct {
	Employee
	ManagerName string
}

type Repository interface {
	FindAll() ([]Employee, error)
	FindByID(id uuid.UUID) (*Employee, error)
	FindByIDWithManager(id uuid.UUID) (*EmployeeWithManager, error)
	FindByDepartmentIDs(departmentIDs []uuid.UUID) ([]Employee, error)
	FindWithFilters(filters ListFilters) ([]Employee, int64, error)
	Create(emp *Employee) error
	Update(emp *Employee) error
	Delete(id uuid.UUID) error
}
