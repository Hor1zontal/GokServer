/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/10/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package rank

import (
	"aliens/common/character"
	basecache "gok/cache"
	"gok/module/cluster/center"
	"gok/service"
	"aliens/common/util"
	"aliens/log"
	"gok/constant"
	//	"aliens/common/util"
	"gok/module/game/cache"
	"gok/module/game/db"
	"time"
)

const (
	TOP_LIMIT int32 = 25
)

var RankRefreshTime time.Time = util.GetMontyBegin(constant.RANK_REFRESH_HOUR)

//人口排行管理
var SaleRankManager = &RankManager{basecache.RANK_SALE_KEY, TOP_LIMIT, constant.STATISTIC_TYPE_SALE}

//国土排行管理
var LootFaithRankManager = &RankManager{basecache.RANK_LOOT_FAITH_KEY, TOP_LIMIT, constant.STATISTIC_TYPE_LOOT_FAITH}


//攻击建筑constant.STATISTIC_TYPE_ATK_BUILDING
var AtkBuildingRankManager = &RankManager{basecache.RANK_ATK_BUIDING_KEY, TOP_LIMIT, constant.STATISTIC_TYPE_ATK_BUILDING}

//抢夺信仰
var LootBelieverRankManager = &RankManager{basecache.RANK_LOOT_BELIEVER_KEY, TOP_LIMIT, constant.STATISTIC_TYPE_LOOT_BELIEVER}



func GetRankManager(rankType int32) *RankManager {
	if rankType == constant.STATISTIC_TYPE_SALE {
		return SaleRankManager
	} else if rankType == constant.STATISTIC_TYPE_LOOT_FAITH {
		return LootFaithRankManager
	} else if rankType == constant.STATISTIC_TYPE_ATK_BUILDING {
		return AtkBuildingRankManager
	} else if rankType == constant.STATISTIC_TYPE_LOOT_BELIEVER {
		return LootBelieverRankManager
	} else {
		//exception.GameException(exception.UNEXPECT_RANK_TYPE)
	}
	return nil
}

func GetStarRankManager(rankType int32, starType int32) *RankManager {
	if rankType == constant.STATISTIC_TYPE_STAR_ONLINE {
		return &RankManager{basecache.RANK_STAR_ONLINE_KEY + character.Int32ToString(starType), TOP_LIMIT, constant.STATISTIC_TYPE_STAR_ONLINE}
	}
	return nil
}
//更新日排行数据,同时写入到数据库
//func UpdateDayRank() {
//	if center.IsMaster(service.SERVICE_USER) {
//		log.Info("%v 开始更新日排行榜信息 ...", time.Now())
//		//PeopleRankManager.UpdateDayRank()
//		//LandRankManager.UpdateDayRank()
//		log.Info("%v 结束更新日排行榜信息 ...", time.Now())
//	}
//}

func Init() {
	//if cache.UserCache.SetNX(basecache.FLAG_LOADRANK, 1) {
	//	log.Debug("start load rank data to redis cache...")
	//	var roles []*db.DBRole
	//	db.DatabaseHandler.QueryAll(&db.DBRole{}, &roles)
	//	for _, role := range roles {
	//		for _, statistic := range role.Statistics {
	//			manager := GetRankManager(statistic.ID)
	//			if manager != nil {
	//				if statistic.UpdateTime.After(RankRefreshTime) {
	//					manager.UpdateRank(role.UserID, int64(statistic.Value))
	//				} else {
	//					manager.UpdateRank(role.UserID, 0)
	//				}
	//			}
	//		}
	//
	//		//更新昵称
	//		//cache.UserCache.SetUserNickname(role.UserID, role.NickName)
	//		//cache.UserCache.SetUserRole(role.UserID, role.ID)
	//	}
	//	log.Debug("end load rank data to redis cache")
	//}

	if cache.UserCache.SetNX(basecache.FLAG_LOADRANK, 1) {
		log.Debug("start load rank data to redis cache...")
		count := 0
		var roles []*db.DBRole
		err := db.DatabaseHandler.QueryAllLimit(&db.DBRole{}, &roles, 10000, func(data interface{}) bool {
			for _, role := range roles {
				for _, statistic := range role.Statistics {
					manager := GetRankManager(statistic.ID)
					if manager != nil {
						if statistic.UpdateTime.After(RankRefreshTime) {
							manager.UpdateRank(role.UserID, int64(statistic.Value))
						} else {
							manager.UpdateRank(role.UserID, 0)
						}
					}
				}

				for _, starStatistic := range role.StarStatistics {
					for _, statistic := range starStatistic.Statistics {
						manager := GetStarRankManager(statistic.ID, starStatistic.Type)
						if manager != nil {
							//manager.UpdateRank(role.UserID, int64(starStatistic.))
							if statistic.UpdateTime.After(RankRefreshTime) {
								manager.UpdateRank(role.UserID, int64(statistic.Value))
							} else {
								manager.UpdateRank(role.UserID, 0)
							}
						}
					}
				}

				//更新昵称
				//cache.UserCache.SetUserNickname(role.UserID, role.NickName)
				//cache.UserCache.SetUserRole(role.UserID, role.ID)
			}
			currLen := len(roles)
			count += currLen
			return currLen == 0
		})
		if err != nil {
			log.Debug("load rank err: %v", err)
		}

		log.Debug("end load rank to redis cache count:%v", count)
	}

	//UpdateDayRank()

	//if center.IsMaster(service.SERVICE_USER) {
	//	//从数据库加载日排行版信息到缓存
	//	//landDayRank := &db.DBDayRank{ID: constant.STATISTIC_TYPE_LAND}
	//	//db.DatabaseHandler.QueryOne(landDayRank)
	//
	//	//LandRankManager.RestoreDayRank(landDayRank.Data)
	//
	//	//peopleDayRank := &db.DBDayRank{ID: constant.STATISTIC_TYPE_PEOPLE}
	//	//db.DatabaseHandler.QueryOne(peopleDayRank)
	//
	//	//PeopleRankManager.RestoreDayRank(peopleDayRank.Data)
	//}

}

//刷新排名
func CleanRank() {
	//每个月1号的1点才需要清除日志
	if center.IsMaster(service.SERVICE_USER) && time.Now().Day() == 1 {
		log.Debug("clean rank start..")
		RankRefreshTime = time.Now()
		SaleRankManager.CleanRank()
		LootFaithRankManager.CleanRank()
		AtkBuildingRankManager.CleanRank()
		LootBelieverRankManager.CleanRank()
		log.Debug("clean rank end..")
	}
}

//func getRoleLand(role *db.DBRole) int64 {
//	for _, resource := range role.Resource {
//		if resource.ID == constant.TID_RESOURCE_IDLELAND {
//			return int64(resource.Limit)
//		}
//	}
//	return 0
//}
//
//func getRolePeople(role *db.DBRole) int64 {
//	var total int64 = 0
//	for _, profession := range role.Profession {
//		total += profession.Value
//	}
//	if role.Unemployed != nil {
//		total += role.Unemployed.Value
//	}
//
//	return total
//}
