/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/11/2
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package constant

const (
	//以下均为公众号管理后台设置项
	//appID          = "wxac5324a28ce7dcbb"
	//encodingAESKey = "pT9QLIRnXraLU9EkG9R7u7RDHXeW8rN3dfwrJz1jbDF"
	//token = "HanamichiSakuragi"

	KEYWORD_EVENT = "event"

	EVENT_SUBSCRIBE = "subscribe"
	EVENT_UNSUBSCRIBE = "unsubscribe"
	EVENT_CLICK = "CLICK"
	EVENT_ACTIVE_PRIVILEGE = "ACTIVE_PRIVILEGE"
	EVENT_DAY_GIFT = "DAY_GIFT"
	EVENT_ACTIVE_PUSH = "ACTIVE_PUSH"

	EVENT_PUSH_ONEDAY_EXPIRE = "OneDayMsg"
	EVENT_PUSH_EXPIRE = "MsgInvalid"
	EVENT_RELIC_AIDFULL = "RelicAidFull" //圣物帮助次数已满
	EVENT_RELIC_STEAL = "RelicSteal" //圣物求助被抢夺
	EVENT_RELIC_AID = "RelicAid" //收到圣物求助

	EVENT_BELIEVER_STEAL = "BelieverSteal" //抢信徒
	EVENT_FAITHLOSE = "FaithLose" //抢心眼
	EVENT_BUILD_ATTACK = "BuildAttack"

	EVENT_BUILD_COMPLETION = "BuildCompletion"
	EVENT_BUILD_REPAIR = "BuildRepair"


	EVENT_BUILD_DAMAGE = "BuildDamage" //建筑即将损毁



	RelicSteal
)
