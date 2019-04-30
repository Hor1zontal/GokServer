/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/11/2
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package notify

import (
	"gok/module/passport/cache"
	"aliens/common/character"
	"strings"
	"gok/module/passport/conf"
	"gok/module/passport/wx"
	"aliens/log"
)

//var EventCallbacks = make(map[string]func(event string, uid int32))

func Init() {
	//
	cache.UserCache.SubscribeExpire(handleExpire)
}


func handleExpire(pattern, channel, value string) {
	results := strings.Split(value, ":")
	if results == nil || len(results) < 2 {
		return
	}
	event := results[0]
	uid := character.StringToInt32(results[1])
	if uid > 0 {
		PushEventMsg(event, uid, 0)
	} else if len(results[1]) > 0 {
		//直接通过openID推送
		PushEventMsgByOpenID(event, results[1], 0)
	}
}


func PushEventMsg(event string, uid int32, delay int) {
	content := conf.DATA.WechatResponseData[event]

	log.Debug("push wechat event %v - %v - %v delay : %v", event, uid, content != "", delay)
	if content == "" {
		log.Debug("unexpect wechat event %v", event)
		return
	}
	if delay > 0 {
		cache.UserCache.SetExpireData(event + ":" + character.Int32ToString(uid), true, delay)
	} else {
		err := wx.PushCustomMessage(uid, content)
		if err != nil {
			log.Error("send wechat message err : %v", err)
		}
	}

	//helper.SendToClient(responseWriter, GetErrorResponse(ResponseResult(0)))
}

func PushEventMsgByOpenID(event string, openID string, delay int) {
	content := conf.DATA.WechatResponseData[event]
	if content == "" {
		log.Debug("unexpect wechat event %v", event)
		return
	}
	log.Info("delay:%v", delay)
	log.Debug("wechat push %v-%v", openID, event)
	if delay > 0 {
		cache.UserCache.SetExpireData(event + ":" + openID, true, delay)
	} else {
		err := wx.PushCustomMessageByOpenID(openID, content)
		if err != nil {
			log.Error("send wechat message err : %v", err)
		}
	}

	//helper.SendToClient(responseWriter, GetErrorResponse(ResponseResult(0)))
}

