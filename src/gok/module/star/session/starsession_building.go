package session

import (
	"aliens/common/character"
	"gok/constant"
	"gok/module/star/cache"
	"gok/module/star/conf"
	"gok/module/star/db"
	"gok/module/star/util"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"gok/service/rpc"
	"time"
)

//升级星球建筑
func (this *StarSession) UpgradeBuilding(message *protocol.BuildStarBuilding) (*db.DBBuilding, bool, int32, int32, *protocol.CivilizationInfo) {
	buildingType := message.GetBuildingType()
	upgradeLevel := message.GetLevel()
	//believerIDs := message.GetBelieverId()
	faith := message.GetFaith()
	//guide := message.GetGuide()
	//禁用过程中只允许引导操作
	//if this.Disable && !guide {
	//	exception.GameException(exception.ILL_BUILD_CONDITION)
	//}
	building := this.getBuilding(buildingType)

	if building == nil || building.Level+1 != upgradeLevel {
		exception.GameException(exception.ILL_BUILD_CONDITION)
	}
	if building.Level >= conf.GetMaxBuildingLevel(this.Type, buildingType) {
		exception.GameException(exception.BUILDING_MAXLEVEL)
	}

	if this.IsOtherBuildingState(buildingType) {
		exception.GameException(exception.STAR_OTHER_BUILDING_EXIST_STATE)
	}

	buildingConf := conf.GetBuildingConf(this.Type, building.Type, upgradeLevel)
	if buildingConf == nil {
		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
	}

	if this.CivilizationLevel < buildingConf.RequireCivilizationLevel {
		exception.GameException(exception.CIVILIZATION_LEVEL_NOT_ENOUGH)
	}

	cost := int32(buildingConf.UpgradeConsumption * (1 + this.Buff[constant.BUFF_UNR_COST]))
	//cost := int32(buildingConf.UpgradeConsumption)
	if cost > faith {
		exception.GameException(exception.FAITH_NOT_ENOUGH)
	}
	//this.ensureBelievers(believerIDs)
	//var accTime int32 = 0
	//for _, believerID := range believerIDs {
	//	accConf := conf.GetBuildingUpgradeAcc(conf.GetGameObjectLevel(believerID), upgradeLevel)
	//	if accConf == nil {
	//		exception.GameException(exception.STAR_ACC_INVALID_BELIEVER)
	//	}
	//	accTime += accConf.DecreaseTime
	//}


	levelInfo := building.GetCurrLevelInfo()
	//levelInfo.AddBelieverCosts(believerIDs)
	levelInfo.FaithCost += cost
	//this.AddStatisticsValue(constant.STAR_STATISTIC_TYPE_BUILDNUM, 1, 0)

	//changeBeliever := this.decBelieverByLevel(buildingConf.UpgradeBelieverLevel, buildingConf.UpgradeBelieverNumber)
	//PushBelieverInfo(this.Owner, believer)
	buildingRatio := 1 / (1 + this.Buff[constant.BUFF_UNR_TIME])
	building.UpdateUpgradeTime(int32(float32(buildingConf.BuildTime) * buildingRatio), this.Buff)
	//building.UpdateUpgradeTime(buildingConf.BuildTime)

	//building.UpgradeBelieverCost += believerNum
	upgrade, civilizationInfo, _ := this.dealUpgradeBuilding(building, false)
	var powerLimit int32 = 0
	if upgrade {
		powerLimit = this.GetPowerLimit()
	}
	//this.decBelievers(believerIDs)
	this.setDirty()
	return building, upgrade, powerLimit, cost, civilizationInfo
}

func (this *StarSession) AccUpgradeStarBuilding(uid int32, buildingType int32, believerIDs []string, guide bool) *protocol.AccUpdateStarBuildRet {
	//禁用过程中是不允许升级的
	//if this.Disable && !guide {
	//	exception.GameException(exception.STAR_UPDATE_BUILD_SP_FAILED)
	//}

	building := this.getBuilding(buildingType)
	if building == nil {
		exception.GameException(exception.STAR_UPDATE_BUILD_SP_FAILED)
	}

	//hjl 损坏和修理中是不能升级的
	if building.IsBroken() || building.IsRepairing() {
		exception.GameException(exception.STAR_UPDATE_BUILD_SP_FAILED)
	}

	//不在升级过程中
	if !building.IsUpgrading() {
		exception.GameException(exception.STAR_UPDATE_BUILD_SP_FAILED)
	}

	//if this.IsOtherBuildingState(buildingType) {
	//	exception.GameException(exception.STAR_OTHER_BUILDING_EXIST_STATE)
	//}

	buildingConf := conf.GetBuildingConf(this.Type, buildingType, building.Level+1)
	//找不到配置表
	if buildingConf == nil {
		exception.GameException(exception.STAR_UPDATE_BUILD_SP_FAILED)
	}

	this.ensureBelievers(believerIDs)
	var accTime int32 = 0
	for _, believerID := range believerIDs {
		accConf := conf.GetBuildingUpgradeAcc(conf.GetGameObjectLevel(believerID), building.Level+1)
		if accConf == nil {
			exception.GameException(exception.STAR_ACC_INVALID_BELIEVER)
		}
		accTime += accConf.DecreaseTime
	}

	//buff 加成
	accTime = int32(float32(accTime) * (1 + this.Buff[constant.BUFF_BELIEVER_TIME]))

	levelInfo := building.GetCurrLevelInfo()
	levelInfo.AddBelieverCosts(believerIDs)
	//building.UpgradeBelieverCost += believerNum
	upgrade, civilizationInfo := this.dealAccUpgradeTime(building, accTime)
	this.decBelievers(believerIDs)

	//加速成功走完引导，取消禁用
	//this.Disable = false
	this.setDirty()
	resp := &protocol.AccUpdateStarBuildRet{
		Done:         upgrade,
		Uid:          uid,
		BuildingType: buildingType,
		Level:        building.Level,
		PowerLimit:   this.GetPowerLimit(),
		UpdateTime:   building.UpdateTime.Unix(),
		//BelieverNum:building.UpgradeBelieverCost),
		//ItemGroove: building.BuildGroovesProtocol(),
		CivilizationInfo:civilizationInfo,
	}
	return resp
}

func (this *StarSession) UpgradeStarBuildEnd(uid int32, buildingType int32) *protocol.UpdateStarBuildEndRet {
	building := this.getBuilding(buildingType)
	if building == nil {
		exception.GameException(exception.STAR_UPDATE_BUILD_END_FAILED)
	}
	//hjl 损坏和修理中是不能升级的
	if building.IsBroken() {
		exception.GameException(exception.STAR_BUILDING_IS_BROKEN)
	}
	if building.IsRepairing() {
		exception.GameException(exception.STAR_BUILDING_IS_REPAIRE)
	}
	if !building.IsUpgrading() { //不在升级过程中哦
		exception.GameException(exception.STAR_BUILDING_IS_NOT_UPDATE)
	}

	upgrade, civilization, itemID := this.dealUpgradeBuilding(building, true)

	var powerReward int32 = 0
	if upgrade {
		powerReward = conf.GetBuildingPowerReward(this.Type, buildingType, building.Level)
		this.DealItemGroupOpen()
		this.setDirty()
	}

	resp := &protocol.UpdateStarBuildEndRet{
		Done:         upgrade,
		Uid:          uid,
		PowerReward:  powerReward,
		BuildingType: buildingType,
		Level:        building.Level,
		UpdateTime:   building.UpdateTime.Unix(),
		PowerLimit:   this.GetPowerLimit(),
		//ItemGroove:   building.BuildGroovesProtocol(),
		CivilizationInfo: civilization,
		ItemID: 	  itemID,
	}
	return resp
}

func (this *StarSession) RepairStarBuilding(uid int32, buildingType int32, faith int32, /* believerIDs []string*/) *protocol.RepairStarBuildRet {
	building := this.getBuilding(buildingType)
	if building == nil {
		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
	}
	if !building.IsBroken() {
		exception.GameException(exception.STAR_REPAIR_BUILD_FAILED)
	}

	if this.IsOtherBuildingState(buildingType) {
		exception.GameException(exception.STAR_OTHER_BUILDING_EXIST_STATE)
	}

	buildingConf := conf.GetBuildingConf(this.Type, building.Type, building.Level)
	if buildingConf == nil { //找不到配置表
		exception.GameException(exception.STAR_REPAIR_BUILD_FAILED)
	}

	cost := int32(float32(buildingConf.RepairConsumption) * (1 + this.Buff[constant.BUFF_UNR_COST]))
	//cost := buildingConf.RepairConsumption
	if faith < cost {
		exception.GameException(exception.FAITH_NOT_ENOUGH)
	}

	//this.ensureBelievers(believerIDs)
	//var accTime int32 = 0
	//for _, believerID := range believerIDs {
	//	accConf := conf.GetBuildingUpgradeAcc(conf.GetGameObjectLevel(believerID), building.Level)
	//	if accConf == nil {
	//		exception.GameException(exception.STAR_ACC_INVALID_BELIEVER)
	//	}
	//	accTime += accConf.DecreaseTime
	//}

	levelInfo := building.GetCurrLevelInfo()
	//levelInfo.AddBelieverCosts(believerIDs)
	levelInfo.FaithCost += cost

	repairTimeRatio := 1 / (1 + this.Buff[constant.BUFF_UNR_TIME])
	building.UpdateRepairTime(int32(float32(buildingConf.RepairTime) * repairTimeRatio), this.Buff) //更新修理时间

	//building.UpdateRepairTime(buildingConf.RepairTime) //更新修理时间
	repair, civilizationInfo := this.dealUpdateRepairing(building, false)
	//this.decBelievers(believerIDs)

	this.setDirty()
	resp := &protocol.RepairStarBuildRet{
		Done:          repair,
		Uid:           uid,
		Cost:          cost,
		BuildingType:  buildingType,
		RepairTime:    building.RepairTime.Unix(),
		BuildingLevel: building.Level,
		CivilizationInfo:civilizationInfo,
	}
	return resp
}

func (this *StarSession) AccRepairStarBuilding(buildingType int32, believerIDs []string) (bool, int64, int32, int32, *protocol.CivilizationInfo) {
	building := this.getBuilding(buildingType)
	if building == nil {
		exception.GameException(exception.STAR_REPAIR_BUILD_FAILED)
	}

	if !building.IsRepairing() { //没有进入维修状态哦
		exception.GameException(exception.STAR_REPAIR_BUILD_FAILED)
	}


	//if this.IsOtherBuildingState(buildingType) {
	//	exception.GameException(exception.STAR_OTHER_BUILDING_EXIST_STATE)
	//}

	buildingConf := conf.GetBuildingConf(this.Type, buildingType, building.Level)
	if buildingConf == nil { //找不到配置表
		exception.GameException(exception.STAR_REPAIR_BUILD_FAILED)
	}
	var accTime int32 = 0

	this.ensureBelievers(believerIDs)
	for _, believerID := range believerIDs {
		accConf := conf.GetBuildingUpgradeAcc(conf.GetGameObjectLevel(believerID), building.Level)
		if accConf == nil {
			exception.GameException(exception.STAR_ACC_INVALID_BELIEVER)
		}
		accTime += accConf.DecreaseTime
	}

	accTime = int32(float32(accTime) * (1 + this.Buff[constant.BUFF_BELIEVER_TIME]))

	levelInfo := building.GetCurrLevelInfo()
	levelInfo.AddBelieverCosts(believerIDs)


	repairSucc, repairTime, civilizationInfo := this.dealAccRepairTime(building, accTime)
	//修理成功需要清0
	this.decBelievers(believerIDs)
	this.setDirty()
	return repairSucc, repairTime, 0, building.Level, civilizationInfo
}

func (this *StarSession) RepairBuildingEnd(buildingType int32) (bool, *protocol.CivilizationInfo, int64, int32) {
	building := this.getBuilding(buildingType)
	if building == nil {
		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
	}
	done, civilizationInfo := this.dealUpdateRepairing(building, true)
	return done, civilizationInfo, building.RepairTime.Unix(), building.Level
}

func (this *StarSession) dealUpdateRepairing(building *db.DBBuilding, repaired bool) (bool, *protocol.CivilizationInfo) {
	repaired = building.UpdateRepairing(repaired)
	var civilizationInfo *protocol.CivilizationInfo = nil
	if repaired {
		config := conf.GetBuildingConf(this.Type, building.Type, building.Level)
		if config != nil {
			civilizationInfo = this.TakeInCivilization(config.CivRepairIncome)
		}
		cache.StarCache.DelStarHelpRepair(this.ID, building.Type)
	}
	return repaired, civilizationInfo
}

func (this *StarSession) dealAccRepairTime(building *db.DBBuilding, accTime int32) (bool, int64, *protocol.CivilizationInfo) {
	repaired := building.AccRepairTime(accTime,false, this.Buff)
	var civilizationInfo *protocol.CivilizationInfo = nil

	if repaired {
		config := conf.GetBuildingConf(this.Type, building.Type, building.Level)
		if config != nil {
			civilizationInfo = this.TakeInCivilization(config.CivRepairIncome)
		}
	}
	return repaired, building.RepairTime.Unix(), civilizationInfo
}

func (this *StarSession) dealAccUpgradeTime(building *db.DBBuilding, accTime int32) (bool, *protocol.CivilizationInfo) {
	upgrade, firstUpgrade := building.AccUpgradeTime(accTime, false, this.Buff)
	var civilizationInfo *protocol.CivilizationInfo = nil
	//第一次升级才奖励文明度
	if upgrade {
		if firstUpgrade {
			config := conf.GetBuildingConf(this.Type, building.Type, building.Level)
			if config != nil {
				civilizationInfo = this.TakeInCivilization(config.CivBuildIncome)
			}
		}
		this.updateBuildingLevel(true)
	}
	return upgrade, civilizationInfo
	//AccUpgradeTime(accTime int32) bool
}

func (this *StarSession) dealUpgradeBuilding(building *db.DBBuilding, upgrade bool) (bool, *protocol.CivilizationInfo, int32) {
	upgrade, firstUpgrade := building.UpdateUpgrading(upgrade, this.Buff)
	var civilizationInfo *protocol.CivilizationInfo = nil
	var itemID int32
	//第一次升级才需要加文明度
	if upgrade {
		this.updateBuildingLevel(true)
		if firstUpgrade {
			config := conf.GetBuildingConf(this.Type, building.Type, building.Level)
			if config != nil {
				civilizationInfo = this.TakeInCivilization(config.CivBuildIncome)
			}
			// 到达第一个圣物组合中圣物的等级时指定给该组合中的圣物
			index := conf.GetIndexOfFirstGroupItem(this.buildingTotalLevel)
			if index >= 0 {
				groupID := this.getItemGroupFirstItemID()
				groupBase := conf.Base.ItemGroup[groupID]
				if groupBase == nil {
					exception.GameException(exception.ITEM_GOURP_BASE_NOT_FOUND)
				}
				itemID = groupBase.Content[index]
			}
		}
	}
	return upgrade, civilizationInfo, itemID
}

func (this *StarSession) CancelRepairStarBuilding(buildingType int32) *protocol.CancelRepairStarBuildRet{
	building := this.getBuilding(buildingType)
	if building == nil {
		exception.GameException(exception.STAR_REPAIR_BUILD_FAILED)
	}

	if !building.IsRepairing() { //没有进入维修状态哦
		exception.GameException(exception.STAR_REPAIR_BUILD_FAILED)
	}


	//if this.IsOtherBuildingState(buildingType) {
	//	exception.GameException(exception.STAR_OTHER_BUILDING_EXIST_STATE)
	//}

	buildingConf := conf.GetBuildingConf(this.Type, buildingType, building.Level)
	if buildingConf == nil { //找不到配置表
		exception.GameException(exception.STAR_BUILDING_CONFIG_NOT_FOUND)
	}
	if building.RepairTime.Before(time.Now()) {
		exception.GameException(exception.STAR_BUILDING_CANCEL_REPAIR_FAILED)
	}
	building.CleanRepairTime()
	building.SetBroken(this.Buff, true)
	return &protocol.CancelRepairStarBuildRet{
		BackFaith:int32(buildingConf.RepairConsumption/2),
		BuildingType:building.Type,
		BuildingLevel:building.Level,
		RepairTime:building.RepairTime.Unix(),
		BrokenTime:building.BrokenTime,
	}
}

func (this *StarSession) CancelUpgradeStarBuilding(buildingType int32) *protocol.CancelUpgradeStarBuildRet{
	building := this.getBuilding(buildingType)
	if building == nil {
		exception.GameException(exception.STAR_REPAIR_BUILD_FAILED)
	}

	if !building.IsUpgrading() { //没有进入升级状态哦
		exception.GameException(exception.STAR_BUILDING_IS_NOT_UPDATE)
	}


	//if this.IsOtherBuildingState(buildingType) {
	//	exception.GameException(exception.STAR_OTHER_BUILDING_EXIST_STATE)
	//}

	buildingConf := conf.GetBuildingConf(this.Type, buildingType, building.Level+1)
	if buildingConf == nil { //找不到配置表
		exception.GameException(exception.STAR_BUILDING_CONFIG_NOT_FOUND)
	}
	if building.UpdateTime.Before(time.Now()) {
		exception.GameException(exception.STAR_BUILDING_CANCEL_UPGRADE_FAILED)
	}
	building.CleanUpgradeTime()
	return &protocol.CancelUpgradeStarBuildRet{
		BackFaith:int32(buildingConf.UpgradeConsumption/2),
		BuildingType:building.Type,
		BuildingLevel:building.Level,
		UpdateTime:building.UpdateTime.Unix(),
	}
}

func (this *StarSession) updateBuildingLevel(updateSearch bool) int32 {
	var level int32 = 0
	for _, building := range this.Building {
		level += building.Level
	}
	this.buildingTotalLevel = level
	//星球没有被禁用才能被搜索到
	if updateSearch && this.Owner > 0 {
		if level > this.BuildingExMaxLevel {
			this.BuildingExMaxLevel = level
			cache.StarCache.SetBuildingExMaxLevel(this.ID, level)
			//解锁星球flag
			this.checkLockFlags()
		}
		rpc.SearchServiceProxy.UpdateData(constant.SEARCH_OPT_UPDATE_BUILDING, this.ID, this.buildingTotalLevel)
		//searcher.Index.UpdateStarBuildingLevel(this.ID, this.buildingTotalLevel)
	}
	cache.StarCache.SetBuildingAllLevel(this.ID, this.buildingTotalLevel)
	return this.buildingTotalLevel
}

func (this *StarSession) GetHelpBuildsInfo(buildType int32) []*protocol.HelpRepairBuildInfo{
	var retInfo []*protocol.HelpRepairBuildInfo
	helpBuildMapping := cache.StarCache.GetAllStarBuildRepair(this.ID)
	if helpBuildMapping == nil {
		return retInfo
	}
	if buildType == constant.ALL_BUILDING_TYPE {
		for buildTypeStr, helperID := range helpBuildMapping {
			retInfo = append(retInfo, this.getHelpBuild(character.StringToInt32(buildTypeStr), character.StringToInt32(helperID)))
		}
	} else {
		helperID, ok := helpBuildMapping[character.Int32ToString(buildType)]
		if !ok {
			exception.GameException(exception.HELP_REPAIR_BUILD_NOT_FOUND)
		}
		retInfo = append(retInfo, this.getHelpBuild(buildType, character.StringToInt32(helperID)))
	}
	return retInfo
}

func (this *StarSession) getHelpBuild(buildType int32, helperID int32) *protocol.HelpRepairBuildInfo {
	building := this.getBuilding(buildType)
	if building == nil {
		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
	}
	if building.IsBroken() || building.IsRepairing() {
		return &protocol.HelpRepairBuildInfo{
			BuildingType:buildType,
			RepairTime:building.RepairTime.Unix(),
			BrokenTime:building.BrokenTime,
			HelpID:helperID,
		}
	}
	cache.StarCache.DelStarHelpRepair(this.ID, buildType)
	return nil
}

func (this *StarSession) HelpRepairBuild(buildType int32, helperID int32) (bool, int32, int32){
	building := this.getBuilding(buildType)
	if building.IsRepairing() || building.IsBroken() {
		if cache.StarCache.GetStarBuildRepair(this.ID, buildType) == character.Int32ToString(0) {
			building.CalFaithUpdateTime(time.Now(), this.Buff)
			building.RepairTime = time.Now()
			building.BrokenTime = 0
			this.setDirty()
			cache.StarCache.SetStarBuildRepair(this.ID, buildType, helperID)
			return true, building.GetStarType(), building.Level
		}
	}
	return false, 0, 0
}

func(this *StarSession) PublicHelpRepairBuild(buildType int32) int32 {

	ret := cache.StarCache.SetStarBuildRepair(this.ID, buildType, 0)
	if ret {
		//this.SetRobotHelpRepairTime(buildType)
		return buildType
	}
	return 0
}

func (this *StarSession) AtkStar(message *protocol.AtkStarBuilding) (int32, int32, bool, bool) { //请求进攻目标星球建筑 wjl20170622
	if message.GetAttackUid() == message.GetDestUid() { //不允许自己打自己哟
		exception.GameException(exception.STAR_CANNOT_ATK_SELF)
	}
	if this.DecBuildingShield() {
		return 0, 0, true, false
	}
	//获取目标建筑物 ( 注：由于客户端 type 在本地写死， 无法获取到ID所以只能暂时先用 type 来获取建筑数据，以后得要修改
	building := this.getBuilding(message.GetBuildingID())
	if building == nil || building.Level == 0{ //找不到建筑物哦
		return 50, 0, false, false
	}

	//基础奖励
	var baseFaith int32 = 50 * building.Level + 1

	//机器人直接返回信仰
	if this.IsRobot() {
		return baseFaith + int32(message.GetFaithRatio()*float32(building.Faith)), 0, false, false
	}

	faith := int32(message.GetFaithRatio() * float32(building.Faith))


	var itemID int32 = 0
	//prob := (10 + float64(message.GetBuildingLevel())) / (100 + float64(this.getBuildingAllLevel()))
	////成功
	//if rand.Float64() <= prob {
	//	itemID = building.TakeoutRandomItem()
	//}
	if this.IsBuildingAllMaxLevel() || !this.IsFlagUnlock(constant.STAR_FLAG_MUTUAL) {
		return baseFaith + faith, itemID, false, true
	}
	building.TakeoutFaith(faith)

	if building.IsUpgrading() {
		//清理建造时间
		building.CleanUpgradeTime()
	}
	if building.IsRepairing() {
		building.CleanRepairTime()
		// 清掉求助维修帮助的人
		if cache.StarCache.ExistStarHelpRepair(this.ID, building.Type) {
			cache.StarCache.SetStarBuildRepair(this.ID, building.Type, 0)
			//this.SetRobotHelpRepairTime(building.Type)
		}
	}

	firstBroken := building.SetBroken(this.Buff, false)
	if firstBroken {
		rpc.PassportServiceProxy.WechatEventPush(this.Owner, constant.EVENT_BUILD_ATTACK, 0)
	}

	this.setDirty()

	uid := this.Owner
	push := util.BuildBuildingInfoPush(uid, this.ID, this.Type, building)
	rpc.UserServiceProxy.Push(uid, push)
	return baseFaith + faith, itemID, false, false
}

func (this *StarSession) GetAllBuildingFaith(buildingType int32) []*protocol.BuildingFaith{
	var buildingsFaith []*protocol.BuildingFaith
	buildings := this.Building
	if buildingType == constant.ALL_BUILDING_TYPE {
		buildingsFaith = make([]*protocol.BuildingFaith, len(buildings))
		for index, building := range buildings {
			faith, faithUpdateTime := this.GetBuildingFaith(building.Type, false)
			buildingFaith := &protocol.BuildingFaith{
				BuildingType:building.Type,
				BuildingFaith:faith,
				FaithUpdateTime:faithUpdateTime,
			}
			buildingsFaith[index] = buildingFaith
		}
	} else {
		buildingsFaith = make([]*protocol.BuildingFaith, 1)
		faith, faithUpdateTime := this.GetBuildingFaith(buildingType, false)
		buildingFaith := &protocol.BuildingFaith{
			BuildingType:buildingType,
			BuildingFaith:faith,
			FaithUpdateTime:faithUpdateTime,
		}
		buildingsFaith[0] = buildingFaith
	}
	return buildingsFaith
}

//处理建筑信仰值  draw 是否领取
func (this *StarSession) GetBuildingFaith(buildingType int32, draw bool) (int32, int64) {
	building := this.getBuilding(buildingType)
	if building == nil {
		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
	}
	result := building.CheckoutFaithBuff(this.Type, draw, this.Buff)
	this.setDirty()
	return result, building.FaithUpdateTime.Unix()
}

func (this *StarSession) DrawAllBuildingFaith() int32 {

	var faith int32 = 0
	for _, building := range this.Building {
		faith += building.CheckoutFaithBuff(this.Type, true, this.Buff)
	}
	this.setDirty()
	return faith
}

func (this *StarSession) updateBuildingExMax() {
	var exLevel int32 = 0
	for _, building := range this.Building {
		exLevel += building.GetMaxLevel()
	}
	this.BuildingExMaxLevel = exLevel
	cache.StarCache.SetBuildingExMaxLevel(this.ID, exLevel)
}

//星球上的建筑是否都达到最大等级
func (this *StarSession) IsBuildingAllMaxLevel() bool {
	for _, v := range this.Building {
		if v.Level < constant.MAX_BUILDING_LEVEL {
			return false
		}
	}
	return true
}

//是否有损坏的建筑
func (this *StarSession) hasBrokenBuilding() bool {
	for _, v := range this.Building {
		if v.Level > 0 && v.IsBroken() {
			return true
		}
	}
	return false
}

//是否拥有建筑
func (this *StarSession) hasBuilding() bool {
	for _, v := range this.Building {
		if v.Level > 0 || v.Exist {
			return true
		}
	}
	return false
}

func (this *StarSession) IsOtherBuildingState(buildingType int32) bool {
	for _, building := range this.Building {
		if building.Type == buildingType {
			continue
		}
		if building.IsRepairing() || building.IsUpgrading() {
			return true
		}
	}
	return false
}

func (star *StarSession) HasBuilding() bool {
	return star.buildingTotalLevel > 0
}

func (star *StarSession) CheckBuildingLevel() bool {
	//if star.IsBuildingAllMaxLevel() {
	//	exception.GameException(exception.SEARCH_BUILDING_MAXLEVEL)
	//}
	return star.buildingTotalLevel > 0
}

func (this *StarSession) UpdateLevelNum() []int32 {
	levelNum := []int32{0,0,0,0,0}
	for _,build := range this.Building {
		for i := 0 ; i <  int(build.GetMaxLevel()) ; i++ {
			levelNum[i]+=1
		}
	}
	return levelNum
}

//根据建筑ID获取星球建筑
func (this *StarSession) getBuildingByID(id int32) *db.DBBuilding {
	for _, obj := range this.Building {
		if obj.ID == id {
			return obj
		}
	}
	return nil
}

//根据建筑类型获取星球建筑
func (this *StarSession) getBuilding(buildingType int32) *db.DBBuilding {
	for _, roleStarBuilding := range this.Building {
		if roleStarBuilding.Type == buildingType {
			return roleStarBuilding
		}
	}
	return nil
}

func (this *StarSession) getBuildingAllLevel() int32 {
	var level int32 = 0
	for _, roleStarBuilding := range this.Building {
		level += roleStarBuilding.Level
	}
	return level
}
