/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2017/5/5
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package session

import (
	"aliens/common/util"
	"aliens/log"
	basecache "gok/cache"
	"gok/constant"
	clustercache "gok/module/cluster/cache"
	"gok/module/cluster/center"
	"gok/module/star/cache"
	"gok/module/star/conf"
	"gok/module/star/db"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"gok/service/rpc"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"time"
)

type StarFilter struct {
}

//const (
//	FILTER_IGNORE   int32 = 0  //忽略过滤条件
//	FILTER_MATCH    int32 = 1  //拥有过滤条件
//	FILTER_UNMATCH  int32 = 2  //不拥有过滤条件的属性
//)

var StarManager = &starManager{

	//ownStarIDMapping:       make(map[int32][]*protocol.StarInfo), //用户拥有的星球序号
	//ownActiveStarIDMapping: make(map[int32]*protocol.StarInfo),   //用户当前使用星球信息
	//
	//ownStarMapping:       make(map[int32]*StarSession), //用户已完成的星球
	//ownActiveStarMapping: make(map[int32]*StarSession), //wjl 20170531 隶属于玩家的星球数据列表[ 注：只有未完成的星球才会存在于这里 ]( key 星球id, value: 星球数据 )
	//originalStarMapping:  make(map[int32]*StarSession), //wjl 20170531 原始星球的 map容器( key 星球id, value: 星球数据 )
	userStars: 				 make(map[int32]*UserStarManager),
	guideRobotMapping:       make(map[int32]*UserStarManager),
	robotMapping:            make(map[int32]*UserStarManager),
	eventRobotRandomMapping: make(map[int32]map[int32]map[int32]int32),
}

type starManager struct {
	//sync.RWMutex //事务锁
	robotStarID int32

	//ownStarIDMapping       map[int32][]*protocol.StarInfo //用户拥有的星球序号
	//ownStarMapping       map[int32]*StarSession //用户未使用的星球id和数据映射关系
	//ownActiveStarMapping map[int32]*StarSession //用户当前使用的星球id和数据映射关系

	userStars				map[int32]*UserStarManager     //用户id-用户星球信息映射关系
	//originalStarMapping  map[int32]*StarSession //未被探索的原始星球和id的映射表

	guideRobotMapping       map[int32]*UserStarManager             //引导机器人星球
	robotMapping            map[int32]*UserStarManager             //
	eventRobotRandomMapping map[int32]map[int32]map[int32]int32 //随机事件ID - 随机级别 - 机器人星球id - 随机权重
}



type RobotLevelRandom struct {
	level         int32           //等级
	weightMapping map[int32]int32 //星球id - 权重
}

func (this *starManager) genRobot() {
	//生成引导机器人
	for _, starType := range conf.Base.StarWeCanArrive {
		this.genGuideRobotStar(starType)
	}

	//生成事件机器人
	for _, robotRule := range conf.Base.RobotRule {
		for _, starType := range conf.Base.StarWeCanArrive {
			this.genEventRobotStar(robotRule, starType)
		}
	}
}

func (this *starManager) updateActiveStarInfo(starSession *StarSession, updateRemote bool) *UserStarManager {
	uid := starSession.Owner
	if uid <= 0 {
		return nil
	}
	result:= this.userStars[uid]
	if result == nil {
		result = NewUserStarManager(uid)
		this.userStars[uid] = result
	}
	result.UpdateStarSession(starSession)
	if updateRemote {
		this.updateRemoteCache(starSession, updateRemote)
	}
	cache.StarCache.SetUserIDNode(uid, center.GetServerNode())
	return result
}

func (this *starManager) updateRemoteCache(starSession *StarSession, update bool) {
	cache.StarCache.SetOwner(starSession.ID, starSession.Owner)
	cache.StarCache.SetType(starSession.ID, starSession.Type)

	starSession.updateBelieverCount(update)
	starSession.updateBuildingLevel(update)

	starSession.updateEventTimesCache()

	rpc.SearchServiceProxy.UpdateHelpData(starSession.Owner, constant.SEARCH_OPT_CHANGE_STAR, starSession.Type)
}

//生成事件机器人
func (this *starManager) genEventRobotStar(base *conf.RobotRuleBase, starType int32) {
	starSession := this.allocStar(starType, this.genRobotId())
	starSession.Owner = starSession.ID
	starSession.addCustomBeliever(2, true, 1)
	starSession.addCustomBeliever(2, false, 1)

	if len(base.BuildingType) == len(base.BuildingLevel) {
		for index, buildingType := range base.BuildingType {
			building := starSession.getBuilding(buildingType)

			level := base.BuildingLevel[index]
			if building == nil {
				log.Error("unexpect building type %v", buildingType)
				continue
			}
			building.Exist = level != 0
			building.Level = level

			buildingConf := conf.GetBuildingConf(starSession.Type, building.Type, building.Level)
			if buildingConf != nil {
				building.Faith = buildingConf.FaithLimit
			}

		}
	} else {
		log.Error("robot building config is not valid!")
	}

	if len(base.BelieverLevel) == len(base.BelieverNum) {
		for index, believerLevel := range base.BelieverLevel {
			believerNum := base.BelieverNum[index]
			if believerNum > 0 {
				maleNum := believerNum / 2
				femaleNum := believerNum - maleNum
				starSession.addCustomBeliever(believerLevel, true, maleNum)
				starSession.addCustomBeliever(believerLevel, false, femaleNum)
			}
		}

	} else {
		log.Error("robot believer config is not valid!")
	}

	for _, randomType := range base.RandomType {
		randomMapping := this.eventRobotRandomMapping[randomType]
		if randomMapping == nil {
			randomMapping = make(map[int32]map[int32]int32)
			this.eventRobotRandomMapping[randomType] = randomMapping
		}

		weightMapping := randomMapping[base.Level]
		if weightMapping == nil {
			weightMapping = make(map[int32]int32)
			randomMapping[base.Level] = weightMapping
		}

		weightMapping[starSession.ID] = base.Weight
	}

	this.robotMapping[starSession.ID] = NewRobotStarManager(starSession.ID, starSession)
}

func (this *starManager) genRobotId() int32 {
	this.robotStarID--
	//log.Debug("robot id %v", this.robotStarID)
	return this.robotStarID
}

//生成引导机器人
func (this *starManager) genGuideRobotStar(starType int32) {
	starSession := this.allocStar(starType, this.genRobotId())
	starSession.Owner = starSession.ID
	starSession.addCustomBeliever(2, true, 1)
	starSession.addCustomBeliever(2, false, 1)
	for _, building := range starSession.Building {
		building.Exist = true
		building.Level = rand.Int31n(5) + 1
	}
	this.guideRobotMapping[starSession.ID] = NewRobotStarManager(starSession.ID, starSession)
}

//分配一个原始星球
func (this *starManager) allocStar(starType int32, id int32) *StarSession {
	//从配置表中初始化星球建筑信息
	buildingTypes := conf.GetStarBuildingTypes(starType)
	buildings := make([]*db.DBBuilding, len(buildingTypes))
	for index, buildingType := range buildingTypes {
		buildings[index] = &db.DBBuilding{
			ID:   conf.GenStarBuildingID(starType, buildingType.(int32)),
			Type: buildingType.(int32),
		}
	}
	dbStar := &db.DBStar{
		Type:       starType,
		CreateTime: time.Now(),
		Building:   buildings,
	}

	dbStar.ID = id
	//if id != 0 {
	//	dbStar.ID = id
	//} else {
	//	//id, err := db.DatabaseHandler.GenId(dbStar)
	//	//if err != nil {
	//	//	exception.GameException(exception.DATABASE_EXCEPTION)
	//	//}
	//	//dbStar.ID = id
	//}
	session := newStarSession(dbStar)
	//session.initBuffListener()
	return session
}

//随机获取一个原始星球 wjl 20170531
//func (this *starManager) randOriginalStar() *db.DBStar {
//	l := len(this.originalStarMapping);
//	if l == 0 { //如果长度为0.... 呵呵呵 好吧 那就临时创建几个星球出来
//		for i := 0; i < 12; i++ {
//			randomStarType := conf.RandomStarType()
//			session := this.allocStar(randomStarType, 0)
//			//session := this.allocStar(  int32( i%8+1 ) )
//			if session == nil {
//				continue
//			}
//			session.Owner = 0                                //默认没有拥有任何用户
//			//lpc.DBServiceProxy.Insert(session.DBStar, db.DatabaseHandler)
//			err := db.DatabaseHandler.Insert(session.DBStar)        //更新到数据库中
//			if err != nil {
//				exception.GameException(exception.DATABASE_EXCEPTION)
//			}
//			this.originalStarMapping[ session.ID ] = session //放入原始星球列表内
//		}
//		l = len(this.originalStarMapping) //重新在获取一次长度
//	}
//	rd := rand.Int() % l //来来来 随机一个值吧
//	i := 0
//	for _, v := range this.originalStarMapping {
//		if i != rd {
//			i++
//			continue
//		}
//		return v.DBStar
//	}
//	return nil
//}

//通过用户id获取用户当前使用的星球
func (this *starManager) GetUserActiveStar(uid int32) *StarSession {
	if uid < 0 {
		return this.getRobotStar(uid)
	}
	userStars := this.userStars[uid]
	if userStars == nil {
		return nil
	}
	return userStars.active
}

//通过用户id和星球id获取星球信息
//func (this *starManager) GetStarByID(uid int32, starID int32) *StarSession {
//	if uid < 0 {
//		return this.getRobotStar(uid)
//	}
//	//if v, ok := this.originalStarMapping[ starID ]; ok { //先查原始星球
//	//	return v
//	//}
//	userStars := this.userStars[uid]
//	if userStars == nil {
//		return nil
//	}
//	if starID == 0 {
//		return userStars.active
//	}
//	return userStars.ownStarMapping[starID]
//}

func (this *starManager) GetStarInfoByUser(uid int32) []*protocol.StarInfo {
	userStars := this.userStars[uid]
	if userStars == nil {
		return nil
	}
	return userStars.starInfos
}

func (this *starManager) GetUserStarsByUser(uid int32) *UserStarManager {
	return this.userStars[uid]
}

func (this *starManager) getRobot(id int32) *UserStarManager {
	result := this.robotMapping[id]
	if result != nil {
		return result
	}
	return this.guideRobotMapping[id]
}

func (this *starManager) getRobotStar(id int32) *StarSession {
	robot := this.getRobot(id)
	if robot == nil {
		return nil
	}
	return robot.active
}

func (this *starManager) getUserActiveStarID(uid int32) int32 {
	starInfo := this.userStars[uid]
	if starInfo == nil {
		return 0
	}
	return starInfo.GetActiveStarID()
}

//根据过滤条件随机星球
//func (this *starManager) randomOwnStar(uid int32, num int) []*db.DBStar {
//	starCount := len(this.userStars)
//	//没有任何数据了哦 好吧 返回空吧
//	if starCount == 0 {
//		return nil
//	}
//	results := []*db.DBStar{}
//	randomScope := []int32{}
//	for id, stars := range this.userStars {
//		if stars.active == nil {
//			continue
//		}
//		//不允许把自己的星球给出去
//		if uid == stars.active.Owner {
//			continue
//		}
//		randomScope = append(randomScope, id)
//	}
//	randomLength := int32(len(randomScope))
//	if randomLength == 0 {
//		return results
//	}
//	resultIndex := util.RandIntervalN(0, randomLength-1, num)
//	for _, index := range resultIndex {
//		if index >= randomLength {
//			continue
//		}
//		uid := randomScope[index]
//		results = append(results, this.GetUserActiveStar(uid).DBStar)
//	}
//	return results
//}

//分配一个用户星球
func (this *starManager) allocUserStar(uid int32, starType int32, disable bool, seq int32) *StarSession {
	starSession := this.allocStar(starType, 0)

	starSession.Seq = seq
	starSession.Disable = disable
	starSession.addCustomBeliever(1, true, 12)
	starSession.addCustomBeliever(1, false, 12)
	starSession.Owner = uid
	starSession.OwnTime = time.Now()
	//starSession.State = db.STAR_STATE_DEV  //用户初始化星球未开发中
	starSession.BelieverUpdateTime = time.Now()
	starSession.Active = true

	err := db.DatabaseHandler.Insert(starSession.DBStar)
	if err != nil {
		exception.GameException(exception.DATABASE_EXCEPTION)
	}

	//starSession.addCustomBeliever(2, true, 3)
	//starSession.addCustomBeliever(2, false, 3)
	//starSession.addCustomBeliever(3, true, 3)
	//starSession.addCustomBeliever(3, false, 3)
	//starSession.addCustomBeliever(4, true, 3)
	//starSession.addCustomBeliever(4, false, 3)
	//starSession.addCustomBeliever(5, true, 3)
	//starSession.addCustomBeliever(5, false, 3)
	//for _, building := range starSession.Building {
	//	building.Exist = true
	//	building.Level = rand.Int31n(5) + 1
	//}
	this.updateActiveStarInfo(starSession, true)

	clustercache.Cluster.AddStarSummary(starType)

	return starSession
}

func (this *starManager) UpdateAllStarFlag(key int32, value int32) bool {
	var stars []*db.DBStar
	//db.DatabaseHandler.UpdateOne()
	var count = 0
	err := db.DatabaseHandler.QueryAllLimit(&db.DBStar{}, &stars, 10000, func(data interface{}) bool {
		for _, star := range stars {
			if star.Owner != 0{
				var have bool
				for _, flag := range star.Flags {
					if flag.Flag == key {
						flag.Value = value
						flag.UpdateTime = time.Now()
						have = true
					}
				}
				if !have {
					star.Flags = append(star.Flags, &db.DBStarFlag{Flag:key,Value:value,UpdateTime:time.Now()})
				}
				db.DatabaseHandler.UpdateOne(star)
			}
		}
		currLen := len(stars)
		count += currLen
		return currLen == 0
	})
	if err != nil {
		return false
		log.Debug("load star err: %v", err)
	}
	log.Debug("end update star flag count:%v", count)
	return true
}

func (this *starManager) Init() {
	//this.Lock()
	//defer this.Unlock()
	//数据库加载数据到缓存
	if cache.StarCache.SetNX(basecache.FLAG_LOADSTAR, 1) {
		log.Debug("start load star data to redis cache...")
		count := 0
		var stars []*db.DBStar
		//QueryAllLimitEx(data interface{}, result interface{}, limit int, callback func(result)) error
		err := db.DatabaseHandler.QueryAllLimit(&db.DBStar{}, &stars, 10000, func(data interface{}) bool {
			for _, star := range stars {
				if star.Owner != 0{
					starSession := newStarSession(star)
					this.updateRemoteCache(starSession, false)
					if starSession.Active {
						cache.StarCache.SetUserActiveStar(starSession.Owner, starSession.ID)
						cache.StarCache.SetUserActiveStarType(starSession.Owner, starSession.Type)
						cache.StarCache.SetCivilLevel(starSession.ID, starSession.CivilizationLevel)

						//lpc.StatisticsHandler.AddStatisticData(&model.StatisticLogout{
						//	UserID:star.Owner,
						//	BelieverTotal:starSession.believerTotalLevel,
						//	BuildingTotal:starSession.buildingTotalLevel,
						//	Civil:starSession.CivilizationLevel,
						//})
					}
				}
			}
			currLen := len(stars)
			count += currLen
			return currLen == 0
		})
		if err != nil {
			log.Debug("load star err: %v", err)
		}

		log.Debug("end load star data to redis cache count:%v", count)
	}

	if cache.StarCache.SetNX(basecache.FLAG_LOADSTAR_SUMMARY, 1) {
		clustercache.Cluster.CleanSummary()
		count := 0
		var stars []*db.DBStar
		//QueryAllLimitEx(data interface{}, result interface{}, limit int, callback func(result)) error
		db.DatabaseHandler.QueryAllLimit(&db.DBStar{}, &stars, 10000, func(data interface{}) bool {
			for _, star := range stars {
				clustercache.Cluster.AddStarSummary(star.Type)
			}
			currLen := len(stars)
			count += currLen
			return currLen == 0
		})
	}
	//生成机器人
	this.genRobot()
}

func (this *starManager) RemoveStar(uid int32) {
	delete(this.userStars, uid)
}

//加载星球数据 没有从数据库中加载
func (this *starManager) LoadUserStar(uid int32) *UserStarManager {
	//log.Debug("deal star %v", uid)
	if uid < 0 {
		return this.getRobot(uid)
	}
	userStar := this.userStars[uid]
	if userStar != nil {
		return userStar
	}
	var userStarManager *UserStarManager
	stars := this.loadDataFromRemote(uid)
	if stars == nil {
		stars = this.loadDataFromDB(uid)
	}
	if stars == nil || len(stars) == 0 {
		this.userStars[uid] = NewUserStarManager(uid)
		return this.userStars[uid]
	}
	for _, star := range stars {
		starSession := newStarSession(star)
		userStarManager = this.updateActiveStarInfo(starSession, false)
	}
	return userStarManager
}

func (this *starManager) loadDataFromRemote(uid int32) []*db.DBStar {
	node := center.GetServerNode()
	remoteNode := cache.StarCache.GetUserIDNode(uid)
	if remoteNode == "" || node == remoteNode {
		return nil
	}
	log.Debug("star %v data transfer %v -> %v", uid, remoteNode, node)
	resp := rpc.StarServiceProxy.TransmitUserStarInfo(uid, remoteNode)
	if resp == nil || resp.Stars == nil || len(resp.Stars) == 0 {
		return nil
	}
	stars := make([]*db.DBStar, len(resp.Stars))
	for index, starData := range resp.Stars {
		star := &db.DBStar{}
		err := bson.Unmarshal(starData, star)
		if err != nil {
			log.Debug("%v", err.Error())
			return nil
		}
		stars[index] = star
	}
	return stars
}

func (this *starManager)loadDataFromDB(uid int32) []*db.DBStar {
	//log.Debug("load from DB")
	var stars []*db.DBStar
	db.DatabaseHandler.QueryAllCondition(&db.DBStar{}, "owner", uid, &stars)
	return stars
}

func (this *starManager) Close() {
	log.Debug("star session release start count(%v) ....", len(this.userStars ))
	startTime := time.Now()
	//this.Lock()
	//defer this.Unlock()
	for _, userStars := range this.userStars {
		userStars.DealUnSave()
		cache.StarCache.SetUserIDNode(userStars.id, "")
	}
	log.Debug("star session release end duration(%v) ...", time.Now().Sub(startTime).Seconds())
}

//分配一个用户星球
func (this *starManager) AllocUserStar(uid int32, starType int32) *StarSession {
	//starType = conf.RandomStarType()
	//this.Lock()
	//defer this.Unlock()
	//需要走引导后才不禁用
	return this.allocUserStar(uid, starType, true, 1)
}

//处理损坏过期的星球，需要还原为未建造
func (this *starManager) DealStarTimer() {
	//if constant.DEBUG {
	//	log.Debug("current star cache total %v", len(this.userStars))
	//}
	//this.RLock()
	//defer this.RUnlock()

	//上传搜索的索引数据
	rpc.SearchServiceProxy.SyncSearchData()

	for _, starSession := range this.userStars {
		starSession.DealStarTimer()
	}
}

//根据事件类型和用户等阶获取机器人星球
func (this *starManager) RandomEventRobot(eventType int32, levels []int32) []*protocol.Target {
	levelMapping := this.eventRobotRandomMapping[eventType]
	if levelMapping == nil {
		log.Debug("robot not match eventID : %v", eventType)
		return nil
	}

	results := []*protocol.Target{}
	for _, level := range levels {
		weightMapping := levelMapping[int32(level)]
		if weightMapping == nil {
			log.Debug("robot not match eventID : %v level : %v", eventType, level)
			continue
		}
		robotID := util.RandomWeight(weightMapping)
		robotStar := this.getRobotStar(robotID)
		if robotStar == nil {
			continue
		}
		results = append(results, robotStar.BuildTarget())
	}
	return results
}

func (this *starManager) RandomGuideRobot(num int) []*protocol.Target {
	results := []*protocol.Target{}
	current := 0
	for _, robotStar := range this.guideRobotMapping {
		results = append(results, robotStar.GetActiveStar().BuildTarget())
		current ++
		if current == num {
			break
		}
	}
	return results
}

func (this *starManager) GetRobotByID(robotID int32) *protocol.Target {
	return this.getRobotStar(robotID).BuildTarget()
}

//func (this *starManager) SearchStarInfo( uid int32 ) *db.DBStar{//探索星球 wjl 20170531
//	this.RLock()
//	defer this.RUnlock()
//	var star *db.DBStar = nil
//
//	rd := rand.Int()%2
//
//	if rd == 0{//随机如果是0的话 那就搜索已经被占领的星球
//		stars := this.randomOwnStar(uid, 1)
//		if stars != nil && len(stars) > 0 {
//			star = stars[0]
//		}
//	}
//	//if star == nil{//如果为空 好吧 那就搜索一个原始星球出来
//	//	star = this.randOriginalStar()
//	//}
//	if star == nil{
//		exception.GameException(exception.STAR_NOTFOUND)
//	}
//	return star
//}



//func (this *starManager)OccupyStar( uid int32, starID int32 )( *StarSession, *StarSession ){//占领星球 wjl20170605
//	star_src := this.getUserActiveStar( uid )
//	if star_src == nil{
//		exception.GameException(exception.STAR_OCCUPY_FAILED)
//	}
//
//	star_dest := this.getStarByID( starID )
//	if star_dest == nil{
//		exception.GameException(exception.STAR_OCCUPY_FAILED)
//	}
//
//	if !star_src.IsBuildingAllMaxLevel() {
//		exception.GameException(exception.STAR_OCCUPY_FAILED)
//	}
//
//	star_dest.ChangeOwner(uid)
//
//	this.Lock()
//	defer this.Unlock()
//	//删除老的星球记录
//	delete( this.originalStarMapping, starID)
//	delete( this.ownActiveStarMapping, star_src.ID)
//
//	this.ownActiveStarMapping[ starID ] = star_dest
//	cache.StarCache.SetUserActiveStar(uid, starID, star_dest.Type)
//
//
//	//this.ownStarMapping_User[ uid ] = starID
//	db.DatabaseHandler.DeleteOne( star_src.DBStar )//删除数据
//	searcher.Index.RemoveStarIndex(star_src.ID)
//
//	star_dest.BelieverUpdateTime = time.Now()
//	star_dest.setDirty()
//
//	return star_src, star_dest
//}

//func (this *starManager) ResetBuildingGrooves( uid int32, buildingType int32, lockGroove []int32 ) ([]*protocol.ItemGroove, []int32) {
//	star := this.getUserActiveStar( uid )
//	if star == nil {
//		exception.GameException(exception.STAR_NOTFOUND)
//	}
//	grooves, items := star.ResetBuildingGrooves(buildingType, lockGroove)
//
//	if (items != nil) {
//		for _, item := range items {
//			searcher.Index.RemoveItemIndex(star.Type, item, uid)
//		}
//	}
//
//	return grooves, items
//}
