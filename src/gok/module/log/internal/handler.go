package internal

import (
	"gok/module/log/db"
	"gok/constant"
	"gok/module/log/cache"
	"github.com/name5566/leaf/timer"
	"aliens/log"
	"time"
)

var dayTimeStr string = time.Now().Format("2006-01-02")

func init() {
	// 向当前模块注册客户端发送的消息处理函数 handleMessage
	skeleton.RegisterChanRPC(constant.LOG_COMMAND, handle)

	//每天凌晨12点执行一次
	cron, err := timer.NewCronExpr("0 0 * * *")
	if err != nil {
		log.Error("init dump timer error : %v", err)
	}
	skeleton.CronFunc(cron, dayDump)
}

func dayDump() {
	dayTimeStr = time.Now().Add(time.Hour).Format("2006-01-02")

	dumpTime := time.Now().Add(-time.Hour)
	dumpTimeStr := dumpTime.Format("2006-01-02")
	if cache.LogCache.SetNX("dump" + dumpTimeStr, 1) {
		log.Debug("dump day log %v start....", dumpTime)
		db.DatabaseHandler.Insert(&db.DayChargeRecord{ID:dumpTimeStr, Total:cache.LogCache.GetDayCharge(dumpTime)})
		db.DatabaseHandler.Insert(&db.DayRegisterRecord{ID:dumpTimeStr, Total:cache.LogCache.GetDayRegister(dumpTime)})
		db.DatabaseHandler.Insert(&db.DayLoginRecord{ID:dumpTimeStr, Total:cache.LogCache.GetDayLogin(dumpTime)})
		log.Debug("dump day log %v end....", dumpTime)
	}
}

func handle(args []interface{}) {
	dbLog := args[0]
	switch dbLog.(type) {
		case *db.LoginRecord:
			record := dbLog.(*db.LoginRecord)
			cache.LogCache.IncrDayLogin(record.UserID, record.LoginTime)
			break
		case *db.RegisterRecord:
			cache.LogCache.IncrDayRegister(dbLog.(*db.RegisterRecord).Time)
			break
		case *db.OrderRecord:
			order := dbLog.(*db.OrderRecord)
			cache.LogCache.IncrDayCharge(order.UserID, order.Time)
			cache.LogCache.IncrByDayChargeFee(order.Time, order.Amount)
			break
	}
	err := db.DatabaseHandler.Insert(dbLog)
	if err != nil {
		log.Debug("insert log error %v", err)
	}
}