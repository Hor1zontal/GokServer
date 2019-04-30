/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/7/12
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package conf

import (
	"encoding/json"
	"math/rand"
)

var moduleConfigLoaders map[ModuleEnum]func(data []byte)interface{}

//获取事件步骤配置
func GetEventStepConfig(eventType int32, step int32) *StepData {
	eventData := DATA.EVENT_DATA[eventType]
	if (eventData == nil) {
		return nil
	}
	if (eventData.Steps == nil || len(eventData.Steps) < int(step)) {
		return nil
	}
	return eventData.Steps[step - 1]
}


func GetEventName(eventType int32) string {
	eventData := DATA.EVENT_DATA[eventType]
	if (eventData == nil) {
		return ""
	}
	return eventData.Name
}

type ModuleEnum int32

const (
	MODULE_STORYBOARD 	ModuleEnum = 100
	//MODULE_DISPLAY	 	ModuleEnum = 101
	//MODULE_SELECT_AREA  ModuleEnum = 101
	MODULE_GET_FAITH    ModuleEnum = 101 //获取信仰
	MODULE_GET_BELIEVER ModuleEnum = 102 //获取信徒
	MODULE_LOOT_FAITH   ModuleEnum = 103 //抢夺信仰
	MODULE_ATTACK_BUILD     ModuleEnum = 104 //攻击建筑
	MODULE_LOOT_BELIEVER    ModuleEnum = 105 //抢夺信徒
	MODULE_RANDOM_TARGET	ModuleEnum = 1
	//MODULE_RECRUIT	 	ModuleEnum = 2
	//MODULE_BUFF	  	ModuleEnum = 3
	//MODULE_SWITCH_STAR	ModuleEnum = 4
	//MODULE_ATTACK		ModuleEnum = 5
	//MODULE_VOTE		ModuleEnum = 6
	MODULE_CAPTURE		ModuleEnum = 7
	MODULE_CUSTOM_DATA	ModuleEnum = 8
)

//初始化模块数据加载句柄
func initLoader() {
	moduleConfigLoaders =  make(map[ModuleEnum]func(data []byte)interface{})
	moduleConfigLoaders[MODULE_RANDOM_TARGET] = loadConfigRandomTarget
	//moduleConfigLoaders[MODULE_RECRUIT] = loadConfigRecruit
	//moduleConfigLoaders[MODULE_BUFF] = loadConfigBuff
	//moduleConfigLoaders[MODULE_SWITCH_STAR] = loadConfigSwitchStar

	moduleConfigLoaders[MODULE_GET_FAITH] = loadConfigGetFaith
	moduleConfigLoaders[MODULE_GET_BELIEVER] = loadConfigGetBeliever
	moduleConfigLoaders[MODULE_LOOT_FAITH] = loadConfigGetFaith
	moduleConfigLoaders[MODULE_ATTACK_BUILD] = loadConfigAttackBuild
	moduleConfigLoaders[MODULE_LOOT_BELIEVER] = loadConfigLootBeliever

	//moduleConfigLoaders[MODULE_SELECT_AREA] = loadConfigSelectArea

	//moduleConfigLoaders[MODULE_VOTE] = loadConfigVote
	moduleConfigLoaders[MODULE_STORYBOARD] = loadConfigStoryboard
	moduleConfigLoaders[MODULE_CAPTURE] = loadConfigCaptureBeliever
	moduleConfigLoaders[MODULE_CUSTOM_DATA] = loadConfigCustomData
}

//根据模块id和配置数据加载配置对象
func loadData(module ModuleEnum, data []byte) interface{} {
	loader := moduleConfigLoaders[module]
	if (loader != nil) {
		return loader(data)
	}
	return nil
}

//--------------------------module 1----------------------------------
func loadConfigRandomTarget(data []byte) interface{} {
	var newData *DataRandomTarget
	json.Unmarshal(data, &newData)
	return newData
}

//随机目标模块
type DataRandomTarget struct {
	Friend		int32 	`json:"friend"`		//0 忽略此过滤条件  1 必须是好友 	 2不能是好友
	Building	int32 	`json:"building"`	//0 忽略此过滤条件  1 有建筑 	 2不能有建筑
	Repair		int32 	`json:"repair"`		//0 忽略此过滤条件  1 有修理中的建筑 2不能有修理中的建筑
	Broken		int32 	`json:"broken"`		//0 忽略此过滤条件  1 有损坏的建筑   2不能有损坏的建筑
	Believer	int32 	`json:"believer"`   //0 忽略此过滤条件  1 信徒超过目标   2信徒没超过目标
}

//--------------------------module 2----------------------------------
func loadConfigRecruit(data []byte) interface{} {
	var newData *DataRecruit
	json.Unmarshal(data, &newData)
	return newData
}

type DataRecruit struct {
	RecruitLimit int32 `json:"recruitLimit"`		//招募人数上限
}

//--------------------------module 3----------------------------------
func loadConfigBuff(data []byte) interface{} {
	var newData *DataBuff
	json.Unmarshal(data, &newData)
	return newData
}

type DataBuff struct {
	BuffFields    []*BuffField  `json:"buffFields"`		//buff字段
}

type BuffField struct {
	Field		string 	`json:"field"`		    //buff字段名  可以是客户端用来表现的id
	RefField 	string	`json:"refField"`  	    //buff字段关联的字段名,  当前例子为:1个mob在30内产生1个dump
	Interval	int32	`json:"interval"`     	//buff间隔时间
	ChangeValue	int32	`json:"changeValue"`   	//buff间隔时间内变更的值,正数为增加,负数为减少
	InitialType int32 	`json:"initialType"`   	//buff字段初始化类型，0取initial字段，1招募+当前用户人数
	Initial		int32	`json:"initial"`       	//buff字段初始值
	//Condition	int32   `json:"condition"`      //字段条件  0无 1 达到阈值作为模块结束条件(|或条件) 2 达到阈值作为模块结束条件(&与条件)
	Threshold	*int32	`json:"threshold"`      //模块结束的buff字段阈值 无，代表不用监控阈值
}

//--------------------------module 101----------------------------------
//func loadConfigSelectArea(data []byte) interface{} {
//	var newData *DataSwitchStar
//	json.Unmarshal(data, &newData)
//	return newData
//}
//
//type DataSwitchStar struct {
//	StarOwner 	int32 		`json:"starOwner"`	//切换的星球对象 0 玩家自己   1目标玩家
//}

func GetInt32Mapping(keys []int32, values []int32) map[int32]int32 {
	result := make(map[int32]int32)
	if len(keys) < len(values) {
		return result
	}
	for index, value := range values {
		result[keys[index]] = value
	}
	return result
}

func GetFloat32Mapping(keys []float32, values []int32) map[float32]int32 {
	result := make(map[float32]int32)
	if len(keys) < len(values) {
		return result
	}
	for index, value := range values {
		result[keys[index]] = value
	}
	return result
}


//--------------------------module 101------------------------------------
func loadConfigGetFaith(data []byte) interface{} {
	var newData *DataGetFaith
	json.Unmarshal(data, &newData)

	faithScopeMapping := make(map[int32][]int32)
	maxLen := len(newData.FaithInterval)
	for index, level := range newData.BelieverLevel {
		if index >= maxLen {
			break
		}
		faithScopeMapping[level] = newData.FaithInterval[index]
	}

	newData.FaithScopeMapping = faithScopeMapping
	//newData.BelieverWeightMapping = GetInt32Mapping(newData.BelieverLevel, newData.BelieverWeight)
	//newData.FaithRatioMapping = GetInt32Mapping(newData.FaithRatio, newData.FaithWeight)
	return newData
}

type DataGetFaith struct {
	//FaithRatio 	    []int32 		`json:"faithRatio"`	    //奖励信仰  随机出的用户信徒等级*系数
	//FaithWeight 	[]int32 		`json:"faithWeight"`	//随机到的信仰权重
	BelieverLevel 	[]int32 		`json:"believerLevel"`	//信徒等级
	FaithInterval 	[][]int32 		`json:"faithInterval"`	//信徒等级
	//BelieverWeight 	[]int32 		`json:"believerWeight"`	//信徒等级对应的权重
	FaithScopeMapping 	map[int32][]int32 		`json:"-"`	//信徒信仰的随机区间
	//BelieverWeightMapping 	map[int32]int32 		`json:"-"`	//信徒等级随机的权重关系
}

func (this *DataGetFaith) RandomFaith(level int32) int32 {
	scope := this.FaithScopeMapping[level]
	if scope == nil || len(scope) != 2 {
		return  0
	}
	return scope[0] + rand.Int31n(scope[1] - scope[0])
}

//--------------------------module 102------------------------------------
func loadConfigGetBeliever(data []byte) interface{} {
	var newData *DataGetBeliever
	json.Unmarshal(data, &newData)
	newData.PeopleWeightMapping = GetInt32Mapping(newData.People, newData.Weight)
	newData.BelieverWeightMapping = GetInt32Mapping(newData.BelieverLevel, newData.BelieverWeight)
	return newData
}

type DataGetBeliever struct {
	People 	[]int32 		`json:"people"`	//获取到的原始人数量
	Weight 	[]int32 		`json:"weight"`	//获取到的原始人权重

	BelieverLevel 	[]int32 `json:"believerLevel"`	//获取到的原始人数量
	BelieverWeight 	[]int32 `json:"believerWeight"`	//获取到的原始人权重

	PeopleWeightMapping    map[int32]int32  `json:"-"`	//获取到的信徒数量随机权重对应关系
	BelieverWeightMapping  map[int32]int32  `json:"-"`	//获取到的信徒等级权重对应关系
}

//--------------------------module 103------------------------------------
//func loadConfigLootFaith(data []byte) interface{} {
//	var newData *DataLootFaith
//	json.Unmarshal(data, &newData)
//	return newData
//}
//
//type DataLootFaith struct {
//	ratio 	[]int32 		`json:"ratio"`	//
//}


//--------------------------module 104------------------------------------
func loadConfigAttackBuild(data []byte) interface{} {
	var newData *DataAttackBuild
	json.Unmarshal(data, &newData)
	newData.FaithRatioMapping = GetFloat32Mapping(newData.FaithRatio, newData.FaithWeight)
	return newData
}

type DataAttackBuild struct {
	SuccessRatio 	float32 		`json:"successRatio"`	//攻击成功的系数
	FaithRatio 	    []float32 		`json:"faithRatio"`	    //奖励信仰  随机出的用户信徒等级*系数
	FaithWeight 	[]int32 		`json:"faithWeight"`	//随机到的建筑信仰权重
	FaithRatioMapping 	map[float32]int32 		`json:"-"`	//信徒信仰系数的权重关系
	FailedRatio		float32 		`json:"failedRatio"`	//攻击失败的系数
}


//--------------------------module 105------------------------------------
func loadConfigLootBeliever(data []byte) interface{} {
	var newData *DataLootBeliever
	json.Unmarshal(data, &newData)
	newData.BelieverWeightMapping = GetInt32Mapping(newData.BelieverLevel, newData.BelieverWeight)
	return newData
}

type DataLootBeliever struct {
	SuccessRatio 	float32 		`json:"successRatio"`	//抢夺成功的系数
	BelieverLevel 	[]int32 		`json:"believerLevel"`	//信徒等级
	BelieverWeight 	[]int32 		`json:"believerWeight"`	//信徒等级对应的权重
	BelieverWeightMapping 	map[int32]int32 		`json:"-"`	//信徒等级随机的权重关系
}


//--------------------------module 4------------------------------------
func loadConfigSwitchStar(data []byte) interface{} {
	var newData *DataSwitchStar
	json.Unmarshal(data, &newData)
	return newData
}

type DataSwitchStar struct {
	StarOwner 	int32 		`json:"starOwner"`	//切换的星球对象 0 玩家自己   1目标玩家
}


//-------------------------module 6------------------------------

func loadConfigVote(data []byte) interface{} {
	var newData *DataVote
	json.Unmarshal(data, &newData)
	return newData
}


type DataVote struct {
	OptionNum  int  	  `json:"optionNum"`		//投票选项数量
	Options    []*VoteOption  `json:"options"`		//投票选项
}

type VoteOption struct {
	Summary  string       `json:"summary"`  //投票的摘要信息
	EventID  int32        `json:"eventID"`  //投票结束后处理的的事件id
}


//--------------------------module 100----------------------------------
func loadConfigStoryboard(data []byte) interface{} {
	var newData *DataStoryboard
	json.Unmarshal(data, &newData)
	return newData
}

type DataStoryboard struct {
	StoryId 	string 		`json:"storyId"`	//story编号，story编辑器导出的故事编号
}

//--------------------------module 7----------------------------------
func loadConfigCaptureBeliever(data []byte) interface{} {
	var newData *DataCapture
	json.Unmarshal(data, &newData)
	return newData
}

//随机目标模块
type DataCapture struct {
	CaptureLimit 	int32  	`json:"captureLimit"`	//抓捕次数
	BelieverProbs 	[]int32 `json:"believerProbs"`	//抓捕的信徒等级对应的概率   100为上限
}


//-------------------------module 8---------------------------------
func loadConfigCustomData(data []byte) interface{} {
	var newData *DataCustom
	json.Unmarshal(data, &newData)
	return newData
}

type DataCustom struct {
	Initial []*CustomData  `json:"initial"`
}

func (this *DataCustom) CopyData() []*CustomData {
	result := []*CustomData{}
	if (this.Initial != nil) {
		for _, data := range this.Initial {
			result = append(result, &CustomData{data.Key, data.Value})
		}
	}
	return result
}

type CustomData struct {
	Key 	string  `json:"key"`	//初始化数据key
	Value 	int32   `json:"value"`	//初始化数据value
}

