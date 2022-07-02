package service

import (
	"boe-backend/internal/db"
	"boe-backend/internal/devicemanager"
	"boe-backend/internal/orm"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type reqBrightness struct {
	data float64
}

func CtlScreenshotHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var device orm.Device
	device.ID = id
	db.GetInstance().Find(&device)
	if d := devicemanager.GetDeviceByMac(device.Mac); d != nil {
		err := d.CtlScreenshot()
		if err != nil {
			log.Println(err)
			c.JSON(200, gin.H{
				"code":    500,
				"message": "send msg fail",
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
		})
	}
	c.JSON(200, gin.H{
		"code":    404,
		"message": "device offline",
	})
}

func GetScreenshotHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var device orm.Device
	device.ID = id
	db.GetInstance().Find(&device)

	if data, ok := devicemanager.Screenshots.Get(device.Mac); ok {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"data":    data.(string),
		})
	} else {
		c.JSON(200, gin.H{
			"code":    404,
			"message": "device offline",
		})
	}
}

func ChangeBrightnessHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req reqBrightness
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  400,
			"error": "Bad request parameter",
		})
		return
	}

	var device orm.Device
	device.ID = id
	db.GetInstance().Find(&device)
	if d := devicemanager.GetDeviceByMac(device.Mac); d != nil {
		err := d.ChangeBrightness(req.data)
		if err != nil {
			log.Println(err)
			c.JSON(200, gin.H{
				"code":    500,
				"message": "send msg fail",
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
		})
	}
	c.JSON(200, gin.H{
		"code":    404,
		"message": "device offline",
	})
}
