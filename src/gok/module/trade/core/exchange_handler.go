package core

import (
	"time"
	"gok/service/exception"
	"gok/module/trade/db"
	"gok/service/msg/protocol"
	"aliens/common/character"
	"gok/service/lpc"
	"gok/module/trade/cache"
)

var Exchanger = GetExchangeManager()

func GetExchangeManager() *ExchangeManager {
	return &ExchangeManager{
		//sales:make(map[int32]*protocol.Sale),
		//goods:make(map[string]*protocol.Goods),
		//userGoods:make(map[int32][]*protocol.Goods),
	}
}

type ExchangeManager struct {
	//sync.RWMutex

	//sales map[int32]*protocol.Sale  //拍卖的物品

	//goods map[string]*protocol.Goods //圣物货架的物品

	//userGoods map[int32][]*protocol.Goods //
}

//func (this *ExchangeManager) Init() {
//	//this.Lock()
//	//defer this.Unlock()
//	var sales []*db.UserSale
//	db.DatabaseHandler.QueryAll(&db.UserSale{}, &sales)
//	for _, sale := range sales {
//		this.sales[sale.UserID] = sale.BuildProtocol()
//	}
//
//	var goods []*db.UserGoods
//	db.DatabaseHandler.QueryAll(&db.UserGoods{}, &goods)
//	for _, good := range goods {
//		goodsID := getGoodsID(good.GetUserID(), good.GetItemID())
//		this.addGoods(good.GetUserID(), good.GetItemID(), goodsID, good.BuildProtocol())
//	}
//}

//发布商品
func (this *ExchangeManager) PublicSale(uid int32, itemID int32) *protocol.Sale {
	//this.Lock()
	//defer this.Unlock()
	sale := this.GetSale(uid)

	if sale != nil {
		exception.GameException(exception.ITEM_PUBLIC_REPEAT)
	}
	dbSale := &db.UserSale{UserID:uid, ItemID:itemID, PublicTime:time.Now()}
	saleInfo := dbSale.BuildProtocol()
	if cache.TradeCache.SetUserPublicSale(uid, saleInfo) {
		lpc.DBServiceProxy.ForceUpdate(dbSale, db.DatabaseHandler)
	}
	return saleInfo
}

//移除商品
func (this *ExchangeManager) RemoveSale(uid int32, itemID int32) *protocol.Sale {
	//this.Lock()
	//defer this.Unlock()
	sale := this.GetSale(uid)
	if sale == nil || sale.GetItemID() != itemID {
		exception.GameException(exception.PUBLIC_ITEM_NOTFOUND)
	}
	if cache.TradeCache.RemoveUserPublicSale(uid) {
		lpc.DBServiceProxy.Delete(&db.UserSale{UserID:uid}, db.DatabaseHandler)
	}
	//delete(this.sales, uid)
	return sale
}

func (this *ExchangeManager) GetSale(uid int32) *protocol.Sale {
	//this.RLock()
	//defer this.RUnlock()
	return cache.TradeCache.GetUserPublicSale(uid)
}

func getGoodsID(uid int32, itemID int32) string {
	if uid <= 0 || itemID <= 0 {
		exception.GameException(exception.INVALID_PARAM)
	}
	return character.Int32ToString(uid) + "_" + character.Int32ToString(itemID)
}

func (this *ExchangeManager) GetGoodsInfo(uid int32) []*protocol.Goods {
	//this.RLock()
	//defer this.RUnlock()
	return cache.TradeCache.HGetCacheGoods(uid)
}

func (this *ExchangeManager) BuyGoods(uid int32, itemID int32, num int32, social int32) (int32, int32) {
	//this.Lock()
	//defer this.Unlock()
	//goodsID := getGoodsID(uid, itemID)
	//goods := this.goods[goodsID]
	goods := cache.TradeCache.HGetCacheGood(uid, itemID)
	if goods == nil {
		exception.GameException(exception.GOODS_NOT_FOUND)
	}
	needSocial := goods.GetPrice() * num
	if needSocial > social {
		exception.GameException(exception.GAYPOINT_NOT_ENOUGH)
	}
	remain := goods.GetNum() - num
	if remain < 0 {
		exception.GameException(exception.GOODS_NOT_ENOUGH)
	} else if remain == 0 {
		//用完要删除对象
		//this.removeGoods(uid, itemID, goodsID)
		if cache.TradeCache.HRemoveCacheGood(uid, itemID) {
			lpc.DBServiceProxy.Delete(&db.UserGoods{CompID:db.CompID{uid, itemID}}, db.DatabaseHandler)
		}
		//db.DatabaseHandler.DeleteOne()
	} else {
		if cache.TradeCache.HSetCacheGood(uid, itemID, goods) {
			lpc.DBServiceProxy.Update(buildGoods(uid, goods), db.DatabaseHandler)
		}

		//db.DatabaseHandler.UpdateOne()
	}
	return remain, needSocial
}

func (this *ExchangeManager) PublicGoods(uid int32, goods *protocol.Goods) bool {
	//this.Lock()
	//defer this.Unlock()
	//goodsID := getGoodsID(uid, goods.GetId())

	//goodsCache := this.goods[goodsID]
	goodsCache := cache.TradeCache.HGetCacheGood(uid, goods.GetId())
	if goodsCache != nil {
		exception.GameException(exception.GOODS_ALREADY_EXIST)
	}
	//this.addGoods(uid, goods.GetId(), goodsID, goods)
	if cache.TradeCache.HSetCacheGood(uid, goods.GetId(), goods) {
		lpc.DBServiceProxy.Insert(buildGoods(uid, goods), db.DatabaseHandler)
	}
	return true
}

func buildGoods(uid int32, goods *protocol.Goods) *db.UserGoods {
	return &db.UserGoods{CompID:db.CompID{uid, goods.GetId()}, Num:goods.GetNum(), Price:goods.GetPrice()}
}


//撤销货物
func (this *ExchangeManager) CancelGoods(uid int32, itemID int32) *protocol.Goods {
	//this.Lock()
	//defer this.Unlock()
	//goodsID := getGoodsID(uid, itemID)

	goodsCache := cache.TradeCache.HGetCacheGood(uid, itemID)
	if goodsCache == nil {
		exception.GameException(exception.GOODS_NOT_FOUND)
	}
	//this.removeGoods(uid, itemID, goodsID)
	if cache.TradeCache.HRemoveCacheGood(uid, itemID) {
		lpc.DBServiceProxy.Delete(&db.UserGoods{CompID:db.CompID{uid, itemID}}, db.DatabaseHandler)
	}
	//db.DatabaseHandler.DeleteOne()
	return goodsCache
}


//func (this *ExchangeManager) addGoods(uid int32, itemID int32, goodsID string, goods *protocol.Goods) {
//	this.goods[goodsID] = goods
//	userGoods := this.userGoods[uid]
//	if userGoods == nil {
//		userGoods = []*protocol.Goods{}
//	}
//	userGoods = append(userGoods, goods)
//	this.userGoods[uid] = userGoods
//}
//
//func (this *ExchangeManager) removeGoods(uid int32, itemID int32, goodsID string) {
//	delete(this.goods, goodsID)
//	userGoods := this.userGoods[uid]
//	if userGoods == nil {
//		return
//	}
//	for index, goods := range userGoods {
//		if goods.GetId() == itemID {
//			this.userGoods[uid] = append(userGoods[:index], userGoods[index+1:]...)
//		}
//	}
//}

//func (this *ExchangeManager) PublicGoods(uid int32, ) *protocol.Sale {

//func (this *ExchangeManager) GetSaleItem(uid int32) int32 {
//	this.RLock()
//	defer this.RUnlock()
//	sale := this.sales[uid]
//	if (sale == nil) {
//		return 0
//	}
//	return sale.ItemID
//}