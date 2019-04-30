///*******************************************************************************
// * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
// * All rights reserved.
// * Date:
// *     2017/5/12
// * Contributors:
// *     aliens idea(xiamen) Corporation - initial API and implementation
// *     jialin.he <kylinh@gmail.com>
// *******************************************************************************/
package service
//
//import (
//	//"aliens/game/command"
//	"github.com/name5566/leaf/gate"
//	"runtime/debug"
//	"gok/service/msg/command"
//	"sync"
//)
//
//func OpenChannelHandler(networkChannel IMessageChannel, handler IMessageHandler, maxMessage int) *ChannelHandler {
//	channelHandler := &ChannelHandler{
//		networkChannel: networkChannel,
//		handler:        handler,
//	}
//	channelHandler.OpenInnerChannel(maxMessage)
//	return channelHandler
//}
//
//type IChannelHandler interface {
//	SyncMessage(message interface{}) //接收消息
//	GateClose(gate gate.Agent)         //关闭网关
//}
//
//type ChannelHandler struct {
//	sync.RWMutex
//	networkChannel IMessageChannel //往外写的消息管道
//	innerChannel   IMessageChannel //往内写的消息管道
//	handler        IMessageHandler //服务处理容器
//}
//
////消息管道是否开启
//func (this *ChannelHandler) IsChannelOpen() bool {
//	return this.innerChannel != nil
//}
//
////收取系统消息
//func (this *ChannelHandler) SetHandler(handler IMessageHandler) {
//	this.handler = handler
//}
//
//func (this *ChannelHandler) GetNetworkChannel() IMessageChannel {
//	this.RLock()
//	defer this.RUnlock()
//	return this.networkChannel
//}
//
//func (this *ChannelHandler) SetNetworkChannel(networkChannel IMessageChannel) {
//	this.Lock()
//	defer this.Unlock()
//	this.networkChannel = networkChannel
//}
//
////收取系统消息
//func (this *ChannelHandler) SyncCommand(cmd command.Command, param ...interface{}) {
//	if this.innerChannel != nil {
//		this.innerChannel.WriteMsg(command.NewSystemMessage(cmd, param...))
//	}
//}
//
//func (this *ChannelHandler) GateClose(gate gate.Agent) {
//	this.Close()
//}
//
////收取消息
//func (this *ChannelHandler) SyncMessage(message interface{}) {
//	if this.innerChannel != nil {
//		this.innerChannel.WriteMsg(message)
//	}
//}
//
////往连接客户端写消息
//func (this *ChannelHandler) SendMessage(message interface{}) {
//	this.RLock()
//	defer this.RUnlock()
//	if this.networkChannel != nil {
//		this.networkChannel.WriteMsg(message)
//	}
//}
//
////是否在线
//func (this *ChannelHandler) IsOnline() bool {
//	this.RLock()
//	defer this.RUnlock()
//	return this.networkChannel != nil
//}
//
//func (this *ChannelHandler) handleMessage(msg interface{}) {
//	defer func() {
//		//处理消息异常
//		if err := recover(); err != nil {
//			debug.PrintStack()
//			this.Close()
//		}
//	}()
//	if this.handler != nil {
//		this.handler.handleMessage(msg)
//	}
//}
//
////打开收消息管道
//func (this *ChannelHandler) OpenInnerChannel(maxMessage int) {
//	if this.innerChannel != nil {
//		return
//	}
//	server := &MessageChannel{
//		messageLimit:   maxMessage,
//		messageHandler: this,
//	}
//	server.Open()
//	this.innerChannel = server
//}
//
////关闭收消息管道
//func (this *ChannelHandler) CloseInnerChannel() {
//	if this.innerChannel != nil {
//		this.innerChannel.Close()
//		this.innerChannel = nil
//	}
//}
//
////关闭收消息管道
//func (this *ChannelHandler) CloseNetworkChannel() {
//	this.Lock()
//	defer this.Unlock()
//	if this.networkChannel != nil {
//		this.networkChannel.Close()
//		this.networkChannel = nil
//	}
//}
//
////关闭所有管道
//func (this *ChannelHandler) Close() {
//	this.CloseInnerChannel()
//	this.CloseNetworkChannel()
//}
