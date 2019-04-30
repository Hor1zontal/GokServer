/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/3/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package cache

//缓存数据管理 上层不关心数据的存取规则



//存储事件记录数据
//func (this *UserCacheManager) SetEventRecord(eventID string, data string) bool {
//	return this.redisClient.SetData(EVENT_RECORD_PREFIX + eventID , data)
//}

////获取事件记录数据
//func (this *UserCacheManager) GetEventRecord(eventID string) string {
//	return this.redisClient.GetData(EVENT_RECORD_PREFIX + eventID)
//}

//-----------------------------------投票接口——--------------------------------------

//删除投票
//func (this *UserCacheManager) DelVote(eventID string) bool {
//	return this.redisClient.DelData(EVENT_VOTE_PREFIX + eventID)
//}
//
////更新投票数量  返回:更新后的总票数
//func (this *UserCacheManager) UpdateEventVote(eventID string, option string, voteNum int) int {
//	return this.redisClient.HIncrby(EVENT_VOTE_PREFIX + eventID, option, voteNum)
//}
//
////投票选项是否存在
//func (this *UserCacheManager) IsEventVoteOptionExist(eventID string, option string) bool {
//	return this.redisClient.HFieldExists(EVENT_VOTE_PREFIX + eventID, option)
//}
//
////获取事件的投票数量  投票选项和票数的映射数据
//func (this *UserCacheManager) GetEventVote(eventID string) map[string]int {
//	return this.redisClient.HGetAllInt(EVENT_VOTE_PREFIX + eventID)
//}
//
////更新用户投票记录
//func (this *UserCacheManager) UpdateEventVoteUser(eventID string, userRecord string) bool {
//	return this.redisClient.SAddData(EVENT_VOTE_USER_PREFIX + eventID, userRecord)
//}
//
////用户是否已经投过票了
//func (this *UserCacheManager) IsUserVote(eventID string, userRecord string) bool {
//	//list := list.New()
//	//list.inser
//	return this.redisClient.SContains(EVENT_VOTE_USER_PREFIX + eventID, userRecord)
//}

