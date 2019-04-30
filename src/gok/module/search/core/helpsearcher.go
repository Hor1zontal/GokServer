/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/10/26
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

import (
	"time"
	"gok/constant"
	"gok/module/search/cache"
	"gok/module/search/db"
	"aliens/log"
	"gok/service/lpc"
	"gok/module/search/conf"
)

var ItemHelpSearcher = &HelpSearcher{}

const (
	NotFound int = -1

	//ReceiveInterval = 3 * time.Minute
	MaxSend = 50

	OnlineInterval = 2 * 60 * time.Second
	Day1 = 24 * time.Hour
	Day3 = 3 * 24 * time.Hour

	ActiveTypeNone uint8 = 0
	ActiveTypeOnline uint8 = 1
	ActiveTypeDay1 uint8 = 2
	ActiveTypeDay3 uint8 = 3
)

//帮助对象搜索
type HelpSearcher struct {
	conditions []*Condition //过滤条件  按顺序排列
	indies 	map[int32]*db.UserIndex
	conditionLen int

}

func(this *HelpSearcher) Init() {
	this.conditions = []*Condition{}
	this.addCondition(ActiveTypeOnline,true, false)
	this.addCondition(ActiveTypeOnline,true, true)
	this.addCondition(ActiveTypeOnline,false, false)
	this.addCondition(ActiveTypeOnline,false,true)
	this.addCondition(ActiveTypeDay1,true, false)
	this.addCondition(ActiveTypeDay1,true, true)
	this.addCondition(ActiveTypeDay1,false, false)
	this.addCondition(ActiveTypeDay1,false, true)
	this.addCondition(ActiveTypeDay3,true, false)
	this.addCondition(ActiveTypeDay3,true, true)
	this.addCondition(ActiveTypeDay3,false, false)
	this.addCondition(ActiveTypeDay3,false, true)
	this.conditionLen = len(this.conditions)

	this.indies = make(map[int32]*db.UserIndex)

	var userIndies []*db.UserIndex
	db.DatabaseHandler.QueryAll(&db.UserIndex{}, &userIndies)
	//currTime := time.Now()
	for _, userIndex := range userIndies {
		//if this.UpdateIndex(currTime, userIndex) {
		//	this.indies[userIndex.Uid] = userIndex
		//}
		this.indies[userIndex.Uid] = userIndex
	}
}


func (this *HelpSearcher) Clean() {
	currTime := time.Now()
	log.Debug("start clean help search data... %v", currTime)
	//清除3天没有活跃的玩家
	for uid, index := range this.indies {
		offlineTime := currTime.Sub(index.ActiveTime)
		if offlineTime > Day3 {
			delete(this.indies, uid)
			lpc.DBServiceProxy.Delete(index, db.DatabaseHandler)
		}
	}
	log.Debug("end clean help search data..., duration %v(s)", time.Now().Sub(currTime))
}

func (this *HelpSearcher) addCondition(ActiveType uint8, Star bool, Receive bool) {
	//for _, starType := range conf.Base.StarWeCanArrive {
	//
	//}
	this.conditions = append(this.conditions, &Condition{
		//StarType:starType,
		EqualStar:  Star,
		Receive:    Receive,
		ActiveType: ActiveType,
		//mapping:make(map[int32]struct{}),
	})
}

//检索条件
type Condition struct {

	EqualStar bool //是否相同星球
	//StarType 	int32   //星球类型
	Receive     bool  	//在特定时长内是否接收过求助消息

	ActiveType  uint8 	//0在线 1 离线一天内 2离线三天内

	//mapping  	map[int32]struct{}  //满足条件的玩家id
}


func (this *HelpSearcher) EnsureUserIndex(uid int32) *db.UserIndex {
	index := this.indies[uid]
	if index == nil {
		starType := cache.StarCache.GetUserActiveStarType(uid)
		if starType > 0 {
			index = &db.UserIndex{Uid:uid, StarType:starType}
			this.indies[uid] = index
		}
	}
	return index
}

//func (this *HelpSearcher) UpdateIndex(currTime time.Time, userIndex *db.UserIndex) bool {
//	userActiveType := ActiveTypeNone
//	offlineTime := currTime.Sub(userIndex.ActiveTime)
//	if offlineTime < OnlineInterval {
//		userActiveType = ActiveTypeOnline
//	} else if offlineTime < Day1 {
//		userActiveType = ActiveTypeDay1
//	} else if offlineTime < Day3 {
//		userActiveType = ActiveTypeDay3
//	}
//
//	userReceive := currTime.Sub(userIndex.ReceiveTime) <= time.Minute
//
//	//matchConditionID := NotFound
//	//for conditionID, condition := range this.conditions {
//	//	if condition.StarType == userIndex.StarType && condition.ActiveType == userActiveType && condition.Receive == userReceive {
//	//		condition.mapping[userIndex.Uid] = struct{}{}
//	//		matchConditionID = conditionID
//	//		break
//	//	}
//	//}
//	//
//	//if matchConditionID == NotFound {
//	//	return false
//	//}
//	//
//	////索引变化需要删除索引
//	//if userIndex.ConditionID > 0 && userIndex.ConditionID != matchConditionID {
//	//	this.RemoveIndex(userIndex.ConditionID, userIndex.Uid)
//	//}
//	//
//	////更新索引
//	//userIndex.ConditionID = matchConditionID
//	return true
//}

//func (this *HelpSearcher) RemoveIndex(conditionID int, uid int32) {
//	if conditionID >= this.conditionLen {
//		return
//	}
//	condition := this.conditions[conditionID]
//	delete(condition.mapping, uid)
//}



func (this *HelpSearcher) Opt(uids []int32, opt int32, param int32, sync bool)  {
	for _, uid := range uids {
		index := this.EnsureUserIndex(uid)
		if index == nil {
			return
		}
		switch opt {
		case constant.SEARCH_OPT_UPDATE_RECEIVE_HELP:
			index.ReceiveTime = time.Now()
			break
		case constant.SEARCH_OPT_UPDATE_ACTIVE:
			index.ActiveTime = time.Now()
			break
		case constant.SEARCH_OPT_CHANGE_STAR:
			index.StarType = param
			break
		default:
		}
		if !sync {
			lpc.DBServiceProxy.ForceUpdate(index, db.DatabaseHandler)
		}
	}
}


//搜索帮助目标
func (this *HelpSearcher) RandomTargets(uid int32, starType int32, count int) []int32 {
	currTime := time.Now()
	matchUser := make(map[int32]*db.UserIndex)
	for _, condition := range this.conditions {
		for userID, userIndex  := range this.indies {
			if uid == userID {
				continue
			}
			userActiveType := ActiveTypeNone
			offlineTime := currTime.Sub(userIndex.ActiveTime)
			if offlineTime < OnlineInterval {
				userActiveType = ActiveTypeOnline
			} else if offlineTime < Day1 {
				userActiveType = ActiveTypeDay1
			} else if offlineTime < Day3 {
				userActiveType = ActiveTypeDay3
			}

			userReceive := currTime.Sub(userIndex.ReceiveTime).Seconds() <= conf.Base.ReceiveNewsInterval
			equalStar := starType == userIndex.StarType
			if condition.EqualStar == equalStar && condition.ActiveType == userActiveType && condition.Receive == userReceive {
				//matchUser = append(matchUser, userID)
				matchUser[userID] = userIndex

				if len(matchUser) >= MaxSend {
					return trans(matchUser)
				}
			}

		}
	}

	return trans(matchUser)
}

func trans(datas map[int32]*db.UserIndex) []int32 {
	result := make([]int32, 0)
	for uid, data := range datas {
		result = append(result, uid)
		data.ReceiveTime = time.Now()
	}
	return result
}