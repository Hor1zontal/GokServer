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

import (
	"time"
	"aliens/common/character"
	"gok/service/msg/protocol"
)

//缓存数据管理 上层不关心数据的存取规则

type COMMENT_PROP_KEY string	//评论属性
//
//func time2int32(timestamp time.Time) int32{
//	return int32(timestamp.Unix())
//}

const (
	MOMENTS_SALE_PREFIX = "moments:sale:"
	MOMENTS_USER_PUBLIC_PREFIX = "moments:public:"   //存储用户的朋友圈发布记录
	MOMENTS_USER_RECEIVE_PREFIX = "moments:receive:"  //存储用户的朋友圈接收记录
	MOMENTS_INFO_PREFIX = "moments:info:"  //存储朋友圈的详细信息

	//FRIEND_INFO_PREFIX = "fr_"  //存储好友列表信息
	//FRIEND_REQUEST_INFO_PREFIX = "frq_"  //存储好友申请列表信息
	//FRIEND_DATA_PRIFIX = "fd_" //存储好友详细数据
	FOLLOWER_INFO_PREFIX = "social:follower"  //存储关注的人员信息
	FOLLOWING_INFO_PREFIX = "social:following" //存储被关注的人员信息

	FLAG_LOADFOLLOW string = "flag:follow"  	 	    //标识，是否加载关注信息到内存
	FLAG_LOADMOMENT_TIMELINE string = "flag:timeLine"  	//标识，是否加载朋友圈时间线
	FLAG_LOADMOMENT_INFO string = "flag:momentsInfo"  	 	//标识，是否加载朋友圈详细信息
)


type CommunityCacheManager struct {
	*cacheManager
}

func NewCommunityCacheManager() *CommunityCacheManager {
	return &CommunityCacheManager{
		&cacheManager{},
	}
}


//---------------------------------弹幕接口----------------------

////添加弹幕
//func (this *CommunityCacheManager) AddEventRecordDANMAKU(eventID string, stage string, content string) int {
//	length, _ := this.redisClient.LPush(EVENT_RECORD_DANMAKU_PREFIX + eventID + stage, content)
//	return length
//}
//
////获取所有的弹幕记录
//func (this *CommunityCacheManager) GetAllEventRecordDANMAKU(eventID string, stage string) []string {
//	key := EVENT_RECORD_DANMAKU_PREFIX + eventID + stage
//	return this.redisClient.LRangeAll(key)
//}
//
////获取事件记录弹幕   limit 数量
//func (this *CommunityCacheManager) GetEventRecordDANMAKU(eventID string, stage string, limit int) []string {
//	key := EVENT_RECORD_DANMAKU_PREFIX + eventID + stage
//	return this.redisClient.LRange(key, 0, limit - 1)
//}


//-----------------------------------投票接口——--------------------------------------

//删除投票
//func (this *CommunityCacheManager) DelVote(eventID string) bool {
//	return this.redisClient.DelData(EVENT_VOTE_PREFIX + eventID)
//}
//
////更新投票数量  返回:更新后的总票数
//func (this *CommunityCacheManager) UpdateEventVote(eventID string, option string, voteNum int) int {
//	return this.redisClient.HIncrby(EVENT_VOTE_PREFIX + eventID, option, voteNum)
//}
//
////投票选项是否存在
//func (this *CommunityCacheManager) IsEventVoteOptionExist(eventID string, option string) bool {
//	return this.redisClient.HFieldExists(EVENT_VOTE_PREFIX + eventID, option)
//}
//
////获取事件的投票数量  投票选项和票数的映射数据
//func (this *CommunityCacheManager) GetEventVote(eventID string) map[string]int {
//	return this.redisClient.HGetAllInt(EVENT_VOTE_PREFIX + eventID)
//}
//
////更新用户投票记录
//func (this *CommunityCacheManager) UpdateEventVoteUser(eventID string, userRecord string) bool {
//	return this.redisClient.SAddData(EVENT_VOTE_USER_PREFIX + eventID, userRecord)
//}
//
////用户是否已经投过票了
//func (this *CommunityCacheManager) IsUserVote(eventID string, userRecord string) bool {
//	//list := list.New()
//	//list.inser
//	return this.redisClient.SContains(EVENT_VOTE_USER_PREFIX + eventID, userRecord)
//}

//-----------------------------------朋友圈接口---------------------------------

//设置拍卖物品和朋友圈的映射关系
func (this *CommunityCacheManager) SetSaleMoments(saleID int32, momentsID string) {
	this.redisClient.SetData(MOMENTS_SALE_PREFIX + character.Int32ToString(saleID), momentsID)
}

//获取拍卖物品和朋友圈的隐射关系
func (this *CommunityCacheManager) GetSaleMoments(saleID int32) string {
	return this.redisClient.GetData(MOMENTS_SALE_PREFIX + character.Int32ToString(saleID))
}

func (this *CommunityCacheManager) DeleteSaleMoments(saleID int32) {
	this.redisClient.DelData(MOMENTS_SALE_PREFIX + character.Int32ToString(saleID))
}

//获取朋友圈详细信息
//func (this *CommunityCacheManager) GetMomentsInfo(momentsID string) map[string]string {
//	return this.redisClient.HGetAll(MOMENTS_INFO_PREFIX + momentsID)
//}
//
////保存朋友圈信息
//func (this *CommunityCacheManager) SetMomentsInfo(fields map[interface{}]interface{}) string {
//	momentsID := util.Rand().Hex()
//	if this.redisClient.HMSet(MOMENTS_INFO_PREFIX + momentsID, fields) {
//		return momentsID
//	}
//	return ""
//}

func (this *CommunityCacheManager) GetMomentsInfo(momentsID string) *protocol.MomentInfo {
	data := this.redisClient.GetBytesData(MOMENTS_INFO_PREFIX + momentsID)
	if data == nil || len(data) == 0 {
		return nil
	}
	result := &protocol.MomentInfo{}
	err := result.Unmarshal(data)
	if err != nil {
		return nil
	}
	return result
}

//保存朋友圈信息
func (this *CommunityCacheManager) SetMomentsInfo(info *protocol.MomentInfo) bool {
	data, _ := info.Marshal()
	return this.redisClient.SetData(MOMENTS_INFO_PREFIX + info.Id, data)
}

func (this *CommunityCacheManager) RemoveMomentsInfo(momentsID string) bool {
	return this.redisClient.DelData(MOMENTS_INFO_PREFIX + momentsID)
}


//获取用户发布的所有朋友圈id
//func (this *CommunityCacheManager) GetAllUserPublicMomentIDs(uid string) []interface{}  {
//	return this.redisClient.ZAll(MOMENTS_USER_PUBLIC_PREFIX + uid)
//}
//
//func (this *CommunityCacheManager) GetAllUserReceiveMomentIDs(uid string) []interface{}  {
//	return this.redisClient.ZAll(MOMENTS_USER_RECEIVE_PREFIX + uid)
//}

//获取用户指定期间发布指定数量的朋友圈消息索引信息
func (this *CommunityCacheManager) GetUserPublicMomentIDs(uid int32, limitTime int32, offset int32, count int32) []string  {
	return this.redisClient.ZRevRangeByScoreBeforeLimit(MOMENTS_USER_PUBLIC_PREFIX + character.Int32ToString(uid),
		limitTime, offset, count)
}

//更新自己发布的time line
func (this *CommunityCacheManager) SetUserPublicMomentID(uid int32, momentsID string, createTimestamp time.Time) bool {
	return this.redisClient.ZAdd(MOMENTS_USER_PUBLIC_PREFIX + character.Int32ToString(uid),
		createTimestamp.Unix(), momentsID)
}

func (this *CommunityCacheManager) RemoveUserPublicMomentID(uid string, momentsID string) bool {
	return this.redisClient.ZRem(MOMENTS_USER_PUBLIC_PREFIX + uid, momentsID)
}

//获取用户指定期间收到的指定数量朋友圈消息索引信息
func (this *CommunityCacheManager) GetUserReceiveMomentIDs(uid int32, limitTime int32, offset int32, count int32) []string {
	return this.redisClient.ZRevRangeByScoreBeforeLimit(MOMENTS_USER_RECEIVE_PREFIX + character.Int32ToString(uid),
		limitTime, offset, count)
}

func (this *CommunityCacheManager) RemoveUserReceiveMomentID(uid string, momentsID string) bool {
	return this.redisClient.ZRem(MOMENTS_USER_RECEIVE_PREFIX + uid, momentsID)
}

//获取接收到的time line
func (this *CommunityCacheManager) SetUserReceiveMomentID(uid int32, momentsID string, createTimestamp time.Time) bool {
	return this.redisClient.ZAdd(MOMENTS_USER_RECEIVE_PREFIX + character.Int32ToString(uid),
		createTimestamp.Unix(), momentsID)
}

//删除用户发布的信息
func (this *CommunityCacheManager) RemoveUserReceive(uid string, receiveID string) bool {
	momentsID := this.redisClient.ZAll(MOMENTS_USER_PUBLIC_PREFIX + receiveID)
	if momentsID == nil || len(momentsID) == 0 {
		return false
	}
	return this.redisClient.ZRems(MOMENTS_USER_RECEIVE_PREFIX + uid, momentsID)
}

//新增朋友圈信息
//func (this *CacheManager) AddMomentsInfo() string {
//	momentsID := "uuid"
//	json.Marshal()
//	this.redisClient.HSet(momentsID)
//}


//新增事件评论
//eventID 事件id  data评论信息(昵称 时间戳 内容)
//func (this *CacheManager) AddEventComment(eventID int32, data string, timestamp int32) bool {
//	//采用uuid
//	//自动生成评论id
//	commentID := ""
//	if (this.redisClient.HSet(EVENT_RECORD_COMMENT_PREFIX + character.Int32ToString(eventID), commentID, data)) {
//		return this.redisClient.ZAdd(EVENT_RECORD_COMMENT_PREFIX + character.Int32ToString(eventID), timestamp , commentID)
//	}
//	return false
//}

//获取事件评论信息
//func (this *CacheManager) GetEventComment(eventID int32, commentID int32) string {
//	//采用uuid
//	//自动生成评论id
//	return this.redisClient.HGet(COMMENT_EVENT_PREFIX + character.Int32ToString(eventID),  commentID )
//}

//------------好友信息索引----------------
//func (this *CommunityCacheManager) AddNicknameUIDMapping(nickname string, uid int32) bool {
//	return this.redisClient.HSet(nickname, character.Int32ToString(uid), time.Now().Unix())
//}
//
//func (this *CommunityCacheManager) GetAllUIDByNickname(nickname string) map[string]int64 {
//	return this.redisClient.HGetAllInt64(nickname)
//}



//------------好友接口--------------------

//func getFriendKey(id string) string {
//	return FRIEND_INFO_PREFIX + id
//}
//
//func getFriendReqKey(id string) string {
//	return FRIEND_REQUEST_INFO_PREFIX + id
//}
//
//func getFriendDataKey(id string) string {
//	return FRIEND_DATA_PRIFIX + id
//}
//
//func (this *CommunityCacheManager) GetFriendArray(id string) []int32 {
//	results := []int32{}
//	friends := this.redisClient.HGetAllInt64(getFriendKey(id))
//	if friends == nil {
//		return	results
//	}
//
//	for friendID, _ := range friends {
//		results = append(results, character.StringToInt32(friendID))
//	}
//	return results
//}
//
////获取所有好友
//func (this *CommunityCacheManager) GetFriends(id string) map[string]int64 {
//	return this.redisClient.HGetAllInt64(getFriendKey(id))
//}
//
////添加好友
//func (this *CommunityCacheManager) AddFriend(id string, friendID string){
//	this.redisClient.HSet(getFriendKey(id), friendID, time.Now().Unix())
//}
//
//
//func (this *CommunityCacheManager) ExistFriend(id string, friendID string) bool {
//	return this.redisClient.HFieldExists(getFriendKey(id), friendID)
//}
//
//
////删除好友
//func (this *CommunityCacheManager) DelFriend(id string, friendID string) bool {
//	return this.redisClient.HDel(getFriendKey(id), friendID) != nil
//}
//
////更新好友自定义数据
//func (this *CommunityCacheManager) SetFriendData(id string, data interface{}){
//	this.redisClient.SetData(getFriendDataKey(id), data)
//}
//
////获取好友自定义数据
//func (this *CommunityCacheManager) GetFriendData(id string) string {
//	return this.redisClient.GetData(getFriendDataKey(id))
//}

//-------------------------好友申请-----------------------------

//获取所有的好友请求信息 好友id-时间戳
//func (this *CommunityCacheManager) GetFriendRequests(id string) map[string]int64 {
//	return this.redisClient.HGetAllInt64(getFriendReqKey(id))
//}
//
////是否存在好友请求
//func (this *CommunityCacheManager) ExistFriendRequest(id string, friendID string) bool {
//	return this.redisClient.HFieldExists(getFriendReqKey(id), friendID)
//}
//
////添加好友申请
//func (this *CommunityCacheManager) AddFriendRequest(id string, friendID string){
//	this.redisClient.HSet(getFriendReqKey(id), friendID, time.Now().Unix())
//}
//
////删除好友申请
//func (this *CommunityCacheManager) DelFriendRequest(id string, friendID string) bool {
//	return this.redisClient.HDel(getFriendReqKey(id), friendID) != nil
//}


//----------------社交关系--------------------


//-----------------关注--------------------

//关注数据key
func getFollowerKey(id string) string {
	return FOLLOWER_INFO_PREFIX + id
}

//粉丝数据key
func getFollowingKey(id string) string {
	return FOLLOWING_INFO_PREFIX + id
}


//获取关注列表
func (this *CommunityCacheManager) GetFollowers(id string) map[string]int64 {
	return this.redisClient.HGetAllInt64(getFollowerKey(id))
}

func (this *CommunityCacheManager) GetFollowerCount(id string) int {
	len, _ := this.redisClient.HLen(getFollowerKey(id))
	return len
}

//是否在关注列表中
func (this *CommunityCacheManager) ExistFollower(id string, followerID string) bool {
	return this.redisClient.HFieldExists(getFollowerKey(id), followerID)
}

//添加关注人员
func (this *CommunityCacheManager) AddFollower(id string, followerID string){
	this.AddFollower1(id, followerID, time.Now())
}

func (this *CommunityCacheManager) AddFollower1(id string, followerID string, addTime time.Time){
	this.redisClient.HSet(getFollowerKey(id), followerID, addTime.Unix())
	this.redisClient.HSet(getFollowingKey(followerID), id, time.Now().Unix())
}

//删除关注人员
func (this *CommunityCacheManager) DelFollower(id string, followerID string) {
	this.redisClient.HDel(getFollowerKey(id), followerID)
	this.redisClient.HDel(getFollowingKey(followerID), id)
}

//获取被关注列表
func (this *CommunityCacheManager) GetFollowings(id string) map[string]int64 {
	return this.redisClient.HGetAllInt64(getFollowingKey(id))
}

//是否在被关注列表中
func (this *CommunityCacheManager) ExistFollowing(id string, FollowingID string) bool {
	return this.redisClient.HFieldExists(getFollowingKey(id), FollowingID)
}