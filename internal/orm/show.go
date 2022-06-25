package orm

import "time"

type Show struct {
	ID        int       `gorm:"column:id;primary_key"`
	Name      string    `gorm:"column:name"`
	Duration  int       `gorm:"column:duration"`
	Author    string    `gorm:"column:author"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}

func (t *Show) TableName() string {
	return "show"
}

type PlayPeriodAndShow struct {
	ID           int       `gorm:"column:id;primary_key"`
	PlayPeriodID int       `gorm:"column:play_period_id"`
	ShowID       int       `gorm:"column:show_id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	DeletedAt    time.Time `gorm:"column:deleted_at"`
}

func (t *PlayPeriodAndShow) TableName() string {
	return "play_period_show"
}
