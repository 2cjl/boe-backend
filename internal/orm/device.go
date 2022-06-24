package orm

import "time"

type Device struct {
	ID             int       `gorm:"column:id;primary_key"`
	Name           string    `gorm:"column:name"`
	OrganizationID int       `gorm:"column:organization_id"`
	Mac            string    `gorm:"column:mac"`
	Sn             string    `gorm:"column:sn"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	DeletedAt      time.Time `gorm:"column:deleted_at"`
}
