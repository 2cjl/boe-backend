package types

type AddDeviceReq struct {
	ID   int
	Name string `binding:"required"`
	Mac  string `binding:"required"`
}
