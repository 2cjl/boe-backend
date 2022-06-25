package service

import (
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"log"
)

func AddDeviceHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	log.Println(info)
	c.JSON(200, gin.H{})
}

func GetDeviceListHandler(c *gin.Context) {

}
