package ginapi

import (
	"net/http"

	"api-employees-and-departments/internal/domain/department"
	"api-employees-and-departments/internal/infrastructure/logging"
	"api-employees-and-departments/internal/interfaces/api/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DepartmentHandler struct {
	service *department.Service
}

func NewDepartmentHandler(s *department.Service) *DepartmentHandler {
	return &DepartmentHandler{service: s}
}

// GetAll godoc
// @Summary List all departments
// @Tags departments
// @Accept json
// @Produce json
// @Success 200 {array} dto.DepartmentResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /departments [get]
func (h *DepartmentHandler) GetAll(c *gin.Context) {
	departments, err := h.service.GetAllDepartments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_server_error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.ToDepartmentResponseList(departments))
}

// GetByID godoc
// @Summary Get department by ID with full hierarchy and manager name
// @Tags departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} dto.DepartmentWithHierarchyResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /departments/{id} [get]
func (h *DepartmentHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid department ID format",
		})
		return
	}

	deptWithHierarchy, err := h.service.GetDepartmentWithHierarchy(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Department not found",
		})
		return
	}

	c.JSON(http.StatusOK, h.toHierarchyResponse(deptWithHierarchy))
}

func (h *DepartmentHandler) toHierarchyResponse(dept *department.DepartmentWithHierarchy) dto.DepartmentWithHierarchyResponse {
	subdepartments := make([]dto.DepartmentWithHierarchyResponse, 0, len(dept.Subdepartments))
	for _, sub := range dept.Subdepartments {
		subdepartments = append(subdepartments, h.toHierarchyResponse(&sub))
	}

	return dto.DepartmentWithHierarchyResponse{
		ID:                 dept.Department.ID,
		Name:               dept.Department.Name,
		ManagerID:          dept.Department.ManagerID,
		ManagerName:        dept.ManagerName,
		ParentDepartmentID: dept.Department.ParentDepartmentID,
		Subdepartments:     subdepartments,
		CreatedAt:          dept.Department.CreatedAt,
		UpdatedAt:          dept.Department.UpdatedAt,
	}
}

// Create godoc
// @Summary Create a new department
// @Description Create a new department. Use null (without quotes) for parent_department_id if creating a root department.
// @Tags departments
// @Accept json
// @Produce json
// @Param department body dto.CreateDepartmentRequest true "Department data"
// @Success 201 {object} dto.DepartmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /departments [post]
func (h *DepartmentHandler) Create(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	dept := dto.ToDepartmentEntity(&req)
	if err := h.service.CreateDepartment(dept); err != nil {
		logging.Error("Failed to create department",
			zap.Error(err),
			zap.String("name", req.Name),
			zap.String("request_id", getRequestID(c)),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "creation_failed",
			Message: err.Error(),
		})
		return
	}

	logging.Info("Department created successfully",
		zap.String("department_id", dept.ID.String()),
		zap.String("name", dept.Name),
		zap.String("request_id", getRequestID(c)),
	)

	c.JSON(http.StatusCreated, dto.ToDepartmentResponse(dept))
}

// Update godoc
// @Summary Update a department
// @Description Update a department. Use null (without quotes) for parent_department_id to remove parent reference.
// @Tags departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Param department body dto.UpdateDepartmentRequest true "Department data"
// @Success 200 {object} dto.DepartmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /departments/{id} [put]
func (h *DepartmentHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid department ID format",
		})
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	dept := dto.ToDepartmentEntityFromUpdate(&req)
	if err := h.service.UpdateDepartment(id, dept); err != nil {
		logging.Error("Failed to update department",
			zap.Error(err),
			zap.String("department_id", id.String()),
			zap.String("request_id", getRequestID(c)),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	logging.Info("Department updated successfully",
		zap.String("department_id", id.String()),
		zap.String("request_id", getRequestID(c)),
	)

	c.JSON(http.StatusOK, dto.ToDepartmentResponse(dept))
}

// Delete godoc
// @Summary Delete a department
// @Tags departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /departments/{id} [delete]
func (h *DepartmentHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid department ID format",
		})
		return
	}

	if err := h.service.DeleteDepartment(id); err != nil {
		logging.Error("Failed to delete department",
			zap.Error(err),
			zap.String("department_id", id.String()),
			zap.String("request_id", getRequestID(c)),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "deletion_failed",
			Message: err.Error(),
		})
		return
	}

	logging.Info("Department deleted successfully",
		zap.String("department_id", id.String()),
		zap.String("request_id", getRequestID(c)),
	)

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary List departments with filters and pagination
// @Tags departments
// @Accept json
// @Produce json
// @Param filters body dto.ListDepartmentsRequest true "Filter and pagination params"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /departments/list [post]
func (h *DepartmentHandler) List(c *gin.Context) {
	var req dto.ListDepartmentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Build filters
	filters := department.ListFilters{
		Name:        req.Name,
		ManagerName: req.ManagerName,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}

	// Parse parent department ID if provided
	if req.ParentDepartmentID != nil && *req.ParentDepartmentID != "" {
		parentID, err := uuid.Parse(*req.ParentDepartmentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_filter",
				Message: "Invalid parent department ID format",
			})
			return
		}
		filters.ParentDepartmentID = &parentID
	}

	departments, total, err := h.service.ListDepartments(filters)
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
		Data:       dto.ToDepartmentResponseList(departments),
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}
