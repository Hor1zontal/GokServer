/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/11/14
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package util

import (
	"aliens/common/cache"
	"aliens/database/dbconfig"
	"aliens/log"
	"gok/constant"
	"net"
	"os"
	"strconv"
)

var ip = GetIP()

func GetAddress(port int) string {
	return ip + ":" + strconv.Itoa(port)
}

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	return ""
}

func GetRedisEnvValue(redisClient *cache.RedisCacheClient) {
	redisAddr := os.Getenv(constant.REDIS_ADDRESS)
	if redisAddr != "" {
		//log.Info("redisAddr:%v",redisAddr)
		redisClient.Address = redisAddr
	}
}

func GetMongodbEnvValue(db *dbconfig.DBConfig) {
	mongodbAddr := os.Getenv(constant.MONGODB_ADDRESS)
	mongodbUser := os.Getenv(constant.MONGODB_ROOT_USER)
	mongodbPsw := os.Getenv(constant.MONGODB_ROOT_PSW)
	if mongodbAddr != "" {
		db.Address = "mongodb://" + mongodbUser + ":" + mongodbPsw + "@" + mongodbAddr + "/admin"
		log.Info("mongodbAddr:%v", db.Address)
	}
}

func GetZkEnvValue() (string, bool) {
	zkAddr := os.Getenv(constant.ZOOKEEPER_ADDRESS)
	if zkAddr != "" {
		log.Info("zkAddr:%v",zkAddr)
		return zkAddr, true
	}
	return "", false
}