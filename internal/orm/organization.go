package orm

import (
	"gorm.io/gorm"
	"time"
)

type Organization struct {
	ID        int            `gorm:"column:id;primary_key"`
	Name      string         `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}
