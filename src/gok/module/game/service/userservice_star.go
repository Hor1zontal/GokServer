/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/5/14
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"gok/constant"
	"gok/module/game/conf"
	"gok/module/game/global"
	"gok/module/game/user"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"gok/service/rpc"
)

//获取角色星球信息
type GetStarInfoService struct {
}

func (service *GetStarInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	result := rpc.StarServiceProxy.Call(user.GetID(), request).GetGetStarInfoRet()
	response.GetStarInfoRet = result
}

//随机出能选择的星球
type GetStarsSelectService struct {

}

func (service *GetStarsSelectService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	result, succ := user.EnsureStarsSelect()
	if !succ {
		request.GetStarsSelect.Uid = user.GetID()
		result = rpc.StarServiceProxy.HandleMessage(request).GetGetStarsSelectRet().GetStarsType()
		user.SetStarsSelect(result)
	}
	response.GetStarsSelectRet = &protocol.GetStarsSelectRet{StarsType:result}
}

//选择星球
type SelectStarService struct {
}

func (service *SelectStarService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetSelectStar()
	if !user.EnsureSelectInStars(message.GetStarType()) {
		exception.GameException(exception.STAR_SELECT_NOT_FOUND)
	}
	resp := rpc.StarServiceProxy.AllocUserNewStar(user.GetID(), message.GetStarType())

	user.UpdateStarInfo(resp.GetCurrentStar().GetStarID(), resp.GetCurrentStar().GetType())
	user.EnsureItemRandom(resp.GetCurrentStar().GetType(), true)
	user.UpdateStarFlags(resp.GetCurrentStar().GetStarFlags())

	user.TakeInFaith(resp.GetFaith(),constant.OPT_TYPE_DRAW_BUILDING_FAITH, 0 )
	user.SetPowerLimit(conf.DATA.InitPowerLimit)
	user.CleanMallItems()
	user.CleanStarsSelect()

	if !resp.IsFirst {
		newsFeed := BuildNewsFeed(user.GetID(), constant.NEWSFEED_TYPE_NEXT_STAR, resp.GetLastStarType(), 0, 0)
		global.BroadcastMessage(&protocol.GlobalMessage{NewsFeed: newsFeed})
	}

	response.SelectStarRet = &protocol.SelectStarRet{
		Star:resp.GetStar(),
		CurrentStar:resp.GetCurrentStar(),
		LastStarType:resp.GetLastStarType(),
		Faith:resp.GetFaith(),
		Items:resp.GetItems(),
	}
}

type GetStarInfoDetailService struct {
}

func (service *GetStarInfoDetailService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //查询星球详细信息
	message := request.GetGetStarInfoDetail() //获取消息包
	resp := rpc.StarServiceProxy.Call(message.GetUid(), request)
	response.GetStarInfoDetailRet = resp.GetGetStarInfoDetailRet()
}

type GetStarShieldService struct {
}

func (service *GetStarShieldService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //查询星球详细信息
	//message := request.GetGetStarInfoDetail() //获取消息包
	resp := rpc.StarServiceProxy.Call(user.GetID(), request)
	//resp := rpc.StarServiceProxy.GetStarInfoDetail(message.GetUid(), message.GetStarID(), message.GetDestUid())
	response.GetStarShieldRet = resp.GetGetStarShieldRet()
}

type StarStatisticsService struct {
}

func (service *StarStatisticsService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //查询星球详细信息
	resp := rpc.StarServiceProxy.Call(user.GetID(), request)
	response.GetStarStatisticsRet = resp.GetGetStarStatisticsRet()
}

type StarHistoryService struct {
}

func (service *StarHistoryService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //查询星球详细信息
	//request.GetStarHistory.Uid = user.GetID()
	resp := rpc.StarServiceProxy.Call(user.GetID(), request)
	response.GetStarHistoryRet = resp.GetGetStarHistoryRet()
}

type GetStarInfoComplete struct {
	//获取已完成星球数据 20170607
}

func (service *GetStarInfoComplete) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {

}


//type NextStarService struct {
//	//请求占领星球
//}
//
//func (service *NextStarService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求占领星球
//	message := request.GetNextStar()
//	message.Uid = user.GetID()
//	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetNextStarRet()
//
//	user.UpdateStarInfo(resp.GetStar().GetStarID(), resp.GetStar().GetType())
//	user.EnsureItemRandom(resp.GetStar().GetType(),true)
//
//	//user.Record(resp.GetLastStarType())
//	//items := user.TakeinAllTempItem()
//	user.TakeInFaith(resp.GetFaith(), constant.OPT_TYPE_DRAW_BUILDING_FAITH, 0)
//	user.SetPowerLimit(0)
//	//user.RestoreStarStatistics()
//	//resp.Items = items
//	user.DataManager. SetDirty()
//
//	newsFeed := BuildNewsFeed(user.GetID(), constant.NEWSFEED_TYPE_NEXT_STAR, resp.GetLastStarType(), 0, 0)
//	global.BroadcastMessage(&protocol.GlobalMessage{NewsFeed: newsFeed})
//
//	//rpc.CommunityServiceProxy.PublicMoments(user.GetID(), constant.MOMENTS_TYPE_STAR, resp.GetLastStarType())
//
//	response.NextStarRet = resp
//	response.NextStarRet.Result = true
//}
//type GetStarRecordService struct {
//}
//
//func (service *GetStarRecordService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求获取星球记录
//	resp := rpc.StarServiceProxy.GetRecordStarInfo(user.GetRecordOri(), user.GetRecordUser())
//	for _, v := range resp.StarsOri {
//		if v.GetOwnID() != 0 {
//			v.OwnNikeName = cache.UserCache.GetUserNickname(v.GetOwnID()))
//		}
//	}
//
//	for _, v := range resp.StarsUser {
//		if v.GetOwnID() != 0 {
//			v.OwnNikeName = cache.UserCache.GetUserNickname(v.GetOwnID()))
//		}
//	}
//	response.GetStarRecordInfoRet = resp
//}

//type SetStarRecordService struct {
//}
//
//func (service *SetStarRecordService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求设置星球记录
//	message := request.GetSetStarRecordInfo()
//	if user.GetFaith() < conf.GetGameBase().CostRecordStarFaith { //需要消耗的信仰值不够
//		exception.GameException(exception.FAITH_NOT_ENOUGH)
//	}
//	user.UpdateStarRecord(message.GetRecordType(), message.GetRecordID()) //更新收藏数据
//
//	response.SetStarRecordInfoRet = &protocol.SetStarRecordInfoRet{
//		RecordType: message.GetRecordType()),
//		RecordID:   message.GetRecordID()),
//	}
//	user.TakeOutFaith(conf.DATA.CostRecordStarFaith, constant.OPT_TYPE_RECORD_STAR, 0)
//	user.WriteMsg(BuildRoleSocialPush(user))
//}
//
//type DelStarRecordService struct {
//}
//
//func (service *DelStarRecordService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求删除星球记录
//	message := request.GetDelStarRecordInfo()
//	ret := user.DelStarRecord(message.GetRecordType(), message.GetRecordID()) //删除星球记录
//	if !ret {
//		exception.GameException(exception.STAR_NOTFOUND)
//	}
//
//	response.DelStarRecordInfoRet = &protocol.DelStarRecordInfoRet{
//		RecordType: message.GetRecordType()),
//		RecordID:   message.GetRecordID()),
//	}
//}
//
//type MoveStarRecordService struct {
//	//请求移动星球记录 wjl 20170607
//}
//
//func (service *MoveStarRecordService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetMoveStarRecordInfo()
//	ret := user.MoveStarRecord(message.GetSrcType(), message.GetSrcID(), message.GetDestType(), message.GetDestID())
//	if !ret {
//		exception.GameException(exception.STAR_RECORD_MOVE_FAILED)
//	}
//	response.MoveStarRecordInfoRet = &protocol.MoveStarRecordInfoRet{
//		SrcType:  message.GetSrcType()),
//		SrcID:    message.GetSrcID()),
//		DestType: message.GetDestType()),
//		DestID:   message.GetDestID()),
//	}
//}
//
//type ReplaceStarRecordService struct {
//	//替换指定星球
//}
//
//func (service *ReplaceStarRecordService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetReplaceStarRecordInfo()
//	if user.GetFaith() < conf.GetGameBase().CostRecordStarFaith { //需要消耗的信仰值不够
//		exception.GameException(exception.FAITH_NOT_ENOUGH)
//	}
//	user.ReplaceStarRecord(message.GetRecordType(), message.GetRecordID(), message.GetReplaceRecordID())
//
//	response.ReplaceStarRecordInfoRet = &protocol.ReplaceStarRecordInfoRet{
//		RecordType:      message.GetRecordType()),
//		RecordID:        message.GetRecordID()),
//		ReplaceRecordID: message.GetReplaceRecordID()),
//	}
//	user.TakeOutFaith(conf.GetGameBase().CostRecordStarFaith, constant.OPT_TYPE_RECORD_STAR, 0)
//	user.WriteMsg(BuildRoleSocialPush(user))
//}
//
//type OccupyStarService struct {
//	//请求占领星球
//}
//
//func (service *OccupyStarService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求占领星球
//	message := request.GetOccupyStar()
//	message.Uid = user.GetID())
//	resp := rpc.StarServiceProxy.OccupyStar(message.GetStarID(), message.GetUid())
//	user.OccupyStar(resp.StarOld, resp.Star)
//	response.OccupyStarRet = resp
//}


//掠夺圣物
type LootItemService struct {
}

func (service *LootItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetLootItem()
	searchID := message.GetSearchID()
	searchResult := user.GetSearch(searchID)
	if searchResult == nil {
		exception.GameException(exception.SEARCH_RESULT_INVALID)
	}
	message.ItemID = searchResult.ItemID

	result := rpc.StarServiceProxy.Call(searchResult.GetId(), request).GetLootItemRet()

	//抢夺成功需要通知被抢夺的玩家
	if result.GetResult() {
		notifyID := message.GetLootID()
		if !isRobot(notifyID) {
			newsFeedMessage := BuildNewsFeedMessage(user.GetID(), constant.NEWSFEED_TYPE_BE_LOOT_ITEM, message.GetItemID(), result.GetBuilding(), 0)
			rpc.UserServiceProxy.PersistCall(notifyID, newsFeedMessage)

			//物品被抢需要全服通知
			if message.GetItemID() > 0 {
				newsFeed := BuildNewsFeed(user.GetID(), constant.NEWSFEED_TYPE_BE_LOOT_ITEM, message.GetItemID(), result.GetBuilding(), 0)
				global.BroadcastMessage(&protocol.GlobalMessage{NewsFeed: newsFeed})

				AddUserNewsFeed(user, BuildNewsFeed(notifyID, constant.NEWSFEED_TYPE_LOOT_ITEM, message.GetItemID(), result.GetBuilding(), 0))
			}
		}

		if message.GetItemID() > 0 {
			user.TakeInItem(message.GetItemID(), 1, constant.OPT_TYPE_EVENT_LOOT, notifyID)
		}
	}

	user.RemoveSearch(searchID)
	
	response.LootItemRet = result
}


//请求维修星球建筑
type AccRepairStarBuildingService struct {
}

func (service *AccRepairStarBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求维修某个建筑物
	resp := rpc.StarServiceProxy.Call(user.GetID(),request).GetAccRepairStarBuildRet()
	if resp.GetDone() {
		resp.ItemID = user.DealRepairedBuilding1(resp.GetBuildingType(), resp.GetBuildingLevel(), true)
	}
	user.DataManager.SetDirty()
	response.AccRepairStarBuildRet = resp
}

//请求开始维修星球建筑
type RepairStarBuildingService struct {
}

func (service *RepairStarBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求开始维修建筑物
	message := request.GetRepairStarBuild()
	message.Faith = user.GetFaith()

	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetRepairStarBuildRet()
	user.TakeOutFaith(resp.GetCost(), constant.OPT_TYPE_REPAIR_BUILDING, message.GetBuildingType()) //消耗相应的信仰值
	rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_EXPEND_FAITH_REPAIRE, float64(resp.GetCost()), 0)

	if resp.GetDone() {
		resp.ItemID = user.DealRepairedBuilding1(resp.GetBuildingType(), resp.GetBuildingLevel(), true)
	}
	user.DataManager.SetDirty()
	//user.WriteMsg(BuildRoleSocialPush(user))
	AddRoleSocialPush(user, response)
	response.RepairStarBuildRet = resp
}

type UpgradeStarBuildingEndService struct {
	//请求判断建筑物升级是否结束
}

func (service *UpgradeStarBuildingEndService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //请求判断建筑物升级是否结束
	//是否引导操作
	//if user.IsUpgradeGuiding() {
	//	exception.GameException(exception.INVALID_FLAG_AUTHORITY)
	//}
	//oldMaxLevel := cache.StarCache.GetBuildingExMaxLevel(user.GetStarId())
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetUpdateStarBuildEndRet()

	if resp.GetDone() {
		resp.ItemID = user.DealUpgradeBuilding1(resp.GetItemID(), resp.GetPowerLimit(), resp.GetPowerReward(), resp.GetBuildingType(), resp.GetLevel(), true)
	}
	user.DataManager.SetDirty()
	response.UpdateStarBuildEndRet = resp
}

type RepairStarBuildEndService struct {
	//请求判断建筑物修理是否结束
}

func (service *RepairStarBuildEndService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //请求判断建筑物升级是否结束
	message := request.GetRepairStarBuildEnd()
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetRepairStarBuildEndRet()
	if resp.GetDone() {
		resp.ItemID = user.DealRepairedBuilding1(message.GetBuildingType(), resp.GetBuildingLevel(), true)
	}
	user.DataManager.SetDirty()
	response.RepairStarBuildEndRet = resp
}

//请求加速升级建筑
type AccUpgradeStarBuildingService struct {
}

func (service *AccUpgradeStarBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //请求判断建筑物升级加速
	message := request.AccUpdateStarBuild
	message.Uid = user.GetID()
	//是否引导操作
	//message.Guide = user.IsUpgradeGuiding()

	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetAccUpdateStarBuildRet()
	if resp.GetDone() {
		resp.ItemID = user.DealUpgradeBuilding1( 0, resp.GetPowerLimit(), 0, resp.GetBuildingType(), resp.GetLevel(), true)
	}
	user.DataManager.SetDirty()
	response.AccUpdateStarBuildRet = resp
}


type BuildingFaithInfoService struct {
}

func (service *BuildingFaithInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //请求判断建筑物升级加速
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetGetBuildingFaithRet()
	response.GetBuildingFaithRet = resp
}

type ReceiveBuildingFaithService struct {
}

func (service *ReceiveBuildingFaithService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //请求判断建筑物升级加速
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetReceiveBuildingFaithRet()
	user.TakeInFaith(resp.GetBuildingFaith(), constant.OPT_TYPE_DRAW_BUILDING_FAITH, request.GetReceiveBuildingFaith().GetBuildingType())
	user.DataManager.SetDirty()
	//user.WriteMsg(BuildRoleSocialPush(user))
	AddRoleSocialPush(user, response)
	response.ReceiveBuildingFaithRet = resp
}

//升级信徒
type UpgradeBelieverService struct {
}

func (service *UpgradeBelieverService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	request.UpgradeBeliever.Faith = user.GetFaith()
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetUpgradeBelieverRet()
	
	upgradeBeliever := resp.GetResult()
	cost := resp.GetCost()
	//addCivilization := response.GetUpgradeBelieverRet().GetCivilization()

	if upgradeBeliever != nil {
		if user.GetBelieverFlag(upgradeBeliever.GetId()) {
			resp.ItemID = user.RandomItem(constant.RANDOM_BELIEVER, true, true, constant.OPT_TYPE_UPGRADE_BELIEVER)
		} else {
			user.UpdateBelieverFlag(upgradeBeliever.GetId(), true)
		}
		user.UpdateLastUpgradeBelieverID(upgradeBeliever.GetId())
	}

	//if addCivilization > 0 {
	//	user.TakeInCivilization(addCivilization, user.GetStarSize())
	//	user.WriteMsg(user.BuildRoleCivilizationPush())
	//}

	if cost>0 {
		user.TakeOutFaith(cost, constant.OPT_TYPE_UPGRADE_BELIEVER, 0)
		rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_EXPEND_FAITH_BELIEVER, float64(cost), 0)
	}
	//user.DataManager.SetDirty()
	response.UpgradeBelieverRet = resp

}

//取消升级建筑
type CancelUpgradeStarBuildService struct {
}

func (service *CancelUpgradeStarBuildService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetCancelUpgradeStarBuildRet()
	if resp.GetBackFaith() > 0 {
		user.TakeInFaith(resp.GetBackFaith(), constant.OPT_TYPE_CANCEL_BUILDING, 0)
	}
	AddRoleSocialPush(user, response)
	response.CancelUpgradeStarBuildRet = resp
}

//取消修理建筑
type CancelRepairStarBuildService struct {
}

func (service *CancelRepairStarBuildService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetCancelRepairStarBuildRet()
	if resp.GetBackFaith() > 0 {
		user.TakeInFaith(resp.GetBackFaith(), constant.OPT_TYPE_CANCEL_BUILDING, 0)
	}
	AddRoleSocialPush(user, response)
	response.CancelRepairStarBuildRet = resp
}

type GetHelpRepairInfoService struct {

}

func (service *GetHelpRepairInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.GetHelpRepairInfoRet = rpc.StarServiceProxy.Call(user.GetID(), request).GetGetHelpRepairInfoRet()
}

type StarSettleService struct {

}

func (service *StarSettleService)  Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetStarSettleRet()
	resp.RankNum, resp.TotalRank = user.GetUserStarRank(constant.STATISTIC_TYPE_STAR_ONLINE, int32(resp.GetPassTime()))
	response.StarSettleRet = resp
}

type StarFlagInfoService struct {

}

func (service *StarFlagInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.StarFlagInfoRet = rpc.StarServiceProxy.Call(user.GetID(), request).GetStarFlagInfoRet()
}

type UpdateStarFlagService struct {

}

func (service *UpdateStarFlagService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.UpdateStarFlagRet = rpc.StarServiceProxy.Call(user.GetID(), request).GetUpdateStarFlagRet()
}