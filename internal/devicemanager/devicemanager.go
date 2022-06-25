package devicemanager

import (
	"boe-backend/internal/db"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
	"time"
)

const (
	// device->backend
	typePing       = "ping"
	typeDeviceInfo = "device_info"
	typeSyncPlan   = "sync_plan"
	typeHello      = "hello"

	// backend->device
	typePong       = "pong"
	typePlanList   = "plan_list"
	typeDeletePlan = "delete_plan"
	typeHi         = "hi"
)

var (
	devices    = make(map[string]*Device)
	deviceLock sync.Mutex
)

type Device struct {
	ID           string
	DeviceName   string
	Organization string
	Mac          string

	LastHeartbeat time.Time
	RunningTime   int
	PlanID        int

	Conn *websocket.Conn
}

// InitInfo 初始化设备信息
func (d *Device) InitInfo() {
	log.Printf("device(%s) init", d.Mac)

	device := db.GetDeviceByMac(d.Mac)
	if device == nil {
		return
	}
	d.ID = strconv.Itoa(device.ID)
	d.DeviceName = device.Name
	organization := db.GetOrganizationById(device.OrganizationID)
	if organization != nil {
		d.Organization = organization.Name
	}
}

func (d *Device) Receive() {
	defer func() {
		log.Println("delete")
		deviceLock.Lock()
		delete(devices, d.Mac)
		deviceLock.Unlock()
	}()
	defer d.Conn.Close()
	for {
		_, msg, err := d.Conn.ReadMessage()
		log.Println("err:", err)
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
		case typeHello:
			log.Printf("hello")
			result = map[string]interface{}{
				"type": typeHi,
			}

			mac := m["mac"].(string)
			d.Mac = mac
			deviceLock.Lock()
			if devices[mac] != nil {
				result["msg"] = "fail: mac conflict"
				deviceLock.Unlock()
				log.Println("fail: mac conflict")
				return
			}
			devices[mac] = d
			deviceLock.Unlock()
			result["msg"] = "ok"
			d.InitInfo()

		case typePing:
			result = map[string]interface{}{
				"type": typePong,
			}
			d.LastHeartbeat = time.Now()
			d.RunningTime = m["running_time"].(int)
			d.PlanID = m["plan_id"].(int)
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
		err = d.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("device(%s)write for ws fail: %s\n", d.Mac, err.Error())
			return
		}
	}
}
