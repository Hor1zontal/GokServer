/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/5/27
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package rpc

import (
	"gok/service"
	"gok/service/msg/protocol"
)

var TradeServiceProxy = &tradeHandler{&rpcHandler{serviceType:service.SERVICE_TRADE}}

type tradeHandler struct {
	*rpcHandler
}

//请求RPC调用接口层
func (this *tradeHandler) Call(uid int32, request *protocol.C2GS) *protocol.GS2C {
	request.Param = uid
	//node := center.TradeHashring.GetNode(character.Int32ToString(uid))
	return this.HandleMessage(request)
}

func (this *tradeHandler) NoErrorCall(uid int32, request *protocol.C2GS) *protocol.GS2C {
	request.Param = uid
	return this.NoErrorHandleMessage(request)
}

func (this *tradeHandler) AddSale(uid int32, itemID int32) *protocol.AddSaleRet {
	request := &protocol.C2GS{
		Sequence: []int32{530},
		AddSale: &protocol.AddSale{
			Id:     uid,
			ItemID: itemID,
		},
	}
	return this.HandleMessage(request).GetAddSaleRet()
}

func (this *tradeHandler) RemoveSale(uid int32, itemID int32) {
	request := &protocol.C2GS{
		Sequence: []int32{531},
		RemoveSale: &protocol.RemoveSale{
			Id:     uid,
			ItemID: itemID,
		},
	}
	this.HandleMessage(request)
}

func (this *tradeHandler) GetSale(uid int32) *protocol.Sale {
	request := &protocol.C2GS{
		Sequence: []int32{532},
		GetSale: &protocol.GetSale{
			Id: uid,
		},
	}
	return this.HandleMessage(request).GetGetSaleRet().GetSale()
}

func (this *tradeHandler) GetSales(uids []int32) []*protocol.Sale {
	request := &protocol.C2GS{
		Sequence: []int32{533},
		GetSales: &protocol.GetSales{
			Id: uids,
		},
	}
	return this.HandleMessage(request).GetGetSalesRet().GetSales()
}