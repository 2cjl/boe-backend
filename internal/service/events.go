package service

import (
	"boe-backend/internal/db"
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"log"
)

func GetAllEvents(context *gin.Context) {
	t, _ := context.Get(jwtx.IdentityKey)
	var user = t.(*jwtx.TokenUserInfo)
	log.Print(user)
	var events = db.GetAllEvents(user.OrganizationID)
	context.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"events": events,
		},
	})
}
