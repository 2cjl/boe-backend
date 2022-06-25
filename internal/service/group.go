package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func AddGroupHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)

	var group orm.Group
	err := c.ShouldBindJSON(&group)
	if err != nil || strconv.Itoa(group.OrganizationID) != info.OrganizationID {
		c.JSON(200, gin.H{
			"code":  400,
			"error": "Bad request parameter",
		})
		return
	}

	res := db.GetInstance().Create(&group)
	if res.Error != nil {
		log.Println(res.Error)
		c.JSON(200, gin.H{
			"code":  500,
			"error": "add group fail",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

func GetGroupListHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	list := db.GetAllGroups(info.OrganizationID)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    list,
	})
}

func GetGroupInfoHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	log.Println(info)
	c.JSON(200, gin.H{})
}
