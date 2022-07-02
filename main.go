package main

import (
	"boe-backend/internal/devicemanager"
	"boe-backend/internal/util"
	"boe-backend/internal/util/config"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	port          int
	websocketPort int
	upgrader      = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
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
	//r.Use(util.Cors())

	// 注册所有路由
	util.RegisterRouter(r)

	http.HandleFunc("/", getWebSocketHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	err := http.ListenAndServe(fmt.Sprintf(":%d", websocketPort), nil)
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
	device := &devicemanager.Device{
		Conn: conn,
	}
	device.Receive()
}
