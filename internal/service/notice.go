package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	"boe-backend/internal/types"
	"github.com/gin-gonic/gin"
)

func GetNotice(c *gin.Context) {
	var notice orm.Notice
	// 简化，全局一篇公告
	db.GetInstance().Where("id = ?", 1).First(&notice)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    notice,
	})
}

func UpdateNotice(c *gin.Context) {
	var req types.UpdateNoticeRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "req error",
		})
	}

	var notice orm.Notice
	notice.ID = 1
	notice.Name = req.Name
	notice.Content = req.Content
	// 简化，全局一篇公告
	db.GetInstance().Updates(&notice)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    notice,
	})
}
