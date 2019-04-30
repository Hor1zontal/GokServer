/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/9/4
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package internal

import (
	"github.com/name5566/leaf/timer"
	"aliens/log"
	"gok/module/search/core"
)

func init() {
	//每天凌晨5点清除不活跃的玩家索引数据
	cron, err := timer.NewCronExpr("0 5 * * *")
	if err != nil {
		log.Error("init searcher timer error : %v", err)
	}


	skeleton.CronFunc(cron, core.StarSearcher.Clean)


	//每5分钟更新一下搜索过滤列表
	filterCron, err := timer.NewCronExpr("*/5 * * * *")
	if err != nil {
		log.Error("init searcher filter expire timer error : %v", err)
	}

	skeleton.CronFunc(filterCron, core.StarSearcher.DealFilterExpire)

}