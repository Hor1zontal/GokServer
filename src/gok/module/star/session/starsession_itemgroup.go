package session

import (
	"aliens/common/character"
	"gok/constant"
	"gok/module/star/conf"
	"gok/module/star/db"
	"gok/module/statistics/model"
	"gok/service/exception"
	"gok/service/lpc"
	"gok/service/msg/protocol"
	"gok/service/rpc"
	"math/rand"
)

/*------------------------------- 圣物组合操作相关函数 ----------------------------------*/

//圣物组合开启
func (this *StarSession) ItemGroupsInit() {
	//var group
	confLen := len(conf.Base.GroupUnlockBase)
	for i := 1; i <= confLen; i++ {
		groupID := this.RandomGroupIDByNum(conf.Base.GroupUnlockBase[int32(i)].UnlockType)
		this.SetItemGroupAbleActive(groupID)
	}
	this.setDirty()
}

func (this *StarSession) AddItemsGroupRecord(groupID int32, items []int32 ,satifyNum int32) bool {
	itemGroup := this.GetItemGroup(groupID)
	record := &db.DBStarItemGroupRecord{Items:items,Num:satifyNum}//number
	itemGroup.Records = append(itemGroup.Records,record)
	return true

}

func (this *StarSession) SetItemGroupDone(groupID int32) *protocol.CivilizationInfo {
	var civilizationInfo *protocol.CivilizationInfo = nil
	index, itemGroup := this.GetItemGroupWithIndex(groupID)
	if !itemGroup.Done {
		itemGroup.Done = true

		if index == 0 {
			this.updateFirstGroupFlagDone()
		}
		unlockBase := conf.Base.GroupUnlockBase[int32(index + 1)]
		if unlockBase != nil && unlockBase.CivilizationIncome > 0 {
			civilizationInfo = this.TakeInCivilization(unlockBase.CivilizationIncome)
		}

		lpc.StatisticsHandler.AddStatisticData(&model.StatisticItemGroup{
			UserID:this.Owner,
			StarType:this.Type,
			GroupID:itemGroup.ID,
			Seq:index + 1,
			Type:model.ITEM_GROUP_TYPE_DONE,
		})

		this.DealItemGroupOpen()
	}
	return civilizationInfo
}

//func (this *StarSession) SetItemGroupActive(groupID int32) *db.DBStarItemGroup{
//	itemGroup := this.GetItemGroup(groupID)
//	if !itemGroup.Active {
//		itemGroup.Active = true
//	}
//	return itemGroup
//}

func (this *StarSession) SetItemGroupAbleActive(groupID int32) *db.DBStarItemGroup {
	itemGroup := this.GetItemGroup(groupID)
	if itemGroup == nil {
		itemGroup = &db.DBStarItemGroup{ID:groupID,Done:false,Records:[]*db.DBStarItemGroupRecord{}}
		this.ItemGroups = append(this.ItemGroups, itemGroup)
	}
	return itemGroup
}


func (this *StarSession) IsDoneItemGroup(groupID int32) bool {
	itemGroup := this.GetItemGroup(groupID)
	if itemGroup == nil {
		return false
	}
	return itemGroup.Done
}

func (this *StarSession) GetProtocolItemGroups() []*protocol.ItemGroup {
	result := []*protocol.ItemGroup{}
	for _,itemGroup := range this.ItemGroups{
		result = append(result, itemGroup.BuildProtocol())
	}
	return result
}

func (this *StarSession) DealItemGroupOpen()  {
	//levelNum := this.UpdateLevelNum()
	//根据建筑解锁情况获取当前激活进度
	requireID := conf.GetUnlockGroupConfigID(this.BuildingExMaxLevel)
	//requireLen := conf.
	currentLen := this.getCurrentLen()

	for ; currentLen < requireID; currentLen++ {
		length := int(character.Int32Min(int32(len(this.ItemGroups)), int32(currentLen + 1)))
		var lastItemGroup *db.DBStarItemGroup
		for i:=0; i< length; i ++ {
			itemGroup := this.ItemGroups[i]
			//上一个圣物组合是否激活
			isLastItemGroupDone := lastItemGroup == nil || lastItemGroup.Done
			if isLastItemGroupDone && !itemGroup.Active {
				itemGroup.Active = true

				lpc.StatisticsHandler.AddStatisticData(&model.StatisticItemGroup{
					UserID:this.Owner,
					StarType:this.Type,
					GroupID:itemGroup.ID,
					Seq:i + 1,
					Type:model.ITEM_GROUP_TYPE_ACTIVE,
				})
				if i == 0 {
					this.updateFirstGroupFlagUnlock()
				}
				PushItemGroupOpenInfo(this.Owner,itemGroup.ID)
			}
			lastItemGroup = itemGroup
		}


	}
}

func (starSession *StarSession) ActiveGroup(groupID int32, itemIDs []int32) *protocol.ActiveGroupRet {
	var civilizationInfo *protocol.CivilizationInfo = nil
	result := &protocol.ActiveGroupRet{}
	groupBase := conf.Base.ItemGroup[groupID]
	if groupBase == nil {
		exception.GameException(exception.ITEM_GOURP_BASE_NOT_FOUND)
	}
	if starSession.IsDoneItemGroup(groupID) {
		exception.GameException(exception.ITEM_GROUP_IS_DONE)
	}
	satifyNum := groupBase.ContainsItemsNum(itemIDs)
	if int(satifyNum) != len(itemIDs) {
		starSession.AddItemsGroupRecord(groupID, itemIDs, satifyNum)
		result.Result = false
	} else {
		civilizationInfo = starSession.SetItemGroupDone(groupID)
		starSession.AddItemsGroupRecord(groupID, itemIDs, satifyNum)
		starSession.SetBuff(groupID)
		result.Result = true
	}
	starSession.setDirty()
	result.SatifyNum = satifyNum
	result.CivilizationInfo = civilizationInfo
	return result
}

func (this *StarSession) RandomGroupIDByNum(numType int32) int32 {

	var groupID int32 = 0
	groups := conf.Base.StarGroupMapping[this.Type]
	if groups == nil {
		exception.GameException(exception.ITEM_GOURP_BASE_NOT_FOUND)
	}
	var tempData3 = []int32{}
	var tempData5 = []int32{}
	for _, group := range groups {
		itemGroup := this.GetItemGroup(group.ID)
		if itemGroup == nil {
			if group.Rarity == constant.NUM_THREE {
				tempData3 = append(tempData3, group.ID)
			} else if group.Rarity == constant.NUM_FIVE {
				tempData5 = append(tempData5, group.ID)
			}
		}
	}
	if numType == constant.UNLOCK_NUM_THREE {
		randomIndex := rand.Intn(len(tempData3))
		groupID = tempData3[randomIndex]
	} else if numType == constant.UNLOCK_NUM_FIVE {
		randomIndex := rand.Intn(len(tempData5))
		groupID = tempData5[randomIndex]
	}
	return groupID
}

// return 正在尝试的组合ID，正在尝试的组合的序列号(index+1)，完成的组合的数量
func (this *StarSession) GetCurrentActiveGroup() (int32, int32, int32) {
	var groupID int32
	var activeIndex int32
	var doneNum int32
	for index, group := range this.ItemGroups {
		if !group.Active {
			break
		}
		if !group.Done {
			groupID = group.ID
			activeIndex = int32(index)
		} else {
			doneNum = int32(index + 1)
		}
	}
	return groupID, activeIndex, doneNum
}

func (this *StarSession) GetCurrentActiveGroupItems() ([]int32, int32) {
	groupID, _, _ := this.GetCurrentActiveGroup()
	if groupID == 0 {
		return nil, 0
	}
	if conf.Base.ItemGroup[groupID] == nil {
		exception.GameException(exception.ITEM_GOURP_BASE_NOT_FOUND)
	}
	return conf.Base.ItemGroup[groupID].Content, groupID
}

func (this *StarSession) getCurrentLen() int {
	for index, group := range this.ItemGroups {
		if group.Active == false {
			return index
		}
	}
	return 0
}

func PushItemGroupOpenInfo(uid int32, groupID int32) {
	rpc.UserServiceProxy.Push(uid, &protocol.GS2C{
		Sequence:[]int32{1041},
		ItemGroupOpenPush:&protocol.ItemGroupOpenPush{GroupID:groupID},
	})
}

func (this *StarSession) GetItemGroup(groupID int32) *db.DBStarItemGroup {
	_, result := this.GetItemGroupWithIndex(groupID)
	return result
}

func (this *StarSession) GetItemGroupWithIndex(groupID int32) (int, *db.DBStarItemGroup) {
	for index, itemGroup := range this.ItemGroups {
		if groupID == itemGroup.ID {
			return index, itemGroup
		}
	}
	return -1, nil
}

func (this *StarSession) getItemGroups() []*db.DBStarItemGroup {
	return this.ItemGroups
}

func (this *StarSession) getItemGroupFirstItemID() int32 {
	if this.ItemGroups[0] == nil {
		return 0
	}
	return this.ItemGroups[0].ID
}