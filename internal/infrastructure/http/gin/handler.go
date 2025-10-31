package ginapi

import (
	"net/http"

	"api-employees-and-departments/internal/domain/employee"
	"api-employees-and-departments/internal/infrastructure/logging"
	"api-employees-and-departments/internal/interfaces/api/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EmployeeHandler struct {
	service *employee.Service
}

func NewEmployeeHandler(s *employee.Service) *EmployeeHandler {
	return &EmployeeHandler{service: s}
}

// GetAll godoc
// @Summary List all employees
// @Tags employees
// @Accept json
// @Produce json
// @Success 200 {array} dto.EmployeeResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /employees [get]
func (h *EmployeeHandler) GetAll(c *gin.Context) {
	employees, err := h.service.GetAllEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_server_error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.ToEmployeeResponseList(employees))
}

// GetByID godoc
// @Summary Get employee by ID with manager name
// @Tags employees
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} dto.EmployeeWithManagerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /employees/{id} [get]
func (h *EmployeeHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid employee ID format",
		})
		return
	}

	empWithManager, err := h.service.GetEmployeeWithManager(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Employee not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.EmployeeWithManagerResponse{
		ID:           empWithManager.Employee.ID,
		Name:         empWithManager.Employee.Name,
		CPF:          empWithManager.Employee.CPF,
		RG:           empWithManager.Employee.RG,
		DepartmentID: empWithManager.Employee.DepartmentID,
		ManagerName:  empWithManager.ManagerName,
		CreatedAt:    empWithManager.Employee.CreatedAt,
		UpdatedAt:    empWithManager.Employee.UpdatedAt,
	})
}

// Create godoc
// @Summary Create a new employee
// @Tags employees
// @Accept json
// @Produce json
// @Param employee body dto.CreateEmployeeRequest true "Employee data"
// @Success 201 {object} dto.EmployeeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /employees [post]
func (h *EmployeeHandler) Create(c *gin.Context) {
	var req dto.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Warn("Employee creation validation failed",
			zap.Error(err),
			zap.String("request_id", getRequestID(c)),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	emp := dto.ToEmployeeEntity(&req)
	if err := h.service.CreateEmployee(emp); err != nil {
		logging.Error("Failed to create employee",
			zap.Error(err),
			zap.String("name", req.Name),
			zap.String("cpf", req.CPF),
			zap.String("request_id", getRequestID(c)),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "creation_failed",
			Message: err.Error(),
		})
		return
	}

	logging.Info("Employee created successfully",
		zap.String("employee_id", emp.ID.String()),
		zap.String("name", emp.Name),
		zap.String("request_id", getRequestID(c)),
	)

	c.JSON(http.StatusCreated, dto.ToEmployeeResponse(emp))
}

// Update godoc
// @Summary Update an employee
// @Tags employees
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param employee body dto.UpdateEmployeeRequest true "Employee data"
// @Success 200 {object} dto.EmployeeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /employees/{id} [put]
func (h *EmployeeHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid employee ID format",
		})
		return
	}

	var req dto.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	emp := dto.ToEmployeeEntityFromUpdate(&req)
	if err := h.service.UpdateEmployee(id, emp); err != nil {
		logging.Error("Failed to update employee",
			zap.Error(err),
			zap.String("employee_id", id.String()),
			zap.String("request_id", getRequestID(c)),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	logging.Info("Employee updated successfully",
		zap.String("employee_id", id.String()),
		zap.String("request_id", getRequestID(c)),
	)

	c.JSON(http.StatusOK, dto.ToEmployeeResponse(emp))
}

// Delete godoc
// @Summary Delete an employee
// @Tags employees
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /employees/{id} [delete]
func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid employee ID format",
		})
		return
	}

	if err := h.service.DeleteEmployee(id); err != nil {
		logging.Error("Failed to delete employee",
			zap.Error(err),
			zap.String("employee_id", id.String()),
			zap.String("request_id", getRequestID(c)),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "deletion_failed",
			Message: err.Error(),
		})
		return
	}

	logging.Info("Employee deleted successfully",
		zap.String("employee_id", id.String()),
		zap.String("request_id", getRequestID(c)),
	)

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary List employees with filters and pagination
// @Tags employees
// @Accept json
// @Produce json
// @Param filters body dto.ListEmployeesRequest true "Filter and pagination params"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /employees/list [post]
func (h *EmployeeHandler) List(c *gin.Context) {
	var req dto.ListEmployeesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Build filters
	filters := employee.ListFilters{
		Name:     req.Name,
		CPF:      req.CPF,
		RG:       req.RG,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// Parse department ID if provided
	if req.DepartmentID != nil && *req.DepartmentID != "" {
		deptID, err := uuid.Parse(*req.DepartmentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_filter",
				Message: "Invalid department ID format",
			})
			return
		}
		filters.DepartmentID = &deptID
	}

	employees, total, err := h.service.ListEmployees(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "list_failed",
			Message: err.Error(),
		})
		return
	}

	// Calculate total pages
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data:       dto.ToEmployeeResponseList(employees),
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// getRequestID retrieves the request ID from Gin context
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("RequestID"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
