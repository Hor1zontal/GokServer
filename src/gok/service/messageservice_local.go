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
	"gok/service/msg/protocol"
	//	"runtime/debug"
	//"github.com/golang/protobuf/proto"
	//"gok/service/exception"
	////"aliens/log"
	//"runtime"
	//"github.com/name5566/leaf/log"
	"gok/service/exception"
	"runtime/debug"
	"aliens/log"
	"gok/constant"
	"time"
	"gok/service/lpc"
)

var ExceptionSeq = []int32{1000}

//新建本地服务管理类
func NewLocalService(serviceType string) *LocalService {
	service := &LocalService{
		//open: true,
		//dealTotal : 0,
		serviceType: serviceType,
		services:    make(map[int32]IServiceHandler),
		filter: nil,
	}
	return service
}

type IServiceHandler interface {
	Request(request *protocol.C2GS, response *protocol.GS2C, channel IMessageChannel)
}

//override IMessageService
type LocalService struct {
	//sync.RWMutex
	//open           bool //服务是否开放
	//dealTotal      int64 //当前处理中的消息数量
	serviceType    string
	services       map[int32]IServiceHandler
	filter         func(message *protocol.C2GS) bool
	callbackFilter func(seq int32) bool
}

func (this *LocalService) GetType() string {
	return this.serviceType
}

//新增过滤器
func (this *LocalService) SetFilter(filter func(message *protocol.C2GS) bool) {
	this.filter = filter
}

func (this *LocalService) SetCallbackFilter(callbackFilter func(seq int32) bool) {
	this.callbackFilter = callbackFilter
}


func (this *LocalService) GetServiceSeq() []int32 {
	results := []int32{}
	for seq, _ := range this.services {
		results = append(results, seq)
	}
	return results
}

func (this *LocalService) CleanFilter() {
	this.filter = nil
}

//注册消息服务处理句柄
func (this *LocalService) RegisterHandler(seq int32, service IServiceHandler) {
	this.services[seq] = service
}

//是否能够处理指定编号的消息
func (this *LocalService) CanHandle(seq int32) bool {
	return this.services[seq] != nil
}

//阻塞式调用消息接口
func (this *LocalService) HandleMessage(message *protocol.C2GS) *protocol.GS2C {
	response := &protocol.GS2C{}
	this.channelMessage(message, response, nil)
	return response
}

//当前处理中的消息数据
//func (this *LocalService) GetDealTotal() int64{
//	return this.dealTotal
//}

func (this *LocalService) Close() {
	//this.RLock()
	//defer this.RUnlock()
	//this.open = false
	//timeout := 100 * time.Millisecond
	////10次定时后还没处理完毕直接超时返回
	//for i:= 0; i < 10; i ++ {
	//	time.Sleep(time.Duration(i) * timeout)
	//	if (this.dealTotal <= 0) {
	//		return
	//	}
	//	log.Debug("%v [%v] undeal message : %v", time.Now(), this.serviceType, this.dealTotal)
	//}
}

//处理管道消息，响应消息写传入管道
func (this *LocalService) HandleChannelMessage(message *protocol.C2GS, channel IMessageChannel) {
	this.channelMessage(message, &protocol.GS2C{}, channel)
}

func (this *LocalService) channelMessage(message *protocol.C2GS, response *protocol.GS2C, channel IMessageChannel) {
	if len(message.GetSequence()) == 0 {
		return
	}
	if this.filter != nil && this.filter(message) {
		return
	}
	startTime := time.Now()

	response.Session = message.Session
	response.Sequence = message.Sequence
	// 输出收到的消息的内容
	defer func() {
		//处理消息异常
		if err := recover(); err != nil {
			response.Sequence = ExceptionSeq
			switch err.(type) {
			case exception.GameCode:
				response.ResultPush = &protocol.ResultPush{
					Result: int32(err.(exception.GameCode)),
				}
				break
			default:
				log.Error("%v - %v", err, string(debug.Stack()))
				//log.Error("")
				//debug.PrintStack()
				response.ResultPush = &protocol.ResultPush{
					Result: exception.SERVICE_INTERNAL_ERR,
				}
			}
		}
		//log.Debug("response %v-%v", this.serviceType, response)
		if channel != nil {
			if this.callbackFilter == nil || !this.callbackFilter(message.GetSequence()[0]) {
				channel.WriteMsg(response)
			}
		}
		if constant.ES_LOG {
			callDuration := time.Now().Sub(startTime)
			lpc.StatisticsHandler.AddServiceStatistic(this.serviceType, message.Sequence[0], callDuration.Seconds())
		}
	}()



	//处理多条消息 TODO 要考虑消息服务的前后依赖关系，比如失败了能否继续
	for _, seq := range message.GetSequence() {
		if service, ok := this.services[seq]; ok {
			//log.Debug("request %v-%v", this.serviceType, message)
			service.Request(message, response, channel)
		} else {
			exception.GameException(exception.SERVICE_NO_INVALID)
		}
	}
}
