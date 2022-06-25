package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
)

type PlayPeriod struct {
}

type CreatePlanRequest struct {
	// 名称
	Name string
	// 播放模式
	Mode string
	// 开始时间
	StartDate string
	// 结束时间
	EndDate string
	// 该计划对应的时间段
	PlayPeriod PlayPeriod
}

func CreatePlan(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	var user = t.(*jwtx.TokenUserInfo)

	var req CreatePlanRequest

	err := c.BindJSON(&req)
	if err != nil {
		return
	}

	var plan orm.Plan
	plan.Author = user.ID
	plan.EndDate = req.EndDate
	plan.StartDate = req.StartDate
	plan.Name = req.Name
	plan.Mode = req.Mode
	// 初始设置为未发布状态
	plan.State = "NOT_RELEASED"

	var dbInstance = db.GetInstance()
	dbInstance.Create(&plan)

	c.JSON(200, gin.H{
		"message": "success",
	})
}
