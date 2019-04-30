/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/10/22
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"gok/service/msg/protocol"
	"gok/module/game/user"
	"gok/service/rpc"
	"gok/constant"
	"gok/service/exception"
	"gok/module/game/conf"
	"gok/module/game/cache"
	"time"
	"gok/module/game/global"
	"gok/module/game/db"
)





type PublicItemHelpService struct {
}

func (service *PublicItemHelpService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetPublicItemHelp()
	if message.GetItemID() <= 0 {
		exception.GameException(exception.ITEM_HELP_ITEM_NOTFOUND)
	}
	var cost int32 = 0
	if !message.GetIsWatchAd() {
		cost = conf.DATA.HelpRequestPrice
		user.EnsureGayPoint(cost)
		response.PublicItemHelpRet = rpc.TradeServiceProxy.Call(user.GetID(), request).GetPublicItemHelpRet()
		user.TakeOutGayPoint(cost, constant.OPT_TYPE_PUBLIC_HELP, 0)
	} else {
		response.PublicItemHelpRet = rpc.TradeServiceProxy.Call(user.GetID(), request).GetPublicItemHelpRet()
	}

	id := response.PublicItemHelpRet.GetItemHelp().GetId()

	//推送关注玩家
	sendMessage := &protocol.C2GS{
		Sequence:[]int32{523},
		AddNewsFeed: BuildNewsFeed1(user.GetID(), constant.NEWSFEED_TYPE_PUBLIC_ITEMHELP, message.GetItemID(), 0, 0, []string{id}),
	}
	//消息发布给所有玩家
	global.PersistCallFollowings(user.GetID(), sendMessage)
	user.SetHelpPublicTime(time.Now(), message.GetItemID())
	response.PublicItemHelpRet.Cost = cost
}

//
//type CancelItemHelpService struct {
//}
//
//func (service *CancelItemHelpService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	//message := request.GetCancelItemHelp()
//	response.CancelItemHelpRet = rpc.TradeServiceProxy.Call(user.GetID(), request).GetCancelItemHelpRet()
//	if response.CancelItemHelpRet.GetResult() {
//		cache.UserCache.DelBeHelpItemCacheUid(user.GetID(), request.GetCancelItemHelp().GetItemID())
//	}
//	user.CleanHelpPublicTime()
//}


type DrawItemHelpService struct {
}

func (service *DrawItemHelpService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	result := rpc.TradeServiceProxy.Call(user.GetID(), request).GetDrawItemHelpRet()

	if result.GetItemNum() > 0 {
		user.TakeInItem(result.GetItemID(), result.GetItemNum(), constant.OPT_TYPE_DRAW_HELP, 0)
	}

	if request.GetDrawItemHelp().GetCancel() {
		cache.UserCache.DelBeHelpItemCacheUid(user.GetID(), request.GetDrawItemHelp().GetItemID())
		user.CleanHelpPublicTime()


		itemHelp := result.GetItemHelp()
		if itemHelp != nil {
			//推送关注玩家完成圣物组合
			sendMessage := &protocol.C2GS{
				Sequence:[]int32{523},
				AddNewsFeed: BuildNewsFeed1(user.GetID(), constant.NEWSFEED_TYPE_DONE_ITEMHELP, itemHelp.GetItemID(), itemHelp.GetHelpNum(), itemHelp.GetLootNum(), []string{itemHelp.GetId()}),
			}
			global.PersistCallFollowings(user.GetID(), sendMessage)
		}

	}

	response.DrawItemHelpRet = result
}


type GetItemHelpService struct {
}

func (service *GetItemHelpService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetGetItemHelp()
	response.GetItemHelpRet = rpc.TradeServiceProxy.Call(message.GetUid(), request).GetGetItemHelpRet()

	itemHelp := response.GetItemHelpRet.ItemHelp

	var changeNewsFeed *db.DBNewsFeed = nil
	//刷新消息
	if message.GetNewsFeedID() != "" {
		if itemHelp != nil {
			if itemHelp.Overdue {
				changeNewsFeed = user.DoneItemHelpNewsFeed(message.GetNewsFeedID(), itemHelp.GetHelpNum(), itemHelp.GetLootNum())
			}
		} else {
			changeNewsFeed = user.DoneItemHelpNewsFeed(message.GetNewsFeedID(), 0, 0)
		}
	}

	if changeNewsFeed != nil {
		user.WriteMsg(BuildNewsFeedPush(changeNewsFeed.BuildProtocol()))
	}
}


type LootHelpItemService struct {
}

func (service *LootHelpItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetLootHelpItem()
	message.Power = user.GetPower()
	message.Costs = conf.DATA.StealCost
	message.Limit = conf.DATA.StealSingleLimit
	message.Probs = conf.DATA.StealRate
	message.AddProb = float32(cache.StarCache.GetBuffRelicSteal(user.GetID()))
	//message.AllLimit = conf.DATA.StealAllLimit

	followEach := rpc.CommunityServiceProxy.IsEachFollow(user.GetID(), message.GetUid())
	message.EachFollow = followEach

	result := rpc.TradeServiceProxy.Call(user.GetID(), request).GetLootHelpItemRet()

	id := result.GetItemHelp().GetId()

	//法力偷取圣物
	if !message.GetIsWatchAd() {
		user.TakeOutPower(result.GetCost(), constant.OPT_TYPE_LOOT_HELP)
	}

	if result.GetResult() {
		itemID := result.GetItemHelp().GetItemID()
		user.TakeInItem(itemID, 1, constant.OPT_TYPE_LOOT_HELP, 0)

		sendMessage := BuildNewsFeedMessage1(user.GetID(), constant.NEWSFEED_TYPE_BE_LOOT_ITEMHELP, itemID, result.GetItemHelp().GetItemNum(), 0, []string{id})
		//推送被抢夺求组圣物
		rpc.UserServiceProxy.PersistCall(message.GetUid(), sendMessage)
		
		rpc.PassportServiceProxy.WechatEventPush(message.GetUid(), constant.EVENT_RELIC_STEAL, 0)
	}

	if !followEach {
		rpc.CommunityServiceProxy.FollowEach(message.GetUid(), user.GetID())
		result.EachFollow = true
	}

	//推送玩家数据变更
	pushMessage := &protocol.GS2C{
		Sequence:[]int32{1081},
		ItemHelpPush: result.GetItemHelp(),
	}
	rpc.UserServiceProxy.Push(message.GetUid(), pushMessage)

	response.LootHelpItemRet = result
}


type HelpItemService struct {
}

func (service *HelpItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	result := user.HelpItemAsk(false, request)
	response.HelpItemRet = result
}


type GetHelpItemHistoryService struct {
}

func (service *GetHelpItemHistoryService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.GetHelpItemHistoryRet = rpc.TradeServiceProxy.HandleMessage(request).GetGetHelpItemHistoryRet()
}



