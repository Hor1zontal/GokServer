/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/10/22
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

import (
	"gok/service/msg/protocol"
	"gok/service/exception"
	"gok/module/trade/db"
	"gok/module/trade/cache"
	"gok/service/lpc"
	"gok/constant"
	"math/rand"
	"aliens/common"
	"gopkg.in/mgo.v2/bson"
	"aliens/log"
	"time"
)

var HelpHandler = GetHelpManager()

func GetHelpManager() *HelpManager {
	return &HelpManager{
		//sales:make(map[int32]*protocol.Sale),
		//goods:make(map[string]*protocol.Goods),
		//userGoods:make(map[int32][]*protocol.Goods),
	}
}

type HelpManager struct {
	//helps map[int32]*protocol.ItemHelp
}

func (this *HelpManager) GetItemHelp(uid int32, id string) *protocol.ItemHelp {
	current := cache.TradeCache.GetItemHelp(uid)
	if id == "" {
		return current
	}

	if current != nil && current.GetId() == id {
		return current
	}

	historyData := &db.ItemHelpHistory{ID:id, Uid:uid}
	//从历史中查询
	err := db.DatabaseHandler.QueryOne(historyData)
	if err != nil {
		return nil
	}

	itemHelp := &protocol.ItemHelp{}
	err = itemHelp.Unmarshal(historyData.Data)
	if err != nil {
		log.Error("itemHelp unmarshal err: %v", err.Error())
		return nil
	}
	itemHelp.Overdue = true
	return itemHelp
}

func (this *HelpManager) SetItemHelp(uid int32, itemHelp *protocol.ItemHelp) bool {
	result := cache.TradeCache.SetItemHelp(uid, itemHelp)
	if result {
		lpc.DBServiceProxy.ForceUpdate(itemHelp, db.DatabaseHandler)
	} else {
		exception.GameException(exception.DATABASE_EXCEPTION)
	}
	return result
}

func (this *HelpManager) CancelItemHelp(uid int32, itemHelp *protocol.ItemHelp) bool {
	itemHelp.Overdue = true
	data, err := itemHelp.Marshal()
	if err != nil {
		log.Error("itemHelp unmarshal err: %v", err.Error())
	} else {
		lpc.DBServiceProxy.Insert(&db.ItemHelpHistory{
			ID:itemHelp.GetId(),
			Uid:itemHelp.GetUid(),
			Data:data,
			Time:time.Now(),
		}, db.DatabaseHandler)
	}
	cache.TradeCache.RemoveItemHelp(uid)
	condition := bson.D{{"uid", uid}}
	lpc.DBServiceProxy.DeleteCondition(&protocol.ItemHelp{Uid:uid}, condition, db.DatabaseHandler)
	return true
}


func (this *HelpManager) PublicItemHelp(uid int32, itemID int32) *protocol.ItemHelp {
	help := this.GetItemHelp(uid, "")
	if help != nil {
		exception.GameException(exception.ITEM_HELP_PUBLIC_REPEAT)
	}
	itemHelp := &protocol.ItemHelp{
		Id:util.GenUUID(),
		Uid:uid,
		ItemID:itemID,
	}
	this.SetItemHelp(uid, itemHelp)
	return itemHelp
}

func (this *HelpManager) DrawItemHelp(uid int32, itemID int32, cancel bool) (int32, *protocol.ItemHelp) {
	help := this.GetItemHelp(uid, "")
	if help == nil  || help.GetItemID() != itemID {
		exception.GameException(exception.ITEM_HELP_NOT_FOUND)
	}
	/*if help.GetItemNum() <= 0 {
		exception.GameException(exception.ITEM_HELP_ITEM_NOT_ENOUGH)
	}*/
	draw := help.GetItemNum()
	help.ItemNum = 0
	help.DrawNum += draw

	if help.Events != nil {
		for _, event := range help.Events {
			event.Draw = true
		}
	}

	if cancel {
		this.CancelItemHelp(uid, help)
	} else {
		this.SetItemHelp(uid, help)
	}
	return draw, help
}


//func (this *HelpManager) GetItemHelp(uid int32) *protocol.ItemHelp {
//	return this.helps[uid]
//}


func (this *HelpManager) LootItemHelp(uid int32, lootID int32, itemID int32, power int32,
		costs []int32, limit int32, probs []float32, addProb float32, isWatchAd bool, eachFollow bool) (bool, int32, *protocol.ItemHelp) {
	help := this.GetItemHelp(uid, "")
	if help == nil  || help.GetItemID() != itemID {
		exception.GameException(exception.ITEM_HELP_NOT_FOUND)
	}

	if help.GetItemNum() == 0 {
		exception.GameException(exception.ITEM_HELP_ITEM_NOTFOUND)
	}

	allLimit := int32(len(probs) - 1)

	if help.GetLootNum() >= allLimit {
		exception.GameException(exception.ITEM_HELP_LOOT_LIMIT)
	}

	//有帮助过 不能再偷取
	if existHelpEvent(help, lootID, constant.ITEMHELP_EVENT_HELP) {
		exception.GameException(exception.ITEM_HELP_LOOT_ALREADY_HELP)
	}

	//当前玩家的偷取总次数
	var lootCount int32 = 0
	successCount := getHelpEventCount(help, lootID, constant.ITEMHELP_EVENT_LOOT_SUCCESS)
	failedCount := getHelpEventCount(help, lootID, constant.ITEMHELP_EVENT_LOOT_FAILED)

	//成功次数限制
	if successCount >= limit {
		exception.GameException(exception.ITEM_HELP_LOOT_LIMIT)
	}
	lootCount += successCount
	lootCount += failedCount

	var cost int32 = 0
	if !isWatchAd {
		cost = getLatestCost(lootCount, costs)
	}
	if cost > power {
		exception.GameException(exception.POWER_NOT_ENOUGH)
	}

	//偷取次数越多，概率越低
	prob := getLatestProb(help.GetLootNum(), probs)

	if prob != 0 && addProb != 0 {
		prob = prob * (1 + addProb)
	}

	//偷取成功标识
	result := prob > rand.Float32()

	if result {
		help.LootNum += 1
		help.ItemNum -= 1
		addHelpEvent(help, lootID, constant.ITEMHELP_EVENT_LOOT_SUCCESS, eachFollow)
	} else {
		addHelpEvent(help, lootID, constant.ITEMHELP_EVENT_LOOT_FAILED, eachFollow)
	}
	this.SetItemHelp(uid, help)
	//lpc.StatisticsHandler.AddStatisticData(&model.StatisticSteal{Uid:lootID, TargetID:uid, ItemID: help.GetItemID(), Count:help.GetLootNum(), ItemNum:help.GetItemNum(), Success:result, Prob:prob})

	return result, cost, help
}

//获取偷取的概率
func getLatestProb(count int32, probs []float32) float32 {
	var result float32 = 0
	if probs == nil {
		return result
	}

	for index, prob := range probs {
		result = prob
		if int32(index) == count {
			return result
		}
	}
	return result
}

func getLatestCost(count int32, costs []int32) int32 {
	var result int32 = 0
	if costs == nil {
		return result
	}

	for index, cost := range costs {
		result = cost
		if int32(index) == count {
			return result
		}
	}
	return result
}

func (this *HelpManager) GetHelpItemHistory(uid int32, skip int32, limit int32, getCount bool) (int32, []*protocol.ItemHelp) {
	var datas []*db.ItemHelpHistory

	var count int32 = 0

	if getCount {
		result, err := db.DatabaseHandler.QueryConditionCount(&db.ItemHelpHistory{}, "uid", uid)
		if err != nil {
			log.Debug("db err : %v", err)
			exception.GameException(exception.DATABASE_EXCEPTION)
		}

		count = int32(result)
	}

	db.DatabaseHandler.QueryAllConditionSkipLimit(&db.ItemHelpHistory{}, "uid", uid, &datas, int(skip), int(limit), "-time")
	if datas == nil || len(datas) == 0 {
		return count, nil
	}

	results := make([]*protocol.ItemHelp, 0)
	for _, data := range datas {
		itemHelp := &protocol.ItemHelp{}
		err := itemHelp.Unmarshal(data.Data)
		if err == nil {
			itemHelp.Overdue = true
			results = append(results, itemHelp)
		}
	}
	return  count, results

}

func (this *HelpManager) HelpItemHelp(uid int32, helpID int32, itemID int32, limit int32, eachFollow bool) *protocol.ItemHelp {
	help := this.GetItemHelp(uid, "")
	if help == nil || help.GetItemID() != itemID {
		exception.GameException(exception.ITEM_HELP_NOT_FOUND)
	}
	if help.GetHelpNum() >= limit {
		exception.GameException(exception.ITEM_HELP_LIMIT)
	}

	//有偷取过不能再帮助
	if existHelpEvent(help, helpID, constant.ITEMHELP_EVENT_LOOT_FAILED) ||
		existHelpEvent(help, helpID, constant.ITEMHELP_EVENT_LOOT_SUCCESS)  {
		exception.GameException(exception.ITEM_HELP_HELP_ALREADY_LOOT)
	}

	addHelpEvent(help, helpID, constant.ITEMHELP_EVENT_HELP, eachFollow)

	help.ItemNum += 1
	help.HelpNum += 1
	this.SetItemHelp(uid, help)
	return help
}

//获取事件次数
func getHelpEventCount(help *protocol.ItemHelp, uid int32, eventType int32) int32 {
	if help.Events == nil {
		return 0
	}
	var result int32 = 0
	for _, event := range help.Events {
		if event.GetUid() == uid && event.GetType() == eventType {
			result ++
		}
	}
	return result
}

func existHelpEvent(help *protocol.ItemHelp, uid int32, eventType int32) bool {
	if help.Events == nil {
		return false
	}
	for _, event := range help.Events {
		if event.GetUid() == uid && event.GetType() == eventType {
			return true
		}
	}
	return false
}


func addHelpEvent(help *protocol.ItemHelp, uid int32, eventType int32, eachFollow bool) *protocol.ItemEvent {
	var result *protocol.ItemEvent = nil
	if help.Events == nil {
		help.Events = []*protocol.ItemEvent{}
	}

	//for _, event := range help.Events {
	//	if event.GetUid() == uid && event.GetType() == eventType {
	//		result = event
	//		break
	//	}
	//}
	//
	//if result == nil {
	//	result = &protocol.ItemEvent{
	//		Uid:uid,
	//		Type:eventType,
	//		Count:1,
	//	}
	//	help.Events = append(help.Events, result)
	//} else {
	//	result.Count += 1
	//}

	result = &protocol.ItemEvent{
		Uid:uid,
		Type:eventType,
		Draw:false,
		IsNew: !eachFollow,  //之前没有互相关注的为新的交互关系
	}
	help.Events = append(help.Events, result)
	return result
}



