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
)

type httpServiceFactory struct {
}

func (this *httpServiceFactory) CreateService(data []byte) core.IService {
	service := &HTTPService{}
	json.Unmarshal(data, service)
	return service
}

func PublicHTTPService(serviceType string, address string) *HTTPService {
	if !center.ClusterCenter.IsConnect() {
		panic(serviceType + " cluster center is not connected")
	}
	service := &HTTPService{
		id:          center.GetServerNode(),
		serviceType: serviceType,
		Address:     address,
	}
	//center.ClusterCenter.AddServiceFactory(service.serviceType, &HTTPServiceFactory{})
	//websocket服务启动成功,则发布到中心服务器
	if !center.ClusterCenter.PublicService(service) {
		panic(service.serviceType + " http service can not be public")
	}
	return service
}

func SubscribeHTTPService(serviceType string) {
	center.ClusterCenter.AddServiceFactory(serviceType, &httpServiceFactory{})
	center.ClusterCenter.SubscribeService(serviceType)
}

type HTTPService struct {
	Address     string `json:"address"` //服务访问地址 写入到中心服务器供外部调用
	id          string //服务ID
	serviceType string //服务类型
}

func (this *HTTPService) GetDesc() string {
	return "http service"
}

func (this *HTTPService) GetID() string {
	return this.id
}

func (this *HTTPService) GetType() string {
	return this.serviceType
}

func (this *HTTPService) SetID(id string) {
	this.id = id
}

func (this *HTTPService) SetType(serviceType string) {
	this.serviceType = serviceType
}

//启动服务
func (this *HTTPService) Start() bool {
	return true
}

//连接服务
func (this *HTTPService) Connect() bool {
	return true
}

//比较服务是否冲突
func (this *HTTPService) Equals(other core.IService) bool {
	otherService, ok := other.(*HTTPService)
	if !ok {
		return false
	}
	return this.serviceType == otherService.serviceType && this.Address == otherService.Address
}

//服务是否本进程启动的
func (this *HTTPService) IsLocal() bool {
	return true
}

//关闭服务
func (this *HTTPService) Close() {
}

//向服务请求消息
func (this *HTTPService) Request(in interface{}) (interface{}, error) {
	return nil, nil
}

func (this *HTTPService) AsyncRequest(in interface{}) error { //异步请求
	return nil
}
