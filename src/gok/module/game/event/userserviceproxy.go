/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package event

import (
	"gok/service/msg/protocol"
	"gok/module/game/db"
)

//func SendDeleteEventAssist(uid int32, eventId int32) {
//	//向用户模块发送删除协助消息
//	message := &protocol.C2GS{
//		Sequence:[]int32{504},
//		DeleteEventAssist:eventId,
//	}
//	rpc.UserServiceProxy.PersistCall(uid, cache.ClusterCache.GetUserNode(uid), message)
//}

//发送触发任务
//func SendTaskTrigger(uid int32, eventID int32, taskType int32) {
//	message := &protocol.C2GS{
//		Sequence:[]int32{503},
//		TriggerTask:&protocol.TriggerTask{
//			EventID:eventID,
//			TaskType:taskType,
//		},
//	}
//	rpc.UserServiceProxy.PersistCall(uid, cache.ClusterCache.GetUserNode(uid), message)
//}

func SendEventDone(event *EventSession) {
	reward := &protocol.Reward{Believer:[]*protocol.BelieverInfo{}}
	for _, module := range event.StepModules {
		handler := module.Handler
		if (handler != nil) {
			rewardHandler, ok := handler.(db.ModuleRewardHandler)
			if (ok) {
				rewardHandler.AppendReward(reward)
			}
		}
	}


	doneMessage := &protocol.EventDone{
		EventID:event.ID,
		Reward:reward,
		Guide:event.Guide,
	}
	target := getSelectTarget(event)
	if target != nil {
		doneMessage.TargetID = target.Uid
	}

	message := &protocol.C2GS{
		Sequence:[]int32{505},
		EventDone:doneMessage,
	}
	event.handler.SyncMessage(message)
	//userSession := event.handler.GetUserSession(event.Caller.Uid)
	//service.ServiceManager.HandleChannelMessage(message, userSession)
	//rpc.UserServiceProxy.PersistCall(event.Caller.Uid, message)

	//target := getSelectTarget(event)
	//if (target != nil) {
	//	rpc.UserServiceProxy.PersistCall(target.Uid, cache.ClusterCache.GetUserNode(target.Uid),message)
	//}
}

//发送协助请求
//func SendAssistRequest(uid int32, msg string, event db.EventContext){
//	message := &protocol.C2GS{
//		Sequence:[]int32{502},
//		AssistEventRequest:&protocol.AssistEventRequest{
//			EventID:event.GetID()),
//			Uid:event.GetCaller().Uid),
//			Nickname:event.GetCaller().Nickname),
//			Msg:msg),
//		},
//	}
//	//向用发送消息
//	rpc.UserServiceProxy.PersistCall(uid, message)
//}

//func PublicMoments(uid int32, momentType int32, refID int32, title string) {
//	message := &protocol.C2GS{
//		Sequence:[]int32{243},
//		PublicMoment:&protocol.PublicMoment{
//			Uid:uid),
//			Type:momentType),
//			RefID:refID),
//			Title:title),
//		},
//	}
//	rpc.UserServiceProxy.PersistCall(uid, message)
//	//return this.HandleMessage(message).GetPublicEventRet()
//}