/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/5/12
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package cache

import (
	basecache "gok/cache"
	"gok/module/community/conf"
	"gok/module/community/db"
	"aliens/common/character"
	"aliens/log"
)

var CommunityCache = basecache.NewCommunityCacheManager()

func Init() {
	CommunityCache.Init1(conf.Server.RedisAddress, conf.Server.RedisPassword,
		conf.Server.RedisMaxActive, conf.Server.RedisMaxIdle, conf.Server.RedisIdleTimeout)

	if CommunityCache.SetNX(basecache.FLAG_LOADFOLLOW, 1) {
		log.Debug("start load follow data to redis cache...")
		count := 0
		var follows []*db.DFollow
		err := db.DatabaseHandler.QueryAllLimit(&db.DFollow{}, &follows, 10000, func(data interface{}) bool {
			for _, follow := range follows {
				CommunityCache.AddFollower1(character.Int32ToString(follow.ID.SubID1), character.Int32ToString(follow.ID.SubID2), follow.AddTime)
			}
			currLen := len(follows)
			count += currLen
			return currLen == 0
		})
		if err != nil {
			log.Debug("load follow err: %v", err)
		}
		log.Debug("end load follow data to redis cache:%v", count)
	}


	//if CommunityCache.SetNX(basecache.FLAG_LOADFOLLOW, 1) {
	//	log.Debug("start load follow data to redis cache...")
	//	var follows []*db.DFollow
	//	db.DatabaseHandler.QueryAll(&db.DFollow{}, &follows)
	//	for _, follow := range follows {
	//		CommunityCache.AddFollower1(character.Int32ToString(follow.ID.SubID1), character.Int32ToString(follow.ID.SubID2), follow.AddTime)
	//	}
	//	log.Debug("end load follow data to redis cache")
	//}

	if CommunityCache.SetNX(basecache.FLAG_LOADMOMENT_TIMELINE, 1) {
		log.Debug("start load moment timeline data to redis cache...")
		count := 0
		var moments []*db.DMoments
		err := db.DatabaseHandler.QueryAllLimit(&db.DMoments{}, &moments, 10000, func(data interface{}) bool {
			for _, moment := range moments {
				//更新自己发布的朋友圈消息
				CommunityCache.SetUserPublicMomentID(moment.Uid, moment.ID, moment.CreateTime)
				CommunityCache.SetUserReceiveMomentID(moment.Uid, moment.ID, moment.CreateTime)

				//更新关注发布者的所有用户的朋友圈时间线
				publicUID := character.Int32ToString(moment.Uid)
				followings := CommunityCache.GetFollowings(publicUID)
				for following := range followings {
					CommunityCache.SetUserReceiveMomentID(character.StringToInt32(following), moment.ID, moment.CreateTime)
				}
			}
			currLen := len(moments)
			count += currLen
			return currLen == 0
		})
		if err != nil {
			log.Debug("load moment timeline err: %v", err)
		}
		log.Debug("end load moment timeline data to redis cache count:%v", count)
	}

	//if CommunityCache.SetNX(basecache.FLAG_LOADMOMENT_TIMELINE, 1) {
	//	log.Debug("start load moment data to redis cache...")
	//	var moments []*db.DMoments
	//	db.DatabaseHandler.QueryAll(&db.DMoments{}, &moments)
	//	for _, moment := range moments {
	//		//更新自己发布的朋友圈消息
	//		CommunityCache.SetUserPublicMomentID(moment.Uid, moment.ID, moment.CreateTime)
	//		CommunityCache.SetUserReceiveMomentID(moment.Uid, moment.ID, moment.CreateTime)
	//
	//		//更新关注发布者的所有用户的朋友圈时间线
	//		publicUID := character.Int32ToString(moment.Uid)
	//		followings := CommunityCache.GetFollowings(publicUID)
	//		for following, _ := range followings {
	//			CommunityCache.SetUserReceiveMomentID(character.StringToInt32(following), moment.ID, moment.CreateTime)
	//		}
	//	}
	//	log.Debug("end load moment data to redis cache")
	//}
}

func Close() {
	CommunityCache.Close()
}
