package orm

import (
	"gorm.io/gorm"
	"time"
)

type Show struct {
	ID        int            `gorm:"column:id;primary_key"`
	Name      string         `gorm:"column:name"`
	Duration  int            `gorm:"column:duration"`
	Author    string         `gorm:"column:author"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (t *Show) TableName() string {
	return "show"
}
