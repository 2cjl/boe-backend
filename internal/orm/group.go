package orm

import (
	"gorm.io/gorm"
	"time"
)

type Group struct {
	ID             int            `gorm:"column:id;primary_key"`
	Name           string         `gorm:"column:name"`
	Describe       string         `gorm:"column:describe"`
	OrganizationID int            `gorm:"column:organization_id"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (t *Group) TableName() string {
	return "group"
}

type GroupDevice struct {
	ID        int            `gorm:"column:id;primary_key"`
	DeviceID  int            `gorm:"column:device_id"`
	GroupID   int            `gorm:"column:group_id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (t *GroupDevice) TableName() string {
	return "group_device"
}
