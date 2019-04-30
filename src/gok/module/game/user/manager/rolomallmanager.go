package manager

import (
	"gok/service/msg/protocol"
	"time"
	"gok/service/exception"
	"gok/service/rpc"
	"aliens/common/character"
	"math/rand"
	"github.com/name5566/leaf/util"
	"gok/module/game/db"
	"gok/module/game/conf"
	util2 "aliens/common/util"
)

type RoleMallManager struct {
	uid             int32
	mallItems       []*protocol.MallItem
	mallItemsReTime time.Time
}

func (this *RoleMallManager) Init(role *db.DBRole) {
	this.uid = role.UserID
	this.mallItems = role.MallItems
	this.mallItemsReTime = role.MallItemsReTime
}

//更新数据库内存
func (this *RoleMallManager) Update(role *db.DBRole) {
	role.MallItems = this.mallItems
	role.MallItemsReTime = this.mallItemsReTime
	role.UserID = this.uid
}

func (this *RoleMallManager) CleanMallItems() {
	this.mallItems = nil
}

func (this *RoleMallManager) RefreshMallItems(auto bool, starType int32, costCurrent int32, costNext int32) ([]*protocol.MallItem, time.Time) {
	//var items []*protocol.MallItem = nil

	if auto {
		if time.Now().After(this.mallItemsReTime) {
			this.mallItems = randomMallItems(this.uid, starType)
			this.reRefreshTime()
		}
	} else {
		this.mallItems = randomMallItems(this.uid, starType)
	}
	return this.mallItems, this.mallItemsReTime
}

func (this *RoleMallManager) CalcReMallItemCost(reCount int32) (int32, int32) {
	mallItemsCost := conf.DATA.ShopReCost
	if mallItemsCost == nil || len(mallItemsCost) != 3{
		exception.GameException(exception.ITEM_BASE_NOT_FOUND)
	}
	costFirst := mallItemsCost[0]
	costInc := mallItemsCost[1]
	costFinal := mallItemsCost[2]
	cost := costFirst + costInc*reCount
	if cost > costFinal {
		cost = costFinal
	}
	costNext := cost + costInc
	if costNext > costFinal {
		costNext = costFinal
	}
	return cost, costNext
}

func randomMallItems(uid int32, starType int32) []*protocol.MallItem {
	starItems := conf.DATA.ITEM_RANDOM_DATA[starType].Items
	itemsNum := conf.DATA.ShopItemsNum
	if starItems == nil || itemsNum == 0 || len(starItems) < itemsNum {
		exception.GameException(exception.ITEM_BASE_NOT_FOUND)
	}
	itemIDs := rpc.StarServiceProxy.GetCurrentGroupItems(uid).GetItemIDs()
	var baseItemsNum = 0
	shopBaseNum := conf.DATA.ShopBaseItemsNum
	if shopBaseNum == nil || len(shopBaseNum) != 2 {
		exception.GameException(exception.ITEM_BASE_NOT_FOUND)
	}
	if itemIDs != nil {
		if len(itemIDs) == 3 {
			baseItemsNum = shopBaseNum[0]
		}
		if len(itemIDs) == 5 {
			baseItemsNum = shopBaseNum[1]
		}
	}
	var baseItems []int32
	var randomItems []int32
	if baseItemsNum != 0 {
		baseItems = randomItemFromItems(baseItemsNum, itemIDs)
	}

	randomItems = eliminateArrayFromArray(baseItems, starItems)

	resultItems := character.AppendArray(baseItems, randomItemFromItems(itemsNum - baseItemsNum, randomItems))
	resultLen := len(resultItems)
	result := make([]*protocol.MallItem, resultLen)

	resultIndex := rand.Perm(len(resultItems))

	for indexItem, indexID := range resultIndex {
		result[indexID] = &protocol.MallItem{
			ID:int32(indexID),
			ItemID:resultItems[indexItem],
			Num:randomArrayScope(conf.DATA.RelicQuantity),
			BuyTimes:randomArrayScope(conf.DATA.BuyRelic),
			GroupCost:int32(randomArrayScope(conf.DATA.RelicPrice)/10)*10,
		}
	}
	return result
}

//从数组中剔除掉数组
func eliminateArrayFromArray(desArray []int32, srcArray []int32) []int32 {

	if len(desArray) > len(srcArray) {
		return nil
	}
	result := character.CopyArray(srcArray)
	for _, desNum := range desArray {
		for index, srcNum := range result {
			if srcNum == desNum {
				result = append(result[:index], result[index+1:]...)
				break
			}
		}
	}
	return result
}

func randomItemFromItems(num int, items []int32) []int32 {
	if num > len(items) {
		return  nil
	}
	randomItems := character.CopyArray(items)
	result := make([]int32, num)
	for i := 0; i < num; i++ {
		index := rand.Intn(len(randomItems))
		result[i] = randomItems[index]
		randomItems = append(randomItems[:index], randomItems[index+1:]...)
	}
	return result
}

func randomArrayScope(array []int32) int32 {
	if array == nil || len(array) != 2 {
		array = []int32{0,0}
	}
	return util.RandInterval(array[0], array[1])
}

func (this *RoleMallManager) reRefreshTime( /*reFromNow bool*/) {
	//refreshTime := conf.DATA.ShopRe
	//if reFromNow {
	//	this.mallItemsReTime = time.Now().Add(time.Duration(int64(time.Second) * int64(refreshTime)))
	//} else {
	//	reTime := this.mallItemsReTime
	//	intervalTime := time.Now().Sub(reTime).Seconds()
	//	ratio := int64(intervalTime/refreshTime) + 1
	//	this.mallItemsReTime = reTime.Add(time.Duration(int64(time.Second) * ratio *int64(refreshTime)))
	//}
	refreshTime := conf.DATA.ShopReHour
	if refreshTime == 0 {
		exception.GameException(exception.ITEM_BASE_NOT_FOUND)
	}
	ratio := int(time.Now().Hour()/refreshTime) + 1
	this.mallItemsReTime = util2.GetTodayHourTime(ratio * refreshTime)
}

func (this *RoleMallManager) GetMallItems(starType int32) ([]*protocol.MallItem, time.Time) {
	//if this.mallItems== nil || len(this.mallItems) == 0 {
	//	this.mallItems = randomMallItems(this.uid, starType)
	//	this.reRefreshTime()
	//} else {
	//	if time.Now().After(this.mallItemsReTime) {
	//		this.mallItems = randomMallItems(this.uid, starType)
	//		this.reRefreshTime()
	//	}
	//}
	if this.mallItems != nil && len(this.mallItems) != 0 && time.Now().Before(this.mallItemsReTime) {
		return this.mallItems, this.mallItemsReTime
	}
	this.mallItems = randomMallItems(this.uid, starType)
	this.reRefreshTime()
	return this.mallItems, this.mallItemsReTime
}

func (this *RoleMallManager) BuyMallItem(id int32) *protocol.MallItem {
	//item := this.mallItems[id]
	if this.mallItems[id] == nil {
		exception.GameException(exception.MALL_ITEM_ID_NOT_FOUND)
	}
	if this.mallItems[id].BuyTimes <= 0 {
		exception.GameException(exception.MALL_ITEM_TIMES_NOT_ENOUGH)
	}

	this.mallItems[id].BuyTimes--
	return this.mallItems[id]
}
