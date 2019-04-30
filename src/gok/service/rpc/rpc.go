/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/4
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package rpc

import (
	"gok/service/msg/protocol"
	"gok/service"
	"aliens/log"
)


type rpcHandler struct {
	serviceType string
	callStatistics map[string]int64
}

func (this *rpcHandler) AddCall(node string) {
	if this.callStatistics == nil {
		this.callStatistics = make(map[string]int64)
	}
	this.callStatistics[node] += 1
}

func (this *rpcHandler) PrintCallStatistics() {
	log.Debug("%v call statistics %v", this.serviceType, this.callStatistics)
	this.callStatistics = make(map[string]int64)
}


//rpc 消息应该能继续往上抛
func (this *rpcHandler) HandleMessage(message *protocol.C2GS) *protocol.GS2C {
	return service.ServiceManager.AssertHandleRemoteMessage(this.serviceType, message)
}

//如果是错误码不会继续往外抛
func (this *rpcHandler) NoErrorHandleMessage(message *protocol.C2GS) *protocol.GS2C {
	return service.ServiceManager.HandleRemoteMessage(this.serviceType, message)
}

//rpc 转发到指定节点消息
func (this *rpcHandler) AssertHandleNodeMessage(message *protocol.C2GS, node string) *protocol.GS2C {
	return service.ServiceManager.AssertSyncHandleRemoteMessage(this.serviceType, node, message)
}

func (this *rpcHandler) HandleNodeMessage(message *protocol.C2GS, node string) *protocol.GS2C {
	return service.ServiceManager.SyncHandleRemoteMessage(this.serviceType, node, message)
}

//rpc 异步处理消息
func (this *rpcHandler) AsyncHandleMessage(message *protocol.C2GS) error {
	return service.ServiceManager.AsyncHandleMessage(this.serviceType, message)
}

func (this *rpcHandler) AsyncHandleNodeMessage(message *protocol.C2GS, node string) error {
	return service.ServiceManager.AsyncHandleRemoteMessage(this.serviceType, node, message)
}

func (this *rpcHandler) AsyncBroadcastAllRemoteIgnore(message *protocol.C2GS, ignoreNode string) bool {
	return service.ServiceManager.AsyncBroadcastAllRemoteIgnore(this.serviceType, ignoreNode, message)
}

func (this *rpcHandler) AsyncBroadcastAllMessage(message *protocol.C2GS) bool {
	return service.ServiceManager.AsyncBroadcastAllRemote(this.serviceType, message)
}

func (this *rpcHandler) AsyncHandlePriorityRemoteMessage(message *protocol.C2GS, node string) {
	service.ServiceManager.AsyncHandlePriorityRemoteMessage(this.serviceType, node, message)
}


