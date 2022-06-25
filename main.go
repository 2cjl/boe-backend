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

	http.HandleFunc("/ws", getWebSocketHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", websocketPort), nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func getWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
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
