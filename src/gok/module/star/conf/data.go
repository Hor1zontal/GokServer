package conf

import (
	"encoding/json"
	"strings"
	"aliens/log"
	util2 "aliens/common/util"
	"gok/module/cluster/center"
)

const (
	CONFIG_ARRAY_SPLIT string = "|"
)


func Close() {
}

func Init() {
	center.ConfigCenter.SubscribeConfig("gamebase", updateGameBaseData)
	center.ConfigCenter.SubscribeConfig("item", updateItemData)
	center.ConfigCenter.SubscribeConfig("itemgroup", updateItemGroupData)
	center.ConfigCenter.SubscribeConfig("buff", updateBuffBase)
	center.ConfigCenter.SubscribeConfig("believerupgrade", updateBelieverUpgradeBase)
	center.ConfigCenter.SubscribeConfig("believerupgrademodulus", updateCivilizationRatioBase)


	center.ConfigCenter.SubscribeConfig("robotrule", updateRobotRuleBase)
	center.ConfigCenter.SubscribeConfig("civilization", updateTaskCivilizationData)
	center.ConfigCenter.SubscribeConfig("gameobject", updateGameObjectBase)
	center.ConfigCenter.SubscribeConfig("buildingsbase", updateBuildingBaseData)
	center.ConfigCenter.SubscribeConfig("starbase", updateStarBaseData)
	center.ConfigCenter.SubscribeConfig("believerbuildexpedite", updateBelieverBuildingExp)
	center.ConfigCenter.SubscribeConfig("reward", updateCivilizationRewardData)

	center.ConfigCenter.SubscribeConfig("itemgroupunlock", updateItemGroupUnlockData)
	center.ConfigCenter.SubscribeConfig("believerprice", updateBelieverPriceBase)
	center.ConfigCenter.SubscribeConfig("flagunlock",updateFlagUnlockData)
}

//
//type WeightData interface {
//	GetWeight() int32
//}
//
//func RandomWeightData(weightMapping []WeightData) WeightData {
//	var totalWeight int32 = 0
//	for _, weightData := range weightMapping {
//		weight := weightData.GetWeight()
//		if weight <= 0 {
//			continue
//		}
//		totalWeight += weight
//	}
//	if (totalWeight <= 0) {
//		return nil
//	}
//	randomValue := rand.Int31n(totalWeight) + 1
//	var currentValue int32 = 0
//	for _, weightData := range weightMapping {
//		weight := weightData.GetWeight()
//		if weight <= 0 {
//			continue
//		}
//		currentValue += weight
//		if currentValue >= randomValue {
//			return weightData
//		}
//	}
//	return nil
//}

var Base struct {
	*GameBase
	Star                map[int32]*StarBase
	Building            map[int32]*BuildingBase
	GameObject          map[string]*GameObjectBase       //星球对象基础配置
	BelieverUpgrade     map[string][]*BelieverUpgradeBase  //合成配置
	CivilizationRatio	map[int32]*CivilizationRatioBase //文明度加成系数
	CivilizationReward	map[int32]*CivilizationReward	//文明度等级奖励

	BLimit              map[int32]int32                 //建筑的最大等级
	BelieverBuildingExp []*BelieverBuildingExpedite     //信徒建筑加成表
	Item                map[int32]*ItemBase
	Buff                map[int32]*BuffBase
	RobotRule           []*RobotRuleBase


	CivilMaxSeq   int32
	CivilData     map[int32]map[int32]*Civilization //星球序号-文明度等级-当前等级阈值
	CivilMaxLevel map[int32]int32                   //文明度最大等级

	ItemGroup        map[int32]*ItemGroupBase
	StarGroupMapping map[int32][]*ItemGroupBase //星球ID - 圣物组合
	GroupTypeMapping map[int32]int32 // groupID - buffType
	GroupUnlockBase  map[int32]*ItemGroupUnlockBase //
	GroupUnlockRequire []*ItemGroupUnlockBase

	CivilBelieverMapping map[int32]int32
	BelieverFaithMapping map[int32]int32

	FlagKeyLevelMapping map[int32]int32
}


//type GroupUnlockNum struct {
//	GroupNum		int32
//	BuildingNum		int32
//}

type BelieverBuildingExpedite struct {
	BuildingLevel int32 `json:"buildingLevel"`
	BelieverLevel int32 `json:"believerLevel"`
	DecreaseTime  int32 `json:"decreaseTime"`
}

type Civilization struct {
	Seq       int32 `json:"starOrdinal"`
	Level     int32 `json:"level"`
	Threshold int32 `json:"threshold"`
	BuffID    int32 `json:"buff"`
	Reward    int32 `json:"reward"`
	Gift	  int32 `json:"gift"`
}


type RobotRuleBase struct {
	Level         int32   `json:"level"`      //等阶号
	RandomType    []int32 `json:"randomType"` //1抢信仰任务 2抢信徒任务 3拆建筑任务
	Weight        int32   `json:"weight"`     //权重
	BuildingType  []int32 `json:"buildingType"`
	BuildingLevel []int32 `json:"buildingLevel"`
	BelieverLevel []int32 `json:"believerLevel"`
	BelieverNum   []int32 `json:"believerNum"`
}

type ItemGroupUnlockBase struct {
	ID			  	int32 	`json:"id"`
	UnlockType 		int32 	`json:"unlock"`
	BuildingRequire int32 `json:"buildingRequire"`
	CivilizationIncome int32  `json:"civilizationIncome"`
}

//建筑基础表
type BuildingBase struct {
	ID                   int32   //编号
	BuildID              int32   //建筑类型
	Level                int32   //建筑等级
	StarID               int32   //所属星球编号
	BuildTime            int32   //建造时长（分钟）
	RepairTime           int32   //修复时长（分钟）
	PowerAcquired        int32   //获得神力点数
	UpgradeConsumption   float32 //建造消耗信仰值数量
	RepairConsumption    int32   //修复建筑消耗信仰值数量
	BuiAttSucAcquired    int32   //建造中攻击成功获得信仰值
	BuiAttFaiAcquired    int32   //建造中攻击失败获得信仰值
	RepAttSucAcquired    int32   //修复中攻击成功获得信仰值
	RepAttFaiAcquired    int32   //修复中攻击失败获得信仰值
	DoneAttSucAcquired   int32   //完成后攻击成功获得信仰值
	DoneAttFaiAcquired   int32   //完成后攻击失败获得信仰值
	PowerLimit           int32   //建筑增加的法力值上限
	PowerReward          int32   //建筑获得的法力值
	BoomTime             int32   //建筑爆炸CD时间
	UpgradeBelieverLevel int32   //建造需要消耗的信徒的等级
	//UpgradeBelieverNumber   int32  //建造需要消耗的信徒的数量
	RepairBelieverLevel int32 //修复需要消耗信徒
	//RepairBelieverNumber 	int32  //修理建筑需要消耗的信徒数量

	ExpediteUpgradeBelieverLimit int32            //加速建造最多使用的信徒数量
	ExpediteRepairBelieverLimit  int32            //加速修理最多使用的信徒数量
	UpdateFaithTime              int32            //建筑更新一次信仰值的时间
	UpdateFaithNum               int32            //建筑更新一次的信仰值
	FaithLimit                   int32            //建筑存储的信仰值上限
	ReceiveFaithPercent          int32            //领取信仰值的最小百分比
	ReceiveFaithMin              int32 `json:"-"` //领取信仰值的最小值

	RequireCivilizationLevel int32 //建造的文明等级的限制
	CivBuildIncome           int32 //建造完成后可获得文明度
	CivRepairIncome          int32 //维修完成后可获得文明度
}

type BelieverUpgradeBase struct {
	Id       int32 `json:"id"`
	SelectID string //选中的信徒
	MatchID  string //匹配的信徒
	//UpgradeID	string     		//升级的id
	Cost                     int32 //信仰花费
	RequireCivilizationLevel []int32 //需要的最小文明等级
	RequireBuildingLevel []int32

	CivilizationIncome []int32
	UpgradeNum         []int32
	UpgradeID          []string
	Weight             []int32
	//UpgradeID 			string
	//Weight 				int32
	//Num 				int32
	//CivBuildIncome  int32

	//UpgradeResult []string   `json:"-"`	//升级的结果
	MaxWeight int32 `json:"-"`  //最大的权重
	RandomResult map[int32]util2.WeightData `json:"-"`
	RandomRequire map[int32]*UpgradeRequire `json:"-"`
}

type CivilizationRatioBase struct {
	Level   int32      `json:"level"`
	Modulus []float64  `json:"modulus"`
}

type UpgradeRequire struct {
	RequireCivilLevel 		int32
	RequireBuildingLevel    int32
}

type UpgradeResult struct {
	UpgradeID          []string
	Weight             int32
	Num                int32
	CivilizationIncome int32
}

func (this *UpgradeResult) GetWeight() int32 {
	return this.Weight
}

const (
	GAMEOBJECT_TYPE_NPC      int32 = 1
	GAMEOBJECT_TYPE_ARTICLE  int32 = 2
	GAMEOBJECT_TYPE_BELIEVER int32 = 3
	GAMEOBJECT_TYPE_ENEMY    int32 = 4
)

type GameObjectBase struct {
	Id   int32 `json:"id"` //
	ID   string            //对象id
	Type int32             //对象类型  "1-剧情NPC 2-场景物件 3-信徒 4-怪物"
	LV   int32             //对象等级
	Sex  int32             //对象性别 "0-无性别 1-男性 2-女性"
	Hp   int32             //对象血量
	Atk  int32             //对象攻击力
}

//游戏基础数据配置表
type GameBase struct {
	BelieverBuffInterval float64 //生成信徒的时间间隔
	StarBelieverLimit    int32   //星球的信徒数量上限
	//InitPower			 int32	 //出事法力值
	InitPowerLimit		 int32   //初始法力上限
	//ExpediteBuildTime int32   //消耗一个五级信徒加速的建造时间
	//ExpediteRepairTime int32  //消耗一个无级信徒加速的修理时间
	RelicFaithPowerPercent float64 //一个圣物加成的速率
	GrooveEffectTime       float64 //槽的CD时间
	BelieverAddEffectTime  float64 //信徒增加后缩短的生效时间
	StarWeCanArrive        []int32 //新账号的随机星球列表
	EggTouchReduceTime	   int32
	EggTouchCost 		   int32

	EggActivation			int64 //蛋首次激活信徒刷新时间（秒）
	EggActivationBeliever	int64 //蛋首次激活信徒刷新个数

	FirstGroupRelicDrop		[]int32
}

type BuffBase struct {
	ID    int32   `json:"id"`      //buff id
	Type  int32   `json:"buffType"` //buff 类型
	Ratio float32 `json:"buffNum"`  //buff 系数
}

//星球基础表
type StarBase struct {
	Type       int32               //星球类型
	Name       string              //星球名称
	Consume    int32               //前置id
	BelieverID string              //星球能够产生的平民 用|分割
	Believers  []string `json:"-"` //星球能够产生的平民 分割好后的数组
}

type ItemBase struct {
	ID       int32 `json:"id"`       //物品id
	Type     int32 `json:"type"`     //物品类型
	Color    int32 `json:"color"`    //物品颜色
	StarType int32 `json:"starType"` //星球类型
	GetWay   int32 `json:"getWay"`   //获取途径
	BuffID   int32 `json:"buffID"`   //buffID
	//GetWayID    int32  `json:"getWayID"` 	 //获取
}

type ItemGroupBase struct {
	ID 			    int32  `json:"id"`       	//图鉴id
	StarType 		int32  `json:"starType"`   //星球类型
	Rarity 			int32  `json:"rarity"`      //稀有度
	Reward 			int32  `json:"reward"`      //奖励信仰值
	BuffID 			int32  `json:"buffID"`      //buff
	Content 		[]int32  `json:"content"`  //图鉴的组成物品
	//CivilizationIncome int32 `json:"civilizationIncome"` //第一次完成的文明度奖励
}

type CivilizationReward struct {
	ID				int32 `json:"id"`			//文明度等级id
	LinkName		int32 `json:"linkName"` 	//文明度等级
	FaithNum		int32 `josn:"faithNum"`		//信仰数
	DiamondNum		int32 `json:"diamondNum"`	//钻石数
	RelicPointNum	int32 `json:"relicPointNum"`//圣物碎片数
	BelieverLevel	[]int32 `json:"believerLevel"`//信徒等级
	BelieverNum		[]int32 `json:"believerNum"`	//圣徒数量
}

type BelieverPrice struct {
	ID 				int32 `json:"id"`
	BelieverLevel	int32 `json:"beliLevel"`
	CiviliLevel		int32 `json:"civLevel"`
	Faith			int32 `json:"faith"`
}

type FlagUnlock struct {
	ID				int32 `json:"id"`
	FlagKey  		int32 `json:"flagValue"`
	BuildingLevel 	int32 `json:"buildingLevel"`
}

//func updateItemGroupData(data []byte) {
//	var datas []*ItemGroupBase
//	json.Unmarshal(data, &datas)
//	results := make(map[int32]*ItemGroupBase)
//	itemGroupMapping := make(map[int32]int32)
//	for _, data := range datas {
//		results[data.ID] = data
//		for _, itemID := range data.Content {
//			itemGroupMapping[itemID] = data.ID
//		}
//	}
//	Base.ItemGroup = results
//	Base.StarGroupMapping = itemGroupMapping
//}

func updateFlagUnlockData(data []byte) {
	var datas []*FlagUnlock
	json.Unmarshal(data, &datas)
	results := make(map[int32]int32)
	for _, data := range datas {
		results[data.FlagKey] = data.BuildingLevel
	}
	Base.FlagKeyLevelMapping = results
}

func updateItemGroupData(data []byte) {
	var datas []*ItemGroupBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]*ItemGroupBase)
	starGroupMapping := make(map[int32][]*ItemGroupBase)
	groupTypeMapping := make(map[int32]int32)
	for _, data := range datas {
		results[data.ID] = data
		mapping := starGroupMapping[data.StarType]
		if (mapping == nil) {
			mapping = []*ItemGroupBase{}
		}
		mapping = append(mapping, data)
		starGroupMapping[data.StarType] = mapping
		groupTypeMapping[data.ID] = data.BuffID
	}
	Base.ItemGroup = results
	Base.StarGroupMapping = starGroupMapping

}

func updateItemGroupUnlockData( data []byte) {
	var datas []*ItemGroupUnlockBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]*ItemGroupUnlockBase)
	for _,data := range datas {
		results[data.ID] = data
	}
	Base.GroupUnlockRequire = datas
	Base.GroupUnlockBase = results
}

func updateBelieverPriceBase(data []byte) {
	var datas []*BelieverPrice
	err := json.Unmarshal(data, &datas)
	if err != nil {
		log.Debug("believerPrice unmarshal error: %v",err.Error())
	}
	result1 := make(map[int32]int32)
	result2 := make(map[int32]int32)
	for _, data := range datas {
		result1[data.CiviliLevel] = data.BelieverLevel
		result2[data.BelieverLevel] = data.Faith
	}
	Base.CivilBelieverMapping = result1
	Base.BelieverFaithMapping = result2
}


func updateCivilizationRewardData( data []byte ){
	var datas []*CivilizationReward
	json.Unmarshal(data, &datas)
	results := make(map[int32]*CivilizationReward)
	for _, data := range datas {
		results[data.ID] = data
	}
	Base.CivilizationReward = results
}

func updateBuildingBaseData(data []byte) {
	var datas []*BuildingBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]*BuildingBase)
	for _, data := range datas {
		data.ReceiveFaithMin = data.FaithLimit * data.ReceiveFaithPercent / 100
		results[data.ID] = data
	}
	Base.Building = results;
	InitBuildingLimit() //重新统计建筑等级上限
}

func updateGameBaseData(data []byte) {
	var datas = &GameBase{}
	json.Unmarshal(data, datas)
	Base.GameBase = datas
}

func updateTaskCivilizationData(data []byte) {
	var datas []*Civilization
	json.Unmarshal(data, &datas)
	results := make(map[int32]map[int32]*Civilization)

	civilMaxLevel := make(map[int32]int32)
	for _, data := range datas {
		if data.Seq > Base.CivilMaxSeq {
			Base.CivilMaxSeq = data.Seq
		}

		currLevel := civilMaxLevel[data.Seq]
		if data.Level > currLevel {
			civilMaxLevel[data.Seq] = data.Level
		}
		levelMapping := results[data.Seq]
		if levelMapping == nil {
			levelMapping = make(map[int32]*Civilization)
			results[data.Seq] = levelMapping
		}
		levelMapping[data.Level] = data
	}
	Base.CivilData = results
	Base.CivilMaxLevel = civilMaxLevel
}

func updateRobotRuleBase(data []byte) {
	var datas []*RobotRuleBase
	json.Unmarshal(data, &datas)
	Base.RobotRule = datas
}


func updateCivilizationRatioBase(data []byte) {
	var datas []*CivilizationRatioBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]*CivilizationRatioBase)
	for _, data := range datas {
		results[data.Level] = data
	}
	Base.CivilizationRatio = results
}

func updateBelieverUpgradeBase(data []byte) {
	var datas []*BelieverUpgradeBase
	json.Unmarshal(data, &datas)
	results := make(map[string][]*BelieverUpgradeBase)
	for _, data := range datas {
		dataLen := len(data.CivilizationIncome)
		if len(data.UpgradeID) != dataLen || len(data.Weight) != dataLen || len(data.UpgradeNum) != dataLen {
			log.Error("believer upgrade config length not match! %v", data.Id)
			continue
		}
		upgradeResult := make(map[int32]util2.WeightData)
		upgradeRequire := make(map[int32]*UpgradeRequire)
		var maxWeight int32 = 0

		for i := 0; i < dataLen; i++ {
			upgradeResult[data.Weight[i]] = &UpgradeResult{
				UpgradeID:          strings.Split(data.UpgradeID[i], CONFIG_ARRAY_SPLIT),
				Weight:             data.Weight[i],
				Num:                data.UpgradeNum[i],
				CivilizationIncome: data.CivilizationIncome[i],
			}
			upgradeRequire[data.Weight[i]] = &UpgradeRequire{
				RequireCivilLevel: 		data.RequireCivilizationLevel[i],
				RequireBuildingLevel:   data.RequireBuildingLevel[i],
			}
			if data.Weight[i] > maxWeight {
				maxWeight = data.Weight[i]
			}
		}

		data.MaxWeight = maxWeight
		data.RandomResult = upgradeResult
		data.RandomRequire = upgradeRequire
		//data.UpgradeResult = strings.Split(data.UpgradeID, CONFIG_ARRAY_SPLIT)
		results[data.SelectID+data.MatchID] = append(results[data.SelectID+data.MatchID], data)
	}
	Base.BelieverUpgrade = results
}

func updateGameObjectBase(data []byte) {
	var datas []*GameObjectBase
	json.Unmarshal(data, &datas)
	results := make(map[string]*GameObjectBase)
	for _, data := range datas {
		if (data.Type == GAMEOBJECT_TYPE_BELIEVER || data.Type == GAMEOBJECT_TYPE_ENEMY) {
			results[data.ID] = data
		}
	}
	Base.GameObject = results;
}

func updateStarBaseData(data []byte) {
	var datas []*StarBase
	json.Unmarshal(data, &datas)
	resluts := make(map[int32]*StarBase)
	for _, data := range datas {
		data.Believers = strings.Split(data.BelieverID, CONFIG_ARRAY_SPLIT)
		resluts[data.Type] = data
	}
	Base.Star = resluts
}

func updateBelieverBuildingExp(data []byte) {
	var datas []*BelieverBuildingExpedite
	json.Unmarshal(data, &datas)

	Base.BelieverBuildingExp = datas;
}

func updateBuffBase(data []byte) {
	var datas []*BuffBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]*BuffBase)

	for _, data := range datas {
		results[data.ID] = data
	}
	Base.Buff = results
}

func updateItemData(data []byte) {
	var datas []*ItemBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]*ItemBase)

	for _, data := range datas {
		results[data.ID] = data
	}
	Base.Item = results
}

func (this *ItemGroupBase) ContainsItems(items []int32) bool {
	for _, item := range items {
		if (!Contains(this.Content, item)) {
			return false
		}
	}
	return true
}
func (this *ItemGroupBase) ContainsItemsNum(items []int32) int32 {
	var num int32 = 0
	for _, item := range items {
		if (Contains(this.Content, item)) {
			num++
		}
	}
	return num
}

func (this *ItemGroupBase) ContainsItem(item int32) bool {
	return Contains(this.Content, item)
}

func (this *ItemGroupBase) IsFinish(items []int32) bool {
	for _, item := range this.Content {
		if (!Contains(items, item)) {
			return false
		}
	}
	return true
}


