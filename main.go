package main

import (
	"boe-backend/internal/devicemanager"
	"boe-backend/internal/util"
	"boe-backend/internal/util/config"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	port          int
	websocketPort int
	upgrader      = websocket.Upgrader{}
	devices       = make(map[string]*devicemanager.Device)
	deviceLock    sync.Mutex
)

func init() {
	flag.IntVar(&port, "port", 8080, "")
	flag.IntVar(&websocketPort, "websocket-port", 8081, "")
	flag.Parse()
	config.InitViper()
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(util.Cors())

	// 注册所有路由
	util.RegisterRouter(r)

<<<<<<< HEAD
	http.HandleFunc("/ws", getWebSocketHandler)

	// users
	userRoute := r.Group("/user")
	userRoute.POST("/login", authMiddleware.LoginHandler)
	userRoute.POST("/register", registerHandler)
	// 一组需要验证的路由
	auth := userRoute.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	homeRoute := r.Group("/home")
	homeRoute.Use(authMiddleware.MiddlewareFunc())

	// 首页事件列表路由
	homeRoute.GET("/events", func(context *gin.Context) {
		var organizationId = context.Query("organizationId")
		var events = db.GetAllEvents(organizationId)
		context.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"data": gin.H{
				"events": events,
			},
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	err = http.ListenAndServe(fmt.Sprintf(":%d", websocketPort), nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func getWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
=======
	// websocket 相关
	r.GET("/ws", getWebSocketHandler)
	util.WatchSignalGrace(r, port)
}

func getWebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
>>>>>>> origin/feat/yuzhanglong
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	_, msg, err := conn.ReadMessage()
	m := make(map[string]interface{})
	err = json.Unmarshal(msg, &m)
	if err != nil || m["type"] != "hello" || m["mac"] == nil || devices[m["mac"].(string)] != nil {
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
