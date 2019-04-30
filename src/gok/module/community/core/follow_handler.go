package core

import (
	"gok/service/msg/protocol"
	"aliens/common/character"
	"gok/module/community/cache"
)

//构建好友申请列表
func BuildFollowerList(id string) *protocol.GetFollowerListRet {
	followers := cache.CommunityCache.GetFollowers(id)
	results := []*protocol.UserInfo{}
	for followerID, followerTime := range followers {
		results = append(results, BuildFollowerInfo(id, followerID, followerTime))
	}
	return &protocol.GetFollowerListRet{
		Followers: results,
	}
}

func BuildFollowingList(id string) *protocol.GetFollowingListRet {
	followings := cache.CommunityCache.GetFollowings(id)
	results := []*protocol.UserInfo{}
	for followingID, followerTime := range followings {
		results = append(results, BuildFollowingInfo(followingID, followerTime))
	}
	return &protocol.GetFollowingListRet{
		Followings: results,
	}
}

func BuildFollowerInfo(id string, followerID string, followerTime int64) *protocol.UserInfo {
	return &protocol.UserInfo{
		Id: character.StringToInt32(followerID),
		FollowEachOther: cache.CommunityCache.ExistFollower(followerID, id),
		FollowTime:followerTime,
	}
}

func BuildFollowingInfo(followingID string, followingTime int64) *protocol.UserInfo {
	return &protocol.UserInfo{
		Id: character.StringToInt32(followingID),
		//FollowEachOther: cache.CommunityCache.ExistFollower(refID, id)),
		FollowTime:followingTime,
	}
}
