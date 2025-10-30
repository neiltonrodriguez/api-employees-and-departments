package ginapi

import (
	"net/http"

	"api-employees-and-departments/internal/domain/department"
	"api-employees-and-departments/internal/domain/employee"
	"api-employees-and-departments/internal/interfaces/api/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ManagerHandler struct {
	departmentService *department.Service
	employeeService   *employee.Service
}

func NewManagerHandler(deptService *department.Service, empService *employee.Service) *ManagerHandler {
	return &ManagerHandler{
		departmentService: deptService,
		employeeService:   empService,
	}
}

// GetSubordinateEmployees godoc
// @Summary Get all employees subordinate to a manager (recursive)
// @Tags managers
// @Accept json
// @Produce json
// @Param id path string true "Manager ID (Employee ID)"
// @Success 200 {array} dto.EmployeeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /managers/{id}/employees [get]
func (h *ManagerHandler) GetSubordinateEmployees(c *gin.Context) {
	managerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid manager ID format",
		})
		return
	}

	// Verify that the manager exists
	_, err = h.employeeService.GetEmployeeByID(managerID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Manager not found",
		})
		return
	}

	// Get all department IDs managed by this manager (recursively)
	departmentIDs, err := h.getAllManagedDepartmentIDs(managerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	// Get all employees from these departments
	employees, err := h.employeeService.GetEmployeesByDepartmentIDs(departmentIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToEmployeeResponseList(employees))
}

func (h *ManagerHandler) getAllManagedDepartmentIDs(managerID uuid.UUID) ([]uuid.UUID, error) {
	// Find all departments where this employee is the manager
	departments, err := h.departmentService.GetDepartmentsByManagerID(managerID)
	if err != nil {
		return nil, err
	}

	var result []uuid.UUID
	for _, dept := range departments {
		result = append(result, dept.ID)

		// Recursively get subdepartments
		subdeptIDs, err := h.getAllSubdepartmentIDs(dept.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, subdeptIDs...)
	}

	return result, nil
}

func (h *ManagerHandler) getAllSubdepartmentIDs(parentID uuid.UUID) ([]uuid.UUID, error) {
	children, err := h.departmentService.GetDepartmentsByParentID(parentID)
	if err != nil {
		return nil, err
	}

	var result []uuid.UUID
	for _, child := range children {
		result = append(result, child.ID)

		// Recursive call
		subdeptIDs, err := h.getAllSubdepartmentIDs(child.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, subdeptIDs...)
	}

	return result, nil
}
