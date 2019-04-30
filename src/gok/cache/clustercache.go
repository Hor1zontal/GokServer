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

import "aliens/common/character"

const (
	keyAccessToken string = "wx:accessToken" //

	keyOnline string = "online:" //用户内存所在的节点信息
	keySession string = "session:" //用户内存所在的节点信息
	keyIpWhiteList = "whitelist:ip"
	keyUidWhiteList = "whitelist:uid"

	FLAG_LOADSTAR_SUMMARY string = "flag:star_summary" //标识，是否加载星球统计数据到redis

	keyStarSummary = "summary:star"
)

type ClusterCacheManager struct {
	*cacheManager
}

func NewClusterCacheManager() *ClusterCacheManager {
	return &ClusterCacheManager{
		&cacheManager{},
	}
}

func getUserNodeKey(uid int32) string {
	return keySession + character.Int32ToString(uid)
}

func getUserOnlineNodeKey(uid int32) string {
	return keyOnline + character.Int32ToString(uid)
}

//新增uid白名单
func (this *ClusterCacheManager) AddStarSummary(starType int32) int {
	return this.redisClient.HIncrby(keyStarSummary, character.Int32ToString(starType), 1)
}

//新增uid白名单
func (this *ClusterCacheManager) DelStarSummary(starType int32) int {
	return this.redisClient.HIncrby(keyStarSummary, character.Int32ToString(starType), -1)
}

func (this *ClusterCacheManager) GetStarSummary() map[string]int {
	return this.redisClient.HGetAllInt(keyStarSummary)
}

func (this *ClusterCacheManager) CleanSummary() {
	this.redisClient.DelData(keyStarSummary)
}

//更新微信accessToken
func (this *ClusterCacheManager) SetAccessToken(accessToken string, expire int) bool {
	result := this.redisClient.SetData(keyAccessToken, accessToken)
	this.redisClient.Expire(keyAccessToken, expire)
	return result
}

//获取微信的accessToken
func (this *ClusterCacheManager) GetAccessToken() string {
	return this.redisClient.GetData(keyAccessToken)
}

//设置用户会话所在的服务节点
func (this *ClusterCacheManager) SetUserNode(uid int32, node string) bool {
	return this.redisClient.SetData(getUserNodeKey(uid), node)
}

//设置用户会话所在的服务节点
func (this *ClusterCacheManager) CleanUserNode(uid int32) bool {
	return this.redisClient.DelData(getUserNodeKey(uid))
}


//获取用户会话所在的服务节点
func (this *ClusterCacheManager) GetUserNode(uid int32) string {
	return this.redisClient.GetData(getUserNodeKey(uid))
}

//新增ip白名单
func (this *ClusterCacheManager) AddIpWhiteList(ip string) bool {
	return this.redisClient.SAddData(keyIpWhiteList, ip)
}

//删除ip白名单
func (this *ClusterCacheManager) RemoveIpWhiteList(ip string) bool {
	return this.redisClient.SDelData(keyIpWhiteList, ip)
}

//是否在ip白名单中
func (this *ClusterCacheManager) IsIpWhiteList(ip string) bool {
	return this.redisClient.SContains(keyIpWhiteList, ip)
}

//新增uid白名单
func (this *ClusterCacheManager) AddUidWhiteList(uid int32) bool {
	return this.redisClient.SAddData(keyUidWhiteList, uid)
}

//删除uid白名单
func (this *ClusterCacheManager) RemoveUidWhiteList(uid int32) bool {
	return this.redisClient.SDelData(keyUidWhiteList, uid)
}

//是否在uid白名单中
func (this *ClusterCacheManager) IsUidWhiteList(uid int32) bool {
	return this.redisClient.SContains(keyUidWhiteList, uid)
}

func (this *ClusterCacheManager) GetWhiteListUid() []int{
	return this.redisClient.SMembers(keyUidWhiteList)
}

//设置用户会话所在的服务节点
func (this *ClusterCacheManager) SetUserOnlineNode(uid int32, node string) bool {
	return this.redisClient.SetData(getUserOnlineNodeKey(uid), node)
}

//获取用户会话所在的服务节点
func (this *ClusterCacheManager) GetUserOnlineNode(uid int32) string {
	return this.redisClient.GetData(getUserOnlineNodeKey(uid))
}

func (this *ClusterCacheManager) IsUserOnline(uid int32) bool {
	return this.redisClient.GetData(getUserOnlineNodeKey(uid)) != ""
}
