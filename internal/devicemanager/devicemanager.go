package devicemanager

import (
	"boe-backend/internal/db"
	"boe-backend/internal/orm"
	"boe-backend/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm/clause"
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

	DeviceOffline = "OFFLINE"
	DeviceOnline  = "ONLINE"
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

	// 获取设备信息
	data := make(map[string]interface{})
	data["type"] = typeDeviceInfo
	err := d.writeMsg(data)
	if err != nil {
		log.Println(err)
		return
	}
	err = d.SyncPlan()
	if err != nil {
		log.Println(err)
	}
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
		//d.Conn.SetReadDeadline(time.Now().Add(readTimeout))
		_, msg, err := d.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("device(%s) unexpected close: %v", d.Mac, err)
				return
			}
			log.Printf("device(%s)read for ws fail: %v\n", d.Mac, err)
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
			var mac string
			if v, ok := m["mac"]; ok && v != nil {
				mac = v.(string)
			} else {
				result["msg"] = "fail: mac is empty"
				log.Println("fail: mac is empty")
				return
			}
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
			var info orm.DeviceInfo
			err := mapstructure.Decode(m["info"], &info)
			if err != nil {
				log.Println(err)
				continue
			}
			if d.ID == 0 {
				continue
			}
			info.ID = d.ID
			info.LastHeartbeat = time.Now()
			db.GetInstance().Create(&info)
			db.GetInstance().Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&info)

			continue
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

func (d *Device) SyncPlan() error {
	var plans []*orm.Plan
	var planMsgList []*types.PlanMsg
	result := map[string]interface{}{
		"type": typePlanList,
		"plan": planMsgList,
	}

	ins := db.GetInstance()
	// 获取plan
	if d.ID == 0 {
		log.Printf("device(%s)id is 0!!!\n", d.Mac)
		return errors.New(fmt.Sprintf("device(%s)id is 0!!!\n", d.Mac))
	}
	ins.Where("id in (?)", ins.Table("plan_device").Select("plan_id").Where("device_id = ?", d.ID)).Find(&plans)
	// 对于每个plan获取PlayPeriods,并构造返回值
	for _, plan := range plans {
		err := ins.Model(&plan).Preload("Shows").Association("PlayPeriods").Find(&plan.PlayPeriods)
		if err != nil {
			log.Println(err)
			return err
		}
		var playPeriodMsgList []types.PlayPeriodMsg
		for _, v := range plan.PlayPeriods {
			var p types.PlayPeriodMsg
			p.HTML = v.Html
			p.StartTime = v.StartTime
			p.EndTime = v.EndTime
			p.LoopMode = v.LoopMode
			playPeriodMsgList = append(playPeriodMsgList, p)
		}
		msg := &types.PlanMsg{}
		msg.ID = plan.ID
		msg.Mode = plan.Mode
		msg.StartDate = plan.StartDate
		msg.EndDate = plan.EndDate
		msg.PlayPeriods = playPeriodMsgList

		planMsgList = append(planMsgList, msg)
	}
	log.Println("len:", len(planMsgList))

	log.Println(planMsgList[0])
	return d.writeMsg(result)
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
