package orm

import (
	"gorm.io/gorm"
	"time"
)

type Plan struct {
	ID          int            `gorm:"column:id;primary_key"`
	Name        string         `gorm:"column:name"`
	State       string         `gorm:"column:state"`
	Mode        string         `gorm:"column:mode"`
	StartDate   string         `gorm:"column:start_date"` //2022-06-25
	EndDate     string         `gorm:"column:end_date"`   //2022-06-25
	Author      string         `gorm:"column:author"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
	PlayPeriods []PlayPeriod
}

func (t *Plan) TableName() string {
	return "plan"
}
