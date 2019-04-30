package conf

import (
	"time"
	"aliens/database/dbconfig"
	"gok/module/cluster/util"
	"gok/config"
)

var (
	// gate conf
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 4096
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2
	LittleEndian           = true

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000
)

var Server struct {
	Enable	     	   bool
	Database dbconfig.DBConfig
	LocalCacheTimeout  float64
	UserFreeTimeout    float64
	SyncDBInterval	   float64  //用户在线数据同步到数据库的时间
	RedisAddress       string
	RedisPassword      string
	RedisMaxActive   int
	RedisMaxIdle     int
	RedisIdleTimeout int
	RPCAddress         string	//提供RPC服务的地址,信息需要注册到中心服务器供其他服务调用
	RPCPort            int	//提供RPC服务的端口，本地启动RPC需要指定此端口启动
}

func init() {
	initLoader()
	config.LoadConfigData("conf/gok/game/server.json", &Server)
	if Server.RPCAddress != "" {
		return
	}
	// 默认本地缓存过期时间 30分钟
	if Server.LocalCacheTimeout <= 0 {
		Server.LocalCacheTimeout = 1800
	}
	// 异步更新数据库 默认10分钟
	if Server.SyncDBInterval <= 0 {
		Server.SyncDBInterval = 600
	}
	// 没有操作踢人默认10分钟
	if Server.UserFreeTimeout <= 0 {
		Server.UserFreeTimeout = 600
	}

	Server.RPCAddress = util.GetAddress(Server.RPCPort)
	util.GetMongodbEnvValue(&Server.Database)
}
