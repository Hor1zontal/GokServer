package http

import (
	"gok/module/passport/cache"
	"net/http"
	"gok/service/rpc"
	"aliens/common/character"
	"aliens/common/helper"
	"gok/constant"
)
const (
	NOTICE_TITLE = "title"
	NOTICE_CONTENT = "content"
	NOTICE_START_TIME = "start"
	NOTICE_END_TIME = "end"
	NOTICE_STATUS = "status"
)

func Notice(responseWriter http.ResponseWriter, request *http.Request){
	request.ParseForm()
	code, noticeData := NoticeProxy(request)
	helper.SendToClient(responseWriter, helper.GetDataResponse(code, noticeData))
}

func Flag(responseWriter http.ResponseWriter, request *http.Request){
	request.ParseForm()

	opt := request.FormValue("opt")
	key := request.FormValue("key")
	value := character.StringToInt32(request.FormValue("value"))
	if opt == "set" {
		cache.UserCache.SetCustomFlag(key, int(value))
	} else if opt == "get" {
		value = int32(cache.UserCache.GetCustomFlag(key))
	}

	helper.SendToClient(responseWriter, character.Int32ToString(value))
}

func NoticeProxy(request *http.Request) (int, interface{}){
	request.ParseForm()
	succ := isSignSuccess(request.Form)
	if !succ {
		return RESULT_CODE_INVALID_SIGN, 0
	}
	if request.FormValue(NOTICE_STATUS) == "" {
		return RESULT_CODE_INVALID_PARAM, 0
	}
	noticeStatus := character.StringToInt32(request.FormValue(NOTICE_STATUS))
	if noticeStatus != constant.NOTICE_ON {
		return RESULT_CODE_INVALID_PARAM, 0
	}
	resp := rpc.CenterServiceProxy.GetOnNotices(noticeStatus)
	if resp == nil {
		return RESULT_CODE_SERVER_EXCEPTION, 0
	}
	return RESULT_CODE_SUCCESS, resp.Notices
}
