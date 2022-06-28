package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/devicemanager"
	"boe-backend/internal/orm"
	"boe-backend/internal/types"
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func AddDeviceHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)

	var req types.AddDeviceReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{
			"code":  400,
			"error": "Bad request parameter",
		})
		return
	}

	var device orm.Device
	oid, _ := strconv.Atoi(info.OrganizationID)
	device.OrganizationID = oid
	device.Name = req.Name
	device.Mac = req.Mac
	device.State = devicemanager.DeviceOffline
	db.GetInstance().Create(&device)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    device,
	})
}

func UpdateDevice(c *gin.Context) {
	var req types.AddDeviceReq
	err := c.ShouldBindJSON(&req)
	if err != nil || req.ID == 0 {
		log.Println(err)
		c.JSON(200, gin.H{
			"code":  400,
			"error": "Bad request parameter",
		})
		return
	}
	var device orm.Device
	device.ID = req.ID
	device.Name = req.Name
	device.Mac = req.Mac
	db.GetInstance().Updates(&device)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    device,
	})
}

func GetDeviceListHandler(c *gin.Context) {
	var offset, _ = strconv.Atoi(c.Query("offset"))
	var count, _ = strconv.Atoi(c.Query("count"))
	var dbInstance = db.GetInstance()
	var devices []orm.Device
	dbInstance.Limit(count).Offset(offset).Find(&devices)

	var total int64
	dbInstance.Table("devices").Count(&total)

	names := make(map[int]string)
	var ids []int
	for _, v := range devices {
		ids = append(ids, v.PlanID)
	}
	plans := db.GetPlanByIds(ids)
	for _, plan := range plans {
		names[plan.ID] = plan.Name
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"total":     total,
			"devices":   devices,
			"plansName": names,
		},
	})
}

func GetDeviceInfoHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var deviceInfo orm.DeviceInfo
	var device orm.Device
	deviceInfo.ID = id
	db.GetInstance().Find(&deviceInfo)
	device.ID = deviceInfo.ID
	db.GetInstance().Find(&device)
	if d := devicemanager.GetDeviceByMac(device.Mac); d != nil {
		deviceInfo.LastHeartbeat = d.LastHeartbeat
		deviceInfo.RunningTime = d.RunningTime
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    deviceInfo,
	})
}

func DeleteDevice(c *gin.Context) {
	id := c.Param("id")

	var device orm.Device
	db.GetInstance().Where("id = ?", id).Find(&device).Delete(&device)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}
