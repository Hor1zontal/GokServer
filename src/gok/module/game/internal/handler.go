package internal

import (
	"gok/service/msg/protocol"
	"github.com/name5566/leaf/gate"
	"gok/module/game/user"
	"reflect"
	"gok/constant"
	"github.com/name5566/leaf/timer"
	"aliens/log"
)

func init() {
	//user.Manager.Server = skeleton.ChanRPCServer
	// 向当前模块注册客户端发送的消息处理函数 handleMessage
	cron, err := timer.NewCronExpr("*/1 * * * *")
	if err != nil {
		log.Error("init user timer error : %v", err)
	}

	//每个月一号凌晨五点清理一次排行榜
	//rankCron, err := timer.NewCronExpr("0 5 1 * *")
	//if err != nil {
	//	log.Error("init rank timer error : %v", err)
	//}


	/*------------事件------------------*/
	//eventCron, err := timer.NewCronExpr("*/1 * * * *")
	//if err != nil {
	//	log.Error("init event timer error : %v", err)
	//}

	skeleton.RegisterChanRPC(reflect.TypeOf(&protocol.C2GS{}), handleMessage)
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)

	skeleton.CronFunc(cron, user.Manager.DealUserTimer)
	//skeleton.CronFunc(rankCron, rank.CleanRank)

	//skeleton.CronFunc(cron, user.M.DealEvent)
	//skeleton.RegisterChanRPC("Release", rpcReleaseSession)
}

////处理RPC消息
//func handleRPCMessage(args []interface{}) {
//	message := args[0].(*protocol.C2GS)
//	server := args[1].(protocol.RPCService_RequestServer)
//	service := args[2].(*service.LocalService)
//	response := service.HandleMessage(message)
//	server.Send(response)
//}

//处理客户端网关消息
func handleMessage(args []interface{}) {
	// 消息的发送者
	request := args[0].(*protocol.C2GS)
	session, ok := args[1].(gate.Agent).UserData().(*user.Session)
	if ok {
		session.SyncMessage(request)
	}
}

func rpcNewAgent(args []interface{}) {
	agent := args[0].(gate.Agent)
	if agent.UserData() == nil {
		session := user.Manager.NewSession(agent)
		agent.SetUserData(session)
	}
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	session, ok := a.UserData().(*user.Session)
	if ok {
		session.Logout(constant.LOGOUT_TYPE_NONE)
		//session.SyncCommand(command.KICK, )
		a.SetUserData(nil)
	}
}

//func rpcReleaseSession(args []interface{}) {
//	session := args[0].(*user.Session)
//	user.Manager.LogoutSession(session)
//}




