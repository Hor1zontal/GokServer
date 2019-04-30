/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package cache

const (
	ORDER_KEY_REPFIX string = "order:"   			//用户内存所在的节点信息
	ORDER_PROP_UID = "uid"		  //用户id
	ORDER_PROP_PRODUCTID = "pid"  //商品id
	ORDER_PROP_AMOUNT = "amount"  //价格
)


type CenterCacheManager struct {
	*cacheManager
}

func NewCenterCacheManager() *CenterCacheManager {
	return &CenterCacheManager{
		&cacheManager{},
	}
}

//设置用户会话所在的服务节点
func (this *CenterCacheManager) SetOrder(orderID string, uid int32, productID int32, amount float64, expire int) {
	key := ORDER_KEY_REPFIX + orderID
	props := make(map[interface{}]interface{}, 3)
	props[ORDER_PROP_UID] = uid
	props[ORDER_PROP_PRODUCTID] = uid
	props[ORDER_PROP_AMOUNT] = amount
	this.redisClient.HMSet(key, props)
	this.redisClient.Expire(key, expire)
}

//获取用户会话所在的服务节点
func (this *CenterCacheManager) GetOrderUID(orderID string) int32 {
	key := ORDER_KEY_REPFIX + orderID
	return this.redisClient.HGetInt32(key, ORDER_PROP_UID)
}

func (this *CenterCacheManager) GetOrderProductID(orderID string) int32 {
	key := ORDER_KEY_REPFIX + orderID
	return this.redisClient.HGetInt32(key, ORDER_PROP_PRODUCTID)
}

func (this *CenterCacheManager) GetOrderAmount(orderID string) float64 {
	key := ORDER_KEY_REPFIX + orderID
	return this.redisClient.HGetFloat64(key, ORDER_PROP_UID)
}

func (this *CenterCacheManager) RemoveOrder(orderID string) bool {
	key := ORDER_KEY_REPFIX + orderID
	return this.redisClient.DelData(key)
}

