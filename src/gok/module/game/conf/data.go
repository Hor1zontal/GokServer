package conf

import (


	"encoding/json"
	"gok/service/msg/protocol"
	"aliens/common/character"
	"strings"
	"gok/module/cluster/center"
	"aliens/log"
	"aliens/common/util"
	"math/rand"
)


//var Server struct {
//	Enable	     	   bool
//	Database dbconfig.DBConfig
//	RemoteCacheTimeout float64
//	UserFreeTimeout    float64
//	SyncDBInterval	   float64  //用户在线数据同步到数据库的时间
//	RedisAddress       string
//	RedisPassword      string
//	RedisMaxActive   int
//	RedisMaxIdle     int
//	RedisIdleTimeout int
//	RPCAddress         string	//提供RPC服务的地址,信息需要注册到中心服务器供其他服务调用
//	RPCPort            int	//提供RPC服务的端口，本地启动RPC需要指定此端口启动
//}




func Close() {

}

func init() {

}

func Init() {
	//配置文件转为json
	//BaseParse1(basePath + "base/Initialization.xlsx", reflect.TypeOf(Initialization{}), "Initialization", true)

	//Base.Exp = make(map[int32]*ExpBase)
	//BaseParse(basePath + "base/ExpBase.xlsx", reflect.TypeOf(ExpBase{}), "Exp")

	//	Base.Building = make(map[int32]*BuildingBase)
	//	BaseParse(basePath + "base/BuildingsBase.xlsx", reflect.TypeOf(BuildingBase{}), "Building")

	//Base.Star = make(map[int32]*StarBase)
	//BaseParse(basePath + "base/StarBase.xlsx", reflect.TypeOf(StarBase{}), "Star")

	DATA.BELIEVER_SCORE_MAPPING = make(map[int32]int32)
	DATA.BELIEVER_SCORE_MAPPING[1] = 1
	DATA.BELIEVER_SCORE_MAPPING[2] = 2
	DATA.BELIEVER_SCORE_MAPPING[3] = 6
	DATA.BELIEVER_SCORE_MAPPING[4] = 30
	DATA.BELIEVER_SCORE_MAPPING[5] = 70
	DATA.BELIEVER_SCORE_MAPPING[6] = 150

	center.ConfigCenter.SubscribeConfig("item", updateItemData)
	//center.ConfigCenter.SubscribeConfig("recastcost", updateRecastCostData)
	center.ConfigCenter.SubscribeConfig("itemgroup", updateItemGroupData)
	center.ConfigCenter.SubscribeConfig("shopbase", updateShopBaseData)
	//center.ConfigCenter.SubscribeConfig("expbase", updateExpBaseData)
	center.ConfigCenter.SubscribeConfig("buildingsbase", updateBuildingsBaseData)

	center.ConfigCenter.SubscribeConfig("taskreward", updateTaskRewardData)
	center.ConfigCenter.SubscribeConfig("itemdroprate", updateItemDropRateData)
	center.ConfigCenter.SubscribeConfig("eventweightappend", updateEventWeightAppend)
	center.ConfigCenter.SubscribeConfig("dial", updateDialData)
	center.ConfigCenter.SubscribeConfig("diallimit", updateDiaLimitData)
	center.ConfigCenter.SubscribeConfig("initdial", updateInitDial)
	//center.ConfigCenter.SubscribeConfig("itemgroupunlock", updateItemGroupUnlockData)

	//事件--------------------------------------------------------
	//center.ConfigCenter.SubscribeConfig("eventfilterbase", updateAstrolaTaskData)
	center.ConfigCenter.SubscribeConfig("eventfilterbase",updateEventFilterBase)
	center.ConfigCenter.SubscribeConfig("eventbase", updateEventData)
	center.ConfigCenter.SubscribeConfig("taskbase", updateEventTaskData)
	center.ConfigCenter.SubscribeConfig("gamebase", updateGameBaseData)
	center.ConfigCenter.SubscribeConfig("itembase", updateItemBaseData)

	center.ConfigCenter.SubscribeConfig("astroladirectreward", updateCardReward)



	center.ConfigCenter.SubscribeConfig("robotsetting", updateRobotSetting)


}


var DATA struct {
	*ItemHelpBase
	*GameBase //wjl 20170614 获取游戏基础数据
	//RECAST_DATA            map[int]*RecastCostBase
	EVENT_TASK_MAPPING     map[int32]int32 //事件生成与任务关联映射关系
	SHOP_BASEDATA          map[int32]*ShopBase // wjl 20170619 游戏的商店基础数据
	//EXP_DATA               map[int32]*ExpBase
	BUILDING_DATA          map[int32]*BuildingBase
	ASTROLA_DATA           map[int32][]*AstrolaTask
	EVENT_FILTER_DATA	   map[int32]*EventFilterBase
	EVENT_ID_WEIGHT_MAPPING	map[int32]int32
	TASK_REWARD_DATA       map[int32]map[int32]*TaskRewardBase
	ITEM_RANDOM_DATA       map[int32]*ItemRandom  //星球类型-随机
	IITEM_REWARD  		   map[int32]int32
	ITEM_GROUP             map[int32]*ItemGroupBase
	STAR_ITEMGOURP_MAPPING map[int32][]*ItemGroupBase
	EVENT_WEIGHT_APPEND    map[int32]*EventWeightAppend

	DIAL_DATA                     map[int32]*DialData
	//DIAL_ID_WEIGHT_MAPPING        map[int32]int32
	//DIAL_ID_REQUIRE_LEVEL_MAPPING map[int32]int32
	DIAL_CIVILLEVEL_DIALID_WEIGHT_MAPPING	map[int32]map[int32]int32
	DIAL_LIMIT_DATA_MAPPING		  map[int32]map[int32]*DialLimitData

	GUIDE_DIAL_LEVEL_ID_WEIGHT_MAPPING map[int32]map[int32]int32 // level - id - weight
	GUIDE_ID_DIALID_MAPPING            map[int32]int32           // id - dialID
	GUIDE_MAX_REQUIRE_LEVEL			int32
	GUIDE_MIN_REQUIRE_LEVEL			int32

	EVENT_CARD_MAPPING map[int32]map[int32]util.WeightData
	DEFAULT_CARD_FAITH_SCOPE  map[int32][]int32

	//ITEM_ID		   map[int32]int32//四十个圣物的ID
	ITEM_DROP_RATE		   []*ItemDropRate //轮次-

	//BelieverUpgrade map[string]*BelieverUpgradeBase	//信徒合成配置
	BELIEVER_SCORE_MAPPING 		map[int32]int32
	//GROUP_UNLOCK_BASE 			map[int32]*ItemGroupUnlockBase //id - 条件的映射关系

	//事件--------------------------------------------------------
	EVENT_DATA 	 map[int32]*EventData
	EVENT_TASK_DATA  map[int32]*EventTaskData

	ROBOT_MAPPPING map[int32]string
}

func (this *ItemRandom) AddWeight(itemID int32, buildingWeight int32, taskWeightWeight int32, believer int32) {
	if this.BuildingWeight == nil {
		this.BuildingWeight = make(map[int32]int32)
	}
	if this.TaskWeight == nil {
		this.TaskWeight = make(map[int32]int32)
	}
	if this.BelieverWeight == nil {
		this.BelieverWeight = make(map[int32]int32)
	}
	if this.Items == nil {
		this.Items = []int32{}
	}

	this.BelieverWeight[itemID] = believer
	this.BuildingWeight[itemID] = buildingWeight
	this.TaskWeight[itemID] = taskWeightWeight
	//append(this.Items[],itemID)
	this.Items = append(this.Items, itemID)
}


type CardReward struct {
	ID 				int32 	 `json:"ID"`
	EventType	    int32	 `json:"EventType"`
	RewardType		int32	 `json:"RewardType"`
	Weight			int32	 `json:"Weight"`
	Num 			[]int32  `json:"Num"`
	Multiple        float32  `json:"Multiple"`
}

func (this *CardReward) GetWeight() int32 {
	return this.Weight
}

func (this *CardReward) RandomValue() int32 {
	if len(this.Num) != 2 {
		return 0
	}
	scope := this.Num[1] - this.Num[0]
	var appendValue int32 = 0
	if scope > 0 {
		appendValue = rand.Int31n(scope)
	}

	//log.Debug("random %v - %v ", this.Num, appendValue)
	return this.Num[0] + appendValue
}


//type BelieverUpgradeBase struct {
//	SelectID	string    		//选中的信徒
//	MatchID		string     		//匹配的信徒
//	//UpgradeID	string     		//升级的id
//	//Cost        int32 			//信仰花费
//	Level       int32 			//需要的文明度等级
//}
type DialData struct {
	ID 				int32 	 `json:"id"`
	//Type			int32	 `json:"type"`
	//Num				int32	 `json:"num"`
	//Weight			int32	 `json:"weight"`
	Position 		int32	 `json:"position"`
	ShareMultiple	int32	 `json:"shareMultiple"`
	AdMultiple		int32	 `json:"adMultiple"`
	//BuildingTopLevel	int32    `json:"buildingTopLevel"`
}

type DialLimitData struct {
	ID            	int32	`json:"id"`
	DialID			int32	`json:"dialID"`
	CivilLevel		int32	`json:"civilLevel"`
	Type			int32	`json:"type"`
	Num   			int32	`json:"num"`
	Weight			int32	`json:"weight"`
}

type InitDialData struct {
	ID				int32 	 `json:"id"`
	BuildingLevel	int32 	 `json:"buildingLevel"`
	Weight          int32    `json:"weight"`
	DialID          int32    `json:"dialID"`
}

type Buff struct {
	ID 			    int32    `json:"id"`       	    //BUFF ID
	Type    		int32    `json:"buffType"`      //BUFF 类型
	BuffRatio       float32  `json:"buffNum"`       //BUFF 系数
}

type EventWeightAppend struct {
	ID 			    int32  `json:"id"`         //时间ID
	Type    		int32  `json:"type"`       //0 信仰加成  1信徒加成
	Limit           int32  `json:"limit"`      //区分差异值的上限
	Util            int32  `json:"util"`       //每一次加成权重的数据单位
	Add             int32  `json:"add"`        //每一次加成的权重值
}

type ItemGroupBase struct {
	ID 			    int32  `json:"id"`       	//图鉴id
	StarType 		int32  `json:"starType"`   //星球类型
	Rarity 			int32  `json:"rarity"`      //稀有度
	Reward 			int32  `json:"reward"`      //奖励信仰值
	BuffID 			int32  `json:"buffID"`      //buff
	Content 		[]int32  `json:"content"`  //图鉴的组成物品
	CivilizationIncome int32 `json:"civilizationIncome"` //第一次完成的文明度奖励
}

type EventFilterBase struct {
	ID 				int32	`json:"ID"`
	StarType		int32	`json:"StarType"`
	EventBase		int32 	`json:"EventBase"`
	IsOnly			bool	`json:"IsOnly"`
	RewardType		int32	`json:"RewardType"`
	Weight			int32 	`json:"Weight"`
	Position 		int32 	`json:"position"`
}

//type ItemGroupUnlockBase struct {
//	ID			  	int32 	`json:"id"`
//	UnlockType 		int32 	`json:"unlock"`
//	BuildingRequire []int32 `json:"buildingRequire"`
//}

//func (this *ItemGroupBase) ContainsItems(items []int32) bool {
//	for _, item := range items {
//		if (!Contains(this.Content, item)) {
//			return false
//		}
//	}
//	return true
//}
//func (this *ItemGroupBase) ContainsItemsNum(items []int32) int32 {
//	var num int32 = 0
//	for _, item := range items {
//		if (Contains(this.Content, item)) {
//			num++
//		}
//	}
//	return num
//}
//
//func (this *ItemGroupBase) ContainsItem(item int32) bool {
//	return Contains(this.Content, item)
//}
//
//func (this *ItemGroupBase) IsFinish(items []int32) bool {
//	for _, item := range this.Content {
//		if (!Contains(items, item)) {
//			return false
//		}
//	}
//	return true
//}

type RecastCostBase struct {
	Num 		int  `json:"num"`       //编号
	MoneyType 	int32  `json:"moneyType"` 	 //货币类型
	Cost     	int32  `json:"cost"` 	 //数量
}

type ItemBase struct {
	ID 			int32  `json:"id"`       //物品id
	Type 		int32  `json:"type"` 	 //物品类型
	Color     	int32  `json:"color"` 	 //物品颜色
	StarType    int32  `json:"starType"` //星球类型
	//GetWay      int32  `json:"getWay"` 	 //获取途径

	PrBeliever  int32  `json:"prBeliever"` 	 //获取途径
	PrBuilding  int32  `json:"prBuilding"` 	 //获取途径
	PrTask      int32  `json:"prTask"` 	 //获取途径
	CivilizationIncome int32 `json:"civilizationIncome"` //第一次获得的文明度奖励

	//GetWayID    int32  `json:"getWayID"` 	 //获取
}

type ItemRandom struct {
	Items []int32

	BuildingWeight map[int32]int32
	TaskWeight map[int32]int32
	BelieverWeight map[int32]int32
}

type ItemDropRate struct{
	ID		int32 `json:"id"`		//轮次
	Weight  int32 `json:"weight"`	//权重
	Number	int32 `json:"number"`	//圣物数量
}

type EventTaskData struct {
	ID               int32 			`json:"id"`            //任务id
	RefID            int32 			`json:"refID"`         //任务关联id
	TriggerStep      int32			`json:"triggerStep"`   //任务触发的事件步骤编号,步骤开始的时候触发
	TriggerTarget    int32 		    `json:"triggerTarget"` //任务触发的时间目标 0事件发起者 1事件目标
	//Fields	 []*TaskFieldData       `json:"fields"` 	 //任务检测字段
	Results 	 []*EventTaskResult     `json:"results"`         //任务结果，可以有多个
	Reward		 *Reward		`json:"reward"` 	//任务奖励
}

//type EventTaskData struct {
//	ID               int32 			`json:"id"`            //任务id
//	RefID            int32 			`json:"refID"`         //任务关联id
//	TriggerStep      int32			`json:"triggerStep"`   //任务触发的事件步骤编号,步骤开始的时候触发
//	TriggerTarget    int32 		    `json:"triggerTarget"` //任务触发的时间目标 0事件发起者 1事件目标
//	Results 	 []*EventTaskResult          `json:"results"`         //任务结果，可以有多个
//}

type EventTaskFieldData struct {
	Field 		string `json:"field"`     //任务字段
	Threshold 	int32  `json:"threshold"` //阈值
}

type Reward struct {
	Faith 		int32  `json:"faith"` 	 	 //奖励信仰值
	EventFaith 	bool   `json:"eventFaith"` 	 //是否追加事件中获取的信仰值
	EventBeliever 	bool   `json:"eventBeliever"` 	 //知否追加事件中获取的信徒
}

type EventTaskResult struct {
	Field 			string `json:"field"`            //监控字段
	Threshold 		int32  `json:"threshold"` 	 //阈值
	RewardFaith     	int32  `json:"rewardFaith"` 	 //奖励信仰值
	RewardFieldFaith        string `json:"rewardFieldFaith"` //奖励信仰值,取字段变量值
}

type ExpBase struct {
	ID    int32
	Consumption int32
}

type TaskRewardBase struct {
	TaskID	int32
	EndingID int32
	FaithNum int32
	BelieverLevel string
	BelieverNum   string
	GetRelic int32
	//ItemID	int32
	//ItemNum int32

	BelieverInfo  []*protocol.BelieverInfo

	//ItemMapping 	 map[int32]int32   //奖励的物品信息
}

//建筑基础表
type BuildingBase struct {
	ID 	           	   int32  //编号
	BuildID            int32  //建筑类型
	Level              int32  //建筑等级
	StarID             int32  //所属星球编号
	BuildTime          int32  //建造时长（分钟）
	RepairTime         int32  //修复时长（分钟）
	PowerAcquired      int32  //获得神力点数
	UpgradeConsumption int32  //建造消耗信仰值数量
	RepairConsumption  int32  //修复建筑消耗信仰值数量
	BuiAttSucAcquired  int32  //建造中攻击成功获得信仰值
	BuiAttFaiAcquired  int32  //建造中攻击失败获得信仰值
	RepAttSucAcquired  int32  //修复中攻击成功获得信仰值
	RepAttFaiAcquired  int32  //修复中攻击失败获得信仰值
	DoneAttSucAcquired int32  //完成后攻击成功获得信仰值
	DoneAttFaiAcquired int32  //完成后攻击失败获得信仰值

	FaithLimit	   int32  //建筑信仰上限
	PowerLimit	   int32  //建筑新增的法力值上限
}

//道具表
type Props struct {
	ID        int32     //道具ID
	Name      string    //道具名称
	Type      int32     //道具类型
	Quantity  int32     //数量
	TimeLimit int32     //是否限时
}

//星球基础表
//type StarBase struct {
//	ID           int32     //星球编号
//	Name         string    //星球名称
//	Type         int32     //星球类型
//	Consume      int32     //前置id
//}


type ItemHelpBase struct{
	HelpRequestPrice int32
	HelpAdReward     int32
	StealCost        []int32
	StealRate 		 []float32 //偷取的概率
	StealAllLimit    int32  //所有玩家偷取的上限
	StealSingleLimit int32  //单个玩家偷取的上限
	RequestRelicLimit int32 //求助获取的圣物上限

	BuyRelic		[]int32
	RelicPrice 		[]int32
	RelicQuantity 	[]int32
	PriceSection 	[]float32
	//StealAdReward    int32
	ShopRe            float64
	ShopReHour        int
	ShopReCost        []int32	//圣物商城刷新需要消耗的圣物碎片 [首次,递增,最高]
	ShopItemsNum      int	//商城圣物的数量
	ShopBaseItemsNum  []int	//当前正在激活的圣物组合需要的圣物必在商城中的基础值[三个的组合,五个的组合]
	AdRefresh         int32
	RelicHelpInterval float64
	RelicHelpIntervalLimit float64
	RelicHelpNum      int32
	RelicShareReward  int32

	FirstGroupRelicDrop []int32
}

//游戏基础数据配置表
type GameBase struct{
	StarBelieverLimit int32	  //星球的信徒数量上限
	CostSearchStarMana int32 //搜索星球消耗的魔法值
	CostRecordStarFaith int32 //记录星球消耗的信仰值
	CostAstrola int32   //转动星盘消耗的法力值
	PowerRestoreTime int32//法力恢复时间( 单位 s )
	PowerRestorePoint int32//法力恢复点数
	InitPower int32 //初始化法力值
	InitPowerLimit int32 //初始化法力值上限
	//GetTributeRate float64
	GetRelicsRate []int32
	GetRelicsRateMapping map[int32]int32
	InitFaith int32 //初始化信仰值
	InitDiamond int32 //初始化钻石
	InitLevel int32 //初始化等级
	FriendPointPrices int32 //拍卖物品的价格
	CountdownApplying int64 //过期时间

	CountdownRevenge int
	CountdownSearch float64 //搜索CD时间
	SearchCost int32 //圣物每 SearchCostPerTime 搜索花费
	SearchCostPerTime float64 //没添加一次 SearchCost 的时间 ，单位s

	ShareLimit int32 //每天的分享上限
	ShareReward int32 //分享奖励的信仰值
	//ShareLimitWechat int32 //每天分享微信的次数
	//ShareRewardWechat int32 //每天分享微信的奖励值
	AstrolaTargetCost int32 //刷新目标需要消耗的价格

	TaskRequireTimes   []int32  //转盘转到交互任务的间隔次数
	RevengeMana	int32 //复仇任务消耗的法力值

	DialID	int32 //转盘必随到星盘的ID，如不写则为默认的随机方式

	HelpLimit 	int32
	HelpFaith 	int32
	HelpBonusF		int32
	HelpBeliever	int32
	HelpBonusB		int32

	GetHelpLimit    int32
	HelpReward      int32
	HelpBonusReward int32
	HelperReward    int32
	HelpInterval	int64

	AdBonusReward 	int32 //每天第一次看广告获得的法力值
	AdReward 		int32 //看广告获得的法力值

	PrivilegeReward		int32 //激活特权获得法力值上限
	PrivilegeRewardDis	int32 //第一次激活特权获得法力值

	RobotAvatarPrefix string

	DayGiftFaith	int32
	DayGiftGayPoint int32
	DayGiftPower	int32

	SearchRandomCD float64 //玩家被搜索到的cd时间

	MultipleReward []int32

	MultipleWeight []int32
	MultipleWeightMapping map[int32]int32

	FirstWatchAdPower int32
}

//星盘任务
type AstrolaTask struct{
	ID	int32 //星盘随机ID
	StarType int32 //星球类型
	EventBase int32 //事件ID
	RewardType int32 // 1 信仰相关事件 2 //信徒相关任务
	IsOnly bool //是否唯一
	Weight int32 //随机权重
	//WeightPhase []int32
	//WeightNum []int32
}


//func (this *AstrolaTask) getWeight(condition int32) int32 {
//	index := this.getWeightIndex(condition)
//	if index == -1 || index > len(this.WeightNum) {
//		return 0
//	}
//	return this.WeightNum[index]
//}
//
//func (this *AstrolaTask) getWeightIndex(condition int32) int {
//	var lower int32 = -1
//	for index, upper := range this.WeightPhase {
//		if condition > lower && condition <= upper {
//			return index
//		}
//		lower = upper
//	}
//	return -1
//}

//商店基础配置 wjl 20170619
type ShopBase struct{
	ID int32 //商品的基本ID
	Type int32 //商品类型 // 1 购买钻石 2 月卡 3 购买法力
	Amount int32 //商品获得的数值
	MoneyType int32 //消耗的类型 0x01 信仰值 0x02 钻石 0x03 现金
	Value float64 //消耗需要的数额
}

//func updateBelieverUpgradeBase( data []byte ){
//	var datas []*BelieverUpgradeBase
//	json.Unmarshal(data, &datas)
//	results := make(map[string]*BelieverUpgradeBase)
//	for _, data := range datas {
//		//data.UpgradeResult = strings.Split(data.UpgradeID, "|")
//		results[data.SelectID + data.MatchID] = data
//	}
//	DATA.BelieverUpgrade = results
//}


func updateEventWeightAppend( data []byte ){
	var datas []*EventWeightAppend
	json.Unmarshal(data, &datas)
	results := make(map[int32]*EventWeightAppend)
	for _, data := range datas {
		results[data.ID] = data
	}
	DATA.EVENT_WEIGHT_APPEND = results
}

func updateEventFilterBase( data []byte) {
	var datas []*EventFilterBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]*EventFilterBase)
	mapping := make(map[int32]int32)
	for _, data := range datas {
		mapping[data.ID] = data.Weight
		results[data.ID] = data
	}
	DATA.EVENT_ID_WEIGHT_MAPPING = mapping
	DATA.EVENT_FILTER_DATA = results

	var datas1 []*AstrolaTask
	json.Unmarshal(data, &datas1)
	results1 := make(map[int32][]*AstrolaTask)
	for _, data := range datas1 {
		match := results1[data.RewardType]
		if match == nil {
			match = []*AstrolaTask{}
		}
		match = append(match, data)
		results1[data.RewardType] = match
	}
	DATA.ASTROLA_DATA = results1
}

//func updateAstrolaTaskData( data []byte ){//hjl 更新星盘随机表
//	var datas []*AstrolaTask;
//	json.Unmarshal(data, &datas)
//	results := make(map[int32][]*AstrolaTask)
//	for _, data := range datas {
//		match := results[data.RewardType]
//		if match == nil {
//			match = []*AstrolaTask{}
//		}
//		match = append(match, data)
//		results[data.RewardType] = match
//	}
//	DATA.ASTROLA_DATA = results
//}

func updateRobotSetting(data []byte) {
	var datas []*RobotSetting
	err := json.Unmarshal(data, &datas)
	if err != nil {
		log.Debug("robotSetting unmarshal error: %v",err.Error())
	}
	result := make(map[int32]string)
	for _, data := range datas {
		result[data.ID] = data.Name
	}
	DATA.ROBOT_MAPPPING = result
}

type RobotSetting struct {
	ID 		int32 `json:"robotID"`
	Name	string `json:"robotname"`
}

//更新任务奖励信息
func updateTaskRewardData( data []byte ){
	var datas []*TaskRewardBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]map[int32]*TaskRewardBase)
	for _, data := range datas {
		taskResults := results[data.TaskID]
		if (taskResults == nil) {
			taskResults = make(map[int32]*TaskRewardBase)
			results[data.TaskID] = taskResults
		}
		believerInfo := []*protocol.BelieverInfo{}
		if (data.BelieverLevel != "") {
			levels := strings.Split(data.BelieverLevel, ",")
			nums := strings.Split(data.BelieverNum, ",")
			if (len(levels) == len(nums)) {
				for index, level := range levels {
					id := "b01" + level + character.Int32ToString(1)
					believerInfo = append(believerInfo, &protocol.BelieverInfo{
						Id: id,
						Num:character.StringToInt32(nums[index]),
					})
				}
			}
		}
		data.BelieverInfo = believerInfo
		taskResults[data.EndingID] = data
	}
	DATA.TASK_REWARD_DATA = results
}

func updateEventTaskData(data []byte) {
	var datas []*EventTaskData
	json.Unmarshal(data, &datas)
	results := make(map[int32]*EventTaskData)
	eventTaskMapping := make(map[int32]int32)

	for _, data := range datas {
		results[data.ID] = data
		if data.TriggerStep == 1 {
			eventTaskMapping[data.RefID] = data.ID
		}
	}
	DATA.EVENT_TASK_DATA = results
	DATA.EVENT_TASK_MAPPING = eventTaskMapping
}

func updateItemData(data []byte) {
	var datas []*ItemBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]*ItemRandom)
	rewards := make(map[int32]int32)

	for _, data := range datas {
		rewards[data.ID] = data.CivilizationIncome
		itemRandom := results[data.StarType]
		if itemRandom == nil {
			itemRandom = &ItemRandom{}
			results[data.StarType] = itemRandom
		}
		itemRandom.AddWeight(data.ID, data.PrBuilding, data.PrTask, data.PrBeliever)
	}
	DATA.IITEM_REWARD = rewards
	DATA.ITEM_RANDOM_DATA = results
	//DATA.ITEM_ID = items
}

func updateItemDropRateData(data []byte){
	var datas []*ItemDropRate
	json.Unmarshal(data,&datas)
	DATA.ITEM_DROP_RATE = datas
}

//func updateRecastCostData(data []byte) {
//	var datas []*RecastCostBase
//	json.Unmarshal(data, &datas)
//	results := make(map[int]*RecastCostBase)
//	for _, data := range datas {
//		results[data.Num] = data
//	}
//	DATA.RECAST_DATA = results
//}


func updateItemGroupData(data []byte) {
	var datas []*ItemGroupBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]*ItemGroupBase)
	starGroupMapping := make(map[int32][]*ItemGroupBase)

	for _, data := range datas {
		results[data.ID] = data
		mapping := starGroupMapping[data.StarType]
		if (mapping == nil) {
			mapping = []*ItemGroupBase{}
		}
		mapping = append(mapping, data)
		starGroupMapping[data.StarType] = mapping
	}
	DATA.ITEM_GROUP = results
	DATA.STAR_ITEMGOURP_MAPPING = starGroupMapping
}


func updateItemBaseData( data []byte ){
	var datas  = &ItemHelpBase{}
	json.Unmarshal(data, datas)
	datas.RelicHelpIntervalLimit = 3 * 24 * 60 * 60
	DATA.ItemHelpBase = datas

	//log.Debug("load %v", datas)
}


func updateCardReward( data []byte ){
	var datas []*CardReward
	json.Unmarshal(data, &datas)
	results := make(map[int32]map[int32]util.WeightData)
	for _, data := range datas {
		eventMapping := results[data.EventType]
		if eventMapping == nil {
			eventMapping = make(map[int32]util.WeightData)
			results[data.EventType] = eventMapping
			//if data.RewardType == constant.CARD_REWARD_FAITH {
			//	DATA.DEFAULT_CARD_FAITH_SCOPE = data.Num
			//}
			//results[data.EventType] =
		}
		eventMapping[data.ID] = data
	}
	DATA.EVENT_CARD_MAPPING = results
}

func updateGameBaseData( data []byte ){//wjl20170614 通过zk 获取游戏基础数据
	var datas  = &GameBase{}
	json.Unmarshal(data, datas)
	datas.GetRelicsRateMapping = make(map[int32]int32)
	for index, weight := range datas.GetRelicsRate {
		datas.GetRelicsRateMapping[int32(index)] = weight
	}
	datas.MultipleWeightMapping = make(map[int32]int32)
	for index, weight := range datas.MultipleWeight {
		datas.MultipleWeightMapping[int32(index)] = weight
	}
	DATA.GameBase = datas

	if DATA.CountdownSearch == 0 {
		DATA.CountdownSearch = 1
	}

	if DATA.SearchCostPerTime == 0 {
		DATA.SearchCostPerTime = 60 * 30 //默认半小时 单位s
	}
}

func updateShopBaseData( data []byte ){//wjl 20170619 通过zk 获取游戏商店数据
	var datas []*ShopBase;
	json.Unmarshal(data, &datas)
	results := make(map[int32]*ShopBase)
	for _, data := range datas {
		results[data.ID] = data
	}
	DATA.SHOP_BASEDATA = results;
}

//func updateExpBaseData( data []byte ){//hjl 20170703 通过zk 获取经验数据
//	var datas []*ExpBase
//	json.Unmarshal(data, &datas)
//	results := make(map[int32]*ExpBase)
//	for _, data := range datas {
//		results[data.ID] = data
//	}
//	DATA.EXP_DATA = results;
//}

func updateBuildingsBaseData( data []byte){//wjl 20170619 通过zk 建筑数据
	var datas []*BuildingBase
	json.Unmarshal(data, &datas)
	results := make(map[int32]*BuildingBase)
	for _, data := range datas {
		results[data.ID] = data
	}
	DATA.BUILDING_DATA = results;
}


func updateDialData(data []byte) {
	var datas []*DialData
	json.Unmarshal(data, &datas)
	results := make(map[int32]*DialData)
	for _, data := range datas {
		results[data.ID] = data
	}
	DATA.DIAL_DATA = results
}

func updateDiaLimitData(data []byte) {
	var datas []*DialLimitData
	json.Unmarshal(data, &datas)
	civilWeightMapping := make(map[int32]map[int32]int32)
	resultsMapping := make(map[int32]map[int32]*DialLimitData)
	for _, data := range datas {
		weightMapping := civilWeightMapping[data.CivilLevel]
		if weightMapping == nil {
			weightMapping = make(map[int32]int32)
		}
		weightMapping[data.DialID] = data.Weight
		dataMapping := resultsMapping[data.CivilLevel]
		if dataMapping == nil {
			dataMapping = make(map[int32]*DialLimitData)
		}
		dataMapping[data.DialID] = data

		civilWeightMapping[data.CivilLevel] = weightMapping
		resultsMapping[data.CivilLevel] = dataMapping
	}
	DATA.DIAL_CIVILLEVEL_DIALID_WEIGHT_MAPPING = civilWeightMapping
	DATA.DIAL_LIMIT_DATA_MAPPING = resultsMapping
}

func updateInitDial(data []byte) {
	var datas []*InitDialData
	json.Unmarshal(data, &datas)
	results := make(map[int32]map[int32]int32)
	guideIDMapping := make(map[int32]int32)
	var minRequireLevel int32 = 25 // 25 - 建筑最高等级
	var maxRequireLevel int32 = 0 // 0 - 建筑最低等级
	for _, data := range datas {
		initDialMapping := results[data.BuildingLevel]
		if initDialMapping == nil {
			initDialMapping = make(map[int32]int32)
		}
		initDialMapping[data.ID] = data.Weight
		results[data.BuildingLevel] = initDialMapping
		guideIDMapping[data.ID] = data.DialID
		if data.BuildingLevel > maxRequireLevel {
			maxRequireLevel = data.BuildingLevel
		}
		if data.BuildingLevel < minRequireLevel {
			minRequireLevel = data.BuildingLevel
		}
	}
	DATA.GUIDE_MAX_REQUIRE_LEVEL = maxRequireLevel
	DATA.GUIDE_MIN_REQUIRE_LEVEL = minRequireLevel
	DATA.GUIDE_DIAL_LEVEL_ID_WEIGHT_MAPPING = results
	DATA.GUIDE_ID_DIALID_MAPPING = guideIDMapping
}
//func BaseParse(path string, dataType reflect.Type, name string) {
//	BaseParse1(path, dataType, name, false)
//}

//func BaseParse1(path string, dataType reflect.Type, name string, single bool) {
//	//配置表列字段映射关系
//	mapping := make(map[int]string)
//	//配置表的列类型映射关系
//	typeMapping := make(map[int]string)
//
//	xlFile, err := xlsx.OpenFile(path)
//	if err != nil {
//		log.Error("%v", err)
//		return
//	}
//	//infile := flag.String("input", path, "Path to the xlsx file")
//	//flag.Parse()
//	//xlFile, _ := xlsx.OpenFile(*infile)
//	container := reflect.ValueOf(&Base).Elem().FieldByName(name)
//	//container.Set(reflect.ValueOf())
//	sheet := xlFile.Sheets[0]
//
//	for n, row := range sheet.Rows {
//		if n == 2 {
//			for i, cell := range row.Cells {
//				k, _ := cell.String()
//				typeMapping[i] = k
//			}
//		} else if n == 3 {
//			for i, cell := range row.Cells {
//				k, _ := cell.String()
//				mapping[i] = k
//			}
//		} else if n > 3 {
//			data := reflect.New(dataType).Interface()
//			mutable := reflect.ValueOf(data).Elem()
//			for i, cell := range row.Cells {
//				if (typeMapping[i] == "int") {
//					k, _ := cell.Int()
//					value := mutable.FieldByName(mapping[i])
//					if (value.IsValid()) {
//						value.SetInt(int64(k))
//					}
//
//				} else if (typeMapping[i] == "string") {
//					k, _ := cell.String()
//					value := mutable.FieldByName(mapping[i])
//					if (value.IsValid()) {
//						value.SetString(k)
//					}
//
//				}
//			}
//			if (single) {
//				container.Set(reflect.ValueOf(data))
//			} else {
//				container.SetMapIndex(mutable.FieldByName("ID"), reflect.ValueOf(data))
//			}
//
//		}
//	}
//}

/*----------------事件------------*/
type EventData struct {
	ID               int32 			`json:"id"`	//事件id
	Name             string 		`json:"name"`	//事件名称
	Steps		 []*StepData	 	`json:"steps"`	//事件步骤
}

type StepData struct {
	Module		int32    	`json:"module"`		//事件模块id
	TimeLimit	int32    	`json:"timelimit"`	//时间限制
	Data		interface{}    	`json:"data"`		//模块配置数据,每个模块的配置数据不一样需要根据模块id自动加载
}

func updateEventData(data []byte) {
	var datas []*EventData
	json.Unmarshal(data, &datas)
	results := make(map[int32]*EventData)
	for _, data := range datas {
		results[data.ID] = data
		for _, v:=range data.Steps {
			moduleData, _ := json.Marshal(v.Data)
			v.Data = loadData(ModuleEnum(v.Module), moduleData)
			//log.Debug("%v:  %v",v.Module, v.Data)
		}

	}
	DATA.EVENT_DATA = results
}


