package types

type AddGroupReq struct {
	Name     string `binding:"required"`
	Describe string `binding:"required"`
	Devices  []int  `binding:"required"`
}
