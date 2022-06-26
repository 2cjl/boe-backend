package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"log"
)

type PlayPeriod struct {
	StartTime string
	EndTime   string
	LoopMode  string
	ShowIds   []string
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
	PlayPeriods []PlayPeriod
}

func CreatePlan(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	var user = t.(*jwtx.TokenUserInfo)

	var req CreatePlanRequest

	err := c.BindJSON(&req)
	if err != nil {
		log.Print(err)
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

	var playPeriods []orm.PlayPeriod

	for _, period := range req.PlayPeriods {
		var p orm.PlayPeriod
		p.StartTime = period.StartTime
		p.EndTime = period.EndTime
		p.LoopMode = period.LoopMode
		playPeriods = append(playPeriods, p)
	}

	plan.PlayPeriods = playPeriods

	var dbInstance = db.GetInstance()
	// 保存计划实体
	dbInstance.Create(&plan)

	// 保存时间段
	//var playPeriodList orm.PlayPeriod

	c.JSON(200, gin.H{
		"message": "success",
	})
}
