package internal

import (
	"github.com/name5566/leaf/timer"
	"aliens/log"
	"gok/module/star/session"
)


func init() {
	// 向当前模块注册客户端发送的消息处理函数 handleMessage
	cron, err := timer.NewCronExpr("*/1 * * * *")
	if err != nil {
		log.Error("init star timer error : %v", err)
	}

	skeleton.CronFunc(cron, session.StarManager.DealStarTimer)
}
