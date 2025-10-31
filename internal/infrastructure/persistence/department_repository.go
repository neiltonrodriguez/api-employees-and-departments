package persistence

import (
	"api-employees-and-departments/internal/domain/department"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

// hierarchyRow represents a row from the CTE recursive query
type hierarchyRow struct {
	ID                   uuid.UUID      `gorm:"column:id"`
	Name                 string         `gorm:"column:name"`
	ManagerID            uuid.UUID      `gorm:"column:manager_id"`
	ParentDepartmentID   *uuid.UUID     `gorm:"column:parent_department_id"`
	ManagerName          string         `gorm:"column:manager_name"`
	Level                int            `gorm:"column:level"`
	Path                 pq.StringArray `gorm:"column:path;type:text[]"`
	CreatedAt            time.Time      `gorm:"column:created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at"`
}

// FindHierarchyByID retrieves department hierarchy using PostgreSQL CTE recursive query
func (r *DepartmentRepository) FindHierarchyByID(id uuid.UUID) (*department.DepartmentWithHierarchy, error) {
	var rows []hierarchyRow

	// CTE recursivo para buscar toda a hierarquia de uma vez
	query := `
	WITH RECURSIVE department_tree AS (
		-- Base case: departamento raiz
		SELECT
			d.id,
			d.name,
			d.manager_id,
			d.parent_department_id,
			e.name as manager_name,
			d.created_at,
			d.updated_at,
			0 as level,
			ARRAY[d.id::text] as path
		FROM departments d
		LEFT JOIN employees e ON d.manager_id = e.id AND e.deleted_at IS NULL
		WHERE d.id = $1
		AND d.deleted_at IS NULL

		UNION ALL

		-- Recursive case: filhos recursivamente
		SELECT
			d.id,
			d.name,
			d.manager_id,
			d.parent_department_id,
			e.name as manager_name,
			d.created_at,
			d.updated_at,
			dt.level + 1,
			dt.path || d.id::text
		FROM departments d
		LEFT JOIN employees e ON d.manager_id = e.id AND e.deleted_at IS NULL
		INNER JOIN department_tree dt ON d.parent_department_id = dt.id
		WHERE d.deleted_at IS NULL
		AND NOT d.id::text = ANY(dt.path)  -- Previne ciclos
	)
	SELECT * FROM department_tree
	ORDER BY level, name
	`

	err := r.db.Raw(query, id).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Construir árvore a partir do resultado flat
	return buildTreeFromFlatList(rows), nil
}

// buildTreeFromFlatList constrói a estrutura hierárquica a partir do resultado flat do CTE
func buildTreeFromFlatList(rows []hierarchyRow) *department.DepartmentWithHierarchy {
	if len(rows) == 0 {
		return nil
	}

	// Criar map para acesso rápido por ID
	deptMap := make(map[uuid.UUID]*department.DepartmentWithHierarchy)

	// Primeira passagem: criar todos os nós
	for _, row := range rows {
		dept := &department.DepartmentWithHierarchy{
			Department: department.Department{
				ID:                   row.ID,
				Name:                 row.Name,
				ManagerID:            row.ManagerID,
				ParentDepartmentID:   row.ParentDepartmentID,
				CreatedAt:            row.CreatedAt,
				UpdatedAt:            row.UpdatedAt,
			},
			ManagerName:    row.ManagerName,
			Subdepartments: []department.DepartmentWithHierarchy{},
		}
		deptMap[row.ID] = dept
	}

	// Segunda passagem: construir as relações parent-child
	var root *department.DepartmentWithHierarchy
	for _, row := range rows {
		dept := deptMap[row.ID]
		if row.ParentDepartmentID != nil {
			if parent, exists := deptMap[*row.ParentDepartmentID]; exists {
				parent.Subdepartments = append(parent.Subdepartments, *dept)
			}
		} else {
			root = dept  // Departamento raiz (nível 0)
		}
	}

	return root
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
