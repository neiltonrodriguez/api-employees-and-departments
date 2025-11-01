package employee

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type MockRepository struct {
	mu             sync.RWMutex
	employees      map[uuid.UUID]*Employee
	findAllError   error
	findByIDError  error
	createError    error
	updateError    error
	deleteError    error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		employees: make(map[uuid.UUID]*Employee),
	}
}

func (m *MockRepository) FindAll() ([]Employee, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.findAllError != nil {
		return nil, m.findAllError
	}

	result := make([]Employee, 0, len(m.employees))
	for _, emp := range m.employees {
		result = append(result, *emp)
	}
	return result, nil
}

func (m *MockRepository) FindByID(id uuid.UUID) (*Employee, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.findByIDError != nil {
		return nil, m.findByIDError
	}

	emp, exists := m.employees[id]
	if !exists {
		return nil, errors.New("employee not found")
	}
	return emp, nil
}

func (m *MockRepository) FindByIDWithManager(id uuid.UUID) (*EmployeeWithManager, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	emp, exists := m.employees[id]
	if !exists {
		return nil, errors.New("employee not found")
	}

	return &EmployeeWithManager{
		Employee:    *emp,
		ManagerName: "Mock Manager",
	}, nil
}

func (m *MockRepository) FindByDepartmentIDs(departmentIDs []uuid.UUID) ([]Employee, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Employee, 0)
	for _, emp := range m.employees {
		for _, deptID := range departmentIDs {
			if emp.DepartmentID == deptID {
				result = append(result, *emp)
				break
			}
		}
	}
	return result, nil
}

func (m *MockRepository) FindWithFilters(filters ListFilters) ([]Employee, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Employee, 0)
	for _, emp := range m.employees {
		match := true
		if filters.Name != nil && emp.Name != *filters.Name {
			match = false
		}
		if filters.CPF != nil && emp.CPF != *filters.CPF {
			match = false
		}
		if filters.DepartmentID != nil && emp.DepartmentID != *filters.DepartmentID {
			match = false
		}
		if match {
			result = append(result, *emp)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockRepository) Create(emp *Employee) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.createError != nil {
		return m.createError
	}

	if emp.ID == uuid.Nil {
		emp.ID = uuid.New()
	}

	m.employees[emp.ID] = emp
	return nil
}

func (m *MockRepository) Update(emp *Employee) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.updateError != nil {
		return m.updateError
	}

	if _, exists := m.employees[emp.ID]; !exists {
		return errors.New("employee not found")
	}

	m.employees[emp.ID] = emp
	return nil
}

func (m *MockRepository) Delete(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.deleteError != nil {
		return m.deleteError
	}

	if _, exists := m.employees[id]; !exists {
		return errors.New("employee not found")
	}

	delete(m.employees, id)
	return nil
}

func (m *MockRepository) SetFindAllError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.findAllError = err
}

func (m *MockRepository) SetFindByIDError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.findByIDError = err
}

func (m *MockRepository) SetCreateError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createError = err
}

func (m *MockRepository) SetUpdateError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateError = err
}

func (m *MockRepository) SetDeleteError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deleteError = err
}

func (m *MockRepository) AddEmployee(emp *Employee) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if emp.ID == uuid.Nil {
		emp.ID = uuid.New()
	}
	m.employees[emp.ID] = emp
}

func (m *MockRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.employees = make(map[uuid.UUID]*Employee)
	m.findAllError = nil
	m.findByIDError = nil
	m.createError = nil
	m.updateError = nil
	m.deleteError = nil
}
