package db

import (
	"gok/service/msg/protocol"
	"time"
	"gok/module/game/conf"
	"gok/constant"
	"aliens/common/character"
)


func (this *DBNewsFeed) BuildProtocol() *protocol.NewsFeed {
	return &protocol.NewsFeed{
		Id:this.ID,
		RelateID:this.Uid,
		//RelateNickname:nickname),
		Type:this.Type,
		Time:this.Time.Unix(),
		Param1:this.Param1,
		Param2:this.Param2,
		Param3:this.Param3,
		Ext:this.Ext,
		//IsRevenge:this.IsRevenge),
		DoneRevenge:this.DoneRevenge,
		Read:this.Read,
	}
}

func (this *DBNewsFeed) UpdateDetail() {
	if constant.IsAttackNewsFeed(this.Type) {
		if this.Self == nil {
			this.Self = &protocol.NewsFeedDetail{}
		}
		if this.Type == constant.NEWSFEED_TYPE_LOOT_FAITH {
			this.Self.Faith = this.Self.GetFaith() + this.Param1
		} else if this.Type == constant.NEWSFEED_TYPE_ATK_BUILD {
			this.Self.Faith = this.Self.GetFaith() + this.Param1
			//抢夺圣物
			if this.Param2 != 0 {
				this.Self.ItemNum = this.Self.GetItemNum() + 1
			}
			//攻打建筑
			if this.Ext != nil {
				for _, buildingLevel := range this.Ext {
					addBuildingNum(this.Self, character.StringToInt32(buildingLevel))
				}
			}
		} else if this.Type == constant.NEWSFEED_TYPE_LOOT_BELIEVER {
			if this.Ext != nil {
				for _, believerID := range this.Ext {
					addBelieverNum(this.Self, believerID, 1)
				}
			}
		}
	}

	if constant.IsBeAttackNewsFeed(this.Type) {
		if this.Other == nil {
			this.Other = &protocol.NewsFeedDetail{}
		}

		if this.Type == constant.NEWSFEED_TYPE_BE_LOOT_FAITH {
			this.Other.Faith = this.Self.GetFaith() + this.Param1
		} else if this.Type == constant.NEWSFEED_TYPE_BE_ATK_BUILD {
			this.Other.Faith = this.Self.GetFaith() + this.Param1
			//抢夺圣物
			if this.Param2 != 0 {
				this.Other.ItemNum = this.Self.GetItemNum() + 1
			}
			//攻打建筑
			if this.Ext != nil {
				for _, buildingLevel := range this.Ext {
					addBuildingNum(this.Other, character.StringToInt32(buildingLevel))
				}
			}
		} else if this.Type == constant.NEWSFEED_TYPE_BE_LOOT_BELIEVER {
			if this.Ext != nil {
				for _, believerID := range this.Ext {
					addBelieverNum(this.Other, believerID, 1)
				}
			}
		}
	}
}

func addBelieverNum(detail *protocol.NewsFeedDetail, believerID string, believerNum int32) {
	for _, believer := range detail.BelieverInfo {
		if believer.GetId() == believerID {
			believer.Num = believer.GetNum() + believerNum
		}
	}
	detail.BelieverInfo = append(detail.BelieverInfo, &protocol.BelieverInfo{
		Id:believerID,
		Num:believerNum,
	})
}

func addBuildingNum(detail *protocol.NewsFeedDetail, buildingLevel int32) {
	for _, statistics := range detail.AttackStatistics {
		if statistics.GetLevel() == buildingLevel {
			statistics.Num = statistics.GetNum() + 1
		}
	}
	detail.AttackStatistics = append(detail.AttackStatistics, &protocol.AttackStatistics{
		Level:buildingLevel,
		Num:1,
	})
}

//是否过期
func (this *DBNewsFeed) IsOverdue() bool {
	return time.Now().Sub(this.Time).Seconds() > float64(conf.DATA.CountdownApplying)
}


func (flag *DBRoleFlag) BuildProtocol() *protocol.FlagInfo {
	return &protocol.FlagInfo{
		Id:int32(flag.Flag),
		Value:flag.Value,
		Time:flag.UpdateTime.Unix(),
	}
}
