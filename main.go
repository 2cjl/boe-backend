package main

import (
	"boe-backend/internal/devicemanager"
	"boe-backend/internal/util"
	"boe-backend/internal/util/config"
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
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(util.Cors())

	// 注册所有路由
	util.RegisterRouter(r)

	// websocket 相关
	r.GET("/ws", getWebSocketHandler)
	util.WatchSignalGrace(r, port)
}

func getWebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
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
