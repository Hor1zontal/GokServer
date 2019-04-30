/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/30
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package model

import (
	"github.com/sirupsen/logrus"
	"gok/constant"
)

type IStatisticData interface {
	GetData() logrus.Fields
	GetName() string
}

//星盘任务统计
type StatisticDial struct {
	Uid 		int32 //用户id
	DialID 		int32 //转盘id
	Faith 		int32 //奖励信仰值
	Power       int32 //奖励法力值
	Believer    int32 //奖励平民数量
	Event       int32 //随机到的事件类型
	GayPoint	int32 //奖励圣物碎片
}

func (this *StatisticDial) GetData() logrus.Fields {
	fields := logrus.Fields{}
	fields["uid"] = this.Uid
	fields["dialID"] = this.DialID
	if this.Faith != 0 {
		fields["faith"] = this.Faith
	}
	if this.Power != 0 {
		fields["power"] = this.Power
	}
	if this.Believer != 0 {
		fields["believer"] = this.Believer
	}
	if this.Event != 0 {
		fields["event"] = this.Event
	}
	if this.GayPoint != 0 {
		fields["gayPoint"] = this.GayPoint
	}
	return fields
}

func (this *StatisticDial) GetName() string {
	return "dial"
}

//交互事件统计
type StatisticEvent struct {
	Uid           int32 //用户id
	Event         int32 //事件类型
	Revenge       bool  //是否复仇事件
	Source        int32 //时间出发来源 0 转盘 1复仇
	TargetID      int32 //交互目标id
	//TargetSeq     int32 //交互目标序号 -1代表非转盘随机目标产生的事件 - 报复
	BelieverTotal int32 //目标信徒总等级
	BuildingTotal int32 //目标建筑总等级
	Civil         int32 //目标的文明度等级
	BelieverLoot  int32 //抢夺的信徒数
	FaithLoot     int32 //抢夺的信仰数
	BuildAtk      int32 //毁坏的建筑等级
	FaithGet      int32 //获得的总信仰
	CardType      int32 //翻卡片的类型
}

func (this *StatisticEvent) GetData() logrus.Fields {
	fields := logrus.Fields{}
	fields["uid"] = this.Uid
	fields["event"] = this.Event
	if this.Revenge {
		fields["revenge"] = true
	}
	if this.Source != 0 {
		fields["source"] = this.Source
	}
	fields["targetID"] = this.TargetID
	//fields["targetSeq"] = this.TargetSeq

	if this.TargetID != 0 {
		fields["believerTotal"] = this.BelieverTotal
		fields["buildingTotal"] = this.BuildingTotal
		fields["civil"] = this.Civil
	}

	fields["faithGet"] = this.FaithGet
	fields["cardType"] = this.CardType
	if this.CardType == constant.CARD_REWARD_TARGET {
		switch this.Event {
		case constant.EVENT_ID_LOOT_FAITH:
			fields["faithLoot"] = this.FaithLoot
		case constant.EVENT_ID_LOOT_BELIEVER:
			fields["believerLoot"] = this.BelieverLoot
		case constant.EVENT_ID_ATK_BUILDING:
			fields["buildAtk"] = this.BuildAtk
		}
	}
	return fields
}

func (this *StatisticEvent) GetName() string {
	return "event"
}


type StatisticItem struct {
	UserID   int32     `json:"uid"`              //用户id
	ItemID   int32     `json:"tid"`              //用户id
	RefID    int32     `json:"rid"`              //关联id
	Operation uint8    `json:"opt"`              //操作类型
	Change   int32     `json:"change"`           //改变数量
	Total    int32     `json:"total"`            //改变后的总数
	//Time     time.Time `json:"time"`             //数据变更后的时间
}

func (this *StatisticItem) GetData() logrus.Fields {
	fields := logrus.Fields{}
	fields["uid"] = this.UserID
	fields["tid"] = this.ItemID
	fields["rid"] = this.RefID
	fields["opt"] = this.Operation
	fields["change"] = this.Change
	fields["total"] = this.Total
	return fields
}

func (this *StatisticItem) GetName() string {
	return "item"
}


type StatisticGayPoint struct {
	UserID   int32     `json:"uid"`              //用户id
	RefID    int32     `json:"rid"`              //关联id
	Operation uint8    `json:"opt"`              //操作类型
	Change   int32     `json:"change"`           //改变数量
	Total    int32     `json:"total"`            //改变后的总数
}

func (this *StatisticGayPoint) GetData() logrus.Fields {
	fields := logrus.Fields{}
	fields["uid"] = this.UserID
	fields["rid"] = this.RefID
	fields["opt"] = this.Operation
	fields["change"] = this.Change
	fields["total"] = this.Total
	return fields
}

func (this *StatisticGayPoint) GetName() string {
	return "gaypoint"
}

type StatisticItemGroup struct {
	UserID   int32  //用户id
	StarType int32  //新球类型
	GroupID  int32  //圣物组合id
	Seq      int  //圣物组合索引
	Type     int32  //0 解锁圣物组合 1完成圣物组合
}

const (
	ITEM_GROUP_TYPE_ACTIVE int32 = 0
	ITEM_GROUP_TYPE_DONE int32 = 1
)

func (this *StatisticItemGroup) GetData() logrus.Fields {
	fields := logrus.Fields{}
	fields["uid"] = this.UserID
	fields["star"] = this.StarType
	fields["id"] = this.GroupID
	fields["seq"] = this.Seq
	fields["groupType"] = this.Type
	return fields
}

func (this *StatisticItemGroup) GetName() string {
	return "itemgroup"
}

type StatisticLogout struct {
	UserID   int32     `json:"uid"`              //用户id
	BelieverTotal int32
	BuildingTotal int32
	Faith int32
	Power int32
	ExMaxBuildLv int32
	//Civil   int32 //文明度等级
	Address string
	Star int32
	Mutual int32
	BeMutual int32
}

func (this *StatisticLogout) GetData() logrus.Fields {
	fields := logrus.Fields{}
	fields["uid"] = this.UserID
	fields["believerTotal"] = this.BelieverTotal
	fields["buildingTotal"] = this.BuildingTotal
	fields["faith"] = this.Faith
	fields["power"] = this.Power
	fields["exMaxBuildLv"] = this.ExMaxBuildLv
	//fields["civil"] = this.Civil
	fields["address"] = this.Address
	fields["star"] = this.Star
	fields["mutual"] = this.Mutual
	fields["beMutual"] = this.BeMutual
	return fields
}

func (this *StatisticLogout) GetName() string {
	return "logout"
}


type StatisticWechat struct {
	OpenID    string `json:"openID"` //公众号openid
	Uid       int32  `json:"uid"`    //游戏id 0 游戏内没有账号
	Event 	  string `json:"event"`  //事件id
}

func (this *StatisticWechat) GetData() logrus.Fields {
	fields := logrus.Fields{}
	fields["openID"] = this.OpenID
	fields["uid"] = this.Uid
	fields["event"] = this.Event
	return fields
}

func (this *StatisticWechat) GetName() string {
	return "wechat"
}


////偷取圣物统计
//type StatisticSteal struct {
//	Uid           int32 //偷取的用户id
//	TargetID      int32 //偷取的目标id
//	ItemID  	  int32 //偷取的物品id
//	Count         int32 //当前偷取成功的次数
//	ItemNum       int32 //物品剩余数量
//	Success	      bool     //是否偷取成功
//	Prob          float32  //偷取成功的概率
//}
//
//func (this *StatisticSteal) GetData() logrus.Fields {
//	fields := logrus.Fields{}
//	fields["uid"] = this.Uid
//	fields["targetID"] = this.TargetID
//	fields["count"] = this.Count
//	fields["itemID"] = this.ItemID
//	fields["itemNum"] = this.ItemNum
//	fields["success"] = this.Success
//	fields["prob"] = this.Prob
//	return fields
//}
//
//func (this *StatisticSteal) GetName() string {
//	return "steal"
//}


//type StatisticTarget struct {
//	Uid  	  int32 //用户id
//	TargetID  int32 //交互目标id
//	TargetSeq int32 //交互目标序号
//	Event 	  int32 //事件类型
//	Mutual    bool //是否复仇
//}
//
//func (this *StatisticTarget) GetData() logrus.Fields {
//	fields := logrus.Fields{}
//	fields["uid"] = this.Uid
//	fields["dialID"] = this.DialID
//	if this.Faith != 0 {
//		fields["faith"] = this.Faith
//	}
//	if this.Power != 0 {
//		fields["power"] = this.Power
//	}
//	if this.Believer != 0 {
//		fields["believer"] = this.Believer
//	}
//	if this.Event != 0 {
//		fields["event"] = this.Event
//	}
//	return fields
//}
//
//func (this *StatisticTarget) GetName() string {
//	return "target"
//}


type CallInfo struct {
	count int32 	 //调用次数
	interval float64 //调用时间总长
}

func (this *CallInfo) AddCall(interval float64) {
	this.count ++
	this.interval += interval
}

func (this *CallInfo) IsEmpty() bool {
	return this.count == 0 || this.interval == 0
}

func (this *CallInfo) DumpData() (bool, int32, float64) {
	if this.IsEmpty() {
		return false, 0, 0
	}
	avg := this.interval / float64(this.count)
	count := this.count

	this.count = 0
	this.interval = 0

	return true , count, avg
}
