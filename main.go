package main

import (
	"boe-backend/internal/util"
	"boe-backend/internal/util/config"
	jwtx "boe-backend/internal/util/jwt"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	flag.Parse()
	config.InitViper()
}

func main() {
	r := gin.New()
	//gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})

	// 初始化JWT中间件
	authMiddleware, err := jwtx.GetAuthMiddleware()
	if err != nil {
		log.Fatal(err)
	}

	// users
	router := r.Group("/user")
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/register", registerHandler)
	// 一组需要验证的路由
	auth := router.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	util.WatchSignalGrace(r, 8080)
}

func registerHandler(c *gin.Context) {
	var registerForm jwtx.RegisterForm
	if err := c.ShouldBind(&registerForm); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request parameter",
		})
		return
	}

	//isExist, err := db.IsExistPhone(registerForm.Phone)
	//if err != nil {
	//	log.Println(err)
	//	c.JSON(500, gin.H{
	//		"error": "Server internal error",
	//	})
	//	return
	//}
	//if isExist {
	//	c.JSON(400, gin.H{
	//		"error": "phone already exists",
	//	})
	//	return
	//}
	//
	//id, err := db.Register(registerForm.Username, registerForm.Phone, registerForm.Passwd, registerForm.IdNumber, registerForm.WorkStatus, registerForm.Age)
	//if err != nil {
	//	log.Println(err)
	//	c.JSON(500, gin.H{
	//		"error": "Server internal error",
	//	})
	//	return
	//}

	//fmt.Printf("id (%v, %T)\n", id, id)

	c.JSON(200, gin.H{
		//"token": jwtx.GenerateToken(id),
	})
}
