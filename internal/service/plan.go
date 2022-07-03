package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/devicemanager"
	"boe-backend/internal/orm"
	"boe-backend/internal/types"
	jwtx "boe-backend/internal/util/jwt"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CreatePlan(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	var user = t.(*jwtx.TokenUserInfo)

	var req types.CreatePlanRequest
	var ins = db.GetInstance()

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "request param error",
		})
		return
	}

	var plan orm.Plan
	plan.UserID, _ = strconv.Atoi(user.ID)
	plan.EndDate = req.EndDate
	plan.StartDate = req.StartDate
	plan.Name = req.Name
	plan.Mode = req.Mode
	// 初始设置为未发布状态
	plan.State = "未发布"

	// 开始初始化各个时间段，并复制
	var playPeriods []orm.PlayPeriod

	// 在内存中先初始化实体
	for _, period := range req.PlayPeriods {
		var p orm.PlayPeriod
		p.StartTime = period.StartTime
		p.EndTime = period.EndTime
		p.LoopMode = period.LoopMode
		p.Html = period.Html

		var shows []orm.Show
		// PREF: 在 for 循环之外合并所有传入的 showID 然后通过一次查询的查到所有用的 show
		ins.Find(&shows, period.ShowIds)
		p.Shows = shows

		playPeriods = append(playPeriods, p)
	}

	plan.PlayPeriods = playPeriods
	if req.DeviceIds != nil {
		ins.Find(&plan.Devices, req.DeviceIds)
	}
	var dbInstance = db.GetInstance()
	// 保存计划实体
	dbInstance.Create(&plan)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

// GetPlan 获取计划
func GetPlan(c *gin.Context) {
	var planId = c.Query("planId")
	var dbInstance = db.GetInstance()
	var plan orm.Plan

	dbInstance.Where("id = ?", planId).Find(&plan)

	err := dbInstance.Model(&plan).Association("Author").Find(&plan.Author)
	if err != nil {
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"plan": plan,
		},
	})
}

// GetPlanList 获取计划列表
func GetPlanList(c *gin.Context) {
	var offset, _ = strconv.Atoi(c.Query("offset"))
	var count, _ = strconv.Atoi(c.Query("count"))
	var name = c.Query("name")
	var state = c.Query("state")

	var dbInstance = db.GetInstance()
	var plans []orm.Plan
	dbInstance.Limit(count).Offset(offset).Where("name LIKE ?", "%"+name+"%").Where("state LIKE ?", "%"+state+"%").Preload("Author").Find(&plans)

	dtos := make([]types.PlanDTO, len(plans))
	for i := 0; i < len(plans); i++ {
		dtos[i].Plan = plans[i]
	}
	var total int64
	dbInstance.Model(&orm.Plan{}).Where("deleted_at IS NULL").Count(&total)

	var ids []interface{}
	for _, v := range plans {
		ids = append(ids, v.ID)
	}

	imagesMap := db.GetPlanFirstImagesByIds(ids)
	previews := make(map[int]string)
	for planId, images := range imagesMap {
		if _, ok := previews[planId]; ok {
			continue
		}
		var m []string
		err := json.Unmarshal([]byte(images), &m)
		if err != nil {
			continue
		}
		if len(m) > 0 {
			previews[planId] = m[0]
		}
	}

	for i := 0; i < len(dtos); i++ {
		dtos[i].Preview = previews[dtos[i].ID]
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"total": total,
			"plans": dtos,
		},
	})
}

// GetPlanDetail 获取计划详细信息
func GetPlanDetail(c *gin.Context) {
	var planId = c.Query("planId")
	var dbInstance = db.GetInstance()
	var plan orm.Plan

	dbInstance.Where("id = ?", planId).Find(&plan)

	// 连结用户信息
	dbInstance.Model(&plan).Association("Author").Find(&plan.Author)
	dbInstance.Model(&plan).Preload("Shows").Association("PlayPeriods").Find(&plan.PlayPeriods)
	devices := db.GetDevicesByPlanId(plan.ID)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"plan":    plan,
			"devices": devices,
		},
	})
}

// DeletePlan 删除某个计划
func DeletePlan(c *gin.Context) {
	var planId = c.Query("planId")
	var dbInstance = db.GetInstance()
	var plan orm.Plan
	dbInstance.Where("id = ?", planId).Find(&plan).Delete(&plan)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

// CopyPlan 复制一个计划
func CopyPlan(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	var user = t.(*jwtx.TokenUserInfo)

	var planId = c.Query("planId")
	var dbInstance = db.GetInstance()
	var plan orm.Plan
	// 首先获取原计划
	dbInstance.Where("id = ?", planId).Find(&plan)
	if plan.ID == 0 {
		c.JSON(400, gin.H{
			"error": "plan not exist!",
		})
		return
	}
	var newPlan orm.Plan
	newPlan.UserID, _ = strconv.Atoi(user.ID)
	newPlan.Name = plan.Name
	newPlan.Mode = plan.Mode
	newPlan.State = "未发布"

	dbInstance.Create(&newPlan)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

func PublishPlan(c *gin.Context) {
	var planId = c.Query("planId")
	var dbInstance = db.GetInstance()
	var plan orm.Plan
	dbInstance.Where("id = ?", planId).Find(&plan)
	if plan.ID == 0 {
		c.JSON(400, gin.H{
			"error": "plan not exist!",
		})
		return
	}
	devices := db.GetDevicesByPlanId(plan.ID)

	for _, d := range devices {
		onlineDevice := devicemanager.GetDeviceByMac(d.Mac)
		if onlineDevice == nil {
			c.JSON(200, gin.H{
				"code":    400,
				"message": "device is offline",
			})
			return
		}
		err := onlineDevice.SyncPlan([]*orm.Plan{&plan})
		if err != nil {
			c.JSON(200, gin.H{
				"code":    500,
				"message": "failed to publish",
			})
			db.GetInstance().Table("plan").Where("id = ", plan.ID).Updates(map[string]interface{}{"state": devicemanager.PublishFail})
			return
		}
	}

	db.GetInstance().Table("plan").Where("id = ", plan.ID).Updates(map[string]interface{}{"state": devicemanager.PublishSuccess})
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}
