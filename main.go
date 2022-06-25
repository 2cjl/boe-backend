package main

import (
	"boe-backend/internal/devicemanager"
	"boe-backend/internal/util"
	"boe-backend/internal/util/config"
	jwtx "boe-backend/internal/util/jwt"
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var (
	port       int
	upgrader   = websocket.Upgrader{}
	devices    = make(map[string]*devicemanager.Device)
	deviceLock sync.Mutex
)

func init() {
	flag.IntVar(&port, "port", 8080, "")
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

	r.GET("/ws", getWebSocketHandler)

	// users
	router := r.Group("/user")
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/register", registerHandler)
	// 一组需要验证的路由
	auth := router.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	util.WatchSignalGrace(r, port)
}

func getWebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "server internal error",
		})
		conn.Close()
		return
	}
	_, msg, err := conn.ReadMessage()
	m := make(map[string]interface{})
	err = json.Unmarshal(msg, &m)
	if err != nil || m["type"] != "hello" || m["mac"] == nil || devices[m["mac"].(string)] != nil {
		c.JSON(200, gin.H{
			"code":  "500",
			"error": "hello msg error",
		})
		conn.Close()
		return
	}
	mac := m["mac"].(string)
	device := &devicemanager.Device{
		Mac: mac,
	}
	deviceLock.Lock()
	devices[mac] = device
	deviceLock.Unlock()

	device.Init(conn)
	go device.Receive(func() {
		deviceLock.Lock()
		delete(devices, mac)
		deviceLock.Unlock()
	})
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
