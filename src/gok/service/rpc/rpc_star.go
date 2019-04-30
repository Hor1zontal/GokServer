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
	"gok/service/msg/protocol"
	"gok/service"
	"gok/module/cluster/center"
	"aliens/common/character"
	"gok/service/exception"
)

var StarServiceProxy = &starHandler{&rpcHandler{serviceType:service.SERVICE_STAR_RPC}}

type starHandler struct {
	*rpcHandler
}

const (
	OP_BELIEVER_ADD int32 = 0
	OP_BELIEVER_DEC int32 = 2
)

var INVALID_SERVICE_RESPONSE = &protocol.GS2C{
	Sequence: []int32{1000},
	ResultPush: &protocol.ResultPush{
		Result: int32(exception.SERVICE_INVALID),
	},
}

//请求RPC调用接口层
func (this *starHandler) Call(uid int32, request *protocol.C2GS) *protocol.GS2C {
	request.Param = uid
	node := center.StarHashring.GetNode(character.Int32ToString(uid))
	//if constant.DEBUG {
	//	this.AddCall(node)
	//}
	return this.AssertHandleNodeMessage(request, node)
}

//请求RPC调用接口层
func (this *starHandler) AllocUserNewStar(uid int32, starType int32) *protocol.AllocNewStarRet {
	request := &protocol.C2GS{
		Sequence: []int32{600},
		AllocNewStar: &protocol.AllocNewStar{
			Uid: uid,
			StarType:starType,
		},
	}
	return this.Call(uid, request).GetAllocNewStarRet()
}

func (this *starHandler) GetStarInfoDetail(uid int32, starID int32,  shieldType int32, isConvert bool, allBuildingLevel int32, allBelieverLevel int32) *protocol.GetStarInfoDetailRet {
	request := &protocol.C2GS{
		Sequence: []int32{302},
		GetStarInfoDetail: &protocol.GetStarInfoDetail{
			Uid:uid,
			StarID:starID,
			ShieldType: shieldType,
			IsConvert: isConvert,
			BuildingTotalLevel: allBuildingLevel,
			BelieverTotalLevel: allBelieverLevel,
		},
	}
	return this.Call(uid, request).GetGetStarInfoDetailRet()
}

func (this *starHandler) UpdateStarStatistics(uid int32, statisticsID int32, change float64, param int32) *protocol.UpdateStarStatisticsRet {
	request := &protocol.C2GS{
		Sequence: []int32{370},
		UpdateStarStatistics: &protocol.UpdateStarStatistics{
			Id:     statisticsID,
			Change: change,
			Param: param,
		},
	}
	return this.Call(uid, request).GetUpdateStarStatisticsRet()
}

func (this *starHandler) UpdateBelieverInfo(uid int32, updateInfo []*protocol.BelieverInfo, operation int32, isConvert bool) *protocol.GS2C {
	request := &protocol.C2GS{
		Sequence: []int32{603},
		UpdateBeliever: &protocol.UpdateBeliever{
			Uid:        uid,
			UpdateInfo: updateInfo,
			Operation:  operation,
			IsConvert: isConvert,
		},
	}
	response := this.Call(uid, request)
	if response == nil || response.GetUpdateBelieverRet() == nil {
		return nil
	}
	updateBeliever := response.GetUpdateBelieverRet().GetBeliever()
	if isConvert {
		return response
	}
	if updateBeliever != nil && len(updateBeliever) > 0 {
		return &protocol.GS2C{
			Sequence: []int32{1040},
			BelieverPush: &protocol.BelieverPush{
				Believer: updateBeliever,
			},
		}
	}
	return nil
}

func (this *starHandler) AddCivilization(uid int32, value int32) *protocol.CivilizationInfo {
	return this.Call(uid, &protocol.C2GS{
		Sequence:      []int32{368},
		AddCivilization: &protocol.AddCivilization{CivilizationValue: value},
	}).GetAddCivilizationRet()
}

//func (this *starHandler) GetUserStarInfo(uid int32) *protocol.UserStarInfoRet { //获取用户的当前开发星球信息
//	request := &protocol.C2GS{
//		Sequence: []int32{602},
//		UserStarInfo: &protocol.UserStarInfo{},
//	}
//	return this.PersistCall(uid, request).GetStarInfoRet.GetUserStarInfoRet()
//}

func (this *starHandler) GetLoginStarInfo(uid int32) *protocol.LoginStarInfoRet { //获取用户的法力值上限
	return this.Call(uid, &protocol.C2GS{
		Sequence:      []int32{604},
		LoginStarInfo: &protocol.LoginStarInfo{Uid: uid},
	}).GetLoginStarInfoRet()
}

//func (this *starHandler) SearchStarInfo( uid int32 ) *protocol.SearchStarInfoRet{//探索星球 wjl 20170601
//	request := &protocol.C2GS{
//		Sequence:[]int32{301},
//		SearchStarInfo:&protocol.SearchStarInfo{
//			Uid:uid),//用户ID
//		},
//	}
//	return this.HandleMessage(request).GetSearchStarInfoRet();//消息模块转发到 star模块
//}


func (this *starHandler) RandomEventRobot(eventType int32, level []int32) *protocol.RandomEventRobotRet {
	request := &protocol.C2GS{
		Sequence: []int32{360},
		RandomEventRobot: &protocol.RandomEventRobot{
			EventType:eventType,
			Level:level,
		},
	}
	return this.HandleMessage(request).GetRandomEventRobotRet()
}

func (this *starHandler) RandomGuideRobot(num int32) *protocol.RandomGuideRobotRet {
	request := &protocol.C2GS{
		Sequence: []int32{361},
		RandomGuideRobot: &protocol.RandomGuideRobot{
			Num : num,
		},
	}
	return this.HandleMessage(request).GetRandomGuideRobotRet()
}

//func (this *starHandler) GetRecordStarInfo( starID []int32, userID []int32 )*protocol.GetStarRecordInfoRet{
//
//	request := &protocol.C2GS{
//		Sequence:[]int32{311},
//		GetStarRecordInfo:&protocol.GetStarRecordInfo{
//			Uid:0),
//			starID: starID,
//			UserID: userID,
//		},
//	}
//	return this.HandleMessage(request).GetGetStarRecordInfoRet()
//}

//func (this *starHandler) OccupyStar( starID int32, userID int32 )*protocol.OccupyStarRet{
//
//	request := &protocol.C2GS{
//		Sequence:[]int32{320},
//		OccupyStar:&protocol.OccupyStar{
//			Uid: userID ),
//			starID: starID ),
//		},
//	}
//	return this.HandleMessage(request).GetOccupyStarRet();
//}

//随机事件星球
//func (this *starHandler) RandomStarTarget( request *protocol.C2GS ) *protocol.RandomTargetRet {//随机一个星球的随机目标
//	return this.HandleMessage(request).GetRandomTargetRet()
//}

func (this *starHandler) LootStarBeliever(attackID int32, destID int32, lootBeliever []string) *protocol.LootStarBelieverRet {
	request := &protocol.C2GS{
		Sequence: []int32{270},
		LootStarBeliever: &protocol.LootStarBeliever{
			AttackID:   attackID,
			DestID:     destID,
			BelieverID: lootBeliever,
		},
	}
	return this.Call(destID, request).GetLootStarBelieverRet()
}

func (this *starHandler) TransmitUserStarInfo(uid int32, node string) *protocol.TransmitUserStarRet{

	requset := &protocol.C2GS{
		Sequence: []int32{605},
		TransmitUserStar: &protocol.TransmitUserStar{
			UserID: uid,
		},
	}
	resp := this.HandleNodeMessage(requset, node)

	if resp != nil {
		return resp.GetTransmitUserStarRet()
	}
	return nil
}

func (this *starHandler) HelpRepairBuildPublic(uid int32, buildType int32) *protocol.HelpRepairBuildPublicRet {
	request := &protocol.C2GS{
		Sequence: []int32{606},
		HelpRepairBuildPublic:&protocol.HelpRepairBuildPublic{
			BuildingType:buildType,
		},
	}
	return this.Call(uid, request).GetHelpRepairBuildPublicRet()
}

func (this *starHandler) GetCurrentGroupItems(uid int32) *protocol.GetCurrentGroupItemsRet {
	request := &protocol.C2GS{
		Sequence: []int32{608},
		GetCurrentGroupItems:&protocol.GetCurrentGroupItems{},
	}
	return this.Call(uid, request).GetGetCurrentGroupItemsRet()
}

func (this *starHandler) GetStarFlags(uid int32) *protocol.StarFlagInfoRet {
	request := &protocol.C2GS{
		Sequence: []int32{326},
		StarFlagInfo:&protocol.StarFlagInfo{},
	}
	return this.Call(uid, request).GetStarFlagInfoRet()
}

func (this *starHandler) GetEventRobotTarget(robotID int32) *protocol.GetEventRobotRet {
	request := &protocol.C2GS{
		Sequence: []int32{609},
		GetEventRobot:&protocol.GetEventRobot{Uid:robotID},
	}
	return this.HandleMessage(request).GetGetEventRobotRet()
}

func (this *starHandler) UpdateStarFlag(uid int32, key int32, value int32) *protocol.UpdateStarFlagRet {
	request := &protocol.C2GS{
		Sequence: []int32{327},
		UpdateStarFlag:&protocol.UpdateStarFlag{Flag:key, Value:value},
	}
	return this.Call(uid, request).GetUpdateStarFlagRet()
}

func (this *starHandler) UpdateAllStarFlag(key int32, value int32) *protocol.UpdateAllStarFlagRet {
	request := &protocol.C2GS{
		Sequence: []int32{328},
		UpdateAllStarFlag:&protocol.UpdateAllStarFlag{Key:key, Value:value},
	}
	return this.HandleMessage(request).GetUpdateAllStarFlagRet()
}

func (this *starHandler) GetUidsByCondition(start int64, end int64, buildLv int32, limit int32, skip int32) *protocol.GetOwnersByConditionRet {
	request := &protocol.C2GS{
		Sequence: []int32{610},
		GetOwnersByCondition:&protocol.GetOwnersByCondition{
			Start:start,
			End:end,
			BuildLv:buildLv,
			Limit: limit,
			Skip:skip,
		},
	}
	return this.HandleMessage(request).GetGetOwnersByConditionRet()
}