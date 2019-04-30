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
	"gok/module/log/conf"
)

var LogCache = basecache.NewUserCacheManager()

func Init() {
	LogCache.Init1(conf.Config.RedisAddress, conf.Config.RedisPassword,
		conf.Config.RedisMaxActive, conf.Config.RedisMaxIdle, conf.Config.RedisIdleTimeout)


	//if UserCache.SetNX(basecache.FLAG_LOADROLE, 1) {
	//	var roles []*db.DBRole
	//	db.DatabaseHandler.QueryAll(&db.DBRole{}, &roles)
	//	for _,role := range roles {
	//		UpdateRoleCache(role)
	//	}
	//
	//	//var sales []*db.DBSale
	//	//db.DatabaseHandler.QueryAll(&db.DBSale{}, &sales)
	//	//for _, sale := range sales {
	//	//	UserCache.SetUserSale(sale.UserID, sale.ItemID)
	//	//}
	//}

}

func Close() {
	//UserCache.FlashAll()
	LogCache.Close()
}
