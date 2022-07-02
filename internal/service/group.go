package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	"boe-backend/internal/types"
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func AddGroupHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)

	var group types.AddGroupReq
	err := c.ShouldBindJSON(&group)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  400,
			"error": "Bad request parameter",
		})
		return
	}

	var groupOrm orm.Group
	groupOrm.Name = group.Name
	oid, _ := strconv.Atoi(info.OrganizationID)
	groupOrm.OrganizationID = oid
	groupOrm.Describe = group.Describe
	res := db.GetInstance().Create(&groupOrm)
	if res.Error != nil {
		log.Println(res.Error)
		c.JSON(200, gin.H{
			"code":  500,
			"error": "add groupOrm fail",
		})
		return
	}

	var groupDevice []*orm.GroupDevice
	for _, id := range group.Devices {
		groupDevice = append(groupDevice, &orm.GroupDevice{
			DeviceID: id,
			GroupID:  groupOrm.ID,
		})
	}
	db.GetInstance().Create(&groupDevice)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

func GetGroupListHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	var offset, _ = strconv.Atoi(c.Query("offset"))
	var count, _ = strconv.Atoi(c.Query("count"))

	var groups []orm.Group
	db.GetInstance().Limit(count).Offset(offset).Find(&groups, "organization_id = ?", info.OrganizationID)
	gc := db.GetGroupDeviceCntByGroup(groups)

	m := make(map[int]int)
	for _, v := range gc {
		m[v.ID] = v.Cnt
	}

	groupDTOs := make([]types.GroupDTO, len(groups))

	for i := 0; i < len(groups); i++ {
		groupDTOs[i].Group = groups[i]
		groupDTOs[i].DeviceCnt = m[groups[i].ID]
	}

	var total int64
	db.GetInstance().Table("groups").Where("deleted_at IS NULL AND organization_id = ?", info.OrganizationID).Count(&total)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"groups": groupDTOs,
			"total":  total,
		},
	})
}

func GetGroupDevicesHandler(c *gin.Context) {
	var groupDevices []orm.GroupDevice
	gid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request parameter",
		})
		return
	}
	db.GetInstance().Find(&groupDevices, "group_id = ?", gid)
	devices := db.GetDevicesByGroupDevice(groupDevices)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    devices,
	})
}

func DeleteGroup(c *gin.Context) {
	gid := c.Param("id")

	var dbInstance = db.GetInstance()
	var group orm.Group
	dbInstance.Where("id = ?", gid).Find(&group).Delete(&group)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}
