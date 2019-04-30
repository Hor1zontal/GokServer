/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/5/17
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     zhangyishen
 *******************************************************************************/
package core

import (
	"gok/module/community/db"
	"time"
	"aliens/common"
	"gok/module/community/cache"
	basecache "gok/cache"
	"gok/service/msg/protocol"
	"aliens/common/character"
	"encoding/json"
	"gok/constant"
	"gok/service/lpc"
	"aliens/log"
)

var Moments = &momentsManager{}//moments: make(map[string]*protocol.MomentInfo)}


type momentsManager struct {
	//sync.RWMutex
	//moments map[string]*protocol.MomentInfo
}

func (manager *momentsManager) Init() {
	//manager.Lock()
	//defer manager.Unlock()

	var moments []*db.DMoments
	db.DatabaseHandler.QueryAll(&db.DMoments{}, &moments)

	//if cache.CommunityCache.SetNX(basecache.FLAG_LOADMOMENT_INFO, 1) {
	//	for _,moment := range moments {
	//		cache.CommunityCache.SetMomentsInfo(moment.BuildProtocol())
	//	}
	//}

	if cache.CommunityCache.SetNX(basecache.FLAG_LOADMOMENT_INFO, 1) {
		log.Debug("start load moment info data to redis cache...")
		count := 0
		var moments []*db.DMoments
		err := db.DatabaseHandler.QueryAllLimit(&db.DMoments{}, &moments, 10000, func(data interface{}) bool {
			for _,moment := range moments {
				cache.CommunityCache.SetMomentsInfo(moment.BuildProtocol())
			}
			currLen := len(moments)
			count += currLen
			return currLen == 0
		})
		if err != nil {
			log.Debug("load moment info err: %v", err)
		}

		log.Debug("end load moment info data to redis cache count:%v", count)
	}

}

func (manager *momentsManager) Close() {

}


func (manager *momentsManager) AddMoment1(uid int32, momentType constant.MomentsType, refID int32) *protocol.MomentInfo {
	content, _ := json.Marshal(&db.MomentData{Type:int(momentType), RefID:refID})
	result := manager.AddMoment(uid, string(content))
	if momentType == constant.MOMENTS_TYPE_SALE {
		cache.CommunityCache.SetSaleMoments(refID, result.GetId())
	}
	return result
}

func (manager *momentsManager) RemoveSaleMoments(saleID int32) {
	if saleID == 0 {
		return
	}
	//manager.Lock()
	//defer manager.Unlock()
	//保存朋友圈发布记录
	momentsID := cache.CommunityCache.GetSaleMoments(saleID)
	if momentsID != "" {
		manager.RemoveMoment(saleID, momentsID)
		cache.CommunityCache.DeleteSaleMoments(saleID)
	}
}

func (manager *momentsManager) RemoveMoment(uid int32, momentID string) {
	//manager.Lock()
	//defer manager.Unlock()
	//保存朋友圈发布记录
	//delete(manager.moments, momentID)
	cache.CommunityCache.RemoveMomentsInfo(momentID)
	publicID := character.Int32ToString(uid)
	cache.CommunityCache.RemoveUserPublicMomentID(publicID, momentID)
	cache.CommunityCache.RemoveUserReceiveMomentID(publicID, momentID)
	followings := cache.CommunityCache.GetFollowings(publicID)
	for following, _ := range followings {
		cache.CommunityCache.RemoveUserReceiveMomentID(following, momentID)
	}
}

//生成朋友圈信息--------
func (manager *momentsManager) AddMoment(uid int32, data string) *protocol.MomentInfo {
	//manager.Lock()
	//defer manager.Unlock()
	//moment.ID = db.DatabaseHandler.GenId(moment)
	dbMoment := &db.DMoments{
		ID:util.GenUUID(),
		Uid:uid,
		CreateTime:time.Now(),
		Data:data,
	}
	lpc.DBServiceProxy.Insert(dbMoment, db.DatabaseHandler)
	addMomentsCache(dbMoment)

	moment := dbMoment.BuildProtocol()
	cache.CommunityCache.SetMomentsInfo(moment)
	//保存朋友圈发布记录
	//manager.moments[dbMoment.ID] = moment
	return moment
}


func addMomentsCache(moment *db.DMoments) {
	publicUID := character.Int32ToString(moment.Uid)
	cache.CommunityCache.SetUserPublicMomentID(moment.Uid, moment.ID, moment.CreateTime)
	cache.CommunityCache.SetUserReceiveMomentID(moment.Uid, moment.ID, moment.CreateTime)
	//更新关注发布者的所有用户的朋友圈时间线
	followings := cache.CommunityCache.GetFollowings(publicUID)
	for following, _ := range followings {
		cache.CommunityCache.SetUserReceiveMomentID(character.StringToInt32(following), moment.ID, moment.CreateTime)
	}
}

func (manager *momentsManager) DeleteMoment(id string) {
	//manager.Lock()
	//defer manager.Unlock()
	//moment := manager.moments[id]
	if cache.CommunityCache.RemoveMomentsInfo(id) {
		lpc.DBServiceProxy.Delete(&db.DMoments{ID:id}, db.DatabaseHandler)
	}
}

//获取单个用户的朋友圈消息
func (manager *momentsManager) GetMomentsByIds(ids []string) []*protocol.MomentInfo {
	//manager.RLock()
	//defer manager.RUnlock()
	results := []*protocol.MomentInfo{}
	for _, id := range ids {
		moment := cache.CommunityCache.GetMomentsInfo(id)
		//moment := manager.moments[id]
		if moment != nil {
			results = append(results, moment)
		}
	}
	return results
}
