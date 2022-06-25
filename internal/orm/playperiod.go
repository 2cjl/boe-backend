package orm

import "time"

type PlayPeriod struct {
	ID        int       `gorm:"column:id;primary_key"`
	StartTime string    `gorm:"column:start_time"`
	EndTime   string    `gorm:"column:end_time"`
	LoopMode  string    `gorm:"column:loop_mode"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}

func (t *PlayPeriod) TableName() string {
	return "play_period"
}

type PlayPeriodAndPlan struct {
	ID           int       `gorm:"column:id;primary_key"`
	PlanID       int       `gorm:"column:plan_id"`
	PlayPeriodID int       `gorm:"column:play_period_id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	DeletedAt    time.Time `gorm:"column:deleted_at"`
}

func (t *PlayPeriodAndPlan) TableName() string {
	return "play_period_plan"
}
