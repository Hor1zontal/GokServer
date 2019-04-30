/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/5/6
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package service

import (
	"aliens/log"
	"runtime"
	"fmt"
	"github.com/name5566/leaf/chanrpc"
)


func NewMessageChannel(messageLimit int, handle func(msg interface{})) *MessageChannel {
	channel := &MessageChannel{}
	channel.Open(messageLimit)
	return channel
}

type IMessageChannel interface {
	WriteMsg(msg interface{})
	//SyncMessage(msg interface{})
	//Close()
	//SetUserData(data interface{})
	//UserData() interface{}
}

//type IMessageHandler interface {
//	handleMessage(msg interface{}) //处理消息
//}

type MessageChannel struct {
	server *chanrpc.Server
	//server        chan interface{} //管道
	open           bool
	//messageHandler func(msg interface{}) //消息处理handler
}

//向管道发送消息
func (this *MessageChannel) AcceptMessage(id interface{}, args ...interface{}) {
	if this.server == nil {
		return
	}
	this.server.Go(id, args...)
	//用户消息管道没开，不接受消息
	//if !this.IsOpen() {
	//	return
	//}
	//select {
	//case this.server <- message:
	//default:
	//	log.Debug("message server full %v - %v", this.server, message)
	//	//TODO 消息管道满了需要异常处理
	//}
}

//打开用户消息管道
func (this *MessageChannel) Open(messageLimit int) {
	this.server = chanrpc.NewServer(messageLimit)
	//this.server = make(chan interface{}, this.messageLimit)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				stackInfo := fmt.Sprintf("%s", buf[:n])
				log.Error("panic stack info %s", stackInfo)
			}
		} ()
		for {
			//只要消息管道没有关闭，就一直等待用户请求
			message, open := <-this.server.ChanCall
			if !this.open || !open {
				this.server = nil
				break
			}
			this.server.Exec(message)
			//this.messageHandler(message)
		}
		this.Close()
	}()

	//for {
	//	select {
	//	case <-closeSig:
	//		s.commandServer.Close()
	//		s.server.Close()
	//		for !s.g.Idle() || !s.client.Idle() {
	//			s.g.Close()
	//			s.client.Close()
	//		}
	//		return
	//	case ri := <-s.client.ChanAsynRet:
	//		s.client.Cb(ri)
	//	case ci := <-s.server.ChanCall:
	//		s.server.Exec(ci)
	//	case ci := <-s.commandServer.ChanCall:
	//		s.commandServer.Exec(ci)
	//	case cb := <-s.g.ChanCb:
	//		s.g.Cb(cb)
	//	case t := <-s.dispatcher.ChanTimer:
	//		t.Cb()
	//	}
	//}


	this.open = true
}

//关闭消息管道
func (this *MessageChannel) Close() {
	if !this.open {
		return
	}
	defer func() {
		recover()
	}()
	this.open = false
	if this.server != nil {
		this.server.Close()
		//close(this.server)
	}
}

//消息管道是否打开
func (this *MessageChannel) IsOpen() bool {
	return this.open && this.server != nil
}

//func (this *MessageChannel) SetUserData(data interface{}) {
//	this.userdata = data
//}
//
//func (this *MessageChannel) UserData() interface{} {
//	return this.userdata
//}
