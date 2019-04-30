package manager

import (
	"gok/module/game/db"
	"gok/service/exception"
	"gok/module/game/conf"
	"gok/module/game/cache"
	basecache "gok/cache"
	"gok/service/rpc"
	"time"
	"gok/constant"
	"gok/service/lpc"
	"gok/module/statistics/model"
)

//角色资产管理
type RoleBaseManager struct {
	*db.DBRole
	subscribe bool //是否关注公众号
	weightAppend map[int32]int32
}

//初始化
func (this *RoleBaseManager) Init(role *db.DBRole) {
	this.DBRole = role
	this.weightAppend = make(map[int32]int32)
	this.UpdatePower(time.Now())
}

//更新数据库内存
func (this *RoleBaseManager) Update(role *db.DBRole) {
	role = this.DBRole
}

func (this *RoleBaseManager)GetWeightAppendMapping() map[int32]int32 {
	for _, append := range conf.DATA.EVENT_WEIGHT_APPEND {
		var diff int32 = 0
		var addWeight int32 = 0
		if append.Type == constant.WEIGHT_ADD_TYPE_FAITH {
			diff = append.Limit - this.Faith
		} else if append.Type == constant.WEIGHT_ADD_TYPE_BELIEVER {
			diff = append.Limit - cache.StarCache.GetBelieverCount(this.GetID())
		}

		if diff > 0 {
			addWeight = int32 (diff / append.Util) * append.Add
		}
		this.weightAppend[append.ID] = addWeight
	}
	return this.weightAppend
}

//获取用户id
func (this *RoleBaseManager) GetID() int32 {
	return this.UserID
}

//获取用户游戏图标
//func (this *RoleBaseManager) GetIcon() int32 {
//	return this.Icon
//}

//获取用户游戏图标
func (this *RoleBaseManager) GetNickName() string {
	return this.NickName
}

func (this *RoleBaseManager) GetPublicTime() time.Time {
	return this.SalePublicTime
}

//func (this *RoleBaseManager) GetCivilizationLevel() int32 {
//	return this.CivilizationLevel
//}
//
//func (this *RoleBaseManager) GetCivilizationValue() int32 {
//	return this.CivilizationValue
//}

func (this *RoleBaseManager) RefreshPublicTime() time.Time {
	this.SalePublicTime = time.Now()
	return this.SalePublicTime
}

func (this *RoleBaseManager) GetHelpPublicTime() time.Time {
	return this.HelpPublicTime
}

func (this *RoleBaseManager) SetHelpPublicTime(time time.Time, publicItemID int32)  {
	this.HelpPublicTime = time
	this.HelpPublicItem = publicItemID
}

func (this *RoleBaseManager) GetHelpPublicItem() int32 {
	return this.HelpPublicItem
}


func (this *RoleBaseManager) CleanHelpPublicTime()  {
	this.HelpPublicTime = time.Unix(0, 0)
	this.HelpPublicItem = 0
}

func (this *RoleBaseManager) GetNextRandomItem() int32 {
	return this.NextRandomItem
}

func (this *RoleBaseManager) CleanNextRandomItem() {
	this.NextRandomItem = 0
}

func (this *RoleBaseManager) SetNextRandomItem(itemID int32) {
	this.NextRandomItem = itemID
}

//获取用户游戏图标
//func (this *RoleBaseManager) GetLevel() int32 {
//	return this.Level
//}


//获取用户游戏图标
//func (this *RoleBaseManager) GetExp() int32 {
//	return this.Exp
//}


//获取用户游戏图标
func (this *RoleBaseManager) GetPower() int32 {
	return this.Power
}

func (this *RoleBaseManager) GetPowerTime() time.Time {
	return this.PowerTime
}

func (this *RoleBaseManager) GetNextAdPowerTime() time.Time {
	return this.WatchAdGetPower
}

//获取下一次更新法力值的时间戳
func (this *RoleBaseManager) GetLastPowerTime() time.Time {
	restoreTime := conf.DATA.PowerRestoreTime + cache.StarCache.GetBuffMANAInterval(this.UserID)
	return this.PowerTime.Add(time.Duration(int64(time.Second) * int64(restoreTime)))
}

//
func (this *RoleBaseManager) SetSubscribe(subscribe bool) {
	this.subscribe = subscribe
}

func (this *RoleBaseManager) GetSubscribe() bool {
	return this.subscribe
}

//获取用户游戏图标
func (this *RoleBaseManager) GetPowerLimit() int32 {
	var appendLimit int32 = 0
	//关注加上限
	if this.subscribe {
		appendLimit = conf.DATA.PrivilegeReward
	}
	return this.PowerLimit + appendLimit
}

func (this *RoleBaseManager) GetPowerLimitBase() int32 {
	return this.PowerLimit
}

func (this *RoleBaseManager) SetPowerLimit(powerLimit int32) {
	this.PowerLimit = powerLimit
}

//获取用户信仰值
func (this *RoleBaseManager) GetFaith() int32 {
	return this.Faith
}

//获取用户钻石数量
func (this *RoleBaseManager)GetDiamond() int32{
	return this.Diamond
}

//获取友情点
func (this *RoleBaseManager)GetGayPoint() int32{
	return this.GayPoint
}

func (this *RoleBaseManager)GetDesc() string{
	return this.Desc
}

func (this *RoleBaseManager)SetDesc(desc string){
	this.Desc = desc
	cache.UserCache.SetUserDesc(this.UserID, this.Desc)
}

//增加神力-经验值
//func (this *RoleBaseManager) TakeInExp(exp int32) {
//	for  {
//		if (exp <= 0) {
//			break;
//		}
//		//当前等级还需要增加的经验
//		levelRemainExp := conf.DATA.EXP_DATA[this.Level + 1].Consumption - this.Exp
//		if levelRemainExp > exp {
//			//经验不够升级,直接添加经验
//			this.Exp = this.Exp + exp
//			break;
//		}
//		this.Level = this.Level + 1;
//		cache.UserCache.SetUserAttr(this.UserID, basecache.UPROP_LEVEL, this.Level)
//		this.Exp = 0;
//		exp = exp - levelRemainExp
//	}
//}

func (this *RoleBaseManager) TakeInGayPoint(gaypoint int32, opt constant.OPT, refID int32) bool {
	if gaypoint <= 0 {
		return false
	}
	this.GayPoint += gaypoint
	lpc.StatisticsHandler.AddStatisticData(&model.StatisticGayPoint{
		UserID:this.UserID,
		RefID:refID,
		Operation:uint8(opt),
		Change:gaypoint,
		Total:this.GayPoint,
	})
	return true
}


func (this *RoleBaseManager) TakeOutGayPoint(gaypoint int32, opt constant.OPT, refID int32) bool {
	if this.GayPoint < gaypoint {
		exception.GameException(exception.GAYPOINT_NOT_ENOUGH)
	}
	this.GayPoint -= gaypoint

	lpc.StatisticsHandler.AddStatisticData(&model.StatisticGayPoint{
		UserID:this.UserID,
		RefID:refID,
		Operation:uint8(opt),
		Change:-gaypoint,
		Total:this.GayPoint,
	})
	return true
}

func (this *RoleBaseManager) EnsureGayPoint(gaypoint int32) {
	if this.GayPoint < gaypoint {
		exception.GameException(exception.GAYPOINT_NOT_ENOUGH)
	}
}

//获取信仰值-金币
func (this *RoleBaseManager) TakeInFaith(faith int32, operation constant.OPT, refID int32) bool {
	if faith <= 0 {
		return false
	}
	this.Faith += faith
	cache.UserCache.SetUserAttr(this.UserID, basecache.UPROP_FAITH, this.Faith)
	lpc.LogServiceProxy.AddSocialRecord(this.UserID, constant.SOCIAL_ID_FAITH, 0, operation, faith ,this.Faith)
	return true
}

func (this *RoleBaseManager) EnsureFaith(faith int32) {
	if (this.Faith < faith) {
		exception.GameException(exception.FAITH_NOT_ENOUGH)
	}
}

func (this *RoleBaseManager) AssertFaith(faith int32) {
	if (this.Faith < faith) {
		exception.GameException(exception.FAITH_NOT_ENOUGH)
	}
}

func (this *RoleBaseManager) CanTakeOutFaith(faith int32) bool {
	return this.Faith >= faith
}

//消耗信仰值-金币
func (this *RoleBaseManager) TakeOutFaith(faith int32, operation constant.OPT, refID int32) {
	if this.Faith < faith {
		exception.GameException(exception.FAITH_NOT_ENOUGH)
	}
	this.Faith -= faith
	cache.UserCache.SetUserAttr(this.UserID, basecache.UPROP_FAITH, this.Faith)
	lpc.LogServiceProxy.AddSocialRecord(this.UserID, constant.SOCIAL_ID_FAITH, refID, operation, -faith ,this.Faith)
}

func (this *RoleBaseManager) ForceTakeOutFaith(faith int32, operation constant.OPT, refID int32) int32 {
	oldFaith := this.Faith
	this.Faith -= faith
	if this.Faith < 0 {
		this.Faith = 0
	}
	changeFaith :=  this.Faith - oldFaith
	cache.UserCache.SetUserAttr(this.UserID, basecache.UPROP_FAITH, this.Faith)
	lpc.LogServiceProxy.AddSocialRecord(this.UserID, constant.SOCIAL_ID_FAITH, refID, operation, changeFaith, this.Faith)
	return changeFaith
}

//更新法力值
func (this *RoleBaseManager) UpdatePower(timestamp time.Time) {
	interval := timestamp.Sub(this.PowerTime).Seconds()
	restoreTime := conf.DATA.PowerRestoreTime + cache.StarCache.GetBuffMANAInterval(this.UserID)
	ratio := int32(interval / float64( restoreTime ) )//获取时间间隔
	if ratio <= 0 {
		return
	}
	if this.TakeInPower( ratio * conf.DATA.PowerRestorePoint, true , constant.OPT_TYPE_REFRESH_POWER) {
		duration := time.Duration( int64( restoreTime ) * int64( time.Second ))
		if this.IsPowerFull() {
			this.PowerTime = timestamp
		} else {
			this.PowerTime = this.PowerTime.Add(duration)//更新增加法力值的时间
		}
		rpc.StarServiceProxy.UpdateStarStatistics(this.GetID(), constant.STAR_STATISTIC_GAIN_POWER_AUTO, float64(ratio * conf.DATA.PowerRestorePoint), 0)
	}
}

func (this *RoleBaseManager) IsPowerFull() bool {
	return this.Power >= this.GetPowerLimit()
}

//获取法力值
func (this *RoleBaseManager) TakeInPower(power int32, limit bool, operation constant.OPT) bool {
	if power <= 0 {
		return false
	}

	oldPower := this.Power

	resultPower := this.Power + power
	powerLimit := this.GetPowerLimit()
	//是否有法力值上限限制和超出法力值上限
	if limit && resultPower > powerLimit {
		if this.Power < powerLimit  {
			this.Power = powerLimit
		} else {
			return false
		}
	} else {
		this.Power = resultPower
	}

	//AddSocialRecord(uid int32, socialID SOCIAL_ID, refID int32, operation int32, change int32, total int32) {

	cache.UserCache.SetUserAttr(this.UserID, basecache.UPROP_POWER, this.Power)
	lpc.LogServiceProxy.AddSocialRecord(this.UserID, constant.SOCIAL_ID_POWER, powerLimit, operation, this.Power - oldPower ,this.Power)
	return true
}

//消耗法力值
func (this *RoleBaseManager) EnsurePower(power int32) {
	if this.Power < power {
		exception.GameException(exception.POWER_NOT_ENOUGH)
	}
}


//消耗法力值
func (this *RoleBaseManager) TakeOutPower(power int32, operation constant.OPT) {
	isFull := this.IsPowerFull()
	oldPower := this.Power

	if this.Power < power {
		exception.GameException(exception.POWER_NOT_ENOUGH)
	}
	this.Power -= power

	//从法力值满到未满，需要开启法力值动态更新
	if isFull && !this.IsPowerFull() {
		this.PowerTime = time.Now()
	}

	cache.UserCache.SetUserAttr(this.UserID, basecache.UPROP_POWER, this.Power)
	lpc.LogServiceProxy.AddSocialRecord(this.UserID, constant.SOCIAL_ID_POWER, this.GetPowerLimit(), operation, this.Power - oldPower ,this.Power)
}

//获取钻石
func ( this *RoleBaseManager)TakeInDiamond( diamond int32, operation constant.OPT, refID int32 ) bool {
	if diamond <= 0 {
		return false
	}
	this.Diamond += diamond
	lpc.LogServiceProxy.AddSocialRecord(this.UserID, constant.SOCIAL_ID_DIAMOND, 0, operation, diamond ,this.Diamond)
	return true
}

//消耗钻石
func ( this *RoleBaseManager)TakeOutDiamond( diamond int32, operation constant.OPT, refID int32 ) bool{
	if this.Diamond < diamond {
		return false
	}
	this.Diamond -= diamond
	lpc.LogServiceProxy.AddSocialRecord(this.UserID, constant.SOCIAL_ID_DIAMOND, refID, operation, -diamond ,this.Diamond)
	return true
}

func ( this *RoleBaseManager)EnsureDiamond( diamond int32 ) {
	if this.Diamond < diamond {
		exception.GameException(exception.DIAMOND_NOT_ENOUGH)
	}
}

func ( this *RoleBaseManager)AssertDiamond( diamond int32 ) {
	if this.Diamond < diamond {
		exception.GameException(exception.DIAMOND_NOT_ENOUGH)
	}
}

func (this *RoleBaseManager) SetStarsSelect(starsType []int32) {
	this.StarsSelect = starsType
}

func (this *RoleBaseManager) CleanStarsSelect() {
	this.StarsSelect = []int32{}
}

func (this *RoleBaseManager) EnsureStarsSelect() ([]int32, bool) {
	if this.StarsSelect == nil || len(this.StarsSelect) == 0 {
		return nil, false
	}
	return this.StarsSelect, true
}

func (this *RoleBaseManager) EnsureSelectInStars(selectType int32) bool {
	for _, starType := range this.StarsSelect {
		if 	selectType == starType {
			return true
		}
	}
	return false
}