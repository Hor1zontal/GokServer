/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/31
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package testcase

import (
	"gok/service/msg/protocol"
	"log"
)


func NewTestSuite(name string, circle bool, messages []*protocol.C2GS) *TestSuiteBase {
	result := &TestSuiteBase{name:name, circle:circle, messages:[]*protocol.C2GS{}}
	result.messages = messages
	if result.messages == nil || len(result.messages) == 0 {
		log.Fatal("testcase can not be empty %v", result.name)
	}
	return result
}


type TestSuite interface {
	NextMessage() *protocol.C2GS
	AcceptResult(ret *protocol.GS2C)
}


type TestSuiteBase struct {
	name string

	circle bool

	messages []*protocol.C2GS

	index int

	nextSeq int32

	session int32

	end bool
}

//设置下一条需要测试的消息编号
func (this *TestSuiteBase) SetNextMessage(seq int32) {
	if seq == -1 {
		this.end = true
	} else {
		this.nextSeq = seq
	}
}

func (this *TestSuiteBase) NextMessage() *protocol.C2GS {
	if this.end {
		return nil
	}
	if this.nextSeq >= 0 {
		this.index = int(this.nextSeq)
		this.nextSeq = -1
	}
	if this.index >= len(this.messages) {
		if !this.circle {
			return nil
		}
		this.index = 0
	}
	message := this.messages[this.index]
	//fmt.Printf("%d\n",this.index)
	this.session ++
	this.index ++
	message.Session = this.session

	return message
}


func (this *TestSuiteBase) AcceptResult(ret *protocol.GS2C)  {

}
