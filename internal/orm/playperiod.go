package orm

import (
	"gorm.io/gorm"
	"time"
)

type PlayPeriod struct {
	ID        int            `gorm:"column:id;primary_key"`
	StartTime string         `gorm:"column:start_time"`
	EndTime   string         `gorm:"column:end_time"`
	LoopMode  string         `gorm:"column:loop_mode"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	PlanID    string
	Plan      Plan
}

func (t *PlayPeriod) TableName() string {
	return "play_period"
}
