package conf

import (
	"gok/module/cluster/util"
	"aliens/database/dbconfig"
	"gok/config"
)

var Server struct {
	Enable      bool
	Database dbconfig.DBConfig
	RPCAddress         string	//提供RPC服务的地址,信息需要注册到中心服务器供其他服务调用
	RPCPort            int	//提供RPC服务的端口，本地启动RPC需要指定此端口启动

	RedisAddress	 string
	RedisPassword    string
	RedisMaxActive   int
	RedisMaxIdle     int
	RedisIdleTimeout int
}

func init() {
	config.LoadConfigData("conf/gok/trade/server.json", &Server)
	if Server.RPCAddress != "" {
		return
	}
	Server.RPCAddress = util.GetAddress(Server.RPCPort)
	util.GetMongodbEnvValue(&Server.Database)
}


