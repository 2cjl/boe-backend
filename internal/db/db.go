package db

import (
	"boe-backend/internal/orm"
	"boe-backend/internal/types"
	"boe-backend/internal/util/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func getInstance() {
	if db == nil {
		once.Do(func() {
			cfg := config.GetConfig().Mysql
			source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Passwd, cfg.Host, cfg.Port, cfg.Database)
			_db, err := gorm.Open(mysql.Open(source), &gorm.Config{})
			db = _db
			if err != nil {
				panic(err)
			}
		})
	}
}

func GetInstance() *gorm.DB {
	getInstance()
	return db
}

func Login(phone, passwd string) *orm.User {
	getInstance()
	var u orm.User
	db.Where("phone = ? AND passwd = ?", phone, passwd).First(&u)
	if u.ID == 0 {
		return nil
	}
	return &u
}

func GetDeviceByMac(mac string) *orm.Device {
	getInstance()
	var d orm.Device
	db.Where("mac = ?", mac).First(&d)
	if d.ID == 0 {
		return nil
	}
	return &d
}

func GetOrganizationById(id int) *orm.Organization {
	getInstance()
	var o orm.Organization
	db.First(&o, id)
	if o.ID == 0 {
		return nil
	}
	return &o
}

func GetOrganizationByUser(uid int) *orm.Organization {
	getInstance()
	var o orm.Organization
	db.First(&o, db.Table("users").Select("organization_id").Where("id = ?", uid))
	if o.ID == 0 {
		return nil
	}
	return &o
}

// GetRecentEvents 获取最近的30条事件
func GetRecentEvents(organizationId string) []orm.Event {
	getInstance()
	var events []orm.Event
	db.Limit(30).Order("time desc").Find(&events, "organization_id = ?", organizationId)
	return events
}

type GroupCnt struct {
	ID   int
	Name string
	Cnt  int
}

func GetGroupDeviceCnt(organizationId string) []GroupCnt {
	getInstance()
	rows, err := db.Raw("select group_id id, g.name name, count(*) cnt FROM group_device, `groups` g  WHERE g.id = group_device.group_id AND g.organization_id = ? GROUP BY group_id", organizationId).Rows()

	if err != nil {
		log.Println(err)
		return nil
	}
	var c []GroupCnt
	for rows.Next() {
		var g GroupCnt
		err := rows.Scan(&g.ID, &g.Name, &g.Cnt)
		if err != nil {
			log.Println(err)
			return nil
		}
		c = append(c, g)
	}
	return c
}

func GetGroupDeviceCntByGroup(groupIdList []orm.Group) []GroupCnt {
	if groupIdList == nil || len(groupIdList) == 0 {
		return []GroupCnt{}
	}
	var ids []interface{}
	for _, v := range groupIdList {
		ids = append(ids, v.ID)
	}

	getInstance()
	rows, err := db.Raw(`select group_id id, g.name name, count(*) cnt FROM group_device, groups g  WHERE g.id = group_device.group_id AND g.id in (?`+strings.Repeat(",?", len(ids)-1)+`)`+` GROUP BY group_id`, ids...).Rows()

	if err != nil {
		log.Println(err)
		return nil
	}
	var c []GroupCnt
	for rows.Next() {
		var g GroupCnt
		err := rows.Scan(&g.ID, &g.Name, &g.Cnt)
		if err != nil {
			log.Println(err)
			return nil
		}
		c = append(c, g)
	}
	return c
}

func GetPlanByIds(ids []int) []orm.Plan {
	getInstance()
	var plans []orm.Plan
	if len(ids) == 0 {
		ids = []int{0}
	}
	db.Find(&plans, ids)
	return plans
}

func GetDevicesByGroupDevice(groupDevices []orm.GroupDevice) []orm.Device {
	if groupDevices == nil || len(groupDevices) == 0 {
		return []orm.Device{}
	}
	getInstance()
	var devices []orm.Device
	var ids []int
	for _, v := range groupDevices {
		ids = append(ids, v.DeviceID)
	}
	db.Find(&devices, ids)
	return devices
}

func GetDevicesByPlanId(planID int) (res []types.DeviceDTO) {
	getInstance()
	rows, err := db.Raw("SELECT d.id, d.`name`, d.mac, info.resolution, plan.`name` pname, d.state FROM devices d LEFT JOIN device_info info ON d.id = info.id, plan_device pd, plan WHERE d.id = pd.device_id AND pd.plan_id = ? AND plan.id = d.plan_id", planID).Rows()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var d types.DeviceDTO
		err := rows.Scan(&d.ID, &d.Name, &d.Mac, &d.Resolution, &d.PlanName, &d.State)
		if err != nil {
			log.Println(err)
			return
		}
		res = append(res, d)
	}
	return
}

func GetAllResolutions() (res []string) {
	getInstance()
	rows, err := db.Raw("SELECT resolution FROM `show` GROUP BY resolution").Rows()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var str string
		err := rows.Scan(&str)
		if err != nil {
			log.Println(err)
			return
		}
		res = append(res, str)
	}
	return
}

func GetAllDevice(organizationId string) []orm.Device {
	getInstance()
	var devices []orm.Device
	db.Find(&devices, "organization_id = ?", organizationId)
	return devices
}

func GetDeviceByState(organizationId string, state string) []orm.Device {
	getInstance()
	var devices []orm.Device
	db.Find(&devices, "organization_id = ? AND state = ?", organizationId, state)
	return devices
}

func GetPlanFirstImagesByIds(ids []interface{}) (res map[int]string) {
	res = make(map[int]string)
	if len(ids) == 0 {
		return
	}
	getInstance()
	rows, err := db.Raw("SELECT p.id, s.images FROM `play_period` d, `plan` p, `play_period_show` ds, `show` s WHERE d.plan_id = p.id AND d.id = ds.play_period_id AND ds.show_id = s.id AND p.id IN (?"+strings.Repeat(",?", len(ids)-1)+")", ids...).Rows()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var id int
		var images string
		err := rows.Scan(&id, &images)
		if err != nil {
			return
		}
		res[id] = images
	}
	return
}
