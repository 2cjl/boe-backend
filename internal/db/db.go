package db

import (
	"boe-backend/internal/orm"
	"boe-backend/internal/util/config"
	"errors"
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

func Login(phone, passwd string) (*orm.User, error) {
	getInstance()
	var u orm.User
	db.Where("phone = ? AND passwd = ?", phone, passwd).First(&u)
	if u.ID == 0 {
		return nil, errors.New("does not exist")
	}
	return &u, nil
}
