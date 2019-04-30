package user

import (

)

//处理用户系统消息

//var systemHandlers = make(map[command.Command]func(user *Session, param ...interface{}))
//
//func init() {
//	//RegisterSystemHandler(command.LOGIN, innerLoginService)
//	//RegisterSystemHandler(command.AGENT_CLOSE, innerAgentCloseService)
//	//RegisterSystemHandler(command.RELEASE, innerReleaseService)
//
//	//T人
//	//RegisterSystemHandler(command.KICK, innerKickService)
//	RegisterSystemHandler(command.TIMER, innerTimerService)
//	//RegisterSystemHandler(command.EVENT, innerEventService)
//	//RegisterSystemHandler(command.LOGINPUSH, innerLoginPushService)
//	//RegisterSystemHandler(command.MESSAGEPUSH, messagePushService)
//
//
//}
//
//func RegisterSystemHandler(cmd command.Command, handler func(user *Session, param ...interface{})) {
//	systemHandlers[cmd] = handler
//}
//
//func HandleSystemMessage(cmd command.Command, userContext *Session, param ...interface{}) {
//	handler := systemHandlers[cmd]
//	if handler != nil {
//		handler(userContext, param...)
//	}
//}

//func innerLoginPushService(user *Session, param ...interface{}) {
//	//推送所有玩家上线
//
//}
//
//func messagePushService(user *Session, param ...interface{}) {
//	if len(param) == 1 {
//		message, ok := param[0].(*protocol.GS2C)
//		if !ok {
//			return
//		}
//		if message.GetGlobalMessagePush() != nil && !user.IsGlobalMessagePush() {
//			return
//		}
//		user.WriteMsg(message)
//	}
//}

//func innerLoginService(user *Session, param ...interface{}) {
//	user.LoginSession()
//}

//func innerKickService(user *Session, param ...interface{}) {
//	logoutType := param[0].(constant.LOGOUT_TYPE)
//	user.Logout(logoutType)
//}

//func innerEventService(user *Session, param ...interface{}) {
//	user.DealEvent()
//}

//func innerAgentCloseService(user *Session, param ...interface{}) {
//	user.logout(constant.LOGOUT_TYPE_NONE)
//	//没有登录过, 直接释放
//	if !user.IsAuth() {
//		user.Release(false)
//	}
//}

//释放用户上下文
//func innerReleaseService(user *Session, param ...interface{}) {
//	user.Release(true)
//}

//处理用户超时管理
//func innerTimerService(user *Session, param ...interface{}) {
//	time := time.Now()
//	if user.IsOnline() {
//		if user.IsFreeTimeout(time, conf.Server.UserFreeTimeout) {
//			user.Logout(constant.LOGOUT_TYPE_TIMEOUT)
//		}
//		if user.IsDirty() {
//			if user.IsSyncDBTimeout(time, conf.Server.SyncDBInterval) {
//				user.UpdateData()
//			}
//		}
//	}
//}



