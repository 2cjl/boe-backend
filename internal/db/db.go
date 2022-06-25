package db

import (
	"boe-backend/internal/orm"
	"boe-backend/internal/util/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func getInstance() {
	if db == nil {
		once.Do(func() {
			cfg := config.GetConfig().Mysql
			source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Passwd, cfg.Host, cfg.Port, cfg.Database)
			_db, err := gorm.Open(mysql.Open(source), &gorm.Config{})
			db = _db
			if err != nil {
				panic(err)
			}
		})
	}
}

func GetInstance() *gorm.DB {
	getInstance()
	return db
}

func Login(phone, passwd string) *orm.User {
	getInstance()
	var u orm.User
	db.Where("phone = ? AND passwd = ?", phone, passwd).First(&u)
	if u.ID == 0 {
		return nil
	}
	return &u
}

func GetDeviceByMac(mac string) *orm.Device {
	getInstance()
	var d orm.Device
	db.Where("mac = ?", mac).First(&d)
	if d.ID == 0 {
		return nil
	}
	return &d
}

func GetOrganizationById(id int) *orm.Organization {
	getInstance()
	var o orm.Organization
	db.First(&o, id)
	if o.ID == 0 {
		return nil
	}
	return &o
}

func GetOrganizationByUser(uid int) *orm.Organization {
	getInstance()
	var o orm.Organization
	db.First(&o, db.Table("users").Select("organization_id").Where("id = ?", uid))
	if o.ID == 0 {
		return nil
	}
	return &o
}

// GetAllEvents 获取所有的事件
func GetAllEvents(organizationId string) []orm.Event {
	getInstance()
	var users []orm.Event
	db.Find(&users, "organization_id = ?", organizationId)
	return users
}
