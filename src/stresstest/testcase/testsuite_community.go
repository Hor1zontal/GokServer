package testcase

import (
	"gok/service/msg/protocol"
	"log"
)

type CommunityTestSuite struct {
	*TestSuiteBase
	momentsID	string
}

func NewCommunityTestSuite(name string, circle bool, messages []*protocol.C2GS) *CommunityTestSuite {
	result := &CommunityTestSuite{&TestSuiteBase{name:name, circle:circle, messages:[]*protocol.C2GS{}}, ""}
	result.messages = messages
	if result.messages == nil || len(result.messages) == 0 {
		log.Fatal("testcase can not be empty %v", result.name)
	}
	return result
}

func (this *CommunityTestSuite) NextMessage() *protocol.C2GS {
	message := this.TestSuiteBase.NextMessage()
	if message.GetRemoveMoments() != nil {
		//this.mailID = append(this.mailID[:0],this.mailID[len(this.mailID):]...)
		message.RemoveMoments.MomentsID = append(message.RemoveMoments.MomentsID,  this.momentsID)
	}
	//TODO 修改参数
	return message
}

func (this *CommunityTestSuite) AcceptResult(ret *protocol.GS2C)  {
	if ret.GetPublicMomentRet() != nil {
		this.handlerPublicMomentRet(ret.GetPublicMomentRet())
	}
}

func (this *CommunityTestSuite) handlerPublicMomentRet(ret *protocol.PublicMomentRet) {
	this.momentsID = ret.MomentInfo.Id
}

