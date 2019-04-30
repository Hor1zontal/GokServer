/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/6/27
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package conf

import (
	"encoding/json"
	"aliens/common/util"
	"aliens/log"
	"time"
	"sync"
	"gok/module/cluster/center"
)

var Maintain = MaintainBase{} //维护信息

var lock sync.RWMutex

type MaintainBase struct {
	MaintainState int32 `json:"maintain_state"`     //维护状态    0关闭维护  1开启维护

	ServerState        int32     `json:"state"`               //服务器状态  0开放 1停新 2停服
	OpenTime           string    `json:"open_time"`           //服务器的开放时间
	OpenTimestamp      time.Time `json:"-"`                   //服务器的开放时间戳
	IsCheckVersion     bool      `json:"is_check_version"`    //是否校验客户端版本号

}

func ChangeState(serverState int32, openTimestamp int64, maintainState int32, isCheckVersion bool) bool {
	var time = util.GetTime(openTimestamp)
	var timeStr = time.Format("2006-01-02 15:04:05")
	//request := &MaintainBase{OpenTime:timeStr, MaintainState: maintainState, ServerState: serverState}
	lock.Lock()
	Maintain.ServerState = serverState
	Maintain.OpenTime = timeStr
	Maintain.MaintainState = maintainState
	Maintain.IsCheckVersion = isCheckVersion
	lock.Unlock()
	//Maintain.OpenTimestamp = time
	data, _ := json.Marshal(Maintain)
	result := center.ConfigCenter.PublicConfig("maintain", data)
	if result {
		log.Debug("update maintain success : %v", Maintain)
		return true
	} else {
		log.Debug("update maintain failed : %v", Maintain)
		return false
	}
}

func updateMaintainBase( data []byte)  {
	maintain := &MaintainBase{}
	error := json.Unmarshal(data, maintain)
	if error != nil {
		return
	}
	endTimestamp, err := time.ParseInLocation("2006-01-02 15:04:05", maintain.OpenTime, time.Local)
	if err != nil {
		if maintain.OpenTime != "" {
			log.Debug("invalid server open : %v", err)
		}
	} else {
		maintain.OpenTimestamp = endTimestamp
	}
	lock.Lock()
	Maintain = *maintain
	lock.Unlock()
}