/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/4/17
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package rpc

import (
	"gok/service"
	"gok/service/msg/protocol"
)

var CenterServiceProxy = &centerHandler{&rpcHandler{serviceType:service.SERVICE_CENTER_RPC}}

type centerHandler struct {
	*rpcHandler
}

func (this *centerHandler) GetOnNotices(status int32) *protocol.GetOnNoticesRet{
	request := &protocol.C2GS{
		Sequence: []int32{700},
		GetOnNotices: &protocol.GetOnNotices{
			Status: status,
		},
	}
	return this.HandleMessage(request).GetGetOnNoticesRet()
}