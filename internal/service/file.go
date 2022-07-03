package service

import (
	jwtx "boe-backend/internal/util/jwt"
	"boe-backend/internal/util/miniox"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
)

const (
	accessKey = "LQFcmFN3P_DwpBdzT5f_Php9MQ03qlFRF84zmHQW"
	secretKey = "2_qWBMZ38WWWsYGA-HPpzY8GfMtAY0e5kvKwKviN"
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
	url := miniox.PreSignObject(info.ID + "/" + pathInfo.Path)
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

func GetUploadToken(c *gin.Context) {
	bucket := "yzlyzl123"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	c.JSON(200, gin.H{
		"code":    200,
		"data":    upToken,
		"message": "success",
	})
}
