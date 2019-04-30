package cache

import (
	basecache "gok/cache"
	"gok/module/center/conf"

)

var UserCache = basecache.NewUserCacheManager()
var CenterCache = basecache.NewCenterCacheManager()
var ClusterCache = basecache.NewClusterCacheManager()

func Init() {
	UserCache.Init1(conf.Server.RedisAddress, conf.Server.RedisPassword,
		conf.Server.RedisMaxActive, conf.Server.RedisMaxIdle, conf.Server.RedisIdleTimeout)

	CenterCache.Init1(conf.Server.RedisAddress, conf.Server.RedisPassword,
		conf.Server.RedisMaxActive, conf.Server.RedisMaxIdle, conf.Server.RedisIdleTimeout)

	ClusterCache.Init1(conf.Server.RedisAddress, conf.Server.RedisPassword,
		conf.Server.RedisMaxActive, conf.Server.RedisMaxIdle, conf.Server.RedisIdleTimeout)
}

func Close() {

	//清除所有缓存数据
	UserCache.Close()
	CenterCache.Close()
	ClusterCache.Close()
}