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
	"encoding/json"
	"gok/constant"
	"gok/module/game/conf"
	"gok/module/game/db"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"gok/service/rpc"
)

type LootFaithModule struct{
	conf *conf.DataGetFaith  `bson:"-"`
	Faith  int32	 `bson:"faith"`           	//获取的信仰
}

//从持久化数据源中初始化数据
func (this *LootFaithModule)Init(data []byte, config *conf.StepData) {
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

func (this *LootFaithModule)Start(context db.EventContext)  {
	//this.changeDisplayOwner(context)
}

//处理时间限制到达
func (this *LootFaithModule)HandleTimesUp(context db.EventContext) {
	context.NextStep()
}

func (this *LootFaithModule) AppendReward(reward *protocol.Reward) {
	reward.Faith += this.Faith
}

//处理消息请求
func (this *LootFaithModule)HandleMessage(request *protocol.C2GS, response *protocol.GS2C, context db.EventContext) {
	message := request.GetLootFaith()
	if message != nil {
		target := getSelectTarget(context)
		//无效的进攻目标
		if target == nil || target.Uid == 0 {
			exception.GameException(exception.EVENT_INVALID_TARGET)
		}
		resp := &protocol.LootFaithRet{}
		//response.LootFaithRet = &protocol.LootFaithRet{}

		mutual := false
		mutual = message.GetMutual()
		if !target.IsRobot() && mutual {
			result := rpc.StarServiceProxy.Call(target.Uid, request).GetLootFaithRet()
			//if result.GetShield() || result.GetIsMax() {
				//response.LootFaithRet = &protocol.LootFaithRet{
				//	TargetID:  target.Uid,
				//	Shield: result.GetShield(),
				//	IsMax: result.GetIsMax(),
				//}
			//}
			resp.Shield = result.GetShield()
			resp.IsMax = result.GetIsMax()
			resp.HasBuilding = result.GetHasBuilding()
		}

		this.Faith = message.GetFaith()
		resp.TargetID = target.Uid
		resp.Faith = this.Faith
		response.LootFaithRet = resp
		//统计抢到的信仰数
		context.AppendStatisticData(constant.EVENT_ID_LOOT_FAITH, this.Faith)
		context.NextStep()
	}
}




