package service

import (
	jwtx "boe-backend/internal/util/jwt"
	"boe-backend/internal/util/miniox"
	"github.com/gin-gonic/gin"
	"log"
)

type PathInfo struct {
	Path string `json:"path" jsonschema:"required"`
}

func PreSignHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	var pathInfo PathInfo
	err := c.BindJSON(&pathInfo)
	if err != nil {
		return
	}
	url := miniox.PreSignObject(info.OrganizationID + "/" + pathInfo.Path)
	log.Println(url)
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
