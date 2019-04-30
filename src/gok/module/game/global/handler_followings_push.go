/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/10/15
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package global

import (
	"gok/service/msg/protocol"
	"gok/service/rpc"
)

type Push struct {
	Uid int32
	Msg *protocol.GS2C
	Call *protocol.C2GS
}

func (this *Push) PushFollowings() {
	followings := rpc.CommunityServiceProxy.GetFollowingList(this.Uid).GetFollowings()
	if followings != nil && len(followings) > 0 {
		for _, following := range followings {
			rpc.UserServiceProxy.Push(following.GetId(), this.Msg)
		}
	}
}

func (this *Push) PersistCallFollowings() {
	followings := rpc.CommunityServiceProxy.GetFollowingList(this.Uid).GetFollowings()
	if followings != nil && len(followings) > 0 {
		for _, following := range followings {
			rpc.UserServiceProxy.PersistCall(following.GetId(), this.Call)
		}
	}
}
