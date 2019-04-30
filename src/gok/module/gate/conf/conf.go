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
	"time"
	"gok/config"
)

var (
	// gate conf
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 4096
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2
	LittleEndian           = true
)

var Config struct {
	Enable              bool   //网络模块是否开启
	LogLevel            string //日志级别
	LogPath             string //日志输出路径
	WSAddr              string //
	TCPAddr             string //
	PublicWSAddress     string //当前节点注册到中心服务器的websocket连接地址信息
	PublicTCPAddres     string //当前节点注册到中心服务器的tcp连接地址信息
	MaxConnNum          int    //最大连接数
	//ConsolePort         int
	//ProfilePath         string
	SecretKey	    	string
	//AuthTimeout			float64
	MessageChannelLimit int //单个连接消息请求管道最大消息缓存数量
}


func init() {
	config.LoadConfigData("conf/gok/gate/server.json", &Config)
	// 默认最大连接数20000
	if Config.MaxConnNum <= 0 {
		Config.MaxConnNum = 20000
	}
	if Config.MessageChannelLimit <= 0 {
		Config.MessageChannelLimit = 5
	}
	//if Config.AuthTimeout == 0 {
	//	Config.AuthTimeout = 30
	//}
}
