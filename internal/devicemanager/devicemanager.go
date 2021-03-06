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
	"github.com/patrickmn/go-cache"
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
	typeHello      = "hello"

	// backend->device
	typePong       = "pong"
	typePlanList   = "planList"
	typeDeletePlan = "deletePlan"
	typeHi         = "hi"
	typeBrightness = "brightness"
	typeScreenshot = "screenshot"

	writeTimeout = time.Second * 8

	DeviceOffline = "OFFLINE"
	DeviceOnline  = "ONLINE"

	PublishSuccess = "已发布"
	PublishFail    = "发布失败"
	Unpublished    = "未发布"
)

var (
	devices     = make(map[string]*Device)
	Screenshots = cache.New(time.Minute, 2*time.Minute)
	deviceLock  sync.Mutex
)

func GetDeviceByMac(mac string) *Device {
	deviceLock.Lock()
	d := devices[mac]
	deviceLock.Unlock()
	return d
}

type Device struct {
	ID             int
	DeviceName     string
	OrganizationID int
	Mac            string

	LastHeartbeat time.Time
	RunningTime   uint64
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

	var plans []*orm.Plan
	ins := db.GetInstance()
	ins.Where("id in (?)", ins.Table("plan_device").Select("plan_id").Where("device_id = ?", d.ID)).Where("state = ? OR state = ?", PublishFail, PublishSuccess).Find(&plans)
	log.Println(plans)
	err = d.SyncPlan(plans)
	var ids []int
	for _, p := range plans {
		ids = append(ids, p.ID)
	}
	if err != nil {
		log.Println(err)
		ins.Table("plan").Where("id IN ?", ids).Updates(map[string]interface{}{"state": PublishFail})
		return
	}
	ins.Table("plan").Where("id IN ?", ids).Updates(map[string]interface{}{"state": PublishSuccess})
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
			d.RunningTime = uint64(m["runningTime"].(float64))
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
		case typeScreenshot:
			Screenshots.Set(d.Mac, m["data"].(string), time.Minute)
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

func (d *Device) ChangeBrightness(data float64) error {
	m := make(map[string]interface{})
	m["type"] = typeBrightness
	m["data"] = data
	return d.writeMsg(m)
}

func (d *Device) CtlScreenshot() error {
	m := make(map[string]interface{})
	m["type"] = typeScreenshot
	return d.writeMsg(m)
}

func (d *Device) SyncPlan(plans []*orm.Plan) error {
	if d.ID == 0 {
		return errors.New(fmt.Sprintf("device(%s)id is 0!!!\n", d.Mac))
	}

	var planMsgList []*types.PlanMsg
	result := map[string]interface{}{
		"type": typePlanList,
	}

	ins := db.GetInstance()

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
	result["plan"] = planMsgList
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
	log.Println(string(msg))
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
