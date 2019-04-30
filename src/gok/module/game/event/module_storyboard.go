/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/7/12
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package event

import (
	"gok/module/game/conf"
	"encoding/json"
	"gok/service/msg/protocol"
	"gok/module/game/db"
)

type StoryboardModule struct{

}


//从持久化数据源中初始化数据
func (this *StoryboardModule)Init(data []byte, config *conf.StepData) {
	if (data == nil) {
		return
	}
	json.Unmarshal(data, &this)
}

func (this *StoryboardModule)Start(context db.EventContext)  {

}

//处理时间限制到达
func (this *StoryboardModule)HandleTimesUp(context db.EventContext) {
	context.NextStep()
}

//处理消息请求
func (this *StoryboardModule)HandleMessage(request *protocol.C2GS, response *protocol.GS2C, context db.EventContext) {
	//处理
	if (request.GetDoneEventStep() != nil) {
		context.NextStep()
	}
}

