package manager

import (
	"gok/constant"
	"gok/module/game/db"
	"gok/module/game/rank"
	"gok/module/statistics/model"
	"gok/service/lpc"
	"gok/service/msg/protocol"
	"time"
)

//角色标识
type RoleStatisticsManager struct {
	uid int32
	Statistics map[int32]*db.DBStatistics

	StarStatistic map[int32]map[int32]*db.DBStatistics     // 星球类型 - 统计类型ID

	targets            []*protocol.Target
	roleStatisticEvent *model.StatisticEvent
}

//初始化
func (this *RoleStatisticsManager) Init(role *db.DBRole) {
	this.uid = role.UserID
	this.Statistics = make(map[int32]*db.DBStatistics)
	for _, statistic := range role.Statistics {
		this.Statistics[statistic.ID] = statistic
	}
	this.StarStatistic = make(map[int32]map[int32]*db.DBStatistics)
	for _, statistics := range role.StarStatistics {
		this.StarStatistic[statistics.Type] = make(map[int32]*db.DBStatistics)
		for _, statistic := range statistics.Statistics {
			this.StarStatistic[statistics.Type][statistic.ID] = statistic
		}
	}

	//this.EnsureStatistics(constant.STATISTIC_TYPE_SALE)
	//this.EnsureStatistics(constant.STATISTIC_TYPE_LOOT_FAITH)
	//this.EnsureStatistics(constant.STATISTIC_TYPE_ATK_BUILDING)
	//this.EnsureStatistics(constant.STATISTIC_TYPE_LOOT_BELIEVER)
}

func (this *RoleStatisticsManager) Update(role *db.DBRole) {
	role.Statistics = this.GetStatisticsArray()
	role.StarStatistics = this.GetStarStatisticsArray()
}

func (this *RoleStatisticsManager) BeginEventStatistic(eventType int32, dial int32) {
	this.roleStatisticEvent = &model.StatisticEvent{
		Uid:    this.uid,
		Event:  eventType,
		Source: dial,
	}
}

func (this *RoleStatisticsManager) AppendEventStatistic(eventType int32, refNum int32) {
	if this.roleStatisticEvent != nil {
		switch eventType {
		case constant.EVENT_ID_LOOT_FAITH:
			this.roleStatisticEvent.FaithLoot = refNum
		case constant.EVENT_ID_LOOT_BELIEVER:
			this.roleStatisticEvent.BelieverLoot = refNum
		case constant.EVENT_ID_ATK_BUILDING:
			this.roleStatisticEvent.BuildAtk = refNum
		}
	}
}

func (this *RoleStatisticsManager) AppendEventCardStatistic(cardType int32) {
	if this.roleStatisticEvent != nil {
		this.roleStatisticEvent.CardType = cardType
	}
}

func (this *RoleStatisticsManager) EndEventStatistic(faith int32) {
	if this.roleStatisticEvent != nil {
		this.roleStatisticEvent.FaithGet = faith
		lpc.StatisticsHandler.AddStatisticData(this.roleStatisticEvent)
	}
}

func (this *RoleStatisticsManager) EventStatisticOpenCard(target *protocol.Target) {
	if this.roleStatisticEvent != nil {
		this.roleStatisticEvent.TargetID = target.GetId()
		//this.roleStatisticEvent.TargetSeq = int32(index + 1)
		this.roleStatisticEvent.Revenge = target.GetMutual()
		this.roleStatisticEvent.BelieverTotal = target.GetBelieverTotalLevel()
		this.roleStatisticEvent.BuildingTotal = target.GetBuildingTotalLevel()
	}

}

func (this *RoleStatisticsManager) EventStatisticRandomTarget(targets []*protocol.Target) {
	this.targets = targets
}

func (this *RoleStatisticsManager) AppendEventRevengeStatistic(targetID int32) {
	if this.roleStatisticEvent == nil {
		return
	}

	for _, target := range this.targets {
		if target.GetId() == targetID {
			this.roleStatisticEvent.TargetID = targetID
			//this.roleStatisticEvent.TargetSeq = int32(index + 1)
			this.roleStatisticEvent.Revenge = target.GetMutual()
			this.roleStatisticEvent.BelieverTotal = target.GetBelieverTotalLevel()
			this.roleStatisticEvent.BuildingTotal = target.GetBuildingTotalLevel()
			//this.roleStatisticEvent = nil
			return
		}
	}
}


//获取所有标识
func (this *RoleStatisticsManager) GetStatisticsArray() []*db.DBStatistics {
	result := []*db.DBStatistics{}
	for _, statistic := range this.Statistics {
		result = append(result, statistic)
	}
	return result
}

//获取所有标识
func (this *RoleStatisticsManager) GetStarStatisticsArray() []*db.DBStarStatistics {
	var starResult []*db.DBStarStatistics
	for starType, statistics := range this.StarStatistic {
		var result []*db.DBStatistics
		for _, statistic := range statistics {
			result = append(result, statistic)
		}
		starResult = append(starResult, &db.DBStarStatistics{Type:starType, Statistics:result})
	}
	return starResult
}


//获取角色标识
func (this *RoleStatisticsManager) GetStatisticsValue(key int32) float64 {
	Statistics := this.Statistics[key]
	if Statistics == nil {
		return 0
	}
	return Statistics.Value
}

func (this *RoleStatisticsManager) updateStatistic(statistic *db.DBStatistics, addValue int32, starType int32, coverValue bool) {
	//初始化的时候更新一波排名
	if statistic.UpdateTime.Before(rank.RankRefreshTime) {
		statistic.Value = 0
	}
	if coverValue {
		statistic.Value = float64(addValue)
	} else {
		statistic.Value = statistic.Value + float64(addValue)
	}
	statistic.UpdateTime = time.Now()
	var manager *rank.RankManager
	if starType == 0 {
		manager = rank.GetRankManager(statistic.ID)
	} else {
		manager = rank.GetStarRankManager(statistic.ID, starType)
	}
	if manager != nil {
		manager.UpdateRank(this.uid, int64(statistic.Value))
	}
}

//func (this *RoleStatisticsManager) EnsureStatistics(id int32) *db.DBStatistics {
//	statistic := this.Statistics[id]
//	if statistic == nil {
//		statistic = &db.DBStatistics{
//			ID:         id,
//			Value:      0,
//			UpdateTime: time.Now(),
//		}
//		this.Statistics[id] = statistic
//		//过了刷新时间需要清理统计数据
//		this.updateStatistic(statistic, 0)
//	}
//
//
//	return statistic
//}

//新增统计数据
func (this *RoleStatisticsManager) AddStatisticsValue(id int32, value int32) *db.DBStatistics {
	statistic := this.Statistics[id]
	if statistic == nil {
		statistic = &db.DBStatistics{
			ID:         id,
			Value:      0,
			UpdateTime: time.Now(),
		}
		this.Statistics[id] = statistic
	}

	//过了刷新时间需要清理统计数据
	this.updateStatistic(statistic, value, 0, false)
	return statistic
}

func (this *RoleStatisticsManager) AddStarStatisticValue(id int32, starType int32, value int32) *db.DBStatistics {
	starStatistic := this.StarStatistic[starType]
	if starStatistic == nil {
		this.StarStatistic[starType] = make(map[int32]*db.DBStatistics)
	}
	statistic := this.StarStatistic[starType][id]
	if statistic == nil {
		statistic = &db.DBStatistics{
			ID:			id,
			Value:      0,
			UpdateTime: time.Now(),
		}
		this.StarStatistic[starType][id] = statistic
	}
	//过了刷新时间需要清理统计数据
	this.updateStatistic(statistic, value, starType, true)
	return statistic
}
