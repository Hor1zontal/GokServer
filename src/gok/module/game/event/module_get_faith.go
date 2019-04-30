/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/7/12
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package event

import (
	"gok/module/game/conf"
	"encoding/json"
	"gok/module/game/db"
	"gok/service/msg/protocol"
	"aliens/common/util"
	"aliens/common/character"
)

var condition1 = make(map[int32]int32)
var condition2 = make(map[int32]int32)
var condition3 = make(map[int32]int32)

func init() {
	condition1[1] = 6
	condition1[2] = 3
	condition1[3] = 1

	condition2[1] = 3
	condition2[2] = 6
	condition2[3] = 1

	condition3[1] = 1
	condition3[2] = 6
	condition3[3] = 3
}

type GetFaithModule struct{
	conf *conf.DataGetFaith  `bson:"-"`
	Faith  int32	 `bson:"faith"`           	//获取的信仰
}

//从持久化数据源中初始化数据
func (this *GetFaithModule)Init(data []byte, config *conf.StepData) {
	if (config != nil) {
		conf, ok := config.Data.(*conf.DataGetFaith)
		if (ok) {
			this.conf = conf
		}
	}
	if (data != nil) {
		json.Unmarshal(data, &this)
	}

}

func (this *GetFaithModule) Start(context db.EventContext)  {
	//this.changeDisplayOwner(context)
}

//处理时间限制到达
func (this *GetFaithModule) HandleTimesUp(context db.EventContext) {
	context.NextStep()
}

func (this *GetFaithModule) AppendReward(reward *protocol.Reward) {
	reward.Faith = this.Faith
}

//处理消息请求
func (this *GetFaithModule)HandleMessage(request *protocol.C2GS, response *protocol.GS2C, context db.EventContext) {
	message := request.GetGetFaith()
	if message == nil {
		return
	}

	var faith int32 = 0

	believerNum := len(message.GetBelieverID())
	var condition = condition3
	if believerNum == 1 {
		condition = condition1
	} else if believerNum == 2 {
		condition = condition2
	}

	result := util.RandomWeight(condition)
	if result == 1 {
		faith = character.RandInt32Scop(50, 150)
	} else if result == 2 {
		faith = character.RandInt32Scop(150, 250)
	} else {
		faith = character.RandInt32Scop(250, 600)
	}


	//for _, believerID := range  {
	//	faith += this.conf.RandomFaith(getBelieverLevel(believerID))
	//}


	//level := randomLevel(message.GetBelieverID(), this.conf.BelieverWeightMapping)
	//ratio := util.RandomWeight(this.conf.FaithRatioMapping)
	//this.Faith = level * ratio
	//believerID := randomBeliever(message.GetBelieverID(), level)
	this.Faith = faith

	response.GetFaithRet = &protocol.GetFaithRet{
		Faith: this.Faith,
		//BelieverID: believerID),
	}
	context.NextStep()
}


func randomBeliever(believers []string, level int32) string {
	for _, believer := range believers {
		if (getBelieverLevel(believer) == level) {
			return believer
		}
	}
	return ""
}

func randomLevel(believers []string, levelWeightMapping map[int32]int32) int32 {
	weightMapping := make(map[int32]int32)
	len := len(believers)
	if len == 0 {
		return 0
	}
	for _, believerID := range believers {
		level := getBelieverLevel(believerID)
		weightMapping[level] = levelWeightMapping[level]
	}
	return util.RandomWeight(weightMapping)
}

func getBelieverLevel(believerID string) int32 {
	levelStr := believerID[3:4]
	return character.StringToInt32(levelStr)
}


