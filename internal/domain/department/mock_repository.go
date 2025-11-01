package department

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type MockRepository struct {
	mu                 sync.RWMutex
	departments        map[uuid.UUID]*Department
	findAllError       error
	findByIDError      error
	findHierarchyError error
	createError        error
	updateError        error
	deleteError        error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		departments: make(map[uuid.UUID]*Department),
	}
}

func (m *MockRepository) FindAll() ([]Department, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.findAllError != nil {
		return nil, m.findAllError
	}

	result := make([]Department, 0, len(m.departments))
	for _, dept := range m.departments {
		result = append(result, *dept)
	}
	return result, nil
}

func (m *MockRepository) FindByID(id uuid.UUID) (*Department, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.findByIDError != nil {
		return nil, m.findByIDError
	}

	dept, exists := m.departments[id]
	if !exists {
		return nil, errors.New("department not found")
	}
	return dept, nil
}

func (m *MockRepository) FindByManagerID(managerID uuid.UUID) ([]Department, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Department, 0)
	for _, dept := range m.departments {
		if dept.ManagerID == managerID {
			result = append(result, *dept)
		}
	}
	return result, nil
}

func (m *MockRepository) FindByParentID(parentID uuid.UUID) ([]Department, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Department, 0)
	for _, dept := range m.departments {
		if dept.ParentDepartmentID != nil && *dept.ParentDepartmentID == parentID {
			result = append(result, *dept)
		}
	}
	return result, nil
}

func (m *MockRepository) FindHierarchyByID(id uuid.UUID) (*DepartmentWithHierarchy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.findHierarchyError != nil {
		return nil, m.findHierarchyError
	}

	dept, exists := m.departments[id]
	if !exists {
		return nil, errors.New("department not found")
	}

	return &DepartmentWithHierarchy{
		Department:     *dept,
		ManagerName:    "Mock Manager",
		Subdepartments: []DepartmentWithHierarchy{},
	}, nil
}

func (m *MockRepository) FindWithFilters(filters ListFilters) ([]Department, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Department, 0)
	for _, dept := range m.departments {
		match := true
		if filters.Name != nil && dept.Name != *filters.Name {
			match = false
		}
		if filters.ParentDepartmentID != nil {
			if dept.ParentDepartmentID == nil || *dept.ParentDepartmentID != *filters.ParentDepartmentID {
				match = false
			}
		}
		if match {
			result = append(result, *dept)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockRepository) Create(dept *Department) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.createError != nil {
		return m.createError
	}

	if dept.ID == uuid.Nil {
		dept.ID = uuid.New()
	}

	m.departments[dept.ID] = dept
	return nil
}

func (m *MockRepository) Update(dept *Department) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.updateError != nil {
		return m.updateError
	}

	if _, exists := m.departments[dept.ID]; !exists {
		return errors.New("department not found")
	}

	m.departments[dept.ID] = dept
	return nil
}

func (m *MockRepository) Delete(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.deleteError != nil {
		return m.deleteError
	}

	if _, exists := m.departments[id]; !exists {
		return errors.New("department not found")
	}

	delete(m.departments, id)
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

func (m *MockRepository) SetFindHierarchyError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.findHierarchyError = err
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

func (m *MockRepository) AddDepartment(dept *Department) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if dept.ID == uuid.Nil {
		dept.ID = uuid.New()
	}
	m.departments[dept.ID] = dept
}

func (m *MockRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.departments = make(map[uuid.UUID]*Department)
	m.findAllError = nil
	m.findByIDError = nil
	m.findHierarchyError = nil
	m.createError = nil
	m.updateError = nil
	m.deleteError = nil
}
