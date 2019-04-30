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
	"gok/service/msg/protocol"
	"aliens/common/util"
	"gok/module/game/db"
)

type GetBelieverModule struct{
	conf *conf.DataGetBeliever  `bson:"-"`
	Believers  []*BelieverInfo	`bson:"believer"`           //获取的原始人id
}

//从持久化数据源中初始化数据
func (this *GetBelieverModule)Init(data []byte, config *conf.StepData) {
	if (config != nil) {
		conf, ok := config.Data.(*conf.DataGetBeliever)
		if (ok) {
			this.conf = conf
		}
	}
	if (data != nil) {
		json.Unmarshal(data, &this)
	}
}

func (this *GetBelieverModule)Start(context db.EventContext)  {
	//this.changeDisplayOwner(context)
}

//处理时间限制到达
func (this *GetBelieverModule)HandleTimesUp(context db.EventContext) {
	context.NextStep()
}

func (this *GetBelieverModule) AppendReward(reward *protocol.Reward) {
	if (this.Believers != nil) {
		for _, believer := range this.Believers {
			reward.Believer = append(reward.Believer, believer.getProtocol())
		}
	}
}

//const (
//	PEOPLE_A = "b01" +1 +"1"
//	PEOPLE_B = "b0112"
//)

//处理消息请求
func (this *GetBelieverModule)HandleMessage(request *protocol.C2GS, response *protocol.GS2C, context db.EventContext) {
	if request.GetGetBeliever() != nil {
		peopleNum := util.RandomWeight(this.conf.PeopleWeightMapping)
		response.GetBelieverRet = &protocol.GetBelieverRet{}
		if peopleNum > 0 {
			results := []*protocol.BelieverInfo{}
			this.Believers = []*BelieverInfo{}

			believerInfo := make(map[string]int32)
			for i := 0; i < int(peopleNum); i ++ {
				//PeopleANum := int32(peopleNum / 2)
				//PeopleBNum := peopleNum - PeopleANum

				believerLevel := util.RandomWeight(this.conf.BelieverWeightMapping)
				randomBeliever := buildBelieverID(believerLevel)
				believerInfo[randomBeliever] += 1
			}

			for believerID, believerNum := range believerInfo {
				peopleA := &BelieverInfo{ID: believerID, Num: believerNum}
				this.Believers = append(this.Believers, peopleA)
				results = append(results, peopleA.getProtocol())
			}
			response.GetBelieverRet.Believer = results
		}
		context.NextStep()
	}

}






