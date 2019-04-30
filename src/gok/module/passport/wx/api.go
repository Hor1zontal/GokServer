/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/10/9
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package wx

import (
	clustercache "gok/module/cluster/cache"
	"gok/module/passport/cache"
	"net/http"
	"io/ioutil"
	"strings"
	"errors"
	"fmt"
	"encoding/json"
	"aliens/log"
	"gok/module/passport/wx/model"
	"bytes"
	"gok/service/rpc"
	"gok/service/msg/protocol"
	"gok/constant"
	"aliens/common/util"
	"time"
)


func RefreshToken() {
	//TODO 只需要一个服务器刷新即可
	//只有主节点才需要刷新token
	result, err := FetchAccessToken()
	if err != nil {
		log.Error("refresh accessToken error : %v", err)
		return
	}
	log.Info("refresh accessToken success : %v", result)
	clustercache.Cluster.SetAccessToken(result.AccessToken, int(result.ExpiresIn))
}

//获取wx_AccessToken 拼接get请求 解析返回json结果 返回 AccessToken和err
func FetchAccessToken() (*model.AccessTokenResponse, error) {
	if wxConfig.AppID == "" {
		return nil, errors.New("accessToken param is not configuration")
	}
	//appID string, appSecret string, accessTokenFetchUrl string
	requestLine := strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=",
		wxConfig.AppID,
		"&secret=",
		wxConfig.AppSecret}, "")

	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if bytes.Contains(body, []byte("access_token")) {
		atr := &model.AccessTokenResponse{}
		err = json.Unmarshal(body, &atr)
		if err != nil {
			return nil, err
		}
		return atr, nil
	} else {
		ater := &model.ErrorResponse{}
		err = json.Unmarshal(body, &ater)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", ater.ErrMsg)
	}
}

func BuildTextResponse(fromUserName, toUserName, content, nonce, timestamp string, encrypt bool) []byte {
	if encrypt {
		respBody, e := MakeEncryptResponseBody(fromUserName, toUserName, content, nonce, timestamp)
		if e != nil {
			log.Error("make encrypt response err : %v", e)
		}
		return respBody
	} else {
		respBody, e := MakeTextResponseBody(fromUserName, toUserName, content)
		if e != nil {
			log.Error("make text response err : %v", e)
		}
		return respBody
	}
}

func RemoveUnionID(openID string) {
	if openID == "" {
		return
	}
	log.Debug("remove wx service oumapping %v", openID)
	unionID := cache.UserCache.GetWxServiceOUMapping(openID)

	if unionID != "" {
		uid := cache.UserCache.GetUserIDByUnionID(unionID)
		if uid > 0 {
			rpc.UserServiceProxy.PersistCall(uid, &protocol.C2GS{
				Sequence:[]int32{160},
				ActivePrivilege:false,
			})
		}
	}

	cache.UserCache.CleanWxServiceMapping(openID)

}

func HandleClickEvent(eventID string, openID string) string {
	//激活游戏特权
	if eventID == constant.EVENT_ACTIVE_PRIVILEGE {
		_, err := UpdateUnionID(openID)
		if err != nil {
			log.Error("%v", err)
		}
	} else if eventID == constant.EVENT_DAY_GIFT {
		uid := GetUIDByWechatServiceOpenID(openID)
		if uid > 0 {
			SetDayGiftCache(uid)
		}
	}
	return eventID
}

func GetUIDByWechatServiceOpenID (openID string) int32{
	unionID := cache.UserCache.GetWxServiceOUMapping(openID)
	if unionID == "" {
		return 0
	}
	return cache.UserCache.GetUserIDByUnionID(unionID)
}

func GetWechatServiceOpenIDByUid (uid int32) string {
	unionID := cache.UserCache.GetUnionIDByUserID(uid)
	if unionID == "" {
		return ""
	}
	//通过unionID获取服务号的openid
	return cache.UserCache.GetWxServiceUOMapping(unionID)
}

//通过openid获取unionid
func UpdateUnionID(openID string) (string, error) {
	accessToken := clustercache.Cluster.GetAccessToken()
	if accessToken == ""  {
		return "", errors.New("accessToken can not be found")
	}
	requestLine := strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/user/info?lang=zh_CN&access_token=", accessToken, "&openid=", openID}, "")
	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Debug("get %v unionID : %v",openID, string(body))
	decryptedMapping := make(map[string]interface{})
	json.Unmarshal(body, &decryptedMapping)
	unionID, _ := decryptedMapping["unionid"].(string)
	if unionID != "" {
		log.Debug("update wx service oumapping %v - %v", openID, unionID)
		cache.UserCache.SetWxServiceMapping(openID, unionID)
		uid := cache.UserCache.GetUserIDByUnionID(unionID)
		if uid > 0 {
			rpc.UserServiceProxy.PersistCall(uid, &protocol.C2GS{
				Sequence:[]int32{160},
				ActivePrivilege:true,
			})
		}
	}
	return unionID, nil
}

func GetPushContent(openID string, content string) string {
	return strings.Join([]string{"{\"touser\":\"", openID, "\",\"msgtype\":\"text\",\"text\":{\"content\":\"", content, "\"}"}, "")
}

//推送客服消息
func PushCustomMessage(uid int32, content string) error {
	//unionID := cache.UserCache.GetUnionIDByUserID(uid)
	//if unionID == "" {
	//	return errors.New(fmt.Sprintf("user %v can not found unionID", uid))
	//}
	////通过unionID获取服务号的openid
	//serviceOpenID := cache.UserCache.GetWxServiceUOMapping(unionID)
	////没有关注服务号
	//if serviceOpenID == "" {
	//	return errors.New(fmt.Sprintf("user %v is not subscribe", uid))
	//}
	serviceOpenID := GetWechatServiceOpenIDByUid(uid)
	if serviceOpenID == "" {
		return errors.New(fmt.Sprintf("user %v is not subscribe", uid))
	}
	return PushCustomMessageByOpenID(serviceOpenID, content)
}

func PushCustomMessageByOpenID(openID string, content string) error {
	accessToken := clustercache.Cluster.GetAccessToken()
	if accessToken == ""  {
		return errors.New("accessToken can not be found")
	}
	//推送自定义消息
	requestLine := strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=", accessToken}, "")
	rawData := GetPushContent(openID, content)
	//log.Debug("rawData %v", rawData)
	resp, err := http.Post(requestLine, "raw", strings.NewReader(rawData))
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	ater := &model.ErrorResponse{}
	err = json.Unmarshal(body, &ater)
	if err != nil {
		return err
	}
	if ater.ErrCode != 0 {
		return errors.New(fmt.Sprintf("push message error %v", ater.ErrMsg))
	}
	return nil
}

func SetDayGiftCache (uid int32) {
	if !cache.UserCache.IsUserExist(uid) {
		return
	}
	if !cache.UserCache.ExistUserDayGiftCache(uid) {
		cache.UserCache.SetUserDayGiftCache(uid, false)
		//log.Info("set uid:%v cache false",uid)
		expireTime := util.GetTodayHourTime(24).Sub(time.Now()).Seconds()
		cache.UserCache.SetUserDayGiftExpireTime(uid, int(expireTime))
		//如果在线直接领取
		message := &protocol.C2GS{
			Sequence: []int32{560},
			DrawDayGift: &protocol.DrawDayGift{
				Uid:uid,
			},
		}
		rpc.UserServiceProxy.UserHandleMessage(uid, message)
	}
}
