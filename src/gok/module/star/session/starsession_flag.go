package session

import (
	"gok/constant"
	"gok/module/star/conf"
	"gok/module/star/db"
	"gok/module/star/util"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"gok/service/rpc"
	"time"
)

//-------------------------------STAR_FLAG------------------------------------------------------

//return true:需要更新标识数据 false:不需要更新表示数据
func (this *StarSession) ensureFlag(key int32, value int32) bool {
	data := this.flags[key]
	if data == nil {
		this.flags[key] = &db.DBStarFlag{Flag:key, Value:value, UpdateTime:time.Now()}
		return true
	}
	return false
}

func (this *StarSession) getFlags() []*db.DBStarFlag {
	var result []*db.DBStarFlag
	for _, flag := range this.flags {
		result = append(result, flag)
	}
	return result
}

func (this *StarSession) GetProtocolFlags() []*protocol.FlagInfo {
	var result []*protocol.FlagInfo
	//flags := this.getFlags()
	//for _, flag := range flags {
	//
	//}
	for _, flag := range this.flags {
		result = append(result, flag.BuildProtocol())
	}
	return result
}

//获取角色标识
func (this *StarSession) GetFlag(key int32) *db.DBStarFlag {
	flag := this.flags[key]
	//this.RefreshFlag(flag)
	return flag
}

func (this *StarSession) GetFlagValue(key int32) int32 {
	flag := this.flags[key]
	if flag == nil {
		return DEFAULT_FLAG_VALUE
	}
	//this.RefreshFlag(flag)
	return flag.Value
}

func (this *StarSession) AddFlag(key int32) *db.DBStarFlag{
	flag := this.flags[key]
	if flag == nil {
		flag = &db.DBStarFlag{Flag:key, Value:1, UpdateTime:time.Now()}
		this.flags[key] = flag
	} else {
		flag.Value += 1
		flag.UpdateTime = time.Now()
	}
	this.updateFlag()
	return flag
}

func (this *StarSession) UpdateFlagValue(key int32, value int32) *db.DBStarFlag{
	if this.flags == nil {
		this.flags = make(map[int32]*db.DBStarFlag)
	}
	flag := this.flags[key]
	if flag == nil {
		flag = &db.DBStarFlag{Flag:key, Value:value, UpdateTime:time.Now()}
		this.flags[key] = flag
	} else {
		flag.Value = value
		flag.UpdateTime = time.Now()
	}
	this.updateFlag()
	return flag
}

func (this *StarSession) updateFlag() {
	this.Flags = this.getFlags()
	this.setDirty()
}

func (this *StarSession) initFlag() {
	if this.Owner < 0 {
		//机器人不用初始化flag
		return
	}
	if this.flags == nil {
		this.flags = make(map[int32]*db.DBStarFlag)
	}
	for _, flag := range this.Flags {
		this.flags[flag.Flag] = flag
	}
	if this.isUpdateFlags() {
		this.updateFlag()
	}
	this.checkFirstGroupFlag()
}

func (this *StarSession) checkFirstGroupFlag() {
	_, groupIndex, doneNum := this.GetCurrentActiveGroup()
	if groupIndex == 1 {
		this.updateFirstGroupFlagUnlock()
	}
	if doneNum > 0 {
		this.updateFirstGroupFlagDone()
	}
}

func (this *StarSession) updateFirstGroupFlagUnlock() {
	this.updateFirstGroupFlag(constant.STAR_FLAG_FIRST_GROUP, constant.FLAG_VALUE_GROUP_UNLOCK)
}

func (this *StarSession) updateFirstGroupFlagDone() {
	this.updateFirstGroupFlag(constant.STAR_FLAG_FIRST_GROUP, constant.FLAG_VALUE_GROUP_DONE)
}

func (this *StarSession) updateFirstGroupFlag(key int32, value int32) {
	var flag *db.DBStarFlag
	if this.GetFlagValue(key) != value {
		flag = this.UpdateFlagValue(key, value)
	}
	if flag != nil {
		this.pushUserFlags([]*protocol.FlagInfo{flag.BuildProtocol()})
	}
}

//return 是否需要更新flags的map到db
func (this *StarSession) isUpdateFlags() bool {
	var ret = false
	if conf.Base.FlagKeyLevelMapping == nil {
		exception.GameException(exception.FLAG_UNLOCK_NOT_FOUND)
	}
	for key := range conf.Base.FlagKeyLevelMapping {
		if this.ensureFlag(key, constant.FLAG_VALUE_LOCK) {
			ret = true
		}
	}
	return ret || this.ensureFlag(constant.STAR_FLAG_FIRST_GROUP, 0)
}

func (this *StarSession) checkLockFlags() {
	if this.Owner < 0 {
		return
	}
	if conf.Base.FlagKeyLevelMapping == nil {
		exception.GameException(exception.FLAG_UNLOCK_NOT_FOUND)
	}
	var updateFlags []*protocol.FlagInfo
	for key, requireLevel := range conf.Base.FlagKeyLevelMapping {
		flag := this.unlockFlag(requireLevel, key)
		if flag != nil {
			updateFlags = append(updateFlags, flag.BuildProtocol())
		}
		//log.Info("this.buildingExMaxLevel:%v, requireLevel:%v, flag_key:%v， flag_value:%v", this.buildingExMaxLevel, requireLevel, db.STAR_FLAG(key), this.GetFlagValue(db.STAR_FLAG(key)))
	}
	this.pushUserFlags(updateFlags)
}

func (this *StarSession) unlockFlag(requireLevel int32, key int32) *db.DBStarFlag {
	if this.BuildingExMaxLevel >= requireLevel && this.GetFlagValue(key) != constant.FLAG_VALUE_UNLOCK {
		switch key  {
		case constant.STAR_FLAG_EGG:
			//蛋解锁 以当前时间为信徒上一次自动刷新时间
			currTime := time.Now()
			if conf.Base.EggActivation != 0 && conf.Base.EggActivationBeliever != 0 {
				subTime := conf.Base.EggActivationBeliever*int64(conf.Base.BelieverBuffInterval) - conf.Base.EggActivation
				if subTime > 0 {
					currTime = currTime.Add(-time.Duration(subTime * int64(time.Second)))
				}
			}
			this.BelieverUpdateTime = currTime
		case constant.STAR_FLAG_REVENGE_MSG:
			//复仇解锁 推一条被机器人打的消息
			eventRobots := StarManager.RandomEventRobot(constant.EVENT_ID_ATK_BUILDING, []int32{1})
			if eventRobots == nil || len(eventRobots)== 0 {
				exception.GameException(exception.ROBOT_IS_NOT_EXIST)
			}
			var build *db.DBBuilding
			for _, building := range this.Building {
				if building.Level > 0 {
					build = building
				}
			}
			newsFeedMessage := util.BuildNewsFeedMessage1(eventRobots[0].Id, constant.NEWSFEED_TYPE_GUIDE_BE_ATK_BUILD, 0, 0,  build.Type, nil)
			rpc.UserServiceProxy.PersistCall(this.Owner, newsFeedMessage)
		}
		flag := this.UpdateFlagValue(key, constant.FLAG_VALUE_UNLOCK)
		return flag
	}
	return nil
}

func (this *StarSession)pushUserFlags(updateFlags []*protocol.FlagInfo) {
	if updateFlags != nil || len(updateFlags) > 0 {
		message := util.BuildUpdateFlag(this.Owner, updateFlags)
		rpc.UserServiceProxy.UserHandleMessage(this.Owner, message)
	}
}

func (this *StarSession) IsFlagUnlock(key int32) bool {
	if this.GetFlagValue(key) == constant.FLAG_VALUE_UNLOCK {
		return true
	}
	return false
}