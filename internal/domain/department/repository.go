package department

import "github.com/google/uuid"

type ListFilters struct {
	Name                 *string
	ManagerName          *string
	ParentDepartmentID   *uuid.UUID
	Page                 int
	PageSize             int
}

type Repository interface {
	FindAll() ([]Department, error)
	FindByID(id uuid.UUID) (*Department, error)
	FindByManagerID(managerID uuid.UUID) ([]Department, error)
	FindByParentID(parentID uuid.UUID) ([]Department, error)
	FindWithFilters(filters ListFilters) ([]Department, int64, error)
	Create(dept *Department) error
	Update(dept *Department) error
	Delete(id uuid.UUID) error
}
