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
	"encoding/json"
	"gok/module/game/db"
	"gok/service/msg/protocol"
	"gok/module/game/conf"
)

type CustomDataModule struct{
	initial *conf.DataCustom  `bson:"-"`
	Data    []*conf.CustomData   `bson:"data"`
}

//从持久化数据源中初始化数据
func (this *CustomDataModule)Init(data []byte, config *conf.StepData) {
	if (config != nil) {
		conf, ok := config.Data.(*conf.DataCustom)
		if (ok) {
			this.initial = conf
		}
	}

	if (data == nil) {
		return
	}
	json.Unmarshal(data, &this)
}

func (this *CustomDataModule)Start(context db.EventContext)  {
	//初始化数据
	if (this.initial != nil) {
		this.Data = this.initial.CopyData()
	}
}

//处理时间限制到达
func (this *CustomDataModule)HandleTimesUp(context db.EventContext) {
	context.NextStep()
}

//处理消息请求
func (this *CustomDataModule)HandleMessage(message *protocol.C2GS, response *protocol.GS2C, context db.EventContext) {
	//处理数据保存
	if (message.GetSaveData() != nil) {
		data := message.GetSaveData().GetData()
		if (data != nil) {
			for _, saveData := range data {
				this.updateData(saveData.GetKey(), saveData.GetValue())

			}
		}
	} else if (message.GetDoneEventStep() != nil) {
		context.NextStep()
	}
}

func (this *CustomDataModule) updateData(key string, value int32)  {
	//初始化数据
	if (this.Data == nil) {
		this.Data = []*conf.CustomData{}
	}

	for _, data := range this.Data {
		if data.Key == key {
			data.Value = value
		}
		return
	}
	this.Data = append(this.Data, &conf.CustomData{key, value})
}
