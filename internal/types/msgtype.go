package types

type PlanMsg struct {
	ID          int             `json:"id"`
	StartDate   string          `json:"startDate"`
	EndDate     string          `json:"endDate"`
	Mode        string          `json:"mode"`
	PlayPeriods []PlayPeriodMsg `json:"playPeriods"`
}
type PlayPeriodMsg struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	LoopMode  string `json:"loopMode"`
	HTML      string `json:"html"`
}
