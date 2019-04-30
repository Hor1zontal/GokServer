/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/20
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package main

import (
	"github.com/name5566/leaf/log"
	"os"
	"os/signal"
)

func main() {

	//StarID               int32                 `protobuf:"varint,1,opt,name=starID,proto3" json:"starID,omitempty"`
	//Type                 int32                 `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	//OwnID                int32                 `protobuf:"varint,3,opt,name=ownID,proto3" json:"ownID,omitempty"`
	//Building             []*BuildingInfo       `protobuf:"bytes,4,rep,name=building" json:"building,omitempty"`
	//Believer             []*BelieverInfo       `protobuf:"bytes,5,rep,name=believer" json:"believer,omitempty"`
	//CreateTime           int64                 `protobuf:"varint,6,opt,name=createTime,proto3" json:"createTime,omitempty"`
	//DoneTime             int64                 `protobuf:"varint,7,opt,name=doneTime,proto3" json:"doneTime,omitempty"`
	//Seq                  int32                 `protobuf:"varint,8,opt,name=seq,proto3" json:"seq,omitempty"`
	//CivilizationLv       int32                 `protobuf:"varint,19,opt,name=civilizationLv,proto3" json:"civilizationLv,omitempty"`
	//CivilizationProgress int32                 `protobuf:"varint,20,opt,name=civilizationProgress,proto3" json:"civilizationProgress,omitempty"`
	//CivilizationReward   []*CivilizationReward `protobuf:"bytes,21,rep,name=civilizationReward" json:"civilizationReward,omitempty"`


	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	log.Release("kylin start")
	sig := <-c
	log.Release("kylin closing down (signal: %v)", sig)

	//starInfo := &protocol.GetStarInfoRet{
	//	CurrentStar:&protocol.StarInfoDetail{
	//		StarID:1,
	//		Type:1,
	//		OwnID:1,
	//	},
	//}
	//data := &protocol.GS2C{
	//	Sequence:[]int32{1234},
	//	GetStarInfoRet:starInfo,
	//}
	//
	//starInfoMarshal, _ := starInfo.Marshal()
	//dataMarshal, _ := data.Marshal()
	//log.Debug("%v-%v", len(starInfoMarshal), len(dataMarshal))
}
