package user

import (
	"aliens/log"
	"gok/module/game/cache"
	"gok/module/game/conf"
	"gok/module/game/user/manager"
	"gok/service/msg/protocol"
	"time"
	"gok/constant"
	"github.com/name5566/leaf/gate"
	"gok/module/game/db"
	"github.com/gogo/protobuf/proto"
	//"github.com/name5566/leaf/chanrpc"
	"gok/service"
	"strings"
	"aliens/common/character"
	"gok/service/rpc"
	"gok/module/game/util"
	"gok/service/exception"
	"gopkg.in/mgo.v2/bson"
	"gok/service/lpc"
	"gok/module/game/global"
	"gok/module/statistics/model"
)

func newSession(agent gate.Agent) *Session {
	session := &Session{agent:agent}
	//session.Open()
	return session
}

//用户session
type Session struct {
	//sync.RWMutex

	*manager.DataManager //数据管理句柄
	//server *chanrpc.Server
	agent gate.Agent //网络连接句柄

	lastActiveTime time.Time //上一次活动消息的时间
	lastSyncDBTime time.Time //上一次自动同步到数据库的时间
}


func (this *Session) WriteMsg(msg interface{}) {
	if this.agent != nil {
		this.agent.WriteMsg(msg)
	}
}

func (this *Session) IsOnline() bool {
	return this.agent != nil
}



//func (this *Session) AssertSession() bool {
//	userNode := cache.UserCache.GetUserNode(this.GetID())
//	if userNode == "" {
//		cache.UserCache.SetUserNode(this.GetID(), center.GetServerNode())
//		return true
//	}
//	//用户会话在其他服务器存在，需要释放本地会话
//	valid := (center.GetServerNode() == userNode)
//	if !valid {
//		log.Debug("#AssertSession# userid=%v local=%v remote=%v", this.GetID(), center.GetServerNode(), userNode)
//		//this.Kick(constant.LOGOUT_TYPE_OTHER)
//		//this.WriteMsg(util.BuildKickPush(constant.LOGOUT_TYPE_OTHER))
//		this.BlockKick(constant.LOGOUT_TYPE_OTHER)
//		//节点信息在其他服务器更新，不需要去清除用户节点信息
//		this.Release(false)
//	}
//	return valid
//}

//更新缓存和数据库
func (this *Session) UpdateData() {
	if this.DataManager != nil {
		this.lastSyncDBTime = time.Now()
		this.DataManager.UpdateAll()
		this.CleanDirty()
	}
}

func (this *Session) IsAuth() bool {
	return this.DataManager != nil
}

//加载和处理用户内存数据
func (this *Session) Login(auth interface{}, persistData []byte) {
	switch auth.(type) {
	case int32:
		this.DataManager = &manager.DataManager{}
		//加载用户数据
		this.DataManager.EnsureData(auth.(int32), persistData)
		//this.DataManager.ensure
		this.RoleItemManager.SetNewItemListener(this.DealNewItem)
		//处理离线期间的持久化消息
		var persistMessages []*db.DBMessage
		db.DatabaseHandler.QueryAllCondition(&db.DBMessage{}, "userid", this.GetID(), &persistMessages)
		for _, persistMessage := range persistMessages {
			message := &protocol.C2GS{}
			err := proto.Unmarshal(persistMessage.Data, message)
			if err == nil {
				service.ServiceManager.HandleChannelMessage(message, this)
			}
			//db.DatabaseHandler.DeleteOne(persistMessage)
		}
		//处理完毕统一删除
		condition := bson.D{{"userid", this.GetID()}}
		db.DatabaseHandler.DeleteAllCondition(&db.DBMessage{}, condition)
		break
	case *manager.DataManager:
		this.DataManager = auth.(*manager.DataManager)
		break
	default:
		log.Warn("invalid auth param %v", auth)
	}
	log.Debug("#login#: name=%v userid=%v network=%v", this.GetNickName(), this.GetID(), this.agent.RemoteAddr())
	loginStarInfo := rpc.StarServiceProxy.GetLoginStarInfo(this.GetID())
	this.LoginStarInit(loginStarInfo)
	this.SetLoginTime(time.Now())
	cache.UserCache.SetUserOnlineTimestamp(this.GetID(), time.Now())
	this.SetHandler(this)

	subscribe := cache.UserCache.IsUserSubscribe(this.GetID())
	this.SetSubscribe(subscribe)

	//更新活跃时间
	rpc.SearchServiceProxy.UpdateHelpData(this.GetID(), constant.SEARCH_OPT_UPDATE_ACTIVE, 0)

	//推送所有关注者上线
	onlinePush := util.BuildOnlinePush(this.GetID())
	global.PushFollowings(this.GetID(), onlinePush)
	Manager.LoginSession(this)

	//领取每日礼包
	this.DrawDayGift()
}

func (this *Session) Release() {
	//释放内存的时候需要存储用户数据到数据库
	if this.DataManager.IsDirty() {
		this.UpdateData()
	}
	if this.GetGuideTime().Unix() > 0 {
		lpc.LogServiceProxy.AddGuideRecord(this.GetID(), this.GetGuideTime(), this.GetLogoutTime())
	}
	Manager.RemoveSession(this)

}

func (this *Session) Logout(logoutType constant.LOGOUT_TYPE) {
	//用户已经登出了，不需要再处理登出操作
	if this.agent == nil || this.DataManager == nil {
		return
	}
	//清除离线消息
	this.CleanOfflineMessage()
	this.SetLogoutTime(time.Now())

	offlinePush := util.BuildOfflinePush(this.GetID())
	global.PushFollowings(this.GetID(), offlinePush)

	//lpc.LogServiceProxy.AddLoginRecord(this.GetID(), this.GetNetworkAddress(), this.GetLoginTime(), this.GetLogoutTime())
	lpc.StatisticsHandler.AddStatisticData(&model.StatisticLogout{
		UserID:this.GetID(),
		Address:this.GetNetworkAddress(),
		Faith:this.GetFaith(),
		Power:this.GetPower(),
		BelieverTotal:cache.StarCache.GetBelieverTotalLevel(this.GetStarId()),
		BuildingTotal:cache.StarCache.GetBuildingAllLevel(this.GetStarId()),
		//Civil:cache.StarCache.GetCivilLevel(this.GetStarId()),
		ExMaxBuildLv: cache.StarCache.GetBuildingExMaxLevel(this.GetStarId()),
		Star:this.GetStarType(),
		Mutual:cache.StarCache.GetMutualTimes(this.GetStarId()),
		BeMutual:cache.StarCache.GetBeMutualTimes(this.GetStarId()),
		//Civil:this.ge
	})


	//未初始化,直接清楚会话
	if logoutType != constant.LOGOUT_TYPE_NONE {
		this.WriteMsg(util.BuildKickPush(logoutType))
	}
	log.Debug("#logout#: name=%v userid=%v network=%v", this.GetNickName(), this.GetID(), this.agent.RemoteAddr())
	//cache.UserCache.SetUserOnline(this.GetID(), false)
	Manager.LogoutSession(this, this.IsAuth())
	this.agent.Close()
	this.agent = nil
}


//是否登出时间超时需要释放
func (this *Session) IsLogoutTimeout(timestamp time.Time, timeout float64) bool {
	return !this.IsOnline() && timestamp.Sub(this.GetLogoutTime()).Seconds() > timeout
}

//用户空闲时间时间是否超时
func (this *Session) IsFreeTimeout(timestamp time.Time, timeout float64) bool {
	if timeout <= 0 {
		return false
	}
	interval := timestamp.Sub(this.lastActiveTime).Seconds()
	return interval > timeout
}

//同步数据库时间是否超时
func (this *Session) IsSyncDBTimeout(timestamp time.Time, timeout float64) bool {
	interval := timestamp.Sub(this.lastSyncDBTime).Seconds()
	return interval > timeout
}

func  (this *Session) GetNetworkAddress() string {
	if this.agent == nil {
		return ""
	}
	gate, ok := this.agent.(gate.Agent)
	if !ok {
		return ""
	}
	address := gate.RemoteAddr().String()
	array := strings.Split(address, ":")
	if array == nil || len(array) == 0 {
		return ""
	}
	return array[0]
}

func (this *Session) DialAddBeliever(believerType int32, believerNum int32) []*protocol.BelieverInfo{
	var level int32 = 0
	switch believerType {
	case constant.BELIEVER_L1:level = 1
		break
	case constant.BELIEVER_L2:level = 2
		break
	case constant.BELIEVER_L3:level = 3
		break
	case constant.BELIEVER_L4:level = 4
		break
	case constant.BELIEVER_L5:level = 5
		break
	case constant.BELIEVER_L6:level = 6
		break
	default:
		break
	}
	maleNum := believerNum/2
	femaleNum := believerNum - maleNum
	maleID := "b01" + character.Int32ToString(level) +"1"
	femaleID := "b01" + character.Int32ToString(level) +"2"
	addBeliever := []*protocol.BelieverInfo{{Id:maleID,Num:maleNum},{Id:femaleID,Num:femaleNum}}
	return addBeliever
}

func (this *Session) LoginStarInit(loginStarInfo *protocol.LoginStarInfoRet) {
	star := loginStarInfo.GetCurrStar()
	if star != nil {
		this.UpdateStarInfo(star.GetStarID(), star.GetStarType())
		this.EnsureItemRandom(star.GetStarType(),false)
		this.DealUpgradeBuilding(0, loginStarInfo.GetPowerLimit(), 0, loginStarInfo.GetUpgradedBuilding(), false)
		this.DealRepairedBuilding(loginStarInfo.GetRepairedBuilding(), false)
		this.UpdateStarFlags(loginStarInfo.GetStarFlags())

		//指定时间内没有交互任务加入到被随机到目标的队列中
		if !this.HasMutual(conf.DATA.SearchRandomCD) {
			rpc.SearchServiceProxy.UpdateRandomStar(star.GetStarID())
		}
	}
}

func (this *Session) NewStarInit() {

}

//func (this *Session) LoginStarInit(starInfo *protocol.StarInfo, powerLimit int32, upgradedBuilding []*protocol.BuildingState, repairedBuilding []*protocol.BuildingState) {
//	if starInfo != nil {
//		this.UpdateStarInfo(starInfo.GetStarID(), starInfo.GetStarType())
//		this.EnsureItemRandom(starInfo.GetStarType(),false)
//		this.DealUpgradeBuilding(powerLimit, upgradedBuilding, false)
//		this.DealRepairedBuilding(repairedBuilding, false)
//	}
//}

func (this *Session) DecFlag(flag db.ROLE_FLAG) int32{
	shareCount := this.GetFlagValue(flag)
	if shareCount <= 0 {
		exception.GameException(exception.SHARE_COUNT_NOT_ENOUGH)
	}
	return this.UpdateFlag(flag, shareCount - 1).Value
}

func (this *Session) IsExceedShareTime(confInterval int64) bool {
	timeNow := time.Now()
	if timeNow.After(this.WatchAdGetPower) {
		this.WatchAdGetPower = timeNow.Add(time.Duration(confInterval*int64(time.Second)))
		return true
	}
	return false
}

func (this *Session) IsDrawInviteLoginGift() bool {
	return this.GetFlagValue(db.FLAG_DRAW_INVITE_GIFT) == 0
}

func (this *Session) SetInviteID(inviteID int32) {
	this.InviteID = inviteID
}

func (this *Session) IsFirstLoginByInvite(inviteID int32) bool{
	return inviteID != 0 && this.InviteID == 0
}

func (this *Session) IsBuildNeedRepair(uid int32, buildType int32) bool {
	starID := cache.StarCache.GetUserActiveStar(uid)
	if cache.StarCache.GetStarBuildRepair(starID, buildType) == character.Int32ToString(0) {
		return true
	}
	return false
}

func (this *Session) IsFirstDrawAdReward(flag db.ROLE_FLAG) bool {
	return this.GetFlagValue(flag) == 0
}