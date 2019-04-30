package object

import (
	"stresstest/testcase"
	"gok/service/msg/protocol"
	"stresstest/message"
)

func (this *Player) DialTestSuite() testcase.TestSuite {
	dialTestCases := []*protocol.C2GS{
		message.BuildAddPower(1),
		message.BuildRandomDial(),
		message.BuildRandomTarget(this.Uid),
		message.BuildSelectEventTarget(this.Uid),
		message.BuildIntoEvent(this.Uid),
		message.BuildGetFaith(this.starType),
		message.BuildGetBeliever(this.Uid),
		message.BuildLootFaith(this.Uid),
		message.BuildAtkBuilding(this.Uid),
		message.BuildLootBeliever(this.Uid),
		message.BuildDoneEventStep(this.Uid),
	}
	return testcase.NewDialTestSuite("dial", true, dialTestCases)
}
