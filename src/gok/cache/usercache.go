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
	"aliens/common/character"
	"aliens/log"
	"github.com/gogo/protobuf/proto"
	"time"
)

//缓存数据管理 上层不关心数据的存取规则
const (
	USER_KEY_PREFIX string = "user:"
	USERNAME_KEY_PREFIX string = "username:"

	USER_KEY_DAY_GIFT string = "user_day_gift:"

	USER_TEST_KEY_PREFIX string = "usertestID"

	keyWxServiceOUMapping string = "wx:service_ou:" //公众号openid和unionid的映射关系
	keyWxServiceUOMapping string = "wx:service_uo:" //公众号unionid和openid的映射关系

	keyWxGameUOMapping string = "wx:game_uu:" //小游戏unionid和uid的映射关系

	UPROP_DESC string = "desc"          //用户签名
	UPROP_STATUS string = "status"      //用户状态 0正常 1封号
	UPROP_NICKNAME string = "nname"     //用户昵称
	UPROP_ICON string = "icon"      	//图标
	UPROP_ONLINE string = "online"		//用户是否登录
	UPROP_ONLINE_TIME string = "otime"  //用户上次的登录时间戳
	//UPROP_INTO_EVENT string = "inevent" //用户当前进入的事件副本
	UPROP_FAITH       string = "faith"   //用户信仰值
	UPROP_POWER       string = "power"   //用户法力值
	UPROP_POWER_LIMIT string = "p_limit" //用户法力值上限
	UPROP_LEVEL       string = "level"   //用户神力值
	UPROP_AVATAR      string = "avatar"  //用户头像
	UPROP_UNIONID     string = "unionID" //用户微信unionID
	UPROP_HISTORY_ATK string = "h_atk"   //星球被攻击的历史记录

	UPROP_PUBLIC_ITEM string = "pitem"      	//用户发布的物品
	UPROP_PUBLIC_TIMESTAMP string = "ptstamp"   //用户发布物品的时间戳

	UPROP_SALE string = "h_sale"  //用户挂售的圣物

	//UPROP_BUFF_MANA_LIMIT string = "buffml"      	//用户当前的BUFF法力上限加成
	UPROP_BUFF_MANA_INTERVAL string = "buffmi"      	//用户当前的BUFF法力回复间隔
	UPROP_BUFF_RELIC_STEAL string = "buffrs"      	//偷取圣物的加成buff

	UPROP_TOKEN string = "token:"      	//登录令牌
	USER_HELPITEM_WECHAT_KEY = "itemHelpWeChat:" //微信好友帮助圣物求助

	FLAG_LOADUSER string = "flag:user"   	//标识，是否加载用户数据到缓存
	FLAG_LOADROLE string = "flag:role"  	 	//标识，是否加载角色数据到缓存
	MUTUAL_FLAG_PREFIX = "flag:mutual"

	//MAX_USER_ID string = "mxuid"		//当前服务器最大的用户id

	//USER_MESSAGE_KEY_PREFIX     string = "um_"



)

type UserCacheManager struct {
	*cacheManager
}

func NewUserCacheManager() *UserCacheManager {
	return &UserCacheManager{
		&cacheManager{},
		//NewCommunityCacheManager(),
	}
}


func GetUserKey(uid int32) string {
	return USER_KEY_PREFIX + character.Int32ToString(uid)
}

func GetUserDayGiftKey(uid int32) string {
	return USER_KEY_DAY_GIFT + character.Int32ToString(uid)
}

////获取用户挂售商品
//func (this *UserCacheManager) GetUserSale(uid int32) int32 {
//	return this.redisClient.HGetInt32(GetUserKey(uid), UPROP_SALE)
//}
//
////更新用户挂售商品
//func (this *UserCacheManager) SetUserSale(uid int32, itemID int32) bool {
//	return this.redisClient.HSet(GetUserKey(uid), UPROP_SALE, itemID)
//}
//
//func (this *UserCacheManager) CleanUserSale(uid int32) bool {
//	return this.redisClient.HDel(GetUserKey(uid), UPROP_SALE) != nil
//}

//设置用户id上限值
//func (this *UserCacheManager) SetUIDLimit(uid int)  {
//	this.redisClient.SetData(MAX_USER_ID, uid)
//}
//
////获取用户id上限值
//func (this *UserCacheManager) GetUIDLimit() int {
//	return this.redisClient.GetDataInt32(MAX_USER_ID)
//}


//订阅过期
func (this *UserCacheManager) SubscribeExpire(callback func(pattern, channel, value string)) {
	this.redisClient.PSubscribe(callback,"__keyevent@0__:expired")
}

func (this *UserCacheManager) SetCustomFlag(key string, value int) bool {
	return this.redisClient.SetData("cflag:" + key, value)
}

func (this *UserCacheManager) GetCustomFlag(key string) int {
	return this.redisClient.GetDataInt32("cflag:" + key)
}

func (this *UserCacheManager) SetExpireData(key string, value interface{}, seconds int) bool {
	return this.redisClient.SetExpireData(key, value, seconds)
}

//更新公众号unionid和openid的隐射关系
func (this *UserCacheManager) SetWxServiceMapping(openID string, unionID string) {
	this.redisClient.SetData(keyWxServiceOUMapping + openID, unionID)
	this.redisClient.SetData(keyWxServiceUOMapping + unionID, openID)
}

//清除微信公众号unionid和openid的映射关系
func (this *UserCacheManager) CleanWxServiceMapping(openID string) {
	unionID := this.GetWxServiceOUMapping(openID)
	if unionID != "" {
		this.redisClient.DelData(keyWxServiceUOMapping + unionID)
	}
	this.redisClient.DelData(keyWxServiceOUMapping + openID)
}

//通过微信公众号的openid获取unionid
func (this *UserCacheManager) GetWxServiceOUMapping(openID string) string {
	return this.redisClient.GetData(keyWxServiceOUMapping + openID)
}

//通过微信公众号的unionid获取openid
func (this *UserCacheManager) GetWxServiceUOMapping(unionID string) string {
	return this.redisClient.GetData(keyWxServiceUOMapping + unionID)
}

func (this *UserCacheManager) GetUnionIDByUserID(uid int32) string {
	return this.redisClient.HGet(GetUserKey(uid), UPROP_UNIONID)
}

func (this *UserCacheManager) GetUserIDByUnionID(unionID string) int32 {
	return int32(this.redisClient.GetDataInt32(keyWxGameUOMapping + unionID))
}

func (this *UserCacheManager) SetUserIDUnionIDMapping(uid int32, unionID string) bool {
	this.redisClient.SetData(keyWxGameUOMapping + unionID, uid)
	return this.redisClient.HSet(GetUserKey(uid), UPROP_UNIONID, unionID)
}

func (this *UserCacheManager) CleanUserIDUnionIDMapping(uid int32, unionID string) bool {
	ret := this.redisClient.DelData(keyWxGameUOMapping + unionID)
	err := this.redisClient.HDel(GetUserKey(uid), UPROP_UNIONID)
	if err != nil {
		log.Error("userCache error:%v", err.Error())
	}
	if err == nil && ret {
		return true
	}
	return false
}

//用户是否关注公众号
func (this *UserCacheManager) IsUserSubscribe(uid int32) bool {
	unionID := this.GetUnionIDByUserID(uid)
	if unionID == "" {
		return false
	}
	//通过unionID获取服务号的openid
	serviceOpenID := this.GetWxServiceUOMapping(unionID)
	return serviceOpenID != ""
}

func GetTokenKey(uid int32) string {
	return UPROP_TOKEN + character.Int32ToString(uid)
}

//设置用户会话token
func (this *UserCacheManager) SetUserToken(uid int32, token string, expire int) bool {
	key := GetTokenKey(uid)
	result := this.redisClient.SetData(key, token)
	this.redisClient.Expire(key, expire)
	return result
}

//获取用户会话token
func (this *UserCacheManager) GetUserToken(uid int32) string {
	return this.redisClient.GetData(GetTokenKey(uid))
}

//用户名是否存在
func (this *UserCacheManager) IsUsernameExist(username string) bool {
	return this.GetUidByUsername(username) != 0
}

func (this *UserCacheManager) SetUsernameUidMapping(username string, uid int32) bool {
	return this.redisClient.SetData(USERNAME_KEY_PREFIX + username, uid)
}

func (this *UserCacheManager) GetUidByUsername(username string) int32 {
	return int32(this.redisClient.GetDataInt32(USERNAME_KEY_PREFIX + username))
}

func (this *UserCacheManager) DelUidByUsername(username string) bool {
	return this.redisClient.DelData(USERNAME_KEY_PREFIX + username)
}

//获取用户所有信息数据
func (this *UserCacheManager) HSetUser(uid int32, data interface{}) {
	this.redisClient.HSetData(GetUserKey(uid), data)
}

//设置用户所有信息数据
func (this *UserCacheManager) HGetUser(uid int32, data interface{}) {
	this.redisClient.HGetData(GetUserKey(uid), data)
}

//用户是否存在
func (this *UserCacheManager) IsUserExist(uid int32) bool {
	result, _ := this.redisClient.Exists(GetUserKey(uid))
	return result
}

//用户是否在线
//func (this *UserCacheManager) IsUserOnline(uid int32) bool {
//	return this.GetUserAttrBool(uid, UPROP_ONLINE)
//}

//设置用户是否在线
//func (this *UserCacheManager) SetUserOnline(uid int32, online bool) bool {
//	return this.SetUserAttr(uid, UPROP_ONLINE, online)
//}

func (this *UserCacheManager) GetUserOnlineTimestamp(uid int32) int64 {
	return this.GetUserAttrInt64(uid, UPROP_ONLINE_TIME)
}

func (this *UserCacheManager) SetUserOnlineTimestamp(uid int32, timestamp time.Time) bool {
	return this.SetUserAttr(uid, UPROP_ONLINE_TIME, timestamp.Unix())
}

////用户是否在事件中
//func (this *UserCacheManager) IsInEvent(uid int32, eventID int32) bool {
//	return this.GetUserAttrInt32(uid, UPROP_INTO_EVENT) == eventID
//}
//
////用户进入事件
//func (this *UserCacheManager) InEvent(uid int32, eventID int32) bool {
//	return this.SetUserAttr(uid, UPROP_INTO_EVENT, eventID)
//}
//
////用户离开事件
//func (this *UserCacheManager) OutEvent(uid int32, eventID int32) bool {
//	if (this.InEvent(uid, eventID)) {
//		return this.SetUserAttr(uid, UPROP_INTO_EVENT, 0)
//	}
//	return true
//}

//更新用户当前星球被攻击的次数
func (this *UserCacheManager) UpdateATKHistory(uid int32, num int32) bool {
	this.redisClient.HIncrby(GetUserKey(uid), UPROP_HISTORY_ATK, int(num))
	return true
}

//获取用户当前被攻击的次数
func (this *UserCacheManager) GetATKHistory(uid int32) int32 {
	return this.GetUserAttrInt32(uid, UPROP_HISTORY_ATK)
}

//设置用户属性
func (this *UserCacheManager) SetUserAttr(uid int32, propKey string, value interface{}) bool {
	return this.redisClient.HSet(GetUserKey(uid), propKey, value)
}

func (this *UserCacheManager) SetUserAvatar(uid int32, avatar string) bool {
	return this.redisClient.HSet(GetUserKey(uid), UPROP_AVATAR, avatar)
}

func (this *UserCacheManager) GetUserAvatar(uid int32) string {
	return this.redisClient.HGet(GetUserKey(uid), UPROP_AVATAR)
}

//设置用户会话token
func (this *UserCacheManager) SetUserNickname(uid int32, nickname string) bool {
	return this.redisClient.HSet(GetUserKey(uid), UPROP_NICKNAME, nickname)
}

//获取用户昵称
func (this *UserCacheManager) GetUserNickname(uid int32) string {
	return this.redisClient.HGet(GetUserKey(uid), UPROP_NICKNAME)
}

func (this *UserCacheManager) GetUserDesc(uid int32) string {
	return this.redisClient.HGet(GetUserKey(uid), UPROP_DESC)
}

func (this *UserCacheManager) SetUserDesc(uid int32, desc string) bool {
	return this.redisClient.HSet(GetUserKey(uid), UPROP_DESC, desc)
}

//func (this *UserCacheManager) GetPublicItem(uid int32) int32 {
//	return this.GetUserAttrInt32(uid, UPROP_PUBLIC_ITEM)
//}
//
//func (this *UserCacheManager) RemovePublicItem(uid int32) {
//	this.RemoveUserAttr(uid, UPROP_PUBLIC_ITEM)
//}
//
//func (this *UserCacheManager) SetPublicItem(uid int32, itemID int32) bool {
//	return this.SetUserAttr(uid, UPROP_PUBLIC_ITEM, itemID)
//}
//
//func (this *UserCacheManager) GetPublicItemTimestamp(uid int32) int64 {
//	return this.GetUserAttrInt64(uid, UPROP_PUBLIC_TIMESTAMP)
//}
//
//func (this *UserCacheManager) SetPublicItemTimestamp(uid int32, timestamp int64) bool {
//	return this.SetUserAttr(uid, UPROP_PUBLIC_TIMESTAMP, timestamp)
//}

func (this *UserCacheManager) SetPersistData(id int64, dataKey string, data proto.Marshaler) {
	persistData, _ := data.Marshal()
	key := dataKey + "_" + character.Int64ToString(id)
	this.redisClient.SetData(key, persistData)
}

func (this *UserCacheManager) GetPersistData(id int64, dataKey string, data proto.Unmarshaler) bool {
	key := dataKey + "_" + character.Int64ToString(id)
	persistData := this.redisClient.GetData(key)
	if persistData == "" {
		return false
	}
	err := data.Unmarshal([]byte(persistData))
	if err != nil {
		log.Error("invalid cache format : %v", key)
		return false
	}
	return true
}

//获取用户属性
func (this *UserCacheManager) GetUserAttr(uid int32, propKey string) string {
	return this.redisClient.HGet(GetUserKey(uid), propKey)
}

func (this *UserCacheManager) RemoveUserAttr(uid int32, propKey string) bool {
	return this.redisClient.HDel(GetUserKey(uid), propKey) != nil
}

func (this *UserCacheManager) GetUserAttrInt32(uid int32, propKey string) int32 {
	return this.redisClient.HGetInt32(GetUserKey(uid), propKey)
}

func (this *UserCacheManager) GetUserAttrInt64(uid int32, propKey string) int64 {
	return this.redisClient.HGetInt64(GetUserKey(uid), propKey)
}

func (this *UserCacheManager) GetUserAttrBool(uid int32, propKey string) bool {
	return this.redisClient.HGetBool(GetUserKey(uid), propKey)
}

//func GetUserDataKey(uid int32) string {
//	return  "user_data_" + character.Int32ToString(uid)
//}

//更新用户数据
//func (this *UserCacheManager) UpdateUserData(uid int32, data interface{}) {
//	persistData, _ := bson.Marshal(data)
//	this.redisClient.SetData(GetUserDataKey(uid), persistData)
//}
//
////更新用户数据的过期时间
//func (this *UserCacheManager) UpdateUserExpire(uid int32, data interface{}, seconds int) {
//	this.redisClient.Expire(GetUserDataKey(uid), seconds)
//}

//加载用户数据
//func (this *UserCacheManager) GetUserData(uid int32, data interface{}) bool {
//	persistData := this.redisClient.GetData(GetUserDataKey(uid))
//	if persistData == "" {
//		return false
//	}
//	err := bson.Unmarshal([]byte(persistData), data)
//	if err != nil {
//		log.Error("invalid cache format ")
//		return false
//	}
//	return true
//}

//------------------------用户星球信息-------------------------

////获取用户当前星球类型
//func (this *UserCacheManager) GetUserStarType(uid int32) int32 {
//	return this.redisClient.HGetInt32(GetUserKey(uid), UPROP_STARTYPE)
//}

//---------------------用户未处理的持久化消息----------------------

//func GetUserMessageKey(uid int32) string {
//	return USER_MESSAGE_KEY_PREFIX + character.Int32ToString(uid)
//}
//
////新增用户持久化消息
//func (this *UserCacheManager) AddUserMessage(uid int32, message proto.Marshaler) bool {
//	data, err := message.Marshal()
//	if (err != nil) {
//		return false
//	}
//	_, err1 := this.redisClient.LPush(GetUserMessageKey(uid), data)
//	return err1 == nil
//}
//
//func (this *UserCacheManager) CleanUserMessage(uid int32) {
//	this.redisClient.DelData(GetUserMessageKey(uid))
//}
//
//
////func (this *UserCacheManager) GetAllUserMessageData(uid int64) map[string]string {
////	this.redisClient.LRangeAllByte(GetUserMessageKey(uid))
////}
//
//func (this *UserCacheManager) UpdateUserMessageExpire(uid int32, seconds int) {
//	this.redisClient.Expire(GetUserMessageKey(uid), seconds)
//}
//
////获取所有用户持久化消息
//func (this *UserCacheManager) GetUserMessages(uid int32, id string, message proto.Unmarshaler) bool {
//	data := this.redisClient.HGetBytes(GetUserMessageKey(uid), id)
//	if (data == nil) {
//		return false
//	}
//	err := message.Unmarshal(data)
//	return err == nil
//}


func getMutualKey(id1 int32, id2 int32) string {
	if id1 > id2 {
		return MUTUAL_FLAG_PREFIX + character.Int32ToString(id1) + "_" + character.Int32ToString(id2)
	} else {
		return MUTUAL_FLAG_PREFIX + character.Int32ToString(id2) + "_" + character.Int32ToString(id1)
	}
}

//更新互动
func (this *UserCacheManager) UpdateMutual(id1 int32, id2 int32, timeout int){
	key := getMutualKey(id1, id2)
	this.redisClient.SetData(key, "")
	this.redisClient.Expire(key, timeout)
}

//是否存在互动关系
func (this *UserCacheManager) ExistMutual(id1 int32, id2 int32) bool {
	result, _  := this.redisClient.Exists(getMutualKey(id1, id2))
	return result
}

// ------------------------------------------------------------------------

func getItemHelpWeChatKey(beHelpID int32, itemID int32) string {
	return USER_HELPITEM_WECHAT_KEY + character.Int32ToString(beHelpID) + "_" + character.Int32ToString(itemID)
}

func (this *UserCacheManager) SetHelpItemWechatUid(uid int32, beHelpID int32, itemID int32) bool {
	return this.redisClient.SAddData(getItemHelpWeChatKey(beHelpID, itemID), uid)
}



func (this *UserCacheManager) ExistHelpItemWechatUid(uid int32, beHelpID int32, itemID int32) bool {
	return this.redisClient.SContains(getItemHelpWeChatKey(beHelpID, itemID), uid)
}

func (this *UserCacheManager) ExistBeHelpItemCacheUid(beHelpID int32, itemID int32) bool {
	ret, err := this.redisClient.Exists(getItemHelpWeChatKey(beHelpID, itemID))
	if err != nil {
		log.Error("ExistBeHelpItemCacheUid uid: %v error: %v", beHelpID, err.Error())
		return false
	}
	return ret
}

func (this *UserCacheManager) DelBeHelpItemCacheUid(beHelpID int32, item int32) bool {
	return this.redisClient.DelData(getItemHelpWeChatKey(beHelpID, item))
}

//----------------------------------每日礼包--------------------------------------

func (this *UserCacheManager) SetUserDayGiftCache(uid int32, draw bool) bool{
	return this.redisClient.SetData(GetUserDayGiftKey(uid), draw)
	//return this.redisClient.SetExpireData(GetUserDayGiftKey(uid), false, int(util.GetTodayHourTime(24).Unix()))
}

func (this *UserCacheManager) GetUserDayGiftCache(uid int32) string {
	return this.redisClient.GetData(GetUserDayGiftKey(uid))
}

func (this *UserCacheManager) ExistUserDayGiftCache(uid int32) bool {
	ret, err := this.redisClient.Exists(GetUserDayGiftKey(uid))
	if err != nil {
		log.Error("error:%v", err.Error())
		return false
	}
	return ret
}

func (this *UserCacheManager) SetUserDayGiftExpireTime(uid int32, seconds int) bool {
	return this.redisClient.Expire(GetUserDayGiftKey(uid), seconds)
}

//--------------------------------测试账号------------------------------------------
func (this *UserCacheManager) SetTestUserID(uid int32) bool {
	return this.redisClient.SAddData(USER_TEST_KEY_PREFIX, uid)
}

func (this *UserCacheManager) ExistTestUserID(uid int32) bool {
	return this.redisClient.SContains(USER_TEST_KEY_PREFIX, uid)
}

func (this *UserCacheManager) DelTestUserID(uid int32) bool{
	return this.redisClient.SDelData(USER_TEST_KEY_PREFIX, uid)
}

func (this *UserCacheManager) GetTestUserID() []int {
	return this.redisClient.SMembers(USER_TEST_KEY_PREFIX)
}