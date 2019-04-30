/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/7/27
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"sync"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"errors"
)

var ServiceManager = initMessageServiceManager()

var INVALID_SERVICE_RESPONSE = &protocol.GS2C{
	Sequence: ExceptionSeq,
	ResultPush: &protocol.ResultPush{
		Result: int32(exception.SERVICE_INVALID),
	},
}

func initMessageServiceManager() *MessageServiceManager {
	manager := &MessageServiceManager{
		remoteServices: make(map[string]*RemoteService),
		localServices: make(map[string]*LocalService),
	}
	return manager
}

//服务组件名字，后续采用此KEY在zookeeper上查询服务来rpc调用
const (
	SERVICE_VISITOR       string = "visitor" //处理未验权的网络连接消息
	SERVICE_USER          string = "user"
	SERVICE_USER_RPC      string = "userrpc"
	SERVICE_PASSPORT_RPC  string = "passport"
	SERVICE_LOG_RPC  	  string = "log"
	SERVICE_TRADE     	  string = "trade"
	SERVICE_COMMUNITY_RPC string = "community"
	SERVICE_STAR_RPC      string = "star"
	SERVICE_SEARCH_RPC    string = "search"
	SERVICE_CENTER_RPC    string = "center"
	SERVICE_MAIL_RPC	  string = "mail"
)

//消息服务,抽象层，可以为local也可以是remote
//type IMessageService interface {
//	GetType() string                                                      //获取消息服务类型
//	CanHandle(seq int32) bool                                             //能否处理指定编号的消息
//	HandleMessage(message *protocol.C2GS) *protocol.GS2C                  //阻塞调用服务接口
//	HandleChannelMessage(message *protocol.C2GS, channel IMessageChannel) //异步调用服务接口 channel为消息回写管道
//}

//type IRemoteService interface {
//	AsyncBroadcastAll(message *protocol.C2GS) bool
//	AsyncHandleRemoteMessage(serviceID string, message *protocol.C2GS) *protocol.GS2C         //阻塞调用指定服务接口，
//	AsyncHandlePriorityRemoteMessage(serviceID string, message *protocol.C2GS) *protocol.GS2C //阻塞调用指定服务接口,优先发送到serviceID节点，没有会分配一个节点处理
//}

//服务容器,管理本地加载的服务句柄
type MessageServiceManager struct {
	sync.RWMutex
	remoteServices map[string]*RemoteService //远程服务句柄  处理业务消息
	localServices  map[string]*LocalService  //本地服务句柄
}

//func (this *MessageServiceManager) allocService(serviceType string) *RemoteService {
//
//}

//注册本地服务
func (this *MessageServiceManager) RegisterLocalService(service *LocalService) {
	this.localServices[service.GetType()] = service
}

//订阅远程消息服务
func (this *MessageServiceManager) SubscribeRemoteService(serviceType string) *RemoteService {
	//log.Debug("register service %v", service.GetType())
	service := this.remoteServices[serviceType]
	if service == nil {
		service = NewRemoteService(serviceType)
		this.remoteServices[serviceType] = service
	}
	return service
}

//异步的根据消息编号找到指定的服务去处理,优先使用本地服务
func (this *MessageServiceManager) HandleChannelMessage(message *protocol.C2GS, channel IMessageChannel) bool {
	seq := message.GetSequence()[0]
	//优先用户服务处理
	for _, service := range this.localServices {
		if service.CanHandle(seq) {
			service.HandleChannelMessage(message, channel)
			return true
		}
	}
	for _, service := range this.remoteServices {
		if service.CanHandle(seq) {
			service.HandleChannelMessage(message, channel)
			return true
		}
	}
	return false
}

//处理模块消息,确保响应消息和请求消息编号一致，如果是错误码会继续往外面抛
func (this *MessageServiceManager) AssertHandleRemoteMessage(serviceType string, message *protocol.C2GS) *protocol.GS2C {
	service := this.remoteServices[serviceType]
	if service == nil {
		return INVALID_SERVICE_RESPONSE
	}
	response := service.HandleMessage(message)
	if response.GetResultPush() != nil {
		exception.GameException(exception.GameCode(response.GetResultPush().GetResult()))
	}
	return response
}

//如果是错误码不会继续往外抛
func (this *MessageServiceManager) HandleRemoteMessage(serviceType string, message *protocol.C2GS) *protocol.GS2C {
	service := this.remoteServices[serviceType]
	if service == nil {
		return INVALID_SERVICE_RESPONSE
	}
	response := service.HandleMessage(message)
	if response.GetResultPush() != nil {
		return nil
	}
	return response
}

func (this *MessageServiceManager) AsyncHandleMessage(serviceType string, message *protocol.C2GS) error {
	service := this.remoteServices[serviceType]
	if service == nil {
		return errors.New("remote service can not alloc " + serviceType)
	}
	return service.AsyncHandleMessage(message)
}


//处理远程模块消息
func (this *MessageServiceManager) AsyncHandleRemoteMessage(serviceType string, serviceID string, message *protocol.C2GS) error {
	service := this.remoteServices[serviceType]
	if service == nil {
		return errors.New("remote service can not alloc " + serviceType)
	}
	return service.AsyncHandleRemoteMessage(serviceID, message)
}

func (this *MessageServiceManager) AssertSyncHandleRemoteMessage(serviceType string, serviceID string, message *protocol.C2GS) *protocol.GS2C {
	service := this.remoteServices[serviceType]
	if service == nil {
		return INVALID_SERVICE_RESPONSE
	}
	response := service.HandleRemoteMessage(serviceID, message)
	if response.GetResultPush() != nil {
		exception.GameException(exception.GameCode(response.GetResultPush().GetResult()))
	}
	return response
}

func (this *MessageServiceManager) SyncHandleRemoteMessage(serviceType string, serviceID string, message *protocol.C2GS) *protocol.GS2C {
	service := this.remoteServices[serviceType]
	if service == nil {
		return INVALID_SERVICE_RESPONSE
	}
	response := service.HandleRemoteMessage(serviceID, message)
	if response.GetResultPush() != nil {
		return nil
	}
	return response
}

func (this *MessageServiceManager) AsyncBroadcastAllRemote(serviceType string, message *protocol.C2GS) bool {
	service := this.remoteServices[serviceType]
	if service == nil {
		return false
	}
	return service.AsyncBroadcastAll(message)
}

func (this *MessageServiceManager) AsyncBroadcastAllRemoteIgnore(serviceType string, serviceID string, message *protocol.C2GS) bool {
	service := this.remoteServices[serviceType]
	if service == nil {
		return false
	}
	return service.AsyncBroadcastAllIgnore(serviceID, message)
}

//优先发送到指定的serviceID,如果没有发送到其他节点
func (this *MessageServiceManager) AsyncHandlePriorityRemoteMessage(serviceType string, serviceID string, message *protocol.C2GS) error {
	service := this.remoteServices[serviceType]
	if service == nil {
		return errors.New("remote service can not alloc " + serviceType)
	}
	return service.AsyncHandlePriorityRemoteMessage(serviceID, message)
}