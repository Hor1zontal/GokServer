package internal

import (
	"gok/constant"
	"github.com/name5566/leaf/timer"
	"aliens/log"
	"gok/module/statistics/model"
	"gok/module/cluster/center"
	"github.com/sirupsen/logrus"
	"gok/module/statistics/elastics"
	"gok/module/statistics/conf"
)


var esHandler = elastics.NewESHandler(conf.Config.ES.Name)

//服务调用统计信息、一分钟一次
var serviceStatistics = make(map[string]map[int32]*model.CallInfo) //服务名 - 服务编号 - 调用信息

var serviceFields = logrus.Fields{}
var onlineFields = logrus.Fields{}
var dialFields = logrus.Fields{}

func init() {
	skeleton.RegisterChanRPC(constant.INTERNAL_STATISTICS_SERVICE_CALL, handleServiceStatic)
	skeleton.RegisterChanRPC(constant.INTERNAL_STATISTICS_ONLINE, handleOnlineStatic)
	skeleton.RegisterChanRPC(constant.STATISTICS_DATA, handleCustomStatic)


	cron, err := timer.NewCronExpr("*/1 * * * *")
	if err != nil {
		log.Error("init service statistics timer error : %v", err)
	}

	//每天凌晨12点执行一次
	//dayCron, err := timer.NewCronExpr("0 0 * * *")
	//if err != nil {
	//	log.Error("init dump timer error : %v", err)
	//}
	//skeleton.CronFunc(dayCron, esHandler.UpdateDayPrefix)
	skeleton.CronFunc(cron, handleTimer)
}

func handleCustomStatic(args []interface{}) {
	data, ok := args[0].(model.IStatisticData)
	if ok {
		esHandler.HandleDayESLog(data.GetName(), "", data.GetData())
	}
}

func handleOnlineStatic(args []interface{}) {
	userCount := args[0].(int)   //用户数量
	visitorCount := args[1].(int)   //空连接数量
	onlineFields["node"] = center.GetServerNode()
	onlineFields["u_count"] = userCount
	onlineFields["v_count"] = visitorCount
	esHandler.HandleDayESLog(constant.ES_LOG_ONLINE, "", onlineFields)
}

//处理服务信息统计
func handleServiceStatic(args []interface{}) {
	service := args[0].(string)   //服务名称
	serviceNo := args[1].(int32)   //服务处理编号
	interval := args[2].(float64)  //服务处理时间间隔
	callInfos := serviceStatistics[service]
	if callInfos == nil {
		callInfos = make(map[int32]*model.CallInfo)
		serviceStatistics[service] = callInfos
	}

	callInfo := callInfos[serviceNo]
	if callInfo == nil {
		callInfo = &model.CallInfo{}
		callInfos[serviceNo] = callInfo
	}
	callInfo.AddCall(interval)
}

func handleTimer() {
	for service, callInfos := range serviceStatistics {
		for serviceNo, callInfo := range callInfos {
			serviceFields["service"] = service
			serviceFields["node"] = center.GetServerNode()
			serviceFields["no"] = serviceNo
			result, count, avg := callInfo.DumpData()
			if !result {
				continue
			}
			serviceFields["count"] = count
			serviceFields["avg"] = avg
			esHandler.HandleDayESLog(constant.ES_LOG_SERVICE, "", serviceFields)
		}
	}
}