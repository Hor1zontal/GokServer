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
	"gok/service/exception"
	"gok/service/rpc"
	"gok/constant"
)

type LootBelieverModule struct{
	conf *conf.DataLootBeliever  `bson:"-"`
	Believer  []*BelieverInfo	`bson:"believer"`           //获取的原始人id
}


//从持久化数据源中初始化数据
func (this *LootBelieverModule)Init(data []byte, config *conf.StepData) {
	if (config != nil) {
		conf, ok := config.Data.(*conf.DataLootBeliever)
		if (ok) {
			this.conf = conf
		}
	}
	if (data != nil) {
		json.Unmarshal(data, &this)
	}

}

func (this *LootBelieverModule)Start(context db.EventContext)  {
	//this.changeDisplayOwner(context)
}

//处理时间限制到达
func (this *LootBelieverModule)HandleTimesUp(context db.EventContext) {
	context.NextStep()
}

func (this *LootBelieverModule) AppendReward(reward *protocol.Reward) {
	if this.Believer != nil {
		for _, believer := range this.Believer {
			reward.Believer = append(reward.Believer, believer.getProtocol())
		}
	}
}

//处理消息请求
func (this *LootBelieverModule)HandleMessage(request *protocol.C2GS, response *protocol.GS2C, context db.EventContext) {
	message := request.GetLootBeliever()
	//believerID := ""
	if message != nil {
		target := getSelectTarget(context)
		//无效的进攻目标
		if target == nil || target.Uid == 0 {
			exception.GameException(exception.EVENT_INVALID_TARGET)
		}

		lootBelievers := message.GetBelieverID()
		shield := false
		isMax := false
		mutual := message.GetMutual() //解锁交互
		if !target.IsRobot() && mutual {
			resp := rpc.StarServiceProxy.LootStarBeliever(context.GetCaller().Uid, target.Uid, message.GetBelieverID())
			isMax = resp.GetIsMax()
			shield = resp.GetShield()
			lootBelievers = resp.GetBelieverID()
			if lootBelievers != nil && len(lootBelievers) > 0 {
				rpc.StarServiceProxy.UpdateStarStatistics(context.GetCaller().Uid, constant.STAR_STATISTIC_LOOT_BELIEVER_NUM, float64(len(lootBelievers)), 0)
			}

		}
		//统计抢信徒的个数
		if lootBelievers != nil && len(lootBelievers) > 0 {
			context.AppendStatisticData(constant.EVENT_ID_LOOT_BELIEVER, int32(len(lootBelievers)))
		}

		believerInfo := make(map[string]int32)
		for _, believerID := range message.GetBelieverID() {
			believerInfo[believerID] = believerInfo[believerID] + 1
		}

		this.Believer = []*BelieverInfo{}
		for believerID, believerNum := range believerInfo {
			this.Believer = append(this.Believer, &BelieverInfo{
				ID:believerID,
				Num:believerNum,
			})
		}

		response.LootBelieverRet = &protocol.LootBelieverRet{
			BelieverID: lootBelievers,
			TargetID:  	target.Uid,
			Shield: 	shield,
			IsMax:		isMax,
		}
		context.NextStep()
	}
}




