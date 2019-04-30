package constant


const(




	// [工作表1:是否获得圣物]	
	FALSE int32 = 0 //不可
	TRUE int32 = 1 //可以
	

	// [Buff:buff的具体类型]	
	BUFF_MANA int32 = 1 //回法间隔
	BUFF_COST int32 = 2 //合成费打折
	BUFF_RELIC_PROTECT int32 = 3 //圣物保护
	BUFF_UNR_TIME int32 = 4 //升级维修速度
	BUFF_UNR_COST int32 = 5 //升级维修费用
	BUFF_MANA_LIMIT int32 = 6 //法力上限
	BUFF_FAITH_LIMIT int32 = 7 //信仰仓上限
	BUFF_FAITH int32 = 8 //信仰生成量
	BUFF_BELIEVER int32 = 9 //出生间隔
	BUFF_RELIC_GET int32 = 10 //获得圣物几率
	BUFF_BELIEVER_PROTECT int32 = 11 //阻挡抢信徒
	BUFF_FAITH_PROTECT int32 = 12 //阻挡抢信仰
	BUFF_BUILDING_PROTECT int32 = 13 //阻挡拆建筑
	BUFF_BELIEVER_TIME int32 = 14 //信徒加速时长
	BUFF_RELIC_STEAL int32 = 15 //偷取圣物的成功率
	BUFF_BELIEVER_CT int32 = 16 //信徒暴击率
	


	// [TalkObject:说话物体]	
	type_believer int32 = 0 //信徒
	type_building int32 = 1 //建筑
	


	// [20170728:性别]	
	MALE int32 = 1 //男
	FEMALE int32 = 2 //女
	


	// [Item:物件类型]	
	type_relic int32 = 1 //圣物
	
	// [Item:颜色]	
	RED int32 = 1 //红
	YELLOW int32 = 2 //黄
	GREEN int32 = 3 //绿
	BLUE int32 = 4 //蓝
	PUPPLE int32 = 5 //紫
	

	// [钻石买VIP月卡:商品类型]	
	DIAMPOND int32 = 1 //钻石
	MANA int32 = 3 //法力
	BELIEVER int32 = 4 //信徒
	
	// [钻石买VIP月卡:货币类型]	
	M_GOLD int32 = 1 //金币
	M_DIAMOND int32 = 2 //钻石
	M_RMB int32 = 3 //人民币
	
	// [钻石买VIP月卡:标签]	
	L_NONE int32 = 0 //无
	L_HOT int32 = 1 //热销
	L_CHEAP int32 = 2 //超值
	
	// [钻石买VIP月卡:分页]	
	TAB_DIAMOND int32 = 0 //钻石
	TAB_ITEM int32 = 1 //道具
	

	// [ItemGroup:组合内的圣物数量]	
	NUM_THREE int32 = 0 //三
	NUM_FIVE int32 = 1 //五
	

	// [ItemGroupUnlock:解锁的组合类型]	
	UNLOCK_NUM_THREE int32 = 0 //三
	UNLOCK_NUM_FIVE int32 = 1 //五
	

	FAITH int32 = 0
	BELIEVER_L1 int32 = 1
	BELIEVER_L2 int32 = 2
	BELIEVER_L3 int32 = 3
	BELIEVER_L4 int32 = 4
	BELIEVER_L5 int32 = 5
	BELIEVER_L6 int32 = 6
	GIFT_MANA int32 = 7
	QUEST_ROB_FAITH int32 = 8
	QUEST_ROB_BELIEVER int32 = 9
	QUEST_ATT_BUILDING int32 = 10
	QUEST_GET_FAITH int32 = 11
	QUEST_GET_BELIEVER int32 = 12
	DIAL_GAYPOINT int32 = 13

)

