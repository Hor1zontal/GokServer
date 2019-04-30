/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/4/17
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package rpc

import (
	"gok/service"
	"gok/service/msg/protocol"
	"gok/module/cluster/cache"
)

var PassportServiceProxy = &PassportService{&rpcHandler{serviceType:service.SERVICE_PASSPORT_RPC}}

type PassportService struct {
	*rpcHandler
}

func (this *PassportService) ChannelLogin(channel string, channelUID string, openID string, nick string, avatar string) *protocol.ChannelLoginRet {
	request := &protocol.C2GS{
		Sequence:[]int32{650},
		ChannelLogin:&protocol.ChannelLogin{
			Channel:channel,
			ChannelUID:channelUID,
			OpenID:openID,
			Nickname:nick,
			Avatar:avatar,
		},
	}
	response := this.HandleMessage(request)
	return response.GetChannelLoginRet()
}

func (this *PassportService) UserState(uid int32, state int32) {
	request := &protocol.C2GS{
		Sequence:[]int32{651},
		ChangeUserState:&protocol.ChangeUserState{
			Uid:uid,
			State:state,
		},
	}
	this.HandleMessage(request)
}

func (this *PassportService) RefreshClientVersion() *protocol.RefreshClientVersionRet{
	request := &protocol.C2GS{
		Sequence: []int32{652},
		RefreshClientVersion: &protocol.RefreshClientVersion{
		},
	}
	return this.HandleMessage(request).GetRefreshClientVersionRet()
}


func (this *PassportService) WechatEventPush(uid int32, event string, delay int32) {
	//用户在线不需要推送微信公众号消息
	if cache.Cluster.IsUserOnline(uid) {
		return
	}

	request := &protocol.C2GS{
		Sequence: []int32{653},
		WechatEventPush: &protocol.WechatEventPush{
			Event:event,
			Uid:uid,
			Delay:delay,
		},
	}
	this.AsyncHandleMessage(request)
}

func (this *PassportService) CleanTestAccount(uid int32) *protocol.CleanTestAccountRet{
	request := &protocol.C2GS{
		Sequence: []int32{654},
		CleanTestAccount:&protocol.CleanTestAccount{
			Uid:uid,
		},
	}
	return this.HandleMessage(request).GetCleanTestAccountRet()
}

func (this *PassportService) QueryByUsername(username string) *protocol.QueryByUsernameRet {
	request := &protocol.C2GS{
		Sequence:[]int32{655},
		QueryByUsername:&protocol.QueryByUsername{
			Username:username,
		},
	}
	return this.HandleMessage(request).GetQueryByUsernameRet()
}