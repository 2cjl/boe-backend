package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	"boe-backend/internal/types"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CreateAccount(c *gin.Context) {
	var registerForm types.CreateAccountForm
	if err := c.ShouldBind(&registerForm); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request parameter",
		})
		return
	}
	var probablySameUser = GetUserByUserName(registerForm.Username)
	if probablySameUser.ID != 0 {
		c.JSON(400, gin.H{
			"error": "The username already exists",
		})
		return
	}

	organization := db.GetOrganizationById(registerForm.Organization)
	if organization == nil {
		c.JSON(400, gin.H{
			"error": "nonexistent organization",
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
	var name = c.Query("name")

	var users []orm.User

	var dbInstance = db.GetInstance()
	dbInstance.Limit(count).Offset(offset).Preload("Organization").Where("username LIKE ?", "%"+name+"%").Find(&users)

	var total int64
	dbInstance.Model(&orm.User{}).Where("username LIKE ?", "%"+name+"%").Where("deleted_at IS NULL").Count(&total)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"total":   total,
		"users":   users,
	})
}

func BanUser(c *gin.Context) {
	var req types.BanUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "request param error",
		})
		return
	}
	var user = GetUserById(req.UserId)
	if user.ID == 0 {
		c.JSON(400, gin.H{
			"error": "user not existed",
		})
		return
	}
	user.Status = "停用"
	var dbInstance = db.GetInstance()
	dbInstance.Save(&user)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

func DeleteUser(c *gin.Context) {
	var req types.DeleteUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "request param error",
		})
		return
	}
	var user = GetUserById(req.UserId)
	if user.ID == 0 {
		c.JSON(400, gin.H{
			"error": "user not existed",
		})
		return
	}
	var dbInstance = db.GetInstance()
	dbInstance.Delete(&user)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

func UpdateAccount(c *gin.Context) {
	var updateAccountFormData types.UpdateAccountForm
	if err := c.ShouldBind(&updateAccountFormData); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request parameter",
		})
		return
	}

	if updateAccountFormData.Organization != 0 {
		organization := db.GetOrganizationById(updateAccountFormData.Organization)
		if organization == nil {
			c.JSON(400, gin.H{
				"error": "nonexistent organization",
			})
			return
		}
	}

	// 获取待编辑的用户
	var user = GetUserById(updateAccountFormData.UserId)
	if user.ID == 0 {
		c.JSON(400, gin.H{
			"error": "user not existed",
		})
		return
	}

	// 当用户名被修改时，去校验
	if user.Username != updateAccountFormData.Username {
		var probablySameUser = GetUserByUserName(updateAccountFormData.Username)
		if probablySameUser.ID != 0 {
			c.JSON(400, gin.H{
				"error": "The username already exists",
			})
			return
		}
	}

	// 开始执行更新操作 (只会更新非 0 字段)
	user.Username = updateAccountFormData.Username
	user.Passwd = updateAccountFormData.Passwd
	user.OrganizationID = updateAccountFormData.Organization
	user.Email = updateAccountFormData.Email
	user.Phone = updateAccountFormData.Phone
	user.RealName = updateAccountFormData.RealName
	user.Status = updateAccountFormData.Status
	var dbInstance = db.GetInstance()
	dbInstance.Updates(&user)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}
