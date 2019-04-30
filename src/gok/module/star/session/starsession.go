/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/5/26
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package session

import (
	"aliens/common/character"
	timeutil "aliens/common/util"
	"aliens/log"
	"gok/constant"
	"gok/module/cluster/center"
	"gok/module/star/cache"
	"gok/module/star/conf"
	"gok/module/star/db"
	"gok/module/star/util"
	"gok/service/exception"
	"gok/service/lpc"
	"gok/service/msg/protocol"
	"gok/service/rpc"
	"math/rand"
	"time"
)

//星球会话
type StarSession struct {
	//sync.RWMutex    //操作锁
	*db.DBStar      //星球数据
	powerBuff          int32 //法力值buff缓存值,方便验证脏数据是否通知
	dirty              bool
	buildingTotalLevel int32             //建筑总等级
	//buildingExMaxLevel int32           //建筑曾经达到的最大等级
	believerCount      int32             //信徒总数量
	believerTotalLevel int32             //信徒总等级
	Buff               map[int32]float32 //buff / type - buffNum
	flags			   map[int32]*db.DBStarFlag
}

const DEFAULT_FLAG_VALUE int32 = 0

func (this *StarSession) isDirty() bool {
	return this.dirty
}

func (this *StarSession) cleanDirty() {
	this.dirty = false
}

func (this *StarSession) setDirty() {
	this.dirty = true
	cache.StarCache.SetUserIDNode(this.Owner, center.GetServerNode())
	//lpc.DBServiceProxy.Opt(this.DBStar, db.DatabaseHandler)
	//db.UpdateHandler.UpdateQueue(database.OP_UPDATE, this.DBStar)
	//db.DatabaseHandler.UpdateOne(this.DBStar)
}

func newStarSession(dbStar *db.DBStar) *StarSession {
	result := &StarSession{DBStar: dbStar}
	return result
}

//初始化
func (this *StarSession) Init() {
	//this.Lock()
	//defer this.Unlock()
	this.initFlag()
	this.updateBuildingExMax()
	//解锁星球flag
	this.checkLockFlags()
	this.updateBelieverCount(false)
	this.updateBuildingLevel(false)
	this.updateBuff()
	this.correctCivilizationValue()
	//for _, building := range this.Building {
	//	for _, groove := range building.ItemGrooves {
	//		if groove.ItemID != 0 {
	//			searcher.Index.AddItemIndex(this.Type, groove.ItemID, this.Owner)
	//		}
	//	}
	//}
	//this.updateBuff()
	//this.initBuffListener()
}

//处理星球的建筑过期
func (this *StarSession) DealStarTimer(now time.Time) {
	//超过释放时间需要释放
	//this.Lock()
	//defer this.Unlock()
	changeBuildings := []*db.DBBuilding{}
	//itemIDs := []int32{}
	for _, building := range this.Building {
		//机器人帮助维修时间到
		//if building.IsRobotHelpOverdue() {
		//	robotID := character.RandInt32Scop(-78, -1)
		//	ret, starType, buildingLevel := this.HelpRepairBuild(building.Type, robotID)
		//	if ret {
		//		message := &protocol.C2GS{
		//			Sequence: []int32{523},
		//			AddNewsFeed: &protocol.NewsFeed{
		//				Id: 		bson.NewObjectId().Hex(),
		//				RelateID:   robotID,
		//				Type: 		constant.NEWSFEED_TYPE_BE_HELP_REPAIR,
		//				Time:		time.Now().Unix(),
		//				Param1: 	starType,
		//				Param2:		building.Type,
		//				Param3:		buildingLevel,
		//				Ext:		nil,
		//			},
		//		}
		//		rpc.UserServiceProxy.PersistCall(this.Owner, message)
		//	}
		//}

		//没有推送过消息、需要更新推送消息 推送标识在用户上线清除
		if !this.Push {
			if building.IsRepairing() {
				//修理好了需要推送公众号
				if building.IsRepaired(now) {
					rpc.PassportServiceProxy.WechatEventPush(this.Owner, constant.EVENT_BUILD_REPAIR, 0)
					this.Push = true
				}
			} else if building.IsBroken() {
				if building.IsBrokenRemainOverdue(int64(time.Hour.Seconds())) {
					rpc.PassportServiceProxy.WechatEventPush(this.Owner, constant.EVENT_BUILD_DAMAGE, 0)
					this.Push = true
				}
			} else if building.IsUpgrading() {
				if building.IsUpgraded(now) {
					rpc.PassportServiceProxy.WechatEventPush(this.Owner, constant.EVENT_BUILD_COMPLETION, 0)
					this.Push = true
				}
			}
		}

		//损坏过期
		if building.IsBrokenOverdue() {
			//itemIDs = append(itemIDs, building.GetItems()...)
			building.Reset()
			//if cache.StarCache.ExistStarHelpRepair(this.Type, building.Type) {
			//}
			cache.StarCache.DelStarHelpRepair(this.ID, building.Type)
			this.updateBuildingLevel(true)
			changeBuildings = append(changeBuildings, building)
		}
	}

	//推送建筑损坏归零
	if len(changeBuildings) != 0 {
		this.setDirty()
		//通知用户模块更新法力值上限
		message := &protocol.C2GS{
			Sequence: []int32{520},
			BuildingReset: &protocol.BuildingReset{
				PowerLimit: this.GetPowerLimit(),
				Building:   util.BuildStarBuildingInfos(changeBuildings),
				//ItemID:     itemIDs,
			},
		}
		rpc.UserServiceProxy.PersistCall(this.Owner, message)
	}



	//异步更新数据库
	if this.isDirty() {
		lpc.DBServiceProxy.Update(this.DBStar, db.DatabaseHandler)
		this.cleanDirty()
	}
}

func (this *StarSession) DealUnSave() {
	//this.Lock()
	//defer this.Unlock()
	if this.isDirty() {
		db.DatabaseHandler.UpdateOne(this.DBStar)
		//lpc.DBServiceProxy.Opt(this.DBStar, db.DatabaseHandler)
		this.cleanDirty()
	}
}



//处理平民的自动添加
//func (this *StarSession) DealBelieverBuff() {
//	//星球未占领和用户不在线不增加buff
//	if (this.Owner == 0 || !cache.StarCache.IsUserOnline(this.Owner)) {
//		return
//	}
//	results := conf.GetStarInitBeliever(this.Type)
//	if (results == nil) {
//		return
//	}
//	if (!this.canAddBeliever(1)) {
//		return
//	}
//	result := this.getMinBeliever(results)
//	if (result == "") {
//		return
//	}
//
//	believer := this.addBeliever(result, 1)
//	this.setDirty()
//	PushBelieverInfo(this.Owner, believer)
//}

func (this *StarSession) IsRobot() bool {
	return this.ID < 0
}

func (star *StarSession) BuildTarget() *protocol.Target {
	return  &protocol.Target{
		Id: star.Owner,
		StarType: star.Type,
		BelieverTotalLevel: star.believerTotalLevel,
		BuildingTotalLevel: star.buildingTotalLevel,
	}
}

//func (this *StarSession) ChangeOwner(uid int32) {
//	if this.Owner != 0 { //已经拥有占领者了
//		exception.GameException(exception.STAR_OCCUPY_FAILED)
//	}
//	//设置星球被占领了
//
//	this.Owner = uid
//}

//func (this *StarSession) ActiveGrooveEffect(buildingType int32, grooveID int32) (int64, bool) {
//	building := this.getBuilding(buildingType)
//	if building == nil {
//		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
//	}
//	building.CheckoutFaithBuff(this.Type, false)
//	timestamp, effect := building.ActiveGrooveEffect(grooveID)
//	this.setDirty()
//	return timestamp, effect
//}

func (this *StarSession) IsDoneStar() bool {
	if !this.IsBuildingAllMaxLevel() {
		exception.GameException(exception.STAR_OCCUPY_FAILED)
	}
	return true
}

func (this *StarSession) LoginStarInfo() *protocol.LoginStarInfoRet {

	result := &protocol.LoginStarInfoRet{CurrStar: &protocol.StarInfo{StarID: this.ID, StarType: this.Type}}

	upgradedBuilding := []*protocol.BuildingState{}
	repairedBuilding := []*protocol.BuildingState{}
	allBuilding := []*protocol.BuildingState{}
	for _, building := range this.Building {
		allBuilding = append(allBuilding, &protocol.BuildingState{
			StarType:this.Type,
			BuildingType:building.Type,
			BuildingLevel:building.Level,
		})
		upgradeDone, _, _ := this.dealUpgradeBuilding(building, false)
		if upgradeDone {
			upgradedBuilding = append(upgradedBuilding, &protocol.BuildingState{
				StarType:      this.Type,
				BuildingType:  building.Type,
				BuildingLevel: building.Level,
			})
		}
		repairedDone, _ := this.dealUpdateRepairing(building, false)
		if repairedDone {
			repairedBuilding = append(repairedBuilding, &protocol.BuildingState{
				StarType:      this.Type,
				BuildingType:  building.Type,
				BuildingLevel: building.Level,
			})
		}
	}
	result.PowerLimit = this.GetPowerLimit()
	result.RepairedBuilding = repairedBuilding
	result.UpgradedBuilding = upgradedBuilding
	result.AllBuilding = allBuilding
	return result
}

//func (this *StarSession) AccBuildingGrooveEffect(buildingType int32, grooveID int32, believerIDs []string) (int64, bool) {
//	num := len(believerIDs)
//	for _, believerID := range believerIDs {
//		level := conf.GetGameObjectLevel(believerID)
//		if level != constant.MAX_BELIEVER_LEVEL {
//			exception.GameException(exception.STAR_BELIEVER_NOT_ENOUGH)
//		}
//	}
//
//	this.ensureBelievers(believerIDs)
//
//	building := this.getBuilding(buildingType)
//	if building == nil {
//		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
//	}
//	building.CheckoutFaithBuff(this.Type, false)
//	timestamp, effect := building.AccGrooveEffect(grooveID, num)
//
//	this.decBelievers(believerIDs)
//	this.setDirty()
//	return timestamp, effect
//}

//

//获取星球统计信息
func (this *StarSession) GetStatisticsInfo() []*protocol.Statistics {
	result := []*protocol.Statistics{}
	for _, statistics := range this.Statistics {
		result = append(result, statistics.BuildProtocol())
	}
	return result
}

func (this *StarSession) GetStarHistory() []*protocol.History {
	result := []*protocol.History{}
	for _, history := range this.History {
		result = append(result, history.BuildProtocol())
	}
	return result
}

//计算法力值上限
func (this *StarSession) GetPowerLimit() int32 {
	var powerLimit int32 = 0
	for _, roleStarBuilding := range this.Building {
		powerLimit += conf.GetBuildingPowerLimit(this.Type, roleStarBuilding.Type, roleStarBuilding.Level)
	}
	if powerLimit > 0 {
		powerLimit += conf.Base.InitPowerLimit
	}
	powerLimit += int32(this.Buff[constant.BUFF_MANA_LIMIT])
	return powerLimit
}

//处理建筑信仰值  draw 是否领取
//func (this *StarSession) TakeoutBuildingItem(buildingType int32, grooveID int32) int32 {
//
//	building := this.getBuilding(buildingType)
//	if (building == nil) {
//		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
//	}
//	//去除物品需要计算信仰值
//	building.CheckoutFaithBuff(this.Type, false)
//	result := building.TakeoutItem(grooveID)
//
//	searcher.Index.RemoveItemIndex(this.Type, result, this.Owner)
//	this.setDirty()
//	return result
//}

func (this *StarSession) DecBelieverShield() bool {
	return decShield(this.BelieverShield)
}

func (this *StarSession) DecFaithShield() bool {
	return decShield(this.FaithShield)
}

func (this *StarSession) DecBuildingShield() bool {
	return decShield(this.BuildingShield)
}

func (this *StarSession) GetShieldInfo() *protocol.GetStarShieldRet {
	result := &protocol.GetStarShieldRet{}
	if this.FaithShield != nil {
		calShield(this.FaithShield)
		result.FaithShield = this.FaithShield.BuildProtocol()
	}
	if this.BelieverShield != nil {
		calShield(this.BelieverShield)
		result.BelieverShield = this.BelieverShield.BuildProtocol()
	}
	if this.BuildingShield != nil {
		calShield(this.BuildingShield)
		result.BuildingShield = this.BuildingShield.BuildProtocol()
	}
	return result
}


//抢夺圣物
func (this *StarSession) LootBuildingItem(buildingID int32, itemID int32) *db.DBBuilding {
	building := this.getBuilding(buildingID)
	if building == nil {
		return nil
	}
	if building.IsBroken() {
		return nil
	}
	prob := (26 + float64(building.Level)) / (100 + float64(this.getBuildingAllLevel()))
	prob = prob * (1 + float64(this.Buff[constant.BUFF_RELIC_PROTECT]))
	//概率成功
	if prob < rand.Float64() {
		return nil
	}
	building.SetBroken(this.Buff, false)
	this.setDirty()
	message := util.BuildBuildingInfoPush(this.Owner, this.ID, this.Type, building)
	rpc.UserServiceProxy.Push(this.Owner, message)

	//for _, groove := range building.ItemGrooves {
	//	if groove.ItemID == itemID {
	//		groove.ItemID = 0
	//		building.SetBroken()
	//		building.UpdateBuff()
	//		this.setDirty()
	//
	//		message := util.BuildBuildingInfoPush(this.Owner, this.ID, this.Type, building)
	//		rpc.UserServiceProxy.Push(this.Owner, cache.StarCache.GetUserNode(this.Owner), message)
	//		return building
	//	}
	//}

	//for _, groove := range building.ItemGrooves {
	//	if groove.ItemID == itemID {
	//		groove.ItemID = 0
	//		building.SetBroken()
	//		building.UpdateBuff()
	//		this.setDirty()
	//
	//		message := util.BuildBuildingInfoPush(this.Owner, this.ID, this.Type, building)
	//		rpc.UserServiceProxy.Push(this.Owner, cache.StarCache.GetUserNode(this.Owner), message)
	//		return building
	//	}
	//}

	return nil
}

//func (this *StarSession) ReplaceBuildingItem(buildingType int32, grooveIDs []int32, itemIDs []int32, takeoutItems []int32) []*protocol.ItemGroove {
//	building := this.getBuilding(buildingType)
//	if building == nil {
//		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
//	}
//
//	//建筑中取出的圣物
//	realTakeoutItems := []int32{}
//	var persistItems map[int32]*db.DBItemGroove = nil
//
//	//优先从别的建筑拔除圣物
//	for _, building := range this.Building {
//		building.CheckoutFaithBuff(this.Type, false)
//		if building.Type == buildingType {
//			persistItems = building.TakeoutExistItem(takeoutItems)
//			for _, persistItem := range persistItems {
//				realTakeoutItems = append(realTakeoutItems, persistItem.ItemID)
//			}
//		} else {
//			buildingTakeoutItems := building.TakeoutExistItemID(takeoutItems)
//			realTakeoutItems = character.AppendArray(realTakeoutItems, buildingTakeoutItems)
//		}
//	}
//
//	//取出的物品
//	realTakeoutItems = character.AppendArray(realTakeoutItems, building.TakeinItems(grooveIDs, itemIDs, persistItems))
//
//	////放入背包的圣物
//	//bagTakeinItem := character.GetArrayDeff(realTakeoutItems, itemIDs)
//	////从背包中需要扣除的圣物
//	//bagTakeoutItem := character.GetArrayDeff(itemIDs, realTakeoutItems)
//	for _, item := range takeoutItems {
//		searcher.Index.AddItemIndex(this.Type, item, this.Owner)
//	}
//	for _, item := range itemIDs {
//		searcher.Index.RemoveItemIndex(this.Type, item, this.Owner)
//	}
//
//
//	this.setDirty()
//	return building.BuildGroovesProtocol()
//}

//func (this *StarSession) ActiveBuildingGroup(buildingType int32) []int32 {
//
//	building := this.getBuilding(buildingType)
//	if building == nil {
//		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
//	}
//	returnItems := building.ActiveGrooves()
//	return returnItems
//}

//func (this *StarSession) ResetBuildingGrooves(buildingType int32, localGroove []int32) ([]*protocol.ItemGroove, []int32) {
//	this.Lock()
//	defer this.Unlock()
//	building := this.getBuilding(buildingType)
//	if (building == nil) {
//		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
//	}
//	//去除物品需要计算信仰值
//	building.CheckoutFaithBuff(this.Type, false)
//	returnItems := building.ResetGrooves(localGroove)
//	this.setDirty()
//	return building.BuildGroovesProtocol(), returnItems
//}

//func (this *StarSession) getBuffRatio(buffType int32) float32 {
//	var ratio float32 = 0
//	for _, building := range this.Building {
//		ratio += building.GetBuffRatio(buffType)
//	}
//	return ratio
//}

func (this *StarSession) GetBuff() map[int32]float32 {
	return this.Buff
}

func (this *StarSession) SetBuff(groupID int32) {
	buffID := conf.Base.ItemGroup[groupID].BuffID
	buff := conf.Base.Buff[buffID]
	if buff == nil {
		log.Error("invalid buff config groupID %v - buffID %v", groupID, buffID)
		return
	}

	switch buff.Type {
	case constant.BUFF_BELIEVER_PROTECT:
		if this.BelieverShield == nil {
			this.BelieverShield = &db.Shield{Value:int32(buff.Ratio),Limit:int32(buff.Ratio),UpdateTime:time.Now()}
		}
		break
	case constant.BUFF_FAITH_PROTECT:
		if this.FaithShield == nil {
			this.FaithShield = &db.Shield{Value:int32(buff.Ratio),Limit:int32(buff.Ratio),UpdateTime:time.Now()}
		}
		break
	case constant.BUFF_BUILDING_PROTECT:
		if this.BuildingShield == nil {
			this.BuildingShield = &db.Shield{Value:int32(buff.Ratio),Limit:int32(buff.Ratio),UpdateTime:time.Now()}
		}
		break
	case constant.BUFF_MANA:
		newPowerBuff := int32(buff.Ratio)
		if newPowerBuff != this.powerBuff {
			this.powerBuff = newPowerBuff
			//需要实时更新法力值间隔,供用户模块调用
			cache.StarCache.SetBuffMANAInterval(this.Owner, this.powerBuff)
		}
		break
	case constant.BUFF_RELIC_STEAL:
		cache.StarCache.SetBuffRelicSteal(this.Owner, float64(buff.Ratio))
		break
	}

	this.Buff[buff.Type] = buff.Ratio
}

func (this *StarSession) GetShield(shieldType int32) *protocol.Shield {
	switch shieldType {
		case constant.SHIELD_TYPE_BELIEVER:
			if this.BelieverShield != nil {
				calShield(this.BelieverShield)
				return this.BelieverShield.BuildProtocol()
			}
			break
		case constant.SHIELD_TYPE_BULDING:
			if this.BuildingShield != nil {
				calShield(this.BuildingShield)
				return this.BuildingShield.BuildProtocol()
			}
			break
		case constant.SHIELD_TYPE_FAITH:
			if this.FaithShield != nil {
				calShield(this.FaithShield)
				return this.FaithShield.BuildProtocol()
			}
			break
		default:
	}
	return nil
}

func decShield(shield *db.Shield) bool {
	if shield == nil {
		return false
	}
	calShield(shield)
	isShieldLimit := shield.Value == shield.Limit
	if shield.Value > 0 {
		shield.Value -= 1
		//当从满变成不满的时候，又开始计算
		if isShieldLimit {
			shield.UpdateTime = time.Now()
		}
		return true
	}
	return false
}

func calShield(shield *db.Shield) {
	//达到上限不需要刷新计算
	if shield.Value == shield.Limit {
		return
	}
	addCount, refreshTime := timeutil.RefreshTime(shield.UpdateTime, time.Now(), constant.SHIELD_ADD_TIME)
	shield.Value += addCount
	shield.UpdateTime = refreshTime
	if shield.Value > shield.Limit {
		shield.Value = shield.Limit
	}
}


//func (this *StarSession) initBuffListener() {
//	for _, building := range this.Building {
//		building.BuffChangeListener = this.buffChange
//	}
//}

//buff 变更监听处理
//func (this *StarSession) buffChange(buildingType int32, buff map[int32]int32) {
//	if this.Owner == 0 {
//		return
//	}
//	newPowerBuff := int32(this.getBuffRatio(constant.BUFF_MANA))
//	if newPowerBuff != this.powerBuff {
//		this.powerBuff = newPowerBuff
//		//需要实时更新法力值间隔,供用户模块调用
//		cache.StarCache.SetBuffMANAInterval(this.Owner, this.powerBuff)
//	}
//	if buildingType != 0 {
//		message := util.BuildBuildingBuffInfoPush(buildingType, buff)
//		rpc.UserServiceProxy.Push(this.Owner, cache.StarCache.GetUserNode(this.Owner), message)
//	}
//}

//func (this *StarSession) updateBuff() {
//	for _, building := range this.Building {
//		building.UpdateBuff()
//	}
//	this.buffChange(0, nil)
//}

func (this *StarSession) updateBuff() {
	if this.Buff == nil {
		this.Buff = make(map[int32]float32)
	}
	for _, itemGroup := range this.ItemGroups {
		if !itemGroup.Done {
			continue
		}
		groupID := itemGroup.ID
		this.SetBuff(groupID)
	}
	this.setDirty()
}


func (this *StarSession) updateEventTimesCache() {
	var mutual int32 = 0
	var beMutual int32 = 0

	for _, statistic := range this.Statistics {
		switch statistic.ID {
		case constant.STAR_STATISTIC_LOOT_BELIEVER, constant.STAR_STATISTIC_LOOT_FAITH, constant.STAR_STATISTIC_ATK_BUILDING:
			mutual += int32(statistic.Value)
		case constant.STAR_STATISTIC_BE_LOOT_BELIEVER, constant.STAR_STATISTIC_BE_LOOT_FAITH, constant.STAR_STATISTIC_BE_ATK_BUILDING:
			beMutual += int32(statistic.Value)
		}
	}
	cache.StarCache.SetMutualTimes(this.ID, mutual)
	cache.StarCache.SetBeMutualTimes(this.ID, beMutual)
}

//是否有修理中的建筑
//func (this *StarSession) hasRepairingBuilding() bool {
//	for _, v := range this.Building{
//		if v.Level > 0 && v.IsRepairing() {
//			return true
//		}
//	}
//	return false
//}

//获取所有标识
func (this *StarSession) getStatisticsArray() []*db.Statistics {
	result := []*db.Statistics{}
	for _, statistic := range this.Statistics {
		result = append(result, statistic)
	}
	return result
}

//获取角色标识
func (this *StarSession) getStatisticsValue(key int32) float64 {
	return this.getStatistics(key).Value
}

func (this *StarSession) updateStatistic(statistic *db.Statistics, addValue float64) float64 {
	//初始化的时候更新一波排名
	statistic.Value = statistic.Value + addValue
	statistic.UpdateTime = time.Now()
	return statistic.Value
}

func (this *StarSession) getStatistics(id int32) *db.Statistics {
	for _, curr := range this.Statistics {
		if curr.ID == id {
			return curr
		}
	}
	statistic := &db.Statistics{
		ID:         id,
		Value:      0,
		UpdateTime: time.Now(),
	}
	this.Statistics = append(this.Statistics, statistic)
	return statistic

}

//新增统计数据
func (this *StarSession) AddStatisticsValue(id int32, value float64, param int32) *db.Statistics {
	statistic := this.getStatistics(id)
	//过了刷新时间需要清理统计数据
	newValue := int32(this.updateStatistic(statistic, value))

	if id == constant.STAR_STATISTIC_LOOT_BELIEVER || id == constant.STAR_STATISTIC_LOOT_FAITH || id == constant.STAR_STATISTIC_ATK_BUILDING {
		cache.StarCache.SetMutualTimes(this.ID, cache.StarCache.GetMutualTimes(this.ID) + 1)
		if constant.IsHistoryAttackThreshold(newValue) {
			this.addStarHistory(constant.STAR_HISTORY_ATTACK, newValue, 0, "")
		}
	} else if id == constant.STAR_STATISTIC_BE_LOOT_BELIEVER || id == constant.STAR_STATISTIC_BE_LOOT_FAITH || id == constant.STAR_STATISTIC_BE_ATK_BUILDING {
		cache.StarCache.SetBeMutualTimes(this.ID, cache.StarCache.GetBeMutualTimes(this.ID) + 1)
		if constant.IsHistoryAttackThreshold(newValue) {
			this.addStarHistory(constant.STAR_HISTORY_BE_ATTACK, newValue, 0, "")
		}
	} else if id == constant.STAR_STATISTIC_TYPE_EVENT {
		if !this.containsStarHistory(constant.STAR_HISTORY_NEW_EVENT, param) {
			this.addStarHistory(constant.STAR_HISTORY_NEW_EVENT, param, 0, "")
		}
	}

	return statistic
}

func (this *StarSession) addStarHistory(historyID int32, param1 int32, param2 int32, param3 string) bool {
	record := &db.History{
		ID:     historyID,
		Param1: param1,
		Param2: param2,
		Param3: param3,
		Time:   time.Now(),
	}
	this.History = append(this.History, record)
	return true
}


func (this *StarSession) containsStarHistory(historyID int32, param int32) bool {
	for _, history := range this.History {
		if history.ID == historyID && history.Param1 == param {
			return true
		}
	}
	return false
}




////离线更新星球在线时长
//func (this *StarSession) UpdateOnlineTime(beginTime time.Time) {
//	if (!this.starBeginTime.IsZero()) {
//		beginTime = this.starBeginTime
//		this.starBeginTime = time.Time{}
//	}
//	this.starStatistics.OnlineTime += int32(time.Now().Sub(beginTime).Seconds())
//}

//
//func (this *StarSession) GetAttackNum() int32 {
//	return cache.UserCache.GetATKHistory(this.uid)
//}
//
//func (this *StarSession) UpdateBuildNum(addNum int32) {
//	this.starStatistics.BuildNum += addNum
//}
//
//func (this *StarSession) UpdateUpgradeBelieverNum(addNum int32) {
//	this.starStatistics.UpgradeBelieverNum += addNum
//}
//
//func (this *StarSession) GetOnlineTime() int32 {
//	return this.starStatistics.OnlineTime
//}

//
//func (this *StarSession) GetBuildNum() int32 {
//	return this.starStatistics.BuildNum
//}
//
//func (this *StarSession) GetUpgradeBelieverNum() int32 {
//	return this.starStatistics.UpgradeBelieverNum
//}

func (this *StarSession) GetHistoryTime() []int64 {
	var historyTime = []int64{}
	for _, reward := range this.CivilizationReward {
		historyTime = append(historyTime, reward.Time.Unix())
	}
	return historyTime
}

//设置机器人帮助修理的时间
//func (this *StarSession) SetRobotHelpRepairTime(buildType int32) {
//	building := this.getBuilding(buildType)
//	leftTime := building.BrokenTime - time.Now().Unix()
//	buildingConf := conf.GetBuildingConf(building.GetStarType(), building.Type, building.Level)
//	if buildingConf == nil {
//		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
//	}
//	shareTime := int64(buildingConf.BoomTime) - leftTime
//	//log.Info("shareTime: %v", shareTime)
//	if shareTime + 30*60 < leftTime {
//		randRobotInternal := character.RandInt32Scop(0,int32(leftTime))
//		this.getBuilding(buildType).RobotHelpTime = time.Unix(time.Now().Unix() + 30*60 + int64(randRobotInternal), 0)
//		//log.Info("RobotHelpTime: %v", this.getBuilding(buildType).RobotHelpTime)
//	}
//	//this.getBuilding(buildType).RobotHelpTime = time.Unix(time.Now().Unix() + 15 , 0)
//	//log.Info("RobotHelpTime: %v", this.getBuilding(buildType).RobotHelpTime)
//}

func (this *StarSession) GetSettleTime() int64 {
	return time.Now().Unix() - this.CreateTime.Unix()
}

/*-----------------------------压测接口--------------------------------------------------*/
func (this *StarSession) SetBuildingsLevel(totalLevel int32) []int32 {
	//for _, build := range this.Building {
	//	build.Level = totalLevel
	//}
	for _, build := range this.Building {
		build.Level = 0
	}

	for totalLevel > 0 {
		for _, build := range this.Building {
			build.Level ++
			totalLevel --
			if totalLevel <= 0 {
				break
			}
		}
	}
	this.updateBuildingLevel(true)
	this.setDirty()
	result := make([]int32,5)
	for index, build := range this.Building {
		result[index] = build.Level
	}
	return result
}

func (this *StarSession) SetBelieversLevel(totalLevel int32) []*protocol.BelieverInfo{
	//this.Believer = []*db.DBBeliever{}
	for _ , believer := range this.Believer {
		believer.Num = 0
	}
	var level int32
	for level = 5; level > 0; level-- {
		num := int32(totalLevel/level)
		totalLevel = totalLevel%level
		if num > 0 {
			believerID := "b01" + character.Int32ToString(level) + "1"
			believerID = replaceStar(this.Type, believerID)
			this.addBeliever(believerID, num)
		}
		if totalLevel == 0 {
			break
		}
	}
	this.updateBuildingLevel(true)
	this.setDirty()


	var result []*protocol.BelieverInfo
	for _, believer := range this.Believer {
		result = append(result, believer.BuildProtocol())
	}
	return result

}