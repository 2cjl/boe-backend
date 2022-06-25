package service

import (
	"boe-backend/internal/orm"
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"log"
)

func AddGroupHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)

	var group orm.Group
	err := c.ShouldBindJSON(&group)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  400,
			"error": "Bad request parameter",
		})
		return
	}

	log.Println(info)
	c.JSON(200, gin.H{})
}

func GetGroupListHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	log.Println(info)
	c.JSON(200, gin.H{})
}

func GetGroupInfoHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	log.Println(info)
	c.JSON(200, gin.H{})
}
