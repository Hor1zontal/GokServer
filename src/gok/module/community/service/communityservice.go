//管理网络消息的处理
package service

import (
	"gok/service/msg/protocol"
	baseservice "gok/service"
	"gok/module/community/cache"
	"gok/module/community/conf"
	"gok/service/exception"
	"gok/module/community/core"
	"aliens/common/character"
	"gok/module/community/db"
	"time"
	"gok/constant"
	"github.com/name5566/leaf/chanrpc"
	"gok/service/lpc"
)

var CommunityLocalService = baseservice.NewLocalService(baseservice.SERVICE_COMMUNITY_RPC)
var CommunityRPCService *baseservice.GRPCService = nil

//初始化事件服务消息
func Init(chanRpc *chanrpc.Server) {

	CommunityLocalService.RegisterHandler(220, new(FollowService))       //关注用户
	CommunityLocalService.RegisterHandler(221, new(UnFollowService))     //取消关注
	CommunityLocalService.RegisterHandler(222, new(FollowerListService)) //获取关注列表
	CommunityLocalService.RegisterHandler(224, new(FollowingListService)) //获取关注列表
	CommunityLocalService.RegisterHandler(225, new(FollowEachService))

	CommunityLocalService.RegisterHandler(226, new(FollowStateService)) //

	//好友相关
	//CommunityLocalService.RegisterHandler(230, new(GetFriendListService)) //获取好友列表
	////CommunityLocalService.RegisterHandler(231,new(GetFriendInfoService))
	//CommunityLocalService.RegisterHandler(232, new(AddFriendRequestService))
	//CommunityLocalService.RegisterHandler(233, new(DeleteFriendService))
	//CommunityLocalService.RegisterHandler(234, new(AcceptFriendRequestService))
	//CommunityLocalService.RegisterHandler(235, new(RefuseFriendRequestService))
	//CommunityLocalService.RegisterHandler(237, new(GetFriendRequestListService))

	CommunityLocalService.RegisterHandler(240, new(GetReceiveMomentsService))
	CommunityLocalService.RegisterHandler(241, new(GetPublicMomentsService))
	CommunityLocalService.RegisterHandler(243, new(PublicMomentService))
	CommunityLocalService.RegisterHandler(244, new(RemoveMomentsService))

	//配置了RPC，需要发布服务到ZK
	CommunityRPCService = baseservice.PublicRPCService1(CommunityLocalService, conf.Server.RPCAddress, conf.Server.RPCPort, chanRpc)
}

func Close() {
	CommunityLocalService.Close()
	CommunityRPCService.Close()
}

//生成朋友圈消息
type PublicMomentService struct {
}

func (service *PublicMomentService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetPublicMoment()
	//生成朋友圈消息
	moment := core.Moments.AddMoment1(message.GetUid(), constant.MomentsType(message.GetType()), message.GetRefID())
	response.PublicMomentRet = &protocol.PublicMomentRet{
		MomentInfo:moment,
	}
}

//获取收到的朋友圈消息
type GetReceiveMomentsService struct {
}

func (service *GetReceiveMomentsService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGetReceiveMoments()
	momentIds := cache.CommunityCache.GetUserReceiveMomentIDs(message.GetUid(), int32(message.GetBeforeTime()), message.GetOffset(), message.GetCount())
	moments := core.Moments.GetMomentsByIds(momentIds)
	response.GetReceiveMomentsRet = &protocol.GetReceiveMomentsRet{
		Moments: moments,
	}
}

//获取指定用户发布的朋友圈消息
type GetPublicMomentsService struct {
}

func (service *GetPublicMomentsService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGetPublicMoments()
	momentIds := cache.CommunityCache.GetUserPublicMomentIDs(message.GetUid(), int32(message.GetBeforeTime()), message.GetOffset(), message.GetCount())
	moments := core.Moments.GetMomentsByIds(momentIds)
	response.GetPublicMomentsRet = &protocol.GetPublicMomentsRet{
		Moments: moments,
	}
}

type RemoveMomentsService struct {
}

func (service *RemoveMomentsService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetRemoveMoments()
	core.Moments.RemoveSaleMoments(message.GetSaleID())
}



//获取好友列表
//type GetFriendListService struct {
//}
//
//func (service *GetFriendListService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
//	response.GetFriendListRet = core.GetFriendList(request.GetGetFriendList().GetId())
//}

//获取好友信息
//type GetFriendInfoService struct {
//}
//
//func (service *GetFriendInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
//	//message := request.GetGetFriendDetailInfo()
//	//message.GetId()
//}

////获取好友申请列表
//type GetFriendRequestListService struct {
//}
//
//func (service *GetFriendRequestListService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
//	response.GetFriendRequestListRet = core.BuildFriendRequestList(request.GetGetFriendRequestList().GetId())
//}

//添加好友
//type AddFriendRequestService struct {
//}
//
//func (service *AddFriendRequestService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
//	message := request.GetAddFriendRequest()
//	userId := message.GetId()
//	requestId := message.GetRequestID()
//
//	if (userId == requestId) {
//		exception.GameException(exception.CANNOT_ADD_MYSELF)
//	}
//
//	//已经是好友
//	if (cache.CommunityCache.ExistFriend(character.Int32ToString(userId), character.Int32ToString(requestId))) {
//		exception.GameException(exception.REPEATE_ADD_FRIEND)
//	}
//
//	req := &db.DFriendReq{
//		ID:      db.CompID{requestId, userId},
//		AddTime: time.Now(),
//	}
//	//TODO 推送好友申请  存储数据库
//	db.DatabaseHandler.Insert(req)
//
//	//添加好友请求
//	cache.CommunityCache.AddFriendRequest(character.Int32ToString(requestId), character.Int32ToString(userId))
//
//	response.AddFriendRequestRet = &protocol.AddFriendRequestRet{
//		Result: true),
//	}
//}

//接受好友申请
//type AcceptFriendRequestService struct {
//}
//
//func (service *AcceptFriendRequestService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
//	message := request.GetAcceptFriendRequest()
//	if (message.GetId() == message.GetRequestID()) {
//		exception.GameException(exception.CANNOT_ADD_MYSELF)
//	}
//
//	userId := character.Int32ToString(message.GetId())
//	requestID := character.Int32ToString(message.GetRequestID())
//
//	if (!cache.CommunityCache.ExistFriendRequest(userId, requestID)) {
//		exception.GameException(exception.FRIEND_REQUEST_NOFOUND)
//	}
//
//	if (!cache.CommunityCache.ExistFriend(userId, requestID)) {
//		addTime := time.Now()
//		relation := db.CompID{message.GetId(), message.GetRequestID()}
//		invRelation := db.CompID{message.GetRequestID(), message.GetId()}
//		db.DatabaseHandler.Insert(&db.DFriend{
//			ID:      relation,
//			AddTime: addTime,
//		})
//		db.DatabaseHandler.Insert(&db.DFriend{
//			ID:      invRelation,
//			AddTime: addTime,
//		})
//		//更新好友关系
//		cache.CommunityCache.AddFriend(userId, requestID)
//		cache.CommunityCache.AddFriend(requestID, userId)
//
//		//删除好友请求
//
//		db.DatabaseHandler.DeleteOne(&db.DFriendReq{ID: relation})
//		db.DatabaseHandler.DeleteOne(&db.DFriendReq{ID: invRelation})
//		cache.CommunityCache.DelFriendRequest(userId, requestID)
//		cache.CommunityCache.DelFriendRequest(requestID, userId)
//	}
//	response.AcceptFriendRequestRet = &protocol.AcceptFriendRequestRet{Result: true)}
//}
//
////拒接好友申请
//type RefuseFriendRequestService struct {
//}
//
//func (service *RefuseFriendRequestService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
//	message := request.GetRefuseFriendRequest()
//	userId := character.Int32ToString(message.GetId())
//	requestID := character.Int32ToString(message.GetRequestID())
//
//	result := cache.CommunityCache.DelFriendRequest(userId, requestID)
//	if (!result) {
//		exception.GameException(exception.FRIEND_REQUEST_NOFOUND)
//	}
//	db.DatabaseHandler.DeleteOne(&db.DFriendReq{ID: db.CompID{message.GetId(), message.GetRequestID()}})
//
//	response.RefuseFriendRequestRet = &protocol.RefuseFriendRequestRet{Result: true)}
//
//}
//
////删除好友
//type DeleteFriendService struct {
//}
//
//func (service *DeleteFriendService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
//	message := request.GetDeleteFriend()
//	userId := character.Int32ToString(message.GetId())
//	friendID := character.Int32ToString(message.GetFriendID())
//
//	db.DatabaseHandler.DeleteOne(&db.DFriend{ID: db.CompID{message.GetId(), message.GetFriendID()}})
//	db.DatabaseHandler.DeleteOne(&db.DFriend{ID: db.CompID{message.GetFriendID(), message.GetId()}})
//
//	cache.CommunityCache.DelFriend(userId, friendID)
//	cache.CommunityCache.DelFriend(friendID, userId)
//
//	response.DeleteFriendRet = &protocol.DeleteFriendRet{Result: true)}
//}

type FollowService struct {
}

func (service *FollowService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetFollow()
	FollowUser(message.GetId(), message.GetFollowerID())
	//id := character.Int32ToString(message.GetId())
	//followerID := character.Int32ToString(message.GetFollowerID())
	//
	//if id == followerID {
	//	exception.GameException(exception.CANNOT_FOLLOW_MYSELF)
	//}
	//
	//if cache.CommunityCache.GetFollowerCount(id) >= conf.FOLLOW_LIMIT {
	//	exception.GameException(exception.MAX_FOLLOW)
	//}
	//
	////已关注
	//if cache.CommunityCache.ExistFollower(id, followerID) {
	//	exception.GameException(exception.REPEATE_FOLLOW)
	//}
	//
	//dbFollow := &db.DFollow{
	//	ID:      db.CompID{message.GetId(), message.GetFollowerID()},
	//	AddTime: time.Now(),
	//}
	//lpc.DBServiceProxy.Insert(dbFollow, db.DatabaseHandler)
	////db.DatabaseHandler.Insert(dbFollow)
	//cache.CommunityCache.AddFollower(id, followerID)

	response.FollowRet = &protocol.FollowRet{
		Follower: core.BuildFollowerInfo(character.Int32ToString(message.GetId()), character.Int32ToString(message.GetFollowerID()), time.Now().Unix()),
	}
}


func FollowUser(uid int32, followerid int32){
	gameCode := FollowUser1(uid, followerid)
	if gameCode != exception.NONE {
		exception.GameException(gameCode)
	}
}

func FollowUser1(uid int32, followerid int32) exception.GameCode {
	id := character.Int32ToString(uid)
	followerID := character.Int32ToString(followerid)
	if id == followerID {
		return exception.CANNOT_FOLLOW_MYSELF
	}

	if cache.CommunityCache.GetFollowerCount(id) >= conf.FOLLOW_LIMIT {
		return exception.MAX_FOLLOW
	}

	//已关注
	if cache.CommunityCache.ExistFollower(id, followerID) {
		return exception.REPEATE_FOLLOW
	}

	dbFollow := &db.DFollow{
		ID:      db.CompID{SubID1:uid, SubID2:followerid},
		AddTime: time.Now(),
	}
	lpc.DBServiceProxy.Insert(dbFollow, db.DatabaseHandler)
	//db.DatabaseHandler.Insert(dbFollow)
	cache.CommunityCache.AddFollower(id, followerID)
	return exception.NONE
}


type UnFollowService struct {
}

func (service *UnFollowService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetUnfollow()
	id := character.Int32ToString(message.GetId())
	followerID := character.Int32ToString(message.GetUnfollowerID())
	//lpc.DBServiceProxy.Delete(&db.DFriend{ID: db.CompID{message.GetId(), message.GetUnfollowerID()}}, db.DatabaseHandler)
	cache.CommunityCache.DelFollower(id, followerID)
	lpc.DBServiceProxy.Delete(&db.DFollow{ID: db.CompID{message.GetId(), message.GetUnfollowerID()}}, db.DatabaseHandler)
	//清除朋友圈消息
	cache.CommunityCache.RemoveUserReceive(id, followerID)

	response.UnfollowRet = &protocol.UnfollowRet{Result: true}
}

type FollowerListService struct {
}

func (service *FollowerListService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	response.GetFollowerListRet = core.BuildFollowerList(character.Int32ToString(request.GetGetFollowerList().GetId()))
}

type FollowingListService struct {
}

func (service *FollowingListService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	response.GetFollowingListRet = core.BuildFollowingList(character.Int32ToString(request.GetGetFollowingList().GetId()))
}

type FollowEachService struct {

}

func (service *FollowEachService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	id := request.Follow.GetId()
	followerID := request.Follow.GetFollowerID()
	FollowUser1(id, followerID)
	FollowUser1(followerID,id)
}

type FollowStateService struct {

}

func (service *FollowStateService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGetFollowState()
	uid1 := character.Int32ToString(message.GetUid1())
	uid2 := character.Int32ToString(message.GetUid2())

	follower := cache.CommunityCache.ExistFollower(uid1, uid2)
	following := cache.CommunityCache.ExistFollowing(uid1, uid2)

	response.GetFollowStateRet = &protocol.GetFollowStateRet{
		Follower:follower,
		Following:following,
	}
}

