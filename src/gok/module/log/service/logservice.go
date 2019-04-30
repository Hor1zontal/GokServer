package service
//
//import (
//
//	"gok/service/msg/protocol"
//	baseservice "gok/service"
//	"aliens/log"
//	"gok/module/log/conf"
//)
//
//var LogService = baseservice.NewLocalService(baseservice.SERVICE_LOG_RPC)
//
//var PassportRPCService *baseservice.GRPCService = nil
//
//func Close() {
//	PassportRPCService.Close()
//}
//
//func init() {
//	LogService.RegisterHandler(550, new(AddLogService)) //新增订单日志
//	PassportRPCService = baseservice.PublicRPCService(LogService, conf.Config.RPCAddress, conf.Config.RPCPort)
//	//order_record addOrderRecord = 550;
//}
//
//
////登录账号服务器请求
//type AddLogService struct {
//}
//
//func (service *AddLogService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
//	message := request.GetAddLog()
//	log.Debug("add log %v", message.OrderRecord)
//}