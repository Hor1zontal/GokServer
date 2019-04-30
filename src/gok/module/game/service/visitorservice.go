//管理网络消息的处理
package service

import (
	"gok/service/msg/protocol"
	"time"
	//"gok/passport/cache"
	"gok/module/game/user"
	baseservice "gok/service"
	"gok/service/exception"
	"gok/module/cluster/center"
	"gok/constant"
	"aliens/log"
	"gok/service/rpc"
	"gok/module/cluster/cache"
	basecache "gok/module/game/cache"
)

var VisitorService = baseservice.NewLocalService(baseservice.SERVICE_VISITOR)

func init() {
	VisitorService.RegisterHandler(11, new(LoginServerService))
}

//登录服务处理类
type LoginServerService struct {
}

func (service *LoginServerService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	//达到上限不允许登录
	//if conf.server.MaxSession > 0 && user.Manager.GetUserTotal() >= conf.server.MaxSession {
	//	network.WriteMsg(util.BuildKickPush(constant.LOGOUT_TYPE_MAX_SESSION))
	//	network.Close()
	//}

	loginRequest := request.GetLoginServer()
	session := network.(*user.Session)

	loginRequest.Ip = session.GetNetworkAddress()
	uid := loginRequest.GetUserId()
	if uid <= 0 || loginRequest.GetToken() == "" {
		exception.GameException(exception.USER_INVALID_TOKEN)
	}

	//账号服务器验证token
	request.Sequence[0] = 9
	rpc.PassportServiceProxy.HandleMessage(request)
	request.Sequence[0] = 11

	oldSession := user.Manager.GetAuthSession(uid)
	if oldSession == session {
		exception.GameException(exception.LOGIN_REPEAT)
	}

	if oldSession != nil {
		session.Login(oldSession.DataManager, nil)
		if oldSession.IsOnline() {
			oldSession.Logout(constant.LOGOUT_TYPE_OTHER)
		}
	} else {
		//其他服务器节点的用户会话的需要转移到当前服务器
		persistData := TransferRemoteUserSession(uid)
		session.Login(uid, persistData)
	}
	//更新登录的版本号
	//session.SetClientVersion(loginRequest.GetVersion())
	response.LoginServerRet = &protocol.LoginSeverRet{
		UserId:     uid,
		ServerTime: time.Now().Unix(),
		RoleInfo:   BuildRoleInfo(session),
	}
	inviteType := loginRequest.InviteType
	inviteID := loginRequest.GetInviteID()
	HandleWeChatInviteMessage(session, inviteType, inviteID)
}

func HandleWeChatInviteMessage(session *user.Session, inviteType *protocol.InviteType,inviteID int32) {
	if inviteID == session.GetID() {
		return
	}
	if inviteType.Type == constant.PUBLIC_WECHAT_TYPE_HELP {
		//通过法力求助分享进来
		if inviteType.RefType == constant.WECHAT_HELP_REF_POWER {
			if session.IsFirstLoginByInvite(inviteID) {
				rpc.CommunityServiceProxy.FollowEach(session.GetID(), inviteID)
				newsFeedMessage := BuildNewsFeedMessage(session.GetID(), constant.NEWSFEED_TYPE_SHARE_SUCC, 0, 0, 0)

				AddUserNewsFeed(session, BuildNewsFeed(inviteID, constant.NEWSFEED_TYPE_MUTUAL_FOLLOW, 0, 0, 0))
				rpc.UserServiceProxy.PersistCall(inviteID, newsFeedMessage)
				session.SetInviteID(inviteID)
			}
		}
		//通过建筑维修求助分享进来
		if inviteType.RefType == constant.WECHAT_HELP_REF_REPAIR {
			if session.IsBuildNeedRepair(inviteID, inviteType.GetRefNum()) {
				resp := rpc.StarServiceProxy.Call(inviteID, &protocol.C2GS{
					Sequence:[]int32{607},
					HelpRepairBuild:&protocol.HelpRepairBuild{
						BuildingType:inviteType.GetRefNum(),
						HelperID:session.GetID(),
					},
				}).GetHelpRepairBuildRet()
				if resp.GetResult() {
					newsFeedMessage := BuildNewsFeedMessage(session.GetID(), constant.NEWSFEED_TYPE_BE_HELP_REPAIR, resp.GetStarType(), inviteType.GetRefNum(), resp.GetBuildingLevel())
					rpc.UserServiceProxy.PersistCall(inviteID, newsFeedMessage)
				}
			}
		}
		//通过圣物求助分享进来
		if inviteType.RefType == constant.WECHAT_HELP_REF_ITEM {
			if inviteType.RefNum == 0 {
				return
			}
			if !basecache.UserCache.ExistBeHelpItemCacheUid(inviteID, inviteType.GetRefNum()) || basecache.UserCache.ExistHelpItemWechatUid(session.GetID(), inviteID, inviteType.GetRefNum()) {
				return
			}
			request := &protocol.C2GS{
				Sequence: []int32{480},
				HelpItem:&protocol.HelpItem{
					Uid:inviteID,
					ItemID:inviteType.GetRefNum(),
				},
			}
			result := session.HelpItemAsk(true, request)
			if result != nil {
				AddUserNewsFeed(session, BuildNewsFeed(inviteID, constant.NEWSFEED_TYPE_HELP_ITEMHELP, result.GetItemHelp().GetItemID(), result.GetItemHelp().GetItemNum(), 0))
				basecache.UserCache.SetHelpItemWechatUid(session.GetID(), inviteID, inviteType.GetRefNum())
			}
		}
	}
}

func TransferRemoteUserSession(uid int32) []byte {
	node := cache.Cluster.GetUserNode(uid)
	//不是本节点需要发送T人消息到远程节点
	if node != "" && node != center.GetServerNode() {
		log.Debug("remote transfer user session data %v - %v", node, center.GetServerNode())
		return rpc.UserServiceProxy.TransferUserSession(uid, node)
	}
	return nil
}
