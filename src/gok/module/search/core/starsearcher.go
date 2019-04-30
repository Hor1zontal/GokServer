package core

import (
	"math/rand"
	"aliens/common/set"
	"time"
	"gok/service/msg/protocol"
	"gok/constant"
	"gok/module/search/cache"
	"aliens/common/util"
	basecache "gok/cache"
	"gok/service/rpc"
	"gok/module/search/db"
	"gok/service/lpc"
	"gok/module/search/conf"
	"aliens/log"
)

var StarSearcher = &Searcher{}

func init() {
	StarSearcher.starMapping = make(map[int32]*StarIndex)
	StarSearcher.believerEvent = newEventCategory()
	StarSearcher.buildingEvent = newEventCategory()
	//StarSearcher.InitEventFilter(constant.EVENT_ID_LOOT_BELIEVER)
	//StarSearcher.InitEventFilter(constant.EVENT_ID_ATK_BUILDING)
	//StarSearcher.InitEventFilter(constant.EVENT_ID_LOOT_FAITH)
	//StarSearcher.eventFilters[] = &EventFilter{eventID,filters:make(map[int32]*StarCategory)}
	//StarSearcher.believerStarLevelMapping = make(map[int32]int32)
	//StarSearcher.buildingStarLevelMapping = make(map[int32]int32)
}

//var VALUE = struct{}{}

const (
	//cleanTimeout = 5 * 24 * time.Hour
	//randomCD = time.Hour
)

func(this *Searcher) Init() {
	//事件类型 - 等级范围 - 配置
	for eventType, base := range conf.Base.RandomTarget {
		filterType := constant.GetEventMatchRule(eventType)
		for _,  levelBase := range base {
			if filterType == constant.MATCH_RULE_BELIEVER {
				StarSearcher.believerEvent.addLevelMap(levelBase)

			} else if filterType == constant.MATCH_RULE_BUILDING {
				StarSearcher.buildingEvent.addLevelMap(levelBase)
			}
		}
	}

	now := time.Now()
	var starBuildingInfoArray []*db.StarBuildingInfo
	db.DatabaseHandler.QueryAll(&db.StarBuildingInfo{}, &starBuildingInfoArray)
	for _, starBuildingInfo := range starBuildingInfoArray {
		if starBuildingInfo.UpdateTime.Unix() < 0 || now.Sub(starBuildingInfo.UpdateTime).Seconds() < conf.Base.SearchActiveTime {
			StarSearcher.UpdateStarBuildingLevel(starBuildingInfo.ID, starBuildingInfo.TotalLevel, false)
		}
	}

	var starBelieverInfoArray []*db.StarBelieverInfo
	db.DatabaseHandler.QueryAll(&db.StarBelieverInfo{}, &starBelieverInfoArray)
	for _, starBelieverInfo := range starBelieverInfoArray {
		if starBelieverInfo.UpdateTime.Unix() < 0 || now.Sub(starBelieverInfo.UpdateTime).Seconds() < conf.Base.SearchActiveTime {
			StarSearcher.UpdateStarBelieverLevel(starBelieverInfo.ID, starBelieverInfo.TotalLevel, false)
		}
	}
}


//索引根节点
type Searcher struct {
	//拥有圣物的星球索引
	starMapping map[int32]*StarIndex //星球类型 - 搜索句柄

	believerEvent *EventCategory

	buildingEvent *EventCategory

}

func (this *Searcher) UpdateRandomStar(starID int32, opts map[int32]int32)  {
	for opt, _ := range opts {
		if opt == constant.SEARCH_OPT_UPDATE_BUILDING {
			this.buildingEvent.addRandomStar(starID)
		} else if opt == constant.SEARCH_OPT_UPDATE_BELIEVER {
			this.believerEvent.addRandomStar(starID)
		}
	}
}

//func (this *Searcher) InitEventFilter(eventID int32) {
//	this.eventFilters[eventID] = &EventFilter{eventID:eventID, filters:make(map[int32]time.Time)}
//}

func (this *Searcher) OPT(starID int32, opt int32, updateData int32, sync bool) {
	switch opt {
		case constant.SEARCH_OPT_UPDATE_BUILDING:
			this.UpdateStarBuildingLevel(starID, updateData, !sync)
			break
		case constant.SEARCH_OPT_UPDATE_BELIEVER:
			this.UpdateStarBelieverLevel(starID, updateData, !sync)
			break
		case constant.SEARCH_OPT_REMOVE_STAR:
			this.RemoveStarIndex(starID)
			break
		default:
	}
}


func (this *Searcher) BuildTarget(starID int32) *protocol.Target {
	if starID == 0 {
		return nil
	}
	if starID < 0 {
		return rpc.StarServiceProxy.GetEventRobotTarget(starID).GetTarget()
	}

	starProps := cache.StarCache.GetStarProps(starID)
	if starProps == nil || len(starProps) == 0 {
		return nil
	}
	uid := int32(starProps[basecache.STAR_PROP_OWNER])
	activeStarID := cache.StarCache.GetUserActiveStar(uid)

	//索引数据不是当前完的星球，需要过滤
	if starID != activeStarID {
		return nil
	}

	believerTotalLevel := int32(starProps[basecache.STAR_PROP_BELIEVER_TOTAL_LEVEL])
	buildingTotalLevel := int32(starProps[basecache.STAR_PROP_BUILDING_TOTAL_LEVEL])

	if believerTotalLevel == 0 {
		category := this.believerEvent.levelMapping[1]
		believerTotalLevel = util.RandInt32Scop(category.min, category.max)
	}
	if buildingTotalLevel == 0 {
		category := this.buildingEvent.levelMapping[1]
		buildingTotalLevel = util.RandInt32Scop(category.min, category.max)
	}
	return &protocol.Target{
		Id: uid,
		StarType: int32(starProps[basecache.STAR_PROP_TYPE]),
		BelieverTotalLevel: believerTotalLevel,
		BuildingTotalLevel: buildingTotalLevel,
	}
}

func (this *Searcher) RandomEventStar(eventType int32, uid int32, mutualID int32, robotFilter bool, num int, always bool) []*protocol.Target {
	filterType := constant.GetEventMatchRule(eventType)
	filters := []int32{} //随机过滤列表
	results := []*protocol.Target{} //随机结果

	//过滤自己
	myStarId := cache.StarCache.GetUserActiveStar(uid)
	if myStarId != 0 {
		filters = append(filters, myStarId)
	}

	testUID := conf.Base.RandomTest[eventType]
	if testUID != 0 {
		starID := cache.StarCache.GetUserActiveStar(testUID)
		if starID != 0 {
			target := this.BuildTarget(starID)
			if target != nil {
				results = append(results, target)
				filters = append(filters, starID)
			}
		}
	}

	//优先随机复仇对象
	if mutualID != 0 {
		mutualStarID := cache.StarCache.GetUserActiveStar(mutualID)
		target := this.BuildTarget(mutualStarID)

		if target != nil {
			results = append(results, target)
			filters = append(filters, mutualStarID)
		}
	}

	remainNum := num - len(results)
	if remainNum <= 0 {
		return results
	}

	robotLevels := []int32{}
	//剩余的数量从陌生人中随机
	for level := 1; level <= remainNum; level++ {
		var starID int32 = 0
		if filterType == constant.MATCH_RULE_BELIEVER {
			starID = this.believerEvent.RandomStar(filters, always)
		} else {
			starID = this.buildingEvent.RandomStar(filters, always)
		}
		//if matchStar != nil && cache.StarCache.ExistMutual(matchStar.Owner, uid) {
		//	matchStar = nil
		//}
		target := this.BuildTarget(starID)

		//log.Debug("random target %v-%v", starID, target)
		if target != nil {
			results = append(results, target)
			filters = append(filters, starID)
		} else {
			if starID != 0 {
				this.RemoveStarIndex(starID)
			}
			robotLevels = append(robotLevels, int32(level))
		}
	}

	if len(robotLevels) == 0 && robotFilter {
		return results
	}

	//剩余的数量从机器人总随机
	eventRobots := rpc.StarServiceProxy.RandomEventRobot(eventType, robotLevels).GetTargets()
	if eventRobots != nil && len(eventRobots) > 0 {
		for _, eventRobot := range eventRobots {
			results = append(results, eventRobot)
		}
	}
	return results
}


func (this *Searcher) SearchUsers(filterUID int32, starType int32, items []int32) []*protocol.SearchResult {
	starIndex := this.starMapping[starType]
	if starIndex == nil {
		return nil
	}
	return starIndex.SearchUser(filterUID, items)
}

func (this *Searcher) RemoveStarIndex(starID int32) {
	this.buildingEvent.removeStar(starID)
	lpc.DBServiceProxy.Delete(&db.StarBuildingInfo{ID:starID}, db.DatabaseHandler)

	this.believerEvent.removeStar(starID)
	lpc.DBServiceProxy.Delete(&db.StarBelieverInfo{ID:starID}, db.DatabaseHandler)
}

func (this *Searcher) DealFilterExpire() {
	now := time.Now()
	this.buildingEvent.DealFilterExpire(now)
	this.believerEvent.DealFilterExpire(now)
}


func (this *Searcher) Clean() {
	now := time.Now()
	log.Debug("start clean search data... %v", now)

	this.buildingEvent.CleanExpire(now)
	this.believerEvent.CleanExpire(now)

	//for _, starCategory := range this.starBuildingMapping {
	//	for starID, updateTime := range starCategory.mapping {
	//		if now.Sub(updateTime) > cleanTimeout {
	//			delete(starCategory.mapping, starID)
	//		}
	//	}
	//}
	//
	//for _, starCategory := range this.starBelieverMapping {
	//	for starID, updateTime := range starCategory.mapping {
	//		if now.Sub(updateTime) > cleanTimeout {
	//			delete(starCategory.mapping, starID)
	//		}
	//	}
	//}
	log.Debug("end clean search data..., duration %v(s)", time.Now().Sub(now))
}

func (this *Searcher) UpdateStarBelieverLevel(starID int32, believerLevel int32, update bool) {
	result := this.believerEvent.UpdateData(starID, believerLevel)
	if result && update {
		lpc.DBServiceProxy.ForceUpdate(&db.StarBelieverInfo{starID, believerLevel, time.Now()}, db.DatabaseHandler)
	}
}


func (this *Searcher) UpdateStarBuildingLevel(starID int32, buildingLevel int32, update bool) {
	result := this.buildingEvent.UpdateData(starID, buildingLevel)
	if result && update {
		lpc.DBServiceProxy.ForceUpdate(&db.StarBuildingInfo{starID, buildingLevel, time.Now()}, db.DatabaseHandler)
	}
}

//func (this *Searcher) Opt(filterUID int32, starType int32, items []int32) []*protocol.SearchResult {
//	starIndex := this.starMapping[starType]
//	if (starIndex == nil) {
//		return nil
//	}
//	return starIndex.SearchUser(filterUID, items)
//}


//星球索引
type StarIndex struct {
	mapping map[int32]*ItemIndex
}

//物品索引
type ItemIndex struct {
	mapping    *set.HashSet
	randoms    []interface{}
	updateTime time.Time
}


func (this *StarIndex) SearchUser(filterUID int32, items []int32) []*protocol.SearchResult {
	result := []*protocol.SearchResult{}

	for _, item := range items {
		if (item == 0) {
			continue
		}
		itemIndex := this.mapping[item]
		if (itemIndex == nil) {
			continue
		}
		uid := itemIndex.RandomUser(filterUID)
		if (uid != 0) {
			if (containsResult(uid, result)) {
				continue
			}
			result = append(result,
				&protocol.SearchResult{
					Id:       uid,
					ItemID:   item})
		}
	}
	return result
}

func containsResult(uid int32 , results []*protocol.SearchResult) bool {
	for _, result := range results {
		if (result.GetId() == uid) {
			return true
		}
	}
	return false
}

func (this *ItemIndex) RandomUser(filterUID int32) int32 {
	now := time.Now()
	if (now.Sub(this.updateTime).Minutes() > 2) {
		this.randoms = this.mapping.Elements()
	}
	total := len(this.randoms)
	if (total == 0) {
		return 0
	}
	randomIndex := rand.Intn(total)
	result := this.randoms[randomIndex].(int32)
	if (result != filterUID) {
		return result
	}
	if (randomIndex < total - 1) {
		return this.randoms[randomIndex + 1].(int32)
	} else if (randomIndex > 0) {
		return this.randoms[randomIndex - 1].(int32)
	}
	return 0
}

//
func (this *Searcher) AddItemIndex(starType int32, itemID int32, uid int32) {
	if itemID == 0 {
		return
	}
	starIndex := this.starMapping[starType]
	if (starIndex == nil) {
		starIndex = &StarIndex{mapping: make(map[int32]*ItemIndex)}
		this.starMapping[starType] = starIndex
	}
	starIndex.AddItemIndex(itemID, uid)
}

func (this *StarIndex) AddItemIndex(itemID int32, uid int32) {
	itemIndex := this.mapping[itemID]
	if (itemIndex == nil) {
		itemIndex = &ItemIndex{mapping: set.NewHashSet()}
		this.mapping[itemID] = itemIndex
	}
	itemIndex.AddItemIndex(uid)
}

func (this *ItemIndex) AddItemIndex(uid int32) {
	this.mapping.Add(uid)
}

func (this *Searcher) RemoveItemsIndex(starType int32, itemID int32, uid int32) {
	starIndex := this.starMapping[starType]
	if (starIndex == nil) {
		return
	}
	starIndex.RemoveItemIndex(itemID, uid)
}

func (this *Searcher) RemoveItemIndex(starType int32, itemID int32, uid int32) {
	if itemID == 0 {
		return
	}
	starIndex := this.starMapping[starType]
	if (starIndex == nil) {
		return
	}
	starIndex.RemoveItemIndex(itemID, uid)
}

func (this *StarIndex) RemoveItemIndex(itemID int32, uid int32) {
	itemIndex := this.mapping[itemID]
	if (itemIndex == nil) {
		return
	}
	itemIndex.RemoveItemIndex(uid)
}

func (this *ItemIndex) RemoveItemIndex(uid int32) {
	this.mapping.Remove(uid)
}
