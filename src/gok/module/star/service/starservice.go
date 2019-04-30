//管理网络消息的处理
package service

import (
	"aliens/common/character"
	"aliens/log"
	"github.com/name5566/leaf/chanrpc"
	"gok/constant"
	"gok/module/star/conf"
	"gok/module/star/db"
	"gok/module/star/session"
	"gok/module/star/util"
	baseservice "gok/service"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"time"
)

var StarLocalService = baseservice.NewLocalService(baseservice.SERVICE_STAR_RPC)
var StarRPCService *baseservice.GRPCService = nil

func Init(chanRpc *chanrpc.Server) {
	StarRPCService = baseservice.PublicRPCService1(StarLocalService, conf.Server.RPCAddress, conf.Server.RPCPort, chanRpc)
	//用户RPC远程服务
	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_USER_RPC)
	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_STAR_RPC)
	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_PASSPORT_RPC)
	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_SEARCH_RPC)
}

func Close() {
	StarLocalService.Close()
	StarRPCService.Close()
}

//初始化事件服务消息
func init() {
	RegisterStarService(18, new(GetStarsSelectService))
	RegisterStarService(20, new(GetStarInfoService))
	RegisterStarService(21, new(UpgradeBuildingService)) //请求升级星球建筑物
	//RegisterStarService(47, new(SelectAreaService))

	RegisterStarService(67, new(LootFaithService))       //掠夺信仰
	RegisterStarService(68, new(AtkStarBuildingService)) //被破坏星球建筑

	RegisterStarService(263, new(LootItemService))         //掠夺圣物
	RegisterStarService(270, new(LootStarBelieverService)) //被破坏星球建筑

	//RegisterStarService(301, new(SearchStarInfoService))    //请求探索星球 wjl20170531
	RegisterStarService(302, new(GetStarInfoDetailService)) //请求星球详细信息 wjl20170603
	RegisterStarService(303, new(GetStarShieldService))     //获取防护罩信息

	//RegisterStarService(311, new(GetStarRecordService))     //请求获取星球记录数据 wjl 20170601
	//RegisterStarService(320, new(OccupyStarService))        //请求占领星球记录 wjl 20170605
	//RegisterStarService(321, new(NextStarService)) //跳转下一个星球

	//	StarServiceManager.RegisterLocalService(330, new(RandomTargetStarService))     //请求随机目标星球 wjl 20170621
	RegisterStarService(321, new(GetHelpRepairInfoService))
	RegisterStarService(325, new(StarSettleService))
	RegisterStarService(326, new(StarFlagInfoService)) // 星球标识
	RegisterStarService(327, new(UpdateStarFlagService)) // 更新星球标识
	RegisterStarService(328, new(updateAllStarFlagService)) //更新所有星球标识

	RegisterStarService(339, new(CancelUpgradeStarBuildService))
	RegisterStarService(340, new(CancelRepairStarBuildService))
	RegisterStarService(341, new(AccRepairStarBuildingService)) //请求加速维修星球建筑
	RegisterStarService(342, new(RepairStarBuildingService))    //请求开始维修星球建筑
	RegisterStarService(350, new(UpdateStarBuildingEndService)) //请求升级建筑物结束
	RegisterStarService(351, new(AccUpgradeBuildingService))    //请求升级建筑物加速

	RegisterStarService(355, new(BuildingFaithInfoService))    //查看建筑信仰值
	RegisterStarService(356, new(ReceiveBuildingFaithService)) //领取建筑信仰值
	RegisterStarService(357, new(RepairStarBuildEndService))   //请求修理建筑结束

	RegisterStarService(360, new(RandomEventRobotService)) ////随机事件机器人
	RegisterStarService(361, new(RandomGuideRobotService))  //随机引导机器人
	RegisterStarService(609, new(GetEventRobotService)) //根据id获取事件机器人

	//RegisterStarService(360, new(TakeInItemBuildingService))  //建筑放入圣物
	//RegisterStarService(361, new(TakeoutItemBuildingService)) //建筑取出圣物
	//RegisterStarService(362, new(ResetBuildingGrooveService)) //领取建筑信仰值
	//RegisterStarService(363, new(ActiveBuildingGroupService))     //物品放入建筑

	//RegisterStarService(364, new(AccBuildingGrooveEffectService)) //物品放入建筑
	//RegisterStarService(365, new(UpdateGrooveEffectService))      //更新槽CD

	RegisterStarService(368, new(AddCivilizationService))        //添加文明度
	RegisterStarService(369, new(DrawCivilizationRewardService)) //领取文明度奖励

	RegisterStarService(370, new(UpdateStatisticsService)) //更新统计数据
	RegisterStarService(371, new(StatisticsInfoService))   //请求星球统计信息

	RegisterStarService(380, new(StarHistoryService)) //请求星球统计信息

	RegisterStarService(400, new(UpgradeBelieverService)) //升级信徒
	RegisterStarService(402, new(AutoAddBelieverService)) //自动加信徒
	RegisterStarService(403, new(UpdateBelieverInfoService)) //客户端同步信徒信息

	RegisterStarService(451, new(GetItemGroupService))        //获取图鉴信息
	RegisterStarService(452, new(DrawItemGroupRewardService)) //领取图鉴奖励
	RegisterStarService(453, new(ActiveGroupService))         //图鉴开启新物品

	RegisterStarService(600, new(AllocNewStarService)) //分配一个初始化的星球
	//RegisterStarService(602, new(UserStarInfoService))     //获取用户星球信息
	RegisterStarService(603, new(UpdateBelieverService))        //更新信徒信息
	RegisterStarService(604, new(LoginStarInfoService))         //获取星球信息
	RegisterStarService(605, new(TransmitUserStarService))      //
	RegisterStarService(606, new(PublicHelpRepairBuildService)) // 发布求助修理建筑
	RegisterStarService(607, new(HelpRepairBuildService))       //帮助维修建筑
	RegisterStarService(608, new(GetCurrentGroupItemsService))  //获取当前正在解锁的圣物组合的圣物ID
	RegisterStarService(610, new(GetOwnersByConditionService))
	//----------------测试用例-----------------------------
	RegisterStarService(667, new(SetBuildsService))
	RegisterStarService(668, new(SetBelieversService))
}



//注册用户消息服务
func RegisterStarService(seq int32, service IStarService) {
	StarLocalService.RegisterHandler(seq, &StarServiceProxy{proxy: service})
}

//服务代理
type StarServiceProxy struct {
	proxy IStarService
}

func (this StarServiceProxy) Request(request *protocol.C2GS, response *protocol.GS2C, channel baseservice.IMessageChannel) {
	uid := request.GetParam()
	if uid == 0 {
		this.proxy.Request(request, response,nil, nil)
	} else {
		userStars := session.StarManager.LoadUserStar(uid)

		if userStars == nil || !userStars.HaveActiveStar(){
			if request.GetLoginStarInfo() == nil && request.GetGetStarInfo() == nil && request.GetAllocNewStar() == nil{
				exception.GameException(exception.STAR_NOTFOUND)
			}
		}
		userStars.UpdateActiveTime()
		star := userStars.GetActiveStar()
		this.proxy.Request(request, response, userStars, star)
	}
}

type IStarService interface {
	Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession)
}

//分配一个初始星球
type AllocNewStarService struct {
}

func (service *AllocNewStarService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetAllocNewStar()

	if userStarManager != nil && userStarManager.HaveActiveStar() && !userStarManager.GetActiveStar().IsDoneStar() {
		exception.GameException(exception.SELECT_STAR_REPEAT)
	}

	var star *session.StarSession
	var lastStarType int32
	var faith int32
	var isFirst bool
	if userStarManager.HaveActiveStar() {
		star, lastStarType, faith = userStarManager.NextStar(message.GetStarType())
		isFirst = false
	} else {
		// 新手星球结束后的第一个星球
		star = session.StarManager.AllocUserStar(message.GetUid(), message.GetStarType())
		isFirst = true
	}
	//resp.Star = session.StarManager.GetStarInfoByUser(star.Owner)
	//response.LoginStarInfoRet = resp

	response.AllocNewStarRet = &protocol.AllocNewStarRet{
		Star:userStarManager.GetStarInfo(),
		CurrentStar:star.BuildProtocol(),
		LastStarType:lastStarType,
		Faith:faith,
		IsFirst:isFirst,
	}
}

type StarHistoryService struct {
}

func (service *StarHistoryService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) { //查询星球详细信息
	message := request.GetGetStarHistory()
	if message.GetStarID() > 0 {
		starSession = userStarManager.GetStar(message.GetStarID())
	}
	if starSession == nil {
		exception.GameException(exception.STAR_NOTFOUND)
	}

	response.GetStarHistoryRet = &protocol.GetStarHistoryRet{
		StarID:  starSession.ID,
		Uid:     starSession.Owner,
		History: starSession.GetStarHistory(),
		Time:	 starSession.GetHistoryTime(),
	}
}

type UpgradeBuildingService struct {
}

func (service *UpgradeBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	building, upgrade, powerLimit, cost, civilizationInfo := starSession.UpgradeBuilding(request.GetBuildStarBuilding())

	//starSession.DealItemGroupOpen(building.Type,building.Level)

	result := &protocol.BuildStarBuildingRet{
		Done:             upgrade,
		Result:           true,
		UpdateTime:       building.UpdateTime.Unix(),
		PowerLimit:       powerLimit,
		Cost:             cost,
		CivilizationInfo: civilizationInfo,
	}
	//升级需要返回建筑最新信息
	if upgrade {
		result.Builidng = building.BuildProtocol()
	}
	response.BuildStarBuildingRet = result
}

type GetStarsSelectService struct {

}

func (service *GetStarsSelectService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetGetStarsSelect()
	var starsType []int32
	confStars := conf.Base.StarWeCanArrive
	if confStars == nil {
		exception.GameException(exception.GAME_BASE_NOT_FOUND)
	}
	if message.GetNum() == -1 {
			starsType = make([]int32, len(confStars))
			starsType = confStars
	} else {
		starsType = conf.RandomStarsType(confStars, message.GetNum(), session.StarManager.GetStarInfoByUser(message.GetUid()))
	}
	response.GetStarsSelectRet = &protocol.GetStarsSelectRet{StarsType:starsType}
}

type GetStarInfoService struct {
}

func (service *GetStarInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetGetStarInfoDetail()
	if message.GetStarID() > 0 {
		starSession = userStarManager.GetStar(message.GetStarID())
	}
	response.GetStarInfoRet = &protocol.GetStarInfoRet{}
	if starSession != nil {
		currStar := starSession.BuildProtocol()
		response.GetStarInfoRet.CurrentStar = currStar
		response.GetStarInfoRet.Star = userStarManager.GetStarInfo()
	}
}

//type SearchStarInfoService struct {
//	//20170531 wjl 搜索星球信息
//}
//
//func (service *SearchStarInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
//	message := request.GetSearchStarInfo()
//	star := session.StarManager.SearchStarInfo(message.GetUid()) //搜索星球
//	if star != nil {
//		buildingInfo := util.BuildStarBuildingInfo(star)
//		response.SearchStarInfoRet = &protocol.SearchStarInfoRet{
//			Star: &protocol.StarInfoDetail{
//				starID:   star.ID),
//				Type:     star.Type),
//				Building: buildingInfo,s
//				OwnID:    star.Owner),
//			},
//		}
//	}
//}

type GetStarInfoDetailService struct {
	//20170603 wjl 获取星球详细数据
}

func (service *GetStarInfoDetailService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetGetStarInfoDetail()
	if message.GetStarID() > 0 {
		starSession = userStarManager.GetStar(message.GetStarID())
	}
	if starSession == nil {
		exception.GameException(exception.STAR_NOTFOUND)
	}
	starInfo := starSession.BuildProtocol()
	if message.IsConvert {
		if !starSession.CheckBuildingLevel() {
			AllocateBuildingsByLevel(starInfo.GetBuilding(), message.GetBuildingTotalLevel())
		}
		if !starSession.CheckBelieverLevel() {
			AllocateBelieversByLevel(starInfo.GetBeliever(), message.GetBelieverTotalLevel(), starInfo.Type)
		}
	}
	response.GetStarInfoDetailRet = &protocol.GetStarInfoDetailRet{
		Star: starInfo,
		Shield: starSession.GetShield(message.GetShieldType()),
	}
}

func AllocateBuildingsByLevel(buildings []*protocol.BuildingInfo, levelTotal int32) {
	num := len(buildings)
	for i := 0; i < int(levelTotal); i ++ {
		buildings[i % num].Level ++
	}
}

func AllocateBelieversByLevel(believers []*protocol.BelieverInfo, levelTotal int32, starType int32) {
	//var addTotalLevel int32
	levelMax := constant.BELIEVER_L3
	for i:= levelMax; i > 0; i-- {
		num := int32(levelTotal/i)
		levelTotal = levelTotal%i
		if num > 0 {
			believerID := getBelieverID(starType, i, 1 + rand.Int31n(2))
			believers = append(believers, &protocol.BelieverInfo{Id:believerID,Num:num})
		}
	}
}

func getBelieverID (starType int32, level int32, sex int32) string {
	typeStr := character.Int32ToString(starType)
	if starType < 10 {
		typeStr = "0" + typeStr
	}
	return "b" + typeStr + character.Int32ToString(level) + character.Int32ToString(sex)
}



type GetStarShieldService struct {
	//获取星球防护罩数据
}

func (service *GetStarShieldService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	star := userStarManager.GetActiveStar()
	if star == nil {
		exception.GameException(exception.STAR_NOTFOUND)
	}
	response.GetStarShieldRet = star.GetShieldInfo()
}

type StatisticsInfoService struct {
	//获取星球统计数据
}

func (service *StatisticsInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	response.GetStarStatisticsRet = &protocol.GetStarStatisticsRet{
		Uid:        starSession.Owner,
		StarID:     starSession.ID,
		Statistics: starSession.GetStatisticsInfo(),
	}
}

//type GetStarRecordService struct {
//	//20170601 wjl 获取星球记录信息
//}
//
//func (service *GetStarRecordService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
//	message := request.GetGetStarRecordInfo()
//
//	var star *session.StarSession = nil
//
//	oriStarInfo := []*protocol.StarInfoBase{}
//	for _, v := range message.GetStarID() {
//		star = session.StarManager.GetStarByID(v)
//		if star != nil {
//			star_info := &protocol.StarInfoBase{
//				starID:      star.ID),
//				Type:        star.Type),
//				OwnID:       star.Owner),
//				OwnNikeName: ""),
//			}
//			oriStarInfo = append(oriStarInfo, star_info)
//		}
//	}
//	userStarInfo := []*protocol.StarInfoBase{}
//	for _, v := range message.GetUserID() {
//		star = session.StarManager.GetUserActiveStar(v)
//		if star != nil {
//			star_info := &protocol.StarInfoBase{
//				starID:      star.ID),
//				Type:        star.Type),
//				OwnID:       star.Owner),
//				OwnNikeName: ""),
//			}
//			userStarInfo = append(userStarInfo, star_info)
//		}
//	}
//	response.GetStarRecordInfoRet = &protocol.GetStarRecordInfoRet{
//		StarsOri:  oriStarInfo,
//		StarsUser: userStarInfo,
//	}
//}

//type OccupyStarService struct {
//}
//
//func (service *OccupyStarService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
//	message := request.GetOccupyStar()
//	srcStar, destStar := session.StarManager.OccupyStar(message.GetUid(), message.GetStarID())
//
//	if srcStar != nil {
//		if destStar != nil {
//			srcStarInfo := util.BuildStarDetailInfo(srcStar.DBStar)
//			destStarInfo := util.BuildStarDetailInfo(destStar.DBStar)
//
//			response.OccupyStarRet = &protocol.OccupyStarRet{
//				Uid:     message.GetUid()),
//				StarOld: srcStarInfo,
//				Star:    destStarInfo,
//			}
//		}
//
//	}
//}

//type NextStarService struct {
//}
//
//func (service *NextStarService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
//	//message := request.GetNextStar()
//	_, lastStarType, faith := userStarManager.NextStar()
//	//nextStarInfo := nextStar.BuildProtocol()
//	response.NextStarRet = &protocol.NextStarRet{
//		LastStarType: lastStarType,
//		Star:         nil,
//		Faith:        faith,
//	}
//}

type LootFaithService struct {
	//进攻星球
}

func (service *LootFaithService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	// 交互没有解锁或者建筑满级
	max := starSession.IsBuildingAllMaxLevel() || !starSession.IsFlagUnlock(constant.STAR_FLAG_MUTUAL)
	shield := starSession.DecFaithShield()
	response.LootFaithRet = &protocol.LootFaithRet{
		Shield: shield,
		HasBuilding:starSession.HasBuilding(),
		IsMax:max,
	}
}

type AtkStarBuildingService struct {
	//进攻星球
}

func (service *AtkStarBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetAtkStarBuilding()
	success := message.GetSuccess()
	faith, itemID, shield, max := starSession.AtkStar(message)
	if shield {
		success = false
	}
	response.AtkStarBuildingRet = &protocol.AtkStarBuildingRet{
		Success: success,
		Faith:   faith,
		ItemID:  itemID,
		Shield:  shield,
		IsMax:   max,
	}
}

type CancelUpgradeStarBuildService struct {
}

func (service *CancelUpgradeStarBuildService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetCancelUpgradeStarBuild()
	response.CancelUpgradeStarBuildRet = starSession.CancelUpgradeStarBuilding(message.GetBuildingType())
}

type CancelRepairStarBuildService struct {
}

func (service *CancelRepairStarBuildService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetCancelRepairStarBuild()
	response.CancelRepairStarBuildRet = starSession.CancelRepairStarBuilding(message.GetBuildingType())
}

//快速维修星球建筑，需要消耗信徒
type AccRepairStarBuildingService struct {
}

func (service *AccRepairStarBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetAccRepairStarBuild()
	repairSuccess, repairTime, believerCost, buildingLevel, civilizationInfo := starSession.AccRepairStarBuilding(message.GetBuildingType(), message.GetBelieverId())
	response.AccRepairStarBuildRet = &protocol.AccRepairStarBuildRet{
		Done:             repairSuccess,
		Uid:              message.GetUid(),
		BuildingType:     message.GetBuildingType(),
		RepairTime:       repairTime,
		BelieverNum:      believerCost,
		BuildingLevel:    buildingLevel,
		CivilizationInfo: civilizationInfo,
	}
}

type RepairStarBuildingService struct {
	//请求开始维修建筑
}

func (service *RepairStarBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetRepairStarBuild()
	resp := starSession.RepairStarBuilding(message.GetUid(), message.GetBuildingType(), message.GetFaith() /*, message.GetBelieverId()*/)
	response.RepairStarBuildRet = resp
}

//请求升级建筑物结束
type UpdateStarBuildingEndService struct {
}

func (service *UpdateStarBuildingEndService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetUpdateStarBuildEnd()
	resp := starSession.UpgradeStarBuildEnd(message.GetUid(), message.GetBuildingType())
	response.UpdateStarBuildEndRet = resp
}

type AccUpgradeBuildingService struct {
	//请求升级建筑物加速
}

func (service *AccUpgradeBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetAccUpdateStarBuild()
	resp := starSession.AccUpgradeStarBuilding(message.GetUid(), message.GetBuildingType(), message.GetBelieverId(), message.GetGuide())
	response.AccUpdateStarBuildRet = resp
}

type UpgradeBelieverService struct {
	//升级信徒
}

func (service *UpgradeBelieverService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetUpgradeBeliever()
	believer, cost, civilization := starSession.UpgradeBeliever(message.GetSelectID(), message.GetMatchID(), message.GetFaith())

	response.UpgradeBelieverRet = &protocol.UpgradeBelieverRet{
		Result:           util.BuildBelieverInfo(believer),
		Cost:             cost,
		CivilizationInfo: civilization,
	}
}

type AutoAddBelieverService struct {
}

func (service *AutoAddBelieverService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetAutoAddBeliever()
	believer := starSession.AutoAddBeliever(message.Times)
	response.AutoAddBelieverRet = &protocol.AutoAddBelieverRet{
		Believer:     believer,
		BelieverTime: starSession.BelieverUpdateTime.Unix(),
	}
}

type UpdateBelieverInfoService struct {

}

func (service *UpdateBelieverInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	believers := starSession.GetStarBelieversInfo()
	var believersInfo []*protocol.BelieverInfo
	for _, believer := range believers {
		believersInfo = append(believersInfo, believer.BuildProtocol())
	}
	response.UpdateBelieverInfoRet = &protocol.UpdateBelieverInfoRet{
		Believer:believersInfo,
	}
}

type UpdateBelieverService struct {
	//更新信徒
}

func (service *UpdateBelieverService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetUpdateBeliever()
	believers, faith := starSession.UpdateBeliever(message.GetOperation(), message.GetUpdateInfo(), message.GetIsConvert())

	results := make([]*protocol.BelieverInfo, len(believers))
	for index, believer := range believers {
		results[index] = util.BuildBelieverInfo(believer)
	}
	response.UpdateBelieverRet = &protocol.UpdateBelieverRet{
		Result:   true,
		Believer: results,
		Faith: faith,
	}
}

type LoginStarInfoService struct {
}

func (service *LoginStarInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	if userStarManager.HaveActiveStar() {
		resp := starSession.LoginStarInfo()
		resp.Star = session.StarManager.GetStarInfoByUser(starSession.Owner)
		resp.StarFlags = starSession.GetProtocolFlags()
		response.LoginStarInfoRet = resp
		starSession.Push = false

	} else {
		response.LoginStarInfoRet = &protocol.LoginStarInfoRet{Star:session.StarManager.GetStarInfoByUser(request.GetLoginStarInfo().Uid)}
	}

}

//查看建筑信仰值信息
type BuildingFaithInfoService struct {
	//更新信徒
}

func (service *BuildingFaithInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetGetBuildingFaith()
	//faith, faithUpdateTime := starSession.GetBuildingFaith(message.GetBuildingType(), false)
	buildingsFaith := starSession.GetAllBuildingFaith(message.GetBuildingType())
	response.GetBuildingFaithRet = &protocol.GetBuildingFaithRet{
		Uid:             message.Uid,
		BuildingsFaith:  buildingsFaith,
		//BuildingType:    message.BuildingType,
		//BuildingFaith:   faith,
		//FaithUpdateTime: faithUpdateTime,
	}
}

//领取建筑信仰值
type ReceiveBuildingFaithService struct {
}

func (service *ReceiveBuildingFaithService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetReceiveBuildingFaith()
	faith, faithUpdateTime := starSession.GetBuildingFaith(message.GetBuildingType(), true)
	starSession.AddStatisticsValue(constant.STAR_STATISTIC_GAIN_FAITH_BUILDING, float64(faith), 0)
	response.ReceiveBuildingFaithRet = &protocol.ReceiveBuildingFaithRet{
		Result:          true,
		BuildingFaith:   faith,
		FaithUpdateTime: faithUpdateTime,
	}
}

type RepairStarBuildEndService struct {
}

func (service *RepairStarBuildEndService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetRepairStarBuildEnd()
	result, civilizationInfo, repairTime, buildingLevel := starSession.RepairBuildingEnd(message.GetBuildingType())
	response.RepairStarBuildEndRet = &protocol.RepairStarBuildEndRet{
		Done:             result,
		RepairTime:       repairTime,
		BuildingLevel:    buildingLevel,
		CivilizationInfo: civilizationInfo,
	}
}

type RandomEventRobotService struct {
}

func (service *RandomEventRobotService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetRandomEventRobot()
	targets := session.StarManager.RandomEventRobot(message.GetEventType(), message.GetLevel())
	response.RandomEventRobotRet = &protocol.RandomEventRobotRet{
		Targets: targets,
	}
}


type RandomGuideRobotService struct {
}

func (service *RandomGuideRobotService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetRandomGuideRobot()
	targets := session.StarManager.RandomGuideRobot(int(message.GetNum()))
	response.RandomGuideRobotRet = &protocol.RandomGuideRobotRet{
		Targets:targets,
	}
}

type GetEventRobotService struct {
}

func (service *GetEventRobotService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetGetEventRobot()
	target := session.StarManager.GetRobotByID(message.GetUid())
	response.GetEventRobotRet = &protocol.GetEventRobotRet{Target:target}
}
//添加圣物
//type TakeInItemBuildingService struct {
//}
//
//func (service *TakeInItemBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
//	message := request.GetPutItemToBuilding()
//	grooves := starSession.ReplaceBuildingItem(message.GetBuildingType(), message.GetItemGroove(), message.GetItemID(), message.GetTakeoutItem())
//	response.PutItemToBuildingRet = &protocol.PutItemToBuildingRet{
//		Result:       true,
//		BuildingType: message.GetBuildingType(),
//		ItemGroove:   grooves,
//	}
//}

//去除圣物
//type TakeoutItemBuildingService struct {
//}
//
//func (service *TakeoutItemBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
//	message := request.GetRemoveItemFromBuilding()
//	itemID := starSession.TakeoutBuildingItem(message.GetBuildingType(), message.GetItemGroove())
//	response.RemoveItemFromBuildingRet = &protocol.RemoveItemFromBuildingRet{
//		BuildingType: message.GetBuildingType(),
//		ItemID:       itemID,
//		Result:       true,
//	}
//
//}

//type ResetBuildingGrooveService struct{//重置槽
//
//}
//
//func ( service *ResetBuildingGrooveService )Request(request *protocol.C2GS, response *protocol.GS2C, star *session.StarSession) {
//	message := request.GetResetBuildingGroove()
//	grooves, returnItems := session.StarManager.ResetBuildingGrooves(message.GetUid(), message.GetBuildingType(), message.GetLockGroove())
//	response.ResetBuildingGrooveRet = &protocol.ResetBuildingGrooveRet{
//		BuildingType : message.BuildingType,
//		ItemGroove : grooves,
//		ItemID:returnItems,
//	}
//}

//type ActiveBuildingGroupService struct {
//}

//func (service *ActiveBuildingGroupService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
//	message := request.GetActiveBuildingGroup()
//	returnItems := starSession.ActiveBuildingGroup(message.GetBuildingType())
//	response.ActiveBuildingGroupRet = &protocol.ActiveBuildingGroupRet{
//		BuildingItems: returnItems,
//	}
//}

//type AccBuildingGrooveEffectService struct {
//}
//
//func (service *AccBuildingGrooveEffectService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
//	message := request.GetAccBuildingGrooveEffect()
//	timestamp, effect := starSession.AccBuildingGrooveEffect(message.GetBuildingType(), message.GetItemGroove(), message.GetBeliever())
//	response.AccBuildingGrooveEffectRet = &protocol.AccBuildingGrooveEffectRet{
//		Effect:          effect,
//		EffectTimestamp: timestamp,
//	}
//}

//type UpdateGrooveEffectService struct {
//}
//
//func (service *UpdateGrooveEffectService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
//	message := request.GetUpdateGrooveEffect()
//	timestamp, effect := starSession.ActiveGrooveEffect(message.GetBuildingType(), message.GetItemGroove())
//	response.UpdateGrooveEffectRet = &protocol.UpdateGrooveEffectRet{
//		Effect:          effect,
//		EffectTimestamp: timestamp,
//	}
//}

type AddCivilizationService struct {
}

func (service *AddCivilizationService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	response.AddCivilizationRet = starSession.TakeInCivilization(request.GetAddCivilization().GetCivilizationValue())
}

type DrawCivilizationRewardService struct {
}

func (service *DrawCivilizationRewardService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetDrawCivilizationReward()
	if message.GetStarID() > 0 {
		starSession = userStarManager.GetStar(message.GetStarID())
	}
	if starSession == nil {
		exception.GameException(exception.STAR_NOTFOUND)
	}
	//starSession.Lock()
	//defer starSession.Unlock()
	reward, faithValue, diamondValue, relicPointValue, believer, believerNum := starSession.DrawStarReward(message.GetDrawLevel())
	response.DrawCivilizationRewardRet = &protocol.DrawCivilizationRewardRet{
		StarID:      message.GetStarID(),
		DrawLevel:   message.GetDrawLevel(),
		Reward:      reward,
		Faith:       faithValue,
		Diamond:     diamondValue,
		GayPiont:    relicPointValue,
		Believer:    believer,
		BelieverNum: believerNum,
	}
}

type UpdateStatisticsService struct {
}

func (service *UpdateStatisticsService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetUpdateStarStatistics()
	result := starSession.AddStatisticsValue(message.GetId(), message.GetChange(), message.GetParam())
	response.UpdateStarStatisticsRet = &protocol.UpdateStarStatisticsRet{
		Value: result.Value,
	}
}



type LootItemService struct {
	//重置槽
}

func (service *LootItemService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetLootItem()
	building := starSession.LootBuildingItem(message.GetBuildingID(), message.GetItemID())

	var buildingType int32 = 0
	if building != nil {
		buildingType = building.Type
	}
	response.LootItemRet = &protocol.LootItemRet{
		Result:   building != nil,
		Building: buildingType,
	}
}

//抢夺信徒
type LootStarBelieverService struct {
}

func (service *LootStarBelieverService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	believerID, shield, max := starSession.LootBeliever(request.GetLootStarBeliever().GetBelieverID())
	response.LootStarBelieverRet = &protocol.LootStarBelieverRet{
		BelieverID: believerID,
		Shield:     shield,
		IsMax:		max,
	}
}

//获取圣物组合
type GetItemGroupService struct {
}

func (service *GetItemGroupService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	if starSession.ItemGroups == nil || len(starSession.ItemGroups) == 0 {
		starSession.ItemGroupsInit()
	}
	response.GetItemGroupRet = &protocol.GetItemGroupRet{
		ItemGroup: starSession.GetProtocolItemGroups(),
	}
}

//领取圣物组合奖励
type DrawItemGroupRewardService struct {
}

func (service *DrawItemGroupRewardService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetGetGroupReward()

	groupBase := conf.Base.ItemGroup[message.GetGroupID()]
	if (groupBase == nil) {
		exception.GameException(exception.ITEM_GOURP_BASE_NOT_FOUND)
	}
	group := starSession.GetItemGroup(message.GetGroupID())
	if (group == nil) {
		exception.GameException(exception.ITEM_GOURP_NOT_FINISH)
	}
	//if (group.Reward) {
	//	exception.GameException(exception.ITEM_GOURP_REWARD_REPEAT)
	//}
	if (!group.Done) {
		exception.GameException(exception.ITEM_GOURP_NOT_FINISH)
	}
	response.GetGroupRewardRet = &protocol.GetGroupRewardRet{
		GroupID: message.GroupID,
	}

	//reward

	//group.DoneItemGroup()
	//group.Reward = true

}

//尝试圣物组合
type ActiveGroupService struct {
}

func (service *ActiveGroupService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetActiveGroup()
	response.ActiveGroupRet = starSession.ActiveGroup(message.GetGroupID(), message.GetItemID())

}

type TransmitUserStarService struct {
}

func (service *TransmitUserStarService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager1 *session.UserStarManager, starSession1 *session.StarSession) {
	message := request.GetTransmitUserStar()
	userStarManager := session.StarManager.GetUserStarsByUser(message.GetUserID())
	resp := &protocol.TransmitUserStarRet{}
	var stars [][]byte
	for _, starSession := range userStarManager.GetOwnStars() {
		ownStar, err := bson.Marshal(starSession.DBStar)
		if err != nil {
			log.Debug("%v",err.Error())
		}
		stars = append(stars, ownStar)
	}
	//log.Debug("%v",stars)
	session.StarManager.RemoveStar(message.GetUserID())
	resp.Stars = stars
	response.TransmitUserStarRet = resp
}

type PublicHelpRepairBuildService struct {

}

func (service *PublicHelpRepairBuildService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetHelpRepairBuildPublic()
	result := starSession.PublicHelpRepairBuild(message.GetBuildingType())
	response.HelpRepairBuildPublicRet = &protocol.HelpRepairBuildPublicRet{BuildingType:result}
}

type GetHelpRepairInfoService struct {

}

func (service *GetHelpRepairInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetGetHelpRepairInfo()
	response.GetHelpRepairInfoRet = &protocol.GetHelpRepairInfoRet{HelpRepairBuildInfo: starSession.GetHelpBuildsInfo(message.GetBuildType())}
}

type HelpRepairBuildService struct {

}

func (service *HelpRepairBuildService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetHelpRepairBuild()
	ret, starType, buildingLevel := starSession.HelpRepairBuild(message.GetBuildingType(), message.GetHelperID())
	response.HelpRepairBuildRet = &protocol.HelpRepairBuildRet{
		Result:ret,
		StarType:starType,
		BuildingLevel:buildingLevel,
	}
}

type GetCurrentGroupItemsService struct {

}

func (service *GetCurrentGroupItemsService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	itemIDs, groupID := starSession.GetCurrentActiveGroupItems()
	response.GetCurrentGroupItemsRet = &protocol.GetCurrentGroupItemsRet{
		ItemIDs:itemIDs,
		GroupID:groupID,
	}
}

type GetOwnersByConditionService struct {

}

func (service *GetOwnersByConditionService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetGetOwnersByCondition()
	//time.Parse("2016-01-01",character.Int64ToString(message.GetStart()))
	start := time.Unix(message.GetStart(), 0)
	end := time.Unix(message.GetEnd(), 0)
	buildLv := message.GetBuildLv()
	skip := message.GetSkip()
	limit := message.GetLimit()
	//time
	//end := time.Now()
	query := make(bson.M)
	//timeQuery := make(bson.M)
	//time.Unix(end.Unix() -8*3600, 0).UTC()
	//timeQuery["$lt"] = end
	//timeQuery["lte"] = end
	query["createtime"] = bson.M{"$gt":start, "$lt":end}
	if buildLv >= 0 {
		query["buildMaxLevel"] = buildLv
	}

	var stars []*db.DBStar
	count, err := db.DatabaseHandler.QueryConditionsCount(&db.DBStar{}, query)
	if err != nil {
		log.Debug("queryConditionsCount err:%v", err.Error())
	}
	err1 := db.DatabaseHandler.QueryAllConditionsSkipLimit(&db.DBStar{}, query, &stars, int((skip-1)*limit), int(limit), "-createtime")
	if err1 != nil {
		log.Debug("QueryAllConditionsSkipLimit err:%v", err1.Error())
	}
	userData := make([]*protocol.UserStarData, len(stars))
	for index, star := range stars {
		results := make([]*protocol.BelieverInfo, len(star.Believer))
		for index, believer := range star.Believer {
			results[index] = util.BuildBelieverInfo(believer)
		}
		statistic := make([]*protocol.Statistics, len(star.Statistics))
		for index, statis := range star.Statistics {
			statistic[index] = statis.BuildProtocol()
		}
		userData[index] = &protocol.UserStarData{Believers:results, Uid:star.Owner, StarStatis:statistic}
	}

	response.GetOwnersByConditionRet = &protocol.GetOwnersByConditionRet{
		Count:int32(count),
		UserData:userData,
	}
}

type StarSettleService struct {

}

func (service *StarSettleService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	starSession.IsDoneStar()
	passTime := starSession.GetSettleTime()
	response.StarSettleRet = &protocol.StarSettleRet{PassTime:passTime}

}

type StarFlagInfoService struct {

}

func (service *StarFlagInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	response.StarFlagInfoRet = &protocol.StarFlagInfoRet{
		Flags:starSession.GetProtocolFlags(),
	}
}

type UpdateStarFlagService struct {

}

func (service *UpdateStarFlagService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetUpdateStarFlag()
	flag := starSession.UpdateFlagValue(message.GetFlag(), message.GetValue())
	response.UpdateStarFlagRet = &protocol.UpdateStarFlagRet{Flag:flag.BuildProtocol()}
}

type updateAllStarFlagService struct {

}

func (service *updateAllStarFlagService) Request(request *protocol.C2GS, response *protocol.GS2C, userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetUpdateAllStarFlag()
	result := session.StarManager.UpdateAllStarFlag(message.GetKey(), message.GetValue())
	response.UpdateAllStarFlagRet = &protocol.UpdateAllStarFlagRet{Result:result}
}

//--------------------------------------压测----------------------------------------------
type SetBuildsService struct {

}

func (service *SetBuildsService) Request(request *protocol.C2GS, response *protocol.GS2C,userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetSetBuildings()
	//star := userStarManager.GetStar(message.GetStarID())
	star := userStarManager.GetActiveStar()
	result := star.SetBuildingsLevel(message.GetLevel())
	response.SetBuildingsRet = &protocol.SetBuildingsRet{
		Result:result,
	}
}

type SetBelieversService struct {

}

func (service *SetBelieversService) Request(request *protocol.C2GS, response *protocol.GS2C,userStarManager *session.UserStarManager, starSession *session.StarSession) {
	message := request.GetSetBelievers()
	star := userStarManager.GetActiveStar()
	result := star.SetBelieversLevel(message.GetLevel())
	response.SetBelieversRet = &protocol.SetBelieversRet{
		Believers:result,
	}
}