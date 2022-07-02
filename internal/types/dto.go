package types

import "boe-backend/internal/orm"

type PlanDTO struct {
	orm.Plan
	Preview string
}

type DeviceDTO struct {
	ID         int
	Name       string
	Mac        string
	Resolution string
	PlanName   string
	State      string
}

type ShowDTO struct {
	orm.Show
	Preview string
}

type GroupDTO struct {
	orm.Group
	DeviceCnt int
}
