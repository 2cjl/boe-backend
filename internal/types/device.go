package types

type AddDeviceReq struct {
	Name string `binding:"required"`
	Mac  string `binding:"required"`
}
