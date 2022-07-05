package orm

type Notice struct {
	ID      int    `gorm:"column:id;primary_key"`
	Name    string `gorm:"column:name"`
	Content string `gorm:"column:content"`
}

func (t *Notice) TableName() string {
	return "notice"
}
