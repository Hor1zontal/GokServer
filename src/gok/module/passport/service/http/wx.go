package http

import (
	"aliens/common/character"
	"aliens/common/helper"
	"aliens/log"
	"fmt"
	"gok/constant"
	"gok/module/passport/conf"
	"gok/module/passport/notify"
	"gok/module/passport/wx"
	"gok/module/passport/wx/model"
	statismodel "gok/module/statistics/model"
	"gok/service/lpc"
	"net/http"
	"strings"
)

//type WechatInfo struct {
//	FromUserName string `xml:"FromUserName,CDATA"`
//	Event	string 		`xml:"Event,CDATA"`
//	EventKey  string 	`xml:"EventKey,CDATA"`
//}
//
//type CDATAText struct {
//	Text string `xml:",innerxml"`
//}
//
//
//type TextResponseBody struct {
//	XMLName      xml.Name `xml:"xml"`
//	ToUserName   CDATAText
//	FromUserName CDATAText
//	CreateTime   string
//	MsgType      CDATAText
//	Content      CDATAText
//}

//刷新token
func Refresh(responseWriter http.ResponseWriter, request *http.Request){
	wx.RefreshToken()
	helper.SendToClient(responseWriter, GetErrorResponse(ResponseResult(0)))
}



func Wechat(responseWriter http.ResponseWriter, request *http.Request){
	request.ParseForm()
	method := strings.Join(request.Form["method"], "")
	uid := strings.Join(request.Form["uid"], "")
	openid := strings.Join(request.Form["openid"], "")
	result := ""
	if method == "getUO" {
		result = wx.GetWechatServiceOpenIDByUid(character.StringToInt32(uid))
	} else if method == "getOU" {
		result = character.Int32ToString(wx.GetUIDByWechatServiceOpenID(openid))
	}
	fmt.Fprintf(responseWriter, result)
}

func Wx(responseWriter http.ResponseWriter, request *http.Request){
	request.ParseForm()
	timestamp := strings.Join(request.Form["timestamp"], "")
	nonce := strings.Join(request.Form["nonce"], "")
	signature := strings.Join(request.Form["signature"], "")
	encryptType := strings.Join(request.Form["encrypt_type"], "")
	msgSignature := strings.Join(request.Form["msg_signature"], "")
	echoStr := strings.Join(request.Form["echostr"], "")


	//content := string(body)
	//log.Debug("handle wx request %v, body %v", request.Form, content)
	if !wx.ValidateUrl(timestamp, nonce, signature){
		log.Info("Wechat Service: this http request is not from Wechat platform!")
		return
	}

	//第一次验证url需要直接返回echoStr
	if echoStr != "" {
		fmt.Fprintf(responseWriter, echoStr)
	}

	var requestBody *model.TextRequestBody

	isEncrypt := encryptType == "aes"

	if isEncrypt {
		encryptRequestBody := wx.ParseEncryptRequestBody(request)
		// Validate mstBody signature
		if !wx.ValidateMsg(timestamp, nonce, encryptRequestBody.Encrypt, msgSignature) {
			log.Error("Wechat Service: validate msg error")
			return
		}
		plainData, err := wx.AesDecrypt(encryptRequestBody.Encrypt)
		if err != nil {
			log.Error("Wechat Service: descrypt err : %v", err)
			return
		}
		//封装struct
		requestBody = wx.ParseEncryptTextRequestBody(plainData)
	} else {
		requestBody = wx.ParseTextRequestBody(request)
	}

	log.Debug("request body : %v", requestBody)

	if requestBody == nil {
		fmt.Fprintf(responseWriter, "success")
		return
	}

	//if requestBody.MsgType == "text" && requestBody.Content == "【收到不支持的消息类型，暂无法显示】" {
	//	//安全模式下向用户回复消息也需要加密
	//	respBody, e := wx.MakeEncryptResponseBody(requestBody.ToUserName, requestBody.FromUserName, "一些回复给用户的消息", nonce, timestamp)
	//	if e != nil {
	//		log.Error("encrypt response error : %v", e)
	//		return
	//	}
	//	fmt.Fprintf(responseWriter, string(respBody))
	//
	//}

	openID := requestBody.FromUserName
	if requestBody.MsgType == constant.KEYWORD_EVENT {

		static := &statismodel.StatisticWechat{
			OpenID: openID,
			Event: requestBody.Event,
		}

		switch requestBody.Event {
		case constant.EVENT_SUBSCRIBE:
			//unionID, err := wx.UpdateUnionID(openID)
			//if err != nil {
			//	log.Error("%v", err)
			//}
			//关注公众号更新过期时间

			updateExpire(openID)
			//uid := UpdateExpire(unionID)
			//static.Uid = uid
			response := wx.BuildTextResponse(requestBody.ToUserName, openID, wx.GetTextResponse(constant.EVENT_SUBSCRIBE), nonce, timestamp, isEncrypt)
			if response != nil {
				responseWriter.Write(response)
			}
		case constant.EVENT_UNSUBSCRIBE:
			wx.RemoveUnionID(openID)
		case constant.EVENT_CLICK:
			static.Event = requestBody.EventKey
			key := wx.HandleClickEvent(requestBody.EventKey, openID)
			//激活推送更新过期时间
			updateExpire(openID)
			//if requestBody.EventKey == constant.EVENT_ACTIVE_PUSH || requestBody.EventKey == constant.EVENT_DAY_GIFT {
			//	//unionID := cache.UserCache.GetWxServiceOUMapping(openID)
			//	//notify.PushEventMsgByOpenID(constant.EVENT_PUSH_EXPIRE, openID, int((40 * time.Hour).Seconds()))
			//	//uid := UpdateExpire(unionID)
			//	//static.Uid = uid
			//}
			response := wx.BuildTextResponse(requestBody.ToUserName, openID, wx.GetTextResponse(key), nonce, timestamp, isEncrypt)
			if response != nil {
				responseWriter.Write(response)
			} else {
				fmt.Fprintf(responseWriter, "success")
			}
		default:

			//TODO
			fmt.Fprintf(responseWriter, "success")
		}

		lpc.StatisticsHandler.AddStatisticData(static)
		//某个类型的消息暂时后台不作处理，也需要向微信服务器做出响应
		//return service.NewSimpleError(service.SERVER_WRITE_TEXT, "success")
	}
}

func updateExpire(openID string) {
	notify.PushEventMsgByOpenID(constant.EVENT_PUSH_ONEDAY_EXPIRE, openID, conf.DATA.WechatExpireTime1)
	notify.PushEventMsgByOpenID(constant.EVENT_PUSH_EXPIRE, openID, conf.DATA.WechatExpireTime2)
}



//func UpdateExpire(unionID string) int32 {
//	//unionID := cache.UserCache.GetWxServiceOUMapping(openID)
//	if unionID != "" {
//		uid := cache.UserCache.GetUserIDByUnionID(unionID)
//		if uid > 0 {
//			notify.PushEventMsg(constant.EVENT_PUSH_EXPIRE, uid, int((40 * time.Hour).Seconds()))
//		}
//		return uid
//	}
//	return 0
//}


//func DayGift(responseWriter http.ResponseWriter, request *http.Request){
//	uid := character.StringToInt32(request.FormValue("uid"))
//	log.Info("uid:%v",uid)
//	wx.SetDayGiftCache(uid)
//	helper.SendToClient(responseWriter, helper.GetResponse(0))
//}
