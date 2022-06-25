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
	userRoute.POST("/register", service.RegisterHandler)

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
	// 首页事件列表路由
	homeRoute.GET("/events", service.GetAllEvents)
}
