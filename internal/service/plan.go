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
	ShowIds   []int
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

	// 开始初始化各个时间段，并复制
	var playPeriods []orm.PlayPeriod

	// 在内存中先初始化实体
	for _, period := range req.PlayPeriods {
		var p orm.PlayPeriod
		p.StartTime = period.StartTime
		p.EndTime = period.EndTime
		p.LoopMode = period.LoopMode

		var shows []orm.Show
		var ins = db.GetInstance()
		// PREF: 在 for 循环之外合并所有传入的 showID 然后通过一次查询的查到所有用的 show
		ins.Find(&shows, period.ShowIds)
		p.Shows = shows

		playPeriods = append(playPeriods, p)
	}

	plan.PlayPeriods = playPeriods

	var dbInstance = db.GetInstance()
	// 保存计划实体
	dbInstance.Create(&plan)

	c.JSON(200, gin.H{
		"message": "success",
	})
}
