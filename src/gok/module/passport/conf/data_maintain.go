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
	"aliens/log"
	"time"
	"gok/service/msg/protocol"
	"gok/constant"
	"gok/service/rpc"
	"sync"
)

var Maintain = MaintainBase{} //维护信息

var lock sync.RWMutex

var MaintainNotify func() = nil

type MaintainBase struct {
	MaintainState int32 `json:"maintain_state"`     //维护状态    0关闭维护  1开启维护

	ServerState        int32     `json:"state"`               //服务器状态  0开放 1停新 2停服
	OpenTime           string    `json:"open_time"`           //服务器的开放时间
	OpenTimestamp      time.Time `json:"-"`                   //服务器的开放时间戳
	IsCheckVersion     bool      `json:"is_check_version"`    // 是否校验客户端版本号
}

//是否在维护期间
func IsMaintain() bool {
	lock.RLock()
	defer lock.RUnlock()
	return Maintain.MaintainState == constant.SERVER_MAINTAIN_STATE_OPEN
}

func IsServerOpen(time time.Time) bool {
	lock.RLock()
	defer lock.RUnlock()
	if Maintain.ServerState == constant.SERVER_STATE_CLOSE {
		return false
	}
	//没有配置开服时间
	if Maintain.OpenTimestamp.Unix() <= 0 {
		return true
	}
	//是否过了开服时间
	return time.After(Maintain.OpenTimestamp)
}

func GetServerState() int32 {
	lock.RLock()
	defer lock.RUnlock()
	return Maintain.ServerState
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
	//maintain.OpenTime = maintain.OpenTimestamp.Format("1月2日15:04")
	//log.Debug("open time : %v", maintain.OpenTime)

	lock.Lock()
	Maintain = *maintain
	lock.Unlock()

	var kickType = constant.LOGOUT_TYPE_NONE
	if !IsServerOpen(time.Now()) {
		kickType = constant.LOGOUT_TYPE_SERVER_CLOSE
	} else if IsMaintain() {
		kickType = constant.LOGOUT_TYPE_SERVER_MAINTAIN
	}

	if kickType == constant.LOGOUT_TYPE_NONE {
		return
	}
	//广播所有服务器停服
	request := &protocol.C2GS{
		Sequence: []int32{500},
		Kickoff: &protocol.KickOff{
			Uid:      -1,
			KickType: int32(kickType),
		},
	}
	rpc.UserServiceProxy.AsyncBroadcastAllMessage(request)
}