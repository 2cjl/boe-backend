package service

import (
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"log"
)

func HomeAllHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	log.Println(info)
	c.JSON(200, gin.H{
		"error": "server internal error",
	})
}
