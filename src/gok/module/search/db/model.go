/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/12
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package db

import "time"

//星球检索信息
type StarBuildingInfo struct {
	ID          int32        	  `bson:"_id"`        //星球id
	TotalLevel  int32      		  `bson:"total"`        //建筑总等级
	UpdateTime  time.Time		  `bson:"updateTime"`        //建筑总等级
}

type StarBelieverInfo struct {
	ID         int32     `bson:"_id"`        //星球id
	TotalLevel int32     `bson:"total"`        //信徒总等级
	UpdateTime time.Time `bson:"updateTime"`        //建筑总等级
}

type UserIndex struct {
	Uid int32    		 `bson:"_id"`
	StarType  int32 	 `bson:"starType"`//星球类型
	ActiveTime time.Time `bson:"activeTime"`//上一次的活跃时间
	ReceiveTime time.Time `bson:"receiveTime"`//收到圣物求助的时间
	//ConditionID int `bson:"-"`//当前所处的索引id
}
