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
	if err != nil || strconv.Itoa(group.OrganizationID) != info.OrganizationID {
		c.JSON(200, gin.H{
			"code":  400,
			"error": "Bad request parameter",
		})
		return
	}

	var groupOrm orm.Group
	groupOrm.Name = group.Name
	groupOrm.OrganizationID = group.OrganizationID
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
	list := db.GetAllGroups(info.OrganizationID)
	gc := db.GetGroupDeviceCnt(info.OrganizationID)

	m := make(map[int]int)
	for _, v := range gc {
		m[v.ID] = v.Cnt
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"groups":    list,
			"deviceCnt": m,
		},
	})
}

func GetGroupInfoHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)
	log.Println(info)
	c.JSON(200, gin.H{})
}
