package db

import (
	"aliens/common/character"
	"time"
	"gok/service/msg/protocol"
	"gok/service/exception"
	"gok/constant"
	"gok/module/star/conf"
)

//建筑是否损坏
func (this *DBBuilding) IsBroken() bool {
	return this.BrokenTime != 0
}

//建筑损坏过期
func (this *DBBuilding) IsBrokenOverdue() bool {
	return this.BrokenTime != 0 && this.BrokenTime <= time.Now().Unix()
}

//还差对应的时间会不会损毁
//func (this *DBBuilding) IsBrokenRemainOverdue(remain int64) bool {
//	return this.BrokenTime != 0 && (this.BrokenTime - remain) <= time.Now().Unix()
//}

//还差对应的时间会不会损毁
func (this *DBBuilding) IsBrokenRemainOverdue(remain int64) bool {
	return this.BrokenTime != 0 && (this.BrokenTime - remain) <= time.Now().Unix()
}

//机器人帮助维修到
//func (this *DBBuilding) IsRobotHelpOverdue() bool {
//	return time.Now().After(this.RobotHelpTime)
//}

//func (this *DBBuilding) CanRobotHelpRepair() bool {
//	//leftTime := this.BrokenTime - time.Now().Unix()
//	//buildingConf := conf.GetBuildingConf(this.GetStarType(), this.Type, this.Level)
//	//if buildingConf == nil {
//	//	exception.GameException(exception.STAR_BUILDING_NOTFOUND)
//	//}
//	//if this.HelpTime.After(this.HelpTime.Add(time.Duration()))
//
//}

func (this *History) BuildProtocol() *protocol.History {
	return &protocol.History{
		Id:     this.ID,
		Param1: this.Param1,
		Param2: this.Param2,
		Param3: this.Param3,
		Time:   this.Time.Unix(),
	}
}



//重置建筑
func (this *DBBuilding) Reset() {
	this.Level = 0
	//this.UpgradeBelieverCost = 0
	//this.RepairBelieverCost = 0
	this.BrokenTime = 0
	this.Faith = 0
	resetTime := time.Time{}
	this.CreateTime = resetTime
	this.RepairTime = resetTime
	this.UpdateTime = resetTime
	//this.ItemGrooves = []*DBItemGroove{}
}

//获取建筑总放入的物品
//func (this *DBBuilding) GetItems() []int32 {
//	itemIDs := []int32{}
//	for _, itemGroove := range this.ItemGrooves {
//		if (itemGroove.ItemID != 0) {
//			itemIDs = append(itemIDs, itemGroove.ItemID)
//		}
//	}
//	return itemIDs
//}

//func (this *DBBuilding) UpdateRepair() (bool, int64) {
//	repairTime := this.RepairTime.Unix()
//	if repairTime <= 0 {
//		return false, repairTime
//	}
//	if time.Now().After(this.RepairTime) {
//		this.RepairTime = time.Now()
//		//this.RepairBelieverCost =
//		return true, 0
//	} else {
//		return false, repairTime
//	}
//}

//获取生效的圣物数量
//func (this *DBBuilding) GetEffectItemTotal() int {
//	total := 0
//	for _, itemGroove := range this.ItemGrooves {
//		if itemGroove.ItemID != 0 && itemGroove.Effect {
//			total ++
//		}
//	}
//	return total
//}
//
//func (this *DBBuilding) GetItemTotal() int {
//	total := 0
//	for _, itemGroove := range this.ItemGrooves {
//		if (itemGroove.ItemID != 0) {
//			total ++
//		}
//	}
//	return total
//}

func (this *DBBuilding) TakeoutFaith(faith int32) {
	if (this.Faith >= faith) {
		this.Faith -= faith
	}
}

//设置建筑损坏
func (this *DBBuilding) SetBroken(buff map[int32]float32, force bool) bool {
	if this.IsBroken() && !force{ //建筑物还在损坏状态
		return false
		//exception.GameException(exception.STAR_ATK_BUILD_FAILED)
	}
	//设置损坏维护时间
	conf := conf.GetBuildingConf(this.GetStarType(), this.Type, this.Level)
	durationTime := constant.BUILDING_BROKEN_TIME
	if conf != nil {
		durationTime = int(conf.BoomTime)
	}

	brokenTime := time.Now().Add(time.Duration(durationTime) * time.Second)
	this.BrokenTime = brokenTime.Unix() //升级信仰值时间改成升级后

	this.CalFaithUpdateTime(brokenTime, buff)

	if !this.IsUpgrading() || this.RemainTime > 0 {
		return true
	}
	remainTime := this.UpdateTime.Unix() - time.Now().Unix()
	if remainTime > 0 {
		this.RemainTime = remainTime
	}
	return true
}

//是否正在升级中
func (this *DBBuilding) IsUpgrading() bool {
	return this.UpdateTime.Unix() > 0
}

func (this *DBBuilding) IsUpgraded(now time.Time) bool {
	return now.After(this.UpdateTime)
}

func (this *DBBuilding) CleanUpgradeTime() {
	this.UpdateTime = time.Time{}
}

func (this *DBBuilding) CleanRepairTime() {
	this.RepairTime = time.Time{}
}

func (this *DBBuilding) UpdateUpgrading(upgrade bool, buff map[int32]float32) (bool, bool) {
	if !this.IsUpgrading() { //表示正在升级中
		return false, false
	}
	if !upgrade {
		this.CalFaithUpdateTime(this.UpdateTime, buff)
		return false, false
	}
	if time.Now().After(this.UpdateTime) { //更新时间到了
		this.CheckoutFaithBuff1(this.GetStarType(), false, this.UpdateTime, buff) //信仰结算需要结算到升级成功前
		//this.CheckoutFaithBuff1(this.GetStarType(), false, time.Now())
		first := this.UpgradeLevel()
		return true, first
	}
	return false, false
}

//func (building *DBBuilding) BuildGroovesProtocol() []*protocol.ItemGroove {
//	grooves := []*protocol.ItemGroove{}
//	for _, groove := range building.ItemGrooves {
//		grooves = append(grooves, &protocol.ItemGroove{
//			//Color:  groove.Color),
//			ItemID:          groove.ItemID,
//			Effect:          groove.Effect,
//			EffectTimestamp: groove.EffectTime.Unix(),
//		})
//	}
//	return grooves
//}

func (this *DBStarItemGroup) BuildProtocol() *protocol.ItemGroup {
	records := []*protocol.ItemGroupRecord{}
	for _, record := range this.Records {
		records = append(records, &protocol.ItemGroupRecord{ItemID:record.Items, Num:record.Num})
	}

	result := &protocol.ItemGroup{
		GroupID : this.ID,
		//ItemID : this.Items,
		//GetReward: this.Reward,
		//Time: this.UpdateTime.Unix(),
		Done: this.Done,
		ItemGroupRecord:records,
		Active:this.Active,
		//ItemGroupRecord:this.Records,
	}

	return result
}

func (building *DBBuilding) BuildProtocol() *protocol.BuildingInfo {
	levelInfo := []*protocol.LevelInfo{}
	for _, info := range building.LevelInfo {
		levelInfo = append(levelInfo, info.BuildProtocol())
	}

	buildingInfo := &protocol.BuildingInfo{
		Id:         building.ID,
		Type:       building.Type,
		Level:      building.Level,
		RepairTime: building.RepairTime.Unix(),
		UpdateTime: building.UpdateTime.Unix(),
		BrokenTime: building.BrokenTime, //待维修时间
		//RepairBeliever:  building.RepairBelieverCost),
		//UpgradeBeliever: building.UpgradeBelieverCost),
		//ItemGrooves: building.BuildGroovesProtocol(),
		Exist:       building.Exist,
		LevelInfo:   levelInfo,
	}
	//if building.Buff != nil {
	//	buffIDArray := []int32{}
	//	buffNumArray := []int32{}
	//	for buffID, buffNum := range building.Buff {
	//		buffIDArray = append(buffIDArray, buffID)
	//		buffNumArray = append(buffNumArray, buffNum)
	//	}
	//	buildingInfo.BuffID = buffIDArray
	//	buildingInfo.BuffNum = buffNumArray
	//}
	return buildingInfo
}

func (this *DBStarFlag) BuildProtocol() *protocol.FlagInfo {
	return &protocol.FlagInfo{Id:int32(this.Flag), Value:this.Value, Time:this.UpdateTime.Unix()}
}

//func (this *DBItemGroove) Reset(filterColor []int32) {
//	randomColor := rand.Int31n(constant.MAX_RANDOM_COLOR) + 1
//	this.ItemID = 0
//	if (len(filterColor) == 0) {
//		this.Color = randomColor
//
//	} else {
//		for currColor := randomColor; currColor <= constant.MAX_RANDOM_COLOR; currColor++ {
//			if (!character.ContainsInt32(currColor, filterColor)) {
//				this.Color = currColor
//				return
//			}
//		}
//
//		for currColor := randomColor; currColor > 0; currColor-- {
//			if (!character.ContainsInt32(currColor, filterColor)) {
//				this.Color = currColor
//				return
//			}
//		}
//	}
//
//}

func (this *DBBuilding) GetCurrLevelInfo() *LevelInfo {
	return this.GetLevelInfo(this.Level)
}

func (this *DBBuilding) GetLevelInfo(level int32) *LevelInfo {
	for _, info := range this.LevelInfo {
		if info.ID == level {
			return info
		}
	}
	result := &LevelInfo{ID: this.Level}
	this.LevelInfo = append(this.LevelInfo, result)
	return result
}

//获取历史最高等级
func (this *DBBuilding) GetMaxLevel() int32 {
	var maxLevel int32 = 0
	for _, info := range this.LevelInfo {
		currLevel := info.ID
		if info.Time > 0 {
			currLevel++
		}

		if currLevel > maxLevel{
			maxLevel = currLevel
		}
	}
	return maxLevel
}

func (this *LevelInfo) BuildProtocol() *protocol.LevelInfo {
	believerCost := []*protocol.BelieverCost{}
	for _, info := range this.BelieverCost {
		believerCost = append(believerCost, &protocol.BelieverCost{BelieverID: info.ID, Num: info.Num})
	}
	return &protocol.LevelInfo{
		Level:        this.ID,
		Time:         this.Time,
		FaithCost:    this.FaithCost,
		BelieverCost: believerCost,
	}
}

func (this *LevelInfo) GetBelieverCost(believerID string) *BelieverCost {
	for _, info := range this.BelieverCost {
		if info.ID == believerID {
			return info
		}
	}
	return nil
}

func (this *LevelInfo) AddBelieverCosts(believerIDs []string) {
	for _, believerID := range believerIDs {
		this.AddBelieverCost(believerID, 1)
	}
}

func (this *LevelInfo) AddBelieverCost(believerID string, num int32) {
	cost := this.GetBelieverCost(believerID)
	if cost == nil {
		cost = &BelieverCost{ID: believerID, Num: num}
		this.BelieverCost = append(this.BelieverCost, cost)
	}
	cost.Num += num
}

//升级 return 是否第一次升级
func (this *DBBuilding) UpgradeLevel() bool {
	this.UpdateTime = time.Time{}

	levelInfo := this.GetLevelInfo(this.Level)

	this.Level += 1
	if !this.Exist {
		this.Exist = true
	}
	if levelInfo.Time <= 0 {
		levelInfo.Time = time.Now().Unix()
		return true
	}
	return false

	//if len(this.ItemGrooves) >= int(this.Level) {
	//	return
	//}
	//newGroove := &DBItemGroove{}
	////newGroove.Reset(this.GetFilterColor())
	//this.ItemGrooves = append(this.ItemGrooves, newGroove)
	//if (this.Level == 1) {
	//	newGroove := &DBItemGroove{}
	//	newGroove.Reset()
	//	this.ItemGrooves = append(this.ItemGrooves, newGroove)
	//}
	//	this.UpgradeBelieverCost = 0
}

//获取需要过滤的颜色
//func (this *DBBuilding) GetFilterColor() []int32 {
//	colorCount := make(map[int32]int)
//	for _, groove := range this.ItemGrooves {
//		colorCount[groove.Color] = colorCount[groove.Color] + 1
//	}
//
//	filterColor := []int32{}
//	for color, count := range colorCount {
//		if (count == 2) {
//			filterColor = append(filterColor, color)
//		}
//	}
//	return filterColor
//}

//升级更新升级时间
func (this *DBBuilding) UpdateUpgradeTime(upgradeTime int32, buff map[int32]float32) {
	this.UpdateTime = time.Now().Add(time.Duration(upgradeTime) * time.Second)
	//升级信仰值时间改成升级后
	this.CalFaithUpdateTime(this.UpdateTime, buff)
	if this.Level == 0 {
		this.CreateTime = time.Now()
		this.Faith = 0
	}
}

//加速减少升级时间
func (this *DBBuilding) AccUpgradeTime(accTime int32, upgrade bool, buff map[int32]float32) (bool, bool) {
	////获取相差时间
	this.UpdateTime = time.Unix(this.UpdateTime.Unix()-int64(accTime), 0)
	this.CalFaithUpdateTime(this.UpdateTime, buff)

	if !upgrade {
		return false, false
	}
	//升级时间到了
	if time.Now().After(this.UpdateTime) {
		first := this.UpgradeLevel()
		return true, first
	}
	return false, false
}

//是否在修理中
func (this *DBBuilding) IsRepairing() bool {
	return this.RepairTime.Unix() > 0
}

func (this *DBBuilding) IsRepaired(now time.Time) bool {
	return now.After(this.RepairTime)
}

//更新维修状态
func (building *DBBuilding) UpdateRepairing(repaired bool) bool {
	if !building.IsRepairing() {
		return false
	}
	if !repaired {
		return false
	}
	if time.Now().After(building.RepairTime) {
		building.RepairTime = time.Time{}
		return true
	} else {
		return false
	}
}

//开始修理更新修理时间
func (this *DBBuilding) UpdateRepairTime(repairTime int32, buff map[int32]float32) {
	this.RepairTime = time.Now().Add(time.Duration(repairTime) * time.Second) //设置修理时间
	this.CalFaithUpdateTime(this.RepairTime, buff)
	this.BrokenTime = 0
	//升级中修理需要更新升级时间
	if this.IsUpgrading() && this.RemainTime != 0 {
		this.UpdateTime = this.RepairTime.Add(time.Duration(this.RemainTime))
		this.RemainTime = 0
	}
}

//加速减少修理时间
func (this *DBBuilding) AccRepairTime(accTime int32, repaired bool, buff map[int32]float32) bool {
	this.RepairTime = time.Unix(this.RepairTime.Unix()-int64(accTime), 0)
	this.CalFaithUpdateTime(this.RepairTime, buff)
	////升级中修理需要更新升级时间
	//if this.IsUpgrading() {
	//	this.UpdateTime = time.Unix(this.UpdateTime.Unix()-int64(accTime), 0)
	//}
	//return this.UpdateRepairing(), this.RepairTime.Unix()
	//	if time.Now().After(this.RepairTime) {
	//		this.RepairTime = time.Time{}
	////		this.RepairBelieverCost = 0
	//		return true, 0
	//	}
	//	return false, this.RepairTime.Unix()


	if !repaired {
		return false
	}
	//升级时间到了
	if time.Now().After(this.RepairTime) {
		this.RepairTime = time.Time{}
		return true
	}
	return false
}

//获取建筑所属的星球
func (this *DBBuilding) GetStarType() int32 {
	idStr := character.Int32ToString(this.ID)
	return character.StringToInt32(idStr[0:len(idStr)-1])
}

//结算并更新新研制更新时间戳
func (this *DBBuilding) CalFaithUpdateTime(newTime time.Time, buff map[int32]float32) {
	this.CheckoutFaithBuff(this.GetStarType(), false, buff)
	this.FaithUpdateTime = newTime
}

//func (this *DBBuilding) ActiveGrooveEffect(grooveID int32) (int64, bool) {
//	groove := this.GetGroove(grooveID)
//	if groove == nil {
//		exception.GameException(exception.GROOVE_NOT_FOUND)
//	}
//	if time.Now().After(groove.EffectTime) {
//		groove.Effect = true
//		this.UpdateBuff()
//	}
//	return groove.EffectTime.Unix(), groove.Effect
//}
//
//func (this *DBBuilding) AccGrooveEffect(grooveID int32, num int) (int64, bool) {
//	groove := this.GetGroove(grooveID)
//	if groove == nil {
//		exception.GameException(exception.GROOVE_NOT_FOUND)
//	}
//	if groove.Effect {
//		exception.GameException(exception.GROOVE_ALREADY_EFFECT)
//	}
//	groove.EffectTime = groove.EffectTime.Add(-util.GetDuraton(float64(num) * conf.Base.BelieverAddEffectTime))
//	if time.Now().After(groove.EffectTime) {
//		groove.Effect = true
//		this.UpdateBuff()
//	}
//	return groove.EffectTime.Unix(), groove.Effect
//}

//更新建筑BUFF
//func (this *DBBuilding) UpdateBuff() {
//	this.Buff = make(map[int32]int32)
//	groups := make(map[int32][]int32)
//
//	for _, item := range this.ItemGrooves {
//		if !item.Effect || item.ItemID == 0 {
//			continue
//		}
//		buffID := conf.GetItemBuff(item.ItemID)
//		if buffID != 0 {
//			this.Buff[buffID] += 1
//		}
//
//		groupID := conf.Base.StarGroupMapping[item.ItemID]
//		if groupID != 0 {
//			groupItems := groups[groupID]
//			if groupItems == nil {
//				groupItems = []int32{}
//			}
//			groupItems = append(groupItems, item.ItemID)
//			groups[groupID] = groupItems
//		}
//	}
//
//	for groupID, groupItems := range groups {
//		itemLen := len(groupItems)
//		if itemLen < 3 {
//			continue
//		}
//		groupBase := conf.Base.ItemGroup[groupID]
//		if groupBase == nil {
//			continue
//		}
//		if groupBase.BuffID == nil || len(groupBase.BuffID) == 0 {
//			continue
//		}
//		this.Buff[groupBase.BuffID[0]] += 1
//
//		//全部完成追加BUFF
//		if itemLen == len(groupBase.Content) && len(groupBase.BuffID) >= 2 {
//			this.Buff[groupBase.BuffID[1]] += 1
//		}
//		//if !groupBase.IsFinish(groupItems) {
//		//	continue
//		//}
//
//	}
//
//	//通知buff改变
//	if this.BuffChangeListener != nil {
//		this.BuffChangeListener(this.Type, this.Buff)
//	}
//}

//func (this *DBBuilding) GetBuffRatio(getBuffType int32) float32 {
//	//建筑不完好没有buff
//	//if !this.IsIntact() {
//	//	return 0
//	//}
//	if this.Buff == nil {
//		this.Buff = make(map[int32]int32)
//	}
//	for buffID, buffNum := range this.Buff {
//		buffType, buffRatio := conf.GetBuffBase(buffID)
//		if getBuffType == buffType {
//			return float32(buffNum) * buffRatio
//		}
//	}
//	return 0
//}

//结算所有BUFF
func (this *DBBuilding) CheckoutBuff(starType int32, buff map[int32]float32) {
	this.CheckoutFaithBuff(starType, false, buff)
}

func (this *DBBuilding) IsIntact() bool {
	return !this.IsBroken() && !this.IsRepairing() && !this.IsUpgrading()
}

func (this *DBBuilding) CheckoutFaithBuff(starType int32, draw bool, buff map[int32]float32) int32 {
	return this.CheckoutFaithBuff1(starType, draw, time.Now(), buff)
}

//BUFF变更的时候信仰需要结算
func (this *DBBuilding) CheckoutFaithBuff1(starType int32, draw bool, updateTime time.Time, buff map[int32]float32) int32 {
	if this.Level == 0 {
		return 0
	}
	//建筑不完好不能领取信仰
	if !this.IsIntact() && draw {
		exception.GameException(exception.STAR_BUILDING_FAITH_CANNOTDRAW)
	}
	buildingConf := conf.GetBuildingConf(starType, this.Type, this.Level)
	if buildingConf == nil {
		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
	}
	duration := int32(updateTime.Sub(this.FaithUpdateTime).Seconds())
	interval := buildingConf.UpdateFaithTime
	if duration > interval {
		//圣物加成
		//faithRatio := 1 + conf.Base.RelicFaithPowerPercent * float64(this.GetEffectItemTotal())
		faithRatio := float64(1 + buff[constant.BUFF_FAITH])
		addFaith := int32(float64((duration/interval)*buildingConf.UpdateFaithNum) * faithRatio)
		//addFaith := int32(float64((duration/interval)*buildingConf.UpdateFaithNum))
		//addFaith := buildingConf.UpdateFaithNum
		remainTime := duration % interval
		this.Faith += addFaith
		this.FaithUpdateTime = time.Unix(updateTime.Unix()-int64(remainTime), 0)
	}
	faithLimitRatio := float64(1 + buff[constant.BUFF_FAITH_LIMIT])
	limit := int32(float64(buildingConf.FaithLimit) * faithLimitRatio)
	//limit := buildingConf.FaithLimit
	if this.Faith > limit {
		this.Faith = limit
	}
	result := this.Faith
	if draw {
		if this.Faith < buildingConf.ReceiveFaithMin {
			exception.GameException(exception.STAR_BUILDING_FAITH_CANNOTDRAW)
		}
		this.Faith = 0
	}
	return result
}

//func (this *DBBuilding) GetGroove(grooveID int32) *DBItemGroove {
//	if (this.ItemGrooves == nil) {
//		return nil
//	}
//	if (len(this.ItemGrooves) < int(grooveID)) {
//		return nil
//	}
//	return this.ItemGrooves[grooveID]
//}
//
////取出存在的物品
//func (this *DBBuilding) TakeoutExistItemID(itemIDs []int32) []int32 {
//	exists := []int32{}
//	for _, groove := range this.ItemGrooves {
//		if groove.ItemID == 0 {
//			continue
//		}
//		if character.ContainsInt32(groove.ItemID, itemIDs) {
//			groove.ItemID = 0
//			exists = append(exists, groove.ItemID)
//		}
//	}
//	this.UpdateBuff()
//	return exists
//}
//
//func (this *DBBuilding) TakeoutExistItem(itemIDs []int32) map[int32]*DBItemGroove {
//	exists := make(map[int32]*DBItemGroove)
//	for _, groove := range this.ItemGrooves {
//		if groove.ItemID == 0 {
//			continue
//		}
//		if character.ContainsInt32(groove.ItemID, itemIDs) {
//			exists[groove.ItemID] = &DBItemGroove{EffectTime: groove.EffectTime, Effect: groove.Effect}
//			groove.ItemID = 0
//		}
//	}
//	if len(exists) > 0 {
//		this.UpdateBuff()
//	}
//	return exists
//}

//func (this *DBBuilding) TakeinItems(grooveIDs []int32, itemIDs []int32, persistItems map[int32]*DBItemGroove) []int32 {
//	if len(grooveIDs) != len(itemIDs) {
//		exception.GameException(exception.INVALID_PARAM)
//	}
//	result := []int32{}
//	for index, grooveID := range grooveIDs {
//		takeoutItem := this.ReplaceItem(grooveID, itemIDs[index], persistItems)
//		if takeoutItem != 0 {
//			result = append(result, takeoutItem)
//		}
//	}
//	this.UpdateBuff()
//	return result
//}
//
//func (this *DBBuilding) ReplaceItem(grooveID int32, itemID int32, persist map[int32]*DBItemGroove) int32 {
//	groove := this.GetGroove(grooveID)
//	if groove == nil {
//		exception.GameException(exception.GROOVE_NOT_FOUND)
//	}
//
//	if itemID == groove.ItemID {
//		return itemID
//	}
//
//	var persistGroove *DBItemGroove = nil
//
//	if persist != nil {
//		//需要保留的状态
//		persistGroove = persist[itemID]
//	}
//
//	if persistGroove != nil {
//		groove.Effect = persistGroove.Effect
//		groove.EffectTime = persistGroove.EffectTime
//	} else {
//		cdTime := time.Duration(conf.Base.GrooveEffectTime) * time.Second
//		groove.EffectTime = time.Now().Add(cdTime)
//		groove.Effect = false
//	}
//	takeoutItem := groove.ItemID
//	groove.ItemID = itemID
//	return takeoutItem
//}
//
//func (building *DBBuilding) TakeoutRandomItem() int32 {
//	grooves := []*DBItemGroove{}
//	for _, groove := range building.ItemGrooves {
//		if (groove.ItemID != 0) {
//			grooves = append(grooves, groove)
//		}
//	}
//	grooveLen := len(grooves)
//	if (grooveLen == 0) {
//		return 0
//	}
//	groove := grooves[rand.Intn(grooveLen)]
//	itemID := groove.ItemID
//	groove.ItemID = 0
//	building.UpdateBuff()
//	return itemID
//}
//
//func (this *DBBuilding) TakeoutItem(grooveID int32) int32 {
//	groove := this.GetGroove(grooveID)
//	if (groove == nil) {
//		exception.GameException(exception.GROOVE_NOT_FOUND)
//	}
//	if (groove.ItemID == 0) {
//		exception.GameException(exception.ITEM_NOT_FOUND)
//	}
//	removeItemID := groove.ItemID
//	groove.ItemID = 0
//	this.UpdateBuff()
//	return removeItemID
//}

//重置槽
//func (this *DBBuilding) ResetGrooves(lockGroove []int32) []int32 {
//
//	returnItems := []int32{}
//	for index, groove := range this.ItemGrooves {
//		if (character.ContainsInt32(int32(index), lockGroove)) {
//			continue
//		}
//		if (groove.ItemID != 0) {
//			returnItems = append(returnItems, groove.ItemID)
//		}
//		groove.Reset(this.GetFilterColor())
//	}
//	return returnItems
//}

//重置槽
//func (this *DBBuilding) ActiveGrooves() []int32 {
//	returnItems := []int32{}
//	for _, groove := range this.ItemGrooves {
//		if (groove.ItemID != 0) {
//			if (groove.Effect) {
//				returnItems = append(returnItems, groove.ItemID)
//			}
//		}
//	}
//	return returnItems
//}

func (this *Statistics) BuildProtocol() *protocol.Statistics {
	return &protocol.Statistics{
		Id:    this.ID,
		Value: this.Value,
	}
}

func (this *Shield) BuildProtocol() *protocol.Shield {
	return &protocol.Shield{
		Value:this.Value,
		Limit:this.Limit,
		UpdateTime:this.UpdateTime.Unix(),
	}
}

func (star *DBStar) BuildProtocol() *protocol.StarInfoDetail{//写入星球详细数据  wjl 20170606
	if star == nil{
		return nil
	}
	buildingInfos := []*protocol.BuildingInfo{}
	for _, building := range star.Building {
		buildingInfos = append(buildingInfos, building.BuildProtocol())
	}
	believerInfos := []*protocol.BelieverInfo{};
	for _, believer := range star.Believer {
		believerInfos = append(believerInfos, believer.BuildProtocol())
	}
	civilizationInfos := []*protocol.CivilizationReward{};
	for _, reward := range star.CivilizationReward {
		civilizationInfos = append(civilizationInfos, &protocol.CivilizationReward{Level:reward.Level, Draw:reward.Draw})
	}
	message := &protocol.StarInfoDetail{
		StarID:  star.ID,
		Type:  star.Type,
		Seq:   star.Seq,
		OwnID: star.Owner,//拥有者的用户ID
		Building:buildingInfos,
		Believer:believerInfos,
		CreateTime:star.CreateTime.Unix(),
		DoneTime:star.DoneTime.Unix(),
		CivilizationLv:star.CivilizationLevel,
		CivilizationProgress:star.CivilizationValue,
		CivilizationReward:civilizationInfos,
		StarFlags: star.BuildFlagsProtocol(),
	}
	return message
}

func (this *DBBeliever) BuildProtocol() *protocol.BelieverInfo {
	return &protocol.BelieverInfo{
		Id: this.ID,
		Num: this.Num,
	}
}

func (this *DBStar) BuildFlagsProtocol() []*protocol.FlagInfo {
	var result []*protocol.FlagInfo
	for _, flag := range this.Flags {
		//this.RefreshFlag(flag)
		result = append(result, flag.BuildProtocol())
	}
	return result
}

