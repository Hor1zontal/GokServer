//管理网络消息的处理
package service

import (
	baseservice "gok/service"
	"gok/module/search/conf"
	"github.com/name5566/leaf/chanrpc"
	"gok/service/msg/protocol"
	"gok/module/search/core"
	"gok/service/rpc"
	"gok/module/cluster/center"
	"gok/constant"
)

var SearchRPCService *baseservice.GRPCService = nil

func Init(chanRpc *chanrpc.Server) {
	localService := baseservice.NewLocalService(baseservice.SERVICE_SEARCH_RPC)
	localService.RegisterHandler(40,  new(RandomTargetService))
	//ocalService.RegisterHandler(260, new(SearchItemService))       //搜索圣物
	localService.RegisterHandler(540, new(UpdateSearchDataService))  //更新搜索数据
	localService.RegisterHandler(541, new(UpdateHelpSearchDataService))
	localService.RegisterHandler(542, new(RandomHelpTargetService))
	localService.RegisterHandler(543, new(UpdateRandomStarService)) //更新随机星球



	SearchRPCService = baseservice.PublicRPCService1(localService, conf.Server.RPCAddress, conf.Server.RPCPort, chanRpc)

	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_STAR_RPC)
	//需要知道主节点做更新
	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_SEARCH_RPC)
}

func Close() {
	SearchRPCService.Close()
}


//随机生成目标用户
type RandomTargetService struct {
}

func (service *RandomTargetService) Request(request *protocol.C2GS, response *protocol.GS2C, messageChannel baseservice.IMessageChannel) {
	message := request.GetRandomTarget()
	message.Filter = append(message.Filter, message.GetUid())
	randomCount := int(message.GetNum())  //随机数量
	robotFilter := message.GetRobotFilter() //是否过滤机器人

	targets := core.StarSearcher.RandomEventStar(message.GetEventType(), message.GetUid(), message.GetMutualID(), robotFilter, randomCount, message.GetAlwaysTarget())

	response.RandomTargetRet = &protocol.RandomTargetRet{
		EventID: message.GetEventID(),
		Targets: targets,
	}
	//response.RandomTargetRet.Targets =
	//if message.GetGuide() {
	//	response.RandomTargetRet.Targets = core.StarSearcher.RandomRobot(randomCount)
	//} else {
	//
	//}
	//stars := session.StarManager.RandomTargetStar(message, int(message.GetNum()))

}


//type SearchItemService struct {
//	//重置槽
//}
//
//func (service *SearchItemService) Request(request *protocol.C2GS, response *protocol.GS2C, messageChannel baseservice.IMessageChannel) {
//	message := request.GetSearchItem()
//	results := core.StarSearcher.SearchUsers(message.GetId(), message.GetStarType(), message.GetItemIDs())
//	response.SearchItemRet = &protocol.SearchItemRet{
//		Strangers: results,
//	}
//}


//更新索引数据
type UpdateSearchDataService struct {

}

func (service *UpdateSearchDataService) Request(request *protocol.C2GS, response *protocol.GS2C, messageChannel baseservice.IMessageChannel) {
	message := request.GetUpdateSearchData()
	opts := message.GetOpts()
	if opts == nil || len(opts) == 0 {
		return
	}

	for starID, opt := range opts {
		_, ok := opt.GetOpt()[constant.SEARCH_OPT_REMOVE_STAR]
		//有删除操作直接处理删除操作即可
		if ok {
			core.StarSearcher.OPT(starID, constant.SEARCH_OPT_REMOVE_STAR, 0, message.GetSync())
		} else {
			for optType, updateData := range opt.GetOpt() {
				core.StarSearcher.OPT(starID, optType, updateData, message.GetSync())
			}
		}
	}

	//不是同步消息，需要把操作同步给其他节点
	if !message.GetSync() {
		message.Sync = true
		//同步更新其他搜索服务器
		rpc.SearchServiceProxy.SyncOtherSearchService(request, center.GetServerNode())
		//this.AsyncBroadcastAllMessage(request)
	}
}

//更新索引数据
type UpdateHelpSearchDataService struct {

}

func (service *UpdateHelpSearchDataService) Request(request *protocol.C2GS, response *protocol.GS2C, messageChannel baseservice.IMessageChannel) {
	message := request.GetUpdateSearchHelpData()
	core.ItemHelpSearcher.Opt(message.GetUid(), message.GetOpt(), message.GetParam(), message.GetSync())
	//不是同步消息，需要把操作同步给其他节点
	if !message.GetSync() {
		message.Sync = true
		//同步更新其他搜索服务器
		rpc.SearchServiceProxy.SyncOtherSearchService(request, center.GetServerNode())
		//this.AsyncBroadcastAllMessage(request)
	}
}


type RandomHelpTargetService struct {

}

func (service *RandomHelpTargetService) Request(request *protocol.C2GS, response *protocol.GS2C, messageChannel baseservice.IMessageChannel) {
	message := request.GetRandomHelpTarget()
	targets := core.ItemHelpSearcher.RandomTargets(message.GetUid(), message.GetStarType(), int(message.GetCount()))

	if len(targets) > 0 {
		//同步给其他玩家收到帮组消息的索引
		message := rpc.SearchServiceProxy.BuildHelpDatas(targets, constant.SEARCH_OPT_UPDATE_RECEIVE_HELP, 0, true)
		rpc.SearchServiceProxy.SyncOtherSearchService(message, center.GetServerNode())
	}

	response.RandomHelpTargetRet = &protocol.RandomHelpTargetRet{
		Targets:targets,
	}
}


type UpdateRandomStarService struct {

}

func (service *UpdateRandomStarService) Request(request *protocol.C2GS, response *protocol.GS2C, messageChannel baseservice.IMessageChannel) {
	message := request.GetUpdateRandomStar()
	core.StarSearcher.UpdateRandomStar(message.GetStarID(), message.GetOpt())
}

