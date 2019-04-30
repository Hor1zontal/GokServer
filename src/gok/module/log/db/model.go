/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package db

import (
	"time"
)


type ServiceRecord struct {
	ID string 			//服务id
	Count int32 		//调用次数
	Interval float32    //平均响应时间

}

type ItemRecord struct {
	ID       int32     `bson:"_id" gorm:"AUTO_INCREMENT"` //记录id
	UserID   int32     `bson:"uid"`              //用户id
	ItemID   int32     `bson:"tid"`              //用户id
	RefID    int32     `bson:"rid"`              //关联id
	Operation uint8    `bson:"opt"`              //操作类型
	Change   int32     `bson:"change"`           //改变数量
	Total    int32     `bson:"total"`            //改变后的总数
	Time     time.Time `bson:"time"`             //数据变更后的时间
}

//引导的完成时间
type GuideRecord struct {
	ID       int32     `bson:"_id"` 			 //用户id
	Time     time.Time `bson:"time"`             //引导完成的时间
	Duration float64   `bson:"duration"`         //引导完成后的在线时间
}

//SOCIAL_ID_FAITH SOCIAL_ID = 1    //信仰
//SOCIAL_ID_POWER SOCIAL_ID = 2    //法力
//SOCIAL_ID_DIAMOND SOCIAL_ID = 3   //钻石
//SOCIAL_ID_GAYPOINT SOCIAL_ID = 4  //圣物碎片

type FaithRecord struct {
	ID       int32     `bson:"_id" gorm:"AUTO_INCREMENT"` //记录id
	UserID   int32     `bson:"uid"`              //用户id
	RefID    int32     `bson:"rid"`              //关联id
	Change   int32     `bson:"change"`           //改变数量
	Operation uint8    `bson:"opt"`             //操作类型
	Total    int32     `bson:"total"`            //改变后的总数
	Time     time.Time `bson:"time"`             //数据变更后的时间
}

type PowerRecord struct {
	ID       int32     `bson:"_id" gorm:"AUTO_INCREMENT"` //记录id
	UserID   int32     `bson:"uid"`              //用户id
	RefID    int32     `bson:"rid"`              //关联id
	Change   int32     `bson:"change"`           //改变数量
	Operation uint8    `bson:"opt"`             //操作类型
	Total    int32     `bson:"total"`            //改变后的总数
	Time     time.Time `bson:"time"`             //数据变更后的时间
}

type DiamondRecord struct {
	ID       int32     `bson:"_id" gorm:"AUTO_INCREMENT"` //记录id
	UserID   int32     `bson:"uid"`              //用户id
	RefID    int32     `bson:"rid"`              //关联id
	Change   int32     `bson:"change"`           //改变数量
	Operation uint8    `bson:"opt"`             //操作类型
	Total    int32     `bson:"total"`            //改变后的总数
	Time     time.Time `bson:"time"`             //数据变更后的时间
}

type GayPointRecord struct {
	ID       int32     `bson:"_id" gorm:"AUTO_INCREMENT"` //记录id
	UserID   int32     `bson:"uid"`              //用户id
	RefID    int32     `bson:"rid"`              //关联id
	Change   int32     `bson:"change"`           //改变数量
	Operation uint8    `bson:"opt"`             //操作类型
	Total    int32     `bson:"total"`            //改变后的总数
	Time     time.Time `bson:"time"`             //数据变更后的时间
}

//登录日志
type LoginRecord struct {
	ID         int64      `bson:"_id" gorm:"AUTO_INCREMENT"`        //
	UserID     int32      `bson:"userid"`     //用户id
	Ip         string     `bson:"ip"`         //登录ip
	LoginTime  time.Time  `bson:"loginTime"`  //登录时间
	LogoutTime time.Time  `bson:"logoutTime"` //登出时间
}

type LogoutRecord struct {
	ID         int64      `bson:"_id" gorm:"AUTO_INCREMENT"`        //
	UserID     int32      `bson:"userid"`     //用户id
	Time  	   time.Time  `bson:"time"`  //登录时间
}

type RegisterRecord struct {
	ID         int64      `bson:"_id" gorm:"AUTO_INCREMENT"`        //
	UserID     int32      `bson:"userid"`     //用户id
	Channel    string     `bson:"channel"`         //登录ip
	Time  time.Time  `bson:"time"`  //登录时间
}

type OrderRecord struct {
	ID         string    `bson:"_id"`        //订单id
	UserID     int32     `bson:"userid"`     //用户id
	ProductID  int32     `bson:"productid"`  //充值商品id
	Amount     float64   `bson:"amount"`     //充值金额
	//State 	   int32  	 `bson:"state"`      //订单状态 0未支付  1已支付完成
	Time time.Time `bson:"createTime"` //充值时间
}


//---------------------------------日结算---------------------------------
//日活跃日志
type DayLoginRecord struct {
	ID string 	`bson:"_id"`        //2018-01-10
	Total int `bson:"total"`
}

//日新增
type DayRegisterRecord struct {
	ID string 	`bson:"_id"`        ////2018-01-10
	Total int `bson:"total"`
}

//日充值
type DayChargeRecord struct {
	ID string 	`bson:"_id"` //2018-01-10
	Total int `bson:"total"`
	TotalMoney float64 `bson:"totalMoney"`
}
