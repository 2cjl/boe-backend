package orm

import (
	"gorm.io/gorm"
	"time"
)

type PlayPeriod struct {
	ID        int `gorm:"column:id;primary_key"`
	PlanID    string
	StartTime string         `gorm:"column:start_time"`
	EndTime   string         `gorm:"column:end_time"`
	LoopMode  string         `gorm:"column:loop_mode"`
	Html      string         `gorm:"column:html"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	// 该计划对应的所有节目
	Shows []Show `gorm:"many2many:play_period_show;"`
}

func (t *PlayPeriod) TableName() string {
	return "play_period"
}
