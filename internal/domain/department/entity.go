package department

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Department struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name               string         `gorm:"type:varchar(255);not null" json:"name"`
	ManagerID          uuid.UUID      `gorm:"type:uuid;not null" json:"manager_id"`
	ParentDepartmentID *uuid.UUID     `gorm:"type:uuid" json:"parent_department_id,omitempty"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Department) TableName() string {
	return "departments"
}
