package manager

import (
	"aliens/common/util"
	"gok/constant"
	"gok/module/game/cache"
	"gok/module/game/conf"
	"gok/module/game/db"
	gameutil "gok/module/game/util"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"math/rand"
	"time"
)

const (
	SEVEN_DAY = 7 * 24 * time.Hour
	THREE_DAY = 3 * 24 * time.Hour
)

//角色标识
type RoleNewsFeedManager struct {

	uid int32
	newsFeeds []*db.DBNewsFeed
	deals map[string]*protocol.NewsFeed
}

//初始化
func (this *RoleNewsFeedManager) Init(role *db.DBRole) {
	this.uid = role.UserID
	this.newsFeeds = []*db.DBNewsFeed{}
	this.deals = make(map[string]*protocol.NewsFeed)
	this.newsFeeds = role.NewsFeeds

	//for _, newsFeed := range role.NewsFeeds {
	//	this.newsFeeds = append(this.newsFeeds, newsFeed.BuildProtocol(cache.UserCache.GetUserNickname(newsFeed.Uid)))
	//}


	//set set.NewHashSet()
	//for i := len(this.newsFeeds) - 1; i>=0; i-- {
	//	newsFeed := this.newsFeeds
	//	if constant.IsBeAttackNewsFeed(newsFeed.Type) || constant.IsAttackNewsFeed(newsFeed.Type) {
	//
	//	}
	//}


	for _, deal := range role.Deals {
		this.deals[deal.ID] =  deal.BuildProtocol()
	}
}

//更新数据库内存
func (this *RoleNewsFeedManager) Update(role *db.DBRole) {
	//role.NewsFeeds = this.transfer(this.newsFeeds)
	role.NewsFeeds = this.newsFeeds
	role.Deals = this.transferMap(this.deals)
}

//更新数据库内存
//func (this *RoleNewsFeedManager) transfer(datas []*protocol.NewsFeed) []*db.DBNewsFeed {
//	results := []*db.DBNewsFeed{}
//	for _, data := range datas {
//		results = append(results, transferNewsFeed(data))
//	}
//	return results
//}

func (this *RoleNewsFeedManager) transferMap(datas map[string]*protocol.NewsFeed) []*db.DBNewsFeed {
	results := []*db.DBNewsFeed{}
	for _, data := range datas {
		if gameutil.IsOverdue(data) {
			continue
		}

		results = append(results, transferNewsFeed(data))
	}
	return results
}

func transferNewsFeed(data *protocol.NewsFeed) *db.DBNewsFeed {
	return &db.DBNewsFeed{
		ID:     data.GetId(),
		Uid:    data.GetRelateID(),
		Type:   data.GetType(),
		Time:   util.GetTime(data.GetTime()),
		Param1: data.GetParam1(),
		Param2: data.GetParam2(),
		Param3: data.GetParam3(),
		//IsRevenge: data.GetIsRevenge(),
		DoneRevenge: data.GetDoneRevenge(),
		Ext: data.GetExt(),
	}
}

//特定时间内是否有交互任务
func (this *RoleNewsFeedManager) HasMutual(duration float64) bool {
	for _, newsFeed := range this.newsFeeds {
		if constant.IsBeAttackNewsFeed(newsFeed.Type) && time.Now().Sub(newsFeed.Time).Seconds() < duration {
			return true
		}
	}
	return false
}


//随机交互用户
func (this *RoleNewsFeedManager) RandomMutualTarget() int32 {
	threeDayUser := []int32{}
	sevenDayUser := []int32{}
	//newMutualEvents := []*db.DBMutualEvent{} //最近的交互事件
	for _, newsFeed := range this.newsFeeds {
		//是否允许复仇的事件
		if !constant.IsBeAttackNewsFeed(newsFeed.Type) {
			continue
		}
		duration := time.Now().Sub(newsFeed.Time)
		//复仇CD未到,不允许随机
		if duration.Seconds() <= float64(conf.DATA.CountdownRevenge) {
			continue
		}
		if duration <= THREE_DAY {
			threeDayUser = append(threeDayUser, newsFeed.Uid)
			//newMutualEvents = append(newMutualEvents, newsFeed)
		} else if duration <= SEVEN_DAY {
			sevenDayUser = append(sevenDayUser, newsFeed.Uid)
			//newMutualEvents = append(newMutualEvents, newsFeed)
		} else {
			//清楚过期的数据
			break
		}
	}
	if len(threeDayUser) > 0 {
		return threeDayUser[rand.Intn(len(threeDayUser))]
	}
	if len(sevenDayUser) > 0 {
		return sevenDayUser[rand.Intn(len(sevenDayUser))]
	}
	return 0
}

//
func (this *RoleNewsFeedManager) ExistsNewsFeed(uid int32, newsFeedType int32) bool {
	for _, newsFeed := range this.newsFeeds {
		if newsFeed.Uid == uid && newsFeed.Type == newsFeedType {
			return true
		}
	}
	return false
}

func (this *RoleNewsFeedManager) FindItemHelpNewsFeed(id string) *db.DBNewsFeed {
	for _, newsFeed := range this.newsFeeds {
		if newsFeed.Type == constant.NEWSFEED_TYPE_PUBLIC_ITEMHELP && newsFeed.Ext != nil && len(newsFeed.Ext) > 0 {
			if newsFeed.Ext[0] == id {
				return newsFeed
			}
		}
	}
	return nil
}

func (this *RoleNewsFeedManager) DoneItemHelpNewsFeed(id string, helpNum int32, lootNum int32) *db.DBNewsFeed {
	newsFeed := this.GetNewsFeed(id)
	if newsFeed == nil {
		return nil
	}
	//合并
	if newsFeed.Type != constant.NEWSFEED_TYPE_PUBLIC_ITEMHELP {
		return nil
	}

	newsFeed.Type = constant.NEWSFEED_TYPE_DONE_ITEMHELP
	newsFeed.Param2 = helpNum
	newsFeed.Param3 = lootNum

	return newsFeed
}

//完成所有改用户的复仇消息
func (this *RoleNewsFeedManager) DoneRevengeNewsFeed(uid int32) []string {
	result := []string{}
	for _, newsFeed := range this.newsFeeds {
		if constant.IsBeAttackNewsFeed(newsFeed.Type) && newsFeed.Uid == uid {
			newsFeed.DoneRevenge = true
			result = append(result, newsFeed.ID)
		}
	}
	return result
}

func (this *RoleNewsFeedManager) ReadNewsFeed(id string) bool {
	newsFeed := this.GetNewsFeed(id)
	//复仇完成或者别人的复仇交互不允许在复仇
	if newsFeed == nil  {
		exception.GameException(exception.NEWSFEED_NOT_FOUND)
	}
	if newsFeed.Read {
		exception.GameException(exception.NEWSFEED_HAS_READ)
	}
	newsFeed.Read = true
	return true
}

func (this *RoleNewsFeedManager) AddNewsFeed(newsFeed *protocol.NewsFeed) *protocol.NewsFeed {
	if newsFeed.GetType() == constant.NEWSFEED_TYPE_HELP_ITEMHELP {
		//圣物求助直接更新老消息
		result := this.AddTypeNewsFeed(newsFeed.GetRelateID(), newsFeed.GetType(), 1)
		if result != nil {
			return result.BuildProtocol()
		}
	}

	if newsFeed.GetType() == constant.NEWSFEED_TYPE_DONE_ITEMHELP {
		if newsFeed.Ext != nil && len(newsFeed.Ext) > 0 {
			//发布消息还在需要更新发布消息
			publicNewsFeed := this.FindItemHelpNewsFeed(newsFeed.Ext[0])
			if publicNewsFeed != nil {
				publicNewsFeed.Type = newsFeed.Type
				publicNewsFeed.Param1 = newsFeed.Param1
				publicNewsFeed.Param2 = newsFeed.Param2
				publicNewsFeed.Param3 = newsFeed.Param3
				publicNewsFeed.Ext = newsFeed.Ext
				return publicNewsFeed.BuildProtocol()
			}
		}
	}

	if newsFeed.GetType() == constant.NEWSFEED_TYPE_REQUEST_ITEM || newsFeed.GetType() == constant.NEWSFEED_TYPE_BE_REQUEST_ITEM {
		if gameutil.IsOverdue(newsFeed) {
			return nil
		}
		this.deals[newsFeed.GetId()] = newsFeed
	} else {
		dbNewsFeed := transferNewsFeed(newsFeed)
		isAttackNewsFeed := constant.IsAttackNewsFeed(newsFeed.GetType())

		//更新CD 不允许被随机到
		if isAttackNewsFeed {
			cache.UserCache.UpdateMutual(this.uid, newsFeed.GetRelateID(), conf.DATA.CountdownRevenge)
		}

		//任务交互消息每个用户只需要一条
		if constant.IsBeAttackNewsFeed(newsFeed.GetType()) || isAttackNewsFeed {
			replaceNewsFeed := this.FindReplaceNewsFeed(newsFeed.GetRelateID())
			if replaceNewsFeed != nil {
				dbNewsFeed.Self = replaceNewsFeed.Self
				dbNewsFeed.Other = replaceNewsFeed.Other
			}
			dbNewsFeed.UpdateDetail()
		}

		//达到上限 需要删除第一个
		if len(this.newsFeeds) == constant.NEWSFEED_LIMIT {
			this.newsFeeds = this.newsFeeds[1:constant.NEWSFEED_LIMIT]
		}
		this.newsFeeds = append(this.newsFeeds, dbNewsFeed)
	}
	return newsFeed
}

func (this *RoleNewsFeedManager) AddTypeNewsFeed(uid int32, newsFeedType int32, addValue int32) *db.DBNewsFeed {
	for _, newsFeed := range this.newsFeeds {
		if newsFeed.Uid == uid && newsFeed.Type == newsFeedType {
			newsFeed.Param2 += addValue
		}
		return newsFeed
	}
	return nil
}


func (this *RoleNewsFeedManager) FindReplaceNewsFeed(relateID int32) *db.DBNewsFeed {
	for index, newsFeed := range this.newsFeeds {
		if constant.IsBeAttackNewsFeed(newsFeed.Type) || constant.IsAttackNewsFeed(newsFeed.Type) {

		}
		if newsFeed.Uid == relateID {
			this.newsFeeds = append(this.newsFeeds[:index], this.newsFeeds[index+1:]...)
			return newsFeed
		}
	}
	return nil
}

func (this *RoleNewsFeedManager) GetNewsFeed(id string) *db.DBNewsFeed {
	for _, newsFeed := range this.newsFeeds {
		if newsFeed.ID == id {
			return newsFeed
		}
	}
	return nil
}

//func (this *RoleNewsFeedManager) FilterDealOverdue() int32 {
//	var total int32 = 0
//	for id, deal := range this.deals {
//		if gameutil.IsOverdue(deal) {
//			if deal.GetType() == constant.NEWSFEED_TYPE_REQUEST_ITEM {
//				total ++
//			}
//			delete(this.deals, id)
//		}
//	}
//	return total
//}


func (this *RoleNewsFeedManager) GetDeal(id string) *protocol.NewsFeed {
	for _, deal := range this.deals {
		if deal.GetId() == id {
			return deal
		}
	}
	return nil
}

func (this *RoleNewsFeedManager) RemoveDeal(id string)  {
	delete(this.deals, id)
}

//获取所有物品
func (this *RoleNewsFeedManager) GetProtocolNewsFeeds() []*protocol.NewsFeed {
		results := []*protocol.NewsFeed{}
		for _, newsFeed := range this.newsFeeds {
			results = append(results, newsFeed.BuildProtocol())
		}
		return results
}

//获取所有物品
func (this *RoleNewsFeedManager) GetProtocolDeals() []*protocol.NewsFeed {
	results := []*protocol.NewsFeed{}

	for _, deal := range this.deals {
		results = append(results, deal)
	}
	return results
}
