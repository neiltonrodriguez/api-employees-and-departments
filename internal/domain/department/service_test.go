package department

import (
	"errors"
	"testing"
	"time"

	"api-employees-and-departments/internal/domain/cache"
	"api-employees-and-departments/internal/domain/logging"

	"github.com/google/uuid"
)

func TestNewService(t *testing.T) {
	repo := NewMockRepository()
	empRepo := NewMockEmployeeRepository()
	logger := logging.NewMockLogger()
	mockCache := cache.NewMockCache()
	cacheTTL := 5 * time.Minute

	service := NewService(repo, empRepo, logger, mockCache, cacheTTL)

	if service == nil {
		t.Error("NewService() returned nil")
	}
	if service.repo == nil {
		t.Error("Service repository is nil")
	}
	if service.logger == nil {
		t.Error("Service logger is nil")
	}
	if service.cache == nil {
		t.Error("Service cache is nil")
	}
}

func TestGetAllDepartments(t *testing.T) {
	repo := NewMockRepository()
	empRepo := NewMockEmployeeRepository()
	logger := logging.NewMockLogger()
	mockCache := cache.NewMockCache()
	service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

	managerID := uuid.New()
	dept1 := &Department{ID: uuid.New(), Name: "IT", ManagerID: managerID}
	dept2 := &Department{ID: uuid.New(), Name: "HR", ManagerID: managerID}

	repo.AddDepartment(dept1)
	repo.AddDepartment(dept2)

	departments, err := service.GetAllDepartments()

	if err != nil {
		t.Errorf("GetAllDepartments() returned error: %v", err)
	}
	if len(departments) != 2 {
		t.Errorf("GetAllDepartments() returned %d departments, expected 2", len(departments))
	}
}

func TestGetDepartmentByID(t *testing.T) {
	repo := NewMockRepository()
	empRepo := NewMockEmployeeRepository()
	logger := logging.NewMockLogger()
	mockCache := cache.NewMockCache()
	service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

	managerID := uuid.New()
	dept := &Department{ID: uuid.New(), Name: "IT", ManagerID: managerID}
	repo.AddDepartment(dept)

	t.Run("valid ID", func(t *testing.T) {
		result, err := service.GetDepartmentByID(dept.ID)
		if err != nil {
			t.Errorf("GetDepartmentByID() returned error: %v", err)
		}
		if result.ID != dept.ID {
			t.Errorf("GetDepartmentByID() returned wrong department")
		}
	})

	t.Run("nil ID", func(t *testing.T) {
		_, err := service.GetDepartmentByID(uuid.Nil)
		if err == nil {
			t.Error("GetDepartmentByID() should return error for nil ID")
		}
	})

	t.Run("non-existent ID", func(t *testing.T) {
		_, err := service.GetDepartmentByID(uuid.New())
		if err == nil {
			t.Error("GetDepartmentByID() should return error for non-existent ID")
		}
	})
}

func TestGetDepartmentWithHierarchy(t *testing.T) {
	repo := NewMockRepository()
	empRepo := NewMockEmployeeRepository()
	logger := logging.NewMockLogger()
	mockCache := cache.NewMockCache()
	service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

	managerID := uuid.New()
	dept := &Department{ID: uuid.New(), Name: "IT", ManagerID: managerID}
	repo.AddDepartment(dept)

	t.Run("cache miss - fetch from database", func(t *testing.T) {
		mockCache.Reset()
		result, err := service.GetDepartmentWithHierarchy(dept.ID)
		if err != nil {
			t.Errorf("GetDepartmentWithHierarchy() returned error: %v", err)
		}
		if result.ID != dept.ID {
			t.Errorf("GetDepartmentWithHierarchy() returned wrong department")
		}
		if mockCache.SetCalls == 0 {
			t.Error("GetDepartmentWithHierarchy() did not cache result")
		}
	})

	t.Run("cache hit", func(t *testing.T) {
		mockCache.Reset()
		service.GetDepartmentWithHierarchy(dept.ID)
		getCalls := mockCache.GetCalls
		service.GetDepartmentWithHierarchy(dept.ID)
		if mockCache.GetCalls <= getCalls {
			t.Error("GetDepartmentWithHierarchy() did not use cache")
		}
	})

	t.Run("nil ID", func(t *testing.T) {
		_, err := service.GetDepartmentWithHierarchy(uuid.Nil)
		if err == nil {
			t.Error("GetDepartmentWithHierarchy() should return error for nil ID")
		}
	})
}

func TestCreateDepartment(t *testing.T) {
	t.Run("valid department", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		dept := &Department{
			Name:      "IT",
			ManagerID: uuid.New(),
		}

		err := service.CreateDepartment(dept)

		if err != nil {
			t.Errorf("CreateDepartment() returned error: %v", err)
		}
		if dept.ID == uuid.Nil {
			t.Error("CreateDepartment() did not set department ID")
		}
		if logger.CountByLevel("INFO") == 0 {
			t.Error("CreateDepartment() did not log success")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		dept := &Department{
			ManagerID: uuid.New(),
		}

		err := service.CreateDepartment(dept)

		if err == nil {
			t.Error("CreateDepartment() should return error for missing name")
		}
		if logger.CountByLevel("WARN") == 0 {
			t.Error("CreateDepartment() did not log validation warning")
		}
	})

	t.Run("missing manager", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		dept := &Department{
			Name: "IT",
		}

		err := service.CreateDepartment(dept)

		if err == nil {
			t.Error("CreateDepartment() should return error for missing manager")
		}
	})

	t.Run("with valid parent department", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		parentDept := &Department{
			ID:        uuid.New(),
			Name:      "Parent",
			ManagerID: uuid.New(),
		}
		repo.AddDepartment(parentDept)

		dept := &Department{
			Name:               "Child",
			ManagerID:          uuid.New(),
			ParentDepartmentID: &parentDept.ID,
		}

		err := service.CreateDepartment(dept)

		if err != nil {
			t.Errorf("CreateDepartment() returned error: %v", err)
		}
	})

	t.Run("with non-existent parent department", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		nonExistentID := uuid.New()
		dept := &Department{
			Name:               "Child",
			ManagerID:          uuid.New(),
			ParentDepartmentID: &nonExistentID,
		}

		err := service.CreateDepartment(dept)

		if err == nil {
			t.Error("CreateDepartment() should return error for non-existent parent")
		}
		if logger.CountByLevel("ERROR") == 0 {
			t.Error("CreateDepartment() did not log error")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		repo.SetCreateError(errors.New("database error"))

		dept := &Department{
			Name:      "IT",
			ManagerID: uuid.New(),
		}

		err := service.CreateDepartment(dept)

		if err == nil {
			t.Error("CreateDepartment() should return error when repository fails")
		}
		if logger.CountByLevel("ERROR") == 0 {
			t.Error("CreateDepartment() did not log repository error")
		}
	})
}

func TestUpdateDepartment(t *testing.T) {
	t.Run("valid update", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		deptID := uuid.New()
		managerID := uuid.New()

		dept := &Department{ID: deptID, Name: "IT", ManagerID: managerID}
		repo.AddDepartment(dept)

		emp := &Employee{ID: managerID, Name: "Manager", DepartmentID: deptID}
		empRepo.AddEmployee(emp)

		updatedDept := &Department{
			Name:      "IT Updated",
			ManagerID: managerID,
		}

		err := service.UpdateDepartment(deptID, updatedDept)

		if err != nil {
			t.Errorf("UpdateDepartment() returned error: %v", err)
		}
		if logger.CountByLevel("INFO") == 0 {
			t.Error("UpdateDepartment() did not log success")
		}
		if mockCache.DeleteCalls == 0 {
			t.Error("UpdateDepartment() did not invalidate cache")
		}
	})

	t.Run("nil ID", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		dept := &Department{Name: "IT", ManagerID: uuid.New()}

		err := service.UpdateDepartment(uuid.Nil, dept)

		if err == nil {
			t.Error("UpdateDepartment() should return error for nil ID")
		}
	})

	t.Run("non-existent department", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		dept := &Department{Name: "IT", ManagerID: uuid.New()}

		err := service.UpdateDepartment(uuid.New(), dept)

		if err == nil {
			t.Error("UpdateDepartment() should return error for non-existent department")
		}
		if logger.CountByLevel("ERROR") == 0 {
			t.Error("UpdateDepartment() did not log error")
		}
	})

	t.Run("manager not in department", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		deptID := uuid.New()
		managerID := uuid.New()

		dept := &Department{ID: deptID, Name: "IT", ManagerID: managerID}
		repo.AddDepartment(dept)

		emp := &Employee{ID: managerID, Name: "Manager", DepartmentID: uuid.New()}
		empRepo.AddEmployee(emp)

		updatedDept := &Department{
			Name:      "IT Updated",
			ManagerID: managerID,
		}

		err := service.UpdateDepartment(deptID, updatedDept)

		if err == nil {
			t.Error("UpdateDepartment() should return error when manager not in department")
		}
	})
}

func TestDeleteDepartment(t *testing.T) {
	t.Run("valid deletion", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		managerID := uuid.New()
		dept := &Department{ID: uuid.New(), Name: "IT", ManagerID: managerID}
		repo.AddDepartment(dept)

		err := service.DeleteDepartment(dept.ID)

		if err != nil {
			t.Errorf("DeleteDepartment() returned error: %v", err)
		}
		if logger.CountByLevel("INFO") == 0 {
			t.Error("DeleteDepartment() did not log success")
		}
		if mockCache.DeleteCalls == 0 {
			t.Error("DeleteDepartment() did not invalidate cache")
		}
	})

	t.Run("nil ID", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		err := service.DeleteDepartment(uuid.Nil)

		if err == nil {
			t.Error("DeleteDepartment() should return error for nil ID")
		}
	})

	t.Run("non-existent department", func(t *testing.T) {
		repo := NewMockRepository()
		empRepo := NewMockEmployeeRepository()
		logger := logging.NewMockLogger()
		mockCache := cache.NewMockCache()
		service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

		err := service.DeleteDepartment(uuid.New())

		if err == nil {
			t.Error("DeleteDepartment() should return error for non-existent department")
		}
		if logger.CountByLevel("ERROR") == 0 {
			t.Error("DeleteDepartment() did not log error")
		}
	})
}

func TestValidateNoCycle(t *testing.T) {
	repo := NewMockRepository()
	empRepo := NewMockEmployeeRepository()
	logger := logging.NewMockLogger()
	mockCache := cache.NewMockCache()
	service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

	t.Run("no parent - no cycle", func(t *testing.T) {
		err := service.validateNoCycle(uuid.New(), nil)
		if err != nil {
			t.Errorf("validateNoCycle() returned error for no parent: %v", err)
		}
	})

	t.Run("self parent - cycle", func(t *testing.T) {
		deptID := uuid.New()
		err := service.validateNoCycle(deptID, &deptID)
		if err == nil {
			t.Error("validateNoCycle() should return error for self parent")
		}
	})

	t.Run("valid hierarchy - no cycle", func(t *testing.T) {
		dept1 := &Department{ID: uuid.New(), Name: "Dept1", ManagerID: uuid.New()}
		dept2 := &Department{ID: uuid.New(), Name: "Dept2", ManagerID: uuid.New(), ParentDepartmentID: &dept1.ID}
		dept3 := &Department{ID: uuid.New(), Name: "Dept3", ManagerID: uuid.New(), ParentDepartmentID: &dept2.ID}

		repo.AddDepartment(dept1)
		repo.AddDepartment(dept2)
		repo.AddDepartment(dept3)

		newParent := dept2.ID
		err := service.validateNoCycle(dept3.ID, &newParent)
		if err != nil {
			t.Errorf("validateNoCycle() returned error for valid hierarchy: %v", err)
		}
	})

	t.Run("cycle detection", func(t *testing.T) {
		repo.Reset()
		dept1 := &Department{ID: uuid.New(), Name: "Dept1", ManagerID: uuid.New()}
		dept2 := &Department{ID: uuid.New(), Name: "Dept2", ManagerID: uuid.New(), ParentDepartmentID: &dept1.ID}

		repo.AddDepartment(dept1)
		repo.AddDepartment(dept2)

		err := service.validateNoCycle(dept1.ID, &dept2.ID)
		if err == nil {
			t.Error("validateNoCycle() should detect cycle")
		}
	})
}

func TestGetDepartmentsByManagerID(t *testing.T) {
	repo := NewMockRepository()
	empRepo := NewMockEmployeeRepository()
	logger := logging.NewMockLogger()
	mockCache := cache.NewMockCache()
	service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

	managerID := uuid.New()
	dept1 := &Department{ID: uuid.New(), Name: "IT", ManagerID: managerID}
	dept2 := &Department{ID: uuid.New(), Name: "HR", ManagerID: uuid.New()}

	repo.AddDepartment(dept1)
	repo.AddDepartment(dept2)

	departments, err := service.GetDepartmentsByManagerID(managerID)

	if err != nil {
		t.Errorf("GetDepartmentsByManagerID() returned error: %v", err)
	}
	if len(departments) != 1 {
		t.Errorf("GetDepartmentsByManagerID() returned %d departments, expected 1", len(departments))
	}
}

func TestGetDepartmentsByParentID(t *testing.T) {
	repo := NewMockRepository()
	empRepo := NewMockEmployeeRepository()
	logger := logging.NewMockLogger()
	mockCache := cache.NewMockCache()
	service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

	parentID := uuid.New()
	parentDept := &Department{ID: parentID, Name: "Parent", ManagerID: uuid.New()}
	childDept := &Department{ID: uuid.New(), Name: "Child", ManagerID: uuid.New(), ParentDepartmentID: &parentID}

	repo.AddDepartment(parentDept)
	repo.AddDepartment(childDept)

	departments, err := service.GetDepartmentsByParentID(parentID)

	if err != nil {
		t.Errorf("GetDepartmentsByParentID() returned error: %v", err)
	}
	if len(departments) != 1 {
		t.Errorf("GetDepartmentsByParentID() returned %d departments, expected 1", len(departments))
	}
}

func TestListDepartments(t *testing.T) {
	repo := NewMockRepository()
	empRepo := NewMockEmployeeRepository()
	logger := logging.NewMockLogger()
	mockCache := cache.NewMockCache()
	service := NewService(repo, empRepo, logger, mockCache, 5*time.Minute)

	managerID := uuid.New()
	dept1 := &Department{ID: uuid.New(), Name: "IT", ManagerID: managerID}
	dept2 := &Department{ID: uuid.New(), Name: "HR", ManagerID: managerID}

	repo.AddDepartment(dept1)
	repo.AddDepartment(dept2)

	name := "IT"
	filters := ListFilters{
		Name: &name,
	}

	departments, total, err := service.ListDepartments(filters)

	if err != nil {
		t.Errorf("ListDepartments() returned error: %v", err)
	}
	if len(departments) != 1 {
		t.Errorf("ListDepartments() returned %d departments, expected 1", len(departments))
	}
	if total != 1 {
		t.Errorf("ListDepartments() returned total %d, expected 1", total)
	}
}
