/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/7/14
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package event

import (
	"gok/service/msg/protocol"
	"gok/module/game/db"
)


func (this *EventSession) BuildEventProtocol() *protocol.Event {
	event := &protocol.Event{
		EventID: this.GetID(),
		BaseID: this.GetType(),
		Step: this.getCurrentStep(),
		DisplayStarUid:this.DisplaySUid,
	}
	if (this.Caller != nil) {
		event.CallMember = BuildEventMember(this.Caller)
	}
	stepModule := this.GetCurrentStepModule()
	if (stepModule != nil) {
		event.StepModule = BuildModule(stepModule, this.GetID())
	}
	return event
}

func BuildModule(module *db.EventModule, eventID int32) *protocol.Event_Module {
	if (module == nil) {
		return nil
	}
	moduleData := &protocol.Event_Module{
		ModuleID:int32(module.Type),
		StartTimestamp:module.StartTimestamp.Unix(),
		EndTimestamp:module.EndTimestamp.Unix(),

	}

	handler := module.Handler
	if (handler == nil) {
		return moduleData
	}
	timerHandler, ok := handler.(db.ModuleDataHandler)
	if (ok) {
		moduleData.Data = timerHandler.GetData(eventID)
	} else {
		jsonData := module.GetJsonData()
		if (jsonData != nil) {
			moduleData.Data = string(jsonData)
		}
	}
	return moduleData
}

func BuildEventMember(dbEventMember *db.DBEventMember) *protocol.EventMember {
	return &protocol.EventMember{
		Uid:dbEventMember.Uid,
		Nickname:dbEventMember.Nickname,
	}
}




//func BuildRecruitMember(member *RecruitMember) *protocol.EventMember {
//	return &protocol.EventMember{
//		Uid:member.Uid),
//		Nickname:member.Nickname),
//		Msg:member.Msg),
//	}
//}

//func BuildEventMembers(dbEventMember []*db.DBEventMember) []*protocol.EventMember {
//	result := []*protocol.EventMember{}
//	for _, dbMember := range dbEventMember {
//		result = append(result, BuildEventMember(dbMember))
//	}
//	return result
//}
//
//
//func BuildEventFields(dbEventFields []*db.EventField) []*protocol.EventField {
//	result := []*protocol.EventField{}
//	for _, dbEvent := range dbEventFields {
//		result = append(result, BuildEventField(dbEvent))
//	}
//	return result
//}
//

//func BuildEventField(dbEventField *BuffField) *protocol.EventField {
//	return &protocol.EventField{
//		Name:dbEventField.Name),
//		Value:dbEventField.Value),
//	}
//}


func BuildDisplayUidPush(uid int32) *protocol.GS2C {
	return &protocol.GS2C{
		Sequence:[]int32{1013},
		DisplayStarUidPush:&protocol.DisplayStarUidPush{
			DisplayStarUid:uid,
		},
	}
}

func BuildStepPush(eventID int32, eventType int32, step int32, module *db.EventModule) *protocol.GS2C {
	return &protocol.GS2C{
		Sequence:[]int32{1012},
		EventStepPush:&protocol.EventStepPush{
			EventID:eventID,
			BaseID:eventType,
			Step:step,
			StepModule:BuildModule(module, eventID),
		},
	}
}


//func BuildStepEndPush(eventID int32) *protocol.GS2C {
//	return &protocol.GS2C{
//		Sequence:[]int32{1012},
//		EventStepPush:&protocol.EventStepPush{
//			EventID:eventID),
//			Step:0),
//		},
//	}
//}
//
//func BuildVoteEndPush(eventID int32, no int, voteNum int) *protocol.GS2C {
//	return &protocol.GS2C{
//		Sequence:[]int32{1015},
//		VoteFinishPush:&protocol.VoteFinishPush{
//			EventID:eventID),
//			Result:&protocol.Vote{
//				No:int32(no)),
//				VoteNum:int32(voteNum)),
//			},
//		},
//	}
//}


