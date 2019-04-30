package session

import (

	"aliens/common/character"
	timeutil "aliens/common/util"
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

//-----------------------------------信徒接口--------------------------------------
//删除num个指定believerLevel等级的信徒
//func (this *StarSession) decBelieverByLevel( believerLevel int32, num int32 ) []*db.DBBeliever {
//	if (num < 0) {
//		exception.GameException(exception.STAR_BELIEVER_NOT_ENOUGH)
//	}
//	believers := this.getBelieverByLevel(believerLevel)
//	count := getBelieverCount(believers)
//	if (count < num) {
//		exception.GameException(exception.STAR_BELIEVER_NOT_ENOUGH)
//	}
//
//	//平均每种信徒需要扣除的数量
//	perDec := num / int32(len(believers))
//
//	//一轮分配平均扣除,不够的要保存回num
//	for _, believer := range believers {
//		remain := believer.Num - perDec
//		num -= perDec
//		if (remain >= 0) {
//			//数量足够
//			believer.Num = remain
//		} else {
//			//数量不够
//			believer.Num = 0
//			num -= remain
//		}
//	}
//	//平均扣除后还有剩余按顺序完继续扣除
//	if (num > 0) {
//		for _, believer := range believers {
//			remain := believer.Num - num
//			if (remain >= 0) {
//				believer.Num = remain
//			} else {
//				num = -remain
//			}
//		}
//	}
//	return believers
//}

func (star *StarSession) CheckBelieverLevel() bool {
	return star.believerCount > 0
}

func (this *StarSession) ensureBeliever(believerID string, num int32) {
	believer := this.getBeliever(believerID)
	if believer == nil || believer.Num < num {
		exception.GameException(exception.STAR_BELIEVER_NOT_ENOUGH)
	}
}

func (this *StarSession) ensureBelievers(believerIDs []string) {
	believerMapping := make(map[string]int32)
	for _, believerID := range believerIDs {
		believerMapping[believerID] += 1
	}
	for believerID, num := range believerMapping {
		this.ensureBeliever(believerID, num)
	}
}

func (this *StarSession) decBelievers(believerIDs []string) {
	believerMapping := make(map[string]int32)
	for _, believerID := range believerIDs {
		believerMapping[believerID] += 1
	}
	for believerID, num := range believerMapping {
		this.decBeliever(believerID, num)
		this.AddStatisticsValue(constant.STAR_STATISTIC_ACC_USE_BELIEVER, float64(num), 0)
	}
}

//删除num个指定believerID信徒
func (this *StarSession) decBeliever(believerID string, num int32) *db.DBBeliever {
	if num < 0 {
		exception.GameException(exception.STAR_BELIEVER_NOT_ENOUGH)
	}
	believer := this.getBeliever(believerID)
	if believer == nil || believer.Num < num {
		exception.GameException(exception.STAR_BELIEVER_NOT_ENOUGH)
	}

	limitNum := this.getAddBelieverCount()
	if limitNum <= 0 && num+limitNum > 0 {
		this.BelieverUpdateTime = time.Now()
	}

	believer.Num = believer.Num - num
	this.updateBelieverCount(true)
	return believer
}

//强制删除num个指定believerID信徒，不够就扣除完
func (this *StarSession) forceDecBeliever(believerID string, num int32) *db.DBBeliever {
	believer := this.getBeliever(believerID)
	if believer == nil {
		return nil
	}
	if believer.Num == 0 {
		return nil
	}
	limitNum := this.getAddBelieverCount()
	if limitNum <= 0 && num+limitNum > 0 {
		this.BelieverUpdateTime = time.Now()
	}
	if believer.Num < num {
		believer.Num = 0
	} else {
		believer.Num = believer.Num - num
	}
	this.updateBelieverCount(true)
	return believer
}

func getBelieverPrefix(starType int32) string {
	prefix := "b"
	if (starType <= 9) {
		prefix += "0"
	}
	prefix = prefix + character.Int32ToString(starType)
	return prefix
}

func (this *StarSession) addCustomBeliever(believerLevel int32, male bool, num int32) *db.DBBeliever {
	maleStr := "1"
	if !male {
		maleStr = "2"
	}
	believerID := getBelieverPrefix(this.Type) + character.Int32ToString(believerLevel) + maleStr
	return this.addBeliever(believerID, num)
}

//添加信徒
func (this *StarSession) addBeliever(believerID string, num int32) *db.DBBeliever {
	if num < 0 {
		exception.GameException(exception.STAR_BELIEVER_LIMIT)
	}
	believer := this.getBeliever(believerID)
	if believer == nil {
		believer = &db.DBBeliever{ID: believerID, Num: num}
		this.Believer = append(this.Believer, believer)
		return believer
	}
	believer.Num = believer.Num + num
	this.updateBelieverCount(true)
	return believer
}

//将超过当前文明度等级的信徒 转换为 当前星球文明度等级对应的信徒和信仰
func (this *StarSession) convertBelieverCivil( updateInfo []*protocol.BelieverInfo) (map[string]int32, int32){
	var faith int32 = 0
	results := make(map[string]int32)

	for _,believerInfo := range updateInfo {
		believerID := believerInfo.GetId()
		level := conf.GetBelieverLevel(believerID)
		requireLevel := conf.Base.CivilBelieverMapping[this.CivilizationLevel]
		if level > requireLevel {
			addFaithOne := conf.Base.BelieverFaithMapping[level] - conf.Base.BelieverFaithMapping[requireLevel]
			faith += addFaithOne * believerInfo.Num
			id := replaceLevel(requireLevel, believerInfo.GetId())
			results[id] += believerInfo.Num
		} else {
			results[believerInfo.GetId()] += believerInfo.GetNum()
		}
	}
	return results, faith
}


//删除当前星球believerID相同等级信徒num个
func (this *StarSession) decLevelBeliever(believerID string, num int32) *db.DBBeliever {
	believerID = replaceStar(this.Type, believerID)
	return this.forceDecBeliever(believerID, num)
}

//添加当前星球believerID相同等级的信徒num个
func (this *StarSession) addLevelBeliever(believerID string, num int32) *db.DBBeliever {
	believerID = replaceStar(this.Type, believerID)
	return this.addBeliever(believerID, num)
}

func replaceLevel(level int32, believerID string) string {
	levelStr := character.Int32ToString(level)
	return believerID[:3] + levelStr + believerID[4:]
}

func replaceStar(starType int32, believerID string) string {
	typeStr := character.Int32ToString(starType)
	if starType < 10 {
		typeStr = "0" + typeStr
	}
	return "b" + typeStr + believerID[3:]
}

//获取数量最少的信徒类型
func (this *StarSession) getMinBeliever(believerIDs []string) string {
	var minNum int32 = conf.Base.StarBelieverLimit
	var minBelieverID string = believerIDs[0]
	for _, believerID := range believerIDs {
		believer := this.getBeliever(believerID)
		if believer == nil {
			return believerID
		}
		if believer.Num < minNum {
			minNum = believer.Num
			minBelieverID = believer.ID
		}
	}
	return minBelieverID
}

//能否添加信徒
func (this *StarSession) canAddBeliever(num int32) bool {
	return this.believerCount+num <= conf.Base.StarBelieverLimit
}

func (this *StarSession) getAddBelieverCount() int32 {
	return conf.Base.StarBelieverLimit - this.believerCount
}

func (this *StarSession) isBelieverLimit() bool {
	return this.getAddBelieverCount() <= 0
}

//是否拥有信徒
func (this *StarSession) hasBeliever() bool {
	return this.believerCount > 0
}

//获取信徒
func (this *StarSession) getBeliever(believerID string) *db.DBBeliever {
	for _, believer := range this.Believer {
		if believer.ID == believerID {
			return believer
		}
	}
	return nil
}

//获取指定等级的信徒
func (this *StarSession) getBelieverByLevel(believerLevel int32) []*db.DBBeliever {
	results := []*db.DBBeliever{}
	for _, believer := range this.Believer {
		if conf.GetGameObjectLevel(believer.ID) == believerLevel {
			results = append(results, believer)
		}
	}
	return results
}

//自动加信徒
func (this *StarSession) AutoAddBeliever(times int32) []*protocol.BelieverInfo {
	believerUpdateInterval := float64(this.Buff[constant.BUFF_BELIEVER]) + conf.Base.BelieverBuffInterval
	var results []*protocol.BelieverInfo
	//believerReduceTime := float64(times*conf.Base.EggTouchReduceTime)
	//if  believerReduceTime > believerUpdateInterval {
	//	exception.GameException(exception.MANUAL_ADD_BELIEVER_ERROR)
	//}
	//believerUpdateInterval -= believerReduceTime
	//believerUpdateInterval := conf.Base.BelieverBuffInterval
	if !this.IsFlagUnlock(constant.STAR_FLAG_EGG) {
		//蛋没解锁
		this.BelieverUpdateTime = time.Now()
		return results
	}
	interval := time.Now().Sub(this.BelieverUpdateTime).Seconds()
	ratio := int(interval / believerUpdateInterval + 0.2) //获取时间间隔内的信徒个数(2秒误差)
	//不能超过上限
	ratio1 := int(this.getAddBelieverCount())
	limit := false
	if ratio1 <= ratio {
		ratio = ratio1
		limit = true
	}
	if ratio <= 0 {
		return results
	}
	scope := conf.GetStarInitBeliever(this.Type)
	if scope == nil {
		return results
	}
	scopeLen := len(scope)
	if scopeLen == 0 {
		return results
	}

	addBelievers := make(map[string]int)
	if ratio == 1 {
		result := this.getMinBeliever(scope)
		addBelievers[result] = addBelievers[result] + 1
	} else {
		scopeIndex := 0
		for i := 0; i < ratio; i++ {
			result := scope[scopeIndex]
			addBelievers[result] = addBelievers[result] + 1
			scopeIndex++
			if scopeIndex == scopeLen {
				scopeIndex = 0
			}
		}
	}

	for id, num := range addBelievers {
		believer := this.addBeliever(id, int32(num))
		results = append(results, util.BuildBelieverInfo(believer))
	}

	if limit {
		this.BelieverUpdateTime = time.Now()
	} else {
		addDuration := time.Duration(int64(ratio)*int64(believerUpdateInterval)) * time.Second
		this.BelieverUpdateTime = this.BelieverUpdateTime.Add(addDuration)
	}
	this.setDirty()
	return results
}

//升级信徒
func (this *StarSession) UpgradeBeliever(selectID string, matchID string, faith int32) (*db.DBBeliever, int32, *protocol.CivilizationInfo) {
	bases := conf.GetBelieverUpgradeResult(selectID, matchID)

	if bases == nil || len(bases) == 0 {
		exception.GameException(exception.STAR_BELIEVER_UPGRADE_INVALID)
	}

	//results := base.UpgradeResult

	var base *conf.BelieverUpgradeBase = nil
	for _, configBase := range bases {
		if this.CivilizationLevel >= configBase.RequireCivilizationLevel[0] && this.BuildingExMaxLevel >= configBase.RequireBuildingLevel[0]  {
			base = configBase
		}
	}

	if base == nil {
		if this.BuildingExMaxLevel < bases[0].RequireBuildingLevel[0] {
			exception.GameException(exception.STAR_BUILDING_LEVEL_NOT_ENOUGH)
		}
		if this.CivilizationLevel < bases[0].RequireCivilizationLevel[0] {
			exception.GameException(exception.CIVILIZATION_LEVEL_NOT_ENOUGH)
		}
	}


	var newBeliever *db.DBBeliever = nil

	//if results == nil || len(results) == 0 {
	//	exception.GameException(exception.STAR_BELIEVER_UPGRADE_INVALID)
	//}
	cost := base.Cost
	cost = int32(float32(cost) * float32(1 + this.Buff[constant.BUFF_COST]))
	if cost > faith {
		exception.GameException(exception.FAITH_NOT_ENOUGH)
	}

	if selectID == matchID {
		this.decBeliever(selectID, 2)
	} else {
		this.ensureBeliever(selectID, 1)
		this.ensureBeliever(matchID, 1)
		this.decBeliever(selectID, 1)
		this.decBeliever(matchID, 1)
	}


	var civilizationInfo *protocol.CivilizationInfo = nil



	randomResultsTemp := base.GetRandomResult()
	randomResults := make(map[int32]timeutil.WeightData)

	//等级限制
	for weight, require := range base.RandomRequire {
		if this.BuildingExMaxLevel >= require.RequireBuildingLevel && this.CivilizationLevel >= require.RequireCivilLevel {
			randomResults[weight] = randomResultsTemp[weight]
		}
	}
	addProb := this.Buff[constant.BUFF_BELIEVER_CT]
	randomResults = base.RandomDataAddProb(randomResults, addProb)

	randomData := timeutil.RandomWeightData(randomResults)


	if randomData != nil {
		data := randomData.(*conf.UpgradeResult)
		if data.Num > 0 && len(data.UpgradeID) != 0 {
			result := this.getMinBeliever(data.UpgradeID)
			newBeliever = this.addBeliever(result, data.Num)

			level := conf.GetGameObjectLevel(newBeliever.ID)
			if level > 0 {
				ratio := conf.GetUpgradeBelieverRatio(this.CivilizationLevel, level)
				rewardCivil := int32(float64(data.CivilizationIncome) * ratio)
				if rewardCivil > 0 {
					civilizationInfo = this.TakeInCivilization(rewardCivil)
				}
			}

			//消耗两个,合成一个新的信徒
			this.AddStatisticsValue(constant.STAR_STATISTIC_UPGRADE_BELIEVER, float64(data.Num), 0)
			this.setDirty()
		}
	}
	return newBeliever, cost, civilizationInfo
}

func resolveBelieversToFaith(believers map[string]int32) (int32, int32) {
	var faith int32 = 0
	var resolveNum int32 = 0
	for id, num := range believers {
		faith += resolveBelieverToFaith(id, num)
		resolveNum += num
	}
	return faith, resolveNum
}

func resolveBelieverToFaith(believerID string, num int32) int32 {
	return conf.Base.BelieverFaithMapping[conf.GetBelieverLevel(believerID)] * num
}

//根据人数上限部分转换为信仰
func (this *StarSession) addBelieverLimit(updateBelievers map[string]int32) (map[string]int32, int32, int32) {
	results := make(map[string]int32)
	var faith int32 = 0
	var resolveNum int32 = 0
	for believerID, believerNum := range updateBelievers {
		if this.getAddBelieverCount() >= believerNum {
			this.addLevelBeliever(believerID, believerNum)
			results[believerID] += believerNum
		} else {
			ableAddNum := this.getAddBelieverCount()
			if ableAddNum > 0 {
				this.addLevelBeliever(believerID, ableAddNum)
				results[believerID] += ableAddNum
			}
			resolveNum = believerNum - ableAddNum
			faith += resolveBelieverToFaith(believerID, believerNum - ableAddNum)
		}
	}
	return results, faith, resolveNum
}

func (this *StarSession) ConvertBeliever(updateInfo []*protocol.BelieverInfo,) ([]*db.DBBeliever, int32){
	changeBeliever := []*db.DBBeliever{}
	var faith int32 = 0
	var resolveNum int32 = 0
	var faithResolve int32 = 0
	if updateInfo == nil || len(updateInfo) == 0 {
		return changeBeliever, faith
	}
	updateBelievers, faithConvert := this.convertBelieverCivil(updateInfo)

	if this.getAddBelieverCount() <= 0 {
		//全部转换为信仰
		faithResolve, resolveNum = resolveBelieversToFaith(updateBelievers)
	} else {
		var results map[string]int32
		results, faithResolve, resolveNum = this.addBelieverLimit(updateBelievers)
		for believerID, believerNum := range results {
			changeBeliever = append(changeBeliever, &db.DBBeliever{ID: believerID, Num: believerNum})
		}
	}
	if resolveNum > 0 {
		this.AddStatisticsValue(constant.STAR_STATISTIC_RESOLVE_BELIEVER, float64(resolveNum), 0)
	}
	this.setDirty()
	faith = faithConvert + faithResolve
	return changeBeliever, faith
}

func (this *StarSession) UpdateBeliever(operation int32, updateInfo []*protocol.BelieverInfo, isConvert bool) ([]*db.DBBeliever, int32) {
	if isConvert {
		return this.ConvertBeliever(updateInfo)
	}
	changeBeliever := []*db.DBBeliever{}
	for _, believerInfo := range updateInfo {
		if (operation == constant.OP_BELIEVER_ADD) {
			believer := this.addLevelBeliever(believerInfo.GetId(), believerInfo.GetNum())
			if (believer != nil) {
				changeBeliever = append(changeBeliever, believer)
			}
		} else if (operation == constant.OP_BELIEVER_DEC) {
			believer := this.decLevelBeliever(believerInfo.GetId(), believerInfo.GetNum())
			if (believer != nil) {
				changeBeliever = append(changeBeliever, believer)
			}
		}
	}
	this.setDirty()
	return changeBeliever, 0
}

func (this *StarSession) updateBelieverCount(updateSearch bool) int32 {
	var count int32 = 0
	var totalLevel int32 = 0
	for _, believer := range this.Believer {
		count += believer.Num
		level := conf.GetGameObjectLevel(believer.ID)
		if level > 0 && believer.Num > 0 {
			totalLevel += level * believer.Num
		}
	}

	this.believerCount = count
	this.believerTotalLevel = totalLevel

	//星球没有被禁用才需要添加到搜索中
	if updateSearch && this.Owner > 0 {
		rpc.SearchServiceProxy.UpdateData(constant.SEARCH_OPT_UPDATE_BELIEVER, this.ID, this.believerTotalLevel)
	}
	cache.StarCache.SetBelieverCount(this.ID, this.believerCount)
	cache.StarCache.SetBelieverTotalLevel(this.ID, this.believerTotalLevel)
	return this.believerCount
}

//掠夺信徒
func (this *StarSession) LootBeliever(believerID []string) ([]string, bool, bool) {

	if this.IsBuildingAllMaxLevel() || !this.IsFlagUnlock(constant.STAR_FLAG_MUTUAL) {
		return nil, false, true
	}
	result := []string{}
	if this.DecBelieverShield() {
		return nil, true, false
	}
	for _, id := range believerID {
		believer := this.forceDecBeliever(id, 1)
		if believer != nil {
			result = append(result, believer.ID)
			this.AddStatisticsValue(constant.STAR_STATISTIC_BE_LOOT_BELIEVER_NUM, 1, 0)
		}
	}
	return result, false, false
}

func(this *StarSession) GetStarBelieversInfo() []*db.DBBeliever {
	return this.Believer
}