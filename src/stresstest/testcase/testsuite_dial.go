package testcase

import (
	"gok/service/msg/protocol"
	"log"
	"gok/constant"
	testsuit_conf "stresstest/conf"
	aliens_log "aliens/log"
	"aliens/common/util"
)

type DialTestSuite struct {
	*TestSuiteBase
	BaseID   	int32
	buildID  	int32
	IsTarget 	bool
	IsTask   	bool

	target  	*protocol.SelectEventTarget
	TargetWeightMapping map[int32]int32
	RandomTimes	int32
}

const (
	ADD_POWER		= 	0
	RANDOM_DIAL		= 	1
	RANDOM_TARGET	=	2
	SELECT_TARGET	=	3
	INTO_EVENT		=	4
	GET_FAITH		=	5
	GET_BELIEVER	=	6
	LOOT_FAITH		=	7
	ATK_BUILDING	=	8
	LOOT_BELIEVER	=	9
	DONE_EVENT		=	10
)

func NewDialTestSuite(name string, circle bool, messages []*protocol.C2GS) *DialTestSuite {
	targetWeightMapping := make(map[int32]int32)
	for index, weight := range testsuit_conf.DIAL.TargetSelectWeight {
		targetWeightMapping[int32(index)] = weight
	}
	result := &DialTestSuite{&TestSuiteBase{name:name, circle:circle, messages:[]*protocol.C2GS{}},0,0,false, false,&protocol.SelectEventTarget{},	targetWeightMapping, 0}

	result.messages = messages
	if result.messages == nil || len(result.messages) == 0 {
		log.Fatal("testcase can not be empty %v", result.name)
	}
	return result
}

func (this *DialTestSuite) NextMessage() *protocol.C2GS {
	msg := this.TestSuiteBase.NextMessage()
	if msg.GetAtkStarBuilding() != nil {
		msg.AtkStarBuilding.BuildingID = this.buildID
	}
	if msg.GetSelectEventTarget() != nil {
		msg.SelectEventTarget = this.target
	}
	return  msg
}

func (this *DialTestSuite) AcceptResult(ret *protocol.GS2C) {
	if ret.GetRandomDialRet() != nil {
		this.handlerRandomDialRet(ret.GetRandomDialRet())
	}
	if ret.GetRandomTargetRet() != nil {
		this.handlerRandomTargetRet(ret.GetRandomTargetRet())
	}
	if ret.GetSelectEventTargetRet() != nil {
		this.handlerSelectEventTargetRet(ret.GetSelectEventTargetRet())
	}
	if ret.GetIntoEventRet() != nil {
		this.handlerIntoEventRet(ret.GetIntoEventRet())
	}
	if ret.GetGetFaithRet() != nil {
		this.SetNextMessage(DONE_EVENT)
	}
	if ret.GetGetBelieverRet() != nil {
		this.SetNextMessage(DONE_EVENT)
	}
	if ret.GetLootFaithRet() != nil {
		this.SendFirstMessage()
	}
	if ret.GetAtkStarBuildingRet() != nil {
		this.SendFirstMessage()
	}
	if ret.GetLootBelieverRet() != nil {
		this.SendFirstMessage()
	}
}

func (this *DialTestSuite) SendFirstMessage() {
	if this.RandomTimes >= testsuit_conf.DIAL.Times {
		//log.Fatalf()
		aliens_log.Info("testCase Dial is over:random dial %v times", this.RandomTimes)
		//aliens_log.Fatal("testCase Dial is over:random dial %v times", this.RandomTimes)
		this.SetNextMessage(-1)
	}
	this.SetNextMessage(ADD_POWER)
}

func (this *DialTestSuite) handlerRandomDialRet(ret *protocol.RandomDialRet) {
	this.RandomTimes++
	if ret.Task == nil {
		this.SendFirstMessage()
	} else {
		this.BaseID = ret.Task.BaseID/10
		if this.BaseID == constant.TASK_ATK_BUILDING || this.BaseID == constant.TASK_ROB_BELIEVER || this.BaseID == constant.TASK_ROB_FAITH {
			this.SetNextMessage(RANDOM_TARGET)
		} else {
			this.SetNextMessage(INTO_EVENT)
		}
	}
}

func (this *DialTestSuite) handlerRandomTargetRet(ret *protocol.RandomTargetRet) {
	index := util.RandomWeight(this.TargetWeightMapping)
	target := ret.Targets[index]
	this.target.EventID = ret.EventID
	this.target.Nickname =target.Nickname
	this.target.TargetId =target.Id
}

func (this *DialTestSuite) handlerSelectEventTargetRet(ret *protocol.SelectEventTargetRet) {
	for _, build := range ret.StarInfo.Building {
		if build.Level > 0 {
			this.buildID = build.Type
			break
		}
	}
}

func (this *DialTestSuite) handlerIntoEventRet(ret *protocol.IntoEventRet) {
	if this.BaseID == constant.TASK_GET_FAITH {
		this.SetNextMessage(GET_FAITH)
	}
	if this.BaseID == constant.TASK_GET_BELIEVER {
		this.SetNextMessage(GET_BELIEVER)
	}
	if this.BaseID == constant.TASK_ROB_FAITH {
		this.SetNextMessage(LOOT_FAITH)
	}
	if this.BaseID == constant.TASK_ATK_BUILDING {
		this.SetNextMessage(ATK_BUILDING)
	}
	if this.BaseID == constant.TASK_ROB_BELIEVER {
		this.SetNextMessage(LOOT_BELIEVER)
	}
}


