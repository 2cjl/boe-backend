package service

import (
	"boe-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func GetAllEvents(context *gin.Context) {
	var organizationId = context.Query("organizationId")
	var events = db.GetAllEvents(organizationId)
	context.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"events": events,
		},
	})
}
