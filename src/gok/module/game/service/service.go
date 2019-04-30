/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/8/4
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"gok/service"
	"gok/module/game/conf"
	"github.com/name5566/leaf/chanrpc"
)

var UserRPCService *service.GRPCService = nil

func Init(chanRpc *chanrpc.Server) {
	service.ServiceManager.RegisterLocalService(VisitorService)
	service.ServiceManager.RegisterLocalService(UserService)
	//service.ServiceManager.RegisterLocalService(EventService)

	//发布用户服务到服务中心
	service.PublicWBService(service.SERVICE_USER)
	//配置了RPC，需要发布服务到ZK
	UserRPCService = service.PublicRPCService1(LocalUserRPCService, conf.Server.RPCAddress, conf.Server.RPCPort, chanRpc)
	//RPC消息异步调用，不用回写用户管道反馈给客户端
	//LocalUserRPCService.SetWriteBack(false)
	//service.ServiceManager.RegisterLocalService(LocalUserRPCService)
	//用户RPC远程服务
	service.ServiceManager.SubscribeRemoteService(service.SERVICE_USER_RPC)
	//service.ServiceManager.SubscribeRemoteService(service.SERVICE_EVENT_RPC)
	service.ServiceManager.SubscribeRemoteService(service.SERVICE_STAR_RPC)
	service.ServiceManager.SubscribeRemoteService(service.SERVICE_COMMUNITY_RPC)
	service.ServiceManager.SubscribeRemoteService(service.SERVICE_PASSPORT_RPC)
	service.ServiceManager.SubscribeRemoteService(service.SERVICE_CENTER_RPC)
	service.ServiceManager.SubscribeRemoteService(service.SERVICE_MAIL_RPC)
	service.ServiceManager.SubscribeRemoteService(service.SERVICE_SEARCH_RPC)
	service.ServiceManager.SubscribeRemoteService(service.SERVICE_TRADE)
}

func Close() {
	VisitorService.Close()
	LocalUserRPCService.Close()
	UserService.Close()
	UserRPCService.Close()
}



