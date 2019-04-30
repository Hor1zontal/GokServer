package http

import (
	"aliens/common/helper"
	"aliens/log"
	"gok/module/passport/cache"
	"gok/module/passport/core"
	"gok/module/passport/vivo"
	"net/http"
)

func VivoLogin (responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	log.Debug("VivoLogin request: %v", request.Form)
	//channel := request.FormValue(PARAM_CHANNEL)
	//appID := request.FormValue(PARAM_APPID)
	//sign := request.FormValue(PARAM_SIGN)
	openID := request.FormValue(PARAM_OPENID)
	channelUID := openID
	//time := request.FormValue(PARAM_TIME)
	avatar := request.FormValue(PARAM_AVATAR)
	nick := request.FormValue(PARAM_NICK)
	unionID := request.FormValue(PARAM_UNIONID)
	result := core.ChannelLogin("", channelUID, openID, avatar, nick)

	if result.GetResult() == int32(RESULT_CODE_SUCCESS) {
		cacheUnionID := cache.UserCache.GetUnionIDByUserID(result.GetUid())
		if cacheUnionID == "" && unionID != ""{
			cache.UserCache.SetUserIDUnionIDMapping(result.GetUid(), unionID)
		}
		helper.SendToClient(responseWriter, GetSuccessResponse(result.GetUid(), result.GetToken(), "", result.GetNew()))
	} else {
		helper.SendToClient(responseWriter, GetErrorResponse(ResponseResult(result.GetResult())))
	}
}

func VivoGetToken(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	log.Debug("VivoGetToken request: %v", request.Form)
	code := request.FormValue(PARAM_VIVO_CODE)
	result, err := vivo.GetVivoAcessToken(code)
	if err == nil {
		helper.SendToClient(responseWriter, result)
	} else {
		helper.SendToClient(responseWriter, GetErrorResponse(ResultInvalidParam))
	}
}

func VivoRefresh(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	log.Debug("VivoRefresh request: %v", request.Form)
	refreshToken := request.FormValue(PARAM_VIVO_REFRESH_TOKEN)
	result, err := vivo.RefreshToken(refreshToken)
	if err == nil {
		helper.SendToClient(responseWriter, result)
	} else {
		helper.SendToClient(responseWriter, GetErrorResponse(ResultInvalidParam))
	}
}
