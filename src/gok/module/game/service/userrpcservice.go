//管理网络消息的处理
package service

import (
	"gok/module/game/util"
	"gok/service"
	"gok/service/msg/protocol"
	"gok/constant"
	"gok/module/game/cache"
	//"gok/module/gate/rpc"
	"gok/module/game/user"

	"aliens/log"
	"github.com/gogo/protobuf/proto"
	"gok/module/game/db"
	"gok/service/lpc"
	"gopkg.in/mgo.v2/bson"
)

//处理用户RPC消息
var LocalUserRPCService = service.NewLocalService(service.SERVICE_USER_RPC)



//初始化用户服务容器
func init () {
	LocalUserRPCService.SetFilter(PersistFilter)
	//rpc.UserServiceProxy.RegisterPush(UserRPCPush)
	LocalUserRPCService.RegisterHandler(500, new(KickOffService))
	LocalUserRPCService.RegisterHandler(507, new(TransferUserSessionService))
	LocalUserRPCService.RegisterHandler(560, new(DrawDayGiftService))//微信公众号领取每日礼包
	LocalUserRPCService.RegisterHandler(511, new(UpdateUnlockFlagService))
	LocalUserRPCService.RegisterHandler(515, new(QueryByNicknameService))
	LocalUserRPCService.RegisterHandler(1000, new(PushService))
}

func PersistFilter(message *protocol.C2GS) bool {
	uid := message.GetParam()
	//不是用户持久化消息，不需要过滤
	if uid == 0 {
		return false
	}
	//用户id无效不需要处理
	if !cache.UserCache.IsUserExist(uid) {
		return true
	}
	//rpc推送消息处理
	userSession := user.Manager.GetAuthSession(uid)
	message.Param = 0

	if userSession != nil {
		message.Offline = !userSession.IsOnline()
		userSession.SyncMessage(message)
	} else {
		//用户会话不存在，消息持久化存储，等下次会话初始化再处理
		message.Offline = true
		messageData, err := proto.Marshal(message)
		if err != nil {
			log.Debug("marshal persist message %v err : %v", message, err)
			return true
		}
		lpc.DBServiceProxy.Insert(&db.DBMessage{ID: bson.NewObjectId(), UserID:uid, Data:messageData}, db.DatabaseHandler)
	}
	return true
}

//用户T人
type KickOffService struct {
}

func (service *KickOffService) Request(request *protocol.C2GS, response *protocol.GS2C, network service.IMessageChannel) {
	message := request.GetKickoff()
	if message.GetUid() == -1 {
		user.Manager.KickAll(constant.LOGOUT_TYPE(message.GetKickType()))
	} else {
		session := user.Manager.GetOnlineAuthSession(message.GetUid())
		if session == nil {
			return
		}
		session.Logout(constant.LOGOUT_TYPE(message.GetKickType()))
		//session.SyncCommand(command.KICK, )
	}
}

//推送消息
type PushService struct {
}

func (service *PushService) Request(request *protocol.C2GS, response *protocol.GS2C, network service.IMessageChannel) {
	message := request.GetUserPush()
	uid := message.GetUid()
	if uid == -1 {
		user.Manager.BroadcastAll(message.GetMessage())
	} else {
		session := user.Manager.GetOnlineAuthSession(uid)
		if session != nil {
			session.WriteMsg(message.GetMessage())
		}
	}
}

type TransferUserSessionService struct {
}

func (service *TransferUserSessionService) Request(request *protocol.C2GS, response *protocol.GS2C, network service.IMessageChannel) {
	message := request.GetTransferRemoteUserSession()

	session := user.Manager.GetAuthSession(message.GetUid())
	response.TransferRemoteUserSessionRet = &protocol.TransferRemoteUserSessionRet{}

	if session != nil {
		session.Logout(constant.LOGOUT_TYPE_OTHER)
		session.Release()
		response.TransferRemoteUserSessionRet.Session = session.GetPersistData()
	}

}

type DrawDayGiftService struct {

}

func (service *DrawDayGiftService) Request(request *protocol.C2GS, response *protocol.GS2C, network service.IMessageChannel) {
	message := request.GetDrawDayGift()
	session := user.Manager.GetAuthSession(message.GetUid())
	resp := &protocol.DrawDayGiftRet{}
	if session != nil {
		session.DrawDayGift()
		//同步数据给客户端
		session.WriteMsg(session.BuildRoleSocialPush())
	}
	response.DrawDayGiftRet = resp
}

type UpdateUnlockFlagService struct {

}

func (service *UpdateUnlockFlagService) Request(request *protocol.C2GS, response *protocol.GS2C, network service.IMessageChannel) {
	message := request.GetUpdateUnlockFlag()
	session := user.Manager.GetAuthSession(message.GetUid())

	for _, flag := range message.GetFlags() {
		push := session.SetStarFlag(flag.GetId(), flag.GetValue())
		if push {
			session.WriteMsg(util.BuildStarFlagPush(flag))
		}
	}
}

type QueryByNicknameService struct {

}

func (service *QueryByNicknameService) Request(request *protocol.C2GS, response *protocol.GS2C, network service.IMessageChannel) {
	message := request.GetQueryByNickname()
	message.GetNickname()
	var roles []*db.DBRole
	nickname := message.GetNickname()
	db.DatabaseHandler.QueryAllCondition(&db.DBRole{}, "nickname", nickname, &roles)
	var result []int32
	for _, role := range roles {
		result = append(result, role.UserID)
	}
	response.QueryByNicknameRet = &protocol.QueryByNicknameRet{Uids:result}
}