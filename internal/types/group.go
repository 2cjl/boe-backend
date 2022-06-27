package types

type AddGroupReq struct {
	Name           string `binding:"required"`
	Describe       string `binding:"required"`
	OrganizationID int    `binding:"required"`
	Devices        []int  `binding:"required"`
}
