//管理网络消息的处理
package service

import (
	"gok/service/msg/protocol"
	baseservice "gok/service"
	"gok/module/trade/conf"
	"gok/module/trade/core"
	"github.com/name5566/leaf/chanrpc"
)

var TradeLocalService = baseservice.NewLocalService(baseservice.SERVICE_TRADE)
var TradeRPCService *baseservice.GRPCService = nil

//初始化事件服务消息
func Init(chanRpc *chanrpc.Server) {

	TradeLocalService.RegisterHandler(470, new(GoodsInfoService))
	TradeLocalService.RegisterHandler(471, new(BuyGoodsService))
	TradeLocalService.RegisterHandler(472, new(PublicGoodsService))
	TradeLocalService.RegisterHandler(473, new(CancelGoodsService))



	TradeLocalService.RegisterHandler(474, new(GetHelpItemHistoryService))
	TradeLocalService.RegisterHandler(475, new(PublicItemHelpService))
	//TradeLocalService.RegisterHandler(476, new(CancelItemHelpService))
	TradeLocalService.RegisterHandler(477, new(DrawItemHelpService))
	TradeLocalService.RegisterHandler(478, new(GetItemHelpService))
	TradeLocalService.RegisterHandler(479, new(LootHelpItemService))
	TradeLocalService.RegisterHandler(480, new(HelpItemService))


	TradeLocalService.RegisterHandler(530, new(AddSaleService))
	TradeLocalService.RegisterHandler(531, new(RemoveSaleService))
	TradeLocalService.RegisterHandler(532, new(GetSaleService))
	TradeLocalService.RegisterHandler(533, new(GetSalesService))
	
	
	


	//配置了RPC，需要发布服务到ZK
	TradeRPCService = baseservice.PublicRPCService1(TradeLocalService, conf.Server.RPCAddress, conf.Server.RPCPort, chanRpc)
}

func Close() {
	TradeLocalService.Close()
	TradeRPCService.Close()
}

type GoodsInfoService struct {
}

func (service *GoodsInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGetGoodsInfo()
	response.GetGoodsInfoRet = &protocol.GetGoodsInfoRet{
		Goods:core.Exchanger.GetGoodsInfo(message.GetUid()),
	}

}

type BuyGoodsService struct {
}

func (service *BuyGoodsService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetBuyGoods()
	remain, fee := core.Exchanger.BuyGoods(message.GetUid(), message.GetItemid(), message.GetNum(), message.GetSocial())
	response.BuyGoodsRet = &protocol.BuyGoodsRet{
		Num: remain,
		Fee: fee,
	}
}

type PublicGoodsService struct {
}

func (service *PublicGoodsService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetPublicGoods()
	core.Exchanger.PublicGoods(message.GetUid(), message.GetGoods())
	response.PublicGoodsRet = &protocol.PublicGoodsRet{
		Result: true,
	}
}

type CancelGoodsService struct {
}

func (service *CancelGoodsService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetCancelGoods()
	goods := core.Exchanger.CancelGoods(message.GetUid(), message.GetItemid())
	response.CancelGoodsRet = &protocol.CancelGoodsRet{
		Goods: goods,
	}
}



type PublicItemHelpService struct {
}

func (service *PublicItemHelpService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetPublicItemHelp()
	help := core.HelpHandler.PublicItemHelp(request.GetParam(), message.GetItemID())
	response.PublicItemHelpRet = &protocol.PublicItemHelpRet{
		ItemHelp:help,
	}
}


//type CancelItemHelpService struct {
//}
//
//func (service *CancelItemHelpService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
//	//message := request.GetCancelItemHelp()
//	result := core.HelpHandler.CancelItemHelp(request.GetParam())
//	response.CancelItemHelpRet = &protocol.CancelItemHelpRet{
//		Result:result,
//	}
//}


type DrawItemHelpService struct {
}

func (service *DrawItemHelpService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetDrawItemHelp()
	itemNum, help := core.HelpHandler.DrawItemHelp(request.GetParam(), message.GetItemID(), message.GetCancel())
	response.DrawItemHelpRet = &protocol.DrawItemHelpRet{
		ItemID:help.GetItemID(),
		ItemNum:itemNum,
		ItemHelp:help,
	}
}


type GetItemHelpService struct {
}

func (service *GetItemHelpService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGetItemHelp()
	help := core.HelpHandler.GetItemHelp(request.GetParam(), message.GetId())
	response.GetItemHelpRet = &protocol.GetItemHelpRet{
		ItemHelp:help,
	}
}


type LootHelpItemService struct {
}

func (service *LootHelpItemService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetLootHelpItem()
	bool, cost, help := core.HelpHandler.LootItemHelp(message.GetUid(), request.GetParam(), message.GetItemID(), message.GetPower(),
		message.GetCosts(), message.GetLimit(), message.GetProbs(), message.GetAddProb(), message.GetIsWatchAd(), message.GetEachFollow())
	response.LootHelpItemRet = &protocol.LootHelpItemRet{
		Result:bool,
		Cost:cost,
		ItemHelp:help,
	}
}


type HelpItemService struct {
}

func (service *HelpItemService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetHelpItem()
	help := core.HelpHandler.HelpItemHelp(message.GetUid(), request.GetParam(), message.GetItemID(), message.GetLimit(), message.GetEachFollow())
	response.HelpItemRet = &protocol.HelpItemRet{
		Result:true,
		ItemHelp:help,
	}
}

type GetHelpItemHistoryService struct {
}

func (service *GetHelpItemHistoryService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGetHelpItemHistory()
	count, itemHelpHistory := core.HelpHandler.GetHelpItemHistory(message.GetUid(), message.GetSkip(), message.GetLimit(), message.GetCount())
	response.GetHelpItemHistoryRet = &protocol.GetHelpItemHistoryRet{
		ItemHelp:itemHelpHistory,
		Count:count,
	}
}

type AddSaleService struct {
}

func (service *AddSaleService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetAddSale()
	//新增一条出售的朋友圈消息
	//moment := core.Moments.AddMoment1(message.GetId(), constant.MOMENTS_TYPE_SALE, message.GetItemID())
	sale := core.Exchanger.PublicSale(message.GetId(), message.GetItemID())

	response.AddSaleRet = &protocol.AddSaleRet{
		Sale:sale,
	}
}

type RemoveSaleService struct {
}

func (service *RemoveSaleService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetRemoveSale()
	core.Exchanger.RemoveSale(message.GetId(), message.GetItemID())
	response.RemoveSaleRet = &protocol.RemoveSaleRet{
		Result:true,
	}
}

type GetSaleService struct {
}

func (service *GetSaleService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGetSale()
	sale := core.Exchanger.GetSale(message.GetId())
	response.GetSaleRet = &protocol.GetSaleRet{Sale: sale}
}

type GetSalesService struct {
}

func (service *GetSalesService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGetSales()

	results := []*protocol.Sale{}
	for _, uid := range message.GetId() {
		sale := core.Exchanger.GetSale(uid)
		if (sale != nil) {
			results = append(results, sale)
		}
	}
	response.GetSalesRet = &protocol.GetSalesRet{Sales:results}
}


