package persistence

import (
	"api-employees-and-departments/internal/domain/employee"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) employee.Repository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) FindAll() ([]employee.Employee, error) {
	var employees []employee.Employee
	err := r.db.Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepository) FindByID(id uuid.UUID) (*employee.Employee, error) {
	var emp employee.Employee
	err := r.db.First(&emp, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepository) FindByIDWithManager(id uuid.UUID) (*employee.EmployeeWithManager, error) {
	var result struct {
		employee.Employee
		ManagerName string
	}

	err := r.db.Table("employees AS e").
		Select("e.*, m.name AS manager_name").
		Joins("INNER JOIN departments AS d ON e.department_id = d.id").
		Joins("INNER JOIN employees AS m ON d.manager_id = m.id").
		Where("e.id = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &employee.EmployeeWithManager{
		Employee:    result.Employee,
		ManagerName: result.ManagerName,
	}, nil
}

func (r *EmployeeRepository) FindByDepartmentIDs(departmentIDs []uuid.UUID) ([]employee.Employee, error) {
	var employees []employee.Employee
	if len(departmentIDs) == 0 {
		return employees, nil
	}
	err := r.db.Where("department_id IN ?", departmentIDs).Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepository) Create(emp *employee.Employee) error {
	return r.db.Create(emp).Error
}

func (r *EmployeeRepository) Update(emp *employee.Employee) error {
	return r.db.Save(emp).Error
}

func (r *EmployeeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&employee.Employee{}, "id = ?", id).Error
}

func (r *EmployeeRepository) FindWithFilters(filters employee.ListFilters) ([]employee.Employee, int64, error) {
	var employees []employee.Employee
	var total int64

	query := r.db.Model(&employee.Employee{})

	// Apply filters
	if filters.Name != nil && *filters.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*filters.Name+"%")
	}
	if filters.CPF != nil && *filters.CPF != "" {
		query = query.Where("cpf = ?", *filters.CPF)
	}
	if filters.RG != nil && *filters.RG != "" {
		query = query.Where("rg = ?", *filters.RG)
	}
	if filters.DepartmentID != nil {
		query = query.Where("department_id = ?", *filters.DepartmentID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filters.Page - 1) * filters.PageSize
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}
