package conf

import (
	vivomodel "gok/module/passport/vivo/model"
	"time"
	"gok/module/cluster/util"
	"aliens/database/dbconfig"
	"gok/config"
	wxmodel "gok/module/passport/wx/model"
)

var Server struct {
	Enable            bool
	Database          dbconfig.DBConfig
	RedisAddress      string
	RedisPassword     string
	RedisMaxActive    int
	RedisMaxIdle      int
	RedisIdleTimeout  int
	DefaultChannelPWD string
	RPCAddress        string //提供RPC服务的地址,信息需要注册到中心服务器供其他服务调用
	RPCPort           int    //提供RPC服务的端口，本地启动RPC需要指定此端口启动
	TokenExpireTime   int
	HTTPAddress       string
	AppKey            string
	IsSign            bool
	WeChat            wxmodel.Config
	Vivo              vivomodel.Config
	VersionUrl        string
	ServerInfo        string
}

func init() {
	config.LoadConfigData("conf/gok/passport/server.json", &Server)
	if Server.RPCAddress == "" {
		Server.RPCAddress = util.GetAddress(Server.RPCPort)
	}
	if Server.TokenExpireTime <= 0 {
		//默认过期时间一个月
		Server.TokenExpireTime = int(30 * 24 * time.Hour)
	}
	if Server.VersionUrl == "" {
		Server.VersionUrl = "http://www.alienidea.com/cdn/Gok/managerConfig.json"
	}
	//if Server.AppID == "" {
	//	Server.AppID = "wxb0c02c8571ad34fc"
	//	Server.AppSecret = "8fd2f62577d623021890576a11c09828"
	//}
	util.GetMongodbEnvValue(&Server.Database)
}