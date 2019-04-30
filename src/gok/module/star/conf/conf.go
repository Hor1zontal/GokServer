package conf

import (
	"aliens/database/dbconfig"
	"gok/module/cluster/util"
	"gok/config"
)

var (
	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000
)


var Server struct {
	Enable           bool
	Database dbconfig.DBConfig
	FreeTimeout float64
	RedisAddress     string
	RedisPassword    string
	RedisMaxActive   int
	RedisMaxIdle     int
	RedisIdleTimeout int
	RPCAddress       string //提供RPC服务的地址,信息需要注册到中心服务器供其他服务调用
	RPCPort          int    //提供RPC服务的端口，本地启动RPC需要指定此端口启动
}

func init() {
	config.LoadConfigData("conf/gok/star/server.json", &Server)
	if Server.RPCAddress != "" {
		return
	}
	Server.RPCAddress = util.GetAddress(Server.RPCPort)
	// 默认本地星球数据没有更改释放时间一天
	if Server.FreeTimeout <= 0 {
		Server.FreeTimeout = 86400
	}
	util.GetMongodbEnvValue(&Server.Database)
}
