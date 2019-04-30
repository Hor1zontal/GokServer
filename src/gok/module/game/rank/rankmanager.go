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
	commoncache "aliens/common/cache"
	"gok/module/game/cache"
	"gok/service/msg/protocol"
)

type RankManager struct {
	key    string //当前排名数据的id
	//dayKey string //日排名数据的id

	topRank  int32 //排行版数量
	rankType int32 //排行版类型
}

//获取当前排名
func (this *RankManager) UpdateRank(uid int32, value int64) {
	cache.UserCache.UpdateRank(this.key, uid, value)
}

//清除排名
func (this *RankManager) CleanRank() {
	cache.UserCache.DeleteData(this.key)
}

func (this *RankManager) GetRankTotalNum(key string) int32 {
	return int32(cache.UserCache.GetRankTotalCount(key))
}

////更新日排名
//func (this *RankManager) UpdateDayRank() string {
//	data := cache.UserCache.CopyData(this.key, this.dayKey)
//
//	dbDayRank := &db.DBDayRank{
//		ID:              this.rankType,
//		Data:            data,
//		UpdateTimestamp: time.Now(),
//	}
//	db.DatabaseHandler.ForceUpdateOne(dbDayRank)
//	return data
//}

//更新日排名
//func (this *RankManager) RestoreDayRank(restoreData string) {
//	cache.UserCache.RestoreData(this.dayKey, restoreData)
//}

//获取前10排名
func (this *RankManager) GetTopRank() []commoncache.Rank {
	return cache.UserCache.GetTopRank(this.key, this.topRank)
}

//获取当前排名
func (this *RankManager) GetCurrentRank(uid int64) int32 {
	return cache.UserCache.GetRank(this.key, uid)
}

//获取当前分数
func (this *RankManager) GetCurrentScore(uid int64) int64 {
	return cache.UserCache.GetScore(this.key, uid)
}

//获取日排名
//func (this *RankManager) GetDayRank(uid int64) int32 {
//	return cache.UserCache.GetRank(this.dayKey, uid)
//}

func (this *RankManager) GetUserRankData(uid int32, currRank int32, score int64) *protocol.Rank {
	//nickname := cache.UserCache.GetUserNickname(uid)
	//avatar := cache.UserCache.GetUserAvatar(uid)
	//昨天的日排行,排行数据redis是从0开始，需要+1
	//dayRank := this.GetDayRank(uid) + 1
	currRank += 1
	//var delta int32 = 0
	//日排行没有统计到，没有名次变化
	//if dayRank <= 0 {
	//	delta = 0
	//} else {
	//	delta = dayRank - currRank
	//}

	//isNew := false

	//上次在排名外，这次进排名了需要更新new标识
	//if (dayRank <= 0 || dayRank > this.topRank) && currRank <= this.topRank {
	//	isNew = true
	//}

	return &protocol.Rank{
		Uid: uid,
		RankNum:  int32(currRank),
		//Nickname: nickname),
		//Avatar: avatar),
		Value:    score,
		//Delta:    delta),
		//IsNew:    isNew),
	}
}
