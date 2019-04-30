/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/9
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package internal

import (
	"github.com/name5566/leaf/timer"
	"github.com/name5566/leaf/log"
	"gok/module/center/notice"
)

func init() {
	cron, err := timer.NewCronExpr("*/1 * * * *")
	if err != nil {
		log.Error("init order timer error : %v", err)
	}

	skeleton.CronFunc(cron, notice.NoticesManager.DealTimeout)
}