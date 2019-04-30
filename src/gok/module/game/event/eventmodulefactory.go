/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/7/13
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package event

import (
	"gok/module/game/db"
	"gok/service/exception"
	"time"
	"encoding/json"
	"gok/module/game/conf"
)

//新建模块handler
func createModuleHandler(moduleType conf.ModuleEnum) db.ModuleHandler {
	switch moduleType {
	case conf.MODULE_RANDOM_TARGET:
		return &RandomTargetModule{}
	case conf.MODULE_GET_FAITH:
		return &GetFaithModule{}
	case conf.MODULE_GET_BELIEVER:
		return &GetBelieverModule{}
	case conf.MODULE_LOOT_FAITH:
		return &LootFaithModule{}
	case conf.MODULE_ATTACK_BUILD:
		return &AttackBuildingModule{}
	case conf.MODULE_LOOT_BELIEVER:
		return &LootBelieverModule{}
	//case conf.MODULE_RECRUIT:
	//	return &RecruitModule{}
	//case conf.MODULE_BUFF:
	//	return &BuffModule{}
	case conf.MODULE_STORYBOARD:
		return &StoryboardModule{}
	//case conf.MODULE_SWITCH_STAR:
	//	return &SwitchStarModule{}
	//case conf.MODULE_ATTACK:
	//	return &AttackModule{}
	//case conf.MODULE_VOTE:
	//	return &VoteModule{}
	//case conf.MODULE_CAPTURE:
	//	return &CaptureModule{}
	case conf.MODULE_CUSTOM_DATA:
		return &CustomDataModule{}
	}
	return nil
}

func NewEventSession(eventID int32, eventType int32, uid int32, nickname string, guide bool) *EventSession {
	module, _ := createModule(eventType, 1)
	if module == nil {
		exception.GameException(exception.EVENT_CONFIG_EXCEPTION)
	}
	stepModules := []*db.EventModule{}
	stepModules = append(stepModules, module)

	//subEvents := []*db.SubEvent{}
	//subEvents = append(subEvents, &db.SubEvent{Type:eventType, Step:1})

	event := &db.DBEvent{
		ID : eventID,
		Type : eventType,
		Guide: guide,
		//SubEvents:subEvents,
		Caller: &db.DBEventMember{
			Uid:uid,
			Nickname:nickname,
		},
		StepModules:stepModules,
		DisplaySUid:uid,
		CreateTime:time.Now(),
	}
	//eventID, err := db.DatabaseHandler.GenId(event)
	//if err != nil {
	//	exception.GameException(exception.DATABASE_EXCEPTION)
	//}
	//event.ID = eventID
	//lpc.DBServiceProxy.Insert(event, db.DatabaseHandler)

	//db.DatabaseHandler.Insert(event)

	eventSession := &EventSession{
		DBEvent:event,
	}
	if module.Handler != nil {
		module.Handler.Start(eventSession)
	}
	return eventSession
}

//create 从数据库对象中初始化事件对象
func initEventSession(eventData *db.DBEvent) *EventSession {
	eventSession := &EventSession{
		DBEvent:eventData,
	}

	////兼容之前的数据
	//if (len(eventData.SubEvents) == 0) {
	//	subEvents := []*db.SubEvent{}
	//	subEvents = append(subEvents, &db.SubEvent{Type:eventData.Type, Step:1})
	//	eventData.SubEvents = subEvents
	//}

	stepData := conf.GetEventStepConfig(eventSession.GetType(), eventSession.getCurrentStep())
	currModule := eventSession.GetCurrentStepModule()
	if (currModule != nil) {
		handler := createModuleHandler(currModule.Type)
		if (handler != nil) {
			data, err := json.Marshal(currModule.PersistentData)
			if (err == nil) {
				handler.Init(data, stepData)
			}
		}
		currModule.Handler = handler
	}
	return eventSession
}


func createModule(eventType int32, step int32) (*db.EventModule, *conf.StepData)  {
	stepData := conf.GetEventStepConfig(eventType, step)
	if (stepData == nil) {
		return nil, nil
	}
	startTimestamp := time.Now()
	endTimestamp := time.Time{}
	moduleType := stepData.Module
	if (stepData.TimeLimit != 0) {
		endTimestamp = startTimestamp.Add(time.Duration(stepData.TimeLimit) * time.Second)
	}
	handler := createModuleHandler(conf.ModuleEnum(moduleType))
	if (handler != nil) {
		handler.Init(nil, stepData)
	}

	return &db.EventModule{
		Type:conf.ModuleEnum(moduleType),
		StartTimestamp:startTimestamp,
		EndTimestamp:endTimestamp,
		Handler:handler,
		//PersistentData:persistentData,
	}, stepData
}



