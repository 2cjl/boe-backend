package devicemanager

import (
	"boe-backend/internal/db"
	"boe-backend/internal/msgtype"
	"boe-backend/internal/orm"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
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

	DeviceOffline = "offline"
	DeviceOnline  = "online"
)

var (
	devices    = make(map[string]*Device)
	deviceLock sync.Mutex
)

type Device struct {
	ID             int
	DeviceName     string
	OrganizationID int
	Mac            string

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
	d.ID = device.ID
	d.DeviceName = device.Name
	organization := db.GetOrganizationById(device.OrganizationID)
	if organization != nil {
		d.OrganizationID = organization.ID
	}

	// 新增事件
	e := NewDeviceEvent(device.Name, device.OrganizationID, true)
	db.GetInstance().Create(e)
	db.GetInstance().Model(device).Update("state", DeviceOnline)
}

func (d *Device) Receive() {
	defer func() {
		log.Println("delete")
		deviceLock.Lock()
		delete(devices, d.Mac)
		deviceLock.Unlock()

		// id为0 则为未注册设备
		if d.ID != 0 {
			e := NewDeviceEvent(d.DeviceName, d.OrganizationID, false)
			db.GetInstance().Create(e)
			// 更新状态
			db.GetInstance().Model(&orm.Device{}).Where("id = ?", d.ID).Updates(orm.Device{State: DeviceOffline, PlanID: d.PlanID})
			db.GetInstance().Model(&orm.DeviceInfo{ID: d.ID}).Updates(orm.DeviceInfo{LastHeartbeat: d.LastHeartbeat, RunningTime: d.RunningTime})
		}
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
			d.RunningTime = int(m["runningTime"].(float64))
			d.PlanID = int(m["planId"].(float64))
		case typeDeviceInfo:
			var info msgtype.DeviceInfoMsg
			err := mapstructure.Decode(m["info"], &info)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(info)
			continue
		case typeSyncPlan:
			result = map[string]interface{}{
				"type": typePlanList,
				"plan": []int{},
			}
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

func (d *Device) DeletePlan(planIds []int) error {
	data := make(map[string]interface{})
	data["type"] = typeDeletePlan
	data["planIds"] = planIds
	return d.writeMsg(data)
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

func NewDeviceEvent(deviceName string, organizationID int, isOnline bool) *orm.Event {
	e := &orm.Event{
		Time:           time.Now(),
		ObjectType:     "设备",
		OrganizationId: strconv.Itoa(organizationID),
	}
	if isOnline {
		e.Content = deviceName + " 已上线"
	} else {
		e.Content = deviceName + " 已下线"
	}
	return e
}
