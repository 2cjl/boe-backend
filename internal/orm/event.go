package orm

import "time"

type Event struct {
	ID         int       `gorm:"column:id;primary_key"`
	Time       time.Time `gorm:"column:time"`
	ObjectType string    `gorm:"column:object_type"`
	Content    string    `gorm:"column:content"`
}
