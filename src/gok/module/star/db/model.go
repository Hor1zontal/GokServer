/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/5/5
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package db

import (
	"time"
)

type STAR_STATE int32

const (
	STAR_STATE_NONE STAR_STATE = iota
	STAR_STATE_DEV					//开发中
)

type History struct {
	ID          int32  		   	   `bson:"type"`         		 	      //记录ID
	Param1        int32			   `bson:"param1"`         		 	  //参数1
	Param2        int32			   `bson:"param2"`         		 	  //参数2
	Param3        string		   `bson:"param3"`         		 	  //参数3
	Time          time.Time		   `bson:"time"`         		 	  //更新时间
}

//星球
type DBStar struct {
	ID          int32        	  `bson:"_id"  gorm:"AUTO_INCREMENT"`        //星球id标示
	Type        int32             `bson:"type"`       //星球类型
	Seq         int32             `bson:"seq"`        //星球序号
	Owner       int32 		      `bson:"owner" unique:"false"`      //星球拥有的用户id  0代表没有人占领
	//State       STAR_STATE 		  `bson:"state"`      //星球状态  0原始星球 1开发中

	BuildingExMaxLevel int32      `bson:"buildMaxLevel"`       //建筑曾经达到的最大等级
	Building    []*DBBuilding     `bson:"building"`   //星球上的建筑
	Believer    []*DBBeliever 	  `bson:"believer"`   //星球上的信徒

	CivilizationLevel  int32     `bson:"civilevel"`       //当前等级文明度
	CivilizationValue  int32     `bson:"civilValue"`   //当前等级的文明度
	CivilizationReward []*CivilizationReward  `bson:"civilReward"`  //文明度奖励状态

	CreateTime  time.Time         `bson:"createtime"`  //星球创建时间
	OwnTime	    time.Time         `bson:"owntime"`     //星球占领的时间
	DoneTime    time.Time         `bson:"donetime"`     //星球的完成时间
	Statistics  []*Statistics `bson:"statistics"`  //星球统计数据
	History 	[]*History `bson:"history"`  //编年史信息
	Flags		[]*DBStarFlag `bson:"flag"` //星球标识
	BelieverUpdateTime time.Time  `bson:"updatetime"`  //信徒上次的更新时间

	Active   bool  		 `bson:"active"`  //是否启用 目前一个用户只有一个星球启用
	Disable  bool  		 `bson:"disable"` //禁用,不能被随机和探索到

	ItemGroups []*DBStarItemGroup `bson:"itemGroup"` //已开启的图鉴
	//CurrentGroupID int32 `bson:"currentGroupID"`		//当前解锁的圣物组合的ID
	//CurrentGroupIndex int32 `bson:"currentGroupIndex"`	//当前解锁第几个圣物组合

	FaithShield	   *Shield `bson:"FaithShield"`    //信仰防护罩
	BelieverShield *Shield `bson:"BelieverShield"` //信徒防护罩
	BuildingShield *Shield `bson:"BuildingShield"` //建筑防护罩

	Push bool `bson:"push"` //是否推送过消息
}

type Shield struct {
	Value int32				//当前数量
	Limit int32	 			//上限
	UpdateTime time.Time    //上次的刷新时间
}

//图鉴信息
type DBStarItemGroup struct {
	ID int32 `bson:"_id"`

	Done bool `bson:"done"` //是否达成组合
	Active bool `bson:"active"`
	Records []*DBStarItemGroupRecord `bson:"record"`
	//TryTime time.Time `bson:"trytime"`
	//UpdateTime	time.Time   `bson:"updatetime"`    //完成的时间
}

type DBStarItemGroupRecord struct {
	Items []int32 `bson:"items"`
	Num   int32   `bosn:"number"`
}

//月统计信息
type CivilizationReward struct {
	Level     int32     `bson:"level"       gorm:"PRIMARY_KEY"` //文明度等级
	Draw      bool      `bson:"reward"      gorm:"NOT NULL"`    //等级奖励是否领取
	Time	  time.Time `bson:"time"`							//文明度升级时间
}


//月统计信息
type Statistics struct {
	ID         int32     `bson:"_id"        gorm:"PRIMARY_KEY"` //统计标识
	Value      float64   `bson:"value"      gorm:"NOT NULL"`    //统计数值
	UpdateTime time.Time `bson:"updatetime"`                    //上次更新统计数据的时间
}

//统计信息
//type StarStatistics struct {
//	OnlineTime         int32 		`bson:"onlineTime"`         //当前星球的在线时间 单位s
//	EventNum           int32 		`bson:"eventNum"`           //当前星球的事件触发数量
//	AttackNum          int32        `bson:"attackNum"`          //当前星球的攻击数量
//	BuildNum           int32 		`bson:"buildNum"`           //当前星球的建造数量
//	UpgradeBelieverNum int32 		`bson:"upgradeBelieverNum"` //当前星球合成信徒的次数
//}


type SaleItem struct {
	ID          int32        	  `bson:"_id"`        //物品id
	Num         int32              `bson:"num"`       //物品数量
	Price       int32 		      `bson:"price"`      //物品单价
}




//信徒
type DBBeliever struct {
	ID        string    `bson:"_id"`       //信徒的id  参考gameObjectBase 配置
	Num 	  int32 	`bson:"num"`	   //信徒的数量
}

func (this *DBBeliever) Name() string {
	return ""
}

func (this *DBBeliever) GetID() interface{} {
	return this.ID
}

//槽点
type DBItemGroove struct {
	//Color 	int32 `bson:"color" json:"color,omitempty"`
	ItemID  int32 `bson:"itemid" json:"itemid,omitempty"`
	EffectTime time.Time  `bson:"effectTime" json:"effectTime,omitempty"` //槽内圣物生效时间
	Effect  bool `bson:"effect" json:"effect,omitempty"` //槽内圣物生效时间
}

//建筑
type DBBuilding struct {
	ID          int32      `bson:"_id"`        //建筑id  配置表中StartID+BuildID
	Type        int32      `bson:"type"`       //建筑类型
	Level       int32      `bson:"level"`      //建筑等级
	State       int32      `bson:"state"`      //建筑状态
	Exist    	bool       `bson:"exist"`      //建筑是否存在

	//RepairBelieverCost   int32      `bson:"repairCost"`  //修理花费的信徒数量
	//UpgradeBelieverCost  int32      `bson:"upgradeCost"` //升级花费的信徒数量
	Faith       int32		  `bson:"faith"`     //建筑存储的信仰值
	FaithUpdateTime  time.Time  `bson:"fUpdateTime"` //更新信仰值的时间戳
	CreateTime  time.Time  `bson:"createtime"` //建筑解锁时间
	RepairTime  time.Time  `bson:"repairtime"` //维护时间( 结束时间 )
	UpdateTime  time.Time  `bson:"updatetime"` //升级时间( 结束时间 )
	RobotHelpTime	time.Time  `bson:"robotHelptime"`
	BrokenTime  int64      `bson:"brokenTime"`  //待维护时间
	RemainTime  int64      `bson:"remainTime"`  //损坏时升级还剩余的时间，供修理成功的时候升级追加的时间
	//ItemGrooves []*DBItemGroove `bson:"itemgrooves"`  //建筑扩展的槽点dbbu
	LevelInfo    []*LevelInfo `bson:"levelinfo"` //每一级的完成时间
	//Buff               map[int32]int32 `bson:"-"`
	//BuffChangeListener func(int32, map[int32]int32) `bson:"-"`
}

func (this *DBBuilding) Name() string {
	return ""
}

func (this *DBBuilding) GetID() interface{} {
	return this.ID
}

type LevelInfo struct {
	ID       int32   `bson:"_id"`
	Time     int64   `bson:"time"`
	BelieverCost []*BelieverCost `bson:"cost"` //当前等级信徒的消耗情况
	FaithCost   int32	   `bson:"faithcost"`     //建筑修理总共消耗的信仰值
}

//信徒消耗
type BelieverCost struct {
	ID      string  `bson:"_id"`
	Num     int32   `bson:"num"`
}

type DBStarFlag struct {
	Flag       int32 	 `bson:"flag"`       //标识
	Value      int32     `bson:"value"`      //标识值
	UpdateTime time.Time `bson:"updatetime"` //上次标识修改的时间
}