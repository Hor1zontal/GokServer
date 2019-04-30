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
	"gok/service/exception"
	"gok/constant"
	"gok/service/rpc"
	"gok/module/game/db"
	"gok/module/game/words"
	"aliens/common/character"
	"gok/module/game/cache"
)

type RandomTargetModule struct{
	Target    *TargetMember `bson:"target"`     //事件目标
	config    *conf.DataRandomTarget	`bson:"-"`           	//配置信息
	AllBelieverLevel map[int32]int32
	AllBuildingLevel map[int32]int32
	reward *protocol.CardReward
}

type TargetMember struct {
	Uid  		int32 	`bson:"uid"`        //目标玩家id  0代表游戏外玩家
	Nickname  	string 	`bson:"nickname"`   //目标玩家昵称
}

func (this *TargetMember) IsRobot() bool {
	return this.Uid < 0
}

//从持久化数据源中初始化数据
func (this *RandomTargetModule)Init(data []byte, config *conf.StepData) {
	this.AllBelieverLevel = make(map[int32]int32)
	this.AllBuildingLevel = make(map[int32]int32)
	if (config != nil) {
		conf, ok := config.Data.(*conf.DataRandomTarget)
		if (ok) {
			this.config = conf
		}
	}
	if (data == nil) {
		return
	}
	json.Unmarshal(data, &this)
}

func (this *RandomTargetModule)Start(context db.EventContext)  {

}

//处理时间限制到达
func (this *RandomTargetModule)HandleTimesUp(context db.EventContext) {
	context.NextStep()
}

func (this *RandomTargetModule) AppendReward(reward *protocol.Reward) {
	if this.reward != nil {
		if this.reward.Type == constant.CARD_REWARD_FAITH || this.reward.Type == constant.CARD_REWARD_TARGET {
			//目标和其他星球都要追加奖励
			reward.Faith += this.reward.Value
		} else if this.reward.Type <= constant.MAX_BELIEVER_LEVEL {
			believerID := buildBelieverID(this.reward.Type)
			reward.Believer = append(reward.Believer, &protocol.BelieverInfo{Id:believerID, Num:this.reward.Value,})
		}
	}
}


//处理消息请求
func (this *RandomTargetModule)HandleMessage(request *protocol.C2GS, response *protocol.GS2C, context db.EventContext) {
	//处理
	if request.GetSelectEventTarget() != nil {
		message := request.GetSelectEventTarget()
		if this.Target != nil {
			exception.GameException(exception.EVENT_TARGET_ALREADY_SELECT)
		}
		this.Target = &TargetMember{
			Uid:message.GetTargetId(),
			Nickname:message.GetNickname(),
		}

		shieldType := constant.GetShieldType(context.GetType())

		result := rpc.StarServiceProxy.GetStarInfoDetail(message.GetTargetId(), 0, shieldType, true,
			this.AllBuildingLevel[message.GetTargetId()], this.AllBelieverLevel[message.GetTargetId()])

		response.SelectEventTargetRet = &protocol.SelectEventTargetRet{
			Result:true,
			StarInfo:result.Star,
			Shield:result.Shield,
		}
		context.NextStep()
	} else if request.GetRandomTarget() != nil {
		message := request.GetRandomTarget()
		message.EventType = context.GetType()
		message.Guide = context.IsGuide()
		if this.Target != nil {
			exception.GameException(exception.EVENT_TARGET_ALREADY_SELECT)
		}
		if this.config == nil {
			exception.GameException(exception.EVENT_CONFIG_EXCEPTION)
		}

		if message.Guide {
			result := rpc.StarServiceProxy.RandomGuideRobot(message.GetNum())
			response.RandomTargetRet = &protocol.RandomTargetRet{
				Targets:result.GetTargets(),
			}
		} else {
			response.RandomTargetRet = rpc.SearchServiceProxy.HandleMessage(request).GetRandomTargetRet()
			for _, target := range response.RandomTargetRet.Targets {
				UpdateTarget(target)
				this.AllBuildingLevel[target.Id] = target.BuildingTotalLevel
				this.AllBelieverLevel[target.Id] = target.BelieverTotalLevel
			}
		}

	} else if request.GetOpenCard() != nil {
		eventType := context.GetType()

		//翻到真实玩家
		ret := &protocol.OpenCardRet{}
		message := request.GetOpenCard()

		cardReward, otherRewards := conf.RandomCardReward(eventType)
		if cardReward == nil {
			exception.GameException(exception.CARD_REWARD_BASE_NOTFOUND)
		}

		randomTargetRequest := &protocol.RandomTarget{
			EventType:eventType,
			Uid:message.GetUid(),
			Num:1,
		}

		//没有随机到真人不需要加入到cd列表中来
		if cardReward.RewardType != constant.CARD_REWARD_TARGET {
			randomTargetRequest.AlwaysTarget = true
		}

		randomTargetRet := rpc.SearchServiceProxy.RandomTarget(randomTargetRequest)
		if len(randomTargetRet.Targets) < 0 {
			exception.GameException(exception.EVENT_INVALID_TARGET)
		}

		target := randomTargetRet.Targets[0]
		UpdateTarget(target)

		ret.Result = 1
		ret.Reward1 = this.transCardReward(eventType, cardReward, target, true)
		this.reward = ret.Reward1
		if otherRewards != nil && len(otherRewards) == 2 {
			ret.Reward2 = this.transCardReward(eventType, otherRewards[0], target, false)
			ret.Reward3 = this.transCardReward(eventType, otherRewards[1], target, false)
		}

		context.NextStep()

		//不是真人 直接下一步
		if cardReward.RewardType != constant.CARD_REWARD_TARGET {
			//直接进入下一步
			context.NextStep()
		}

		response.OpenCardRet = ret


	}
}

func  (this *RandomTargetModule) transCardReward(eventType int32, reward *conf.CardReward, target *protocol.Target, isResult bool) *protocol.CardReward {
	result := &protocol.CardReward{Type:reward.RewardType, Value:reward.RandomValue()}
	if reward.RewardType == constant.CARD_REWARD_TARGET {
		result.Target = target

		this.Target = &TargetMember{
			Uid:target.GetId(),
			Nickname:target.GetNickname(),
		}

		if isResult {
			shieldType := constant.GetShieldType(eventType)
			rpcRet := rpc.StarServiceProxy.GetStarInfoDetail(target.GetId(), 0, shieldType, true,
				target.BuildingTotalLevel, target.BelieverTotalLevel)

			result.StarInfo = rpcRet.GetStar()
		}
	}
	return result

}

func UpdateTarget(target *protocol.Target) {
	if target.GetId() < 0 {
		//target.Nickname = words.RandomName()
		target.Nickname = conf.DATA.ROBOT_MAPPPING[target.GetId()]
		if target.Nickname == "" {
			target.Nickname = words.RandomName()
		}
		target.Avatar = /*conf.DATA.RobotAvatarPrefix +*/ character.Int32ToString(-target.GetId()) + ".jpg"
	} else {
		target.Avatar = cache.UserCache.GetUserAvatar(target.GetId())
		target.Nickname = cache.UserCache.GetUserNickname(target.GetId())
	}
}



func getSelectTarget(context db.EventContext) *TargetMember{
	targetModule := context.GetModuleHandler(conf.MODULE_RANDOM_TARGET)
	if targetModule == nil {
		return nil
	}
	randomTargetModule, ok := targetModule.(*RandomTargetModule)
	if !ok {
		return nil
	}
	return randomTargetModule.Target
}

