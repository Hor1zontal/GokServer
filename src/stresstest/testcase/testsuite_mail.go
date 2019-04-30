package testcase

import (
	"gok/service/msg/protocol"
	"log"
)

type MailTestSuite struct {
	*TestSuiteBase
	mailID	[]int64
}

func NewMailTestSuite(name string, circle bool, messages []*protocol.C2GS) *MailTestSuite {
	result := &MailTestSuite{&TestSuiteBase{name:name, circle:circle, messages:[]*protocol.C2GS{}}, []int64{}}
	result.messages = messages
	if result.messages == nil || len(result.messages) == 0 {
		log.Fatal("testcase can not be empty %v", result.name)
	}
	return result
}

func (this *MailTestSuite) NextMessage() *protocol.C2GS {
	message := this.TestSuiteBase.NextMessage()
	if message.GetRemoveMail() != nil {
		//this.mailID = append(this.mailID[:0],this.mailID[len(this.mailID):]...)
		message.RemoveMail.MailID = this.mailID
	}
	if message.GetDrawMail() != nil {
		for _, id := range this.mailID {
			message.DrawMail.MailID = id
		}
	}
	//TODO 修改参数
	return message
}

func (this *MailTestSuite) AcceptResult(ret *protocol.GS2C)  {
	if ret.GetGetAllMailRet() != nil {
		this.handlerGetAllMailRet(ret.GetGetAllMailRet())
	}
}

func (this *MailTestSuite) handlerGetAllMailRet(ret *protocol.GetAllMailRet) {
	//this.mailID = []int64{}
	this.mailID = append(this.mailID[:0],this.mailID[len(this.mailID):]...)
	for _, mail := range ret.Mail {
		this.mailID = append(this.mailID, mail.Id)
	}
}

