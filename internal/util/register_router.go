package util

import (
	"boe-backend/internal/service"
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
	"log"
)

func RegisterRouter(r *gin.Engine) {
	// 路由未命中兜底
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})

	// 初始化JWT中间件
	authMiddleware, err := jwtx.GetAuthMiddleware()
	if err != nil {
		log.Fatal(err)
	}

	// === 用户相关路由 ===
	userRoute := r.Group("/user")
	userRoute.POST("/login", authMiddleware.LoginHandler)
	userRoute.POST("/create_account", service.CreateAccount)
	userRoute.POST("/get_users", service.GetUsers)

	// === 验证相关路由 ===
	auth := userRoute.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	// === 首页相关路由 ===
	homeRoute := r.Group("/home")
	homeRoute.Use(authMiddleware.MiddlewareFunc())
	// 首页所有信息
	homeRoute.GET("/group/count", service.GroupCountHandler)
	homeRoute.GET("/device/state", service.DevicesStateHandler)
	homeRoute.GET("/count", service.CountHandler)

	// 首页事件列表路由
	homeRoute.GET("/events", service.GetRecentEvents)

	// === 设备相关路由 ===
	deviceRoute := r.Group("/device")
	deviceRoute.Use(authMiddleware.MiddlewareFunc())
	deviceRoute.POST("", service.AddDeviceHandler)
	deviceRoute.GET("/all", service.GetDeviceListHandler)
	deviceRoute.GET("/info/:id", service.GetDeviceInfoHandler)
	deviceRoute.DELETE("/:id", service.DeleteDevice)

	// === 分组相关路由 ===
	groupRoute := r.Group("/group")
	groupRoute.Use(authMiddleware.MiddlewareFunc())
	groupRoute.GET("/all", service.GetGroupListHandler)
	groupRoute.GET("/devices/:id", service.GetGroupDevicesHandler)
	groupRoute.POST("", service.AddGroupHandler)
	groupRoute.DELETE("/:id", service.DeleteGroup)

	// === 计划相关路由 ===
	planRoute := r.Group("/plan")
	planRoute.Use(authMiddleware.MiddlewareFunc())
	planRoute.POST("/create", service.CreatePlan)
	planRoute.GET("/get_plan", service.GetPlan)
	planRoute.GET("/get_plan_list", service.GetPlanList)
	planRoute.GET("/get_plan_detail", service.GetPlanDetail)
	planRoute.DELETE("/delete_plan", service.DeletePlan)

	// === 文件相关路由 ===
	fileRoute := r.Group("file")
	fileRoute.Use(authMiddleware.MiddlewareFunc())
	fileRoute.POST("/presign", service.PreSignHandler)
}
