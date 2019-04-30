package cache

import (
	"aliens/common/character"
	"gok/service/msg/protocol"
	"aliens/log"
	"aliens/common/util"
)

const (

	USER_SALE_KEY = "sale:"
	USER_ITEMHELP_KEY = "itemHelp:"

	USER_GOODS_KEY = "goods:"
	GOOD_ID_PREFIX = "gid_"

	FLAG_LOAD_SALE = "flag:sale"
	FLAG_LOAD_GOOD = "flag:goods"
	FLAG_LOAD_ITEMHELP = "flag:itemHelp"
)

type TradeCacheManager struct {
	*cacheManager
}

func NewTradeCacheManager() *TradeCacheManager {
	return &TradeCacheManager{
		&cacheManager{},
	}
}

func getUserSaleKey(id int32) string {
	return USER_SALE_KEY + character.Int32ToString(id)
}

func getItemHelpKey(id int32) string {
	return USER_ITEMHELP_KEY + character.Int32ToString(id)
}


func getUserGoodsKey(id int32) string {
	return USER_GOODS_KEY + character.Int32ToString(id)
}

func getGoodIDPrefix(id int32) string {
	return GOOD_ID_PREFIX + character.Int32ToString(id)
}

func (this *TradeCacheManager) ExistItemHelp(uid int32) bool {
	result, _ := this.redisClient.Exists(getItemHelpKey(uid))
	return result
}

func (this *TradeCacheManager) SetItemHelp(uid int32, sale *protocol.ItemHelp) bool{
	data, err := sale.Marshal()
	if err != nil {
		log.Error("itemHelp marshal error: %v", err.Error())
		return false
	}
	return this.redisClient.SetData(getItemHelpKey(uid), data)
}

func (this *TradeCacheManager) GetItemHelp(uid int32) *protocol.ItemHelp {
	data := this.redisClient.GetBytesData(getItemHelpKey(uid))
	if data == nil || len(data) == 0{
		return nil
	}
	itemHelp := &protocol.ItemHelp{}
	err := itemHelp.Unmarshal(data)
	if err != nil {
		log.Error("itemHelp unmarshal err: %v", err.Error())
		return nil
	}
	return itemHelp
}

func (this *TradeCacheManager) GetItemHelpData(uid int32) []byte {
	return this.redisClient.GetBytesData(getItemHelpKey(uid))
}

func (this *TradeCacheManager) RemoveItemHelp(uid int32) bool {
	return this.redisClient.DelData(getItemHelpKey(uid))
}


func (this *TradeCacheManager) SetUserPublicSale(uid int32, sale *protocol.Sale) bool{
	data, err := sale.Marshal()
	if err != nil {
		log.Error("sale bson marshal error: %v", err.Error())
		return false
	}
	return this.redisClient.SetData(getUserSaleKey(uid), data)
}

func (this *TradeCacheManager) GetUserPublicSale(uid int32) *protocol.Sale {
	data := this.redisClient.GetBytesData(getUserSaleKey(uid))
	if data == nil || len(data) == 0{
		return nil
	}
	sale := &protocol.Sale{}
	err := sale.Unmarshal(data)
	if err != nil {
		log.Debug("sale bson unmarshal err: %v", err.Error())
		return nil
	}
	return sale
}

func (this *TradeCacheManager) RemoveUserPublicSale(uid int32) bool {
	return this.redisClient.DelData(getUserSaleKey(uid))
}

func (this *TradeCacheManager) HSetCacheGood(uid int32, itemID int32, good *protocol.Goods) bool{
	data, err := good.Marshal()
	if err != nil {
		log.Debug("good bson marshal error: %v", err.Error())
		return false
	}
	return this.redisClient.HSet(getUserGoodsKey(uid), getGoodIDPrefix(itemID), data)
}

func (this *TradeCacheManager) HGetCacheGood(uid int32, itemID int32) *protocol.Goods {
	data := this.redisClient.HGetBytes(getUserGoodsKey(uid), getGoodIDPrefix(itemID))
	if data == nil || len(data) == 0 {
		return nil
	}
	good := &protocol.Goods{}
	err := good.Unmarshal(data)
	if err != nil {
		log.Debug("good bson unmarshal error: %v", err.Error())
		return nil
	}
	return good
}

func (this *TradeCacheManager) HGetCacheGoods(uid int32) []*protocol.Goods {
	goodsData := this.redisClient.HGetAll(getUserGoodsKey(uid))
	var results []*protocol.Goods
	for _, goodData := range goodsData {
		good := &protocol.Goods{}
		err := good.Unmarshal(util.Str2Bytes(goodData))
		if err != nil {
			log.Debug("good bson unmarshal error: %v", err.Error())
			continue
		}
		results = append(results, good)
	}
	return results
}

func (this *TradeCacheManager) HRemoveCacheGood(uid int32, itemID int32) bool {
	err := this.redisClient.HDel(getUserGoodsKey(uid), getGoodIDPrefix(itemID))
	if err != nil {
		log.Debug("cache delete good error: %v", err.Error())
		return false
	}
	return true
}
