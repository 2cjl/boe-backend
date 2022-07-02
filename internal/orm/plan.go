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
	UserID      int            `gorm:"column:user_id"`
	Author      User           `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
	PlayPeriods []PlayPeriod
	Devices     []Device `gorm:"many2many:plan_device;"`
}

func (t *Plan) TableName() string {
	return "plan"
}

type PlanDevice struct {
	ID       int `gorm:"column:id;primary_key"`
	PlanID   int `gorm:"column:plan_id"`
	DeviceID int `gorm:"column:device_id"`
}

func (t *PlanDevice) TableName() string {
	return "plan_device"
}
