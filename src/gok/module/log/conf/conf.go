/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package conf

import (
	"aliens/database/dbconfig"
	"gok/config"
	"gok/module/cluster/util"
)

var Config struct {
	Enable			bool
	//RPCPort			int32
	Database 		dbconfig.DBConfig
	RedisAddress       string
	RedisPassword      string
	RedisMaxActive   int
	RedisMaxIdle     int
	RedisIdleTimeout int
}

func init() {
	config.LoadConfigData("conf/gok/log/server.json", &Config)
	util.GetMongodbEnvValue(&Config.Database)
}
