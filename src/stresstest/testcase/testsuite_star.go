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

import "gok/service/msg/protocol"





type StarTestSuite struct {
	*TestSuiteBase
}

func (this *StarTestSuite) NextMessage() *protocol.C2GS {
	message := this.TestSuiteBase.NextMessage()


	//TODO 修改参数
	return message
}



func (this *StarTestSuite) AcceptResult(ret *protocol.GS2C)  {

}

