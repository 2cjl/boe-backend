package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetShowListHandler(c *gin.Context) {
	var offset, _ = strconv.Atoi(c.Query("offset"))
	var count, _ = strconv.Atoi(c.Query("count"))
	var dbInstance = db.GetInstance()
	var shows []orm.Show
	dbInstance.Limit(count).Offset(offset).Find(&shows)

	var total int64
	dbInstance.Table("show").Where("deleted_at IS NULL").Count(&total)
	previews := make(map[int]string)
	for _, show := range shows {
		var m []string
		err := json.Unmarshal([]byte(show.Images), &m)
		if err != nil {
			continue
		}
		if len(m) > 0 {
			previews[show.ID] = m[0]
		}
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"total":    total,
			"shows":    shows,
			"previews": previews,
		},
	})
}

func AddShowHandler(c *gin.Context) {
	var show orm.Show
	err := c.ShouldBindJSON(&show)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  400,
			"error": "Bad request parameter",
		})
		return
	}
	show.ID = 0
	db.GetInstance().Create(&show)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    show,
	})
}

func DeleteShow(c *gin.Context) {
	id := c.Param("id")

	var dbInstance = db.GetInstance()
	var show orm.Show
	dbInstance.Where("id = ?", id).Find(&show).Delete(&show)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

func UpdateShow(c *gin.Context) {
	var show orm.Show
	err := c.ShouldBindJSON(&show)
	if err != nil || show.ID == 0 {
		c.JSON(200, gin.H{
			"code":  400,
			"error": "Bad request parameter",
		})
		return
	}
	db.GetInstance().Updates(&show)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    show,
	})
}
