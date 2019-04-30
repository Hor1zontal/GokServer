/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package constant

type SOCIAL_ID uint8

type EsLogIndex = string

const (
	SOCIAL_ID_FAITH SOCIAL_ID = 1    //信仰
	SOCIAL_ID_POWER SOCIAL_ID = 2    //法力
	SOCIAL_ID_DIAMOND SOCIAL_ID = 3   //钻石
	SOCIAL_ID_GAYPOINT SOCIAL_ID = 4  //圣物碎片

	LOG_COMMAND = "log"
	ES_LOG_COMMAND = "es_log"
	ES_DAY_LOG_COMMAND = "es_day_log"



	ES_LOG_STAR EsLogIndex = "star"
	ES_LOG_USER EsLogIndex = "user"
	ES_LOG_SERVICE EsLogIndex = "service" //服务统计信息
	ES_LOG_ONLINE EsLogIndex = "online"   //在线统计信息



	DB_COMMAND_INSERT = "I"
	DB_COMMAND_UPDATE = "U"
	DB_COMMAND_DELETE = "D"
	DB_COMMAND_FUPDATE = "FU"
	DB_COMMAND_CONDITION_UPDATE = "CU"
	DB_COMMAND_CONDITION_DELETE = "CD"
)

type OPT uint8

const (
	OPT_TYPE_TEST                OPT = 0  //测试
	OPT_TYPE_UPGRADE_BUILDING    OPT = 1  //升级建筑
	OPT_TYPE_BUY_SHOP            OPT = 2  //商店购买
	OPT_TYPE_MAIL                OPT = 3  //邮件获取
	OPT_TYPE_REFRESH_POWER       OPT = 4  //刷新法力值
	OPT_TYPE_RANDOM_EVENT        OPT = 5  //随机星盘事件
	OPT_TYPE_SEARCH_STAR         OPT = 6  //搜索星球
	OPT_TYPE_DRAW_BUILDING_FAITH OPT = 7  //领取建筑信仰
	OPT_TYPE_TASK_REWARD         OPT = 8  //星盘任务奖励
	OPT_TYPE_RECORD_STAR         OPT = 9  //记录星球
	OPT_TYPE_REPAIR_BUILDING     OPT = 10 //修理星球建筑
	OPT_TYPE_UPGRADE_BELIEVER    OPT = 11 //升级信徒
	OPT_TYPE_EVENT_LOOT          OPT = 12 //事件抢夺
	OPT_ITEM_GROUP_ITEM_REWARD   OPT = 13 //圣物组合奖励

	OPT_TYPE_TEMP_BAG OPT = 14 //临时背包获取
	OPT_TYPE_GROOVE OPT = 15 //建筑槽操作
	OPT_TYPE_BUILDING_RESET OPT = 16 //建筑被摧毁

	OPT_TYPE_SALE OPT = 17 //交易
	OPT_TYPE_REVENGE OPT = 18 //复仇事件
	OPT_TYPE_SHARE OPT = 19 //分享
	OPT_TYPE_REFRESH_TARGET OPT = 20 //刷新星盘任务目标

	OPT_TYPE_GOODS_CANCEL OPT = 21 //撤销圣物货架物品
	OPT_TYPE_GOODS_PUBLIC OPT = 22 //发布圣物货架物品
	OPT_TYPE_GOODS_BUY OPT = 23 //购买圣物货架物品
	OPT_TYPE_DRAW_CIVIL OPT = 24 //领取文明度奖励

	OPT_TYPE_ACTIVE_GROUP OPT = 25 //激活图鉴组合
	OPT_TYPE_RANDOM_DIAL OPT = 26 //转转盘
	OPT_TYPE_WECHAT_SHARE OPT = 27 //微信分享
	OPT_TYPE_CANCEL_BUILDING OPT = 28 //取消建筑修理/建造状态
	OPT_TYPE_WATCH_AD OPT = 29 //看广告


	OPT_TYPE_DRAW_HELP OPT = 30 //领取圣物求助
	OPT_TYPE_LOOT_HELP OPT = 31 //偷取圣物求助
	OPT_TYPE_HELP_HELP OPT = 32 //援助圣物求助
	OPT_TYPE_PUBLIC_HELP OPT = 33 //援助圣物求助

	OPT_YTPE_REFRESH_MALL OPT = 34 //刷新圣物商城
	OPT_TYPE_BUY_ITEM OPT = 35 // 圣物商城购买

	OPT_TYPE_OPEN_CARD OPT = 36 //随机事件翻卡片

	OPT_TYPE_SEARCH_ITEM OPT = 50  //搜索圣物
	OPT_TYPE_PRIVILEGE_GIFT OPT = 51  //公众号特权礼包
	OPT_TYPE_DAT_GIFT OPT = 52 //公众号每日礼包
)

