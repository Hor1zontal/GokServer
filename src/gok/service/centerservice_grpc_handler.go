/*******************************************************************************
* Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
* All rights reserved.
* Date:
*     2018/8/20
* Contributors:
*     aliens idea(xiamen) Corporation - initial API and implementation
*     jialin.he <kylinh@gmail.com>
*******************************************************************************/
package service

import (
	"github.com/name5566/leaf/chanrpc"
	"gok/service/msg/protocol"
	"time"
)


const (
	commandRequest string = "request"
	commandReceive string = "receive"
)

func newServiceHandler(proxy *LocalService, chanRpc *chanrpc.Server) *serviceHandler {
	handler := &serviceHandler{chanRpc:chanRpc, proxy:proxy}
	if chanRpc != nil {
		chanRpc.Register(commandRequest, handler.request)
		chanRpc.Register(commandReceive, handler.receive)
	}
	return handler
}

type serviceHandler struct {
	chanRpc *chanrpc.Server
	proxy *LocalService
	suspended bool
}

func (this *serviceHandler) request(args []interface{}) {
	message := args[0].(*protocol.C2GS)
	server := args[1].(protocol.RPCService_RequestServer)
	response := this.proxy.HandleMessage(message)
	server.Send(response)
}

func (this *serviceHandler) receive(args []interface{}) {
	message := args[0].(*protocol.C2GS)
	this.proxy.HandleMessage(message)
}

func (this *serviceHandler) Request(request *protocol.C2GS, server protocol.RPCService_RequestServer) error {
	//异步调用支持
	if this.chanRpc != nil {
		this.chanRpc.Call0(commandRequest, request, server)
		return nil
	}
	if this.proxy != nil {
		response := this.proxy.HandleMessage(request)
		return server.Send(response)
	}
	return server.Send(INVALID_SERVICE_RESPONSE)
}

func (this *serviceHandler) Receive(server protocol.RPCService_ReceiveServer) error {
	for {
		if this.suspended {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		request, err := server.Recv()
		if err != nil {
			return err
		}
		if this.chanRpc != nil {
			this.chanRpc.Go(commandReceive, request)
		} else if this.proxy != nil {
			this.proxy.HandleMessage(request)
		}

	}
}
