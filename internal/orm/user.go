package orm

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             int            `gorm:"column:id;primary_key"`
	OrganizationID int            `gorm:"column:organization_id"`
	Phone          string         `gorm:"column:phone"`
	RealName       string         `gorm:"column:real_name"`
	Username       string         `gorm:"column:username"`
	Passwd         string         `gorm:"column:passwd"`
	Email          string         `gorm:"column:email"`
	Status         string         `gorm:"column:status"`
	Avatar         string         `gorm:"column:avatar"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
	Organization   Organization   `gorm:"foreignKey:OrganizationID;references:ID"`
}
