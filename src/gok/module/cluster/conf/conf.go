/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/8/4
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package conf

import (
	"gok/module/cluster/util"
	"gopkg.in/mgo.v2/bson"
	"gok/config"
)

var Config struct {
	Node        	  string	//当前集群节点的标识，信息需要注册到中心服务器
	RedisAddress      string
	RedisPassword     string
	RedisMaxActive    int
	RedisMaxIdle      int
	RedisIdleTimeout  int

	ZKServers   []string //集群中心服务器地址
	ZKName  	string
	LBS         string   //负载均衡策略  polling 轮询

}

func init() {
	config.LoadConfigData("conf/gok/cluster.json", &Config)
	if Config.ZKServers == nil || len(Config.ZKServers) == 0 {
		Config.ZKServers = append(Config.ZKServers, "127.0.0.1:2181")
	}
	if Config.ZKName == "" {
		Config.ZKName = "gok"
	}
	if Config.Node == "" {
		Config.Node = bson.NewObjectId().Hex()
	}
	addr, ret := util.GetZkEnvValue()
	if ret {
		Config.ZKServers[0] = addr
	}
}
