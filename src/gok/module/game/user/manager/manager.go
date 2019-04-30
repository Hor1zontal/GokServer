/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/3/28
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package manager

import (
	"gok/module/game/db"
	"time"
	"gok/module/game/cache"
	"gok/module/game/conf"
	"reflect"
	"gok/module/game/words"
	"gok/service/exception"
	"gok/service/lpc"
	"gopkg.in/mgo.v2/bson"
	"aliens/log"
)


//角色管理容器
type DataManager struct {
	Data	*db.DBRole
	RoleBaseManager
	RoleEventTaskManager
	RoleNewsFeedManager
	//RoleAssistManager
	RoleFlagManager
	RoleStarManager
	RoleStatisticsManager
	RoleItemManager
	//RoleStrangerManager
	RoleTempManager
	RoleMallManager
	RoleDialManager

	dirty bool
}

//初始化
func (this *DataManager) Init() {
	mutable := reflect.ValueOf(this).Elem()
	params := make([]reflect.Value, 1)
	//数据管理类操作副本数据，这样更新的时候能够做比对增量更新
	params[0] = reflect.ValueOf(this.Data)
	for i := 0; i < mutable.NumField(); i++ {
		f := mutable.Field(i)
		initMethod := f.Addr().MethodByName("Init")
		if initMethod.IsValid() {
			initMethod.Call(params)
		}
	}

	//total := this.FilterDealOverdue()
	//if (total > 0) {
	//	this.TakeInGayPoint(total * conf.DATA.FriendPointPrices)
	//}
}

func (this *DataManager) GetPersistData() []byte {
	persistData, _ := bson.Marshal(this.Data)
	return persistData
}

func (this *DataManager) LoadPersistData(persistData []byte, data interface{}) bool {
	if persistData == nil || len(persistData) == 0 {
		return false
	}
	err := bson.Unmarshal(persistData, data)
	if err != nil {
		log.Warn("invalid user persist data %v", err)
		return false
	}
	return true
}

func (this *DataManager) SetDirty() {
	this.dirty = true
}

func (this *DataManager) CleanDirty() {
	this.dirty = false
}

func (this *DataManager) IsDirty() bool{
	return this.dirty
}

//更新到redis和数据库
func (this *DataManager) UpdateAll() {
	this.UpdateLocalCache()
	//this.UpdateLocalCache()
	//更新过期时间
	//cache.UserCache.UpdateUserExpire(this.Data.UserID, this.Data, int(conf.server.RemoteCacheTimeout))
	//cache.UserCache.UpdateUserData(this.Data.UserID, this.Data)
	//数据库更新
	lpc.DBServiceProxy.Update(this.Data, db.DatabaseHandler)
	//db.DatabaseHandler.UpdateOne(this.Data)
}

//更新redis缓存
//func (this *DataManager) UpdateRemoteCache() {
//
//	//cache.UserCache.UpdateUserData(this.Data.UserID, this.Data)
//	//更新过期时间
//	//cache.UserCache.UpdateUserExpire(this.Data.UserID, this.Data, int(conf.Server.RemoteCacheTimeout))
//
//}

//更新本地缓存
func (this *DataManager) UpdateLocalCache() {
	mutable := reflect.ValueOf(this).Elem()
	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(this.Data)
	for i := 0; i < mutable.NumField(); i++ {
		f := mutable.Field(i)
		initMethod := f.Addr().MethodByName("Update")
		if initMethod.IsValid() {
			initMethod.Call(params)
			//子结构的数据mysql需要再执行相应的插入和更新语句
			//if constant.DATABASE_DIALECT == constant.DIALECT_MYSQL {
			//	DealSubDataUpdate(values)
			//}
		}
	}
}

//确保用户含有角色数据
func (this *DataManager) EnsureData(uid int32, persistData []byte) {
	data := &db.DBRole{}

	if !this.LoadPersistData(persistData, data) {
		loadFromDatabase(uid, data)
	}
	//从缓存中加载数据
	//result := cache.UserCache.GetUserData(uid, data)
	////缓存没有，从数据库加载并更新到缓存
	//if !result {
	//	loadFromDatabase(uid, data)
	//	cache.UserCache.UpdateUserData(uid, data)
	//}

	this.Data = data
	//更新过期时间
	//cache.UserCache.UpdateUserExpire(uid, this.Data, int(conf.Server.RemoteCacheTimeout))
	//从数据库数据初始化
	this.Init()
}

func loadFromDatabase(uid int32, data *db.DBRole) {
	//log.Debug("load from database %v", uid)
	err := db.DatabaseHandler.QueryOneCondition(data, "userid", uid)

	nickname := cache.UserCache.GetUserNickname(uid)
	if nickname == "" {
		nickname = words.RandomName()
	}
	//数据库没有数据需要初始化
	if err != nil {
		data.NickName =nickname
		//data.Icon = cache.UserCache.GetUserAttrInt32(uid, basecache.UPROP_ICON)
		data.UserID =uid
		data.RegTime =time.Now()
		//data.LoginTime =time.Now()
		//data.Level =conf.DATA.InitLevel
		data.Exp =0
		data.InviteID = 0
		data.Power = conf.DATA.InitPower
		data.PowerLimit = conf.DATA.InitPowerLimit
		data.PowerTime =time.Now()
		data.Faith =conf.GetGameBase().InitFaith
		data.GayPoint = conf.GetGameBase().FriendPointPrices
		data.Diamond =conf.GetGameBase().InitDiamond
		//分配星球数据
		//rpc.StarServiceProxy.AllocUserNewStar(uid)
		//id, err := db.DatabaseHandler.GenId(data)
		//if err != nil {
		//	exception.GameException(exception.DATABASE_EXCEPTION)
		//}
		//data.ID = id
		//数据库缓存更新队列
		//lpc.DBServiceProxy.Insert(data, db.DatabaseHandler)
		err1 := db.DatabaseHandler.Insert(data)
		if err1 != nil {
			exception.GameException(exception.DATABASE_EXCEPTION)
		}
	}
	cache.UpdateRoleCache(data)
}



