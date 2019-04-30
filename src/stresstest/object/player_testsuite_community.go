package object

import (
	"stresstest/testcase"
	"gok/service/msg/protocol"
	"stresstest/message"
)

func (this *Player) CommunityTestSuite(accountNum int32) testcase.TestSuite {
	communityTestCases := []*protocol.C2GS{
		message.BuildFollow(this.Uid, this.Uid-accountNum),
		message.BuildGetFollowerList(this.Uid),
		message.BuildGetFollowingList(this.Uid),
		message.BuildUnFollow(this.Uid, this.Uid-accountNum),
		message.BuildPublicMoment(this.Uid),
		message.BuildGetReceiveMoments(this.Uid),
		message.BuildGetPublicMoments(this.Uid),
		//BuildRemoveMoment(this.Uid, []string{}),//删除朋友圈
	}
	return testcase.NewTestSuite("community", true, communityTestCases)
}

