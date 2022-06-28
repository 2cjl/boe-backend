package types

type BanUserRequest struct {
	UserId string
}

type DeleteUserRequest struct {
	UserId string
}

type CreateAccountForm struct {
	// 账户名
	Username string `form:"username" json:"username" binding:"required"`
	// 密码
	Passwd string `form:"passwd" json:"passwd" binding:"required"`
	// 所属机构
	Organization int `form:"organization" json:"organization" binding:"required"`
	// 账号状态
	Status string `form:"status" json:"status" binding:"required"`
	// 真实姓名
	RealName string `form:"realName" json:"realName" binding:"required"`
	// 邮箱
	Email string `form:"email" json:"email" binding:"required"`
	// 手机号
	Phone string `form:"phone" json:"phone" binding:"required"`
}

type UpdateAccountForm struct {
	UserId string `form:"userId" json:"userId" binding:"required"`
	// 账户名
	Username string `form:"username" json:"username"`
	// 密码
	Passwd string `form:"passwd" json:"passwd"`
	// 所属机构
	Organization int `form:"organization" json:"organization"`
	// 账号状态
	Status string `form:"status" json:"status"`
	// 真实姓名
	RealName string `form:"realName" json:"realName"`
	// 邮箱
	Email string `form:"email" json:"email"`
	// 手机号
	Phone string `form:"phone" json:"phone"`
}
