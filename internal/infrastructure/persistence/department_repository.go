package persistence

import (
	"api-employees-and-departments/internal/domain/department"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) department.Repository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) FindAll() ([]department.Department, error) {
	var departments []department.Department
	err := r.db.Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepository) FindByID(id uuid.UUID) (*department.Department, error) {
	var dept department.Department
	err := r.db.First(&dept, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *DepartmentRepository) FindByParentID(parentID uuid.UUID) ([]department.Department, error) {
	var departments []department.Department
	err := r.db.Where("parent_department_id = ?", parentID).Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepository) FindByManagerID(managerID uuid.UUID) ([]department.Department, error) {
	var departments []department.Department
	err := r.db.Where("manager_id = ?", managerID).Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepository) Create(dept *department.Department) error {
	return r.db.Create(dept).Error
}

func (r *DepartmentRepository) Update(dept *department.Department) error {
	return r.db.Save(dept).Error
}

func (r *DepartmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&department.Department{}, "id = ?", id).Error
}

func (r *DepartmentRepository) FindWithFilters(filters department.ListFilters) ([]department.Department, int64, error) {
	var departments []department.Department
	var total int64

	query := r.db.Model(&department.Department{})

	// Apply filters
	if filters.Name != nil && *filters.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*filters.Name+"%")
	}
	if filters.ParentDepartmentID != nil {
		query = query.Where("parent_department_id = ?", *filters.ParentDepartmentID)
	}
	// TODO: Manager name filter requires join with employees table

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filters.Page - 1) * filters.PageSize
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&departments).Error; err != nil {
		return nil, 0, err
	}

	return departments, total, nil
}
