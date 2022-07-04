package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	"boe-backend/internal/types"
	jwtx "boe-backend/internal/util/jwt"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetShowListHandler(c *gin.Context) {
	var offset, _ = strconv.Atoi(c.Query("offset"))
	var count, _ = strconv.Atoi(c.Query("count"))
	var name = c.Query("name")

	var dbInstance = db.GetInstance()
	var shows []orm.Show
	dbInstance.Limit(count).Offset(offset).Where("name LIKE ?", "%"+name+"%").Find(&shows)

	var total int64
	dbInstance.Table("show").Where("name LIKE ?", "%"+name+"%").Where("deleted_at IS NULL").Count(&total)
	showDTOs := make([]types.ShowDTO, len(shows))
	for i := 0; i < len(shows); i++ {
		showDTOs[i].Show = shows[i]

		var m []string
		err := json.Unmarshal([]byte(shows[i].Images), &m)
		if err != nil {
			continue
		}
		if len(m) > 0 {
			showDTOs[i].Preview = m[0]
		}
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"total": total,
			"shows": showDTOs,
		},
	})
}

func AddShowHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)

	var show orm.Show
	err := c.ShouldBindJSON(&show)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  400,
			"error": "Bad request parameter",
		})
		return
	}

	var user orm.User
	db.GetInstance().First(&user, info.ID)
	show.Author = user.Username

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
