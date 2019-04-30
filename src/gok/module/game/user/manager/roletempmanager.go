package manager

import (
	"aliens/common/util"
	"gok/module/game/cache"
	"gok/module/game/conf"
	"gok/module/game/db"
	"gok/module/game/event"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"gopkg.in/mgo.v2/bson"
	"time"
)


//用户临时数据存储
type RoleTempManager struct {
	uid           int32

	searchResults map[string]*protocol.SearchResult

	offlineMessage []*protocol.NewsFeed

	switchGlobalMessagePush bool //是否开启全局销售信息推送

	avatar *string

	lastTask *protocol.Task

	lastUpgradeBelieverID string

	//ip string //登录ip
	//allBuilding []*protocol.BuildingState
	//allBuilding map[int32]int32 // type - level
	version string //当前连接的客户端版本号

	activeGroupID int32

	loginTime time.Time

	logoutTime time.Time

	//roleDialData *RandomDialID
	roleEvent *event.EventSession


	guideTime time.Time

	//dialIDWeightMapping		map[int32]int32	// id -- 权重

	//eventIDWeightMapping	map[int32]int32 // id -- 权重

	//dialPastIDs	[]int32 //最近两次抽到的转盘id

	//lastEventID	int32 //

	//saleInfo map[]*protocol.Sale
}

//type RandomDialID struct {
//	data []int32
//	times int32 //多少次必现星盘任务
//}

//func (this *RoleTempManager)GetDialData() *RandomDialID{
//	if this.roleDialData == nil {
//		return nil
//	}
//	return this.roleDialData
//}
//
//func (this *RoleTempManager)RandomDialDataTimes() {
//	this.roleDialData.times = rand.Int31n(5) + 5
//}
//
//func (this *RoleTempManager)AddDialData(dialID int32) (int32, bool) {
//	dialID := conf.RandomDial()
//	dialData := this.roleDialData.data
//	dialData = append(dialData, dialID)
//	if len(this.roleDialData.data) > 10 {
//		index := len(dialData) - 10
//		dialData = append(dialData[:0], dialData[index:]...)
//	}
//	return dialID, true
//}
//func (this *RoleTempManager)isInDialID(dialID int32) bool {
//	//for _, v := range this.roleDialData.data {
//	//	if v == dialID {
//	//		return true
//	//	}
//	//}
//	for index := 10 - this.roleDialData.times; index < 10 ; index++ {
//		if this.roleDialData.data[index] == dialID {
//			return true
//		}
//	}
//	return false
//}

//初始化
func (this *RoleTempManager) Init(role *db.DBRole) {
	this.uid = role.UserID
	this.switchGlobalMessagePush = true
	this.searchResults = make(map[string]*protocol.SearchResult)
	this.offlineMessage = []*protocol.NewsFeed{}

	//this.dialIDWeightMapping = make(map[int32]int32)
	//this.eventIDWeightMapping = make(map[int32]int32)
	//this.dialPastIDs = make([]int32,2)
	//this.initDialEvent()
	//this.lastEventID = 0
	//this.allBuilding = make(map[int32]int32)
	//零时加的星盘任务多少次必现
	//this.roleDialData = &RandomDialID{data:[]int32{}}

}

//更新数据库内存
func (this *RoleTempManager) Update(role *db.DBRole) {

}

//
func (this *RoleTempManager) SetActiveGroupID(activeGroupID int32) {
	this.activeGroupID = activeGroupID
}

//
func (this *RoleTempManager) GetActiveGroupID(activeGroupID int32) bool{
	if this.activeGroupID == activeGroupID {
		return true
	}
	return false
}

//设置客户端版本
func (this *RoleTempManager) SetClientVersion(version string) {
	this.version = version
}

//获取客户端版本
func (this *RoleTempManager) GetClientVersion() string {
	return this.version
}

func (this *RoleTempManager) UpdateLastTask(task *protocol.Task) {
	this.lastTask = task
}

func (this *RoleTempManager) UpdateLastUpgradeBelieverID(believerID string) {
	this.lastUpgradeBelieverID = believerID
}

func (this *RoleTempManager) GetLastTask() *protocol.Task {
	return this.lastTask
}

func (this *RoleTempManager) GetLastUpgradeBelieverID() string {
	return this.lastUpgradeBelieverID
}

func (this *RoleTempManager) SetGuideTime(guideTime time.Time) {
	this.guideTime = guideTime
}

func (this *RoleTempManager) GetGuideTime() time.Time {
	return this.guideTime
}


func (this *RoleTempManager)GetLogoutTime() time.Time{
	return this.logoutTime
}

func (this *RoleTempManager)SetLogoutTime(time time.Time) {
	this.logoutTime = time
}

func (this *RoleTempManager)GetLoginTime() time.Time{
	return this.loginTime
}

func (this *RoleTempManager)SetLoginTime(time time.Time) {
	this.loginTime = time
}


func (this *RoleTempManager) CleanOfflineMessage() {
	this.offlineMessage = []*protocol.NewsFeed{}
}

func (this *RoleTempManager) AddOfflineMessage(message *protocol.NewsFeed) {
	this.offlineMessage = append(this.offlineMessage, message)
}

func (this *RoleTempManager) GetOfflineMessage() []*protocol.NewsFeed {
	return this.offlineMessage
}


func (this *RoleTempManager) GetAvatar() string {
	if this.avatar == nil {
		cacheResult := cache.UserCache.GetUserAvatar(this.uid)
		this.avatar = &cacheResult
	}
	return *this.avatar
}

func (this *RoleTempManager) SwitchGlobalSalePush(open bool) {
	this.switchGlobalMessagePush = open
}

func (this *RoleTempManager) IsGlobalMessagePush() bool {
	return this.switchGlobalMessagePush
}

func (this *RoleTempManager) UpdateSearch(results []*protocol.SearchResult) {
	this.searchResults = make(map[string]*protocol.SearchResult)
	for _, result := range results {
		searchID := bson.NewObjectId().Hex()
		result.SearchID = searchID
		this.searchResults[searchID] = result
	}
}

func (this *RoleTempManager) GetSearch(searchID string) *protocol.SearchResult {
	return this.searchResults[searchID]
}


func (this *RoleTempManager) RemoveSearch(searchID string) {
	delete(this.searchResults, searchID)
}


/*-------------------------------事件-------------------------------*/
func (this *RoleTempManager) GetEvent() *event.EventSession {
	return this.roleEvent
}

func (this *RoleTempManager) DeleteEvent() {
	if this.roleEvent != nil {
		this.roleEvent = nil
	}
}

//生成事件
//参数 事件类型  事件发起者
func (this *RoleTempManager) GenEvent(eventType int32, uid int32, nickname string, guide bool, handler event.UserHandler) *event.EventSession {
	eventSession := event.NewEventSession(uid, eventType, uid, nickname, guide)
	this.roleEvent = eventSession
	this.roleEvent.SetHandler(handler)
	return eventSession
}

func (this *RoleTempManager) SetHandler(handler event.UserHandler) {
	if this.roleEvent != nil {
		this.roleEvent.SetHandler(handler)
	}
}

func (this *RoleTempManager) EnsureEvent() *event.EventSession {
	roleEvent := this.GetEvent()
	if roleEvent == nil {
		exception.GameException(exception.EVENT_NOT_FOUND)
	}
	return roleEvent
}

//func (this *RoleTempManager) DealEvent() {
//	time := time.Now()
//	//for _, roleEvent := range this.GetEvents() {
//	//	roleEvent.DealTimeChange(time)
//	//}
//	this.roleEvent.DealTimeChange(time)
//}
//------------------------------------转盘相关-----------------------------------
//func (this *RoleTempManager) initDialEvent() {
//	//this.resetDialAllWeight()
//	this.resetEventAllWeight()
//}

//func (this *RoleTempManager) resetDialAllWeight() {
//	confWeight := conf.DATA.DIAL_ID_WEIGHT_MAPPING
//	if confWeight == nil {
//		exception.GameException(exception.DIAL_NOT_FOUND)
//	}
//	for id, weight := range confWeight {
//		this.dialIDWeightMapping[id] = weight
//	}
//}

//func (this *RoleTempManager) resetDialWeight(resetID int32) {
//	confWeight := conf.DATA.DIAL_ID_WEIGHT_MAPPING
//	if confWeight == nil {
//		exception.GameException(exception.DIAL_NOT_FOUND)
//	}
//	this.dialIDWeightMapping[resetID] = conf.DATA.DIAL_ID_WEIGHT_MAPPING[resetID]
//}


//func (this *RoleTempManager) isResetDialWeight() (int32, bool) {
//	if this.dialPastIDs[0] == 0 {
//		return 0, false
//	}
//	if this.dialIDWeightMapping[this.dialPastIDs[0]] == 0 {
//		return this.dialPastIDs[0], true
//	}
//	return 0, false
//}
//
//func (this *RoleTempManager) pushDialID(dialID int32) {
//	this.dialPastIDs[0] = this.dialPastIDs[1]
//	this.dialPastIDs[1] = dialID
//}
//-----------------------------星盘相关-----------------------------
//func (this *RoleTempManager) resetEventAllWeight() {
//	confWeight := conf.DATA.EVENT_ID_WEIGHT_MAPPING
//	if confWeight == nil {
//		exception.GameException(exception.EVENT_FILTER_NOT_FOUND)
//	}
//	for id, weight := range confWeight {
//		this.eventIDWeightMapping[id] = weight
//	}
//}

func (this *RoleTempManager) RandomEvent(filterEvents []int32) int32{

	randoms := util.CopyMap(conf.DATA.EVENT_ID_WEIGHT_MAPPING)
	for id :=range randoms {
		for _, filterID := range filterEvents {
			if id == filterID {
				delete(randoms, id)
			}
		}
	}
	if len(randoms) == 0 {
		exception.GameException(exception.CAN_NOT_REVENGE)
	}

	eventID := util.RandomWeight(randoms)

	//if this.lastEventID != 0 {
	//	if this.eventIDWeightMapping[this.lastEventID] == 0 {
	//		this.resetEventWeight(this.lastEventID)
	//	}
	//}
	//this.lastEventID = eventID
	//this.eventIDWeightMapping[eventID] = 0

	if conf.DATA.EVENT_FILTER_DATA[eventID] == nil {
		exception.GameException(exception.EVENT_FILTER_NOT_FOUND)
	}
	return conf.DATA.EVENT_FILTER_DATA[eventID].EventBase
}

//func (this *RoleTempManager) resetEventWeight(lastID int32) {
//	weight := conf.DATA.EVENT_ID_WEIGHT_MAPPING
//	if weight == nil {
//		exception.GameException(exception.EVENT_FILTER_NOT_FOUND)
//	}
//	this.eventIDWeightMapping[lastID] = weight[lastID]
//}


