package object

import (
	"stresstest/testcase"
	"gok/service/msg/protocol"
	"stresstest/message"
)

func (this *Player) UserTestSuite() testcase.TestSuite {
	userTestCases := []*protocol.C2GS{
		message.BuildGetAvatar(this.Uid),
		message.BuildChangeDesc("test123"),
		message.BuildDisplayInfo(0,10),
		message.BuildUpdateDisplay(this.Uid),
		message.BuildRankInfo(this.Uid, 1),
		message.BuildRoleFalgInfo(),
		message.BuildUpdateRoleFlag(1,1),
		message.BuildUpdatePower(),
		message.BuildBelieverFlagInfo(),
		message.BuildGetItem(this.Uid),
		message.BuildSearchUser(this.Uid),
		message.BuildDealList(this.Uid),
		message.BuildGlobalMessage(),
		//message.BuildPublicShare(1,"10000"),
		//message.BuildPublicWechatShare(),
		//message.BuildRequestItem("10201"),
		//message.BuildAcceptItem(""),
		//message.BuildRejectItem(""),
	}
	return testcase.NewTestSuite("user", true, userTestCases)
}