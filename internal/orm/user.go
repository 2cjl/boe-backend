package orm

import "time"

type User struct {
	ID             int       `gorm:"column:id;primary_key"`
	OrganizationID int       `gorm:"column:organization_id"`
	Phone          string    `gorm:"column:phone"`
	Username       string    `gorm:"column:username"`
	Passwd         string    `gorm:"column:passwd"`
	Avatar         string    `gorm:"column:avatar"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	DeletedAt      time.Time `gorm:"column:deleted_at"`
}
