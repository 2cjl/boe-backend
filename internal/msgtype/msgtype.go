package msgtype

type DeviceInfoMsg struct {
	VersionApp  string  `json:"versionApp"`
	Memory      string  `json:"memory"`
	HardwareID  string  `json:"hardwareId"`
	IP          string  `json:"ip"`
	Model       string  `json:"model"`
	Storage     string  `json:"storage"`
	RunningTime float64 `json:"runningTime"`
	Resolution  string  `json:"resolution"`
	Mac         string  `json:"mac"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

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
