package object

import (
	"stresstest/testcase"
	"gok/service/msg/protocol"
	"stresstest/message"
)

func (this *Player) MailTestSuite() testcase.TestSuite {
	mailTestCases := []*protocol.C2GS{
		message.BuildCreateMail(this.Uid, this.starType, "1","0","0","0","0","0","0"),
		message.BuildGetAllMail(this.Uid),
		message.BuildDrawMail(this.Uid, 0),
		//message.BuildRemoveMail(this.Uid, []int64{0}),
	}
	return testcase.NewMailTestSuite("mail", true, mailTestCases)
}


