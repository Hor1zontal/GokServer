package manager

import (
	"aliens/common/util"
	"aliens/log"
	"gok/constant"
	"gok/module/game/cache"
	"gok/module/game/conf"
	"gok/module/game/db"
	"gok/service/exception"
	"gok/service/msg/protocol"
)

type RoleDialManager struct {
	uid          int32
	guideHasDial []int32
	dialTimes    int32	//转转盘的次数

	buildExMaxLevel int32
	guideRandoms map[int32]int32
	multipleDialTimes int32
	multipleDialID int32 //翻倍奖励的DialID

	starFlags map[int32]int32
}


func (this *RoleDialManager) Init(role *db.DBRole) {
	this.uid = role.UserID
	this.randomMultipleDialTimes()
	this.guideHasDial = role.GuideDial
	this.guideRandomsUpdate()
	//this.GuideDialInit(role.GuideDial)
}


func (this *RoleDialManager) Update(role *db.DBRole) {
	role.GuideDial = this.guideHasDial
}


func (this *RoleDialManager) guideRandomsUpdate() {
	if this.guideRandoms == nil {
		this.guideRandoms = make(map[int32]int32)
	}
	this.buildExMaxLevel = this.getExMaxBuildingLevel()
	if this.buildExMaxLevel > conf.DATA.GUIDE_MAX_REQUIRE_LEVEL || this.buildExMaxLevel < conf.DATA.GUIDE_MIN_REQUIRE_LEVEL {
		return
	}
	randomMapping := conf.DATA.GUIDE_DIAL_LEVEL_ID_WEIGHT_MAPPING[this.buildExMaxLevel]
	if randomMapping == nil {
		exception.GameException(exception.DIAL_GUIDE_NOT_FOUND)
	}
	this.guideRandoms = util.CopyMap(randomMapping)
	for id := range this.guideRandoms {
		for _, hasId := range this.guideHasDial {
			if id == hasId {
				delete(this.guideRandoms, id)
			}
		}
	}
}

func (this *RoleDialManager) getExMaxBuildingLevel() int32 {
	starID := cache.StarCache.GetUserActiveStar(this.uid)
	if starID == 0 {
		return 0
	}
	buildingLevel := cache.StarCache.GetBuildingExMaxLevel(starID)
	return buildingLevel
}

func (this *RoleDialManager) RandomDial(civilLevel int32, buildingLevel int32, dialData map[int32]*conf.DialLimitData) int32 {
	//dialID := util.RandomWeight(this.dialIDWeightMapping)
	//resetID,ret := this.isResetDialWeight()
	//if ret {
	//	this.resetDialWeight(resetID)
	//}
	////连续三次抽到一样的
	//if dialID == this.dialPastIDs[0] && dialID == this.dialPastIDs[1] {
	//	this.dialIDWeightMapping[dialID] = 0
	//} else {
	//	this.dialIDWeightMapping[dialID] -= 1
	//}
	//
	//this.pushDialID(dialID)

	//建筑等级前几次转盘根据initDial表随机
	dialID := this.randomGuideDial(buildingLevel)
	if dialID != 0 {
		return dialID
	}
	//普通随机

	confWeight := conf.DATA.DIAL_CIVILLEVEL_DIALID_WEIGHT_MAPPING[civilLevel]
	if confWeight == nil {
		exception.GameException(exception.DIAL_LIMIT_NOT_FOUND)
	}

	dialWeight := util.CopyMap(confWeight)

	for dialID := range dialWeight {
		//var lock bool
		var unlock = true
		if dialData[dialID] == nil {
			exception.GameException(exception.DIAL_LIMIT_NOT_FOUND)
		}
		switch dialData[dialID].Type {
		case constant.DIAL_GAYPOINT: unlock = this.IsFlagUnlock(constant.STAR_FLAG_GAYPOINT)
		case constant.QUEST_ROB_FAITH: unlock = this.IsFlagUnlock(constant.STAR_FLAG_LOOT_FAITH)
		case constant.QUEST_ROB_BELIEVER: unlock = this.IsFlagUnlock(constant.STAR_FLAG_LOOT_BELIEVER)
		case constant.QUEST_ATT_BUILDING: unlock = this.IsFlagUnlock(constant.STAR_FLAG_ATK_BUILDING)
		}
		if !unlock {
			//log.Info("delete type:%v", dialData[dialID].Type)
			delete(dialWeight, dialID)
		}
	}
	dialID = util.RandomWeight(dialWeight)

	//到达5级才能看发广告翻倍
	//log.Info("%v",this.IsFlagUnlock(constant.STAR_FLAG_AD_MULTIPLE))
	if this.IsFlagUnlock(constant.STAR_FLAG_AD_MULTIPLE) {
		this.addDialTimes()
	}

	return dialID
}

func (this *RoleDialManager) randomGuideDial(buildingLevel int32) int32 {
	if buildingLevel != this.buildExMaxLevel {
		this.guideRandomsUpdate()
	}
	if this.guideRandoms != nil && len(this.guideRandoms) != 0 {
		guideID := util.RandomWeight(this.guideRandoms)
		//log.Info("==========GuideID: %v", guideID)
		delete(this.guideRandoms, guideID)
		this.guideHasDial = append(this.guideHasDial, guideID)
		return conf.DATA.GUIDE_ID_DIALID_MAPPING[guideID]
	}
	return 0
}

func (this *RoleDialManager) randomMultipleDialTimes() {
	requireTimes := conf.DATA.MultipleReward
	if requireTimes == nil || len(requireTimes) != 2{
		exception.GameException(exception.GAME_BASE_NOT_FOUND)
	}
	this.multipleDialTimes = util.RandInt32Scop(requireTimes[0],requireTimes[1])
	//this.multipleDialTimes = 0
}

func (this *RoleDialManager) compareDialTimes() bool{
	if this.dialTimes >= this.multipleDialTimes {
		this.dialTimes = 0
		this.randomMultipleDialTimes()
		//this.resetDialAllWeight()
		return true
	}
	return false
}


func (this *RoleDialManager) GetMultipleReward(dialID int32) bool {
	if !this.compareDialTimes() {
		return false
	}
	multipleWeight := conf.DATA.MultipleWeightMapping
	if multipleWeight == nil || len(multipleWeight) != 2 {
		exception.GameException(exception.GAME_BASE_NOT_FOUND)
	}
	index := util.RandomWeight(conf.DATA.MultipleWeightMapping)
	if index == 0 {
		//有倍数奖励
		this.multipleDialID = dialID
		log.Info("dialTimes:%v multipleID%v intervalMultipleTimes%v", this.dialTimes, this.multipleDialID, this.multipleDialTimes)
		return true
	}
	return false
}

//func (this *RoleDialManager) GetDialTimes() int32 {
//	return this.dialTimes
//}

func (this *RoleDialManager) addDialTimes() {
	this.dialTimes += 1
}

func (this *RoleDialManager) DecDialTimes() {
	this.dialTimes -= 1
}

func (this *RoleDialManager) GetMultipleDialID() int32 {
	return this.multipleDialID
}

func (this *RoleDialManager) CleanMultipleDialID() {
	this.multipleDialID = 0
}


//--------------------star flags------------------------------------
func (this *RoleDialManager) UpdateStarFlags(flags []*protocol.FlagInfo) {
	if this.starFlags == nil {
		this.starFlags = make(map[int32]int32)
	}
	for _, flag := range flags {
		this.SetStarFlag(flag.Id, flag.Value)
	}
}

func (this *RoleDialManager) SetStarFlag(key int32, value int32) bool {
	//if constant.IsDialFlag(key) {
	//	this.starFlags[key] = value
	//}
	if this.GetStarFlagValue(key) != value {
		this.starFlags[key] = value
		return true
	}
	return false
}

////转盘是否解锁
//func (this *RoleDialManager) IsDialUnlock() bool {
//	return this.GetStarFlagValue(constant.STAR_FLAG_DIAL) == constant.FLAG_VALUE_UNLOCK
//}

//转盘中某一项是否解锁 return true:已解锁 false:没解锁
func (this *RoleDialManager) IsFlagUnlock(key int32) bool {

	return this.GetStarFlagValue(key) == constant.FLAG_VALUE_UNLOCK
}

func (this *RoleDialManager)LockEventIDs() []int32 {
	var result []int32
	if !this.IsFlagUnlock(constant.STAR_FLAG_LOOT_FAITH) {
		result = append(result, constant.EVENT_ID_LOOT_FAITH)
	}
	if !this.IsFlagUnlock(constant.STAR_FLAG_LOOT_BELIEVER) {
		result = append(result, constant.EVENT_ID_LOOT_BELIEVER)
	}
	if !this.IsFlagUnlock(constant.STAR_FLAG_ATK_BUILDING) {
		result = append(result, constant.EVENT_ID_ATK_BUILDING)
	}
	return result
}

func (this *RoleDialManager) GetStarFlagValue(key int32) int32 {
	return this.starFlags[key]
}