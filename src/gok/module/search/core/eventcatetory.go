/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/11/6
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

import (
	"aliens/common/util"
	"gok/module/search/conf"
	"time"
)

type StarState struct {
	level int32
	activeTime 		time.Time //上次活跃时间
	lastSelectTime  time.Time //上一次被选中的时间
}

func newEventCategory() *EventCategory {
	return &EventCategory{
		levelMapping : make(map[int32]*StarCategory),
		filterMapping: make(map[int32]*StarState),
		active: make(map[int32]time.Time),
	}
}

type EventCategory struct {

	levelMapping map[int32]*StarCategory  // 等阶 - 星球类别

	filterMapping map[int32]*StarState //处于CD 间的用户

	active map[int32]time.Time
}

func (this *EventCategory) addRandomStar(starID int32) {
	//过滤列表中存在不需要加入到活跃随机玩家中
	_, ok := this.filterMapping[starID]
	if ok {
		return
	}
	this.active[starID] = time.Now()
}

func (this *EventCategory) removeStar(starID int32) {
	for _, starCategory := range this.levelMapping {
		delete(starCategory.mapping, starID)
	}
	delete(this.filterMapping, starID)
}

func (this *EventCategory) addLevelMap(levelBase *conf.RandomTargetBase) {
	this.levelMapping[levelBase.Level] =
		&StarCategory{level:levelBase.Level, min:levelBase.Min, max:levelBase.Max, mapping:make(map[int32]time.Time)}
}

func (this *EventCategory) RandomStar(filter []int32, always bool) int32 {
	//优先从活跃玩家中随机一个
	starID, activeTime := RandomStar(this.active, filter)
	var category *StarCategory = nil
	var level int32 = 0

	if starID == 0 {
		weightMapping := make(map[int32]int32, len(this.levelMapping))
		for level, category := range this.levelMapping {
			weightMapping[level] = int32(len(category.mapping))
		}
		level := util.RandomWeight(weightMapping)
		category = this.levelMapping[level]
		if category == nil {
			return 0
		}
		starID, activeTime = RandomStar(category.mapping, filter)
		level = category.level
	}

	//TODO 同步到其他服务器

	//随机后的数据需要添加到CD数据中,always 为false的不需要放入cd中
	if !always && starID > 0 && activeTime != nil {
		delete(this.active, starID)
		if category != nil {
			//随机列表中清除
			delete(category.mapping, starID)
		}

		//加入到过滤列表
		this.filterMapping[starID] = &StarState{
			level: level,
			activeTime: *activeTime,
			lastSelectTime: time.Now(),
		}
	}
	return starID
}


//清除非活跃玩家
func (this *EventCategory) CleanExpire(now time.Time) {
	for _, starCategory := range this.levelMapping {
		for starID, updateTime := range starCategory.mapping {
			if now.Sub(updateTime).Seconds() > conf.Base.SearchActiveTime {
				delete(starCategory.mapping, starID)
			}
		}
	}
}

//处理过滤的过期数据，放回到索引数据中
func (this *EventCategory) DealFilterExpire(now time.Time) {
	for starID, filter := range this.filterMapping {
		if now.Sub(filter.lastSelectTime).Seconds() > conf.Base.SearchRandomCD {
			delete(this.filterMapping, starID)
			if filter.level <= 0 {
				continue
			}
			category := this.levelMapping[filter.level]
			if category != nil {
				category.mapping[starID] = filter.activeTime
			}
		}
	}
}




func (this *EventCategory) UpdateData(starID int32, value int32) bool {
	if starID <= 0 {
		return false
	}
	filter := this.filterMapping[starID]

	for _, starCategory := range this.levelMapping {
		if starCategory.InScope(value) {
			//过滤列表中存在，更新到过滤列表即可
			if filter != nil {
				filter.level = starCategory.level
				filter.activeTime = time.Now()
			} else {
				starCategory.mapping[starID] = time.Now()
			}
		} else {
			delete(starCategory.mapping, starID)
		}
	}
	return true
}
