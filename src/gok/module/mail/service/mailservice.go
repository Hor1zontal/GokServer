package service

import (
	"gok/service/msg/protocol"
	baseservice "gok/service"
	"gok/module/mail/conf"
	"gok/module/mail/mail"
	"github.com/name5566/leaf/chanrpc"
)

var CenterRPCService *baseservice.GRPCService = nil

func Init(chanRpc *chanrpc.Server) {
	//中心服务器订阅用户服务,能够讲查询到的用户服务地址返回给客户端
	var centerService = baseservice.NewLocalService(baseservice.SERVICE_MAIL_RPC)

	//发布登录服务到中心服务器
	centerService.RegisterHandler(70, new(GetAllMailService)) //获取所有邮件
	centerService.RegisterHandler(71, new(DrawMailService))   //领取邮件
	centerService.RegisterHandler(72, new(RemoveMailService)) //删除邮件
	centerService.RegisterHandler(73, new(CreateMailService)) //发送邮件
	centerService.RegisterHandler(74, new(GetMailService)) //获取时间段内的邮件

	////配置了RPC，需要发布服务到ZK
	CenterRPCService = baseservice.PublicRPCService1(centerService, conf.Server.RPCAddress, conf.Server.RPCPort, chanRpc)
	//baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_MAIL_RPC)
}

func Close() {
	CenterRPCService.Close()
}


type GetMailService struct {
}

func (service *GetMailService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGetMail()
	mail.Manager.GetMail(message.GetUid(),message.GetBeforeTime(),message.GetOffset(),message.GetCount())
}

type GetAllMailService struct {
}

func (service *GetAllMailService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetGetAllMail()
	mails := mail.Manager.GetAllMail(message.GetUid())
	response.GetAllMailRet = &protocol.GetAllMailRet{
		Mail:mails,
	}
}

type DrawMailService struct {
}

func (service *DrawMailService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetDrawMail()
	mailInfo := mail.Manager.DrawMail(message.GetUid(), message.GetMailID())
	mail.Manager.RemoveMail(message.GetUid(), message.GetMailID())
	response.DrawMailRet = &protocol.DrawMailRet{
		Mail:mailInfo,
	}
}

type RemoveMailService struct {
}

func (service *RemoveMailService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetRemoveMail()
	mail.Manager.RemoveMails(message.GetUid(), message.GetMailID())
	response.RemoveMailRet = &protocol.RemoveMailRet{
		Result:true,
	}
}

type CreateMailService struct {

}

func (service *CreateMailService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetCreateMail()
	response.CreateMailRet = &protocol.CreateMailRet{
		Mail:mail.Manager.CreateMail(message.GetUid(),message.GetTitle(),message.GetContent(),message.GetMailAttach()),
	}
}