package department

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type EmployeeRepository interface {
	FindByID(id uuid.UUID) (Employee, error)
}

type Employee struct {
	ID           uuid.UUID
	Name         string
	DepartmentID uuid.UUID
}

type Service struct {
	repo        Repository
	employeeRepo EmployeeRepository
}

func NewService(r Repository, empRepo EmployeeRepository) *Service {
	return &Service{
		repo:        r,
		employeeRepo: empRepo,
	}
}

func (s *Service) GetAllDepartments() ([]Department, error) {
	return s.repo.FindAll()
}

func (s *Service) ListDepartments(filters ListFilters) ([]Department, int64, error) {
	return s.repo.FindWithFilters(filters)
}

func (s *Service) GetDepartmentByID(id uuid.UUID) (*Department, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid department id")
	}
	return s.repo.FindByID(id)
}

func (s *Service) GetDepartmentsByManagerID(managerID uuid.UUID) ([]Department, error) {
	return s.repo.FindByManagerID(managerID)
}

func (s *Service) GetDepartmentsByParentID(parentID uuid.UUID) ([]Department, error) {
	return s.repo.FindByParentID(parentID)
}

type DepartmentWithHierarchy struct {
	Department
	ManagerName    string
	Subdepartments []DepartmentWithHierarchy
}

func (s *Service) GetDepartmentWithHierarchy(id uuid.UUID) (*DepartmentWithHierarchy, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid department id")
	}

	// Get the department
	dept, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Get manager name
	manager, err := s.employeeRepo.FindByID(dept.ManagerID)
	if err != nil {
		return nil, errors.New("manager not found")
	}

	// Build hierarchy recursively
	subdepartments, err := s.buildHierarchy(id)
	if err != nil {
		return nil, err
	}

	return &DepartmentWithHierarchy{
		Department:     *dept,
		ManagerName:    manager.Name,
		Subdepartments: subdepartments,
	}, nil
}

func (s *Service) buildHierarchy(parentID uuid.UUID) ([]DepartmentWithHierarchy, error) {
	// Find all direct children
	children, err := s.repo.FindByParentID(parentID)
	if err != nil {
		return nil, err
	}

	result := make([]DepartmentWithHierarchy, 0, len(children))

	for _, child := range children {
		// Get manager name for each child
		manager, err := s.employeeRepo.FindByID(child.ManagerID)
		if err != nil {
			continue // Skip if manager not found
		}

		// Recursively build subdepartments
		subdepartments, err := s.buildHierarchy(child.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, DepartmentWithHierarchy{
			Department:     child,
			ManagerName:    manager.Name,
			Subdepartments: subdepartments,
		})
	}

	return result, nil
}

func (s *Service) CreateDepartment(dept *Department) error {
	if err := s.validateDepartment(dept); err != nil {
		return err
	}

	// Validate no cycles in hierarchy (for create, we use a temporary UUID to simulate the new department)
	if dept.ParentDepartmentID != nil {
		// For new departments, just check if parent exists
		_, err := s.repo.FindByID(*dept.ParentDepartmentID)
		if err != nil {
			return errors.New("parent department not found")
		}
	}

	return s.repo.Create(dept)
}

func (s *Service) UpdateDepartment(id uuid.UUID, dept *Department) error {
	if id == uuid.Nil {
		return errors.New("invalid department id")
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("department not found: %w", err)
	}

	if err := s.validateDepartment(dept); err != nil {
		return err
	}

	// Validate no cycles in hierarchy
	if err := s.validateNoCycle(id, dept.ParentDepartmentID); err != nil {
		return err
	}

	// Validate that manager belongs to the department
	if err := s.validateManagerBelongsToDepartment(dept.ManagerID, id); err != nil {
		return err
	}

	dept.ID = existing.ID
	return s.repo.Update(dept)
}

func (s *Service) DeleteDepartment(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid department id")
	}

	_, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("department not found: %w", err)
	}

	return s.repo.Delete(id)
}

func (s *Service) validateDepartment(dept *Department) error {
	if dept.Name == "" {
		return errors.New("department name is required")
	}
	if dept.ManagerID == uuid.Nil {
		return errors.New("department manager is required")
	}
	return nil
}

func (s *Service) validateManagerBelongsToDepartment(managerID, departmentID uuid.UUID) error {
	// Check if manager exists
	manager, err := s.employeeRepo.FindByID(managerID)
	if err != nil {
		return errors.New("manager not found")
	}

	// Check if manager belongs to the department
	if manager.DepartmentID != departmentID {
		return errors.New("manager must be linked to the same department")
	}

	return nil
}

func (s *Service) validateNoCycle(departmentID uuid.UUID, parentDepartmentID *uuid.UUID) error {
	// If no parent department, no cycle possible
	if parentDepartmentID == nil || *parentDepartmentID == uuid.Nil {
		return nil
	}

	// Check if parent is the same as current department
	if *parentDepartmentID == departmentID {
		return errors.New("department cannot be its own parent")
	}

	// Traverse the hierarchy to detect cycles
	visited := make(map[uuid.UUID]bool)
	currentID := *parentDepartmentID

	for currentID != uuid.Nil {
		// If we've already visited this department, there's a cycle
		if visited[currentID] {
			return errors.New("cycle detected in department hierarchy")
		}

		// If we reached the original department, there's a cycle
		if currentID == departmentID {
			return errors.New("cycle detected in department hierarchy")
		}

		visited[currentID] = true

		// Get the parent of current department
		dept, err := s.repo.FindByID(currentID)
		if err != nil {
			return fmt.Errorf("department not found in hierarchy: %w", err)
		}

		// Move to the next parent
		if dept.ParentDepartmentID == nil {
			break
		}
		currentID = *dept.ParentDepartmentID
	}

	return nil
}
