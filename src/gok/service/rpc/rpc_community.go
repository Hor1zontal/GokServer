/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/5/27
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package rpc

import (
	"gok/service"
	"gok/service/msg/protocol"
	"gok/constant"
)

var CommunityServiceProxy = &communityHandler{&rpcHandler{serviceType:service.SERVICE_COMMUNITY_RPC}}

type communityHandler struct {
	*rpcHandler
}

//发布朋友圈消息
func (this *communityHandler) PublicMoments(uid int32, momentType constant.MomentsType, refID int32) *protocol.PublicMomentRet {
	request := &protocol.C2GS{
		Sequence: []int32{243},
		PublicMoment: &protocol.PublicMoment{
			Uid: uid,
			Type: int32(momentType),
			RefID: refID,
		},
	}
	return this.HandleMessage(request).GetPublicMomentRet()
}


func (this *communityHandler) GetFollowingList(uid int32) *protocol.GetFollowingListRet {
	request := &protocol.C2GS{
		Sequence: []int32{224},
		GetFollowingList: &protocol.GetFollowingList{
			Id:uid,
		},
	}
	return this.HandleMessage(request).GetFollowingListRet
}


func (this *communityHandler) RemoveMoments(uid int32) {
	request := &protocol.C2GS{
		Sequence: []int32{244},
		RemoveMoments: &protocol.RemoveMoments{
			Uid:uid,
			SaleID:uid,
		},
	}
	this.HandleMessage(request)
}

func (this *communityHandler) FollowEach(uid int32, followID int32) {
	request := &protocol.C2GS{
		Sequence: []int32{225},
		Follow: &protocol.Follow{
			Id:uid,
			FollowerID:followID,
		},
	}
	this.HandleMessage(request)
}

func (this *communityHandler) IsEachFollow(uid1 int32, uid2 int32) bool {
	request := &protocol.C2GS{
		Sequence: []int32{226},
		GetFollowState: &protocol.GetFollowState{
			Uid1:uid1,
			Uid2:uid2,
		},
	}
	result := this.HandleMessage(request).GetGetFollowStateRet()
	return result.GetFollower() && result.GetFollowing()
}