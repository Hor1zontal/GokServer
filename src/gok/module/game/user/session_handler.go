/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package user

import (
	"gok/service/msg/protocol"
	"time"
	"gok/constant"
	"gok/service"
)

//func (this *Session) SyncCommand(cmd command.Command, param ...interface{}) {
//	HandleSystemMessage(cmd, this, param...)
//}

func (this *Session) SyncMessage(arg interface{}) {
	message, ok := arg.(*protocol.C2GS)
	if ok {
		if len(message.GetSequence()) == 0 {
			return
		}
		//更新消息活跃时间
		if message.GetSequence()[0] != constant.HEARBEAT_SEQ {
			this.lastActiveTime = time.Now()
		}
		service.ServiceManager.HandleChannelMessage(message, this)
	}
}

func (this *Session) AppendUserEventStatistic(arg interface{}, arg1 interface{}) {
	this.AppendEventStatistic(arg.(int32), arg1.(int32))
}