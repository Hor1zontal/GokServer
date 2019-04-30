package service

import (
	"gok/service/msg/protocol"
	eventModule "gok/module/game/event"
	"gok/service/exception"
	"gok/module/game/user"
	"gok/module/game/conf"
	"gok/constant"
	"gok/service/rpc"
	"aliens/common/character"
)


//翻卡片
type OpenCardService struct {
}

func (service *OpenCardService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	//message := request.GetOpenCard()
	//eventID := message.GetEventID()
	//task := user.EnsureRefTask(eventID)

	event := user.EnsureEvent()

	request.OpenCard.Uid = user.GetID()

	event.HandleModuleMessage(request, response)

	selectCard := response.GetOpenCardRet().GetReward1()
	user.AppendEventCardStatistic(selectCard.GetType())
	if selectCard.GetTarget() != nil {
		user.EventStatisticOpenCard(selectCard.GetTarget())
	}
}



//随机生成目标用户
type RandomTargetService struct {
}

func (service *RandomTargetService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetRandomTarget()
	eventID := message.GetEventID()

	task := user.EnsureRefTask(eventID)

	var randomMutual int32 = 0
	var randomFriend int32 = 0

	//是否复仇任务
	if task.RevengeID != 0 {
		//随机一个复仇对象
		randomMutual = task.RevengeID
		message.Num = 1 //随机一个
		message.MutualID = randomMutual
	} else {
		randomMutual = user.RandomMutualTarget()
		message.MutualID = randomMutual
		message.Num = constant.RANDOM_TARGET_COUNT //随机1个
	}

	if task.RandomCount != 0 {
		user.EnsureDiamond(conf.DATA.AstrolaTargetCost)
	}

	event := user.EnsureEvent()
	event.HandleModuleMessage(request, response)
	//resp := rpc.EventServiceProxy.RandomEventTarget(request)
	resp := response.RandomTargetRet
	if resp.Targets == nil || len(resp.Targets) == 0 {
		return
	}
	for _, target := range resp.Targets {
		//机器人随机名字
		//UpdateTarget(target)
		if randomMutual != 0 && randomMutual == target.GetId() {
			target.Mutual = true
		}
		if randomFriend != 0 && randomFriend == target.GetId() {
			target.Friend = true
		}
	}
	if task.RandomCount != 0 {
		user.TakeOutDiamond(conf.DATA.AstrolaTargetCost, constant.OPT_TYPE_REFRESH_TARGET, 0)
		//user.SetDirty()
		user.WriteMsg(user.BuildRoleSocialPush())
	}
	user.EventStatisticRandomTarget(resp.Targets)
	task.RandomCount += 1
	response.RandomTargetRet = resp
}

type LootFaithService struct {
	//
}

func (service *LootFaithService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求进攻某个建筑物
	//message := request.GetLootFaith()
	//task := user.EnsureRefTask(message.GetEventID())
	//result := rpc.EventServiceProxy.HandleMessage(request)
	mutual := user.IsFlagUnlock(constant.STAR_FLAG_MUTUAL)
	request.LootFaith.Mutual = mutual

	event := user.EnsureEvent()
	event.HandleModuleMessage(request, response)
	ret := response.GetLootFaithRet()

	if ret.GetFaith() != 0 {
		statistic := user.AddStatisticsValue(constant.STATISTIC_TYPE_LOOT_FAITH, ret.GetFaith())
		//user.SetDirty()
		//rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_TYPE_ATTACK, 1, 0)
		rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_GAIN_FAITH_EVENT, float64(ret.GetFaith()), 0)
		rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_LOOT_FAITH, 1, 0)
		//user.WriteMsg(BuildScorePush(statistic))
		AddScorePush(statistic, response)
	}

	if !isRobot(ret.GetTargetID()) {
		if ret.GetShield() {
			newsFeedMessage := BuildNewsFeedMessage(user.GetID(), constant.NEWSFEED_TYPE_BE_SHIELD, constant.SHIELD_TYPE_FAITH, 0, 0)
			rpc.UserServiceProxy.PersistCall(ret.GetTargetID(), newsFeedMessage)
		} else {
			//没有建筑不需要扣除信仰值
			if ret.GetHasBuilding() && !ret.GetIsMax() && mutual {
				newsFeedMessage := BuildNewsFeedMessage(user.GetID(), constant.NEWSFEED_TYPE_BE_LOOT_FAITH, ret.GetFaith(), user.GetStarType(), 0)
				rpc.UserServiceProxy.PersistCall(ret.GetTargetID(), newsFeedMessage)
				rpc.PassportServiceProxy.WechatEventPush(ret.GetTargetID(), constant.EVENT_FAITHLOSE, 0)
			}
			AddUserNewsFeed(user, BuildNewsFeed(ret.GetTargetID(), constant.NEWSFEED_TYPE_LOOT_FAITH, ret.GetFaith(), 0, 0))
		}
	}
	response.LootFaithRet = ret
}

type AtkStarBuildingService struct {
}

func (service *AtkStarBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求进攻某个建筑物

	//message := request.GetAtkStarBuilding()
	//result := rpc.EventServiceProxy.HandleMessage(request)
	mutual := user.IsFlagUnlock(constant.STAR_FLAG_MUTUAL)
	request.AtkStarBuilding.Mutual = mutual

	event := user.EnsureEvent()
	event.HandleModuleMessage(request, response)
	ret := response.GetAtkStarBuildingRet()

	buildingID := request.GetAtkStarBuilding().GetBuildingID()

	var buildingLevel []string = nil
	if ret.GetSuccess() {
		buildingLevel = []string{character.Int32ToString(request.GetAtkStarBuilding().GetBuildingLevel())}
		statistic := user.AddStatisticsValue(constant.STATISTIC_TYPE_ATK_BUILDING, 1)
		//user.SetDirty()
		//rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_TYPE_ATTACK, 1, 0)
		rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_ATK_BUILDING, 1, 0)
		//user.WriteMsg(BuildScorePush(statistic))
		AddScorePush(statistic, response)
	}

	if !isRobot(ret.GetTargetID()) {
		if ret.GetShield() {
			newsFeedMessage := BuildNewsFeedMessage(user.GetID(), constant.NEWSFEED_TYPE_BE_SHIELD, constant.SHIELD_TYPE_BULDING, 0, 0)
			rpc.UserServiceProxy.PersistCall(ret.GetTargetID(), newsFeedMessage)
		} else {
			if !ret.GetIsMax() && mutual {
				newsFeedMessage := BuildNewsFeedMessage1(user.GetID(), constant.NEWSFEED_TYPE_BE_ATK_BUILD, ret.GetFaith(), user.GetStarType(), buildingID, buildingLevel)
				rpc.UserServiceProxy.PersistCall(ret.GetTargetID(), newsFeedMessage)
			}
			AddUserNewsFeed(user, BuildNewsFeed1(ret.GetTargetID(), constant.NEWSFEED_TYPE_ATK_BUILD, ret.GetFaith(), ret.GetItemID(), buildingID, buildingLevel))
			//user.SetDirty()
		}
	}


	response.AtkStarBuildingRet = ret
}

type LootBelieverService struct {
}

func (service *LootBelieverService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	//message := request.GetLootBeliever()
	//task := user.EnsureRefTask(message.GetEventID())

	//result := rpc.EventServiceProxy.HandleMessage(request)
	//message := request.GetLootBeliever()
	mutual := user.IsFlagUnlock(constant.STAR_FLAG_MUTUAL)
	request.LootBeliever.Mutual = mutual

	event := user.EnsureEvent()
	event.HandleModuleMessage(request, response)
	ret := response.GetLootBelieverRet()


	believerIDs := request.GetLootBeliever().GetBelieverID()


	if believerIDs != nil {
		var score int32 = 0
		for _, believerID := range believerIDs {
			score += conf.DATA.BELIEVER_SCORE_MAPPING[getBelieverLevel(believerID)]
		}
		statistic := user.AddStatisticsValue(constant.STATISTIC_TYPE_LOOT_BELIEVER, score)
		//user.SetDirty()
		//rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_TYPE_ATTACK, 1, 0)
		rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_LOOT_BELIEVER, 1, 0)
		//user.WriteMsg(BuildScorePush(statistic))
		AddScorePush(statistic, response)
	}

	if !isRobot(ret.GetTargetID()) {
		if ret.GetShield() {
			newsFeedMessage := BuildNewsFeedMessage(user.GetID(), constant.NEWSFEED_TYPE_BE_SHIELD, constant.SHIELD_TYPE_BELIEVER, 0, 0)
			rpc.UserServiceProxy.PersistCall(ret.GetTargetID(), newsFeedMessage)
		} else {
		 	if !ret.GetIsMax() && mutual {
				newsFeedMessage := BuildNewsFeedMessage1(user.GetID(), constant.NEWSFEED_TYPE_BE_LOOT_BELIEVER, 0, user.GetStarType(), 0, ret.GetBelieverID())
				rpc.UserServiceProxy.PersistCall(ret.GetTargetID(), newsFeedMessage)

				rpc.PassportServiceProxy.WechatEventPush(ret.GetTargetID(), constant.EVENT_BELIEVER_STEAL, 0)
			}
			AddUserNewsFeed(user, BuildNewsFeed1(ret.GetTargetID(), constant.NEWSFEED_TYPE_LOOT_BELIEVER, 0, 0, 0, ret.GetBelieverID()))
			//user.SetDirty()
		}
	}
	response.LootBelieverRet = ret
}

//type GenEventService struct {
//
//}
//
//func (service *GenEventService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetGenEvent()
//	event := user.GenEvent(message.GetEventType(), message.GetUid(), message.GetNickname(), message.GetGuide(), user)
//	//event.SetHandler(user)
//	response.GenEventRet = &protocol.GenEventRet{
//		Event:event.BuildEventProtocol(),
//	}
//}



//选择事件目标
type SelectEventTargetService struct {

}

func (service *SelectEventTargetService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	//message := request.GetSelectEventTarget()
	event := user.EnsureEvent()
	event.HandleModuleMessage(request, response)

	message := request.GetSelectEventTarget()
	if message != nil {
		user.AppendEventRevengeStatistic(message.GetTargetId())
	}
}

//type RemoveEventService struct {
//
//}
//
//func (service *RemoveEventService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	//message := request.GetRemoveEvent()
//	user.DeleteEvent()
//}

//type UpdateEventFieldService struct {
//
//}
//
//func (service *UpdateEventFieldService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	//message := request.GetUpdateEventField()
//	event := user.EnsureEvent()
//	event.HandleModuleMessage(request, response)
//	response.UpdateEventFieldRet = &protocol.UpdateEventFieldRet{
//		Result:true,
//	}
//}



type EventModuleInfoService struct {

}

func (service *EventModuleInfoService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetGetEventModule()
	event := user.EnsureEvent()
	event.HandleModuleMessage(request, response)
	response.GetEventModuleRet = &protocol.GetEventModuleRet{
		EventID:message.EventID,
	}

	module := event.GetModule(conf.ModuleEnum(message.GetModuleID()))
	if (module == nil) {
		exception.GameException(exception.EVENT_MODULE_NOTFOUND)
	}
	response.GetEventModuleRet.StepModule = eventModule.BuildModule(module, event.GetID())
}

//type PublicVoteService struct {
//
//}
//
//func (service *PublicVoteService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetPublicVote()
//	event := user.EnsureEvent()
//	event.HandleModuleMessage(request, response)
//	response.PublicVoteRet = &protocol.PublicVoteRet{
//		EventID:message.EventID,
//		EndTimestamp:event.GetCurrentStepModule().EndTimestamp.Unix(),
//	}
//}

//type AddVoteService struct {
//
//}
//
//func (service *AddVoteService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	//message := request.GetAddEventVote()
//	event := user.EnsureEvent()
//	event.HandleModuleMessage(request, response)
//	response.AddEventVoteRet = &protocol.AddEventVoteRet{
//		Result:true,
//	}
//}

//type CaptureBelieverService struct {
//
//}
//
//func (service *CaptureBelieverService)Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
//	message := request.GetCaptureBeliever()
//	event := session.EventManager.EnsureEvent(message.GetEventID())
//	event.HandleModuleMessage(request, response)
//}

//type SaveDataService struct {
//
//}
//
//func (service *SaveDataService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	//message := request.GetSaveData()
//	event := user.EnsureEvent()
//	event.HandleModuleMessage(request, response)
//	response.SaveDataRet = &protocol.SaveDataRet{
//		Resule:true,
//	}
//}

type GetFaithService struct {

}

func (service *GetFaithService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	//message := request.GetGetFaith()
	event := user.EnsureEvent()
	event.HandleModuleMessage(request, response)
}

type GetBelieverService struct {

}

func (service *GetBelieverService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	//message := request.GetGetBeliever()
	event := user.EnsureEvent()
	event.HandleModuleMessage(request, response)
}

type DoneEventStepService struct {

}

func (service *DoneEventStepService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	//message := request.GetDoneEventStep()
	event := user.EnsureEvent()
	event.HandleModuleMessage(request, response)
	response.DoneEventStepRet = &protocol.DoneEventStepRet{
		Result:true,
	}
}

type IntoEventService struct {

}

func (service *IntoEventService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	//message := request.GetIntoEvent()
	event := user.EnsureEvent()
	//event.IntoEvent(message.GetUid())
	response.IntoEventRet = &protocol.IntoEventRet{
		Event:event.BuildEventProtocol(),
	}
}

//type LeaveEventService struct {
//
//}
//
//func (service *LeaveEventService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetLeaveEvent()
//	event := user.EnsureEvent()
//	event.OutEvent(message.GetUid())
//	response.LeaveEventRet = &protocol.LeaveEventRet{
//		Result:true,
//	}
//}


//type EventInfoService struct {
//
//}
//
//func (service *EventInfoService)Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	//message := request.GetEvent
//	//eventIDArray := message.GetEventID()
//	//for _, eventID := range eventIDArray {
//	//	if (event == nil) {
//	//		continue
//	//	}
//	//}
//	results := []*protocol.Event{}
//
//	event := user.EnsureEvent()
//	eventProtocol := event.BuildEventProtocol()
//	results = append(results, eventProtocol)
//	response.GetEventRet = &protocol.GetEventRet{
//		Event: results,
//	}
//}
