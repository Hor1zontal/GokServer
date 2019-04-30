/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/4/18
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package http

import (
	"net/http"
	"aliens/common/character"
	"gok/service/rpc"
	"gok/service/msg/protocol"
	"gok/module/center/cache"
	"aliens/common/helper"

)


const(
	PARAM_UID = "uid"
	PARAM_TITLE = "title"
	PARAM_CONTENT = "content"
	PARAM_ATTACH = "attach"
)

func Mail(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	code, mailID := MailProxy(request)
	helper.SendToClient(responseWriter, helper.GetDataResponse(code, mailID))
}

func MailProxy(request *http.Request) (int, int64) {
	//memo := request.FormValue(PARAM_MEMO)
	uidStr := request.FormValue(PARAM_UID)
	title := request.FormValue(PARAM_TITLE)
	content := request.FormValue(PARAM_CONTENT)
	attach := request.FormValue(PARAM_ATTACH)
	//sign := request.FormValue(PARAM_SIGN)

	uid := character.StringToInt32(uidStr)
	if uid == 0 {
		return RESULT_CODE_INVALID_PARAM, 0
	}

	//signText := title + content + uidStr + attach + conf.Server.AppKey
	succ := isSignSuccess(request.Form)
	if !succ {
		//测试跳过签名
		return RESULT_CODE_INVALID_SIGN, 0
	}

	if !cache.UserCache.IsUserExist(uid) {
		return RESULT_CODE_USER_NOTFOUND, 0
	}

	//rpc.StarServiceProxy.
	mail := rpc.MailServiceProxy.CreateMail(uid, title, content, attach).GetMail()
	if mail == nil {
		return RESULT_CODE_DATABASE_EXCEPTION, 0
	}
	rpc.UserServiceProxy.Push(uid, &protocol.GS2C{
		Sequence:[]int32{1059},
		MailPush: mail,
	})

	return RESULT_CODE_SUCCESS, mail.Id
}
