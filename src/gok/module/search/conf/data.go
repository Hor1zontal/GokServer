/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/11
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package conf

import (
	"encoding/json"
	"gok/module/cluster/center"
	"time"
)

func Init() {
	center.ConfigCenter.SubscribeConfig("randomtarget", updateRandomTargetBase)
	center.ConfigCenter.SubscribeConfig("testtool", updateRandomTest)
	center.ConfigCenter.SubscribeConfig("gamebase", updateGameBaseData)
}

func Close() {

}

var Base struct {
	*GameBase
	RandomTarget  map[int32]map[int32]*RandomTargetBase //事件类型 - 等级范围 - 配置
	RandomTest map[int32]int32
}

type GameBase struct {
	StarWeCanArrive        []int32 //新账号的随机星球列表
	ReceiveNewsInterval    float64

	SearchRandomCD         float64  //玩家被随机搜索到的CD时间 默认一小时
	SearchActiveTime       float64  //玩家的被搜索到的活跃时间 默认5天
}

type RandomTest struct {
	EventID int32
	UID int32
}

type RandomTargetBase struct {
	Level      int32   `json:"level"`      //等阶号
	RandomType []int32 `json:"randomType"` //1抢信仰任务 2抢信徒任务 3拆建筑任务
	Type       int32   `json:"type"`       //1信徒总等级  2建筑总等级
	Min        int32   `json:"min"`        //下限
	Max        int32   `json:"max"`        //上限
}

func updateGameBaseData(data []byte) {
	var datas = &GameBase{}
	json.Unmarshal(data, datas)
	Base.GameBase = datas

	if Base.SearchRandomCD == 0 {
		Base.SearchRandomCD = time.Hour.Seconds()
	}

	if Base.SearchActiveTime == 0 {
		Base.SearchActiveTime = (5 * 24 * time.Hour).Seconds()
	}
}

func updateRandomTargetBase(data []byte) {
	var datas []*RandomTargetBase
	json.Unmarshal(data, &datas)

	results := make(map[int32]map[int32]*RandomTargetBase)

	for _, randomTarget := range datas {
		for _, eventType := range randomTarget.RandomType {
			eventRandom := results[eventType]
			if eventRandom == nil {
				eventRandom = make(map[int32]*RandomTargetBase)
				results[eventType] = eventRandom
			}
			eventRandom[randomTarget.Level] = randomTarget
		}
	}
	Base.RandomTarget = results
}


func updateRandomTest(data []byte) {
	var datas []*RandomTest
	json.Unmarshal(data, &datas)

	results := make(map[int32]int32)

	for _, randomTest := range datas {
		results[randomTest.EventID] = randomTest.UID
	}
	Base.RandomTest = results
}
