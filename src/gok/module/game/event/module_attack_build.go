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
	"aliens/common/util"
	"encoding/json"
	"gok/constant"
	"gok/module/game/conf"
	"gok/module/game/db"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"gok/service/rpc"
)

var failedMessage = &protocol.AtkStarBuildingRet{
	Success: false,
	Faith:   0,
	ItemID:  0,
	Shield:  false,
}

type AttackBuildingModule struct {
	conf *conf.DataAttackBuild `bson:"-"`

	Faith  int32 `bson:"faith"`  //攻击建筑获取到的信仰值
	ItemID int32 `bson:"itemID"` //攻击建筑抢夺到的物品id
}

//从持久化数据源中初始化数据
func (this *AttackBuildingModule) Init(data []byte, config *conf.StepData) {
	if (config != nil) {
		conf, ok := config.Data.(*conf.DataAttackBuild)
		if (ok) {
			this.conf = conf
		}
	}
	if (data != nil) {
		json.Unmarshal(data, &this)
	}
}

func (this *AttackBuildingModule) Start(context db.EventContext) {
	//this.changeDisplayOwner(context)
}

//处理时间限制到达
func (this *AttackBuildingModule) HandleTimesUp(context db.EventContext) {
	context.NextStep()
}

func (this *AttackBuildingModule) AppendReward(reward *protocol.Reward) {
	reward.Faith += this.Faith
	reward.ItemID = this.ItemID
}

//处理消息请求
func (this *AttackBuildingModule) HandleMessage(request *protocol.C2GS, response *protocol.GS2C, context db.EventContext) {
	message := request.GetAtkStarBuilding()
	if message == nil {
		return
	}

	target := getSelectTarget(context)
	//无效的进攻目标
	if target == nil || target.Uid == 0 {
		exception.GameException(exception.EVENT_INVALID_TARGET)
	}
	attackUid := context.GetCaller().Uid
	if attackUid == target.Uid {
		exception.GameException(exception.STAR_CANNOT_ATK_SELF)
	}

	response.AtkStarBuildingRet = &protocol.AtkStarBuildingRet{Success: message.GetSuccess()}
	message.AttackUid = attackUid
	message.DestUid = target.Uid

	//未校验等级与可打的血量
	if this.conf.FailedRatio == 0 {
		this.conf.FailedRatio = 1
	}
	baseFaith :=  int32(this.conf.FailedRatio * float32(request.GetAtkStarBuilding().GetBuildingHurt()))

	mutual := message.GetMutual() //解锁交互
	if message.GetSuccess() && mutual{
		//成功统计打的建筑等级
		context.AppendStatisticData(constant.EVENT_ID_ATK_BUILDING, request.GetAtkStarBuilding().GetBuildingLevel())

		ratio := util.RandomFloat32Weight(this.conf.FaithRatioMapping)
		message.FaithRatio = ratio
		response.AtkStarBuildingRet = rpc.StarServiceProxy.Call(message.GetDestUid(), request).GetAtkStarBuildingRet()
		response.AtkStarBuildingRet.Faith += baseFaith
	} else {
		failedMessage.Faith = int32(baseFaith)
		response.AtkStarBuildingRet = failedMessage
	}
	response.AtkStarBuildingRet.TargetID = target.Uid

	this.Faith = response.AtkStarBuildingRet.GetFaith()
	this.ItemID = response.AtkStarBuildingRet.GetItemID()
	//完成，切换到下一步
	context.NextStep()
}
