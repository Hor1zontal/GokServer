package	http

import (
	"aliens/common/helper"
	"strconv"
	"encoding/json"
	"aliens/log"
	"net/http"
	"gok/module/passport/core"
	"gok/module/passport/cache"
	"encoding/base64"
	"aliens/common/cipher"
	"gok/module/passport/wx"
	"aliens/common/character"
	"aliens/common/util"
)

const(
	PARAM_UID = "uid"
	PARAM_CHANNEL = "channel"
	PARAM_OPENID = "openid"
	PARAM_TIME = "time"
	PARAM_APPID = "appid"
	PARAM_NICK = "nick"
	PARAM_MEMO = "memo"
	PARAM_AVATAR = "avatar"
	PARAM_SESSIONKEY = "session_key"
	PARAM_RAWDATA = "rawData"
	PARAM_SIGNATURE = "signature"
	PARAM_ENCRYPTED = "encryptedData"
	PARAM_IV = "iv"
	PARAM_UNIONID = "unionId"

	PARAM_VIVO_CODE = "code"
	PARAM_VIVO_REFRESH_TOKEN = "refresh_token"
)

type ResponseResult int

const (
	ResultSuccess                   ResponseResult = iota //0 成功
	ResultInvalidParam                             = 1001 //1001 无效的参数
	ResultInvalidGameServer                        = 1002 //1002 无可用的游戏服务器
	ResultInvalidSign                              = 1003 //1003 无效的签名
	ResultInvalidAuth                              = 1004 //1004 用户被封号
	ResultServerColse							   = 1005 //1005 服务器未开放
	ResultServerCloseNew						   = 1006 //1006 服务器停止新账号注册
	ResultServerMaintain						   = 1007 //1005 服务器维护中

)

func parametersIsset(request *http.Request) bool {
	return request.FormValue(PARAM_CHANNEL) != "" &&
		request.FormValue(PARAM_UID) != "" &&
		request.FormValue(PARAM_APPID) != ""
}

func CustomPush(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	uid := character.StringToInt32(request.FormValue("uid"))
	content := request.FormValue("content")
	wx.PushCustomMessage(uid, content)
	helper.SendToClient(responseWriter, GetErrorResponse(ResponseResult(0)))
}

func ChannelLogin(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	//log.Debug("login request: %v", request.Form)
	if !parametersIsset(request) {
		helper.SendToClient(responseWriter, GetErrorResponse(ResultInvalidParam))
		return
	}

	channel := request.FormValue(PARAM_CHANNEL)
	channelUID := request.FormValue(PARAM_UID)
	//appID := request.FormValue(PARAM_APPID)
	//sign := request.FormValue(PARAM_SIGN)
	openID := request.FormValue(PARAM_OPENID)
	//time := request.FormValue(PARAM_TIME)
	avatar := request.FormValue(PARAM_AVATAR)
	nick := request.FormValue(PARAM_NICK)

	//--------------------------------------------------------------------------------------

	//signText := "channel=" + channel + "&appid=" + appID + "&time=" + time + "&uid=" + channelUID + conf.Server.AppKey
	//signText := channel + appID + time + channelUID + conf.Server.AppKey
	succ := isSignSuccess(request.Form)
	if !succ {
		helper.SendToClient(responseWriter, GetErrorResponse(ResultInvalidSign))
		return
	}

	//result := rpc.PassportServiceProxy.ChannelLogin(channel, channelUID, openID, nick, avatar)
	result := core.ChannelLogin(channel, channelUID, openID, avatar, nick)

	if result.GetResult() == int32(RESULT_CODE_SUCCESS) {
		unionID := cache.UserCache.GetUnionIDByUserID(result.GetUid())
		if unionID == "" {
			sessionKey := request.FormValue(PARAM_SESSIONKEY)
			encryptedData := request.FormValue(PARAM_ENCRYPTED)
			ivData := request.FormValue(PARAM_IV)
			if sessionKey != "" && encryptedData != "" && ivData != "" {
				//rawData := request.FormValue(PARAM_RAWDATA)
				//signature := request.FormValue(PARAM_SIGNATURE)
				//signature2 := sha1.Sum(util.Str2Bytes(rawData + sessionKey))
				//bar := signature2[:]
				//if !bytes.Equal(util.Str2Bytes(signature),bar) {
				//	log.Error("error : signature: %v != signature2: %v", signature, util.Bytes2Str(bar))
				//}
				encrypted, _ := base64.StdEncoding.DecodeString(encryptedData)
				aeskey, _ := base64.StdEncoding.DecodeString(sessionKey)
				iv, _ := base64.StdEncoding.DecodeString(ivData)
				decryptedData, _ := cipher.CBCIvDecrypt(encrypted, aeskey, iv)
				log.Debug("decryptedData: %v", util.Bytes2Str(decryptedData))

				decryptedMapping := make(map[string]interface{})
				json.Unmarshal(decryptedData, &decryptedMapping)
				unionid, ok := decryptedMapping[PARAM_UNIONID].(string)
				if ok  {
					log.Debug("update uid-unionID %v-%v", result.GetUid(), unionid)
					cache.UserCache.SetUserIDUnionIDMapping(result.GetUid(), unionid)
				}
			}
		}
		helper.SendToClient(responseWriter, GetSuccessResponse(result.GetUid(), result.GetToken(), "", result.GetNew()))
	} else {
		helper.SendToClient(responseWriter, GetErrorResponse(ResponseResult(result.GetResult())))
	}
}

type LoginResponse struct {
	Code int   `json:"rcode"`
	Uid int32 `json:"uid"`
	Token string `json:"token"`
	GameServer string `json:"gameserver"`
	New bool `json:"new"`
}


func GetErrorResponse(code ResponseResult) string {
	return string("{\"rcode\":" + strconv.Itoa(int(code)) + "}")
}

func GetSuccessResponse(uid int32, token string, gameServer string, new bool) string {
	result, _ := json.Marshal(&LoginResponse{Code: int(ResultSuccess), Uid: uid, Token: token, GameServer: gameServer, New:new})
	return string(result)
}