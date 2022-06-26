package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/devicemanager"
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
)

func GroupCountHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	gc := db.GetGroupDeviceCnt(info.OrganizationID)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    gc,
	})
}

func DevicesStateHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	var all, online, offline int64
	db.GetInstance().Table("devices").Where("organization_id = ?", info.OrganizationID).Count(&all)
	db.GetInstance().Table("devices").Where("organization_id = ? AND state = ?", info.OrganizationID, devicemanager.DeviceOnline).Count(&online)

	offline = all - online
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"All":     all,
			"Online":  online,
			"Offline": offline,
		},
	})
}

func CountHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	var deviceCnt, showCnt, planCnt int64
	db.GetInstance().Table("devices").Where("organization_id = ?", info.OrganizationID).Count(&deviceCnt)
	db.GetInstance().Table("show").Count(&showCnt)
	db.GetInstance().Table("plan").Count(&planCnt)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"DeviceCnt": deviceCnt,
			"ShowCnt":   showCnt,
			"PlanCnt":   planCnt,
		},
	})
}
