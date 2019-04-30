/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/10/9
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package internal

import (
	"github.com/name5566/leaf/timer"
	"aliens/log"
	"gok/module/passport/wx"
)

func Init() {
	cron, err := timer.NewCronExpr("0 */1 * * *")
	if err != nil {
		log.Error("init wechat accessToken timer error : %v", err)
	}
	//刷新token
	skeleton.CronFunc(cron, wx.RefreshToken)
}

