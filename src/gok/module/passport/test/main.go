/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/10/9
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
)

type ErrorResponse struct {
	ErrCode float64 `json:"errcode"`
	ErrMsg  string  `json:"errmsg"`
}


type WechatInfo struct {
	FromUserName string `xml:"FromUserName,CDATA"`
	Event	string 		`xml:"Event,CDATA"`
	EventKey  string 	`xml:"EventKey,CDATA"`
}

func main() {
	//err := PushCustomMessageByOpenID("14_0xD7MKYIZOk5FzSD6wYGnpH26bxuUKntajOvAqFut2q7SqSsBcDHhrt3CgtyY5rTRFcQkS_pKMjMXsQxis7cGNiGE2zXJx9XVMLNNjFGXRIENC2_YVDRe77gyBo26zkhm4YfCT2GeWHH_ixXGYVfACAYQN", "ohNsowkLoC5sOnjs1_7yEQFR3E-s", "你好啊")
	//if err != nil {
	//	fmt.Println(err)
	//}
	a:=2

	switch a {
	case 1:
		fmt.Println(1)
	case 2:
		fmt.Println(2)
	case 3:
		fmt.Println(3)
	default:
		fmt.Println(4)
	}

	//content := "<xml><ToUserName><![CDATA[gh_8c25670512d1]]></ToUserName>\n<FromUserName><![CDATA[ohNsowkLoC5sOnjs1_7yEQFR3E-s]]></FromUserName>\n<CreateTime>1539135007</CreateTime>\n<MsgType><![CDATA[event]]></MsgType>\n<Event><![CDATA[CLICK]]></Event>\n<EventKey><![CDATA[ACTIVE_PUSH]]></EventKey>\n</xml>"
	//
	//wechatInfo := &WechatInfo{}
	//err1 := xml.Unmarshal([]byte(content), wechatInfo)
	//if err1 != nil {
	//	fmt.Errorf("%v", err1)
	//	//log.Error("parse wechat xml info error : %v", err1)
	//	return
	//}
	//fmt.Println(wechatInfo)
}


func GetPushContent(openID string, content string) string {
	return strings.Join([]string{"{\"touser\":\"", openID, "\",\"msgtype\":\"text\",\"text\":{\"content\":\"", content, "\"}"}, "")
}

func PushCustomMessageByOpenID(accessToken string, openID string, content string) error {
	//推送自定义消息
	requestLine := strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=", accessToken}, "")
	rawData := GetPushContent(openID, content)
	resp, err := http.Post(requestLine, "raw", strings.NewReader(rawData))
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	ater := &ErrorResponse{}
	err = json.Unmarshal(body, &ater)
	if err != nil {
		return err
	}
	if ater.ErrCode != 0 {
		return errors.New(fmt.Sprintf("push message error %v", ater.ErrMsg))
	}
	return nil
}
