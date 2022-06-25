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
	typeDeviceInfo = "deviceInfo"
	typeSyncPlan   = "syncPlan"
	typeHello      = "hello"

	// backend->device
	typePong       = "pong"
	typePlanList   = "planList"
	typeDeletePlan = "deletePlan"
	typeHi         = "hi"

	writeTimeout = time.Second * 8
	readTimeout  = time.Second * 8
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
	RunningTime   float64
	PlanID        float64

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
		d.Conn.SetReadDeadline(time.Now().Add(readTimeout))
		_, msg, err := d.Conn.ReadMessage()
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
			d.RunningTime = m["runningTime"].(float64)
			d.PlanID = m["planId"].(float64)
		case typeDeviceInfo:
			///TODO(vincent)获取设备信息，同步到数据库
		case typeSyncPlan:
			///TODO(vincent)从数据库筛选未安排的计划，返回给设备
		default:
			log.Printf("unknown type:%s\n", m["type"].(string))
			continue
		}
		err = d.writeMsg(result)
		if err != nil {
			log.Printf("device(%s)write for ws fail: %s\n", d.Mac, err.Error())
			return
		}
	}
}

func (d *Device) writeMsg(data map[string]interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	d.Conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	err = d.Conn.WriteMessage(websocket.TextMessage, msg)
	return err
}
