package conf


// [StarBase]	
type StarBase struct {
	Type int32 `json:"Type"`   //星球类型
	LinkName string `json:"linkName"`   //名称
	Name string `json:"Name"`   //显示名称
	ResourcePath string `json:"ResourcePath"`   //星球资源路径
	Desc string `json:"Desc"`   //星球描述
	BelieverID string `json:"BelieverID"`   //生成的信徒id
	RelicPosition string `json:"RelicPosition"`   //圣物祭坛的摆放位置
	BelieverArea string `json:"BelieverArea"`   //信徒生成的区域中心点及半径
	
}

// [Message]	
type Message struct {
	Id int32 `json:"id"`   //编号
	Icon string `json:"icon"`   //图标
	Heading string `json:"heading"`   //标题
	Name string `json:"name"`   //物品名称
	Desc int32 `json:"desc"`   //介绍1
	
}

// [TalkContent]	
type TalkContent struct {
	Id int32 `json:"id"`   //编号
	Content string `json:"content"`   //对话内容
	Direction int32 `json:"direction"`   //对话方向
	Time float32 `json:"time"`   //持续时间(秒)
	
}

// [工作表1]	
type TaskReward struct {
	TaskID int32 `json:"TaskID"`   //事件ID
	EndingID int32 `json:"EndingID"`   //结局ID
	FaithNum int32 `json:"FaithNum"`   //信仰值数量
	GetRelic int32 `json:"GetRelic"`   //是否获得圣物
	BelieverLevel int32 `json:"BelieverLevel"`   //信徒等级
	BelieverNum int32 `json:"BelieverNum"`   //信徒数量
	ItemID int32 `json:"ItemID"`   //道具ID
	ItemNum int32 `json:"ItemNum"`   //道具数量
	BuffID int32 `json:"BuffID"`   //Buff的ID
	
}

// [Buff]	
type Buff struct {
	Id int32 `json:"id"`   //编号
	LinkName string `json:"linkName"`   //调用名称
	BuffType int32 `json:"buffType"`   //buff的具体类型
	BuffNum float32 `json:"buffNum"`   //buff参数
	Display float32 `json:"display"`   //显示的值
	
}

// [BuildingsBase]	
type BuildingsBase struct {
	ID int32 `json:"ID"`   //编号
	Name string `json:"Name"`   //建筑名称
	BuildID int32 `json:"BuildID"`   //建筑类型
	Level int32 `json:"Level"`   //建筑等级
	StarID int32 `json:"StarID"`   //所属星球
	Desc string `json:"Desc"`   //建筑说明
	BuildTime int32 `json:"BuildTime"`   //建造时长(秒)
	RepairTime int32 `json:"RepairTime"`   //修复时长(秒)
	BoomTime int32 `json:"BoomTime"`   //损坏时长(秒)
	UpgradeConsumption int32 `json:"UpgradeConsumption"`   //建造消耗信仰值数量
	UpgradeBelieverLevel int32 `json:"UpgradeBelieverLevel"`   //建造所需信徒等级
	RepairConsumption int32 `json:"RepairConsumption"`   //修理消耗信仰值数量
	RepairBelieverLevel int32 `json:"RepairBelieverLevel"`   //维修所需信徒等级
	PowerLimit int32 `json:"PowerLimit"`   //建成后获得的法力和法力上限
	UpdateFaithTime int32 `json:"UpdateFaithTime"`   //信仰增加的时间间隔(秒)
	UpdateFaithNum int32 `json:"UpdateFaithNum"`   //单位时间增加的信仰值
	FaithLimit int32 `json:"FaithLimit"`   //信仰值上限
	ReceiveFaithPercent int32 `json:"ReceiveFaithPercent"`   //可领取信仰下限 (%)
	RequireCivilizationLevel int32 `json:"RequireCivilizationLevel"`   //建造所需最小文明度等级
	CivBuildIncome int32 `json:"CivBuildIncome"`   //建造完成获得文明度
	CivRepairIncome int32 `json:"CivRepairIncome"`   //维修完成获得文明度
	
}

// [TalkObject]	
type TalkObject struct {
	Id int32 `json:"id"`   //编号
	Type int32 `json:"type"`   //说话物体
	Talker string `json:"talker"`   //说话物体编号
	Comprise int32 `json:"comprise"`   //说话的内容
	CompriseLevelNotEnough int32 `json:"compriseLevelNotEnough"`   //说话的内容
	Mission100talk int32 `json:"mission100talk"`   //任务100中说的话
	Mission110talk int32 `json:"mission110talk"`   //任务110中说的话
	Mission200talk int32 `json:"mission200talk"`   //任务200中说的话
	Mission210talk int32 `json:"mission210talk"`   //任务210中说的话
	
}

// [CivilizationDetail]	
type CivilizationDetail struct {
	Id int32 `json:"id"`   //编号
	StarType int32 `json:"starType"`   //星球类型
	Level int32 `json:"level"`   //星球文明度等级
	Title string `json:"title"`   //文明度等级名称
	Intro string `json:"intro"`   //文明度等级的文字介绍
	
}

// [20170728]	
type GameObject struct {
	Id int32 `json:"id"`   //编号
	ID string `json:"ID"`   //保持唯一有序即可
	Name string `json:"Name"`   //游戏对话中显示的名称
	StarType int32 `json:"StarType"`   //所属星球
	Type int32 `json:"Type"`   //游戏物件类型
	Props string `json:"Props"`   //对应role文件夹
	LV int32 `json:"LV"`   //等级
	Sex int32 `json:"Sex"`   //性别
	Hp int32 `json:"Hp"`   //生命值
	Atk int32 `json:"Atk"`   //攻击力
	Spd int32 `json:"Spd"`   //移动速度
	Desc string `json:"Desc"`   //描述
	
}

// [Civilization]	
type Civilization struct {
	Id int32 `json:"id"`   //编号
	StarOrdinal int32 `json:"starOrdinal"`   //玩家的第几个星球。大于10的和第10个星球相同
	Level int32 `json:"level"`   //文明度等级
	Threshold int32 `json:"threshold"`   //升到本级所需文明度
	Reward int32 `json:"reward"`   //完成该等级进入下一个等级后，可领取该等级的钻石奖励。进入等级8就完成了等级8
	Gift int32 `json:"gift"`   //到达这个等级的奖励
	
}

// [Item]	
type Item struct {
	Id int32 `json:"id"`   //编号
	Name string `json:"name"`   //显示的名称
	Type int32 `json:"type"`   //物件类型
	Icon int32 `json:"icon"`   //对应icon文件
	Display int32 `json:"display"`   //显示排序
	Color int32 `json:"color"`   //颜色
	StarType int32 `json:"starType"`   //所属星球
	Desc string `json:"desc"`   //描述
	PrBeliever float32 `json:"prBeliever"`   //信徒合成获得权重
	PrBuilding float32 `json:"prBuilding"`   //升级维修建筑获得权重
	PrTask float32 `json:"prTask"`   //完成任务获得权重
	CivilizationIncome int32 `json:"civilizationIncome"`   //第一次获得该圣物得到的文明度
	
}

// [钻石买VIP月卡]	
type ShopBase struct {
	ID int32 `json:"ID"`   //编号
	Type int32 `json:"Type"`   //商品类型
	Name string `json:"Name"`   //显示的名称
	Amount int32 `json:"Amount"`   //商品的数量
	MoneyType int32 `json:"MoneyType"`   //货币类型
	Value float32 `json:"Value"`   //货币数量
	ItemLabel int32 `json:"ItemLabel"`   //标签
	Tab int32 `json:"Tab"`   //分页
	Desc string `json:"desc"`   //描述
	
}

// [ItemGroup]	
type ItemGroup struct {
	Id int32 `json:"id"`   //保持唯一有序即可
	Name string `json:"name"`   //游戏中显示的名称
	Display int32 `json:"display"`   //显示排序用的
	StarType int32 `json:"starType"`   //所属星球
	Rarity int32 `json:"rarity"`   //组合内的圣物数量
	Reward int32 `json:"reward"`   //奖励内容
	Content int32 `json:"content"`   //内含圣物
	BuffID int32 `json:"buffID"`   //对应的buff的ID
	CivilizationIncome int32 `json:"civilizationIncome"`   //第一次凑齐圣物组合得到的文明度
	Desc int32 `json:"desc"`   //描述
	SuccessTalk string `json:"successTalk"`   //试组合成功时的发言
	SpeakRole string `json:"speakRole"`   //试组合时发言的角色
	
}

// [ItemGroupUnlock]	
type ItemGroupUnlock struct {
	Id int32 `json:"id"`   //编号
	BuildingRequire int32 `json:"buildingRequire"`   //组合解锁需求lv1建筑数量
	Unlock int32 `json:"unlock"`   //解锁的组合类型
	
}

// [BelieverBuildExpedite.csv]	
type BelieverBuildExpedite struct {
	Id int32 `json:"id"`   //编号
	BuildingLevel int32 `json:"buildingLevel"`   //建筑等级
	BelieverLevel int32 `json:"believerLevel"`   //信徒等级
	DecreaseTime int32 `json:"decreaseTime"`   //减少的时间（单位为秒）
	
}

// [Reward]	
type Reward struct {
	Id int32 `json:"id"`   //奖励ID
	LinkName string `json:"linkName"`   //关联用名称
	FaithNum int32 `json:"faithNum"`   //信仰值数量
	DiamondNum int32 `json:"diamondNum"`   //钻石数量
	RelicPointNum int32 `json:"relicPointNum"`   //圣物碎片数量
	BelieverLevel int32 `json:"believerLevel"`   //信徒等级
	BelieverNum int32 `json:"believerNum"`   //信徒数量
	
}

// [BelieverUpgrade]	
type BelieverUpgrade struct {
	Id int32 `json:"id"`   //编号
	SelectID string `json:"SelectID"`   //选中的ID
	MatchID string `json:"MatchID"`   //配对的ID
	RequireCivilizationLevel int32 `json:"RequireCivilizationLevel"`   //合成所需最小文明度等级
	Cost int32 `json:"Cost"`   //合成费用
	UpgradeID int32 `json:"UpgradeID"`   //合成后的ID
	Weight int32 `json:"Weight"`   //权重
	UpgradeNum int32 `json:"UpgradeNum"`   //合成后数量
	CivilizationIncome int32 `json:"CivilizationIncome"`   //成功合成后获得文明度
	
}




