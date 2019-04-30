/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/4/19
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package cache

import (
	"aliens/log"
	"aliens/common/cache"
	"math"
)

const (

	//交易排行
	RANK_SALE_KEY string = "rank:sale"
	//DAY_RANK_PEOPLE_KEY string = "dpr_"

	//抢夺信仰
	RANK_LOOT_FAITH_KEY string = "rank:lootFaith"
	//DAY_RANK_LAND_KEY string = "dlr_"

	//攻击建筑
	RANK_ATK_BUIDING_KEY string = "rank:attack"

	//抢夺信徒
	RANK_LOOT_BELIEVER_KEY string = "rank:lootBeliever"

	//星球通关时长
	RANK_STAR_ONLINE_KEY string = "rank:onlineTime_"

	FLAG_LOADRANK  string = "flag:rank" //标识，是否加载排行版数据到缓存
)

func (this *UserCacheManager) GetStarOnlineKey() string {
	return RANK_STAR_ONLINE_KEY
}

//更新排行信息
func (this *UserCacheManager) DeleteRank(key string, id interface{}) bool {
	return this.redisClient.ZRem(key, id)
}

func (this *UserCacheManager) DeleteData(key string) bool {
	return this.redisClient.DelData(key)
	//将当前数据更新到日排行
	//dumpData := this.redisClient.Dump(dumpKey)
	//
	//this.redisClient.Restore(restoreKey, dumpData)
	//return dumpData
}

func (this *UserCacheManager) GetRankTotalCount(key string) int{
	return this.redisClient.ZCount(key, 0, math.MaxInt32)
}

//更新排行信息
func (this *UserCacheManager) UpdateRank(key string, id interface{}, num int64) {
	this.redisClient.ZAdd(key, num, id)
}

//获取排行版前信息
func (this *UserCacheManager) GetTopRank(key string, count int32) []cache.Rank {
	result, err := this.redisClient.ZRevRangeWithScore(key, 0, count-1)
	if (err != nil) {
		log.Error("call GetTopRank error : %v", err)
	}
	return result
}

//获取排行版尾信息
func (this *UserCacheManager) GetTailRank(key string, count int32) []cache.Rank {
	return this.redisClient.ZRangeWithScore(key, 0, count-1)
}

//获取指定的排名
func (this *UserCacheManager) GetRank(key string, id int64) int32 {
	return int32(this.redisClient.ZRevRank(key, id))
}

func (this *UserCacheManager) GetScore(key string, id int64) int64 {
	return this.redisClient.ZScore(key, id)
}
