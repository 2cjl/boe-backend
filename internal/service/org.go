package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetSelfOrgHandler(c *gin.Context) {
	t, _ := c.Get(jwtx.IdentityKey)
	info := t.(*jwtx.TokenUserInfo)

	var org orm.Organization
	id, err := strconv.Atoi(info.OrganizationID)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  400,
			"error": "Bad request parameter",
		})
		return
	}
	org.ID = id
	db.GetInstance().Find(&org)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    org,
	})
}
