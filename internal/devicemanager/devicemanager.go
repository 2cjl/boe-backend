package devicemanager

import (
	"boe-backend/internal/db"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

const (
	// device->backend
	typePing       = "ping"
	typeDeviceInfo = "device_info"
	typeSyncPlan   = "sync_plan"

	// backend->device
	typePong       = "pong"
	typePlanList   = "plan_list"
	typeDeletePlan = "delete_plan"
)

type Device struct {
	ID            string
	DeviceName    string
	Organization  string
	Mac           string
	LastHeartbeat time.Time
	RunningTime   int
	PlanID        int

	conn *websocket.Conn
}

func (d *Device) Init(conn *websocket.Conn) {
	d.conn = conn
	device := db.GetDeviceByMac(d.Mac)
	if device == nil {
		return
	}
	d.ID = strconv.Itoa(device.ID)
}

func (d *Device) Receive(closeFun func()) {
	defer closeFun()
	defer d.conn.Close()
	for {
		_, msg, err := d.conn.ReadMessage()
		if err != nil {
			// 网络错误则直接返回，等待客户端重连
			log.Printf("device(%s)read for ws fail: %s\n", d.Mac, err.Error())
			return
		}
		log.Printf("recv: %s", msg)
		m := make(map[string]interface{})
		err = json.Unmarshal(msg, &m)
		if err != nil {
			// 忽略错误格式的数据
			log.Println(err)
			continue
		}
		var result map[string]interface{}
		switch m["type"].(string) {
		case typePing:
			result = map[string]interface{}{
				"type": typePong,
			}
			d.LastHeartbeat = time.Now()
			d.RunningTime = m["running_time"].(int)
		case typeDeviceInfo:
			///TODO(vincent)获取设备信息，同步到数据库
		case typeSyncPlan:
			///TODO(vincent)从数据库筛选未安排的计划，返回给设备
		default:
			log.Printf("unknown type:%s\n", m["type"].(string))
			continue
		}
		msg, err = json.Marshal(result)
		if err != nil {
			log.Println(err)
			continue
		}
		err = d.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("device(%s)write for ws fail: %s\n", d.Mac, err.Error())
			return
		}
	}
}
