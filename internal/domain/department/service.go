package department

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"api-employees-and-departments/internal/domain/cache"
	"api-employees-and-departments/internal/domain/logging"

	"github.com/google/uuid"
)

type EmployeeRepository interface {
	FindByID(id uuid.UUID) (Employee, error)
}

type Employee struct {
	ID           uuid.UUID
	Name         string
	DepartmentID uuid.UUID
}

type Service struct {
	repo         Repository
	employeeRepo EmployeeRepository
	logger       logging.Logger
	cache        cache.Cache
	cacheTTL     time.Duration
	cacheKeys    *cache.CacheKeyBuilder
}

func NewService(r Repository, empRepo EmployeeRepository, logger logging.Logger, c cache.Cache, cacheTTL time.Duration) *Service {
	return &Service{
		repo:         r,
		employeeRepo: empRepo,
		logger:       logger,
		cache:        c,
		cacheTTL:     cacheTTL,
		cacheKeys:    cache.NewCacheKeyBuilder("department"),
	}
}

func (s *Service) GetAllDepartments() ([]Department, error) {
	return s.repo.FindAll()
}

func (s *Service) ListDepartments(filters ListFilters) ([]Department, int64, error) {
	return s.repo.FindWithFilters(filters)
}

func (s *Service) GetDepartmentByID(id uuid.UUID) (*Department, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid department id")
	}
	return s.repo.FindByID(id)
}

func (s *Service) GetDepartmentsByManagerID(managerID uuid.UUID) ([]Department, error) {
	return s.repo.FindByManagerID(managerID)
}

func (s *Service) GetDepartmentsByParentID(parentID uuid.UUID) ([]Department, error) {
	return s.repo.FindByParentID(parentID)
}

func (s *Service) GetDepartmentWithHierarchy(id uuid.UUID) (*DepartmentWithHierarchy, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid department id")
	}

	ctx := context.Background()
	cacheKey := s.cacheKeys.Build("hierarchy", id.String())

	// Try to get from cache first (Cache-Aside Pattern)
	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		var result DepartmentWithHierarchy
		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			s.logger.Debug("Department hierarchy retrieved from cache",
				logging.String("department_id", id.String()),
				logging.String("cache_key", cacheKey),
			)
			return &result, nil
		}
		// If unmarshal fails, continue to fetch from database
		s.logger.Warn("Failed to unmarshal cached department hierarchy",
			logging.String("department_id", id.String()),
			logging.Error(err),
		)
	}

	// Cache miss or error - fetch from database using CTE recursive query
	s.logger.Debug("Cache miss - fetching hierarchy from database using CTE",
		logging.String("department_id", id.String()),
	)

	// Use PostgreSQL CTE recursive to fetch entire hierarchy in a single query
	result, err := s.repo.FindHierarchyByID(id)
	if err != nil {
		s.logger.Error("Failed to fetch department hierarchy",
			logging.String("department_id", id.String()),
			logging.Error(err),
		)
		return nil, fmt.Errorf("department not found: %w", err)
	}

	// Store in cache for future requests
	if jsonData, err := json.Marshal(result); err == nil {
		if err := s.cache.Set(ctx, cacheKey, string(jsonData), s.cacheTTL); err != nil {
			s.logger.Warn("Failed to cache department hierarchy",
				logging.String("department_id", id.String()),
				logging.Error(err),
			)
		} else {
			s.logger.Debug("Department hierarchy cached successfully",
				logging.String("department_id", id.String()),
				logging.String("ttl", s.cacheTTL.String()),
			)
		}
	}

	return result, nil
}

// buildHierarchy is deprecated - replaced by PostgreSQL CTE recursive query in FindHierarchyByID
// Kept here for reference only. This method had O(N*M) complexity with N queries.
// The new CTE approach uses 1 single query for the entire hierarchy.
/*
func (s *Service) buildHierarchy(parentID uuid.UUID) ([]DepartmentWithHierarchy, error) {
	// Find all direct children
	children, err := s.repo.FindByParentID(parentID)
	if err != nil {
		return nil, err
	}

	result := make([]DepartmentWithHierarchy, 0, len(children))

	for _, child := range children {
		// Get manager name for each child
		manager, err := s.employeeRepo.FindByID(child.ManagerID)
		if err != nil {
			continue // Skip if manager not found
		}

		// Recursively build subdepartments
		subdepartments, err := s.buildHierarchy(child.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, DepartmentWithHierarchy{
			Department:     child,
			ManagerName:    manager.Name,
			Subdepartments: subdepartments,
		})
	}

	return result, nil
}
*/

func (s *Service) CreateDepartment(dept *Department) error {
	if err := s.validateDepartment(dept); err != nil {
		s.logger.Warn("Department validation failed",
			logging.String("name", dept.Name),
			logging.Error(err),
		)
		return err
	}

	// Validate no cycles in hierarchy (for create, we use a temporary UUID to simulate the new department)
	if dept.ParentDepartmentID != nil {
		// For new departments, just check if parent exists
		_, err := s.repo.FindByID(*dept.ParentDepartmentID)
		if err != nil {
			s.logger.Error("Parent department not found",
				logging.String("parent_id", dept.ParentDepartmentID.String()),
				logging.Error(err),
			)
			return errors.New("parent department not found")
		}
	}

	if err := s.repo.Create(dept); err != nil {
		s.logger.Error("Failed to create department in repository",
			logging.String("department_id", dept.ID.String()),
			logging.String("name", dept.Name),
			logging.Error(err),
		)
		return err
	}

	// Invalidate cache for parent hierarchy (if exists)
	if dept.ParentDepartmentID != nil {
		s.invalidateHierarchyCache(*dept.ParentDepartmentID)
	}

	s.logger.Info("Department created successfully",
		logging.String("department_id", dept.ID.String()),
		logging.String("name", dept.Name),
		logging.String("manager_id", dept.ManagerID.String()),
	)

	return nil
}

func (s *Service) UpdateDepartment(id uuid.UUID, dept *Department) error {
	if id == uuid.Nil {
		return errors.New("invalid department id")
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.Error("Department not found for update",
			logging.String("department_id", id.String()),
			logging.Error(err),
		)
		return fmt.Errorf("department not found: %w", err)
	}

	if err := s.validateDepartment(dept); err != nil {
		s.logger.Warn("Department update validation failed",
			logging.String("department_id", id.String()),
			logging.Error(err),
		)
		return err
	}

	// Validate no cycles in hierarchy
	if err := s.validateNoCycle(id, dept.ParentDepartmentID); err != nil {
		s.logger.Error("Cycle detected in department hierarchy",
			logging.String("department_id", id.String()),
			logging.Error(err),
		)
		return err
	}

	// Validate that manager belongs to the department
	if err := s.validateManagerBelongsToDepartment(dept.ManagerID, id); err != nil {
		s.logger.Error("Manager validation failed",
			logging.String("department_id", id.String()),
			logging.String("manager_id", dept.ManagerID.String()),
			logging.Error(err),
		)
		return err
	}

	dept.ID = existing.ID
	if err := s.repo.Update(dept); err != nil {
		s.logger.Error("Failed to update department in repository",
			logging.String("department_id", id.String()),
			logging.Error(err),
		)
		return err
	}

	// Invalidate cache for this department and parents
	s.invalidateHierarchyCache(id)
	if existing.ParentDepartmentID != nil {
		s.invalidateHierarchyCache(*existing.ParentDepartmentID)
	}
	if dept.ParentDepartmentID != nil && (existing.ParentDepartmentID == nil || *dept.ParentDepartmentID != *existing.ParentDepartmentID) {
		// If parent changed, invalidate new parent too
		s.invalidateHierarchyCache(*dept.ParentDepartmentID)
	}

	s.logger.Info("Department updated successfully",
		logging.String("department_id", id.String()),
		logging.String("name", dept.Name),
	)

	return nil
}

func (s *Service) DeleteDepartment(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid department id")
	}

	dept, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.Error("Department not found for deletion",
			logging.String("department_id", id.String()),
			logging.Error(err),
		)
		return fmt.Errorf("department not found: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("Failed to delete department in repository",
			logging.String("department_id", id.String()),
			logging.Error(err),
		)
		return err
	}

	// Invalidate cache for this department and parent
	s.invalidateHierarchyCache(id)
	if dept.ParentDepartmentID != nil {
		s.invalidateHierarchyCache(*dept.ParentDepartmentID)
	}

	s.logger.Info("Department deleted successfully",
		logging.String("department_id", id.String()),
		logging.String("name", dept.Name),
	)

	return nil
}

func (s *Service) validateDepartment(dept *Department) error {
	if dept.Name == "" {
		return errors.New("department name is required")
	}
	if dept.ManagerID == uuid.Nil {
		return errors.New("department manager is required")
	}
	return nil
}

func (s *Service) validateManagerBelongsToDepartment(managerID, departmentID uuid.UUID) error {
	// Check if manager exists
	manager, err := s.employeeRepo.FindByID(managerID)
	if err != nil {
		return errors.New("manager not found")
	}

	// Check if manager belongs to the department
	if manager.DepartmentID != departmentID {
		return errors.New("manager must be linked to the same department")
	}

	return nil
}

func (s *Service) validateNoCycle(departmentID uuid.UUID, parentDepartmentID *uuid.UUID) error {
	// If no parent department, no cycle possible
	if parentDepartmentID == nil || *parentDepartmentID == uuid.Nil {
		return nil
	}

	// Check if parent is the same as current department
	if *parentDepartmentID == departmentID {
		return errors.New("department cannot be its own parent")
	}

	// Traverse the hierarchy to detect cycles
	visited := make(map[uuid.UUID]bool)
	currentID := *parentDepartmentID

	for currentID != uuid.Nil {
		// If we've already visited this department, there's a cycle
		if visited[currentID] {
			return errors.New("cycle detected in department hierarchy")
		}

		// If we reached the original department, there's a cycle
		if currentID == departmentID {
			return errors.New("cycle detected in department hierarchy")
		}

		visited[currentID] = true

		// Get the parent of current department
		dept, err := s.repo.FindByID(currentID)
		if err != nil {
			return fmt.Errorf("department not found in hierarchy: %w", err)
		}

		// Move to the next parent
		if dept.ParentDepartmentID == nil {
			break
		}
		currentID = *dept.ParentDepartmentID
	}

	return nil
}

// invalidateHierarchyCache removes cached hierarchy for a department
func (s *Service) invalidateHierarchyCache(departmentID uuid.UUID) {
	ctx := context.Background()
	cacheKey := s.cacheKeys.Build("hierarchy", departmentID.String())

	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn("Failed to invalidate cache for department hierarchy",
			logging.String("department_id", departmentID.String()),
			logging.String("cache_key", cacheKey),
			logging.Error(err),
		)
	} else {
		s.logger.Debug("Cache invalidated for department hierarchy",
			logging.String("department_id", departmentID.String()),
			logging.String("cache_key", cacheKey),
		)
	}
}
