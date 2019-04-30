/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/4/18
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package conf

import (
	"encoding/json"
	"gok/module/cluster/center"
)


func Close() {
}


func Init() {

	center.ConfigCenter.SubscribeConfig("shopbase", updateShopBaseData)
	center.ConfigCenter.SubscribeConfig("gamebase", updateGameBaseData)
	center.ConfigCenter.SubscribeConfig("wechatresponse", updateWechatResponseData)
	center.ConfigCenter.SubscribeConfig("maintain", updateMaintainBase)
}



var DATA struct {
	*GameBase
	ShopBaseData   map[int32]*ShopBase // wjl 20170619 游戏的商店基础数据
	WechatResponseData map[string]string
}


type GameBase struct {
	WechatExpireTime  []int
	WechatExpireTime1 int
	WechatExpireTime2 int
}

//商店基础配置 wjl 20170619
type ShopBase struct{
	ID int32 //商品的基本ID
	Type int32 //商品类型 // 1 购买钻石 2 月卡 3 购买法力
	Amount int32 //商品获得的数值
	MoneyType int32 //消耗的类型 0x01 信仰值 0x02 钻石 0x03 现金
	Value float64 //消耗需要的数额
}

type WechatResponse struct {
	Key string
	Value string
}


func updateGameBaseData( data []byte ){
	var datas  = &GameBase{}
	json.Unmarshal(data, datas)
	if datas.WechatExpireTime != nil && len(datas.WechatExpireTime) == 2 {
		datas.WechatExpireTime1 = datas.WechatExpireTime[0]
		datas.WechatExpireTime2 = datas.WechatExpireTime[1]
	}

	if datas.WechatExpireTime1 == 0 {
		datas.WechatExpireTime1 = 24 * 60 * 60
	}
	if datas.WechatExpireTime2 == 0 {
		datas.WechatExpireTime2 = 46 * 60 * 60
	}

	DATA.GameBase = datas
}


func updateWechatResponseData( data []byte ){
	var datas []*WechatResponse
	json.Unmarshal(data, &datas)
	results := make(map[string]string)
	for _, data := range datas {
		results[data.Key] = data.Value
	}
	DATA.WechatResponseData = results
}

func updateShopBaseData( data []byte ){//wjl 20170619 通过zk 获取游戏商店数据
	var datas []*ShopBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]*ShopBase)
	for _, data := range datas {
		results[data.ID] = data
	}
	DATA.ShopBaseData = results
}
