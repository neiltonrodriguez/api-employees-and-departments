package department

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type MockEmployeeRepository struct {
	mu           sync.RWMutex
	employees    map[uuid.UUID]*Employee
	findByIDError error
}

func NewMockEmployeeRepository() *MockEmployeeRepository {
	return &MockEmployeeRepository{
		employees: make(map[uuid.UUID]*Employee),
	}
}

func (m *MockEmployeeRepository) FindByID(id uuid.UUID) (Employee, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.findByIDError != nil {
		return Employee{}, m.findByIDError
	}

	emp, exists := m.employees[id]
	if !exists {
		return Employee{}, errors.New("employee not found")
	}
	return *emp, nil
}

func (m *MockEmployeeRepository) SetFindByIDError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.findByIDError = err
}

func (m *MockEmployeeRepository) AddEmployee(emp *Employee) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if emp.ID == uuid.Nil {
		emp.ID = uuid.New()
	}
	m.employees[emp.ID] = emp
}

func (m *MockEmployeeRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.employees = make(map[uuid.UUID]*Employee)
	m.findByIDError = nil
}
