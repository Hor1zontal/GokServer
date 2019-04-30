package user

import (
	"aliens/common/character"
	"aliens/common/util"
	"aliens/log"
	"gok/constant"
	"gok/module/game/cache"
	"gok/module/game/conf"
	"gok/module/game/rank"
	gameutil "gok/module/game/util"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"gok/service/rpc"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

//查找角色
//func (this *Session) SearchRole(nickName string) *db.DBRole {
//	role := &db.DBRole{}
//	err := db.DatabaseHandler.QueryOneCondition(role,"nickname",nickName)
//	//err := db.DatabaseHandler.QueryAllCondition(&db.DBRole{}, "nickname", nick, &redPacketLogs)
//	if (err != nil) {
//		exception.GameException(exception.USER_NOTFOUND)
//	}
//	return role
//}

//func (this *Session) SearchRoles(condtion string) []*db.DBRole {
//	roles := []*db.DBRole{}
//	count := 0
//	if strings.Count(condtion,"") == 1 {   //条件为空，推荐些玩家
//		allRoles := []*db.DBRole{}
//		db.DatabaseHandler.QueryAll(&db.DBRole{}, &allRoles)
//		for _, role := range allRoles  {
//			if cache.UserCache.IsUserOnline(role.UserID) && role.UserID != this.GetID() && count < 10{
//				if (!this.ContainFriend(role.UserID) && !this.ContainAddFriendReq(role.UserID)) {
//					roles = append(roles,role)
//					count++
//				}
//			}
//		}
//	} else {
//		role := &db.DBRole{}
//		err := db.DatabaseHandler.QueryOneCondition(role,"nickname",condtion)
//		if err == nil {
//			roles = append(roles,role)
//		}
//		id,err := strconv.Atoi(condtion)
//		if err == nil {
//			err = db.DatabaseHandler.QueryOneCondition(role,"userid",id)
//			if err == nil {
//				roles = append(roles,role)
//			}
//		}
//	}
//	return roles
//}

//func (this *Session) FindRoleByID(userid int32) *db.DBRole {
//	role := &db.DBRole{}
//	err := db.DatabaseHandler.QueryOneCondition(role, "userid", userid)
//	if (err != nil) {
//		exception.GameException(exception.USER_NOTFOUND)
//	}
//	return role
//}


func (user *Session) DealUpgradeBuilding1(itemID int32, powerLimit int32, powerReward int32, buildingType int32, buildingLevel int32, online bool) int32 {
	return user.DealUpgradeBuilding(itemID, powerLimit, powerReward, []*protocol.BuildingState{&protocol.BuildingState{
		StarType:      user.GetStarType(),
		BuildingType:  buildingType,
		BuildingLevel: buildingLevel,
	}}, online)
}

func (user *Session) DealUpgradeBuilding(itemID int32, powerLimit int32, powerReward int32, buildings []*protocol.BuildingState, online bool) int32 {
	addLimit := powerLimit - user.GetPowerLimitBase()
	if addLimit > 0 {
		user.TakeInPower(powerReward, false, constant.OPT_TYPE_UPGRADE_BUILDING)
		user.SetPowerLimit(powerLimit)
	}

	if buildings == nil {
		return 0
	}

	if itemID == 0 {
		// 随机奖励
		itemID = user.RandomItem(constant.RANDOM_BUILDING, online, false, constant.OPT_TYPE_UPGRADE_BUILDING)
	} else {
		// 升级建筑指定随机得到的圣物（star模块传来的值）
		user.TakeInItem(itemID, 1, constant.OPT_TYPE_UPGRADE_BUILDING, 0)
	}

	//建筑完成不发newsFeed
	//for _, upgradeBuilding := range buildings {
	//
	//
	//	newsFeed := BuildNewsFeed(user.GetID(), constant.NEWSFEED_TYPE_BUILDING_UPGRADE_DONE, upgradeBuilding.GetStarType(), upgradeBuilding.GetBuildingType(), 0)
	//	user.AddNewsFeed(newsFeed)
	//	if online {
	//		//user.WriteMsg(user.BuildRoleCivilizationPush())
	//		user.WriteMsg(BuildNewsFeedPush(newsFeed))
	//	} else {
	//		user.AddOfflineMessage(newsFeed)
	//	}
	//}
	if online {
		user.WriteMsg(user.BuildRoleSocialPush())
	}
	return itemID
}


//func (this *Session) BuildRoleCivilizationPush() *protocol.GS2C {
//	return &protocol.GS2C{
//		Sequence:[]int32{1006},
//		CivilizationPush: &protocol.CivilizationPush{
//			CivilizationLv:this.GetCivilizationLevel(),
//			CivilizationProgress:this.GetCivilizationValue(),
//		}}
//}

func (user *Session) DealRepairedBuilding1(buildingType int32, buildingLevel int32, online bool) int32 {
	return user.DealRepairedBuilding([]*protocol.BuildingState{
		&protocol.BuildingState{
			StarType:      user.GetStarType(),
			BuildingType:  buildingType,
			BuildingLevel: buildingLevel,
		},
	}, online)
}

func (user *Session) DealRepairedBuilding(buildings []*protocol.BuildingState, online bool) int32 {
	var itemID int32 = 0
	//自己修理完成不发newsFeed
	//for _, repairedBuilding := range buildings {
	//	itemID = user.RandomItem(constant.RANDOM_BUILDING, online, online, constant.OPT_TYPE_REPAIR_BUILDING)
	//	newsfeed := BuildNewsFeed(user.GetID(), constant.NEWSFEED_TYPE_BUILDING_REPAIR_DONE,
	//		user.GetStarType(), repairedBuilding.GetBuildingType(), repairedBuilding.GetBuildingLevel())
	//	user.AddNewsFeed(newsfeed)
	//
	//	if online {
	//		//user.WriteMsg(user.BuildRoleCivilizationPush())
	//		user.WriteMsg(BuildNewsFeedPush(newsfeed))
	//	} else {
	//		user.AddOfflineMessage(newsfeed)
	//	}
	//}
	return itemID
}

func (user *Session) DealNewItem(itemID int32) {
	reward := conf.DATA.IITEM_REWARD[itemID]
	if reward > 0 {
		civilizationInfo := rpc.StarServiceProxy.AddCivilization(user.GetID(), reward)
		if user.IsOnline() {
			user.WriteMsg(&protocol.GS2C{Sequence: []int32{1006}, CivilizationPush: civilizationInfo})
		}
	}
}

func (session *Session) DealAutoHelp(currTime time.Time) {
	if session.GetHelpPublicItem() > 0 {
		duration := currTime.Sub(session.GetHelpPublicTime()).Seconds()
		//求助后1分钟后发送给其他玩家求助
		if duration >= conf.DATA.ItemHelpBase.RelicHelpInterval && duration < conf.DATA.ItemHelpBase.RelicHelpIntervalLimit  {
			targets := rpc.SearchServiceProxy.RandomHelpTargets(session.GetID(), session.GetStarType(), conf.DATA.ItemHelpBase.RelicHelpNum).GetTargets()
			if targets != nil && len(targets) > 0 {
				filters := rpc.CommunityServiceProxy.GetFollowingList(session.GetID()).GetFollowings()
				sendMessage := &protocol.C2GS{
					Sequence:[]int32{523},
					AddNewsFeed: BuildNewsFeed(session.GetID(), constant.NEWSFEED_TYPE_PUBLIC_ITEMHELP, session.GetHelpPublicItem(), 0, 0),
				}
				for _, target := range targets {
					if !containsFollowing(target, filters) {
						rpc.UserServiceProxy.PersistCall(target, sendMessage)
					}
				}
			}
			session.CleanHelpPublicTime()
		}
	}
}

func containsFollowing(uid int32, filters []*protocol.UserInfo) bool {
	for _, filter := range filters {
		if filter.GetId() == uid {
			return true
		}
	}
	return false
}


func (user *Session) RandomItem(randomType uint8, pushItem bool, pushSocial bool, opt constant.OPT) int32 {
	itemID, gayPoint := user.RandomItems(randomType)
	if itemID != 0 {
		item := user.TakeInItem(itemID, 1, opt, 0)
		if pushItem {
			user.WriteMsg(gameutil.BuildItemPush(item))
		}
	}
	if gayPoint != 0 {
		user.TakeInGayPoint(1, opt, 0)
		if pushSocial {
			user.WriteMsg(user.BuildRoleSocialPush())
		}
	}
	return itemID
}

func (user *Session) RandomItems(randomType uint8) (int32, int32) {
	if !user.IsDoneFirstGroup() {
		//未完成第一个圣物组合不掉圣物和圣物碎片
		return 0, 0
	}
	//获取随机圣物
	index := util.RandomWeight(conf.DATA.GetRelicsRateMapping)
	if index == 1 {
		//获取友情点
		return 0, 1
	} else if index == 2 {
		//获取随机圣物
		var retItem int32 = 0
		if user.GetNextRandomItem() != 0 {
			retItem = user.GetNextRandomItem()
			user.CleanNextRandomItem()
		} else {
			//当在尝试第一个圣物组合时，随机得到的圣物第一个不是当前组合的圣物，下一个一定是当前组合的圣物（且是自己拥有的组合中数量最少的一个）
			retItem = conf.RandStarItem(user.GetStarType(), user.GetItemWeight(), randomType)
			if user.IsActiveFirstGroup() {
				currentItems := rpc.StarServiceProxy.GetCurrentGroupItems(user.GetID()).GetItemIDs()
				if !conf.Contains(currentItems, retItem) {
					user.SetNextRandomItem(user.GetMinNumItemFromItems(currentItems))
				}
			}
		}
		return retItem, 0
	}
	return 0, 0
}

func (user *Session) IsActiveFirstGroup() bool {
	if user.GetStarFlagValue(constant.STAR_FLAG_FIRST_GROUP) == constant.FLAG_VALUE_GROUP_UNLOCK {
		return true
	}
	return false
}

func (user *Session) IsDoneFirstGroup() bool {
	if user.GetStarFlagValue(constant.STAR_FLAG_FIRST_GROUP) == constant.FLAG_VALUE_GROUP_DONE {
		return true
	}
	return false
}

func (this *Session) BuildRoleSocialPush() *protocol.GS2C {
	return &protocol.GS2C{
		Sequence: []int32{1004},
		RoleInfoPush: &protocol.RoleInfoPush{
			//Level:    this.GetLevel(),
			//Exp:      this.GetExp(),
			Power:      this.GetPower(),
			PowerLimit: this.GetPowerLimit(),
			Faith:      this.GetFaith(),
			Diamond:    this.GetDiamond(),
			GayPoint:   this.GetGayPoint(),
			Subscribe:  this.GetSubscribe(),
		}}
}

func BuildNewsFeed(uid int32, newsFeedType int32, param1 int32, param2 int32, param3 int32) *protocol.NewsFeed {
	return BuildNewsFeed1(uid, newsFeedType, param1, param2, param3, nil)
}

func BuildNewsFeed1(uid int32, newsFeedType int32, param1 int32, param2 int32, param3 int32, ext []string) *protocol.NewsFeed {
	return &protocol.NewsFeed{
		Id:       bson.NewObjectId().Hex(),
		RelateID: uid,
		//RelateNickname:nickname),
		//RelateAvatar:avatar),
		Type:   newsFeedType,
		Time:   time.Now().Unix(),
		Param1: param1,
		Param2: param2,
		Param3: param3,
		Ext:    ext,
	}
}

func BuildNewsFeedPush(newsFeed *protocol.NewsFeed) *protocol.GS2C {
	return &protocol.GS2C{
		Sequence:     []int32{1052},
		NewsFeedPush: newsFeed,
	}
}

func (this *Session) HelpItemAsk(fromShare bool, request *protocol.C2GS) *protocol.HelpItemRet {
	message := request.GetHelpItem()
	if !fromShare && !this.CanTakeoutItem(message.GetItemID(), 1) {
		exception.GameException(exception.ITEM_NOT_ENOUGH)
	}
	message.Limit = conf.DATA.RequestRelicLimit

	followEach := rpc.CommunityServiceProxy.IsEachFollow(this.GetID(), message.GetUid())
	message.EachFollow = followEach

	var result *protocol.HelpItemRet
	if fromShare {
		result = rpc.TradeServiceProxy.NoErrorCall(this.GetID(), request).GetHelpItemRet()
		if result == nil {
			return nil
		}
	} else {
		result = rpc.TradeServiceProxy.Call(this.GetID(), request).GetHelpItemRet()
	}

	//圣物换碎片
	if !fromShare {
		this.TakeInGayPoint(1, constant.OPT_TYPE_HELP_HELP, 0)
		this.TakeOutItem(message.GetItemID(), 1, constant.OPT_TYPE_HELP_HELP, 0)
		result.GayPoint = 1
	}
	itemHelp := result.GetItemHelp()

	//帮助达到上限需要推送公众号
	if itemHelp.GetHelpNum() == message.Limit {
		rpc.PassportServiceProxy.WechatEventPush(itemHelp.GetUid(), constant.EVENT_RELIC_AIDFULL, 0)
	} else {
		rpc.PassportServiceProxy.WechatEventPush(itemHelp.GetUid(), constant.EVENT_RELIC_AID, 0)
	}

	//推送被抢夺求组圣物

	rpc.UserServiceProxy.PersistCall(message.GetUid(), &protocol.C2GS{
		Sequence:[]int32{523},
		AddNewsFeed: BuildNewsFeed(this.GetID(), constant.NEWSFEED_TYPE_BE_HELP_ITEMHELP, result.GetItemHelp().GetItemID(), result.GetItemHelp().GetItemNum(), 0),
	})

	if !followEach {
		rpc.CommunityServiceProxy.FollowEach(message.GetUid(), this.GetID())
		result.EachFollow = true
	}

	//推送玩家数据变更
	pushMessage := &protocol.GS2C{
		Sequence:[]int32{1081},
		ItemHelpPush: result.GetItemHelp(),
	}

	rpc.UserServiceProxy.Push(message.GetUid(), pushMessage)
	rpc.StarServiceProxy.UpdateStarStatistics(this.GetID(), constant.STAR_STATISTIC_HELP_ITEM_HELP, 1, 0)

	return result
}

func (this *Session) DrawDayGift() bool{
	if !cache.UserCache.ExistUserDayGiftCache(this.GetID()) {
		return false
	}
	draw, _ := strconv.ParseBool(cache.UserCache.GetUserDayGiftCache(this.GetID()))
	if !draw {
		this.TakeInFaith(conf.DATA.DayGiftFaith, constant.OPT_TYPE_DAT_GIFT, 0)
		this.TakeInGayPoint(conf.DATA.DayGiftGayPoint, constant.OPT_TYPE_DAT_GIFT, 0 )
		this.TakeInPower(conf.DATA.DayGiftPower,false, constant.OPT_TYPE_DAT_GIFT)
		cache.UserCache.SetUserDayGiftCache(this.GetID(), true)
		log.Debug("draw day gift uid:%v, faith:%v, gayPoint:%v, power:%v", this.GetID(),conf.DATA.DayGiftFaith,conf.DATA.DayGiftGayPoint,conf.DATA.DayGiftPower)
		this.WriteMsg(&protocol.GS2C{
			Sequence:[]int32{1090},
			DayGiftPush:&protocol.Reward{
				Faith:conf.DATA.DayGiftFaith,
				GayPoint:conf.DATA.DayGiftGayPoint,
				Power:conf.DATA.DayGiftPower,
			},
		})
		return true
	}
	return false
}

func (this *Session) TakeInDialReward(dial *conf.DialLimitData, multipleRatio int32) *protocol.Reward{
	reward := &protocol.Reward{}
	dialNum := dial.Num*multipleRatio
	switch dial.Type {
	case constant.FAITH:
		this.TakeInFaith(dialNum, constant.OPT_TYPE_RANDOM_DIAL, 0)
		rpc.StarServiceProxy.UpdateStarStatistics(this.GetID(), constant.STAR_STATISTIC_GAIN_FAITH_DIAL, float64(dialNum), 0)
		reward.Faith = dialNum
	case constant.BELIEVER_L1, constant.BELIEVER_L2, constant.BELIEVER_L3, constant.BELIEVER_L4, constant.BELIEVER_L5, constant.BELIEVER_L6:
		addBeliever := this.DialAddBeliever(dial.Type, dialNum)
		rpcResp := rpc.StarServiceProxy.UpdateBelieverInfo(this.GetID(), addBeliever , rpc.OP_BELIEVER_ADD, true).GetUpdateBelieverRet()
		addBeliever = rpcResp.GetBeliever()
		addFaith := rpcResp.GetFaith()
		if addFaith > 0 {
			this.TakeInFaith(addFaith, constant.OPT_TYPE_RANDOM_DIAL, 0)
			rpc.StarServiceProxy.UpdateStarStatistics(this.GetID(), constant.STAR_STATISTIC_GAIN_FAITH_EVENT, float64(addFaith), 0)
			reward.Faith = addFaith
		}
		reward.Believer = addBeliever
	case constant.GIFT_MANA:
		this.TakeInPower(dialNum, false, constant.OPT_TYPE_RANDOM_DIAL)
		rpc.StarServiceProxy.UpdateStarStatistics(this.GetID(), constant.STAR_STATISTIC_GAIN_POWER_DIAL, float64(dialNum), 0)
		reward.Power = dialNum
	case constant.DIAL_GAYPOINT:
		this.TakeInGayPoint(dialNum, constant.OPT_TYPE_RANDOM_DIAL, 0)
		reward.GayPoint = dialNum
	}
	return reward
}

func (this *Session) GetUserStarRank(rankType int32, value int32) (int32, int32) {
	this.AddStarStatisticValue(rankType, this.GetStarType(), value)
	rankManager := rank.GetStarRankManager(rankType, this.GetStarType())
	myRank := rankManager.GetUserRankData(this.GetID(), rankManager.GetCurrentRank(int64(this.GetID())), rankManager.GetCurrentScore(int64(this.GetID())))
	totalNum := rankManager.GetRankTotalNum(cache.UserCache.GetStarOnlineKey() + character.Int32ToString(this.GetStarType()))
	log.Info("%v %v", myRank.RankNum, totalNum)
	return totalNum - myRank.RankNum + 1, totalNum
}