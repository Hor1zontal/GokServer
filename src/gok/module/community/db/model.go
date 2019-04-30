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

//事件记录
type DEventRecord struct {
	ID          int32       `bson:"_id"`            //事件id
	Data        string 		`bson:"data"`           //事件记录数据  json格式

}

//朋友圈评论
type Comment struct {
	ID 		int32   	`bson:"_id"`		//评论id
	Uid		int32		`bson:"uid"`		//评论的玩家id
 	NickName        string 		`bson:"nickname"`	//评论的玩家昵称
	ReplyID 	int32		`bson:"replyid"`	//回复哪条评论  0代表不对哪条评论
	Msg 		string		`bson:"msg"`		//评论的内容
	CreateTime 	time.Time	`bson:"createtime"`	//创建时间
}

//朋友圈数据库映射对象
type DMoments struct {
	ID          string      `bson:"_id"        json:"_id"`        //朋友圈id
	Uid         int32		`bson:"uid"        json:"uid"`	      //发消息用户uid
	CreateTime  time.Time 	`bson:"createTime" json:"createTime"` //朋友圈消息的创建时间
	Data        string 		`bson:"data" json:"data"` //朋友圈消息数据
	//Comments    []*Comment	`bson:"comments"   json:"comments"`   //评论
	//Likes	    []int32		`bson:"likes"      json:"likes"`      //点赞的用户id列表
}

type MomentData struct {
	Type     int   `json:"type"`		//朋友圈消息类型 1 出售物品 2 物品被购买 3 切换下个星球 4 分享圣物 5 分享圣物组合
	RefID    int32   `json:"refID"`		//关联id
}

//组合ID
type CompID struct {
	SubID1 		int32			`bson:"s1"`   //id1
	SubID2		int32			`bson:"s2"`   //id2
}

//好友信息
type DFollow struct {
	ID      CompID    `bson:"_id"`     //好友
	AddTime time.Time `bson:"addtime"` //添加时间
}

//好友信息
//type DFriend struct {
//	ID      CompID    `bson:"_id"`     //好友
//	AddTime time.Time `bson:"addtime"` //添加时间
//}

//好友申请请求
//type DFriendReq struct {
//	ID        CompID      `bson:"_id"`     //好友
//	AddTime   time.Time   `bson:"addTime"`   //添加时间
//}