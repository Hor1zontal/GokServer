/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/5/11
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package event

import (
	"gok/module/game/db"
	"gok/service/msg/protocol"
	"gok/service/rpc"
	"gok/module/game/conf"
)

type UserHandler interface {
	SyncMessage(arg interface{})
	AppendUserEventStatistic(arg interface{}, arg1 interface{})
	//GetUserSession(uid int32) *user.Session
}

type EventSession struct {
	*db.DBEvent
	handler UserHandler
}



//ignoreOutEvent是否只发送给进入事件的玩家
func sendEventMessage(eventId int32, uid int32, message *protocol.GS2C) {
	if (uid == 0) {
		return
	}
	//不发送给未进入事件界面的玩家消息
	//if (!cache.UserCache.IsInEvent(uid, eventId)) {
	//	return
	//}
	rpc.UserServiceProxy.Push(uid, message)
}

func (this *EventSession) SetHandler(handler UserHandler) {
	this.handler = handler
}

func (this *EventSession) AppendStatisticData(eventType int32, refNum int32) {
	this.handler.AppendUserEventStatistic(eventType, refNum)
}

func (this *EventSession) getCurrentSubEvent(moduleType conf.ModuleEnum) *db.EventModule{
	if (this.StepModules == nil) {
		return nil
	}
	for _, module := range this.StepModules {
		if (module.Type == moduleType) {
			return module
		}
	}
	return nil
}


//获取指定类型的模块
func (this *EventSession) getModule(moduleType conf.ModuleEnum) *db.EventModule{
	if (this.StepModules == nil) {
		return nil
	}
	for _, module := range this.StepModules {
		if (module.Type == moduleType) {
			return module
		}
	}
	return nil
}

//获取当前步骤编号
func (this *EventSession) getCurrentStep() int32{
	if (this.StepModules == nil) {
		return 0
	}
	return int32(len(this.StepModules))
}

//获取当前步骤模块
func (this *EventSession) GetCurrentStepModule() *db.EventModule{
	if (this.StepModules == nil) {
		return nil
	}
	return this.StepModules[len(this.StepModules) - 1]
}

func (this *EventSession) GetModuleHandler(moduleType conf.ModuleEnum) db.ModuleHandler {
	module := this.getModule(moduleType)
	if (module == nil) {
		return nil
	}
	return module.Handler
}

//获取json格式数据
//获取json格式数据
func (this *EventSession) SetDirty() {
	//db.UpdateHandler.UpdateQueue(this.DBStar)
	//db.DatabaseHandler.UpdateOne(this.DBEvent)
	//db.UpdateHandler.UpdateQueue(database.OP_UPDATE, this.DBEvent)
	//db.DatabaseHandler.UpdateOne(this.DBEvent)
	//lpc.DBServiceProxy.Opt(this.DBEvent, db.DatabaseHandler)
}

//更新同步持久化数据
func (this *EventSession) SyncPersistentData() {
	for _, module := range this.DBEvent.StepModules {
		module.SyncPersistentData()
	}
}


func (this *EventSession) GetCaller() *db.DBEventMember {
	return this.Caller
}

func (this *EventSession) GetID() int32 {
	return this.ID
}

//获取事件类型
func (this *EventSession) GetType() int32 {
	return this.Type
}

func (this *EventSession) SetDisplaySUid(displaySUid int32) {
	this.DisplaySUid = displaySUid
}


//切换到下一个子事件
func (this *EventSession) AddNextEvent(subEventType int32) bool {
	module, stepData := createModule(subEventType, 1)
	if (module == nil || stepData == nil) {
		SendEventDone(this)
		return false
	}

	//把当前事件数据存入历史记录
	this.SyncPersistentData()
	this.History = append(this.History, &db.SubEvent{
		Type:this.GetType(),
		StepModules:this.StepModules,
	})

	stepModules := []*db.EventModule{}
	stepModules = append(stepModules, module)
	this.StepModules = stepModules
	this.Type = subEventType
	this.SetDirty()
	message := BuildStepPush(this.GetID(), this.GetType(), 1, module)
	this.BroadcastAll(message)
	return true
}


//切换到下个步骤
//TODO 后续要考虑配置导致的永远不会结束的事件需要清除
func (this *EventSession) NextStep() bool {
	newStep := this.getCurrentStep() + 1
	module, stepData := createModule(this.GetType(), newStep)
	//下一步拿不到数据了，通知事件结束
	if (module == nil || stepData == nil) {
		SendEventDone(this)
		return false
	}

	this.StepModules = append(this.StepModules, module)
	if (module.Handler != nil) {
		module.Handler.Start(this)
	}
	message := BuildStepPush(this.GetID(), this.GetType(), this.getCurrentStep(), this.GetCurrentStepModule())
	this.BroadcastAll(message)

	//处理事件触发任务
	//taskData := conf.GetTriggerTaskData(this.GetType(), newStep)
	//if (taskData != nil) {
	//	if (taskData.TriggerTarget == SELF) {
	//		SendTaskTrigger(this.GetCaller().Uid, this.GetID(), taskData.ID)
	//	} else if (taskData.TriggerTarget == TARGET) {
	//		target := getSelectTarget(this)
	//		if (target != nil) {
	//			SendTaskTrigger(target.Uid, this.GetID(), taskData.ID)
	//		}
	//	}
	//}
	this.SyncPersistentData()
	this.SetDirty()
	return true
}
