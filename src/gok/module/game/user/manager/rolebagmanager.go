package manager

import (
	"gok/module/game/db"
	"gok/service/msg/protocol"
	"gok/service/exception"
	"gok/constant"
	"time"
	"gok/module/game/conf"
	"gok/service/lpc"
	"gok/module/statistics/model"
	"aliens/common/character"
)

//角色标识
type RoleItemManager struct {
	uid int32
	items map[int32]*db.DBRoleItem
	//itemGroups map[int32]map[int32]*db.DBRoleItemGroup   //星球类型/组合ID/组合对象
	//currentItemGroups map[int32]struct{}
	//tempItems []int32
	itemRandomsMapping map[int32]int32
	itemRandoms        []*db.DBItemRandom
	newItemListener    func(int32)
	//itemGroupsRecords map[int32]*db.DBRoleItemGroupRecord
}

//初始化
func (this *RoleItemManager) Init(role *db.DBRole) {
	this.uid = role.UserID
	this.items = make(map[int32]*db.DBRoleItem)
	for _, item := range role.Items {
		this.items[item.ID] = item
	}
	//this.itemGroups = make(map[int32]map[int32]*db.DBRoleItemGroup)
	//
	//for _, itemGroup := range role.ItemGroups {
	//	groupIDString := util.Int32ToString(itemGroup.ID)
	//	starType := util.StringToInt32(groupIDString[2 : 3])
	//	if this.itemGroups[starType] == nil {
	//		this.itemGroups[starType] = make(map[int32]*db.DBRoleItemGroup)
	//	}
	//	this.itemGroups[starType][itemGroup.ID] = itemGroup
	//}
	this.itemRandomsMapping = make(map[int32]int32)
	//var itemsIdGroup *db.DBItemRandom
	for _, itemsIdGroup := range role.ItemRandoms {
		for _,id:= range itemsIdGroup.ItemId{
			this.itemRandomsMapping[id] = itemsIdGroup.Weight
		}
	}
	this.itemRandoms = role.ItemRandoms
	//this.itemGroupsRecords = make(map[int32]*db.DBRoleItemGroupRecord)
	//
	//
	// for _, itemGroup := range role.ItemGroupsRecords {
	//	}




	//this.EnsureItemRandom()
	//if role.ItemRandoms == nil || len(role.ItemRandoms) == 0 {
	//	this.itemRandomsMapping = conf.ItemAddRandomWeight(role.)
	//}
	//this.itemRandomsMapping = make(map[int32]int32)

	//for _, itemRandom := range role.ItemRandoms{
	//	this.itemRandomsMapping[] = itemRandom
	//}


	//this.tempItems = role.TempItems
}

//更新数据库内存
func (this *RoleItemManager) Update(role *db.DBRole) {
	role.Items = this.GetItems()
	//role.ItemGroups = this.GetAllItemGroups()
	role.ItemRandoms = this.itemRandoms
	//role.TempItems = this.tempItems
}

func (this *RoleItemManager) EnsureItemRandom(starType int32, force bool) {
	if force {
		this.itemRandomsMapping = conf.ItemAddRandomWeight(starType)
		this.itemRandoms = this.GetItemRandoms()
	} else {
		if this.itemRandomsMapping == nil || len(this.itemRandomsMapping) == 0 || !this.IsCurrentStarItemRandoms(starType){
			this.itemRandomsMapping = conf.ItemAddRandomWeight(starType)
			this.itemRandoms = this.GetItemRandoms()
		}
	}

}

func (this *RoleItemManager) IsCurrentStarItemRandoms(starType int32) bool {
	if len(this.itemRandoms) > 0 {
		if len(this.itemRandoms[0].ItemId) > 0 {
			itemIDStr := character.Int32ToString(this.itemRandoms[0].ItemId[0]) //10301
			itemStarType := character.StringToInt32(itemIDStr[1:3])
			if itemStarType != starType {
				return false
			}
			return true
		}
	}
	return false
}

func (this *RoleItemManager) GetItemWeight() map[int32]int32{
	return this.itemRandomsMapping
}


func (this *RoleItemManager) SetNewItemListener(listener func(int32)) {
	this.newItemListener = listener
}

func (this *RoleItemManager) GetHaveItemIDs() []int32 {
	result := []int32{}
	for _, item := range this.items {
		if (item.Num > 0) {
			result = append(result, item.ID)
		}
	}
	return result
}

//func (this *RoleItemManager) GetTempItems() []int32 {
//	return this.tempItems
//}
//
//func (this *RoleItemManager) AddTempItem(itemID int32) bool {
//	if len(this.tempItems) == constant.MAX_TEMP_ITEM {
//		return false
//	}
//	this.tempItems = append(this.tempItems, itemID)
//	return true
//}


////清除临时背包
//func (this *RoleItemManager) TakeinAllTempItem() []int32 {
//	this.TakeInItems(this.tempItems, constant.OPT_TYPE_TEMP_BAG, 0)
//	result := this.tempItems
//	this.tempItems = []int32{}
//	return result
//}

//func (this *RoleItemManager) TakeoutTempItems(itemIDs []int32) bool {
//	for _, itemID := range itemIDs {
//		if !character.ContainsInt32(itemID, this.tempItems) {
//			return false
//		}
//	}
//
//	for _, itemID := range itemIDs {
//		this.TakeoutTempItem(itemID)
//	}
//	return true
//}

//func (this *RoleItemManager) TakeoutTempItem(itemID int32) bool {
//	for index, tempItem := range this.tempItems {
//		if tempItem == itemID {
//			this.tempItems = append(this.tempItems[:index], this.tempItems[index+1:]...)
//			return true
//		}
//	}
//	return false
//}

//获取所有物品
func (this *RoleItemManager) GetItems() []*db.DBRoleItem {
	result := []*db.DBRoleItem{}
	for _, item := range this.items {
		result = append(result, item)
	}
	return result
}



func (this *RoleItemManager) GetItemRandoms() []*db.DBItemRandom{
	result := []*db.DBItemRandom{}
	weightMapping := make(map[int32]*db.DBItemRandom)
	for itemID, weight := range this.itemRandomsMapping {
		data := weightMapping[weight]
		if data == nil {
			data = &db.DBItemRandom{Weight:weight, ItemId:[]int32{}}
			weightMapping[weight] = data
		}
		data.ItemId = append(data.ItemId, itemID)
		//conf.DATA.ITEM_RANDOM_DATA
		//result = append(result, itemRandom)
	}
	for _, v := range weightMapping{
		result = append(result, v)
	}
	return result
}

//获取物品数量
func (this *RoleItemManager) GetItemValue(key int32) int32 {
	item := this.items[key]
	if (item == nil) {
		return DEFAULT_FLAG_VALUE
	}
	return item.Num
}





//func (this *RoleItemManager) EnsureItemsGroup(groupID int32) *db.DBRoleItemGroup {
//	itemGroup := this.GetItemGroup(groupID)//this.itemGroups[groupID]
//	if itemGroup == nil {
//		itemGroup = &db.DBRoleItemGroup{ID:groupID,Done:false}
//		this.GetItemGroup(groupID) = itemGroup
//	}
//	//records :=
//	if itemGroup.Records == nil {
//		itemGroup.Records = []*db.DBRoleItemGroupRecord{}
//		//records = append(records, record)
//	}
//	return itemGroup
//}




//func (this *RoleItemManager) AddItemsGroupRecord(key int32, items []int32 ,reward bool, done bool, ) bool {
//	record := this.itemGroupsRecords[key].Record
//	newRecord := db.Record{Items:items ,Reward:reward,Done:done}
//
//	record = append(record, newRecord)
//	//this.itemGroups[key].TryTime = time.Now()
//	return true
//}
//
//func (this *RoleItemManager) CompareItemsGroup(items []int32,itemsgroup []*conf.ItemGroupBase) bool {
//	for _, confItems := range itemsgroup {
//		if(len(items) != len(confItems.Content)) {
//			continue
//		}
//		ret := compareItemsWithConfItems(items,confItems.Content)
//		if ret == true {
//			return true
//		}
//	}
//	return false
//}
//
//func compareItemsWithConfItems(items []int32,confItems []int32) bool{
//	for _, item := range items {
//		ret := compareItemWithConfItems(item, confItems)
//		if ret == false {
//			return false
//		}
//	}
//	return true
//}
//func compareItemWithConfItems(item int32,confItems []int32) bool {
//	for _, confItem := range confItems {
//		if item == confItem {
//			return true
//		}
//	}
//	return false
//}
//func (this *RoleItemManager) ActiveItemGroup(key int32, items []int32) *db.DBRoleItemGroup {
//	for _, item := range items {
//		if this.ContainsGroupItem(key, item) {
//			return nil
//		}
//	}
//	itemGroup := this.itemGroups[key]
//	if itemGroup == nil {
//		itemGroup = &db.DBRoleItemGroup{
//			ID:     key,
//			//Items:  items,
//			Reward: false,
//			//TryTime:time.Now(),
//		}
//		this.itemGroups[key] = itemGroup
//	} else {
//		//itemGroup.Items = items
//	}
//	return itemGroup
//}

//func (this *RoleItemManager) AddGroupItems(key int32, items []int32) *db.DBRoleItemGroup {
//	itemGroup := this.itemGroups[key]
//	if (itemGroup == nil) {
//		return nil
//	}
//	for _, item := range items {
//		if (this.ContainsGroupItem(key, item)) {
//			continue
//		}
//		itemGroup.Items = append(itemGroup.Items, item)
//	}
//	return itemGroup
//}

//func (this *RoleItemManager) AddGroupItem(key int32, item int32) bool {
//	itemGroup := this.itemGroups[key]
//	if (itemGroup == nil) {
//		return false
//	}
//	if (this.ContainsGroupItem(key, item)) {
//		return true
//	}
//	itemGroup.Items = append(itemGroup.Items, item)
//	return true
//}
//
//
//func (this *RoleItemManager) ContainsGroupItem(key int32, item int32) bool {
//	itemGroup := this.itemGroups[key]
//	if (itemGroup == nil) {
//		return false
//	}
//	for _, oldItem := range itemGroup.Items {
//		if (oldItem == item) {
//			return true
//		}
//	}
//	return false
//}


func (this *RoleItemManager) CanTakeoutItems(itemIDs []int32) bool {
	for _, itemID := range itemIDs {
		if !this.CanTakeoutItem(itemID, 1) {
			return false
		}
	}
	return true
}

//获取物品数量
func (this *RoleItemManager) CanTakeoutItem(key int32, value int32) bool {
	if (value < 0) {
		return false
	}
	item := this.items[key]
	if (item == nil) {
		return false
	}
	return item.Num >= value
}

func (this *RoleItemManager) TakeOutItems(itemIDs []int32, operation constant.OPT, refID int32) {
	for _, itemID := range itemIDs {
		this.TakeOutItem(itemID, 1, operation, refID)
	}
}


func (this *RoleItemManager) TakeOutItem(key int32, value int32, operation constant.OPT, refID int32) {
	if key == 0 {
		return
	}
	if value < 0 {
		return
	}
	item := this.items[key]
	if item == nil {
		exception.GameException(exception.ITEM_NOT_ENOUGH)
	}
	if item.Num < value {
		exception.GameException(exception.ITEM_NOT_ENOUGH)
	}
	item.Num-=value

	lpc.StatisticsHandler.AddStatisticData(&model.StatisticItem{
		UserID:this.uid,
		ItemID:item.ID,
		RefID:refID,
		Operation:uint8(operation),
		Change:-value,
		Total:item.Num,
	})

	//lpc.LogServiceProxy.AddItemRecord(this.uid, item.ID, refID, operation, -value , item.Num)
}

func (this *RoleItemManager) TakeInItems(items []int32, operation constant.OPT, refID int32) {
	if (items == nil) {
		return
	}
	for _, item := range items {
		this.TakeInItem(item, 1, operation, refID)
	}
}

func (this *RoleItemManager) TakeInItem(key int32, value int32, operation constant.OPT, refID int32) *db.DBRoleItem {
	if key == 0 {
		return nil
	}
	if value < 0 {
		return nil
	}
	item := this.items[key]
	if item == nil {
		item = &db.DBRoleItem{
			ID:  key,
			Num: value,
			UpdateTime: time.Now(),
		}
		this.items[key] = item
		if this.newItemListener != nil {
			this.newItemListener(key)
		}
	} else {
		item.Num += value
	}


	lpc.StatisticsHandler.AddStatisticData(&model.StatisticItem{
		UserID:this.uid,
		ItemID:item.ID,
		RefID:refID,
		Operation:uint8(operation),
		Change:value,
		Total:item.Num,
	})

	//lpc.LogServiceProxy.AddItemRecord(this.uid, item.ID, refID, operation, value , item.Num)
	return item
}

////更新物品数量
//func (this *RoleItemManager) UpdateItem(key int32, value int32) *db.DBRoleItem{
//	item := this.items[key]
//	if (item == nil) {
//		item = &db.DBRoleItem{ID:key, Num:value}
//		this.items[key] = item
//	} else {
//		item.Num = value
//	}
//	return item
//}

//获取所有物品
func (this *RoleItemManager) GetProtocolItems() []*protocol.BagItem {
	result := []*protocol.BagItem{}
	for _, item := range this.items {
		result = append(result, &protocol.BagItem{
			Id:item.ID,
			Num:item.Num,
			Time:item.UpdateTime.Unix(),
		})
	}
	return result
}

//获得当前组合中所拥有最少的圣物的ID
func (this *RoleItemManager) GetMinNumItemFromItems(items []int32) int32 {
	var minNum int32 = 0
	var minID int32 = 0
	for _, itemID := range items {
		if this.items[itemID] == nil {
			minID = itemID
			break
		}
		num := this.items[itemID].Num
		if num == 0 {
			minID = itemID
			break
		}
		if minNum == 0{
			minNum = num
			minID = itemID
		}
		if num < minNum {
			minNum = num
			minID = itemID
		}
	}
	return minID
}