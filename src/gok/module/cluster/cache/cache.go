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
	"gok/module/cluster/conf"
)

var Cluster = basecache.NewClusterCacheManager()

func Init() {
	Cluster.Init1(conf.Config.RedisAddress, conf.Config.RedisPassword,
		conf.Config.RedisMaxActive, conf.Config.RedisMaxIdle, conf.Config.RedisIdleTimeout)
}

func Close() {
	Cluster.Close()
}
