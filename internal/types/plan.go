package types

type PlayPeriod struct {
	StartTime string
	EndTime   string
	LoopMode  string
	Html      string
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
