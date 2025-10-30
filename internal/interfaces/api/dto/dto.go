package dto

import (
	"api-employees-and-departments/internal/domain/department"
	"api-employees-and-departments/internal/domain/employee"
	"time"

	"github.com/google/uuid"
)

// Employee DTOs
type CreateEmployeeRequest struct {
	Name         string     `json:"name" binding:"required"`
	CPF          string     `json:"cpf" binding:"required,len=11"`
	RG           *string    `json:"rg,omitempty"`
	DepartmentID uuid.UUID  `json:"department_id" binding:"required"`
}

type UpdateEmployeeRequest struct {
	Name         string     `json:"name" binding:"required"`
	CPF          string     `json:"cpf" binding:"required,len=11"`
	RG           *string    `json:"rg,omitempty"`
	DepartmentID uuid.UUID  `json:"department_id" binding:"required"`
}

type EmployeeResponse struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	CPF          string     `json:"cpf"`
	RG           *string    `json:"rg,omitempty"`
	DepartmentID uuid.UUID  `json:"department_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type EmployeeWithManagerResponse struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	CPF          string     `json:"cpf"`
	RG           *string    `json:"rg,omitempty"`
	DepartmentID uuid.UUID  `json:"department_id"`
	ManagerName  string     `json:"manager_name"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Department DTOs
type CreateDepartmentRequest struct {
	Name               string     `json:"name" binding:"required"`
	ManagerID          uuid.UUID  `json:"manager_id" binding:"required"`
	ParentDepartmentID *uuid.UUID `json:"parent_department_id,omitempty"`
}

type UpdateDepartmentRequest struct {
	Name               string     `json:"name" binding:"required"`
	ManagerID          uuid.UUID  `json:"manager_id" binding:"required"`
	ParentDepartmentID *uuid.UUID `json:"parent_department_id,omitempty"`
}

type DepartmentResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Name               string     `json:"name"`
	ManagerID          uuid.UUID  `json:"manager_id"`
	ParentDepartmentID *uuid.UUID `json:"parent_department_id,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type DepartmentWithHierarchyResponse struct {
	ID                 uuid.UUID                          `json:"id"`
	Name               string                             `json:"name"`
	ManagerID          uuid.UUID                          `json:"manager_id"`
	ManagerName        string                             `json:"manager_name"`
	ParentDepartmentID *uuid.UUID                         `json:"parent_department_id,omitempty"`
	Subdepartments     []DepartmentWithHierarchyResponse  `json:"subdepartments"`
	CreatedAt          time.Time                          `json:"created_at"`
	UpdatedAt          time.Time                          `json:"updated_at"`
}

// Error Response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Pagination Request
type PaginationRequest struct {
	Page     int `json:"page" binding:"min=1"`
	PageSize int `json:"page_size" binding:"min=1,max=100"`
}

// Employee List Request with filters
type ListEmployeesRequest struct {
	Name         *string    `json:"name,omitempty"`
	CPF          *string    `json:"cpf,omitempty"`
	RG           *string    `json:"rg,omitempty"`
	DepartmentID *string    `json:"department_id,omitempty"`
	Page         int        `json:"page" binding:"required,min=1"`
	PageSize     int        `json:"page_size" binding:"required,min=1,max=100"`
}

// Department List Request with filters
type ListDepartmentsRequest struct {
	Name                 *string `json:"name,omitempty"`
	ManagerName          *string `json:"manager_name,omitempty"`
	ParentDepartmentID   *string `json:"parent_department_id,omitempty"`
	Page                 int     `json:"page" binding:"required,min=1"`
	PageSize             int     `json:"page_size" binding:"required,min=1,max=100"`
}

// Paginated Response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// Converters - Employee
func ToEmployeeEntity(req *CreateEmployeeRequest) *employee.Employee {
	return &employee.Employee{
		Name:         req.Name,
		CPF:          req.CPF,
		RG:           req.RG,
		DepartmentID: req.DepartmentID,
	}
}

func ToEmployeeEntityFromUpdate(req *UpdateEmployeeRequest) *employee.Employee {
	return &employee.Employee{
		Name:         req.Name,
		CPF:          req.CPF,
		RG:           req.RG,
		DepartmentID: req.DepartmentID,
	}
}

func ToEmployeeResponse(emp *employee.Employee) *EmployeeResponse {
	return &EmployeeResponse{
		ID:           emp.ID,
		Name:         emp.Name,
		CPF:          emp.CPF,
		RG:           emp.RG,
		DepartmentID: emp.DepartmentID,
		CreatedAt:    emp.CreatedAt,
		UpdatedAt:    emp.UpdatedAt,
	}
}

func ToEmployeeResponseList(employees []employee.Employee) []EmployeeResponse {
	responses := make([]EmployeeResponse, len(employees))
	for i, emp := range employees {
		responses[i] = *ToEmployeeResponse(&emp)
	}
	return responses
}

// Converters - Department
func ToDepartmentEntity(req *CreateDepartmentRequest) *department.Department {
	return &department.Department{
		Name:               req.Name,
		ManagerID:          req.ManagerID,
		ParentDepartmentID: req.ParentDepartmentID,
	}
}

func ToDepartmentEntityFromUpdate(req *UpdateDepartmentRequest) *department.Department {
	return &department.Department{
		Name:               req.Name,
		ManagerID:          req.ManagerID,
		ParentDepartmentID: req.ParentDepartmentID,
	}
}

func ToDepartmentResponse(dept *department.Department) *DepartmentResponse {
	return &DepartmentResponse{
		ID:                 dept.ID,
		Name:               dept.Name,
		ManagerID:          dept.ManagerID,
		ParentDepartmentID: dept.ParentDepartmentID,
		CreatedAt:          dept.CreatedAt,
		UpdatedAt:          dept.UpdatedAt,
	}
}

func ToDepartmentResponseList(departments []department.Department) []DepartmentResponse {
	responses := make([]DepartmentResponse, len(departments))
	for i, dept := range departments {
		responses[i] = *ToDepartmentResponse(&dept)
	}
	return responses
}
