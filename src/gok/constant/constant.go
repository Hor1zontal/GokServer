package constant

import (
	"math/rand"
)

const (
	ES_LOG                   = true  //是否上传到elastic 平台
	ANALYSIS                 = false //是否开启性能分析
 	DB_DEBUG                 = true	 //是否开启数据库操作
 	DB_SIGAL_TIMEOUT float64 = 1 //数据库操作超时告警阈值 1秒
)

var historyAttackThreshold = []int32{1, 5, 25, 125, 625}

func IsHistoryAttackThreshold(value int32) bool {
	for _, curr := range historyAttackThreshold {
		if curr == value {
			return true
		}
	}
	return false
}

const MAX_GUIDE_STEP int32 = 18 //引导的最大步骤数

type LOGOUT_TYPE int

const (
	LOGOUT_TYPE_NONE         LOGOUT_TYPE = 0 //无
	LOGOUT_TYPE_OTHER        LOGOUT_TYPE = 1 //用户在其他地方登陆
	LOGOUT_TYPE_TIMEOUT      LOGOUT_TYPE = 2 //长时间没有操作
	LOGOUT_TYPE_CLOSE        LOGOUT_TYPE = 3 //服务器关闭
	LOGOUT_TYPE_GATE_CLOSE   LOGOUT_TYPE = 5 //用户状态被改变
	LOGOUT_TYPE_NEWVERSION   LOGOUT_TYPE = 6 //新版本需要重新登录热更
	LOGOUT_TYPE_SERVER_CLOSE LOGOUT_TYPE = 7 //服务器停服
	LOGOUT_TYPE_MAX_SESSION  LOGOUT_TYPE = 8 //服务器会话已满
	LOGOUT_TYPE_SERVER_MAINTAIN LOGOUT_TYPE = 9 //服务器维护中
	LOGOUT_TYPE_CLEAN_ACCOUNT LOGOUT_TYPE = 10 //清除测试账号

	HEARBEAT_SEQ int32 = 3
	//exception_push_seq int32 = 1000
)

const (
	MONEY_TYPE_GOLD int32 = 1
	MONEY_TYPE_DIAMOND int32 = 2
)

const (
	OP_BELIEVER_ADD int32 = 0
	OP_BELIEVER_DEC int32 = 2
)

const (
	INTERNAL_STATISTICS_SERVICE_CALL = 0 //服务调用日志
	INTERNAL_STATISTICS_ONLINE = 1 //在线情况
	STATISTICS_DATA = 2 //统计数据

)

const (
	MAX_RANDOM_COLOR int32 = 5
	MAX_SEARCH_LIMIT int32 = 1
)

const (
	SERVER_MAINTAIN_STATE_CLOSE int32 = 0
	SERVER_MAINTAIN_STATE_OPEN int32 = 1
)

//服务器状态
const (
	SERVER_STATE_OPEN int32 = 0 //开启
	SERVER_STATE_CLOSE_NEW int32 = 1 //停新
	SERVER_STATE_CLOSE int32 = 2 //停服
)

const (
	SHIELD_TYPE_BULDING int32 = 1
	SHIELD_TYPE_BELIEVER int32 = 2
	SHIELD_TYPE_FAITH int32 = 3
)


const (
	STRANGER_TYPE_TASK   = 1 //任务中
	STRANGER_TYPE_LOOT   = 2 //抢夺
	STRANGER_TYPE_FOLLOW = 3 //被关注
	STRANGER_TYPE_BUY_SELL = 4 //被购买物品

	STRANGER_LIMIT = 20
	NEWSFEED_LIMIT = 50

	MUTUALEVENT_LIMIT = 20

	BUILDING_BROKEN_TIME = 86400
	SHIELD_ADD_TIME int32 = 86400 //防护盾恢复时间

	RANDOM_TARGET_COUNT = 1

	MAX_BELIEVER_LEVEL = 6
	MAX_BUILDING_LEVEL = 5

	MAX_TEMP_ITEM = 5


	MALE_STR string = "1" //男
	FEMALE_STR string = "2" //女


	//MUTUAL_TIMEOUT = int(24 * util.SECOEND_OF_HOUR)

)

func RandSex() string {
	result := rand.Float32()
	if result < 0.5 {
		return MALE_STR
	} else {
		return FEMALE_STR
	}
}

const (
	//权重加成条件-信仰
	WEIGHT_ADD_TYPE_FAITH int32 = 0
	//权重加成条件-信徒
	WEIGHT_ADD_TYPE_BELIEVER int32 = 1
)


type MomentsType int

//朋友圈消息类型 1 出售物品 2 物品被购买 3 切换下个星球 4 分享圣物 5 分享圣物组合
const (
	MOMENTS_TYPE_SALE MomentsType = 1
	MOMENTS_TYPE_BUY = 2
	//MOMENTS_TYPE_STAR = 3
	MOMENTS_TYPE_SHARE_ITEM = 4
	MOMENTS_TYPE_SHARE_ITEMGROUP = 5
)

const (
	NEWSFEED_TYPE_BE_REJECT        int32 = 1  //被拒绝索取
	NEWSFEED_TYPE_BE_LOOT_ITEM     int32 = 2  //被抢夺星球圣物
	NEWSFEED_TYPE_BE_ACCEPT        int32 = 3  //被通过索取, 直接下发道具
	NEWSFEED_TYPE_BE_BUY_SALE      int32 = 4  //被购买物品
	NEWSFEED_TYPE_BE_LOOT_FAITH    int32 = 5  //被抢夺信仰	1 信仰	2 星球类型	3 			4
	NEWSFEED_TYPE_BE_ATK_BUILD     int32 = 6  //被攻打建筑	1 信仰	2 星球类型	3 建筑类型	4 建筑等级
	NEWSFEED_TYPE_BE_LOOT_BELIEVER int32 = 7  //被抢夺信徒	1 		2 星球类型	3 			4 被抢的信徒
	NEWSFEED_TYPE_BE_REQUEST_ITEM  int32 = 8  //被索取圣物
	NEWSFEED_TYPE_BE_BUY_GOODS     int32 = 9  //被购买圣物货架物品   1 物品id  2 物品数量 3 消费友情点数量

	NEWSFEED_TYPE_BE_SHIELD int32 = 10 //抵挡攻击 param1 防护罩类型(1建筑 2信徒 3信仰)

	NEWSFEED_TYPE_LOOT_ITEM     int32 = 12  //抢夺星球圣物
	//NEWSFEED_TYPE_BUY_SALE      int32 = 14  //购买物品
	NEWSFEED_TYPE_LOOT_FAITH    int32 = 15  //抢夺信仰
	NEWSFEED_TYPE_ATK_BUILD     int32 = 16  //攻打建筑   //param1 ret.GetFaith() param2 ret.GetItemID() param3 buildingID
	NEWSFEED_TYPE_LOOT_BELIEVER int32 = 17  //抢夺信徒
	NEWSFEED_TYPE_REQUEST_ITEM  int32 = 18  //索取圣物


	NEWSFEED_TYPE_SHARE_ITEM int32 = 20 //分享圣物
	NEWSFEED_TYPE_SHARE_ITEM_GROUP int32 = 21 //分享圣物组合
	NEWSFEED_TYPE_SHARE_BELIEVER int32 = 22 //分享新信徒
	NEWSFEED_TYPE_SHARE_TASK int32 = 23 //分享任务

	NEWSFEED_TYPE_SHARE_SUCC	int32 = 25// 分享成功
	NEWSFEED_TYPE_MUTUAL_FOLLOW int32 = 26 // 互相关注

	NEWSFEED_TYPE_DONE_ITEM_GROUP int32 = 30 //达成圣物组合

	NEWSFEED_TYPE_BUILDING_DESTORY = 31 //建筑被摧毁
	NEWSFEED_TYPE_BUILDING_REPAIR_DONE = 32 //建筑修理完成
	NEWSFEED_TYPE_BUILDING_UPGRADE_DONE = 33 //建筑升级完成

	NEWSFEED_TYPE_NEXT_STAR = 34 //到达下个星球
	NEWSFEED_TYPE_BE_HELP_REPAIR = 35 //被帮助修理建筑

	NEWSFEED_TYPE_HELP_ITEMHELP	int32 = 39 //援助圣物
	NEWSFEED_TYPE_BE_HELP_ITEMHELP int32 = 40 //被援助圣物
	NEWSFEED_TYPE_BE_LOOT_ITEMHELP int32 = 41 //被偷取圣物
	NEWSFEED_TYPE_PUBLIC_ITEMHELP int32 = 42  //请求圣物求助

	NEWSFEED_TYPE_DONE_ITEMHELP int32 = 43  //圣物求助信息已结束 【param1】 itemID [param2] helpNum [param1] lootNum

	NEWSFEED_TYPE_GUIDE_BE_ATK_BUILD int32 = 50


	TASK_ROB_FAITH    int32 = 110
	TASK_ROB_BELIEVER int32 = 210
	TASK_ATK_BUILDING int32 = 310
	TASK_GET_FAITH    int32 = 100
	TASK_GET_BELIEVER int32 = 200

)

func IsGuideRevengeNewsFeed(id int32) bool {
	return true
}

func IsBeAttackNewsFeed(id int32) bool {
	return id == NEWSFEED_TYPE_BE_LOOT_FAITH || id == NEWSFEED_TYPE_BE_ATK_BUILD || id == NEWSFEED_TYPE_BE_LOOT_BELIEVER || id == NEWSFEED_TYPE_GUIDE_BE_ATK_BUILD
}

func IsAttackNewsFeed(id int32) bool {
	return id == NEWSFEED_TYPE_LOOT_FAITH || id == NEWSFEED_TYPE_ATK_BUILD || id == NEWSFEED_TYPE_LOOT_BELIEVER
}

//是否攻击事件
func IsAttackEvent(id int32) bool {
	return id == EVENT_ID_LOOT_BELIEVER || id == EVENT_ID_ATK_BUILDING || id == EVENT_ID_LOOT_FAITH
}


func GetShieldType(eventType int32) int32 {
	if eventType == EVENT_ID_LOOT_BELIEVER {
		return SHIELD_TYPE_BELIEVER
	}
	if eventType == EVENT_ID_ATK_BUILDING {
		return SHIELD_TYPE_BULDING
	}
	if eventType == EVENT_ID_LOOT_FAITH {
		return SHIELD_TYPE_FAITH
	}
	return 0
}



const (
	EVENT_ID_REFRESH_FAITH int32 = 100

	EVENT_ID_REFRESH_BELIEVER int32 = 200
	EVENT_ID_LOOT_BELIEVER int32 = 210
	EVENT_ID_ATK_BUILDING int32 = 310
	EVENT_ID_LOOT_FAITH int32 = 110



	EVENT_REWARD_TYPE_FAITH = 1      //任务奖励类型信仰
	EVENT_REWARD_TYPE_BELIEVER = 2
	EVENT_REWARD_TYPE_POWER = 3


	CARD_REWARD_FAITH = 7
	CARD_REWARD_TARGET = 8 //目标玩家


)

const (
	RANK_REFRESH_HOUR = 5
)


const (
	GUIDE_TASK int32 = 4
	GUIDE_BUILDING_FAITH int32 = 6
)

const (
	MATCH_RULE_BELIEVER = 1
	MATCH_RULE_BUILDING = 2
)

const (
	STATISTIC_TYPE_SALE 			int32 = 1 //统计标识 买卖圣物
	STATISTIC_TYPE_LOOT_FAITH   	int32 = 2 //排行版标识 抢信仰
	STATISTIC_TYPE_ATK_BUILDING     int32 = 3 //排行版标识 攻击建筑
	STATISTIC_TYPE_LOOT_BELIEVER    int32 = 4 //排行版标识 抢信徒
	STATISTIC_TYPE_STAR_ONLINE      int32 = 5

	//------------ 单个星球的统计数据---------------------------
	STAR_STATISTIC_TYPE_EVENT int32 = 10 //随机事件次数
	//STAR_STATISTIC_TYPE_BE_ATTACK int32 = 11 //被攻打次数
	//STAR_STATISTIC_TYPE_ATTACK int32 = 12 //攻打次数
	//STAR_STATISTIC_TYPE_BUILDNUM int32 = 13 //建造次数
	//STAR_STATISTIC_UPGRADE_BELIEVER int32 = 14 //升级信徒
	//STAR_STATISTIC_TYPE_ONLINE int32 = 15 //在线时间

	//交互任务相关统计
	STAR_STATISTIC_LOOT_BELIEVER    int32 = 11 //单星球抢信徒任务次数
	STAR_STATISTIC_LOOT_FAITH       int32 = 12 //单星球抢信仰任务次数
	STAR_STATISTIC_ATK_BUILDING     int32 = 13 //单星球拆建筑次数
	STAR_STATISTIC_BE_LOOT_BELIEVER int32 = 14 //单星球被抢信徒任务次数
	STAR_STATISTIC_BE_LOOT_FAITH    int32 = 15 //单星球被抢信仰任务次数
	STAR_STATISTIC_BE_ATK_BUILDING  int32 = 16 //单星球被拆建筑次数
	//信徒相关统计
	STAR_STATISTIC_ACC_USE_BELIEVER int32 = 21 //单星球加速使用信徒个数
	STAR_STATISTIC_UPGRADE_BELIEVER int32 = 22 //单星球合成信徒总个数
	STAR_STATISTIC_RESOLVE_BELIEVER int32 = 23 //单星球被分解信徒数量
	STAR_STATISTIC_LOOT_BELIEVER_NUM 	int32 = 24 //抢夺信徒个数
	STAR_STATISTIC_BE_LOOT_BELIEVER_NUM int32 = 25 //被抢信徒个数

	//单星球圣物援助次数
	STAR_STATISTIC_HELP_ITEM_HELP int32 = 26 //单星球圣物援助次数
	//单星球合成指定信徒个数
	//单星球抢到指定信徒个数

	//信仰相关统计 [30,40)
	STAR_STATISTIC_EXPEND_FAITH_BELIEVER int32 = 31	//合成信徒消耗的信仰
	STAR_STATISTIC_EXPEND_FAITH_BUILD    int32 = 32 //建造建造消耗的信仰
	STAR_STATISTIC_EXPEND_FAITH_REPAIRE  int32 = 33 //修理建筑消耗的信仰
	STAR_STATISTIC_GAIN_FAITH_EVENT      int32 = 34 //单星球交互任务获得的信仰值数量
	STAR_STATISTIC_GAIN_FAITH_DIAL       int32 = 35 //单星球转盘抽取到的信仰值数量
	STAR_STATISTIC_GAIN_FAITH_BUILDING   int32 = 36 //单星球总建筑收取的信仰值数量
	STAR_STATISTIC_GAIN_FAITH_SHARE		 int32 = 37 //单星球分享收取的信仰值数量

	//其他统计
	STAR_STATISTIC_UNLOCK_STORY			int32 = 41	//单星球解锁的故事个数 ???
	STAR_STATISTIC_ADAPT_STORY			int32 = 42	//单星球改编的故事个数 ???
	STAR_STATISTIC_ONLINE				int32 = 43	//单星球通关时长 ?算不算上下线

	//power [50,60)
	STAR_STATISTIC_EXPEND_POWER_DIAL		int32 = 51 //转盘消耗的法力值
	STAR_STATISTIC_EXPEND_POWER_REVENGE		int32 = 52 //复仇消耗的法力值
	STAR_STATISTIC_GAIN_POWER_AUTO			int32 = 53 //自动恢复获得的法力值
	STAR_STATISTIC_GAIN_POWER_DIAL			int32 = 54 //转盘获得的法力值
	STAR_STATISTIC_GAIN_POWER_AD			int32 = 55 //看广告获得的法力值
)


const (
	ITEMHELP_EVENT_HELP int32 = 0
	ITEMHELP_EVENT_LOOT_FAILED int32 = 1
	ITEMHELP_EVENT_LOOT_SUCCESS int32 = 2
)

const (
	SEARCH_OPT_UPDATE_BUILDING int32 = 1
	SEARCH_OPT_UPDATE_BELIEVER int32 = 2
	SEARCH_OPT_REMOVE_STAR int32 = 3

	SEARCH_OPT_UPDATE_ACTIVE int32 = 5
	SEARCH_OPT_UPDATE_RECEIVE_HELP int32 = 6
	SEARCH_OPT_CHANGE_STAR int32 = 7
)


//1开启事件 2受到攻击 3发起攻击
const (
	STAR_HISTORY_NEW_EVENT int32 = 1
	STAR_HISTORY_BE_ATTACK int32 = 2
	STAR_HISTORY_ATTACK int32 = 3
)

const (
	RANDOM_BUILDING uint8 = 0
	RANDOM_TASK uint8 = 1
	RANDOM_BELIEVER uint8 = 2
)


func GetEventMatchRule(eventID int32) int {
	filterType :=MATCH_RULE_BELIEVER
	if eventID == EVENT_ID_ATK_BUILDING {
		filterType = MATCH_RULE_BUILDING
	}
	return filterType
}

const (
	NOTICE_ON	int32 = 0
	NOTICE_WILL	int32 = 1
	NOTICE_DONE	int32 = 2
)

const (
	ALL_BUILDING_TYPE int32 = -1
)

const (
	PUBLIC_WECHAT_TYPE_HELP 	int32 = 1	//求助
	PUBLIC_WECHAT_TYPE_SHOW 	int32 = 2	//炫耀
)

const (
	WECHAT_HELP_REF_BELIEVER int32 = 1
	WECHAT_HELP_REF_FAITH    int32 = 2
	WECHAT_HELP_REF_POWER    int32 = 3
	WECHAT_HELP_REF_REPAIR   int32 = 4
	WECHAT_HELP_REF_ITEM	 int32 = 5
)

const (
	AD_TYPE_GET_POWER		int32 = 1 // 领法力的广告
	AD_TYPE_GET_GAY_POINT	int32 = 2 // 领圣物碎片的广告
)

const (
	PRIVILEGE_DRAW_POWER int32 = 10
)

const (
	MULTIPLE_DIAL_TYPE_SHARE int32 = 1
	MULTIPLE_DIAL_TYPE_AD int32 = 2
)



func IsDialFlag (key int32) bool {
	switch key {
	case STAR_FLAG_DIAL, STAR_FLAG_LOOT_FAITH, STAR_FLAG_LOOT_BELIEVER, STAR_FLAG_ATK_BUILDING, STAR_FLAG_GAYPOINT:
		return true
	}
	return false
}


const  (
	// 解锁相关 value: 0-lock 1-unlock
	STAR_FLAG_FIRST_UPDATE_POWER int32 = 1 	//第一次消耗法力
	STAR_FLAG_DIAL               int32 = 2 	//转盘
	STAR_FLAG_HISTORY			 int32 = 3 	//星球历史(客户端更)
	STAR_FLAG_EGG                int32 = 4 	//蛋
	STAR_FLAG_GAYPOINT           int32 = 5 	//圣物碎片
	STAR_FLAG_LOOT_FAITH         int32 = 6 	//抢信仰
	STAR_FLAG_LOOT_BELIEVER      int32 = 7 	//抢信徒
	STAR_FLAG_ATK_BUILDING       int32 = 8 	//拆建筑
	STAR_FLAG_REVENGE_MSG        int32 = 9 	//报复消息标识
	STAR_FLAG_REVENGE			 int32 = 10 //复仇
	STAR_FLAG_AD_MULTIPLE		 int32 = 11 //看广告奖励翻倍
	STAR_FLAG_RELIC_ICON		 int32 = 12 //圣物组合
	STAR_FLAG_MUTUAL			 int32 = 13 //交互

	// other
	STAR_FLAG_UNLOCK_WATCHAD   int32 = 30 // 解锁了看广告
	STAR_FLAG_FIRST_WATACH_AD  int32 = 31 // 领取过
	STAR_FLAG_FIRST_GROUP int32 = 32 //解锁了第一个圣物组合 1 -- 解锁 2 -- 完成
)

const (
	FLAG_VALUE_LOCK int32 = 0
	FLAG_VALUE_UNLOCK int32 = 1

	FLAG_VALUE_HAS_FIRST_WATACH_AD int32 = 1

	FLAG_VALUE_GROUP_UNLOCK int32 = 1
	FLAG_VALUE_GROUP_DONE   int32 = 2
)

const (
	QUERY_BY_UID int32 = 1
	QUERY_BY_NICKNAME int32 = 2
	QUERY_BY_USERNAME int32 = 3
)