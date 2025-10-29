package employee

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	CPF          string         `gorm:"type:varchar(11);uniqueIndex;not null" json:"cpf"`
	RG           *string        `gorm:"type:varchar(20);uniqueIndex" json:"rg,omitempty"`
	DepartmentID uuid.UUID      `gorm:"type:uuid;not null" json:"department_id"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Employee) TableName() string {
	return "employees"
}
