package db

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"gok/service/msg/protocol"
	"encoding/json"
	"gok/module/game/conf"
)

type DBMessage struct {
	ID     bson.ObjectId `bson:"_id"`    //消息id
	UserID int32         `bson:"userid" unique:"false"` //用户id
	Data   []byte        `bson:"data"`   //持久化二进制消息
}

//陌生人
type DBStranger struct {
	ID         int32     `bson:"_id"`        //陌生人id
	Type       int32     `bson:"type"`       //交互类型 1
	Param      int32     `bson:"param"`      //交互信息参数
	ActiveTime time.Time `bson:"activeTime"` //交互时间
}

//角色
type DBRole struct {
	ID     int32 `bson:"_id" gorm:"AUTO_INCREMENT"` //角色id
	UserID int32 `bson:"userid" unique:"true"`      //用户id
	InviteID int32 `bson:"inviteid"` //分享者的id
	//Icon	      int32 	       `bson:"icon"`            		 //角色图标id
	NickName    string    `bson:"nickname" rorm:"nname"` //角色名称
	Desc        string    `bson:"desc" rorm:"desc"`      //角色签名
	RegTime     time.Time `bson:"regtime"`               //角色注册时间
	//LoginTime   time.Time `bson:"logintime"`             //角色当前登录时间
	//LogoutTime  time.Time `bson:"logouttime"`            //用户上次的登出时间
	ReleaseTime time.Time `bson:"releasetime"`           //用户内存的释放时间
	//Level         int32            `bson:"level" rorm:"level"`       //当前等级
	Exp             int32     `bson:"exp"`                //当前等级经验值 神力值
	Power           int32     `bson:"power" rorm:"power"` //法力值
	PowerLimit      int32     `bson:"limit"`              //法力值上限  体力
	PowerTime       time.Time `bson:"powertime"`          //法力值上次的恢复时间戳
	SalePublicTime  time.Time `bson:"publictime"`         //用户上次发布圣物的时间
	WatchAdGetPower time.Time `bson:"powerAdTime"`        //用户下次得到20点法力的时间
	HelpPublicTime  time.Time `bson:"helpPublicTime"`     //用户上次发布求助的时间
	HelpPublicItem  int32     `bson:"helpPublicItem"`     //用户上次发布求助的物品

	GuideDial       []int32    `bson:"guideDial"`           //用户前几次随机的initDialID

	MallItemsReTime time.Time `bson:"itemReTime"`//圣物商店里的圣物刷新的时间
	NextRandomItem int32 	  `bson:"nextRandomItem"`	 //下次一定随机到的圣物

 	Faith          int32     `bson:"faith" rorm:"faith"`     //信仰值   金币
	GayPoint       int32     `bson:"gaypoint"`               //友情点
	Diamond        int32     `bson:"diamond" rorm:"diamond"` //钻石

	//StarsRecordOri  []int32   `bson:"stars_recordOri"`    //用户 原始星球的储存
	//StarsRecordUser []int32   `bson:"stars_recordUser"`   //用户 占有星球收藏
	//StarsComplete 	 []*DBRoleStar_Complete `bson:"stars_complete"`  //用户 已经完成的星球
	StarsSelect []int32 `bson:"stars_select"` //用户可选择的星球
	//StarsComplete []int32 `bson:"stars_complete"` //用户 已经完成的星球
	//StarStatistics *DBStarStatistics                          //当前星球的统计数据
	Statistics []*DBStatistics    `bson:"statistics"` //用户 已经完成的星球
	StarStatistics []*DBStarStatistics   `bson:"starStatistics"` //用户-星球数据统计
	//Task       []*DBRoleEventTask `bson:"task"`       //用户任务
	//Assist         []*DBRoleAssist    `bson:"assist"`         //用户协助请求
	Flags         []*DBRoleFlag     `bson:"flag"`         //用户标识
	BelieverFlags []*DBBelieverFlag `bson:"believerFlag"` //用户标识 信徒是否合成过
	//TempItems      []int32            `bson:"tempItem"`       //临时物品
	Items      []*DBRoleItem      `bson:"item"`      //用户物品
	//ItemGroups []*DBRoleItemGroup `bson:"itemGroup"` //图鉴（已存的组合）
	//ItemGroupsRecords []*DBRoleItemGroupRecord `bson:"itemGroupRecord"` //尝试组合记录
	//Strangers 		 []*DBStranger      `bson:"stranger"`           //交互记录
	NewsFeeds []*DBNewsFeed `bson:"newsFeed"` //动态消息
	Deals     []*DBNewsFeed `bson:"deal"`     //交易信息，索取和获取

	MutualEvents []*DBMutualEvent `bson:"mutualEvent"` //
	ItemRandoms []*DBItemRandom `bson:"itemRandom"` //

	MallItems []*protocol.MallItem `bson:"mallItems"`//圣物商店里的圣物

	Display []int32 `bson:"display"` //表现标识
}

//type DBRoleItemGroupRecord struct {
//	Items []int32 `bson:"items"`
//	Num   int32   `bosn:"number"`
//}

type DBItemRandom struct {
	Weight int32
	ItemId []int32
	//Number int32
}


////用户邮件
//type DBUserMail struct {
//
//}
//
////邮件数据
//type DBMail struct {
//	ID  int64  `bson:"_id" gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id"` //邮件id
//	Power int64
//	Faith int64
//	GayPoint
//	Diamond
//	Believer string[]
//}

//message NewsFeed {
//optional int32 id = 1; //关联用户
//optional int32 type = 2; //消息类型 1索取圣物 2拒绝索取 3抢夺圣物
//optional int64 time = 3; //创建时间
//optional int32 param1 = 4; //动态消息参数1
//optional int32 param2 = 5; //动态消息2
//}

//最近的交互事件
type DBMutualEvent struct {
	Uid  int32     `bson:"uid"`  //仇人uid
	Type int32     `bson:"type"` //交互事件类型
	Time time.Time `bson:"time"` //交互事件时间
}

//动态消息
type DBNewsFeed struct {
	ID     string    `bson:"_id"`
	Uid    int32     `bson:"uid"`
	Type   int32     `bson:"type"`
	Time   time.Time `bson:"time"`
	Param1 int32     `bson:"param1"`
	Param2 int32     `bson:"param2"`
	Param3 int32     `bson:"param3"`
	Ext    []string  `bson:"ext"`
	//IsRevenge 	bool `bson:"isRevenge"`
	DoneRevenge bool `bson:"doneRevenge"`
	Read bool `bson:"read"`

	Self  *protocol.NewsFeedDetail
	Other *protocol.NewsFeedDetail
}

//背包物品
type DBRoleItem struct {
	ID         int32     `bson:"_id"`
	Num        int32     `bson:"num"`
	UpdateTime time.Time `bson:"updatetime"` //激活的时间
}

////图鉴信息
//type DBRoleItemGroup struct {
//	ID int32 `bson:"_id"`
//
//	Done bool `bson:"done"` //是否达成组合
//
//	Records []*DBRoleItemGroupRecord `bson:"record"`
//	//TryTime time.Time `bson:"trytime"`
//	//UpdateTime	time.Time   `bson:"updatetime"`    //完成的时间
//}

type DBRoleStar struct {
	ID   int32 `bson:"_id"`   //星球id
	Type int32 `bson:"value"` //星球类型
}

//type DBRoleStar_Complete struct{//wjl 20170605 用户星球完成的记录
//	ID		    int32		`bson:"id"`		//星球id
//	Type 		int32		`bson:"type"`		//星球类型
//	Building	[]*DBBuilding	`bson:"building"`	//星球建筑物
//}
//
////建筑
//type DBBuilding struct {
//	ID          int32      `bson:"_id"`        //建筑id  配置表中StartID+BuildID
//	Type        int32      `bson:"type"`       //建筑类型
//	Level       int32      `bson:"level"`      //建筑等级
//}

type ROLE_FLAG int32

const (
	FLAG_NONE    ROLE_FLAG = iota
	FLAG_MESSAGE  //1 是否有新消息标识
	FLAG_FRIEND   //2 是否有新好友标识

	FLAG_AD_RE_MALL_ITEM_COUNT = 3//可以看广告刷新商城圣物次数
	FLAG_RE_MALL_ITEM_COUNT = 4 //刷新商城圣物次数
	FLAG_SEARCH_COUNT  = 5 //搜索次数
	FLAG_SHARE         = 6 //分享次数
	FLAG_GUIDE_REVENGE = 7 //是否进行过复仇引导
	
	FLAG_SHARE_WECHAT_SUCC = 9 // 分享成功 获取奖励次数
	FLAG_SHARE_WECHAT_SHOW = 10 // 炫耀型分享微信次数
	FLAG_SHARE_WECHAT_HELP = 11 // 求助型分享成功次数

	FLAG_GUIDE = 12 //引导标识

	FLAG_DRAW_INVITE_GIFT = 15 // 是否领取过分享登录奖励

	//FLAG_WATCH_AD_GET_POWER = 16 // 是否有看广告领取法力奖励

	FLAG_NEW_USER_BY_SHARE = 20 //通过分享进来的新用户

	FLAG_DRAW_PRIVILEGE     = 21 //玩家是否领取过公众号激活特权

	FLAG_FIRST_USE_POWER	= 30 //第一次消耗法力 1-用过
)

//月统计信息
type DBStatistics struct {
	ID         int32     `bson:"_id"`        //统计标识
	Value      float64   `bson:"value"`      //统计数值
	UpdateTime time.Time `bson:"updatetime"` //上次更新统计数据的时间
}

type DBStarStatistics struct {
	Type       int32           `bson:"type"`      //星球类型
	Statistics []*DBStatistics `bson:"statistic"` //统计
}

//信徒标识
type DBBelieverFlag struct {
	ID         string    `bson:"_id"`        //信徒id   信徒id
	Value      bool      `bson:"value"`      //标识值   是否第一次合成
	UpdateTime time.Time `bson:"updatetime"` //上次标识修改的时间
}

type DBRoleFlag struct {
	Flag       ROLE_FLAG `bson:"flag"`       //标识
	Value      int32     `bson:"value"`      //标识值
	UpdateTime time.Time `bson:"updatetime"` //上次标识修改的时间
}

////角色协助请求
//type DBRoleAssist struct {
//	EventID     int32      `bson:"eventid"`    //需要协助的事件id
//	Uid         int32      `bson:"uid"`        //需要协助的用户id
//	NickName    string     `bson:"nickname"`   //需要协助的用户昵称
//	Msg         string     `bson:"msg"`        //需要协助的宣言
//	CreateTime  time.Time  `bson:"createtime"` //发布协助请求的时间
//}
//
//func (this *DBRoleAssist) Name() string {
//	return C_ASSIST
//}

//任务字段
//type TaskField struct {
//	Name      string `bson:"name"`      //任务字段名
//	Value     int32  `bson:"value"`     //任务字段值
//	Threshold int32  `bson:"threshold"` //任务字段阈值
//}

type TASK_STATE int32

const (
	//TASK_STATE_NOT_ACCEPT TASK_STATE = iota
	TASK_STATE_RUNNING TASK_STATE = 1
	TASK_STATE_DONE    TASK_STATE = 2
)

//用户事件任务
type DBRoleEventTask struct {
	ID          int32      `bson:"_id"`        //任务id
	Type        int32      `bson:"type"`       //任务类型
	RefID       int32      `bson:"refid"`      //任务关联id 会关联事件
	State       TASK_STATE `bson:"state"`      //任务状态
	EndingID    int32      `bson:"endingid"`   //结局id
	CreateTime  time.Time  `bson:"createtime"` //任务创建时间
	RevengeID   int32      `bson:"revengeID"`  //当是复仇任务时为复仇对象的uid
	RandomCount int32      `bson:"count"`      //随机目标次数

	RewardFaith    int32                    `bson:"-"`  //任务奖励信仰值
	RewardBeliever []*protocol.BelieverInfo `bson:"-"`  //任务奖励信徒
	RewardItem     int32                    `bson:"-"`  //任务奖励圣物
	RewardGayPoint int32                    `bsoon:"-"` //任务奖励友情点
}
//事件
type DBEvent struct {
	ID          int32      	`bson:"_id" gorm:"AUTO_INCREMENT"`            //事件id
	Type        int32 		`bson:"type"`           //当前的子事件类型
	Guide       bool 		`bson:"guide"`           //是否引导事件
	StepModules []*EventModule  `bson:"stepModules"`    //当前的子事件进行过的步骤模块,顺序按照下标
	History     []*SubEvent     `bson:"history"`        //之前进行过的子事件历史记录
	DisplaySUid int32 			`bson:"displayStarUid"` //进入事件显示的星球所属用户id
	Caller      *DBEventMember	`bson:"caller"`         //事件触发者
	CreateTime  time.Time  		`bson:"createtime"`     //事件创建时间
}

type SubEvent struct {
	Type        int32 		`bson:"type"`           //子事件类型
	StepModules []*EventModule  	`bson:"stepModules"`    //进行过的步骤模块,顺序按照下标
}

type DBEventMember struct {
	Uid  		int32 	`bson:"uid"`        //目标玩家id  0代表游戏外玩家
	Nickname  	string 	`bson:"nickname"`   //目标玩家昵称
}

func (this *DBEvent) IsGuide() bool {
	return this.Guide
}

type EventModule struct {
	Type           conf.ModuleEnum		`bson:"type"`           //模块类型
	StartTimestamp time.Time       		`bson:"starttimestamp"` //当期事件模块的开始时间戳
	EndTimestamp   time.Time       		`bson:"endtimestamp"`   //当前事件的结束时间戳,没有为&time.Time{}
	PersistentData []byte 				`bson:"persistentdata"` //持久化的数据源
	Handler        ModuleHandler        `bson:"-"`              //从持久化数据中加载的数据对象，不存储到数据库
}

func (this *EventModule) HandleMessage(request *protocol.C2GS, response *protocol.GS2C, context EventContext)  {
	if (this.Handler != nil) {
		this.Handler.HandleMessage(request, response, context)
	}
}

//处理事件限时到
func (this *EventModule) HandleTimesUp(context EventContext)  {
	if (this.Handler != nil) {
		this.Handler.HandleTimesUp(context)
	} else {
		context.NextStep()
	}
}

func (this *EventModule) GetJsonData() []byte {
	if (this.Handler == nil) {
		return nil
	}
	data, err := json.Marshal(this.Handler)
	if (err != nil) {
		return nil
	}
	return data
}

//同步模块持久化数据
func (this *EventModule) SyncPersistentData() {
	this.PersistentData = this.GetJsonData()
}


//事件模块处理类
type ModuleHandler interface {
	Init(data []byte, config *conf.StepData)  //从持久化数据源中初始化数据
	Start(context EventContext)  //模块启动
	HandleMessage(request *protocol.C2GS, response *protocol.GS2C, context EventContext)  //处理消息请求
	HandleTimesUp(context EventContext) //处理时间限制到达
}

type ModuleTimberHandler interface {
	DealTimer(time time.Time, context EventContext) //处理定时任务
}

type ModuleRewardHandler interface {
	AppendReward(reward *protocol.Reward)
}

type ModuleDataHandler interface {
	GetData(eventID int32) string//处理模块数据转换
}

type EventContext interface {
	GetID() int32                                          //获取事件ID
	GetType() int32                                        //获取事件类型
	IsGuide() bool                                         //是否引导事件
	GetCaller() *DBEventMember                             //获取事件的发起者
	GetModuleHandler(module conf.ModuleEnum) ModuleHandler //获取模块处理句柄
	GetCurrentStepModule() *EventModule
	NextStep() bool                    //切换到下一步
	AddNextEvent(eventType int32) bool //切换到下一个子事件
	SetDirty()                         //设置有脏数据
	BroadcastAll(message *protocol.GS2C)
	AppendStatisticData(eventType int32, refNum int32)
	//SetDisplaySUid(displaySUid int32)
	//PushModuleInfo() //推送当前模块信息
}