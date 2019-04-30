/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/5/17
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     zhangyishen
 *******************************************************************************/
package db

import (
	"time"
)

//组合ID
type CompID struct {
	SubID1 		int32			`bson:"s1"`   //id1
	SubID2		int32			`bson:"s2"`   //id2
}

//交易
type UserSale struct {
	UserID	      int32 	       `bson:"_id"`         		 	 //用户id
	ItemID        int32            `bson:"itemID"`         		 	 //用户发布的物品id
	PublicTime    time.Time        `bson:"publicTime"`         		 //用户发布的物品是时间
	//RefID         string           `bson:"refID"`         		 //
}

//交易物品
type UserGoods struct {
	CompID	      CompID 	       `bson:"_id"`         		 	 //
	Num	          int32 	       `bson:"num"`         		 	 //数量
	Price	      int32 	       `bson:"price"`         		 	 //价格
}

func (this *UserGoods) GetUserID() int32 {
	return this.CompID.SubID1
}

func (this *UserGoods) GetItemID() int32 {
	return this.CompID.SubID2
}

func (this *UserGoods) GetID() interface{} {
	return this.CompID
}

type Trade struct {
	ID	      int32 	`bson:"_id"`
	ItemId    int32		`bson:"itemID"`         		 //当前可以领取的圣物id
	ItemNum   int32     `bson:"itemNum"`         		 //当前可以领取的圣物数量
	DrawNum   int32     `bson:"drawNum"`         		 //已经领取的圣物数量
	LootNum   int32     `bson:"itemNum"`         		 //总共被抢走的数量
	Events	  []*Event  `bson:"events"`         		 //总共被抢走的数量

}

type Event struct {
	Uid        int32
	Type   	   int32  //0援助 1偷窃失败 2偷窃成功
	Count      int32  //操作次数
}


type ItemHelpHistory struct {
	ID	      string 	`bson:"_id"`
	Uid       int32 	`bson:"uid" unique:"false"`
	Data      []byte    `bson:"data"`
	Time      time.Time `bson:"time"`
}