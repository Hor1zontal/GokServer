/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/7/28
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package util

import (
	"gok/module/star/db"
	"gok/service/msg/protocol"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func BuildBelieverInfo(believer *db.DBBeliever) *protocol.BelieverInfo {
	return &protocol.BelieverInfo{
		Id:believer.ID,
		Num:believer.Num,
	}
}

//func BuildStarInfo(star *db.DBStar) *protocol.StarInfo{
//	buildingInfos := BuildStarBuildingInfo(star)
//	believerInfos := BuildStarBelieverInfo(star)
//	return &protocol.StarInfo{
//		starID:star.ID),
//		Type:star.Type),
//		Building:buildingInfos,
//		Believer:believerInfos,
//	}
//}

func BuildBuildingInfoPush(uid int32, starID int32, starType int32, building *db.DBBuilding) *protocol.GS2C {
	message := &protocol.BuildingInfoPush{
		Uid:  uid,//拥有者的用户ID
		StarID:  starID,
		Type:  starType,
		Building:[]*protocol.BuildingInfo{building.BuildProtocol()},
	}
	return &protocol.GS2C{
		Sequence:[]int32{1031},
		BuildingInfoPush:message,
	}
}




//func BuildBuildingBuffInfoPush(buildingType int32, buff map[int32]int32) *protocol.GS2C {
//	buffIDArray := []int32{}
//	buffNumArray := []int32{}
//	for buffID, buffNum := range buff {
//		buffIDArray = append(buffIDArray, buffID)
//		buffNumArray = append(buffNumArray, buffNum)
//	}
//	message := &protocol.BuildingBuffInfo{
//		Type:  buildingType,
//		BuffID:buffIDArray,
//		BuffNum:buffNumArray,
//	}
//	return &protocol.GS2C{
//		Sequence:[]int32{1057},
//		BuildingBuffPush:message,
//	}
//}

//func BuildStarBuildingInfo(star *db.DBStar)*protocol.StarInfoDetail{//写入星球详细数据  wjl 20170606
//	if star == nil{
//		return nil;
//	}
//	building := BuildStarBuildingInfo(star)
//	believers := BuildStarBelieverInfo(star)
//	message := &protocol.StarInfoDetail{
//		starID:  star.ID ),
//		Type:  star.Type ),
//		UserID:  star.Owner ),//拥有者的用户ID
//		OwnNikeName:""),
//		Building:building,
//		Believer:believers,
//	}
//	return message;
//}

//func BuildRoleCivilizationPush(star *db.DBStar)  *protocol.GS2C {
//	return &protocol.GS2C{
//		Sequence: []int32{1006},
//		CivilizationPush: &protocol.CivilizationPush{
//			CivilizationLv:       star.CivilizationLevel,
//			CivilizationProgress: star.CivilizationValue,
//			StarID:               star.ID,
//			StarSeq:              star.Seq,
//		}}
//}

func BuildCivilizationInfo(star *db.DBStar)  *protocol.CivilizationInfo {
	return &protocol.CivilizationInfo{
			CivilizationLv:       star.CivilizationLevel,
			CivilizationProgress: star.CivilizationValue,
			StarID:               star.ID,
			StarSeq:              star.Seq,
		}
}

//func GetBuildingTotalLevel(buildings []*db.DBBuilding) int32 {
//	var result int32 = 0
//	if buildings == nil {
//		return result
//	}
//	for _, building := range buildings {
//		result += building.Level
//	}
//	return result
//}
//
//func BuildCivilizationReward(star *db.DBStar)[]*protocol.CivilizationReward{
//	infos := []*protocol.CivilizationReward{};
//	if star != nil{
//		for _, reward := range star.CivilizationReward {
//			infos = append(infos, &protocol.CivilizationReward{Level:reward.Level, Draw:reward.Draw})
//		}
//	}
//	return infos;
//}
//
//func BuildStarBelieverInfo(star *db.DBStar)[]*protocol.BelieverInfo{
//	infos := []*protocol.BelieverInfo{};
//	if star != nil{
//		for _, believer := range star.Believer {
//			infos = append(infos, BuildBelieverInfo(believer))
//		}
//	}
//	return infos;
//}
//
//func BuildStarBuildingInfo(star *db.DBStar)[]*protocol.BuildingInfo{ //写入星球的建筑物信息
//	infos := []*protocol.BuildingInfo{}
//	if star != nil{
//		infos = BuildStarBuildingInfos(star.Building)
//	}
//	return infos
//}

func BuildStarBuildingInfos(buildings []*db.DBBuilding)[]*protocol.BuildingInfo{//写入星球的建筑物信息
	infos := []*protocol.BuildingInfo{}
	for _, building := range buildings {
		infos = append(infos, building.BuildProtocol())
	}
	return infos
}

func BuildCorrectCivilRewardMessage(power int32) *protocol.C2GS {
	return &protocol.C2GS{
		Sequence:[]int32{510},
		CorrectCivilReward: &protocol.CorrectCivilReward{
			Power:power,
		},
	}
}

func BuildUpdateFlag(uid int32, flags []*protocol.FlagInfo) *protocol.C2GS {
	return &protocol.C2GS{
		Sequence:[]int32{511},
		UpdateUnlockFlag:&protocol.UpdateUnlockFlag{
			Uid:uid,
			Flags:flags,
		},
	}
}

func BuildNewsFeedMessage1(uid int32, newsFeedType int32, param1 int32, param2 int32, param3 int32, ext []string) *protocol.C2GS {
	return &protocol.C2GS{
		Sequence:[]int32{523},
		AddNewsFeed: BuildNewsFeed1(uid, newsFeedType, param1, param2, param3, ext),
	}
}


func BuildNewsFeed1(uid int32, newsFeedType int32, param1 int32, param2 int32, param3 int32, ext []string) *protocol.NewsFeed {
	return &protocol.NewsFeed{
		Id:       bson.NewObjectId().Hex(),
		RelateID: uid,
		//RelateNickname:nickname),
		//RelateAvatar:avatar),
		Type:   newsFeedType,
		Time:   time.Now().Unix(),
		Param1: param1,
		Param2: param2,
		Param3: param3,
		Ext:    ext,
	}
}

