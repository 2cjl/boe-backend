package miniox

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	"boe-backend/internal/types"
	"boe-backend/internal/util/config"
	"encoding/json"
	"log"
	"testing"
)

func TestMinio(t *testing.T) {
	config.InitViper()
	result := map[string]interface{}{
		"type": "1",
	}
	//var planList types.PlanMsg
	var plans []*orm.Plan
	var planMsgList []*types.PlanMsg
	ins := db.GetInstance()
	// 获取plan
	ins.Where("id in (?)", ins.Table("plan_device").Select("plan_id").Where("device_id = ?", 1)).Find(&plans)
	// 对于每个plan获取PlayPeriods,并构造返回值
	for _, plan := range plans {
		err := ins.Model(&plan).Preload("Shows").Association("PlayPeriods").Find(&plan.PlayPeriods)
		if err != nil {
			log.Println(err)
			continue
		}
		var playPeriodMsgList []types.PlayPeriodMsg
		for _, v := range plan.PlayPeriods {
			var p types.PlayPeriodMsg
			p.HTML = v.Html
			p.StartTime = v.StartTime
			p.EndTime = v.EndTime
			p.LoopMode = v.LoopMode
			playPeriodMsgList = append(playPeriodMsgList, p)
		}
		msg := &types.PlanMsg{}
		msg.ID = plan.ID
		msg.Mode = plan.Mode
		msg.StartDate = plan.StartDate
		msg.EndDate = plan.EndDate
		msg.PlayPeriods = playPeriodMsgList

		planMsgList = append(planMsgList, msg)
	}
	result["plan"] = planMsgList

	log.Println(result)
	msg, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(msg))

}
