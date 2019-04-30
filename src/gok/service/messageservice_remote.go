/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/24
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"gok/module/cluster/center"
	"gok/service/msg/protocol"
	"github.com/pkg/errors"
)

func NewRemoteService(serviceType string) *RemoteService {
	service := &RemoteService{
		serviceType: serviceType,
	}
	service.Init()
	return service
}

//type IRemoteService interface {
//	AsyncHandleRemoteMessage(string, *protocol.C2GS) *protocol.GS2C
//}

//远程调度服务 override IMessageService
type RemoteService struct {
	serviceType string //服务类型
}

func (this *RemoteService) Init() {
	center.ClusterCenter.AddServiceFactory(this.serviceType, &grpcServiceFactory{})
	center.ClusterCenter.SubscribeService(this.serviceType) //订阅指定类型服务
}

func (this *RemoteService) CanHandle(seq int32) bool {
	return center.ClusterCenter.CanHandle(this.serviceType, seq)
}

//指定服务器处理消息
//func (this *RemoteService) AsyncHandleRemoteMessage(serviceID string, message *protocol.C2GS) *protocol.GS2C  {
//	service := center.ClusterCenter.GetService(this.serviceType, serviceID)
//	if (service == nil) {
//		return INVALID_SERVICE_RESPONSE
//	}
//	result, _ := service.Request(message)
//	response, ok := result.(*protocol.GS2C)
//	if (ok) {
//		return response
//	}
//	return INVALID_SERVICE_RESPONSE
//}

func (this *RemoteService) AsyncHandleRemoteMessage(serviceID string, message *protocol.C2GS) error {
	service := center.ClusterCenter.GetService(this.serviceType, serviceID)
	if service == nil {
		return errors.New("remote service can not alloc " + this.serviceType)
	}
	return service.AsyncRequest(message)
}

func (this *RemoteService) AsyncBroadcastAllIgnore(serviceID string, message *protocol.C2GS) bool {
	services := center.ClusterCenter.GetAllServiceIgnoreID(this.serviceType, serviceID)
	if services == nil || len(services) == 0 {
		return false
	}
	for _, service := range services {
		service.AsyncRequest(message)
	}
	return true
}

func (this *RemoteService) AsyncHandlePriorityRemoteMessage(serviceID string, message *protocol.C2GS) error {
	service := center.ClusterCenter.GetService(this.serviceType, serviceID)
	if service == nil {
		service = center.ClusterCenter.AllocService(this.serviceType)
	}
	if service == nil {
		return errors.New("can not found removete service " + this.serviceType)
	}
	return service.AsyncRequest(message)
}

func (this *RemoteService) AsyncBroadcastAll(message *protocol.C2GS) bool {
	services := center.ClusterCenter.GetAllService(this.serviceType)
	if services == nil || len(services) == 0 {
		return false
	}
	for _, service := range services {
		service.AsyncRequest(message)
	}
	return true
}


func (this *RemoteService) AsyncHandleMessage(message *protocol.C2GS) error {
	service := center.ClusterCenter.AllocService(this.serviceType)
	if service == nil {
		return errors.New("remote service can not alloc " + this.serviceType)
	}
	return service.AsyncRequest(message)
}


func (this *RemoteService) HandleRemoteMessage(serviceID string, message *protocol.C2GS) *protocol.GS2C {
	service := center.ClusterCenter.GetService(this.serviceType, serviceID)
	if service == nil {
		return INVALID_SERVICE_RESPONSE
	}
	result, _ := service.Request(message)
	response, ok := result.(*protocol.GS2C)
	if ok {
		return response
	}
	return INVALID_SERVICE_RESPONSE
}

func (this *RemoteService) HandleMessage(message *protocol.C2GS) *protocol.GS2C {
	service := center.ClusterCenter.AllocService(this.serviceType)
	if service == nil {
		return INVALID_SERVICE_RESPONSE
	}
	result, _ := service.Request(message)
	response, ok := result.(*protocol.GS2C)
	if ok {
		return response
	}
	return INVALID_SERVICE_RESPONSE
}

//func (this *RemoteService) AsyncHandlePriorityRemoteMessage(serviceID string, message *protocol.C2GS) *protocol.GS2C {
//	service := center.ClusterCenter.GetService(this.serviceType, serviceID)
//	if service == nil {
//		service = center.ClusterCenter.AllocService(this.serviceType)
//	}
//	if service == nil {
//		return INVALID_SERVICE_RESPONSE
//	}
//	result, _ := service.Request(message)
//	response, ok := result.(*protocol.GS2C)
//	if ok {
//		return response
//	}
//	return INVALID_SERVICE_RESPONSE
//}

func (this *RemoteService) HandleChannelMessage(message *protocol.C2GS, channel IMessageChannel) {
	response := this.HandleMessage(message)
	channel.WriteMsg(response)
}

//获取消息服务类型
func (this *RemoteService) GetType() string {
	return this.serviceType
}


