package session

import (
	"gok/module/star/cache"
	"gok/module/star/conf"
	"gok/module/star/db"
	"gok/module/star/util"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"gok/service/rpc"
	"time"
)

func (this *StarSession) correctCivilizationValue() {
	//机器人
	if this.Owner < 0 {
		return
	}
	seq := this.Seq
	var power int32 = 0
	for {
		level := this.CivilizationLevel
		confCivil := conf.GetCivilConfigByLevel(seq, level + 1)
		if confCivil == nil || confCivil.Threshold <= 0 {
			break
		}
		if this.CivilizationValue < confCivil.Threshold  {
			break
		}
		//先升级再领取上一个文明度的奖励
		remainValue := this.CivilizationValue - confCivil.Threshold

		if remainValue >= 0 {
			reward := &db.CivilizationReward{
				Level:this.CivilizationLevel,
				Draw:false,
				Time:time.Now(),
			}
			this.CivilizationReward = append(this.CivilizationReward, reward)
			this.CivilizationLevel += 1
			cache.StarCache.SetCivilLevel(this.ID, this.CivilizationLevel)
			this.CivilizationValue = remainValue
		}

		powerValue, _, _, _, _, _ := this.DrawStarReward(level)
		power += powerValue
	}

	if power > 0 {
		message := util.BuildCorrectCivilRewardMessage(power)
		rpc.UserServiceProxy.UserHandleMessage(this.Owner, message)
		//rpc.UserServiceProxy.DrawCorrectCivilReward(this.Owner, power)
	}
}

func (this *StarSession) TakeInCivilization(value int32) *protocol.CivilizationInfo{
	this.CivilizationValue += value

	remainValue := this.CivilizationValue
	for {
		if remainValue <= 0 {
			break
		}
		config := conf.GetCivilConfigByLevel(this.Seq, this.CivilizationLevel + 1)
		if config == nil || config.Threshold <= 0 {
			break
		}

		remainValue = this.CivilizationValue - config.Threshold
		if remainValue >= 0 {
			reward := &db.CivilizationReward{
				Level:this.CivilizationLevel,
				Draw:false,
				Time:time.Now(),
			}
			this.CivilizationReward = append(this.CivilizationReward, reward)
			this.CivilizationLevel += 1
			cache.StarCache.SetCivilLevel(this.ID, this.CivilizationLevel)
			this.CivilizationValue = remainValue
		}
	}

	return util.BuildCivilizationInfo(this.DBStar)
	//if this.Owner > 0 {
	//	rpc.UserServiceProxy.Push(this.Owner, cache.StarCache.GetUserNode(this.Owner), util.BuildRoleCivilizationPush(this.DBStar))
	//}
}

func (this *StarSession) getCivilizationReward(level int32) *db.CivilizationReward {
	for _, reward := range this.CivilizationReward {
		if reward.Level == level {
			return reward
		}
	}
	return nil
	//reward := &db.CivilizationReward{
	//	Level:level,
	//	Draw:false,
	//}
	//this.CivilizationReward = append(this.CivilizationReward, reward)
	//return reward
}

//处理建筑信仰值  draw 是否领取
func (this *StarSession) DrawStarReward(level int32) (int32, int32, int32, int32, []string, []int32) {
	if this.CivilizationLevel != conf.GetCivilMaxLevel(this.Seq) && this.CivilizationLevel <= level {
		exception.GameException(exception.CIVILIZATION_LEVEL_NOT_ENOUGH)
	}
	reward := this.getCivilizationReward(level)
	if reward == nil {
		exception.GameException(exception.CIVILIZATION_REWARD_NOT_FOUND)
	}
	if reward.Draw {
		exception.GameException(exception.CIVILIZATION_REWARD_ALREADY_DRAW)
	}

	var rewardValue int32 = 0
	var faithValue int32 = 0
	var diamondValue int32 = 0
	var relicPointValue int32 = 0
	var giftValue int32 = 0
	var believer = []string{}
	var believerNum = []int32{}
	reward.Draw = true
	config := conf.GetCivilConfigByLevel(this.Seq, level)
	if config != nil {
		rewardValue = config.Reward
		giftValue = config.Gift
	}
	if giftValue != 0 {
		civilizationReward := conf.GetCivilRewardConfigById(giftValue)
		if civilizationReward == nil {
			exception.GameException(exception.CIVILIZATION_REWARD_NOT_FOUND)
		}
		for index, believerLevel := range civilizationReward.BelieverLevel {
			maleNum := civilizationReward.BelieverNum[index]/2
			femaleNum := civilizationReward.BelieverNum[index] - maleNum
			maleBelievers := this.addCustomBeliever(believerLevel,true,maleNum)
			femaleBelievers := this.addCustomBeliever(believerLevel,false,femaleNum)
			believer = append(believer, maleBelievers.ID)
			believer = append(believer, femaleBelievers.ID)
			believerNum = append(believerNum, maleBelievers.Num)
			believerNum = append(believerNum, femaleBelievers.Num)
		}
		faithValue = civilizationReward.FaithNum
		diamondValue = civilizationReward.DiamondNum
		relicPointValue = civilizationReward.RelicPointNum
	}
	this.setDirty()
	return rewardValue, faithValue, diamondValue, relicPointValue, believer, believerNum
}