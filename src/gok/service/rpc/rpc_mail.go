/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package rpc

import (
	"gok/service/msg/protocol"
	"gok/service"
)

var MailServiceProxy = &MailService{&rpcHandler{serviceType:service.SERVICE_MAIL_RPC}}

type MailService struct {
	*rpcHandler
}

func (this *MailService) CreateMail(uid int32, title string, content string, attach string) *protocol.CreateMailRet{
	request := &protocol.C2GS{
		Sequence: []int32{73},
		CreateMail: &protocol.CreateMail{
			Uid: uid,
			Title: title,
			Content: content,
			MailAttach: attach,
		},
	}
	return this.HandleMessage(request).GetCreateMailRet()
}
