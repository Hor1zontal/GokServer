/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/31
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package object

import (
	"gok/service/msg/protocol"
	"stresstest/testcase"
	"stresstest/message"
)

func (this *Player) StarTestSuite() testcase.TestSuite {



	starTestCases := []*protocol.C2GS{
		//BuildGetStarInfo(),
		message.BuildAddAttach(this.starType, 10000, 0, 0 , 0, true,6,10, true, 6, 90),
		message.BuildGetStarShield(),
		message.BuildGetStarStatistics(this.Uid, this.starID),
		message.BuildGetStarHistory(this.starID),
		message.BuildAddBelievers(this.starType, 1,2),
		message.BuildUpgradeBeliever(this.Uid, this.starType),
		message.BuildGetBuildingFaith(this.Uid,1),
		message.BuildGetItemGroup(this.Uid,this.starType),
	}

	//建筑等级 4个 lv4 1个 lv5
	for i := 1; i <= 5; i++ {//level
		for j := 1; j<= 5; j++ {//buildingType
			starTestCases = append(starTestCases, message.BuildBuildStarBuilding(this.starID,int32(j),int32(i)))
			starTestCases = append(starTestCases, message.BuildAccUpdateStarBuild(this.Uid, int32(j), this.starType, int32(i)))
			starTestCases = append(starTestCases, message.BuildUpdateStarBuildEnd(this.Uid, int32(j)))
			if i == 5 && j == 1 {
				break
			}
		}
	}
	//领取文明度奖励 1~4阶的文明度奖励
	//for i := 0; i < 4; i++ {
	//	starTestCases = append(starTestCases, message.BuildDrawCivilizationReward(this.starID, int32(i)))
	//}

	//试圣物组合
	for i := 1; i <= 10; i++ {
		starTestCases = append(starTestCases, message.BuildActiveGroup(this.Uid, this.starType, int32(i)))
	}

	//starTestCases = append(starTestCases, message.BuildSetBuildings(this.starID, 0))

	return testcase.NewTestSuite("star", false, starTestCases)
}