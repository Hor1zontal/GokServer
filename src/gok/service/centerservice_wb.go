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
	"aliens/center/core"
	"encoding/json"
	"gok/module/cluster/center"
	networkConf "gok/module/gate/conf"
)

type wbServiceFactory struct {
}

func (this *wbServiceFactory) CreateService(data []byte) core.IService {
	service := &WBService{}
	json.Unmarshal(data, service)
	return service
}

func PublicWBService(serviceType string) *WBService {
	if !center.ClusterCenter.IsConnect() {
		panic(serviceType + " cluster center is not connected")
	}
	service := &WBService{
		id:          center.GetServerNode(),
		serviceType: serviceType,
		Address:     networkConf.Config.PublicWSAddress,
	}
	//center.ClusterCenter.AddServiceFactory(service.serviceType, &wbServiceFactory{})
	//websocket服务启动成功,则发布到中心服务器
	if !center.ClusterCenter.PublicService(service) {
		panic(service.serviceType + " wb service can not be public")
	}
	return service
}

func SubscribeWBService(serviceType string) {
	center.ClusterCenter.AddServiceFactory(serviceType, &wbServiceFactory{})
	center.ClusterCenter.SubscribeService(serviceType)
}

type WBService struct {
	Address     string `json:"address"` //服务访问地址 写入到中心服务器供外部调用
	id          string //服务ID
	serviceType string //服务类型
}

func (this *WBService) GetDesc() string {
	return "websocket service"
}

func (this *WBService) GetID() string {
	return this.id
}

func (this *WBService) GetType() string {
	return this.serviceType
}

func (this *WBService) SetID(id string) {
	this.id = id
}

func (this *WBService) SetType(serviceType string) {
	this.serviceType = serviceType
}

//启动服务
func (this *WBService) Start() bool {
	return true
}

//连接服务
func (this *WBService) Connect() bool {
	return true
}

//比较服务是否冲突
func (this *WBService) Equals(other core.IService) bool {
	otherService, ok := other.(*WBService)
	if !ok {
		return false
	}
	return this.serviceType == otherService.serviceType && this.Address == otherService.Address
}

//服务是否本进程启动的
func (this *WBService) IsLocal() bool {
	return true
}

//关闭服务
func (this *WBService) Close() {
}

//向服务请求消息
func (this *WBService) Request(in interface{}) (interface{}, error) {
	return nil, nil
}

func (this *WBService) AsyncRequest(in interface{}) error { //异步请求
	return nil
}
