package department

import "github.com/google/uuid"

type ListFilters struct {
	Name               *string
	ManagerName        *string
	ParentDepartmentID *uuid.UUID
	Page               int
	PageSize           int
}

// DepartmentWithHierarchy represents a department with its full hierarchical structure
type DepartmentWithHierarchy struct {
	Department
	ManagerName    string
	Subdepartments []DepartmentWithHierarchy
}

type Repository interface {
	FindAll() ([]Department, error)
	FindByID(id uuid.UUID) (*Department, error)
	FindByManagerID(managerID uuid.UUID) ([]Department, error)
	FindByParentID(parentID uuid.UUID) ([]Department, error)
	FindHierarchyByID(id uuid.UUID) (*DepartmentWithHierarchy, error)
	FindWithFilters(filters ListFilters) ([]Department, int64, error)
	Create(dept *Department) error
	Update(dept *Department) error
	Delete(id uuid.UUID) error
}
