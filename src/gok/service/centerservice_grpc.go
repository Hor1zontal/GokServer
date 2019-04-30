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
	"aliens/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"gok/module/cluster/center"
	"gok/service/msg/protocol"
	"github.com/name5566/leaf/chanrpc"
	"time"
)

type grpcServiceFactory struct {
}

func (this *grpcServiceFactory) CreateService(data []byte) core.IService {
	service := &GRPCService{}
	json.Unmarshal(data, service)
	return service
}
//
//func PublicRPCService(proxy *LocalService, address string, port int) *GRPCService {
//	return PublicRPCService1(proxy, address, port, nil)
//}

func PublicRPCService1(proxy *LocalService, address string, port int, chanRpc *chanrpc.Server) *GRPCService {
	if !center.ClusterCenter.IsConnect() {
		panic(proxy.serviceType + " cluster center is not connected")
	}
	service := &GRPCService{
		id:          center.GetServerNode(),
		serviceType: proxy.serviceType,
		Address:     address,
		port:        port,
		proxy:       proxy,
		rpcServer:   newServiceHandler(proxy, chanRpc),
	}
	if !service.Start() {
		panic(service.serviceType + " rpc service can not be start")
	}
	//RPC启动成功,则发布到中心服务器
	if !center.ClusterCenter.PublicService(service) {
		panic(service.serviceType + " rpc service can not be start")
	}
	return service
}

type GRPCService struct {
	Address     string        `json:"address"` //服务访问地址 写入到中心服务器供外部调用
	id          string                         //服务ID
	serviceType string                         //服务类型
	port        int                            //RPC服务端口
	proxy       *LocalService                  //调用句柄 //handler protocol.RPCServiceServer //当服务为本地启动
	server      *grpc.Server

	caller         protocol.RPCServiceClient //当服务为远程
	receiveClient  protocol.RPCService_ReceiveClient //当服务为远程,处理推送句柄
	rpcServer      protocol.RPCServiceServer

	client      *grpc.ClientConn
}

func (this *GRPCService) GetDesc() string {
	data, _ := json.Marshal(this.proxy.GetServiceSeq())
	return string(data)
}

func (this *GRPCService) GetID() string {
	return this.id
}

func (this *GRPCService) GetType() string {
	return this.serviceType
}

func (this *GRPCService) SetID(id string) {
	this.id = id
}

func (this *GRPCService) SetType(serviceType string) {
	this.serviceType = serviceType
}

//启动服务
func (this *GRPCService) Start() bool {
	if this.proxy == nil {
		log.Error("service handler can not nil")
		return false
	}
	//this.handler = handler
	this.server = grpc.NewServer()
	address := ":" + strconv.Itoa(this.port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Error("failed to listen: %v", err)
		return false
	}
	protocol.RegisterRPCServiceServer(this.server, this.rpcServer)
	go func() {
		this.server.Serve(lis)
		log.Info("rpc service %v stop", this.serviceType)
	}()
	return true
}

//连接服务
func (this *GRPCService) Connect() bool {
	if this.caller != nil {
		return true
	}
	conn, err := grpc.Dial(this.Address, grpc.WithInsecure())
	if err != nil {
		log.Error("did not connect: %v", err)
		return false
	}
	this.client = conn
	this.caller = protocol.NewRPCServiceClient(this.client)
	receiveClient, err := this.caller.Receive(context.Background())
	if err != nil {
		log.Error("%v", err)
		return false
	}
	this.receiveClient = receiveClient
	//log.Debug("connect gok service %v-%v", this.serviceType, this.id)
	return true
}

//比较服务是否冲突
func (this *GRPCService) Equals(other core.IService) bool {
	otherService, ok := other.(*GRPCService)
	if !ok {
		return false
	}
	return this.serviceType == otherService.serviceType && this.Address == otherService.Address
}

//服务是否本进程启动的
func (this *GRPCService) IsLocal() bool {
	return this.proxy != nil
}

//关闭服务
func (this *GRPCService) Close() {
	if this.client != nil {
		this.client.Close()
		this.client = nil
	}
	if this.server != nil {
		this.server.Stop()
		this.server = nil
	}
	if this.proxy != nil {
		this.proxy.Close()
	}
}

//向服务请求消息
func (this *GRPCService) Request(in interface{}) (interface{}, error) {
	request, ok := in.(*protocol.C2GS)
	if !ok {
		return nil, nil
	}
	if this.IsLocal() {
		return this.proxy.HandleMessage(request), nil
	}
	if this.caller == nil {
		return nil, errors.New("service is not initial")
	}
	client, err := this.caller.Request(newTimeoutContext(), request)
	if err != nil {
		log.Debug("call rpc err %v", err)
		return nil, err
	}
	response, err := client.Recv()
	return response, err
}

func (this *GRPCService) AsyncRequest(in interface{}) error {//异步请求
	request, ok := in.(*protocol.C2GS)
	if !ok {
		return nil
	}
	if this.receiveClient == nil {
		return errors.New("receive client is not initial")
	}
	response := this.receiveClient.Send(request)
	return response
}


func newTimeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
	return ctx
}



