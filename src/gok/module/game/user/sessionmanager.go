package user

import (
	"gok/module/cluster/center"
	"time"
	"gok/constant"
	"gok/service/msg/protocol"
	"aliens/common/collection"
	"github.com/name5566/leaf/gate"
	"aliens/common/util"
	"gok/module/cluster/cache"
	"gok/service/lpc"
	"gok/module/game/conf"
	"aliens/log"
	"gok/service/rpc"
)

var Manager = newUserSessionManager()

func newUserSessionManager() *SessionManager {
	manager := &SessionManager{
		authOnline:  make(map[int32]*Session),
		authOffline: make(map[int32]*Session),
		visitor:     &collection.Map{},
	}
	//心跳精确到s
	//manager.timeWheel = util.NewTimeWheel(time.Second, 60, manager.dealAuthTimeout)
	//manager.timeWheel.Start()
	return manager
}

//用户会话管理类，管理所有连接到此台服务器的用户会话
type SessionManager struct {
	//sync.RWMutex
	//*chanrpc.Server
	visitor     *collection.Map    //存储所有未验权的会话
	authOnline  map[int32]*Session //存储验权通过的在线会话
	authOffline map[int32]*Session //登出后缓存验权通过的会话 会有过期时间和上限限制
	timeWheel   *util.TimeWheel    //验权检查时间轮
}

func Init() {

}

func Close() {
	Manager.Close()
}

/**
 *  新建用户
 */
func (this *SessionManager) NewSession(agent gate.Agent) *Session {
	session := newSession(agent)
	this.visitor.UnsafeSet(session, &struct{}{})
	return session
}

//func (this *SessionManager) dealAuthTimeout(data util.TaskData) {
//	network := data[0].(*session)
//	超过固定时长没有验证权限需要提出
//	if network.IsAuthTimeout() {
//		log.Debug("network auth timeout : %v", network.GetRemoteAddr())
//		network.Close()
//		this.networks.Del(network)
//	}
//}

//func (this *SessionManager) ReleaseSession(session *Session) {
//	if this.Server != nil {
//		this.Go("Release", session)
//	}
//}

func (this *SessionManager) RemoveSession(session *Session) {
	delete(this.authOffline, session.GetID())
	delete(this.authOnline, session.GetID())
	cache.Cluster.SetUserNode(session.GetID(), "")
}

/**
 *  释放用户
 */
func (this *SessionManager) LogoutSession(session *Session, auth bool) {
	if auth {
		oldSession := this.authOnline[session.GetID()]
		if oldSession == session {
			delete(this.authOnline, session.GetID())
			this.authOffline[session.GetID()] = session
			cache.Cluster.SetUserOnlineNode(session.GetID(), "")

			//delete(this.authOnline, session.GetID())
			//this.authOffline[session.GetID()] = session
			//cache.Cluster.SetUserNode(session.GetID(), "")
		}

		//delete(this.authOnline, session.GetID())
		//this.authOffline[session.GetID()] = session
		//cache.Cluster.SetUserOnlineNode(session.GetID(), "")
		//log.Debug("session release %v-%v %v", session.GetID(), len(this.auth), this.visitor.Len())
	} else {
		this.visitor.UnsafeDel(session)
	}
}

func (this *SessionManager) LoginSession(session *Session) {
	this.visitor.UnsafeDel(session)
	delete(this.authOffline, session.GetID())
	this.authOnline[session.GetID()] = session
	cache.Cluster.SetUserNode(session.GetID(), center.GetServerNode())
	cache.Cluster.SetUserOnlineNode(session.GetID(), center.GetServerNode())
}

/**
 *  获取用户
 */
func (this *SessionManager) GetOnlineAuthSession(id int32) *Session {
	return this.authOnline[id]
}

//func (this *SessionManager) GetOfflineAuthSession(id int32) *Session {
//	return this.authOffline[id]
//}

func (this *SessionManager) GetAuthSession(id int32) *Session {
	session := this.authOnline[id]
	if session != nil {
		return session
	}
	return this.authOffline[id]
}


//释放空闲时间过长的用户内存
func (this *SessionManager) DealUserTimer() {
	//if constant.ES_LOG {
	//	rpc.StarServiceProxy.PrintCallStatistics()
	//}
	currTime := time.Now()

	lpc.StatisticsHandler.AddOnlineStatistic(len(this.authOnline), this.visitor.UnsafeLen())

	results := make([]int32, 0)

	for _, session := range this.authOnline {
		if session.IsFreeTimeout(currTime, conf.Server.UserFreeTimeout) {
			session.Logout(constant.LOGOUT_TYPE_TIMEOUT)
		}
		if session.IsDirty() && session.IsSyncDBTimeout(currTime, conf.Server.SyncDBInterval) {
			session.UpdateData()
		}
		results = append(results, session.GetID())
		session.DealAutoHelp(currTime)
	}

	//更新活跃时间
	rpc.SearchServiceProxy.UpdateHelpDatas(results, constant.SEARCH_OPT_UPDATE_ACTIVE, 0)


	for _, offline := range this.authOffline {
		//释放长久没上线的缓存数据
		if offline.IsLogoutTimeout(currTime, conf.Server.LocalCacheTimeout) {
			offline.Release()
		}
	}
}

func (this *SessionManager) Close() {
	log.Debug("session release start online-%v offline-%v ....", len(this.authOnline), len(this.authOffline))
	startTime := time.Now()
	for _, session := range this.authOnline {
		session.Logout(constant.LOGOUT_TYPE_CLOSE)
		session.Release()
	}
	for _, session := range this.authOffline {
		session.Release()
	}
	log.Debug("session release end duration(%v) ...", time.Now().Sub(startTime).Seconds())
}

//广播所有玩家
func (this *SessionManager) BroadcastAll(message *protocol.GS2C) {
	for _, session := range this.authOnline {
		session.WriteMsg(message)
	}
}

func (this *SessionManager) KickAll(kickType constant.LOGOUT_TYPE) {
	for _, session := range this.authOnline {
		session.Logout(kickType)
		//session.SyncCommand(command.KICK, kickType)
	}
}
