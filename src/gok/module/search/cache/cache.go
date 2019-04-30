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
	"gok/module/search/conf"
)

var StarCache = basecache.NewStarCacheManager()

func Init() {
	StarCache.Init1(conf.Server.RedisAddress, conf.Server.RedisPassword,
		conf.Server.RedisMaxActive, conf.Server.RedisMaxIdle, conf.Server.RedisIdleTimeout)
}

func Close() {
	StarCache.Close()
}
