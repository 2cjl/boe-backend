package service

import (
	jwtx "boe-backend/internal/util/jwt"
	"boe-backend/internal/util/miniox"
	"github.com/gin-gonic/gin"
	"log"
)

func PreSignHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	log.Println(info)
	path := c.Param("path")
	url := miniox.PreSignObject(path)
	if url == nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "presign fail",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"data":    url.String(),
		"message": "success",
	})
}
