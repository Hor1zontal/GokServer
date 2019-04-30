/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/4/17
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package order

import (
	"aliens/common"
	"aliens/log"
	"gok/module/center/cache"
)

//临时生成订单缓存，未支付完成的
//var tempOrders = make(map[string]*TempOrder)

//临时订单过期时间 60分钟
const tempOrderTimeout = 60 * 60

type TempOrder struct {
	ID string
	UserID int32
	ProductID int32
	Amount float64
	//CreateTime time.Time
}

func init() {

}

//生成临时订单
func GenTempOrder(uid int32, shopID int32, amount float64) string {
	orderId := util.Rand().Hex()
	cache.CenterCache.SetOrder(orderId, uid, shopID, amount, tempOrderTimeout)
	//order := &TempOrder{ID: orderId}
	//order.UserID = uid
	//order.ProductID = shopID
	//order.Amount = amount
	//order.CreateTime = time.Now()
	//tempOrders[orderId] = order
	log.Debug("gen temp order %v", orderId)
	return orderId
}

func GetTempOrder(orderID string) *TempOrder {
	uid := cache.CenterCache.GetOrderUID(orderID)
	if uid == 0 {
		return nil
	}
	productID := cache.CenterCache.GetOrderProductID(orderID)
	if productID == 0 {
		return nil
	}
	amount := cache.CenterCache.GetOrderAmount(orderID)
	return &TempOrder{ID:orderID, UserID:uid, ProductID:productID, Amount:amount}
}

func RemoveTempOrder(orderID string) {
	cache.CenterCache.RemoveOrder(orderID)
}


//处理订单超时
//func DealTimeout() {
//	now := time.Now()
//	for id, TempOrder := range tempOrders {
//		if now.Sub(TempOrder.CreateTime).Seconds() >= tempOrderTimeout{
//			delete(tempOrders, id)
//		}
//	}
//}
