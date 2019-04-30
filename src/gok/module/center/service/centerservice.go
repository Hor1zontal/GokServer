package service

import (
	"gok/service/msg/protocol"
	baseservice "gok/service"
	"gok/module/center/conf"
	"gok/service/exception"
	"gok/module/center/order"
	"github.com/name5566/leaf/chanrpc"
	"gok/constant"
	"gok/module/center/notice"
)

var CenterRPCService *baseservice.GRPCService = nil

func Init(chanRpc *chanrpc.Server) {
	//中心服务器订阅用户服务,能够讲查询到的用户服务地址返回给客户端
	var centerService = baseservice.NewLocalService(baseservice.SERVICE_CENTER_RPC)

	centerService.RegisterHandler(700,  new(GetOnNoticesService))
	centerService.RegisterHandler(8,  new(GenOrderService))
	//发布登录服务到中心服务器
	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_MAIL_RPC)
	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_TRADE)
	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_USER_RPC)
	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_PASSPORT_RPC)
	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_STAR_RPC)
	////配置了RPC，需要发布服务到ZK
	CenterRPCService = baseservice.PublicRPCService1(centerService, conf.Server.RPCAddress, conf.Server.RPCPort, chanRpc)
}

func Close() {
	CenterRPCService.Close()
}

type GenOrderService struct {
}

func (service *GenOrderService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGenOrder()
	shopBase := conf.DATA.ShopBaseData[message.GetShopID()]
	if shopBase == nil {
		exception.GameException(exception.SHOP_BUY_ERROR)
	}
	orderID := order.GenTempOrder(message.GetUid(), shopBase.ID, shopBase.Value)
	response.GenOrderRet = &protocol.GenOrderRet{
		OrderID:orderID,
	}
}

type GetOnNoticesService struct {

}

func (service *GetOnNoticesService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGetOnNotices()
	//var notices []*protocol.Notice
	notices := make([]*protocol.Notice, 1)
	var noticeID int32 = 0
	if message.Status == constant.NOTICE_ON {
		onNotices := notice.NoticesManager.OnNotices
		if onNotices != nil && len(onNotices) > 0 {
			for _, noticeData := range onNotices {
				if noticeData.ID > noticeID {
					noticeID = noticeData.ID
				}
			}
			notices[0] = onNotices[noticeID].BuildProto()
		}
	}
	response.GetOnNoticesRet = &protocol.GetOnNoticesRet{Notices:notices}
	//BuildProto
}