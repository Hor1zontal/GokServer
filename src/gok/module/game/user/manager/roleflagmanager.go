package manager

import (
	"gok/module/game/db"
	"time"
	"gok/constant"
	"gok/service/msg/protocol"
	"gok/module/game/conf"
	"aliens/common/util"
)

const DEFAULT_FLAG_VALUE int32 = 0

//角色标识
type RoleFlagManager struct {
	flags map[db.ROLE_FLAG]*db.DBRoleFlag
	believerFlags map[string]*db.DBBelieverFlag   //信徒标识
	display []int32
}

//初始化
func (this *RoleFlagManager) Init(role *db.DBRole) {
	this.flags = make(map[db.ROLE_FLAG]*db.DBRoleFlag)
	for _, flag := range role.Flags {
		this.flags[db.ROLE_FLAG(flag.Flag)] = flag
	}

	this.believerFlags = make(map[string]*db.DBBelieverFlag)
	for _, believerFlag := range role.BelieverFlags {
		this.believerFlags[believerFlag.ID] = believerFlag
	}

	this.display = role.Display

	this.EnsureFlag(db.FLAG_SEARCH_COUNT, constant.MAX_SEARCH_LIMIT)
	//this.EnsureFlag(db.FLAG_SHARE, conf.DATA.ShareLimit)
	this.EnsureFlag(db.FLAG_SHARE_WECHAT_HELP, conf.DATA.HelpLimit)
	this.EnsureFlag(db.FLAG_SHARE_WECHAT_SHOW, conf.DATA.ShareLimit)
	this.EnsureFlag(db.FLAG_SHARE_WECHAT_SUCC, conf.DATA.GetHelpLimit)
	this.EnsureFlag(db.FLAG_DRAW_INVITE_GIFT, 0)
	//this.EnsureFlag(db.FLAG_WATCH_AD_GET_POWER, 0)
	this.EnsureFlag(db.FLAG_RE_MALL_ITEM_COUNT, 0)
	this.EnsureFlag(db.FLAG_AD_RE_MALL_ITEM_COUNT, conf.DATA.AdRefresh)
}

//更新数据库内存
func (this *RoleFlagManager) Update(role *db.DBRole) {
	role.Flags = this.GetFlags()
	role.BelieverFlags = this.GetBelieverFlags()
	role.Display = this.display
}

func (this *RoleFlagManager) EnsureFlag(key db.ROLE_FLAG, value int32) {
	data := this.flags[key]
	if (data == nil) {
		this.flags[key] = &db.DBRoleFlag{Flag:key, Value:value, UpdateTime:time.Now()}
	}
}

func (this *RoleFlagManager) UpdateDisplay(display int32) {
	if !util.ContainsInt32(display, this.display) {
		this.display = append(this.display, display)
	}
}

func (this *RoleFlagManager) GetDisplay(min int32, max int32) []int32 {
	result := []int32{}
	if max == -1 {
		return this.display
	}
	for _, displayID := range this.display {
		if displayID <= max && displayID>= min {
			result = append(result, displayID)
		}
	}
	return result
}

func (this *RoleFlagManager) GetBelieverFlags() []*db.DBBelieverFlag {
	result := []*db.DBBelieverFlag{}
	for _, flag := range this.believerFlags {
		result = append(result, flag)
	}
	return result
}

//获取所有标识
func (this *RoleFlagManager) GetFlags() []*db.DBRoleFlag {
	result := []*db.DBRoleFlag{}
	for _, flag := range this.flags {
		result = append(result, flag)
	}
	return result
}

func (this *RoleFlagManager) GetFlagProtocol() []*protocol.FlagInfo {
	result := []*protocol.FlagInfo{}
	for _, flag := range this.flags {
		this.RefreshFlag(flag)
		result = append(result, flag.BuildProtocol())
	}
	return result
}

func (this *RoleFlagManager) RefreshFlag(flag *db.DBRoleFlag) {
	if flag == nil {
		return
	}
	if flag.Flag == db.FLAG_SEARCH_COUNT {
		canAddCount := constant.MAX_SEARCH_LIMIT - flag.Value
		if canAddCount > 0 {
			add := int32(time.Now().Sub(flag.UpdateTime).Seconds() / conf.DATA.CountdownSearch)
			flag.Value += add

			if flag.Value > constant.MAX_SEARCH_LIMIT {
				flag.Value = constant.MAX_SEARCH_LIMIT
			}
			flag.UpdateTime = time.Now()
		}
	}

	//if flag.Flag == db.FLAG_SHARE {
	//	refreshTime := util.GetTodayHourTime(5)
	//	//每天凌晨五点刷新
	//	if flag.UpdateTime.Before(refreshTime) && time.Now().After(refreshTime) {
	//		flag.Value = conf.DATA.ShareLimit
	//
	//	}
	//}

	if flag.Flag == db.FLAG_SHARE_WECHAT_SHOW {
		refreshTime := util.GetTodayHourTime(5)
		//每天凌晨五点刷新
		if flag.UpdateTime.Before(refreshTime) && time.Now().After(refreshTime) {
			flag.Value = conf.DATA.ShareLimit
		}
	}

	if flag.Flag == db.FLAG_SHARE_WECHAT_HELP {
		refreshTime := util.GetTodayHourTime(5)
		//每天凌晨五点刷新
		if flag.UpdateTime.Before(refreshTime) && time.Now().After(refreshTime) {
			flag.Value = conf.DATA.HelpLimit
		}
	}

	if flag.Flag == db.FLAG_SHARE_WECHAT_SUCC {
		refreshTime := util.GetTodayHourTime(5)
		//每天凌晨五点刷新
		if flag.UpdateTime.Before(refreshTime) && time.Now().After(refreshTime) {
			flag.Value = conf.DATA.GetHelpLimit
		}
	}

	//if flag.Flag == db.FLAG_WATCH_AD_GET_POWER {
	//	refreshTime := util.GetTodayHourTime(5)
	//	//每天凌晨五点刷新
	//	if flag.UpdateTime.Before(refreshTime) && time.Now().After(refreshTime) {
	//		flag.Value = 0
	//	}
	//}

	if flag.Flag == db.FLAG_RE_MALL_ITEM_COUNT {
		refreshTime := util.GetTodayHourTime(0)
		//每天0点刷新
		if flag.UpdateTime.Before(refreshTime) && time.Now().After(refreshTime) {
			flag.Value = 0
		}
	}

	if flag.Flag == db.FLAG_AD_RE_MALL_ITEM_COUNT {
		refreshTime := util.GetTodayHourTime(0)
		//每天0点刷新
		if flag.UpdateTime.Before(refreshTime) && time.Now().After(refreshTime) {
			flag.Value = conf.DATA.AdRefresh
		}
	}

}

//获取角色标识
func (this *RoleFlagManager) GetFlag(key db.ROLE_FLAG) *db.DBRoleFlag {
	flag := this.flags[key]
	this.RefreshFlag(flag)
	return flag
}

//获取角色标识
func (this *RoleFlagManager) GetFlagValue(key db.ROLE_FLAG) int32 {
	flag := this.flags[key]
	if flag == nil {
		return DEFAULT_FLAG_VALUE
	}
	this.RefreshFlag(flag)
	return flag.Value
}

//是否引导过程中
func (this *RoleFlagManager) IsUpgradeGuiding() bool {
	return this.GetFlagValue(db.FLAG_GUIDE) <= constant.GUIDE_BUILDING_FAITH
}

func (this *RoleFlagManager) UpdateBelieverFlag(id string, value bool) *db.DBBelieverFlag{
	flag := this.believerFlags[id]
	if (flag == nil) {
		flag = &db.DBBelieverFlag{ID:id, Value:value}
		this.believerFlags[id] = flag
	} else {
		flag.Value = value
	}
	flag.UpdateTime = time.Now()
	return flag
}


func (this *RoleFlagManager) GetBelieverFlag(id string) bool {
	flag := this.believerFlags[id]
	if (flag == nil) {
		return false
	}
	return flag.Value
}

func (this *RoleFlagManager) AddLimitFlag(key db.ROLE_FLAG, addValue int32, limit int32) *db.DBRoleFlag {
	flag := this.flags[key]
	if (flag == nil) {
		flag = &db.DBRoleFlag{Flag:key, Value:addValue, UpdateTime:time.Now()}
		this.flags[key] = flag
	} else {
		if (flag.Value == limit) {
			return nil
		}
		flag.Value += addValue
		if (flag.Value > limit) {
			flag.Value = limit
		}
		flag.UpdateTime = time.Now()
	}
	return flag
}

func (this *RoleFlagManager) AddFlag(key db.ROLE_FLAG) *db.DBRoleFlag{
	flag := this.flags[key]
	if flag == nil {
		flag = &db.DBRoleFlag{Flag:key, Value:1, UpdateTime:time.Now()}
		this.flags[key] = flag
	} else {
		flag.Value += 1
		flag.UpdateTime = time.Now()
	}
	return flag
}

//更新角色标识
func (this *RoleFlagManager) UpdateFlag(key db.ROLE_FLAG, value int32) *db.DBRoleFlag{
	flag := this.flags[key]
	if flag == nil {
		flag = &db.DBRoleFlag{Flag:key, Value:value, UpdateTime:time.Now()}
		this.flags[key] = flag
	} else {
		flag.Value = value
		flag.UpdateTime = time.Now()
	}
	return flag
}

func (this *RoleFlagManager) GetFlagBooleanValue(key db.ROLE_FLAG) bool {
	return this.GetFlagValue(key) != DEFAULT_FLAG_VALUE
}

func (this *RoleFlagManager) UpdateBooleanFlag(key db.ROLE_FLAG, value bool) *db.DBRoleFlag {
	if (value) {
		return this.UpdateFlag(key, 1)
	} else {
		return this.UpdateFlag(key, DEFAULT_FLAG_VALUE)
	}
}



