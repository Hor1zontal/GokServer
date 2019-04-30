package object
import (
	"stresstest/testcase"
	"gok/service/msg/protocol"
	"stresstest/message"
	test_conf "stresstest/conf"
	util2 "github.com/name5566/leaf/util"
)

func (this *Player) InitTestSuite() testcase.TestSuite {
	minBeliever := test_conf.INIT.UserInfo[this.accountType].BelieverRange[0]
	maxBeliever := test_conf.INIT.UserInfo[this.accountType].BelieverRange[1]
	minBuild := test_conf.INIT.UserInfo[this.accountType].BuildRange[0]
	maxBuild := test_conf.INIT.UserInfo[this.accountType].BuildRange[1]
	believerLevel := util2.RandInterval(minBeliever, maxBeliever)
	//fmt.Printf("======Level:%v BelieverRange:[%v,%v]\n", believerLevel, minBeliever, maxBeliever)
	buildLevel := util2.RandInterval(minBuild, maxBuild)
	//fmt.Printf("======Level:%v BuildRange:[%v,%v]\n", buildLevel, minBuild, maxBuild)
	initTestCases := []*protocol.C2GS{
		//BuildCreateMail(this.Uid, this.starType, "10000", "0","0","0","5","6","50"),
		//BuildGetAllMail(this.Uid),
		//BuildInternalAddAttach(this.Uid, this.starType, 1000, 0,0,0,false,true,6,5),

		//trade模块的初始化添加圣物
		//message.BuildAddItems(this.starType, 1,1),
		//message.BuildGetOnNotices(constant.NOTICE_ON),
		//message.BuildGetStarsSelect(),
		//message.BuildSelectStar(1),
		//message.BuildGetStarInfo(),
		//message.BuildGetStarsSelect(3),
		message.BuildSetBelievers(believerLevel),
		message.BuildSetBuildings(this.starID, buildLevel),
	}
	return testcase.NewTestSuite("init", false, initTestCases)
}
