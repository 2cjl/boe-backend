package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CreateAccount(c *gin.Context) {
	var registerForm jwtx.CreateAccountForm
	if err := c.ShouldBind(&registerForm); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request parameter",
		})
		return
	}
	var probablySameUser = GetUserByUserName(registerForm.Username)
	if probablySameUser.ID != 0 {
		c.JSON(400, gin.H{
			"error": "该用户名已存在",
		})
		return
	}

	var user orm.User
	user.Username = registerForm.Username
	user.Passwd = registerForm.Passwd
	user.OrganizationID = registerForm.Organization
	user.Email = registerForm.Email
	user.Phone = registerForm.Phone
	user.RealName = registerForm.RealName
	user.Status = registerForm.Status
	var dbInstance = db.GetInstance()
	dbInstance.Save(&user)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

// GetUserByUserName 根据用户名获取用户
func GetUserByUserName(name string) orm.User {
	var ins = db.GetInstance()
	var user orm.User
	ins.First(&user, "username = ?", name)

	return user
}

// GetUserById 根据用户 ID 获取用户
func GetUserById(id string) orm.User {
	var ins = db.GetInstance()
	var user orm.User
	ins.First(&user, id)
	return user
}

// GetUsers 获取所有用户（分页）
func GetUsers(c *gin.Context) {
	var offset, _ = strconv.Atoi(c.Query("offset"))
	var count, _ = strconv.Atoi(c.Query("count"))
	var users []orm.User

	var dbInstance = db.GetInstance()
	dbInstance.Limit(count).Offset(offset).Find(&users)

	var total int64
	dbInstance.Model(&orm.User{}).Count(&total)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"total":   total,
		"users":   users,
	})
}
