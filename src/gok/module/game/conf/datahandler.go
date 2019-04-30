package conf

import (
	"aliens/common/util"
	"gok/constant"
	"math/rand"
	"gok/service/exception"
	"gok/service/msg/protocol"
)

//获取事件对于步骤触发的任务数据
//func GetTriggerTaskData(eventType int32, step int32) *EventTaskData  {
//	for _, taskData := range DATA.EVENT_TASK_DATA {
//		if (taskData.RefID == eventType && taskData.TriggerStep == step) {
//			return taskData
//		}
//	}
//	return nil
//}

func RandomCardReward(eventType int32) (*CardReward, []*CardReward) {
	mapping := DATA.EVENT_CARD_MAPPING[eventType]
	if mapping == nil {
		return nil, nil
	}

	weightData := util.RandomWeightData(mapping)
	if weightData == nil {
		return nil, nil
	}

	randomResult := weightData.(*CardReward)
	if randomResult == nil {
		return nil, nil
	}

	otherFaithCard := []*CardReward{}

	for _, data := range mapping {
		cardReward := data.(*CardReward)
		if cardReward.ID != randomResult.ID {
			otherFaithCard = append(otherFaithCard, cardReward)
		}
	}
	return randomResult, otherFaithCard
}

//func RandomCardFaith(eventType int32) int32 {
//	mapping := DATA.EVENT_CARD_MAPPING[eventType]
//	if mapping == nil {
//		return 0
//	}
//	for _, data := range mapping {
//		cardReward := data.(*CardReward)
//		if cardReward.RewardType == constant.CARD_REWARD_FAITH {
//			return cardReward.RandomValue()
//		}
//	}
//	return 0
//}

func ContainsAll(arrays []int32, subArrays []int32) bool {
	for _, subArray := range subArrays {
		if (!Contains(arrays, subArray)) {
			return false
		}
	}
	return true
}


func Contains(arrays []int32, data int32) bool {
	if (arrays == nil || len(arrays) == 0) {
		return false
	}
	for _, array := range arrays {
		if (data == array) {
			return true
		}
	}
	return false
}

func GetTriggerTaskData(eventType int32) *EventTaskData {
	taskType := DATA.EVENT_TASK_MAPPING[eventType]
	return DATA.EVENT_TASK_DATA[taskType]
}

func GetTaskData(taskType int32) *EventTaskData {
	return DATA.EVENT_TASK_DATA[taskType]
}

func GetEventTask(eventType int32) *EventTaskData {
	for _, task := range DATA.EVENT_TASK_DATA {
		if task.RefID == eventType {
			return task
		}
	}
	return nil
}

func RandStarItem(starType int32, randomWeight map[int32]int32, randomType uint8) int32 {
	randomItem := DATA.ITEM_RANDOM_DATA[starType]
	if randomItem == nil {
		return 0
	}

	//manager.DataManager{}
	//randomWeight =
	//ItemAddRandomWeight(starType)
	//randomWeight = ItemAddRandomWeight(starType)

	var allowMapping map[int32]int32 = nil
	if randomType == constant.RANDOM_BELIEVER {
		allowMapping = randomItem.BelieverWeight
	} else if randomType == constant.RANDOM_BUILDING {
		allowMapping = randomItem.BuildingWeight
	} else if randomType == constant.RANDOM_TASK {
		allowMapping = randomItem.TaskWeight
	}
	if allowMapping == nil {
		return 0
	}

	//fmt.Printf("%d",util.RandomWeight(randomWeight) )
	itemID :=util.RandomWeight(randomWeight)
	if allowMapping[itemID] == 0 {
		return 0
	}
 	return itemID
}


//获取星球建筑配置信息
func GetBuildingConf(starType int32, buildingType int32, buildingLevel int32) *BuildingBase {
	for _, building := range DATA.BUILDING_DATA {
		if building.StarID == starType &&
			building.BuildID == buildingType &&
			building.Level == buildingLevel {
			return building
		}
	}
	return nil
}

//获取 游戏基础数据
func GetGameBase() *GameBase {
	return DATA.GameBase;
}

//获取商店数据
func GetShopBase(id int32) *ShopBase {
	if v, ok := DATA.SHOP_BASEDATA[ id ]; ok {
		return v
	}
	return nil
}

func GetTaskReward(taskId int32, endingId int32) *TaskRewardBase {
	rewards := DATA.TASK_REWARD_DATA[taskId]
	if (rewards == nil) {
		return nil
	}
	return rewards[endingId]
}

func RandomAttackAstrola(starType int32) int32 {
	randomPool := make(map[int32]int32)
	for _, astrolas := range DATA.ASTROLA_DATA {
		for _, astrola := range astrolas {
			//过滤星球类型不匹配的条件
			if astrola.StarType != 0 && astrola.StarType != starType {
				continue
			}
			if !constant.IsAttackEvent(astrola.EventBase) {
				continue
			}
			randomPool[astrola.EventBase] = 1
		}

	}
	return util.RandomWeight(randomPool)
}

//随机星盘事件
//func RandomAstrola(attackEventFilter bool, rewardTypeFileter int32, believerCount int32, starType int32, buildingLevel int32, weightAppend map[int32]int32) int32 {
//	//var weightCount int32 = 0
//	randomPool := make(map[int32]int32)
//	isBelieverFull := believerCount >= DATA.StarBelieverLimit
//	for rewardType, astrolas := range DATA.ASTROLA_DATA {
//		if rewardTypeFileter != 0 && rewardTypeFileter != rewardType {
//			continue
//		}
//		for _, astrola := range astrolas {
//			//过滤星球类型不匹配的条件
//			if astrola.StarType != 0 && astrola.StarType != starType {
//				continue
//			}
//			if attackEventFilter && !constant.IsAttackEvent(astrola.EventBase) {
//				continue
//			}
//			//过滤已经存在的同类型任务的条件
//			//if astrola.IsOnly {
//			//	taskData := GetTriggerTaskData(astrola.EventBase)
//			//	if (taskData == nil || existTaskTypes.Contains(taskData.ID)) {
//			//		continue
//			//	}
//			//}
//			//randomPool(astrola)
//			//randomPool = append(randomPool, astrola)
//
//			//没有信徒不需要刷信仰任务
//			//if astrola.EventBase == constant.EVENT_ID_REFRESH_FAITH && believerCount == 0 {
//			//	continue
//			//}
//
//			//信徒满了不需要随机刷信徒任务
//			if isBelieverFull && astrola.RewardType == constant.EVENT_REWARD_TYPE_BELIEVER {
//				continue
//			}
//
//			//信徒为0的时候不允许出刷信仰任务
//			if believerCount == 0 && astrola.EventBase == constant.EVENT_ID_REFRESH_FAITH {
//				continue
//			}
//
//			randomPool[astrola.EventBase] = astrola.getWeight(buildingLevel) + weightAppend[astrola.EventBase]
//		}
//
//	}
//	return util.RandomWeight(randomPool)
//
//	////随机权重
//	//randomWeight := util.RandInterval(1, weightCount)
//	//
//	//var currWeight int32 = 0
//	//var nextWeight int32 = 0
//	//for _, astrola := range randomPool {
//	//	nextWeight += astrola.Weight
//	//	if randomWeight > currWeight && randomWeight <= nextWeight {
//	//		return astrola
//	//	}
//	//	currWeight = nextWeight
//	//}
//	//return nil
//}

type RandomData struct {
	data []int32
}

func (this *RandomData) randomData() (int32, bool) {
	randomLen := len(this.data)
	if randomLen == 0 {
		return 0, false
	}
	randomIndex := rand.Intn(randomLen)
	result := this.data[randomIndex]
	this.data = append(this.data[:randomIndex], this.data[randomIndex+1:]...)
	return result, true
}
func getArrayInt32(parmArray []int32) []int32 {
	result := []int32{}
	for _,value := range parmArray {
		result = append(result, value)
	}
	return result
}


func ItemAddRandomWeight(starType int32) map[int32]int32 {
	//var itemsGroup []*db.DBItemRandom
	data := DATA.ITEM_RANDOM_DATA[starType]
	if data == nil {
		return nil
	}
	//randomTempData := []int32{}
	randomTempData := getArrayInt32(data.Items)
	randomData := &RandomData{data:randomTempData}
	
	result := make( map[int32]int32)
	for _, round := range DATA.ITEM_DROP_RATE {
		for i:=0; i <int(round.Number); i ++ {
			itemID, succ := randomData.randomData()
			if !succ {
				return result
			}
			result[itemID] = round.Weight
		}
	}
	return result
}

//func RandomDial() (int32, bool) {
//	randomWeight := DATA.DIAL_ID_WEIGHT_MAPPING
//	if randomWeight == nil {
//		return 0, false
//	}
//	dialId := util.RandomWeight(randomWeight)
//	return dialId, true
//}

//func RandomDialNum(min int32, max int32) int32 {
//	resultNum := min + rand.Int31n(max - min + 1)
//	return resultNum
//}

func GetEventBaseByTaskType(dialType int32) int32 {
	var taskType int32
	switch dialType {
	case constant.QUEST_ROB_FAITH:
		taskType = constant.TASK_ROB_FAITH
		break
	case constant.QUEST_ROB_BELIEVER:
		taskType = constant.TASK_ROB_BELIEVER
		break
	case constant.QUEST_ATT_BUILDING:
		taskType = constant.TASK_ATK_BUILDING
		break
	default:
		exception.GameException(exception.DIAL_NOT_FOUND)
	}
	return taskType
}


/*

			for i=0;i<turnNumber;i++{
		for j=0;j<DATA.ITEM_DROP_RATE[i+1].Number;j++{

			for int32(len(nums)) < DATA.ITEM_DROP_RATE[i+1].Number*(i+1) {
					index := rand.Int31n(itemslen)
					exist := false
					for _,v := range nums{
						if v == index{
							exist = true
							break;
						}
					}
				if !exist{
					nums = append(nums,index)
				}
				results[index] = DATA.ITEM_DROP_RATE[i+1].Weight
				//fmt.Print(i+1,DATA.ITEM_DROP_RATE[i+1].Weight)
			}
		}
		//itemslen -= DATA.ITEM_DROP_RATE[i+1].Number
	}

	for k,v := range results{
		fmt.Printf("%d,%d\n",k,v)
	}

	return results
}
*/
//var items [100]int32
//itemsNumber := int32(len(DATA.ITEM_ID))//圣物数量40
//var itemsGroup []*db.DBItemRandom
//itemsGroup := &db.DBItemRandom{}
/*
for _,v := range starItems.Items {

}
*/
//itemsGroups := &db.DBItemRandom{}
//util.
//for ; i < itemsNumber ; i++{
//	 items[i] = DATA.ITEM_ID[i]
//}
/*
			itemsNumber := int32(len(DATA.ITEM_ID))
			itemsGroup[i].Weight = DATA.ITEM_DROP_RATE[i].Weight

			index := rand.Int31n(itemsNumber)
			//itemsGroup.items[index]
			itemsGroup[i].ItemId[j] = DATA.ITEM_ID[index]
			delete(DATA.ITEM_ID,index)
			*/

/*
	for i=0; i<8;i++  {
		index := rand.Int31n(40-i)
		itemsGroup.Weight = 50
		//itemsGroup.ItemId[i]
		itemsGroup.ItemId[i] =items[index]
		//items[index]
	}
	for i=0; i<8;i++  {
		index := rand.Int31n(35-i)
		itemsGroup.Weight = 50
		//itemsGroup.ItemId[i]
		itemsGroup.ItemId[i] =items[index]
		//items[index]
	}
	for i=0; i<8;i++  {
		index := rand.Int31n(-i)
		itemsGroup.Weight = 50
		//itemsGroup.ItemId[i]
		itemsGroup.ItemId[i] =items[index]
		//items[index]
	}
*/
func GenHelpBeliever(num int32) []*protocol.BelieverInfo{
	maleID := "b0161"
	femaleID := "b0162"
	addBeliever := make([]*protocol.BelieverInfo, 2)
	maleNum := int32(num/2)
	femaleNum := num - maleNum
	addBeliever[0] = &protocol.BelieverInfo{Id:maleID, Num:maleNum}
	addBeliever[1] = &protocol.BelieverInfo{Id:femaleID, Num:femaleNum}
	return addBeliever
}


