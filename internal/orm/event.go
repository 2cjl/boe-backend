package orm

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	ID         int            `gorm:"column:id;primary_key"`
	Time       time.Time      `gorm:"column:time"`
	ObjectType string         `gorm:"column:object_type"`
	Content    string         `gorm:"column:content"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"`
}
