package employee

import (
	"errors"
	"fmt"

	"api-employees-and-departments/internal/domain/validators"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetAllEmployees() ([]Employee, error) {
	return s.repo.FindAll()
}

func (s *Service) ListEmployees(filters ListFilters) ([]Employee, int64, error) {
	return s.repo.FindWithFilters(filters)
}

func (s *Service) GetEmployeeByID(id uuid.UUID) (*Employee, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid employee id")
	}
	return s.repo.FindByID(id)
}

func (s *Service) GetEmployeeWithManager(id uuid.UUID) (*EmployeeWithManager, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid employee id")
	}
	return s.repo.FindByIDWithManager(id)
}

func (s *Service) GetEmployeesByDepartmentIDs(departmentIDs []uuid.UUID) ([]Employee, error) {
	return s.repo.FindByDepartmentIDs(departmentIDs)
}

func (s *Service) CreateEmployee(emp *Employee) error {
	if err := s.validateEmployee(emp); err != nil {
		return err
	}
	return s.repo.Create(emp)
}

func (s *Service) UpdateEmployee(id uuid.UUID, emp *Employee) error {
	if id == uuid.Nil {
		return errors.New("invalid employee id")
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("employee not found: %w", err)
	}

	if err := s.validateEmployee(emp); err != nil {
		return err
	}

	emp.ID = existing.ID
	return s.repo.Update(emp)
}

func (s *Service) DeleteEmployee(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid employee id")
	}

	_, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("employee not found: %w", err)
	}

	return s.repo.Delete(id)
}

func (s *Service) validateEmployee(emp *Employee) error {
	if emp.Name == "" {
		return errors.New("employee name is required")
	}
	if emp.CPF == "" {
		return errors.New("employee CPF is required")
	}
	if !validators.ValidateCPF(emp.CPF) {
		return errors.New("invalid CPF")
	}
	if emp.DepartmentID == uuid.Nil {
		return errors.New("employee department is required")
	}
	return nil
}
