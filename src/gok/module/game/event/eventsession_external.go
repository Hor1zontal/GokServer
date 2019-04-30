/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/5/16
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package event

import (
	"gok/service/msg/protocol"
	"gok/module/game/db"
	"time"
	"gok/module/game/conf"
)

//---------------------------------扩展接口，供外部模块调用,需要看情况加锁-----------------

//向所有人广播事件消息
//ignoreOutEvent 是否不发送消息给未进入事件界面的玩家
func (event *EventSession) BroadcastAll(message *protocol.GS2C) {
	//向发起方发送消息
	if event.Caller != nil {
		sendEventMessage(event.ID, event.Caller.Uid, message)
	}
	//向选择的目标发起消息
	//target := getSelectTarget(event)
	//if target != nil {
	//	sendEventMessage(event.ID, target.Uid, message)
	//}
	//members := getRecruitMembers(event)
	//if (members != nil) {
	//	for _, member := range members {
	//		sendEventMessage(event.ID, member.Uid, message)
	//	}
	//}
}

//获取指定类型的模块
func (this *EventSession) GetModule(moduleType conf.ModuleEnum) *db.EventModule{
	return this.getModule(moduleType)
}

//处理模块消息
func (this *EventSession) HandleModuleMessage(message *protocol.C2GS, response *protocol.GS2C) {
	module := this.GetCurrentStepModule()
	module.HandleMessage(message, response, this)
}

//func (this *EventSession) IntoEvent(uid int32) {
//	//if (!this.isUserInEvent(uid)) {
//	//	exception.GameException(exception.INVALID_EVENT_AUTHORITY)
//	//}
//	cache.UserCache.InEvent(uid, this.ID)
//}
//
//func (this *EventSession) OutEvent(uid int32) {
//	//if (!this.isUserInEvent(uid)) {
//	//	exception.GameException(exception.INVALID_EVENT_AUTHORITY)
//	//}
//	cache.UserCache.OutEvent(uid, this.ID)
//}
//
////时间完成，需要推送事件状态，和通知任务模块
//func (this *EventSession) Done() {
//	this.Lock()
//	defer this.Unlock()
//	SendEventDone(this)
//	EventManager.DeleteEvent(this.ID)
//}

func (this *EventSession) DealTimeChange(time time.Time) {
	module := this.GetCurrentStepModule()
	//handler := module.Handler
	//if (handler != nil) {
	//	timerHandler, ok := handler.(db.ModuleTimberHandler)
	//	if (ok) {
	//		timerHandler.DealTimer(time, this)
	//	}
	//}
	//定时时间到
	if !module.EndTimestamp.IsZero() && module.EndTimestamp.Before(time) {
		module.HandleTimesUp(this)
		//this.NextStep()
	}
}