package conf

import (
	"aliens/database/dbconfig"
	"gok/module/cluster/util"
	"gok/config"
)

var Server struct{
	Enable            bool
	Database dbconfig.DBConfig
	RedisAddress	  string
	RedisPassword      string
	RedisMaxActive   int
	RedisMaxIdle     int
	RedisIdleTimeout int
	RedisExpireTime  int

	RPCAddress        string //提供RPC服务的地址,信息需要注册到中心服务器供其他服务调用
	RPCPort           int    //提供RPC服务的端口，本地启动RPC需要指定此端口启动
}

func init() {
	config.LoadConfigData("conf/gok/mail/server.json", &Server)
	if Server.RPCAddress == "" {
		Server.RPCAddress = util.GetAddress(Server.RPCPort)
	}
	// 邮件过期时间默认一个月
	if Server.RedisExpireTime <= 0 {
		Server.RedisExpireTime = 2592000
	}
	util.GetMongodbEnvValue(&Server.Database)
}
