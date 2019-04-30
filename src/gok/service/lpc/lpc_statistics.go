/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package lpc

import (
	"gok/constant"
	"gok/module/statistics"
	"gok/module/statistics/model"
)

var StatisticsHandler = &statisticsHandler{}

type statisticsHandler struct {

}

func (this *statisticsHandler) AddServiceStatistic(service string, no int32, interval float64) {
	statistics.ChanRPC.Go(constant.INTERNAL_STATISTICS_SERVICE_CALL, service, no, interval)
}

func (this *statisticsHandler) AddOnlineStatistic(userCount int, visitorCount int) {
	statistics.ChanRPC.Go(constant.INTERNAL_STATISTICS_ONLINE, userCount, visitorCount)
}

func (this *statisticsHandler) AddStatisticData(data model.IStatisticData) {
	statistics.ChanRPC.Go(constant.STATISTICS_DATA, data)
}


