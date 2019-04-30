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

//import (
//	"gok/service/msg/protocol"
//	"gok/service"
//	"gok/service/exception"
//)
//
//var EventServiceProxy = &eventHandler{&rpcHandler{serviceType:service.SERVICE_EVENT_RPC}}
//
//
//type eventHandler struct {
//	*rpcHandler
//}
//
//
////随机事件星球
//func (this *eventHandler) RandomEventTarget(request *protocol.C2GS) *protocol.RandomTargetRet {
//	return this.HandleMessage(request).GetRandomTargetRet()
//}
//
//
//func (this *eventHandler) RemoveEvent(eventID int32, uid int32) {
//	message := &protocol.C2GS{
//		Sequence: []int32{33},
//		RemoveEvent: &protocol.RemoveEvent{
//			EventID: eventID,
//			Uid:     uid,
//		},
//	}
//	this.HandleMessage(message)
//}
//
//func (this *eventHandler) GenEvent(eventType int32, uid int32, nickname string, guide bool) int32 {
//	message := &protocol.C2GS{
//		Sequence: []int32{501},
//		GenEvent: &protocol.GenEvent{
//			EventType: eventType,
//			Uid:       uid,
//			Nickname:  nickname,
//			Guide:     guide,
//		},
//	}
//	result := this.HandleMessage(message)
//	if result.GetGenEventRet() == nil {
//		exception.GameException(exception.TASK_BASE_NOTFOUND)
//	}
//	return result.GetGenEventRet().GetEvent().GetEventID()
//}