package msgtype

type PlanMsg struct {
	ID          int    `json:"id"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	PlayPeriods []struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
		LoopMode  string `json:"loopMode"`
		HTML      string `json:"html"`
	} `json:"playPeriods"`
}
