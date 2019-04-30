/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package session

import (
	"gok/service/msg/protocol"
	"gok/service/exception"
	"gok/module/star/conf"
	"time"
	"gok/service/lpc"
	"gok/constant"
	"gok/service/rpc"
	"gok/module/star/db"
	"gok/module/star/cache"
	clustercache "gok/module/cluster/cache"
)


func NewRobotStarManager(uid int32, session *StarSession) *UserStarManager {
	result := NewUserStarManager(uid)
	session.Active = true
	result.UpdateStarSession(session)
	return result
}

func NewUserStarManager(uid int32) *UserStarManager {
	return &UserStarManager{
		id : uid,
		ownStarMapping: make(map[int32]*StarSession), //用户历史星球id和数据映射关系
		starInfos:      []*protocol.StarInfo{},       //用户当前拥有的星球概要信息
		lastActiveTime: time.Now(),
	}
}

type UserStarManager struct {
	id 	int32
	active         *StarSession           //当前正在玩的星球
	ownStarMapping map[int32]*StarSession //用户历史星球id和数据映射关系
	starInfos      []*protocol.StarInfo   //用户当前拥有的星球概要信息
	lastActiveTime time.Time
}

func (this *UserStarManager) UpdateActiveTime()  {
	this.lastActiveTime = time.Now()
}

func(this *UserStarManager) UpdateStarSession(starSession *StarSession) {
	//兼容老数据
	if starSession.Seq == 0 {
		starSession.Seq = 1
	}
	starInfo := &protocol.StarInfo{StarID: starSession.ID, StarType: starSession.Type}
	if this.starInfos == nil {
		this.starInfos = []*protocol.StarInfo{}
	}
	this.starInfos = append(this.starInfos, starInfo)
	this.ownStarMapping[starSession.ID] = starSession
	if starSession.Active {
		starSession.Init()
		this.active = starSession
		cache.StarCache.SetUserActiveStar(starSession.Owner, starSession.ID)
		cache.StarCache.SetUserActiveStarType(starSession.Owner, starSession.Type)
	}
}

func(this *UserStarManager) GetActiveStar() *StarSession {
	return this.active
}

func (this *UserStarManager) GetOwnStars() map[int32]*StarSession{
	return this.ownStarMapping
}

func(this *UserStarManager) GetStar(starID int32) *StarSession {
	return this.ownStarMapping[starID]
}


func(this *UserStarManager) GetStarInfo() []*protocol.StarInfo {
	return this.starInfos
}

func(this *UserStarManager) HaveActiveStar() bool {
	return this.active != nil
}

func(this *UserStarManager) GetActiveStarID() int32 {
	if this.active == nil {
		return 0
	}
	return this.active.ID
}

func (this *UserStarManager) DealUnSave() {
	if this.active != nil {
		this.active.DealUnSave()
	}
}

func (this *UserStarManager) DealStarTimer() {
	now := time.Now()
	if now.Sub(this.lastActiveTime).Seconds() > conf.Server.FreeTimeout {
		StarManager.RemoveStar(this.id)
	}
	if this.active != nil {
		this.active.DealStarTimer(now)
	}
}

//切换下一个星球
func (this *UserStarManager) NextStar(starType int32) (*StarSession, int32 , int32) {
	srcStar := this.active
	if srcStar == nil{
		exception.GameException(exception.STAR_NOTFOUND)
	}

	srcStar.IsDoneStar()

	faith := srcStar.DrawAllBuildingFaith()
	srcStar.AddStatisticsValue(srcStar.Owner, float64(faith), 0)

	//删除老星球数据
	//delete( this.ownActiveStarMapping, srcStar.ID)
	srcStar.Active = false
	srcStar.DoneTime = time.Now()
	cache.StarCache.SetBuffRelicSteal(srcStar.Owner, 0)
	cache.StarCache.SetBuffMANAInterval(srcStar.Owner, 0)

	lpc.DBServiceProxy.Update(srcStar.DBStar, db.DatabaseHandler)
	rpc.SearchServiceProxy.UpdateData(constant.SEARCH_OPT_REMOVE_STAR, srcStar.ID, 0)
	rpc.SearchServiceProxy.SyncSearchData()

	//分配新星球
	uid := srcStar.Owner
	clustercache.Cluster.DelStarSummary(srcStar.Type)
	//nextStarType := conf.GetNextStarType(srcStar.Type)
	nextStar := StarManager.allocUserStar(uid, starType, false, srcStar.Seq + 1)
	return nextStar, srcStar.Type, faith
}

