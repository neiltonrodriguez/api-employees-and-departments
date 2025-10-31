package employee

import (
	"errors"
	"fmt"

	"api-employees-and-departments/internal/domain/logging"
	"api-employees-and-departments/internal/domain/validators"

	"github.com/google/uuid"
)

type Service struct {
	repo   Repository
	logger logging.Logger
}

func NewService(r Repository, logger logging.Logger) *Service {
	return &Service{
		repo:   r,
		logger: logger,
	}
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
		s.logger.Warn("Employee validation failed",
			logging.String("name", emp.Name),
			logging.String("cpf", emp.CPF),
			logging.Error(err),
		)
		return err
	}

	if err := s.repo.Create(emp); err != nil {
		s.logger.Error("Failed to create employee in repository",
			logging.String("employee_id", emp.ID.String()),
			logging.String("name", emp.Name),
			logging.Error(err),
		)
		return err
	}

	s.logger.Info("Employee created successfully",
		logging.String("employee_id", emp.ID.String()),
		logging.String("name", emp.Name),
		logging.String("cpf", emp.CPF),
		logging.String("department_id", emp.DepartmentID.String()),
	)

	return nil
}

func (s *Service) UpdateEmployee(id uuid.UUID, emp *Employee) error {
	if id == uuid.Nil {
		return errors.New("invalid employee id")
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.Error("Employee not found for update",
			logging.String("employee_id", id.String()),
			logging.Error(err),
		)
		return fmt.Errorf("employee not found: %w", err)
	}

	if err := s.validateEmployee(emp); err != nil {
		s.logger.Warn("Employee update validation failed",
			logging.String("employee_id", id.String()),
			logging.Error(err),
		)
		return err
	}

	emp.ID = existing.ID
	if err := s.repo.Update(emp); err != nil {
		s.logger.Error("Failed to update employee in repository",
			logging.String("employee_id", id.String()),
			logging.Error(err),
		)
		return err
	}

	s.logger.Info("Employee updated successfully",
		logging.String("employee_id", id.String()),
		logging.String("name", emp.Name),
	)

	return nil
}

func (s *Service) DeleteEmployee(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid employee id")
	}

	employee, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.Error("Employee not found for deletion",
			logging.String("employee_id", id.String()),
			logging.Error(err),
		)
		return fmt.Errorf("employee not found: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("Failed to delete employee in repository",
			logging.String("employee_id", id.String()),
			logging.Error(err),
		)
		return err
	}

	s.logger.Info("Employee deleted successfully",
		logging.String("employee_id", id.String()),
		logging.String("name", employee.Name),
	)

	return nil
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
