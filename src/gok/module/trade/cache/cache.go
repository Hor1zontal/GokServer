package cache

import (
	"gok/module/trade/db"
	basecache "gok/cache"
	"aliens/log"
	"gok/module/trade/conf"
	"gok/service/msg/protocol"
)

var TradeCache = basecache.NewTradeCacheManager()

func Init() {
	TradeCache.Init1(conf.Server.RedisAddress, conf.Server.RedisPassword,
		conf.Server.RedisMaxActive, conf.Server.RedisMaxIdle, conf.Server.RedisIdleTimeout)

	//if TradeCache.SetNX(basecache.FLAG_LOAD_SALE, 1) {
	//	log.Debug("start load sale data to redis cache...")
	//	var sales []*db.UserSale
	//	db.DatabaseHandler.QueryAll(&db.UserSale{}, &sales)
	//	for _, sale := range sales {
	//		TradeCache.SetUserPublicSale(sale.UserID, sale.BuildProtocol())
	//	}
	//	log.Debug("end load sale data to redis cache")
	//}
	if TradeCache.SetNX(basecache.FLAG_LOAD_SALE, 1) {
		log.Debug("start load sale data to redis cache...")
		count := 0
		var sales []*db.UserSale
		err := db.DatabaseHandler.QueryAllLimit(&db.UserSale{}, &sales, 10000, func(data interface{}) bool {
			for _, sale := range sales {
				TradeCache.SetUserPublicSale(sale.UserID, sale.BuildProtocol())
			}
			currLen := len(sales)
			count += currLen
			return currLen == 0
		})
		if err != nil {
			log.Debug("load sale err: %v", err)
		}
		log.Debug("end load sale data to redis cache:%v", count)
	}


	if TradeCache.SetNX(basecache.FLAG_LOAD_ITEMHELP, 1) {
		log.Debug("start load itemHelp data to redis cache...")
		count := 0
		var itemHelps []*protocol.ItemHelp
		err := db.DatabaseHandler.QueryAllLimit(&protocol.ItemHelp{}, &itemHelps, 10000, func(data interface{}) bool {
			for _, itemHelp := range itemHelps {
				TradeCache.SetItemHelp(itemHelp.GetUid(), itemHelp)
			}
			currLen := len(itemHelps)
			count += currLen
			return currLen == 0
		})
		if err != nil {
			log.Debug("load itemHelp err: %v", err)
		}
		log.Debug("end load itemHelp data to redis cache:%v", count)
	}

	//if TradeCache.SetNX(basecache.FLAG_LOAD_GOOD, 1) {
	//	log.Debug("start load good data to redis cache...")
	//	var goods []*db.UserGoods
	//	db.DatabaseHandler.QueryAll(&db.UserGoods{}, &goods)
	//	for _, good := range goods {
	//		TradeCache.HSetCacheGood(good.GetUserID(), good.GetItemID(), good.BuildProtocol())
	//	}
	//	log.Debug("end load good data to redis cache")
	//}

	if TradeCache.SetNX(basecache.FLAG_LOAD_GOOD, 1) {
		log.Debug("start load good data to redis cache...")
		count := 0
		var goods []*db.UserGoods
		err := db.DatabaseHandler.QueryAllLimit(&db.UserGoods{}, &goods, 10000, func(data interface{}) bool {
			for _, good := range goods {
				TradeCache.HSetCacheGood(good.GetUserID(), good.GetItemID(), good.BuildProtocol())
			}
			currLen := len(goods)
			count += currLen
			return currLen == 0
		})
		if err != nil {
			log.Debug("load good err: %v", err)
		}
		log.Debug("end load good data to redis cache:%v", count)
	}
}

func Close() {
	TradeCache.Close()
}

