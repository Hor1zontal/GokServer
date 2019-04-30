/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/4/24
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package global

import (
	"gok/service/msg/protocol"
	"github.com/name5566/leaf/module"
)

var handler *module.Skeleton = nil

var tempGlobalMessage = []*protocol.GlobalMessage{}

func Init(skeleton *module.Skeleton) {
	handler = skeleton
}

const (
	TEMP_GLOBALMESSAGE_LIMIT = 10
)

//const (
//	MSG_BROADCAST = "b_msg"
//	MSG_FOLLOWING = "f_msg"
//)
//
////广播所有玩家
//func broadcast(args []interface{}) {
//	message, ok := args[0].(*protocol.GlobalMessage)
//	if ok {
//		if len(tempGlobalMessage) == TEMP_GLOBALMESSAGE_LIMIT {
//			tempGlobalMessage = tempGlobalMessage[3:TEMP_GLOBALMESSAGE_LIMIT]
//		}
//		tempGlobalMessage = append(tempGlobalMessage, message)
//		rpc.UserServiceProxy.BroadcastAll(&protocol.GS2C{Sequence:[]int32{1070},GlobalMessagePush:message})
//	}
//}
//

//广播世界消息
func BroadcastMessage(message *protocol.GlobalMessage){
	//Handler.ChanRPCServer.Go(MSG_BROADCAST, message)
}

//推送所有关注消
func PushFollowings(uid int32, message *protocol.GS2C) {
	if handler != nil {
		handler.Go((&Push{Uid:uid, Msg:message}).PushFollowings, nil)
	}
}

func PersistCallFollowings(uid int32, call *protocol.C2GS) {
	if handler != nil {
		handler.Go((&Push{Uid:uid, Call:call}).PersistCallFollowings, nil)
	}
}


//获取最新的世界消息
func GetGlobalMessage() *protocol.GetGlobalMessageRet {
	return &protocol.GetGlobalMessageRet{
		Message:tempGlobalMessage,
	}
}
