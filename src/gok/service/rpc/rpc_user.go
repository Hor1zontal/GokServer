/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/7/18
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package rpc

import (
	"gok/cache"
	"gok/constant"
	clustercache "gok/module/cluster/cache"
	"gok/service"
	"gok/service/msg/protocol"
)

//用户服务代理句柄,后面采取RPC方式调用
var UserServiceProxy = &userHandler{rpcHandler:&rpcHandler{serviceType:service.SERVICE_USER_RPC}, ClusterCacheManager:clustercache.Cluster}

type userHandler struct {
	*rpcHandler
	*cache.ClusterCacheManager
}

//发送指定用户节点处理用户消息 会持久化调用 消息没收到可以持久化
func (this *userHandler) PersistCall(uid int32, message *protocol.C2GS) {
	if uid <= 0 {
		return
	}
	message.Param = uid
	node := this.GetUserNode(uid)
	this.AsyncHandlePriorityRemoteMessage(message, node)
}

//T人
func (this *userHandler) KickOut(uid int32, node string, kickType constant.LOGOUT_TYPE) {
	if uid <= 0 {
		return
	}
	request := &protocol.C2GS{
		Sequence: []int32{500},
		Kickoff:  &protocol.KickOff{
			Uid:uid,
			KickType:int32(kickType),
		},
	}
	//异步T人
	this.AsyncHandleNodeMessage(request, node)
}

//转移会话数据
func (this *userHandler) TransferUserSession(uid int32, node string) []byte {
	if uid <= 0 {
		return nil
	}
	request := &protocol.C2GS{
		Sequence: []int32{507},
		TransferRemoteUserSession:  &protocol.TransferRemoteUserSession{
			Uid:uid,
		},
	}
	//阻塞转移会话数据
	result := this.HandleNodeMessage(request, node).GetTransferRemoteUserSessionRet()
	if result == nil {
		return nil
	}
	return result.GetSession()
}

//给指定用户推送消息
func (this *userHandler) Push(uid int32, message *protocol.GS2C) {
	if uid <= 0 {
		return
	}
	node := this.GetUserOnlineNode(uid)
	if node == "" {
		return
	}
	pushMessage := &protocol.C2GS{
		Sequence:[]int32{1000},
		UserPush: &protocol.UserPush{
			Uid:uid,
			Message:message,
		},
	}
	this.AsyncHandleNodeMessage(pushMessage, node)
}

//给所有用户广播消息
func (this *userHandler) BroadcastAll(message *protocol.GS2C) {
	pushMessage := &protocol.C2GS{
		Sequence:[]int32{1000},
		UserPush: &protocol.UserPush{
			Uid:-1,
			Message:message,
		},
	}
	this.AsyncBroadcastAllMessage(pushMessage)
}

//func (this *userHandler) DrawDayGift(uid int32) {
//	if uid <= 0 {
//		return
//	}
//	node := this.GetUserOnlineNode(uid)
//	if node == "" {
//		return
//	}
//	request := &protocol.C2GS{
//		Sequence: []int32{560},
//		DrawDayGift: &protocol.DrawDayGift{
//			Uid:uid,
//		},
//	}
//	this.AsyncHandleNodeMessage(request, node)
//	//this.PersistCall(uid, request)
//}
//
//func (this *userHandler) UpdateStarUnlockFlag(uid int32, request *protocol.C2GS) {
//	if uid <= 0 {
//		return
//	}
//	node := this.GetUserOnlineNode(uid)
//	if node == "" {
//		return
//	}
//	this.AsyncHandleNodeMessage(request, node)
//}

func (this *userHandler) UserHandleMessage(uid int32, message *protocol.C2GS) {
	if uid <= 0 {
		return
	}
	node := this.GetUserOnlineNode(uid)
	if node == "" {
		return
	}
	this.AsyncHandleNodeMessage(message, node)
}

func (this *userHandler) QueryByNickname(nickname string) *protocol.QueryByNicknameRet{
	request := &protocol.C2GS{
		Sequence:[]int32{515},
		QueryByNickname:&protocol.QueryByNickname{
			Nickname:nickname,
		},
	}
	return this.HandleMessage(request).GetQueryByNicknameRet()
}