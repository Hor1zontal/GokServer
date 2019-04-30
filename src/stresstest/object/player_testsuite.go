/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/31
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package object

import (
	"stresstest/testcase"
	"stresstest/conf"
)

func (this *Player) TestSuite() {
	if conf.INIT.Enable {
		this.AddTestSuite(this.InitTestSuite())
	}
	if conf.USER.Enable {
		this.AddTestSuite(this.UserTestSuite())
	}
	if conf.DIAL.Enable {
		this.AddTestSuite(this.DialTestSuite())
	}
	if conf.MAIL.Enable {
		this.AddTestSuite(this.MailTestSuite())
	}
	if conf.STAR.Enable {
		this.AddTestSuite(this.StarTestSuite())
	}
	if conf.COMMUNITY.Enable {
		this.AddTestSuite(this.CommunityTestSuite(2000))
	}
	if conf.TRADE.Enable {
		this.AddTestSuite(this.TradeTestSuite())
	}


	////accountNum为创建的账号数量

}

func (this *Player) AddTestSuite(suite testcase.TestSuite) {
	this.testSuites = append(this.testSuites, suite)
}
