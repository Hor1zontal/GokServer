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
	"aliens/common/cache"
	"gok/module/cluster/util"
	"time"
)


type cacheManager struct {
	redisClient *cache.RedisCacheClient
}

func (this *cacheManager) Init(redisAddress string) {
	this.Init1(redisAddress, "", 0, 0, 0)
}

func (this *cacheManager) Init1(redisAddress string, password string, maxActive int, maxIdle int, idleTimeout int) {
	if maxActive == 0 {
		maxActive = 2000
	}
	if maxIdle == 0 {
		maxIdle = 1000
	}
	if idleTimeout == 0 {
		idleTimeout = 120
	}
	redisClient := &cache.RedisCacheClient{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		Address:     redisAddress,
		Password:    password,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
	}
	util.GetRedisEnvValue(redisClient)
	redisClient.Start()
	this.redisClient = redisClient
}

func (this *cacheManager) Close() {
	if this.redisClient != nil {
		this.redisClient.Close()
	}
}

func (this *cacheManager) SetNX(key string, value interface{}) bool {
	return this.redisClient.SetNX(key, value)
}

