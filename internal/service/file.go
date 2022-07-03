package service

import (
	jwtx "boe-backend/internal/util/jwt"
	"boe-backend/internal/util/miniox"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
)

const (
	accessKey = "LQFcmFN3P_DwpBdzT5f_Php9MQ03qlFRF84zmHQW"
	secretKey = "2_qWBMZ38WWWsYGA-HPpzY8GfMtAY0e5kvKwKviN"
	bucket    = "yzlyzl123"
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

func GetFilelist(c *gin.Context) {
	prefix := c.Query("prefix") + "/"
	limit := 1000
	delimiter := ""
	//初始列举marker为空
	marker := ""
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	var files []string
	for {
		entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, prefix, delimiter, marker, limit)
		if err != nil {
			fmt.Println("list error,", err)
			break
		}
		//print entries
		for _, entry := range entries {
			files = append(files, entry.Key)
			//fmt.Println(entry.Key)
		}
		if hasNext {
			marker = nextMarker
		} else {
			//list end
			break
		}
	}

	c.JSON(200, gin.H{
		"code":    200,
		"data":    files,
		"message": "success",
	})
}
