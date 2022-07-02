package orm

import (
	"gorm.io/gorm"
	"time"
)

type Device struct {
	ID             int            `gorm:"column:id;primary_key"`
	Name           string         `gorm:"column:name"`
	OrganizationID int            `gorm:"column:organization_id"`
	Mac            string         `gorm:"column:mac"`
	PlanID         int            `gorm:"column:plan_id"`
	State          string         `gorm:"column:state"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
}

type DeviceInfo struct {
	ID            int            `gorm:"column:id;primary_key"`
	Model         string         `gorm:"column:model"`
	IP            string         `gorm:"column:ip"`
	HardwareID    string         `gorm:"column:hardware_id"`
	Latitude      float64        `gorm:"column:latitude"`
	Longitude     float64        `gorm:"column:longitude"`
	LastHeartbeat time.Time      `gorm:"column:last_heartbeat"`
	RunningTime   uint64         `gorm:"column:running_time"`
	Resolution    string         `gorm:"column:resolution"`
	AppVersion    string         `gorm:"column:app_version"`
	Memory        string         `gorm:"column:memory"`
	Storage       string         `gorm:"column:storage"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (i *DeviceInfo) TableName() string {
	return "device_info"
}
