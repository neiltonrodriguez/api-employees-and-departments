package persistence

import (
	"api-employees-and-departments/internal/domain/department"

	"github.com/google/uuid"
)

// EmployeeAdapter adapts EmployeeRepository to department.EmployeeRepository interface
type EmployeeAdapter struct {
	repo *EmployeeRepository
}

func NewEmployeeAdapter(repo *EmployeeRepository) department.EmployeeRepository {
	return &EmployeeAdapter{repo: repo}
}

func (a *EmployeeAdapter) FindByID(id uuid.UUID) (department.Employee, error) {
	emp, err := a.repo.FindByID(id)
	if err != nil {
		return department.Employee{}, err
	}

	return department.Employee{
		ID:           emp.ID,
		Name:         emp.Name,
		DepartmentID: emp.DepartmentID,
	}, nil
}
