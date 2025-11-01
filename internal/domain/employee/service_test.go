package employee

import (
	"errors"
	"testing"

	"api-employees-and-departments/internal/domain/logging"

	"github.com/google/uuid"
)

func TestNewService(t *testing.T) {
	repo := NewMockRepository()
	logger := logging.NewMockLogger()

	service := NewService(repo, logger)

	if service == nil {
		t.Error("NewService() returned nil")
	}
	if service.repo == nil {
		t.Error("Service repository is nil")
	}
	if service.logger == nil {
		t.Error("Service logger is nil")
	}
}

func TestGetAllEmployees(t *testing.T) {
	repo := NewMockRepository()
	logger := logging.NewMockLogger()
	service := NewService(repo, logger)

	deptID := uuid.New()
	emp1 := &Employee{ID: uuid.New(), Name: "John Doe", CPF: "12345678909", DepartmentID: deptID}
	emp2 := &Employee{ID: uuid.New(), Name: "Jane Doe", CPF: "11144477735", DepartmentID: deptID}

	repo.AddEmployee(emp1)
	repo.AddEmployee(emp2)

	employees, err := service.GetAllEmployees()

	if err != nil {
		t.Errorf("GetAllEmployees() returned error: %v", err)
	}
	if len(employees) != 2 {
		t.Errorf("GetAllEmployees() returned %d employees, expected 2", len(employees))
	}
}

func TestGetEmployeeByID(t *testing.T) {
	repo := NewMockRepository()
	logger := logging.NewMockLogger()
	service := NewService(repo, logger)

	deptID := uuid.New()
	emp := &Employee{ID: uuid.New(), Name: "John Doe", CPF: "12345678909", DepartmentID: deptID}
	repo.AddEmployee(emp)

	t.Run("valid ID", func(t *testing.T) {
		result, err := service.GetEmployeeByID(emp.ID)
		if err != nil {
			t.Errorf("GetEmployeeByID() returned error: %v", err)
		}
		if result.ID != emp.ID {
			t.Errorf("GetEmployeeByID() returned wrong employee")
		}
	})

	t.Run("nil ID", func(t *testing.T) {
		_, err := service.GetEmployeeByID(uuid.Nil)
		if err == nil {
			t.Error("GetEmployeeByID() should return error for nil ID")
		}
	})

	t.Run("non-existent ID", func(t *testing.T) {
		_, err := service.GetEmployeeByID(uuid.New())
		if err == nil {
			t.Error("GetEmployeeByID() should return error for non-existent ID")
		}
	})
}

func TestCreateEmployee(t *testing.T) {
	t.Run("valid employee", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		emp := &Employee{
			Name:         "John Doe",
			CPF:          "12345678909",
			DepartmentID: uuid.New(),
		}

		err := service.CreateEmployee(emp)

		if err != nil {
			t.Errorf("CreateEmployee() returned error: %v", err)
		}
		if emp.ID == uuid.Nil {
			t.Error("CreateEmployee() did not set employee ID")
		}
		if logger.CountByLevel("INFO") == 0 {
			t.Error("CreateEmployee() did not log success")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		emp := &Employee{
			CPF:          "12345678909",
			DepartmentID: uuid.New(),
		}

		err := service.CreateEmployee(emp)

		if err == nil {
			t.Error("CreateEmployee() should return error for missing name")
		}
		if logger.CountByLevel("WARN") == 0 {
			t.Error("CreateEmployee() did not log validation warning")
		}
	})

	t.Run("missing CPF", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		emp := &Employee{
			Name:         "John Doe",
			DepartmentID: uuid.New(),
		}

		err := service.CreateEmployee(emp)

		if err == nil {
			t.Error("CreateEmployee() should return error for missing CPF")
		}
	})

	t.Run("invalid CPF", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		emp := &Employee{
			Name:         "John Doe",
			CPF:          "11111111111",
			DepartmentID: uuid.New(),
		}

		err := service.CreateEmployee(emp)

		if err == nil {
			t.Error("CreateEmployee() should return error for invalid CPF")
		}
	})

	t.Run("missing department", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		emp := &Employee{
			Name: "John Doe",
			CPF:  "12345678909",
		}

		err := service.CreateEmployee(emp)

		if err == nil {
			t.Error("CreateEmployee() should return error for missing department")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		repo.SetCreateError(errors.New("database error"))

		emp := &Employee{
			Name:         "John Doe",
			CPF:          "12345678909",
			DepartmentID: uuid.New(),
		}

		err := service.CreateEmployee(emp)

		if err == nil {
			t.Error("CreateEmployee() should return error when repository fails")
		}
		if logger.CountByLevel("ERROR") == 0 {
			t.Error("CreateEmployee() did not log repository error")
		}
	})
}

func TestUpdateEmployee(t *testing.T) {
	t.Run("valid update", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		deptID := uuid.New()
		emp := &Employee{ID: uuid.New(), Name: "John Doe", CPF: "12345678909", DepartmentID: deptID}
		repo.AddEmployee(emp)

		updatedEmp := &Employee{
			Name:         "John Updated",
			CPF:          "12345678909",
			DepartmentID: deptID,
		}

		err := service.UpdateEmployee(emp.ID, updatedEmp)

		if err != nil {
			t.Errorf("UpdateEmployee() returned error: %v", err)
		}
		if logger.CountByLevel("INFO") == 0 {
			t.Error("UpdateEmployee() did not log success")
		}
	})

	t.Run("nil ID", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		emp := &Employee{Name: "John Doe", CPF: "12345678909", DepartmentID: uuid.New()}

		err := service.UpdateEmployee(uuid.Nil, emp)

		if err == nil {
			t.Error("UpdateEmployee() should return error for nil ID")
		}
	})

	t.Run("non-existent employee", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		emp := &Employee{Name: "John Doe", CPF: "12345678909", DepartmentID: uuid.New()}

		err := service.UpdateEmployee(uuid.New(), emp)

		if err == nil {
			t.Error("UpdateEmployee() should return error for non-existent employee")
		}
		if logger.CountByLevel("ERROR") == 0 {
			t.Error("UpdateEmployee() did not log error")
		}
	})

	t.Run("validation error", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		deptID := uuid.New()
		emp := &Employee{ID: uuid.New(), Name: "John Doe", CPF: "12345678909", DepartmentID: deptID}
		repo.AddEmployee(emp)

		updatedEmp := &Employee{CPF: "12345678909", DepartmentID: deptID}

		err := service.UpdateEmployee(emp.ID, updatedEmp)

		if err == nil {
			t.Error("UpdateEmployee() should return error for invalid employee")
		}
		if logger.CountByLevel("WARN") == 0 {
			t.Error("UpdateEmployee() did not log validation warning")
		}
	})
}

func TestDeleteEmployee(t *testing.T) {
	t.Run("valid deletion", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		deptID := uuid.New()
		emp := &Employee{ID: uuid.New(), Name: "John Doe", CPF: "12345678909", DepartmentID: deptID}
		repo.AddEmployee(emp)

		err := service.DeleteEmployee(emp.ID)

		if err != nil {
			t.Errorf("DeleteEmployee() returned error: %v", err)
		}
		if logger.CountByLevel("INFO") == 0 {
			t.Error("DeleteEmployee() did not log success")
		}
	})

	t.Run("nil ID", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		err := service.DeleteEmployee(uuid.Nil)

		if err == nil {
			t.Error("DeleteEmployee() should return error for nil ID")
		}
	})

	t.Run("non-existent employee", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		err := service.DeleteEmployee(uuid.New())

		if err == nil {
			t.Error("DeleteEmployee() should return error for non-existent employee")
		}
		if logger.CountByLevel("ERROR") == 0 {
			t.Error("DeleteEmployee() did not log error")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repo := NewMockRepository()
		logger := logging.NewMockLogger()
		service := NewService(repo, logger)

		deptID := uuid.New()
		emp := &Employee{ID: uuid.New(), Name: "John Doe", CPF: "12345678909", DepartmentID: deptID}
		repo.AddEmployee(emp)
		repo.SetDeleteError(errors.New("database error"))

		err := service.DeleteEmployee(emp.ID)

		if err == nil {
			t.Error("DeleteEmployee() should return error when repository fails")
		}
		if logger.CountByLevel("ERROR") == 0 {
			t.Error("DeleteEmployee() did not log repository error")
		}
	})
}

func TestGetEmployeeWithManager(t *testing.T) {
	repo := NewMockRepository()
	logger := logging.NewMockLogger()
	service := NewService(repo, logger)

	deptID := uuid.New()
	emp := &Employee{ID: uuid.New(), Name: "John Doe", CPF: "12345678909", DepartmentID: deptID}
	repo.AddEmployee(emp)

	t.Run("valid ID", func(t *testing.T) {
		result, err := service.GetEmployeeWithManager(emp.ID)
		if err != nil {
			t.Errorf("GetEmployeeWithManager() returned error: %v", err)
		}
		if result.ID != emp.ID {
			t.Errorf("GetEmployeeWithManager() returned wrong employee")
		}
	})

	t.Run("nil ID", func(t *testing.T) {
		_, err := service.GetEmployeeWithManager(uuid.Nil)
		if err == nil {
			t.Error("GetEmployeeWithManager() should return error for nil ID")
		}
	})
}

func TestGetEmployeesByDepartmentIDs(t *testing.T) {
	repo := NewMockRepository()
	logger := logging.NewMockLogger()
	service := NewService(repo, logger)

	deptID1 := uuid.New()
	deptID2 := uuid.New()

	emp1 := &Employee{ID: uuid.New(), Name: "John Doe", CPF: "12345678909", DepartmentID: deptID1}
	emp2 := &Employee{ID: uuid.New(), Name: "Jane Doe", CPF: "11144477735", DepartmentID: deptID2}
	emp3 := &Employee{ID: uuid.New(), Name: "Bob Smith", CPF: "98765432100", DepartmentID: deptID1}

	repo.AddEmployee(emp1)
	repo.AddEmployee(emp2)
	repo.AddEmployee(emp3)

	employees, err := service.GetEmployeesByDepartmentIDs([]uuid.UUID{deptID1})

	if err != nil {
		t.Errorf("GetEmployeesByDepartmentIDs() returned error: %v", err)
	}
	if len(employees) != 2 {
		t.Errorf("GetEmployeesByDepartmentIDs() returned %d employees, expected 2", len(employees))
	}
}

func TestListEmployees(t *testing.T) {
	repo := NewMockRepository()
	logger := logging.NewMockLogger()
	service := NewService(repo, logger)

	deptID := uuid.New()
	emp1 := &Employee{ID: uuid.New(), Name: "John Doe", CPF: "12345678909", DepartmentID: deptID}
	emp2 := &Employee{ID: uuid.New(), Name: "Jane Doe", CPF: "11144477735", DepartmentID: deptID}

	repo.AddEmployee(emp1)
	repo.AddEmployee(emp2)

	name := "John Doe"
	filters := ListFilters{
		Name: &name,
	}

	employees, total, err := service.ListEmployees(filters)

	if err != nil {
		t.Errorf("ListEmployees() returned error: %v", err)
	}
	if len(employees) != 1 {
		t.Errorf("ListEmployees() returned %d employees, expected 1", len(employees))
	}
	if total != 1 {
		t.Errorf("ListEmployees() returned total %d, expected 1", total)
	}
}
