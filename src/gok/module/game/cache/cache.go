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
	"gok/module/game/conf"
	"gok/module/game/db"
	"aliens/log"
)

var UserCache = basecache.NewUserCacheManager()
var StarCache = basecache.NewStarCacheManager()
var TradeCache = basecache.NewTradeCacheManager()

func Init() {
	StarCache.Init1(conf.Server.RedisAddress, conf.Server.RedisPassword,
		conf.Server.RedisMaxActive, conf.Server.RedisMaxIdle, conf.Server.RedisIdleTimeout)

	UserCache.Init1(conf.Server.RedisAddress, conf.Server.RedisPassword,
		conf.Server.RedisMaxActive, conf.Server.RedisMaxIdle, conf.Server.RedisIdleTimeout)

	TradeCache.Init1(conf.Server.RedisAddress, conf.Server.RedisPassword,
		conf.Server.RedisMaxActive, conf.Server.RedisMaxIdle, conf.Server.RedisIdleTimeout)


	//if UserCache.SetNX(basecache.FLAG_LOADROLE, 1) {
	//	log.Debug("start load role data to redis cache...")
	//	var roles []*db.DBRole
	//	db.DatabaseHandler.QueryAll(&db.DBRole{}, &roles)
	//	for _,role := range roles {
	//		UpdateRoleCache(role)
	//	}
	//	log.Debug("end load role data to redis cache count:%v", len(roles))
	//	//var sales []*db.DBSale
	//	//db.DatabaseHandler.QueryAll(&db.DBSale{}, &sales)
	//	//for _, sale := range sales {
	//	//	UserCache.SetUserSale(sale.UserID, sale.ItemID)
	//	//}
	//}

	if UserCache.SetNX(basecache.FLAG_LOADROLE, 1) {
		log.Debug("start load role data to redis cache...")
		count := 0
		var roles []*db.DBRole
		err := db.DatabaseHandler.QueryAllLimit(&db.DBRole{}, &roles, 10000, func(data interface{}) bool {
			for _, role := range roles {
				UpdateRoleCache(role)
			}
			currLen := len(roles)
			count += currLen
			return currLen == 0
		})
		if err != nil {
			log.Debug("load role err: %v", err)
		}

		log.Debug("end load role to redis cache count:%v", count)
	}

}

func UpdateRoleCache(role *db.DBRole) {

	//UserCache.SetUserOnlineTimestamp(role.UserID, role.LoginTime)
	//UserCache.AddNicknameUIDMapping(role.NickName, role.UserID)
	UserCache.HSetUser(role.UserID, role)
}

func Close() {
	StarCache.Close()
	//UserCache.FlashAll()
	UserCache.Close()
}
