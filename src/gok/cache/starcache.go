/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package cache

import (
	"aliens/common/character"
)

const (


	STAR_KEY_PREFIX                 = "star:"
	USER_STAR_KEY_PREFIX            = "userActiveStar:"
	LAST_USER_STAR_KEY_PREFIX		= "userStarNode:"
	STAR_REPAIR_BUILD_KEY_PREFIX	= "starHelpRepair:"	//求助修理建筑

	//LAST_PROP_NODE				   string = "node" 		   //最后处理的节点
	STAR_PROP_BUILDING_TOTAL_LEVEL string = "bu_all_level"    //星球建筑总等级
	STAR_PROP_BELIEVER_COUNT       string = "be_count"        //星球信徒数量
	STAR_PROP_BELIEVER_TOTAL_LEVEL string = "be_all_level"    //星球信徒总等级
	STAR_PROP_CIVIL_LEVEL string = "civil_level"    //文明度等级
	STAR_PROP_OWNER                string = "owner"           //星球所属用户
	STAR_PROP_TYPE                 string = "type"            //星球类型
	//STAR_PROP_ACTIVE_GROUP_ID      string = "active_group_id" //正在激活的圣物组合ID
	//STAR_PROP_ACTIVE_GROUP_INDEX   string = "active_group_index" //正在激活的圣物组合序号
	STAR_PROP_MUTUAL_TIMES         string = "mutual_times"
	STAR_PROP_BE_MUTUAL_TIMES      string = "be_mutual_times"
	STAR_PROP_BUILDING_MAX_LEVEL   string = "bu_max_level"     //建筑过去最大总等级

	UPROP_STAR string = "star"      	//用户当前星球id
	UPROP_STARTYPE string = "tstar"      	//用户当前的星球类型

	FLAG_LOADSTAR         string = "flag:ustar"         //标识，是否加载星球缓存到redis
)

type StarCacheManager struct {
	*cacheManager
}

func NewStarCacheManager() *StarCacheManager {
	return &StarCacheManager{
		&cacheManager{},
	}
}
//
func getStarKey(starID int32) string {
	return STAR_KEY_PREFIX + character.Int32ToString(starID)
}

func getUserStarKey(uid int32) string {
	return USER_STAR_KEY_PREFIX + character.Int32ToString(uid)
}

func getLastUserStarKey(uid int32) string {
	return LAST_USER_STAR_KEY_PREFIX + character.Int32ToString(uid)
}

func getStarHelpRepairKey(starID int32) string {
	return STAR_REPAIR_BUILD_KEY_PREFIX + character.Int32ToString(starID)
}

func (this *StarCacheManager) SetStarBuildRepair(starID int32, buildType int32, helperID int32) bool {
	return this.redisClient.HSet(getStarHelpRepairKey(starID), character.Int32ToString(buildType), helperID)
}

func (this *StarCacheManager) GetAllStarBuildRepair(starID int32) map[string]string {
	return this.redisClient.HGetAll(getStarHelpRepairKey(starID))
}

func (this *StarCacheManager) GetStarBuildRepair(starID int32, buildType int32) string {
	return this.redisClient.HGet(getStarHelpRepairKey(starID), character.Int32ToString(buildType))
}

func (this *StarCacheManager) ExistStarHelpRepair(starID int32, buildType int32) bool {
	return this.redisClient.HFieldExists(getStarHelpRepairKey(starID), character.Int32ToString(buildType))
}

func (this *StarCacheManager) DelStarHelpRepair(starID int32, buildType int32) error {
	return this.redisClient.HDel(getStarHelpRepairKey(starID), character.Int32ToString(buildType))
}

func (this *StarCacheManager) SetUserIDNode(userID int32, node string) bool {
	return this.redisClient.SetData(getLastUserStarKey(userID), node)
}

func (this *StarCacheManager) GetUserIDNode(userID int32) string {
	return this.redisClient.GetData(getLastUserStarKey(userID))
}

///设置用户星球总等级
func (this *StarCacheManager) SetBuildingAllLevel(starID int32, level int32) bool {
	return this.redisClient.HSet(getStarKey(starID), STAR_PROP_BUILDING_TOTAL_LEVEL, level)
}

//获取用户星球总等级
func (this *StarCacheManager) GetBuildingAllLevel(starID int32) int32 {
	return this.redisClient.HGetInt32(getStarKey(starID), STAR_PROP_BUILDING_TOTAL_LEVEL)
}

func (this *StarCacheManager) SetBuildingExMaxLevel(starID int32, level int32) bool {
	return this.redisClient.HSet(getStarKey(starID), STAR_PROP_BUILDING_MAX_LEVEL, level)
}

func (this *StarCacheManager) GetBuildingExMaxLevel(starID int32) int32 {
	return this.redisClient.HGetInt32(getStarKey(starID), STAR_PROP_BUILDING_MAX_LEVEL)
}

//设置星球信徒总数量
func (this *StarCacheManager) SetBelieverCount(starID int32, count int32) bool {
	return this.redisClient.HSet(getStarKey(starID), STAR_PROP_BELIEVER_COUNT, count)
}

//获取星球信徒总数量
func (this *StarCacheManager) GetBelieverCount(starID int32) int32 {
	return this.redisClient.HGetInt32(getStarKey(starID), STAR_PROP_BELIEVER_COUNT)
}

//设置星球信徒总等级
func (this *StarCacheManager) SetBelieverTotalLevel(starID int32, count int32) bool {
	return this.redisClient.HSet(getStarKey(starID), STAR_PROP_BELIEVER_TOTAL_LEVEL, count)
}

//获取星球信徒总等级
func (this *StarCacheManager) GetBelieverTotalLevel(starID int32) int32 {
	return this.redisClient.HGetInt32(getStarKey(starID), STAR_PROP_BELIEVER_TOTAL_LEVEL)
}

//设置星球信徒总等级
func (this *StarCacheManager) SetCivilLevel(starID int32, count int32) bool {
	return this.redisClient.HSet(getStarKey(starID), STAR_PROP_CIVIL_LEVEL, count)
}

//获取星球信徒总等级
func (this *StarCacheManager) GetCivilLevel(starID int32) int32 {
	return this.redisClient.HGetInt32(getStarKey(starID), STAR_PROP_CIVIL_LEVEL)
}

//设置星球所属的用户id
func (this *StarCacheManager) SetOwner(starID int32, uid int32) bool {
	return this.redisClient.HSet(getStarKey(starID), STAR_PROP_OWNER, uid)
}

//获取星球所属的用户id
func (this *StarCacheManager) GetOwner(starID int32) int32 {
	return this.redisClient.HGetInt32(getStarKey(starID), STAR_PROP_OWNER)
}

//设置星球类型
func (this *StarCacheManager) SetType(starID int32, starType int32) bool {
	return this.redisClient.HSet(getStarKey(starID), STAR_PROP_TYPE, starType)
}

//获取星球类型
func (this *StarCacheManager) GetType(starID int32) int32 {
	return this.redisClient.HGetInt32(getStarKey(starID), STAR_PROP_TYPE)
}

////设置当前正在尝试的圣物组合ID
//func (this *StarCacheManager) SetCurrentGroupID(starID int32, groupID int32) bool {
//	return this.redisClient.HSet(getStarKey(starID), STAR_PROP_ACTIVE_GROUP_ID, groupID)
//}
//
////获取当前正在尝试的圣物组合ID
//func (this *StarCacheManager) GetCurrentGroupID(starID int32) int32 {
//	return this.redisClient.HGetInt32(getStarKey(starID), STAR_PROP_ACTIVE_GROUP_ID)
//}

////设置当前正在尝试的圣物组合的序号
//func (this *StarCacheManager) SetCurrentGroupIndex(starID int32, index int32) bool {
//	return this.redisClient.HSet(getStarKey(starID), STAR_PROP_ACTIVE_GROUP_INDEX, index)
//}
//
////获取当前正在尝试的圣物组合的序号
//func (this *StarCacheManager) GetCurrentGroupIndex(starID int32) int32 {
//	return this.redisClient.HGetInt32(getStarKey(starID), STAR_PROP_ACTIVE_GROUP_INDEX)
//}

//设置主动交互的次数
func (this *StarCacheManager) SetMutualTimes(starID int32, mutual int32) bool {
	return this.redisClient.HSet(getStarKey(starID), STAR_PROP_MUTUAL_TIMES, mutual)
}

//获取主动交互的次数
func (this *StarCacheManager) GetMutualTimes(starID int32) int32 {
	return this.redisClient.HGetInt32(getStarKey(starID), STAR_PROP_MUTUAL_TIMES)
}

//设置被动交互的次数
func (this *StarCacheManager) SetBeMutualTimes(starID int32, beMutual int32) bool {
	return this.redisClient.HSet(getStarKey(starID), STAR_PROP_BE_MUTUAL_TIMES, beMutual)
}

//获取被动交互的次数
func (this *StarCacheManager) GetBeMutualTimes(starID int32) int32 {
	return this.redisClient.HGetInt32(getStarKey(starID), STAR_PROP_BE_MUTUAL_TIMES)
}

//获取星球所有属性
func (this *StarCacheManager) GetStarProps(starID int32) map[string]int {
	return this.redisClient.HGetAllInt(getStarKey(starID))
}

//设置用户星球信息
func (this *StarCacheManager) SetUserActiveStar(uid int32, starId int32) {
	key := getUserStarKey(uid)
	this.redisClient.HSet(key, UPROP_STAR, starId)
}

//获取用户当前星球id
func (this *StarCacheManager) GetUserActiveStar(uid int32) int32  {
	return this.redisClient.HGetInt32(getUserStarKey(uid), UPROP_STAR)
}

//设置用户星球当前星球类型
func (this *StarCacheManager) SetUserActiveStarType(uid int32, starType int32) {
	key := getUserStarKey(uid)
	this.redisClient.HSet(key, UPROP_STARTYPE, starType)
}

//获取用户当前星球类型
func (this *StarCacheManager) GetUserActiveStarType(uid int32) int32  {
	return this.redisClient.HGetInt32(getUserStarKey(uid), UPROP_STARTYPE)
}



//设置法力回复BUFF间隔时间
func (this *StarCacheManager) SetBuffMANAInterval(uid int32, interval int32) bool {
	return this.redisClient.HSet(GetUserKey(uid), UPROP_BUFF_MANA_INTERVAL, interval)
}

//获取法力回复BUFF间隔时间
func (this *StarCacheManager) GetBuffMANAInterval(uid int32) int32 {
	return this.redisClient.HGetInt32(GetUserKey(uid), UPROP_BUFF_MANA_INTERVAL)
}

//设置偷取圣物buff加成概率
func (this *StarCacheManager) SetBuffRelicSteal(uid int32, prob float64) bool {
	return this.redisClient.HSet(GetUserKey(uid), UPROP_BUFF_RELIC_STEAL, prob)
}

//获取偷取圣物buff加成概率
func (this *StarCacheManager) GetBuffRelicSteal(uid int32) float64 {
	return this.redisClient.HGetFloat64(GetUserKey(uid), UPROP_BUFF_RELIC_STEAL)
}



