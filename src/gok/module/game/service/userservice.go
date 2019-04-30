//管理网络消息的处理
package service

import (
	"gok/service/msg/protocol"
	"gok/service"
	"gok/service/exception"
	"gok/module/game/cache"
	"gok/module/game/conf"
	"gok/module/game/db"
	"time"
	"gok/module/game/user"
	"gok/module/game/util"
	"gok/service/rpc"
	"aliens/common/character"
	"gok/constant"
	"math/rand"
	"gok/module/game/words"
	"gok/module/game/rank"
	"gok/module/game/global"
	"gok/service/lpc"
	"gok/module/statistics/model"
)

//处理用户客户端消息
var UserService = service.NewLocalService(service.SERVICE_USER)

//初始化用户服务容器
func init() {
	RegisterUserService(4, new(HeartbeatService))
	RegisterUserService(5, new(ServerTimeService))

	RegisterUserService(8, new(GenOrderService))

	RegisterUserService(15, new(GetAvatarService))
	RegisterUserService(16, new(ChangeDescService))
	RegisterUserService(17, new(GetRoleInfoService))

	RegisterUserService(22, new(DisplayInfoService))
	RegisterUserService(23, new(UpdateDisplayService))

	RegisterUserService(25, new(BuyShopItemService))
	RegisterUserService(26, new(RoleFlagInfoService))
	RegisterUserService(27, new(UpdateRoleFlagService))
	RegisterUserService(28, new(UpdatePowerService))


	//RegisterUserService(36, new(AssistEventRequestListService))
	//RegisterUserService(38, new(RejectAssistEventService))
	//RegisterUserService(47, new(SelectAreaService))

	//RegisterUserService(61, new(AddGayPointService))
	//RegisterUserService(62, new(FullSearchService))


	RegisterUserService(70, new(GetAllMailService)) //获取所有邮件
	RegisterUserService(71, new(DrawMailService))   //领取邮件
	RegisterUserService(72, new(RemoveMailService)) //删除邮件

	//引导任务
	RegisterUserService(80, new(GuideTaskService))
	RegisterUserService(81, new(GuideBuildingFaithService))
	RegisterUserService(82, new(GuideRevengeService)) //引导复仇

	RegisterUserService(101, new(RandomEventService))
	RegisterUserService(102, new(TaskListService))
	//RegisterUserService(103, new(UpdateTaskEndingService))
	RegisterUserService(104, new(CancelTaskService))
	RegisterUserService(105, new(RandomRevengeTaskService))

	RegisterUserService(106, new(RandomDialService))
	RegisterUserService(107, new(MultipleDialRewardService))
	//RegisterUserService(103, new(AcceptTaskService))
	//RegisterUserService(104, new(SubmitTaskService))t

	RegisterUserService(150, new(RankInfoService))

	RegisterUserService(160, new(ActivePrivilegeService))



	RegisterUserService(220, new(FollowService))            //关注用户
	RegisterUserService(221, new(UnFollowService))          //取消关注
	RegisterUserService(222, new(FollowListService))        //获取关注列表
	RegisterUserService(223, new(GetFollowerDetailService)) //获取关注用户详细信息

	//RegisterUserService(232, new(AddFriendRequestService))
	//RegisterUserService(234, new(AcceptFriendRequestService))
	RegisterUserService(236, new(SearchUserService))

	RegisterUserService(245, new(UserDetailService))
	RegisterUserService(250, new(PublicSaleService))
	RegisterUserService(251, new(CancelSaleService))
	RegisterUserService(252, new(BuySaleService))
	RegisterUserService(253, new(GetSaleInfoService))

	//RegisterUserService(255, new(GetStrangerListService))

	RegisterUserService(257, new(ReadNewsFeedService))
	RegisterUserService(258, new(NewsfeedDetailService))
	RegisterUserService(259, new(OfflineMessageService))
	//RegisterUserService(260, new(SearchItemService))
	RegisterUserService(261, new(NewsFeedListService))
	//RegisterUserService(262, new(RequestItemService))

	//RegisterUserService(264, new(AcceptItemRequestService))
	//RegisterUserService(265, new(RejectItemRequestService))
	RegisterUserService(266, new(DealListService))
	//RegisterUserService(267, new(ItemRequestOverdueService))

	RegisterUserService(281, new(GlobalMessageService))
	//RegisterUserService(290, new(PublicShareService))
	RegisterUserService(291, new(PublicWechatShareService))
	RegisterUserService(292, new(DrawWechatShareRewardService))
	RegisterUserService(293, new(GetWechatShareTimeService))
	RegisterUserService(294, new(WatchAdSuccessService))


	//朋友圈相关
	RegisterUserService(240, new(GetReceiveMomentsService))
	RegisterUserService(241, new(GetPublicMomentsService))
	RegisterUserService(243, new(PublicMomentService))

	//关注相关
	//RegisterUserService(220, new(FollowUserService))
	//RegisterUserService(221, new(UnfollowUserService))
	//RegisterUserService(222, new(GetFollowListService))

	//事件相关
	RegisterUserService(30, new(IntoEventService))
	RegisterUserService(34, new(SelectEventTargetService))
	RegisterUserService(40, new(RandomTargetService))

	RegisterUserService(41, new(OpenCardService))


	RegisterUserService(45, new(EventModuleInfoService))
	RegisterUserService(46, new(DoneEventStepService))
	RegisterUserService(65, new(GetFaithService))
	RegisterUserService(66, new(GetBelieverService))
	RegisterUserService(67, new(LootFaithService))
	RegisterUserService(68, new(AtkStarBuildingService))
	RegisterUserService(69, new(LootBelieverService))
	//RegisterUserService(501, new(GenEventService))
	//EventLocalService.RegisterHandler(31, new(PublicEventService))
	//EventLocalService.RegisterHandler(32, new(AcceptEventAssistService))
	//RegisterUserService(33, new(RemoveEventService))
	//RegisterUserService(35, new(LeaveEventService))
	//RegisterUserService(37, new(UpdateEventFieldService))
	//RegisterUserService(39, new(EventInfoService))
	//EventLocalService.RegisterHandler(47, new(SelectAreaService))
	//EventLocalService.RegisterHandler(50, new(PublicVoteService))
	//EventLocalService.RegisterHandler(51, new(AddVoteService))
	//EventLocalService.RegisterHandler(55, new(CaptureBelieverService))
	//RegisterUserService(60, new(SaveDataService))


	//星球相关

	RegisterUserService(18, new(GetStarsSelectService))
	RegisterUserService(19, new(SelectStarService))
	RegisterUserService(20, new(GetStarInfoService))
	RegisterUserService(21, new(UpgradeBuildingService))
	RegisterUserService(263, new(LootItemService)) //掠夺圣物
	//RegisterUserService(301, new(SearchStarInfoService))
	RegisterUserService(302, new(GetStarInfoDetailService))
	RegisterUserService(303, new(GetStarShieldService))     //获取防护罩信息

	//RegisterUserService(311, new(GetStarRecordService))
	//RegisterUserService(312, new(SetStarRecordService))
	//RegisterUserService(313, new(DelStarRecordService))
	//RegisterUserService(314, new(MoveStarRecordService))
	//RegisterUserService(315, new(ReplaceStarRecordService))
	//RegisterUserService(320, new(OccupyStarService))
	//RegisterUserService(321, new(NextStarService))

	RegisterUserService(321, new(GetHelpRepairInfoService))
	RegisterUserService(325, new(StarSettleService)) //星球结算

	RegisterUserService(326, new(StarFlagInfoService)) // 星球标识
	RegisterUserService(327, new(UpdateStarFlagService)) // 更新星球标识

	RegisterUserService(339, new(CancelUpgradeStarBuildService))
	RegisterUserService(340, new(CancelRepairStarBuildService))
	RegisterUserService(341, new(AccRepairStarBuildingService))  //加速修理建筑
	RegisterUserService(342, new(RepairStarBuildingService))     //修理建筑
	RegisterUserService(350, new(UpgradeStarBuildingEndService)) //请求结束建筑更新
	RegisterUserService(351, new(AccUpgradeStarBuildingService)) //加速升级建筑

	RegisterUserService(355, new(BuildingFaithInfoService))    //查看建筑信仰值
	RegisterUserService(356, new(ReceiveBuildingFaithService)) //领取建筑信仰值
	RegisterUserService(357, new(RepairStarBuildEndService))   //修理建筑结束

	//RegisterUserService(360, new(TakeInItemBuildingService))  //物品放入建筑
	//RegisterUserService(361, new(TakeoutItemBuildingService)) //物品从建筑取出
	//RegisterUserService(362, new(ResetBuildingGrooveService))		 //重铸槽
	//RegisterUserService(363, new(ActiveBuildingGroupService))	// 激活图鉴
	//RegisterUserService(364, new(AccBuildingGrooveEffectService))//加速槽的生效结果
	//RegisterUserService(365, new(UpdateGrooveEffectService)) //槽生效触发

	RegisterUserService(369, new(DrawCivilizationRewardService)) //领取文明度奖励

	RegisterUserService(371, new(StarStatisticsService))
	RegisterUserService(380, new(StarHistoryService))

	RegisterUserService(400, new(UpgradeBelieverService))  //升级信徒
	RegisterUserService(401, new(BelieverFlagInfoService)) //信徒标识信息
	RegisterUserService(402, new(AutoAddBelieverService))  //自动增加信徒
	RegisterUserService(403, new(UpdateBelieverInfoService)) //客户端同步信徒信息

	RegisterUserService(450, new(GetItemService)) //获取物品

	RegisterUserService(451, new(GetItemGroupService))        //获取图鉴信息
	RegisterUserService(452, new(DrawItemGroupRewardService)) //领取图鉴奖励
	RegisterUserService(453, new(ActiveGroupService))         //图鉴开启新物品
	//RegisterUserService(454, new(ActiveGroupItemService))//图鉴开启新物品
	//RegisterUserService(455, new(AddItemService))//新增物品
	//RegisterUserService(456, new(TempItemService))  //查看临时背包物品
	//RegisterUserService(457, new(TakeinBagService)) //临时背包放入主背包
	//RegisterUserService(470, new(GoodsInfoService))

	RegisterUserService(471, new(BuyGoodsService))
	RegisterUserService(472, new(PublicGoodsService))
	RegisterUserService(473, new(CancelGoodsService))

	RegisterUserService(474, new(GetHelpItemHistoryService))
	RegisterUserService(475, new(PublicItemHelpService))
	//RegisterUserService(476, new(CancelItemHelpService))
	RegisterUserService(477, new(DrawItemHelpService))
	RegisterUserService(478, new(GetItemHelpService))
	RegisterUserService(479, new(LootHelpItemService))
	RegisterUserService(480, new(HelpItemService))

	RegisterUserService(481, new(RefreshMallItemService))		//刷新商城圣物
	RegisterUserService(482, new(BuyMallItemService))			//购买商城圣物
	RegisterUserService(483, new(AutoRefreshMallItemService))	//自动刷新商城圣物
	RegisterUserService(484, new(GetMallItemService))			//获取商城圣物

	RegisterUserService(485, new(GetActiveGroupService))

	//-----------------------persist message service---------------------------
	//RegisterUserService(502, new(AssistEventRequestService))
	//RegisterUserService(503, new(TriggerTaskService))
	//RegisterUserService(504, new(DeleteEventAssistService))
	RegisterUserService(505, new(EventDoneService))
	RegisterUserService(510, new(DrawCorrectCivilRewardService))
	RegisterUserService(520, new(BuildingResetService))
	//RegisterUserService(521, new(TakeinGaypointService)) //拍卖物品获取友情点
	//RegisterUserService(522, new(AddStrangerService))    //新增陌生人
	RegisterUserService(523, new(AddNewsFeedService)) //新增动态消息
	//RegisterUserService(524, new(AddShopItemService)) //新增道具

	//RegisterUserService(550, new(AddLogService))    //新增动态消息

	//-------------------------center--------------------------------
	RegisterUserService(700,  new(GetOnNoticesService))


	//-------------------------压测-----------------------------------
	RegisterUserService(666, new(AddAttachService))
	RegisterUserService(531, new(RemoveSaleService))
	RegisterUserService(667, new(SetBuildsService))
	RegisterUserService(668, new(SetBelieversService))

	UserService.SetCallbackFilter(callbackFilter)
}

func callbackFilter(seq int32) bool {
	return seq >= 502 && seq <= 524
}

//注册用户消息服务
func RegisterUserService(seq int32, service IUserService) {
	UserService.RegisterHandler(seq, &NetworkServiceProxy{proxy: service})
}

//服务代理
type NetworkServiceProxy struct {
	proxy IUserService
}

func (this NetworkServiceProxy) Request(request *protocol.C2GS, response *protocol.GS2C, network service.IMessageChannel) {
	userSession, ok := network.(*user.Session)
	if ok {
		this.proxy.Request(request, response, userSession)
		if request.GetHeartBeat() == nil {
			userSession.SetDirty()
		}
	}
}

type IUserService interface {
	Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session)
}

//登录账号服务器请求
//type AddLogService struct {
//}
//
//func (service *AddLogService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetAddLog()
//	log.Debug("add log %v", message.OrderRecord)
//	message.OrderRecord.Time = time.Now()
//	response.AddLogRet = message
//}

//建筑星球建筑
type UpgradeBuildingService struct {
}

func (service *UpgradeBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	buildRequest := request.GetBuildStarBuilding()
	//是否引导操作
	//buildRequest.Guide = user.IsUpgradeGuiding()
	buildRequest.Faith = user.GetFaith()

	buildingType := buildRequest.GetBuildingType()

	////消耗信仰值
	//user.EnsureFaith(config.UpgradeConsumption)

	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetBuildStarBuildingRet()

	user.TakeOutFaith(resp.GetCost(), constant.OPT_TYPE_UPGRADE_BUILDING, buildingType)
	rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_EXPEND_FAITH_BUILD, float64(resp.GetCost()), 0)

	//user.SetDirty()
	//获取神力值
	//user.TakeInExp(config.PowerAcquired)

	if resp.GetDone() {
		resp.ItemID = user.DealUpgradeBuilding1(0, resp.GetPowerLimit(), 0, resp.GetBuilidng().GetType(), resp.GetBuilidng().GetLevel(), true)
	} else {
		//user.WriteMsg(BuildRoleSocialPush(user))
		AddRoleSocialPush(user, response)
	}
	response.BuildStarBuildingRet = resp
}

//登录服务处理类
type HeartbeatService struct {
}

func (service *HeartbeatService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.HeartBeatRet = &protocol.HeartBeatRet{}
}

type ServerTimeService struct {
}

func (service *ServerTimeService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.GetServerTimeRet = time.Now().Unix()
}

type GenOrderService struct {
}

func (service *GenOrderService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetGenOrder()
	message.Uid = user.GetID()
	resp := rpc.PassportServiceProxy.HandleMessage(request)
	response.GenOrderRet = resp.GenOrderRet
}

type GetAvatarService struct {
}

func (service *GetAvatarService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetGetAvatar()

	result := []string{}
	for _, uid := range message.GetUid() {
		result = append(result, cache.UserCache.GetUserAvatar(uid))
	}
	response.GetAvatarRet = &protocol.GetAvatarRet{
		Uid:    message.GetUid(),
		Avatar: result,
	}
}

type ChangeDescService struct {
}

func (service *ChangeDescService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	newDesc := request.GetChangeDesc()
	user.SetDesc(newDesc)
	//user.SetDirty()
	response.ChangeDescRet = true
}

type GetRoleInfoService struct {
}

func (service *GetRoleInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	subscribe := cache.UserCache.IsUserSubscribe(user.GetID())
	user.SetSubscribe(subscribe)
	response.GetRoleInfoRet = BuildRoleSocialPush(user).RoleInfoPush
}



type DisplayInfoService struct {
}

func (service *DisplayInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetRoleDisplayInfo()
	response.RoleDisplayInfoRet = &protocol.RoleDisplayInfoRet{
		Id: user.GetDisplay(message.GetMin(), message.GetMax()),
	}
}

type UpdateDisplayService struct {
}

func (service *UpdateDisplayService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetUpdateDisplay()
	user.UpdateDisplay(message.GetId())
	//user.SetDirty()
	response.UpdateDisplayRet = &protocol.UpdateDisplayRet{
		Result: true,
	}
}

//获取角色星球信息
type BuyShopItemService struct {
}

func (service *BuyShopItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetBuyShopItem()
	shopBase := conf.GetShopBase(message.GetId())
	if shopBase == nil {
		exception.GameException(exception.SHOP_BUY_ERROR)
		return
	}
	if shopBase.MoneyType == 0x02 { //需要靠钻石购买
		if !user.TakeOutDiamond(int32(shopBase.Value), constant.OPT_TYPE_BUY_SHOP, shopBase.ID) { //消耗钻石
			exception.GameException(exception.DIAMOND_NOT_ENOUGH)
			return
		}
		//user.SetDirty()
	}

	if shopBase.Type == 0x01 { //表示购买钻书
		//user.TakeInDiamond(shopBase.Amount)
		exception.GameException(exception.USER_NOT_AUTH)
	} else if shopBase.Type == 0x03 { //表示购买法力值
		user.TakeInPower(shopBase.Amount, false, constant.OPT_TYPE_BUY_SHOP)
		//user.SetDirty()
		//user.WriteMsg(BuildRoleSocialPush(user))
		AddRoleSocialPush(user, response)
	}
	response.BuyShopItemRet = &protocol.BuyShopItemRet{}
}

//发送好友请求
//type AddFriendRequestService struct {
//}
//
//func (service *AddFriendRequestService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetAddFriendRequest()
//	userId := message.GetId()
//	requestId := message.GetRequestID()
//
//	result := rpc.CommunityServiceProxy.HandleMessage(request)
//	response.AddFriendRequestRet = result.GetAddFriendRequestRet()
//
//	//请求发送成功需要推送
//	if (response.GetAddFriendRequestRet().GetResult()) {
//		//构建好友申请推送
//		pushMessage := &protocol.GS2C{
//			Sequence: []int32{1022},
//			AddFriendRequestPush: &protocol.FriendRequestInfo{
//				Id:       userId,
//				Nickname: user.GetNickName(),
//				AddTime:  time.Now().Unix(),
//			},}
//		rpc.UserServiceProxy.Push(requestId, cache.UserCache.GetUserNode(requestId), pushMessage)
//	}
//}
//
//type AcceptFriendRequestService struct {
//}
//
//func (service *AcceptFriendRequestService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetAcceptFriendRequest()
//	message.Id = user.UserID)
//	result := rpc.CommunityServiceProxy.HandleMessage(request)
//
//	//推送好友请求的玩家
//	friendPush := &protocol.GS2C{
//		Sequence:      []int32{1021},
//		AddFriendPush: BuildFriendInfo(message.GetId()),
//	}
//	rpc.UserServiceProxy.Push(message.GetRequestID(), cache.UserCache.GetUserNode(message.GetRequestID()), friendPush)
//
//	response.AcceptFriendRequestRet = result.GetAcceptFriendRequestRet()
//}

//获取收到的朋友圈消息
type GetReceiveMomentsService struct {
}

func (service *GetReceiveMomentsService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	request.GetReceiveMoments.Uid = user.GetID()

	if request.GetReceiveMoments.BeforeTime <= 0 {
		request.GetReceiveMoments.BeforeTime = time.Now().Unix()
	}
	result := rpc.CommunityServiceProxy.HandleMessage(request)
	response.GetReceiveMomentsRet = result.GetReceiveMomentsRet
}

//获取指定用户发布的朋友圈消息
type GetPublicMomentsService struct {
}

func (service *GetPublicMomentsService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	request.GetPublicMoments.Uid = user.GetID()
	if request.GetPublicMoments.BeforeTime <= 0 {
		request.GetPublicMoments.BeforeTime = time.Now().Unix()
	}
	//转发到朋友圈服务处理
	result := rpc.CommunityServiceProxy.HandleMessage(request)
	response.GetPublicMomentsRet = result.GetPublicMomentsRet
}

type PublicMomentService struct {
}

func (service *PublicMomentService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	request.PublicMoment.Uid = user.GetID()
	//转发到朋友圈服务处理
	result := rpc.CommunityServiceProxy.HandleMessage(request)
	response.PublicMomentRet = result.PublicMomentRet
}

////关注用户
//type FollowUserService struct {
//}
//
//func (service *FollowUserService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetFollowUser()
//	uid := message.GetUid()
//
//	if user.ContainFollow(uid) {
//		exception.GameException(exception.REPEATE_FOLLOW)
//	}
//	if (user.GetID() == uid) {
//		exception.GameException(exception.CANNOT_FOLLOW_ME)
//	}
//	role := user.FindRoleByID(uid)
//	user.NewFollow(role)
//	UserRPCPush(uid, BuildFollowPush(user.GetID()))
//	//response.FollowUserRet = &protocol.FollowUserRet{Result:true)}
//}
//
////取消关注
//type UnfollowUserService struct {
//}
//
//func (service *UnfollowUserService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetUnfollowUser()
//	uid := message.GetUid()
//
//	if user.ContainFollow(uid) {
//
//	} else { //
//		exception.GameException(exception.UNFOLLOW)
//	}
//}
//
////获取关注列表
//type GetFollowListService struct {
//}
//
//func (service *GetFollowListService) Request(request *protocol.C2GS, response *protocol.GS2C, useri *user.Session) {
//	response.GetFollowListRet = BuildFollowList()
//}

type RandomEventService struct {
}

func (service *RandomEventService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	////if user.HaveTask() {
	////	exception.GameException(exception.TASK_ALREAD_HAVE)
	////}
	//user.CleanTask()
	//user.DeleteEvent()
	////user.EnsurePower(conf.DATA.CostAstrola)
	//rewardTypeFilter := request.GetRandomEventTask().GetFilterType()
	//
	//task := RandomEvent(user, false, rewardTypeFilter)
	//
	////user.TakeOutPower(conf.DATA.CostAstrola, constant.OPT_TYPE_RANDOM_EVENT)
	////user.SetDirty()
	////user.WriteMsg(BuildRoleSocialPush(user))
	//response.RandomEventTaskRet = &protocol.RandomEventTaskRet{
	//	Task:     BuildTask(task),
	//	DecPower: conf.DATA.CostAstrola,
	//}
}

//func RandomEvent(user *user.Session, attackEventFilter bool, rewardFilterType int32) *db.DBRoleEventTask {
//	believerCount := cache.StarCache.GetBelieverCount(user.GetStarId())
//	buildingLevel := cache.StarCache.GetBuildingAllLevel(user.GetStarId())
//	eventBase := conf.RandomAstrola(attackEventFilter, rewardFilterType, believerCount, user.GetStarType(), buildingLevel, user.GetWeightAppendMapping())
//	if eventBase == 0 {
//		exception.GameException(exception.TASK_BASE_NOTFOUND)
//	}
//	taskData := conf.GetTriggerTaskData(eventBase)
//	if taskData == nil {
//		exception.GameException(exception.TASK_TRIGGER_NOTFOUND)
//	}
//	eventSession := user.GenEvent(eventBase, user.GetID(), user.GetNickName(), false, user)
//	rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_TYPE_EVENT, 1, eventBase)
//	return user.NewTask(taskData.ID, eventSession.ID)
//}

func GetEvent(user *user.Session, eventBase int32, source int32) *db.DBRoleEventTask {
	if eventBase == 0 {
		exception.GameException(exception.TASK_BASE_NOTFOUND)
	}
	taskData := conf.GetTriggerTaskData(eventBase)
	if taskData == nil {
		exception.GameException(exception.TASK_TRIGGER_NOTFOUND)
	}
	eventSession := user.GenEvent(eventBase, user.GetID(), user.GetNickName(), false, user)
	user.BeginEventStatistic(eventBase, source)
	return user.NewTask(taskData.ID, eventSession.ID)
}

//type AssistEventRequestListService struct {
//}
//
//func (service *AssistEventRequestListService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	flag := user.UpdateBooleanFlag(db.FLAG_MESSAGE, false)
//	user.WriteMsg(BuildFlagInfoPush(flag))
//	response.AssistEventRequestListRet = &protocol.AssistEventRequestListRet{
//		Request: BuildAssistRequests(user),
//	}
//}
//
//type RejectAssistEventService struct {
//}
//
//func (service *RejectAssistEventService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetRejectAssistEvent()
//	user.DeleteAssist(message.GetEventID())
//	response.RejectAssistEventRet = &protocol.RejectAssistEventRet{
//		Result: true,
//	}
//
//}

type TaskListService struct {
}

func (service *TaskListService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.GetTaskListRet = &protocol.GetTaskListRet{
		Tasks: BuildTasks(user),
	}
}

//type UpdateTaskEndingService struct {
//}
//
//func (service *UpdateTaskEndingService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetSetTaskEnding()
//
//	task := user.GetTask(message.GetTaskId())
//	if task == nil {
//		exception.GameException(exception.TASK_NOT_FOUND)
//	}
//	task.EndingID = message.GetEndingID()
//	response.SetTaskEndingRet = &protocol.SetTaskEndingRet{
//		Result: true,
//	}
//
//}

type CancelTaskService struct {
}

func (service *CancelTaskService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	cancelID := request.GetCancelTask()
	user.CleanTask()
	user.DeleteEvent()
	//if (eventID != 0) {
		//rpc.EventServiceProxy.RemoveEvent(eventID, user.GetID())
	//	user.RoleTempManager.DeleteEvent()
	//}
	response.CancelTaskRet = cancelID
}

type RandomRevengeTaskService struct {
}

func (service *RandomRevengeTaskService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	//if user.HaveTask() {
	//	exception.GameException(exception.TASK_ALREAD_HAVE)
	//}
	message := request.GetRandomRevengeTask()
	guide := message.GetGuide()
	if guide && user.GetFlagValue(db.FLAG_GUIDE_REVENGE) > 0 {
		exception.GameException(exception.INVALID_FLAG_AUTHORITY)
	}

	newsFeed := user.GetNewsFeed(message.GetId())
	//复仇完成或者别人的复仇交互不允许在复仇
	if newsFeed == nil || newsFeed.DoneRevenge {
		exception.GameException(exception.NEWSFEED_NOT_FOUND)
	}

	if time.Now().Sub(newsFeed.Time).Seconds() > float64(conf.DATA.CountdownRevenge) {
		exception.GameException(exception.NEWSFEED_TIMEOUT)
	}

	if !constant.IsBeAttackNewsFeed(newsFeed.Type) {
		exception.GameException(exception.NEWSFEED_NOT_FOUND)
	}

	var cost int32 = 0

	if !guide {
		cost = conf.DATA.RevengeMana
		user.EnsurePower(cost)
	}

	filterEvents := user.LockEventIDs()
	eventBase := user.RandomEvent(filterEvents)
	task := GetEvent(user, eventBase, 1)
	//task := RandomEvent(user, true, 0)
	//task := RandomEvent()
	//更新为复仇任务
	task.RevengeID = newsFeed.Uid

	if cost > 0 {
		user.TakeOutPower(cost, constant.OPT_TYPE_REVENGE)
		rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_EXPEND_POWER_REVENGE, float64(cost), 0)
		//user.WriteMsg(BuildRoleSocialPush(user))
	}
	//user.UpdateEventNum(1)
	response.RandomRevengeTaskRet = &protocol.RandomRevengeTaskRet{
		Task:     BuildTask(task),
		DecPower: cost,
	}
}

type RankInfoService struct {
}

func (service *RankInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetGetRankInfo()

	rankManager := rank.GetRankManager(message.GetType())
	if rankManager == nil {
		exception.GameException(exception.UNEXPECT_RANK_TYPE)
	}

	results := rankManager.GetTopRank()
	topRank := []*protocol.Rank{}
	for index, result := range results {
		uid := character.StringToInt32(result.Member)
		topRank = append(topRank, rankManager.GetUserRankData(uid, int32(index), result.Score))
	}

	myRank := rankManager.GetUserRankData(user.GetID(), rankManager.GetCurrentRank(int64(user.GetID())), rankManager.GetCurrentScore(int64(user.GetID())))
	response.GetRankInfoRet = &protocol.GetRankInfoRet{
		TopRank: topRank,
		MyRank:  myRank,
	}

}





//type AcceptTaskService struct {
//
//}
//
//func (service *AcceptTaskService)Request(request *protocol.C2GS,response *protocol.GS2C ,user *UserContext)  {
//	message := request.GetAcceptTask()
//	user.AcceptTask(message.GetId())
//	response.AcceptTaskRet = &protocol.AcceptTaskRet{
//		ID:message.ID,
//	}
//}
//
//
//type SubmitTaskService struct {
//
//}
//
//func (service *SubmitTaskService)Request(request *protocol.C2GS,response *protocol.GS2C ,user *UserContext)  {
//	message := request.GetSubmitTask()
//	if user.SubmitTask(message.GetId()) == 1 {
//		user.TakeInFaith( 1000 )//完成任务获取1000信仰值
//		user.WriteMsg(BuildRoleSocialPush(user))
//	}
//	response.SubmitTaskRet = &protocol.SubmitTaskRet{
//		ID:message.ID,
//	}
//}



type RandomDialService struct {
}

func (service *RandomDialService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {

	user.EnsurePower(conf.DATA.CostAstrola)
	resp := &protocol.RandomDialRet{Reward:&protocol.Reward{}}
	var randomEventType int32 = 0

	//succ := user.CompareDialTimes()
	//eventID := conf.DATA.EventID
	//if eventID != 0 {
	//	succ = true
	//}
	//if succ {//星盘任务
		//eventSession := user.GenEvent(eventBase, user.GetID(), user.GetNickName(), false, user)
	//eventBase := user.RandomEvent()
	//if eventID != 0 {
	//	if conf.DATA.EVENT_FILTER_DATA[eventID] == nil {
	//		exception.GameException(exception.EVENT_FILTER_NOT_FOUND)
	//	}
	//	eventBase = conf.DATA.EVENT_FILTER_DATA[eventID].EventBase
	//}
	if !user.IsFlagUnlock(constant.STAR_FLAG_DIAL) {
		exception.GameException(exception.DIAL_IS_NOT_UNLOCK)
	}
	buildingLevel := cache.StarCache.GetBuildingExMaxLevel(user.GetStarId())
	civilLevel := cache.StarCache.GetCivilLevel(user.GetStarId())

	civilDialData := conf.DATA.DIAL_LIMIT_DATA_MAPPING[civilLevel]
	if civilDialData == nil {
		exception.GameException(exception.DIAL_LIMIT_NOT_FOUND)
	}

	var dialID int32 = 0
	if conf.DATA.DialID != 0 {
		dialID = conf.DATA.DialID
	} else {
		dialID = user.RandomDial(civilLevel, buildingLevel, civilDialData)
	}
	dialData := civilDialData[dialID]
	if dialData == nil {
		exception.GameException(exception.DIAL_LIMIT_NOT_FOUND)
	}
	dial := conf.DATA.DIAL_DATA[dialID]
	if dial == nil {
		exception.GameException(exception.DIAL_NOT_FOUND)
	}
	resp.Position = dial.Position

	if dialData.Type == constant.QUEST_ROB_FAITH || dialData.Type == constant.QUEST_ROB_BELIEVER || dialData.Type == constant.QUEST_ATT_BUILDING  {
		user.DecDialTimes()
		randomEventType = conf.GetEventBaseByTaskType(dialData.Type)
		task := GetEvent(user, randomEventType, 0)
		resp.Task = BuildTask(task)
	} else {
		resp.Reward = user.TakeInDialReward(dialData, 1)
	}
	user.TakeOutPower(conf.DATA.CostAstrola, constant.OPT_TYPE_RANDOM_DIAL)
	rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_EXPEND_POWER_DIAL, float64(conf.DATA.CostAstrola), 0)
	resp.DecPower = conf.DATA.CostAstrola
	resp.Multiple = user.GetMultipleReward(dialID)

	response.RandomDialRet = resp

	//星际星盘数据
	lpc.StatisticsHandler.AddStatisticData(&model.StatisticDial{
		Uid:user.GetID(),
		DialID:dialID,
		Faith:resp.GetReward().GetFaith(),
		Power:resp.GetReward().GetPower(),
		GayPoint:resp.GetReward().GetGayPoint(),
		Event:randomEventType,
		Believer:getBelieverCount(resp.GetReward().GetBeliever()),
	})
}

type MultipleDialRewardService struct {

}

func (this *MultipleDialRewardService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetMultipleDialReward()
	resp := &protocol.MultipleDialRewardRet{Reward:&protocol.Reward{}}
	dial := conf.DATA.DIAL_DATA[user.GetMultipleDialID()]
	if dial == nil {
		exception.GameException(exception.DIAL_NOT_FOUND)
	}
	civilLevel := cache.StarCache.GetCivilLevel(user.GetStarId())
	dialData := conf.DATA.DIAL_LIMIT_DATA_MAPPING[civilLevel][user.GetMultipleDialID()]
	if dialData == nil {
		exception.GameException(exception.DIAL_LIMIT_NOT_FOUND)
	}
	if dial != nil {
		var multipleRatio int32 = 0
		switch message.GetType() {
		case constant.MULTIPLE_DIAL_TYPE_SHARE:
			multipleRatio = dial.ShareMultiple
		case constant.MULTIPLE_DIAL_TYPE_AD:
			multipleRatio = dial.AdMultiple
		}
		resp.Reward = user.TakeInDialReward(dialData, multipleRatio)
	}
	response.MultipleDialRewardRet = resp
}

func getBelieverCount(info []*protocol.BelieverInfo) int32 {
	if info == nil {
		return 0
	}
	var count int32 = 0
	for _, believer := range info {
		count += believer.GetNum()
	}
	return count
}

type RoleFlagInfoService struct {
}

func (service *RoleFlagInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.RoleFlagInfoRet = BuildFlagInfoResponse(user)
}

type UpdateRoleFlagService struct {
}

func (service *UpdateRoleFlagService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端通知服务器 更新红点状态
	message := request.GetUpdateFlag()
	flagID := db.ROLE_FLAG(message.GetFlag())
	//客户端没有权限修改这个标识
	//if db.FLAG_FREE_BUILDING_GROOVE == flagID {
	//	exception.GameException(exception.INVALID_FLAG_AUTHORITY)
	//}

	if db.FLAG_GUIDE == flagID || db.FLAG_GUIDE_REVENGE == flagID /*|| db.FLAG_WATCH_AD_GET_POWER == flagID */{
		//引导标识不能回退
		if message.GetValue() < user.GetFlagValue(flagID) {
			exception.GameException(exception.INVALID_FLAG_AUTHORITY)
		}
	}
	user.UpdateFlag(flagID, message.GetValue()) //更新状态

	if flagID == db.FLAG_GUIDE && message.GetValue() == constant.MAX_GUIDE_STEP {
		user.SetGuideTime(time.Now())
		//统计完成新手引导的时间
		//lpc.LogServiceProxy.AddGuideRecord(user.GetID())
	}

	response.UpdateFlagRet = &protocol.UpdateFlagRet{ //消息包回复回去
		Result: true,
	}
}

type UpdatePowerService struct {
}

func (service *UpdatePowerService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端通知服务器 更新红点状态
	serverTime := time.Now()
	user.UpdatePower(serverTime)
	response.UpdatePowerRet = &protocol.UpdatePowerRet{
		Power:           user.GetPower(),
		UpdateTimestamp: user.GetLastPowerTime().Unix(),
		ServerTimestamp: serverTime.Unix(),
	}

}

//============================> 20170531 wjl 星球模块
//type SearchStarInfoService struct {
//}
//
//func (service *SearchStarInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求探索星球
//	message := request.GetSearchStarInfo()  //获取消息包
//	message.Uid = user.GetID()) //设置用户ID
//
//	if user.GetPower() < conf.GetGameBase().CostSearchStarMana { //法力值不够
//		exception.GameException(exception.STAR_SEARCH_FAILED)
//	}
//	resp := rpc.StarServiceProxy.SearchStarInfo(message.GetUid())
//
//	response.SearchStarInfoRet = &protocol.SearchStarInfoRet{
//		Star: resp.Star,
//	}
//	user.TakeOutPower(conf.GetGameBase().CostSearchStarMana, constant.OPT_TYPE_SEARCH_STAR) //消耗法力值
//	user.WriteMsg(BuildRoleSocialPush(user))
//}

//type SelectAreaService struct {
//	//请求进攻某个星球
//}
//
//func (service *SelectAreaService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求进攻某个建筑物
//	resp := rpc.EventServiceProxy.SelectEventArea(request)
//	itemID := user.RandomItem(user)
//	cache.UserCache.UpdateATKHistory(request.GetSelectArea().GetDestUid(), 1)
//	if (resp.GetFaith() != 0) {
//		user.TakeInFaith(resp.GetFaith())
//	}
//	resp.ItemID = itemID)
//	response.SelectAreaRet = resp
//}

//type AddGayPointService struct {
//}
//
//func (service *AddGayPointService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求进攻某个建筑物
//	user.TakeInGayPoint(request.GetAddGayPoint())
//	response.AddGayPointRet = true)
//}

//type FullSearchService struct {
//	//请求进攻某个星球
//}
//
//func (service *FullSearchService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求进攻某个建筑物
//	user.UpdateFlag(db.FLAG_SEARCH_COUNT, constant.MAX_SEARCH_LIMIT)
//	response.FullSearchRet = true)
//}


func isRobot(id int32) bool {
	return id <= 0
}



func getBelieverLevel(believerID string) int32 {
	levelStr := believerID[3:4]
	return character.StringToInt32(levelStr)
}

type GetAllMailService struct {
}

func (service *GetAllMailService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求进攻某个建筑物
	request.GetGetAllMail().Uid = user.GetID()
	resp := rpc.MailServiceProxy.HandleMessage(request)
	response.GetAllMailRet = resp.GetGetAllMailRet()

}

type DrawMailService struct {
}

func (service *DrawMailService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求进攻某个建筑物
	request.GetDrawMail().Uid = user.GetID()
	resp := rpc.MailServiceProxy.HandleMessage(request).GetDrawMailRet()

	if resp == nil || resp.GetMail() == nil {
		exception.GameException(exception.MAIL_NOT_FOUND)
	}

	attach := resp.GetMail().GetAttach()
	if attach != nil {
		if attach.GetPower() > 0 {
			user.TakeInPower(attach.GetPower(), false, constant.OPT_TYPE_MAIL)
			//user.SetDirty()
		}
		if attach.GetFaith() > 0 {
			user.TakeInFaith(attach.GetFaith(), constant.OPT_TYPE_MAIL, 0)
			//user.SetDirty()
		}
		if attach.GetGayPoint() > 0 {
			user.TakeInGayPoint(attach.GetGayPoint(), constant.OPT_TYPE_MAIL, 0)
			//user.SetDirty()
		}
		if attach.GetDiamond() > 0 {
			user.TakeInDiamond(attach.GetDiamond(), constant.OPT_TYPE_MAIL, 0)
			//user.SetDirty()
		}
		if attach.GetItem() != nil {
			for _, item := range attach.GetItem() {
				user.TakeInItem(item.GetId(), item.GetNum(), constant.OPT_TYPE_MAIL, 0)
				//user.SetDirty()
			}
		}
		if attach.GetBeliever() != nil || len(attach.GetBeliever()) > 0 {
			rpc.StarServiceProxy.UpdateBelieverInfo(user.GetID(), attach.GetBeliever(), rpc.OP_BELIEVER_ADD, false)
		}
	}
	response.DrawMailRet = resp
}

type RemoveMailService struct {
}

func (service *RemoveMailService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) { //客户端向服务器请求进攻某个建筑物
	request.GetRemoveMail().Uid = user.GetID()
	resp := rpc.MailServiceProxy.HandleMessage(request)
	response.RemoveMailRet = resp.GetRemoveMailRet()
}

type GuideTaskService struct {
}

func (service *GuideTaskService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	if user.GetFlagValue(db.FLAG_GUIDE) != constant.GUIDE_TASK {
		exception.GameException(exception.INVALID_FLAG_AUTHORITY)
	}

	taskBase := conf.GetEventTask(constant.EVENT_ID_LOOT_BELIEVER)
	if taskBase == nil {
		exception.GameException(exception.TASK_BASE_NOTFOUND)
	}

	eventSession := user.GenEvent(constant.EVENT_ID_LOOT_BELIEVER, user.GetID(), user.GetNickName(), true, user)
	task := user.NewTask(taskBase.ID, eventSession.ID)
	//user.SetDirty()

	response.GuideTaskRet = BuildTask(task)
}

type GuideBuildingFaithService struct {
}

func (service *GuideBuildingFaithService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	if user.GetFlagValue(db.FLAG_GUIDE) != constant.GUIDE_BUILDING_FAITH {
		exception.GameException(exception.INVALID_FLAG_AUTHORITY)
	}
	buildingConf := conf.GetBuildingConf(user.GetStarType(), request.GetGuideBuildingFaith(), 1)
	if buildingConf == nil {
		exception.GameException(exception.STAR_BUILDING_NOTFOUND)
	}
	user.TakeInFaith(buildingConf.FaithLimit, constant.OPT_TYPE_DRAW_BUILDING_FAITH, buildingConf.ID)

	flag := user.UpdateFlag(db.FLAG_GUIDE, constant.GUIDE_BUILDING_FAITH+1)
	//user.SetDirty()
	//user.WriteMsg(BuildFlagInfoPush(flag))
	AddFlagInfoPush(flag, response)
	response.GuideBuildingFaithRet = buildingConf.FaithLimit
}

type GuideRevengeService struct {
}

func (service *GuideRevengeService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	if user.GetFlagValue(db.FLAG_GUIDE_REVENGE) > 0 {
		exception.GameException(exception.INVALID_FLAG_AUTHORITY)
	}
	eventBase := conf.RandomAttackAstrola(user.GetStarType())

	message := &protocol.RandomTarget{
		Uid:         user.GetID(),
		EventType:   eventBase,
		Num:         1,
		RobotFilter: true,
	}
	targets := rpc.SearchServiceProxy.RandomTarget(message).GetTargets()
	response.GuideRevengeRet = &protocol.GuideRevengeRet{}
	if len(targets) > 0 {
		target := targets[0]
		if rand.Float64() <= 0.5 {
			newsFeed := BuildNewsFeed(target.GetId(), constant.NEWSFEED_TYPE_BE_LOOT_BELIEVER, 0, 0, 0)
			user.AddNewsFeed(newsFeed)
			//user.SetDirty()
			response.GuideRevengeRet.NewsFeed = newsFeed
		} else {
			newsFeed := BuildNewsFeed(target.GetId(), constant.NEWSFEED_TYPE_BE_LOOT_FAITH, 0, 0, 0)
			user.AddNewsFeed(newsFeed)
			//user.SetDirty()
			response.GuideRevengeRet.NewsFeed = newsFeed
		}
	}
}

type BelieverFlagInfoService struct {
	//信徒标识
}

func (service *BelieverFlagInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.BelieverFlagInfoRet = BuildBelieverFlagInfoResponse(user)
}

type AutoAddBelieverService struct {
}

func (service *AutoAddBelieverService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	request.AutoAddBeliever.Uid = user.GetID()
	result := rpc.StarServiceProxy.Call(user.GetID(), request)
	response.AutoAddBelieverRet = result.GetAutoAddBelieverRet()
}

type UpdateBelieverInfoService struct {

}

func (service *UpdateBelieverInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	result := rpc.StarServiceProxy.Call(user.GetID(),request)
	response.UpdateBelieverInfoRet = result.GetUpdateBelieverInfoRet()
}

type ActiveGroupService struct {
	//尝试图鉴组合
}

func (service *ActiveGroupService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetActiveGroup()
	if !user.CanTakeoutItems(message.ItemID) {
		exception.GameException(exception.ITEM_NOT_ENOUGH)
	}
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetActiveGroupRet()
	if resp.GetResult() {
		user.SetActiveGroupID(message.GroupID)
		//user.SetDirty()
	}
	user.TakeOutItems(message.GetItemID(), constant.OPT_TYPE_ACTIVE_GROUP, message.GetGroupID())
	response.ActiveGroupRet = resp
}

//type ActiveBuildingGroupService struct {
//}
//
//func (service *ActiveBuildingGroupService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//
//}

//type AccBuildingGrooveEffectService struct {
//}
//
//func (service *AccBuildingGrooveEffectService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	resp := rpc.StarServiceProxy.PersistCall(user.GetID(), request)
//	response.AccBuildingGrooveEffectRet = resp.GetAccBuildingGrooveEffectRet()
//}

//type UpdateGrooveEffectService struct {
//}
//
//func (service *UpdateGrooveEffectService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	resp := rpc.StarServiceProxy.PersistCall(user.GetID(), request)
//	response.UpdateGrooveEffectRet = resp.GetUpdateGrooveEffectRet()
//}

type DrawCivilizationRewardService struct {
}

func (service *DrawCivilizationRewardService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	//request.GetDrawCivilizationReward().Uid = user.GetID()
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetDrawCivilizationRewardRet()

	if resp.GetReward() > 0 {
		user.TakeInPower(resp.GetReward(), false, constant.OPT_TYPE_DRAW_CIVIL)
	}
	if resp.GetFaith() > 0 {
		user.TakeInFaith(resp.GetFaith(), constant.OPT_TYPE_DRAW_CIVIL, resp.GetDrawLevel())
	}
	if resp.GetDiamond() > 0 {
		user.TakeInDiamond(resp.GetDiamond(), constant.OPT_TYPE_DRAW_CIVIL, resp.GetDrawLevel())
	}
	if resp.GetGayPiont() > 0 {
		user.TakeInGayPoint(resp.GetGayPiont(), constant.OPT_TYPE_DRAW_CIVIL, resp.GetDrawLevel())
	}
	//user.SetDirty()
	response.DrawCivilizationRewardRet = resp
}

//type TakeInItemBuildingService struct {
//}

//func (service *TakeInItemBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetPutItemToBuilding()
//	//TODO 涉及到从别的建筑中替换,需要校验是否在其他建筑存在
//	//物品不够
//	//if !user.CanTakeoutItems(message.GetItemID()) {
//	//	exception.GameException(exception.ITEM_NOT_ENOUGH)
//	//}
//	takeoutItems := message.GetTakeoutItem()
//	putinItems := message.GetItemID()
//
//	//放入背包的圣物
//	bagTakeinItem := character.GetArrayDeff(takeoutItems, putinItems)
//	//从背包中需要扣除的圣物
//	bagTakeoutItem := character.GetArrayDeff(putinItems, takeoutItems)
//
//	if !user.CanTakeoutItems(bagTakeoutItem) {
//		exception.GameException(exception.ITEM_NOT_ENOUGH)
//	}
//	result := rpc.StarServiceProxy.PersistCall(user.GetID(), request).GetPutItemToBuildingRet()
//	user.TakeInItems(bagTakeinItem, constant.OPT_TYPE_GROOVE, message.GetBuildingType())
//	user.TakeOutItems(bagTakeoutItem, constant.OPT_TYPE_GROOVE, message.GetBuildingType())
//	result.BagTakeinItems = bagTakeinItem
//	result.BagTakeoutItems = bagTakeoutItem
//	response.PutItemToBuildingRet = result
//}

//type TakeoutItemBuildingService struct {
//}
//
//func (service *TakeoutItemBuildingService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	result := rpc.StarServiceProxy.PersistCall(user.GetID(), request)
//	takeoutItem := result.GetRemoveItemFromBuildingRet().GetItemID()
//	user.TakeInItem(takeoutItem, 1, constant.OPT_TYPE_GROOVE, request.GetRemoveItemFromBuilding().GetBuildingType())
//	response.RemoveItemFromBuildingRet = result.GetRemoveItemFromBuildingRet()
//}

//type ResetBuildingGrooveService struct{//信徒标识
//
//}
//
//func ( service *ResetBuildingGrooveService )Request(request *protocol.C2GS, response *protocol.GS2C,  user *user.Session) {
//	message := request.GetResetBuildingGroove()
//	message.Uid = user.GetID())
//
//	cost := conf.DATA.RECAST_DATA[len(message.GetLockGroove())]
//	if (cost != nil) {
//		if (cost.MoneyType == constant.MONEY_TYPE_GOLD) {
//			user.AssertFaith(cost.Cost)
//		}
//		if (cost.MoneyType == constant.MONEY_TYPE_DIAMOND) {
//			user.AssertDiamond(cost.Cost)
//		}
//	}
//
//	result := rpc.StarServiceProxy.HandleMessage(request)
//	returnItems := result.GetResetBuildingGrooveRet().GetItemID()
//	user.TakeInItems(returnItems)
//	if (cost != nil) {
//		if (cost.MoneyType == constant.MONEY_TYPE_GOLD) {
//			user.TakeOutFaith(cost.Cost)
//		}
//		if (cost.MoneyType == constant.MONEY_TYPE_DIAMOND) {
//			user.TakeOutDiamond(cost.Cost)
//		}
//	}
//	user.WriteMsg(BuildRoleSocialPush(user))
//	response.ResetBuildingGrooveRet = result.GetResetBuildingGrooveRet()
//}

type GetItemService struct {
	//获取物品
}

func (service *GetItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.GetBagItemRet = &protocol.GetBagItemRet{
		ItemList: user.GetProtocolItems(),
	}
}

type GetItemGroupService struct {
	//获取圣物组合
}

func (service *GetItemGroupService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetItemGroupRet
	response.GetItemGroupRet = resp
}

type DrawItemGroupRewardService struct {
	//领取图鉴奖励
}

func (service *DrawItemGroupRewardService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetGetGroupReward()
	groupBase := conf.DATA.ITEM_GROUP[message.GetGroupID()]
	if (groupBase == nil) {
		exception.GameException(exception.ITEM_GOURP_BASE_NOT_FOUND)
	}
	resp := rpc.StarServiceProxy.Call(user.GetID(), request).GetGroupRewardRet
	user.TakeInDiamond(groupBase.Reward, constant.OPT_ITEM_GROUP_ITEM_REWARD, groupBase.ID)
	//user.WriteMsg(BuildRoleSocialPush(user))
	AddRoleSocialPush(user, response)
	response.GetGroupRewardRet = resp
}

//type ActiveGroupItemService struct{//获取物品
//
//}
//
//func ( service *ActiveGroupItemService )Request(request *protocol.C2GS, response *protocol.GS2C,  user *user.Session) {
//	message := request.GetActiveGroupItem()
//	groupBase := conf.DATA.ITEM_GROUP[message.GetGroupID()]
//	if (groupBase == nil) {
//		exception.GameException(exception.ITEM_GOURP_BASE_NOT_FOUND)
//	}
//	if (!groupBase.ContainsItem(message.GetItemID())) {
//		exception.GameException(exception.INVALID_ITEM_GOURP_ACTIVE)
//	}
//	if (!user.AddItemGroup(message.GetGroupID(), message.GetItemID())) {
//		exception.GameException(exception.INVALID_ITEM_GOURP_ACTIVE)
//	}
//	response.ActiveGroupItemRet = &protocol.ActiveGroupItemRet{
//		GroupID:message.GroupID,
//	}
//}

//type AddItemService struct{//获取物品
//
//}
//
//func ( service *AddItemService )Request(request *protocol.C2GS, response *protocol.GS2C,  user *user.Session) {
//	itemID := request.GetAddItem()
//	if user.GetFlagBooleanValue(db.FLAG_FREE_BUILDING_GROOVE) {
//		exception.GameException(exception.SERVICE_INVALID)
//	}
//	user.UpdateBooleanFlag(db.FLAG_FREE_BUILDING_GROOVE, true)
//	user.TakeInItem(itemID, 1)
//	response.AddItemRet = true)
//}

//type TempItemService struct {
//	//获取物品
//}
//
//func (service *TempItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	response.GetTempItemRet = &protocol.GetTempItemRet{
//		ItemID: user.GetTempItems(),
//	}
//}
//
//type TakeinBagService struct {
//	//获取物品
//}
//
//func (service *TakeinBagService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	//message := request.GetTakeinBag()
//	//success := user.TakeoutTempItems()
//	//if !success {
//	//	exception.GameException(exception.ITEM_NOT_FOUND)
//	//}
//
//	items := user.TakeinAllTempItem()
//	response.TakeinBagRet = &protocol.TakeinBagRet{
//		ItemID: items,
//	}
//}

type BuyGoodsService struct {
}

func (service *BuyGoodsService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetBuyGoods()
	message.Social = user.GetGayPoint()

	resp := rpc.TradeServiceProxy.HandleMessage(request).BuyGoodsRet

	user.TakeOutGayPoint(resp.GetFee(), constant.OPT_TYPE_GOODS_BUY, message.GetUid())
	user.TakeInItem(message.GetItemid(), message.GetNum(), constant.OPT_TYPE_GOODS_BUY, message.GetUid())
	//user.SetDirty()

	//user.WriteMsg(BuildRoleSocialPush(user))
	AddRoleSocialPush(user, response)
	sendMessage := BuildNewsFeedMessage(user.GetID(), constant.NEWSFEED_TYPE_BE_BUY_GOODS, message.GetItemid(),
		message.GetNum(), resp.GetFee())
	//推送友情点获取消息
	rpc.UserServiceProxy.PersistCall(message.GetUid(), sendMessage)

	response.BuyGoodsRet = resp
}

type PublicGoodsService struct {
}

func (service *PublicGoodsService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetPublicGoods()
	message.Uid = user.GetID()

	if !user.CanTakeoutItem(message.GetGoods().GetId(), message.GetGoods().GetNum()) {
		exception.GameException(exception.ITEM_NOT_ENOUGH)
	}
	resp := rpc.TradeServiceProxy.HandleMessage(request)

	user.TakeOutItem(message.GetGoods().GetId(), message.GetGoods().GetNum(), constant.OPT_TYPE_GOODS_PUBLIC, 0)
	//user.SetDirty()

	response.PublicGoodsRet = resp.PublicGoodsRet
}

type CancelGoodsService struct {
	//获取物品
}

func (service *CancelGoodsService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetCancelGoods()
	message.Uid = user.GetID()

	resp := rpc.TradeServiceProxy.HandleMessage(request).GetCancelGoodsRet()
	if resp.GetGoods() != nil {
		user.TakeInItem(resp.GetGoods().GetId(), resp.GetGoods().GetNum(), constant.OPT_TYPE_GOODS_CANCEL, 0)
		//user.SetDirty()
	}
	response.CancelGoodsRet = resp
}

//---------------------persist service-------------
//事件完成
type EventDoneService struct {
}

func (service *EventDoneService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetEventDone()
	task, isRandomItem := user.DoneEventTask(message, user.GetStarType())
	//user.SetDirty()
	if task == nil {
		return
	}

	var gayPoint int32 = 0
	addBeliever := task.RewardBeliever
	//changeFaith := task.RewardFaith
	rewardItem := task.RewardItem
	if isRandomItem {
		rewardItem, gayPoint = user.RandomItems(constant.RANDOM_TASK)
		task.RewardItem = rewardItem
	}

	if addBeliever != nil && len(addBeliever) > 0 {
		resp := rpc.StarServiceProxy.UpdateBelieverInfo(user.GetID(), addBeliever, rpc.OP_BELIEVER_ADD, true).GetUpdateBelieverRet()
		//task.RewardBeliever = resp.Believer
		//task.RewardFaith = resp.Faith
		//追加信仰
		task.RewardFaith += resp.Faith
		task.RewardBeliever = resp.Believer
		//changeFaith = task.RewardFaith
	}
	if gayPoint != 0 {
		task.RewardGayPoint = gayPoint
		user.TakeInGayPoint(gayPoint, constant.OPT_TYPE_TASK_REWARD, task.Type)
		//user.SetDirty()
	}
	if rewardItem != 0 {
		user.TakeInItem(rewardItem, 1, constant.OPT_TYPE_TASK_REWARD, task.Type)
		//user.SetDirty()
	}
	if task.RewardFaith >= 0 {
		user.TakeInFaith(task.RewardFaith, constant.OPT_TYPE_TASK_REWARD, task.Type)
		//user.SetDirty()
	}

	taskPush := BuildTaskPush(task)

	user.EndEventStatistic(task.RewardFaith)

	user.UpdateLastTask(BuildTask(task))
	//user.SetDirty()
	user.WriteMsg(taskPush)
	//AddTaskPush(task, response)

	if task.RevengeID != 0 {
		changeNewsFeed := user.DoneRevengeNewsFeed(task.RevengeID)
		//user.SetDirty()
		//推送完成复仇的消息
		if len(changeNewsFeed) > 0 {
			user.WriteMsg(&protocol.GS2C{
				Sequence:    []int32{1080},
				DoneRevenge: changeNewsFeed,
			})
		}
	}

	//如果是引导任务，需要更新引导标识
	if message.GetGuide() {
		guideFlag := user.UpdateFlag(db.FLAG_GUIDE, constant.GUIDE_TASK+1)
		//user.SetDirty()
		if guideFlag != nil {
			user.WriteMsg(BuildFlagInfoPush(guideFlag))
			//AddFlagInfoPush(guideFlag, response)
		}
	}

	//user.UpdateTargetRule(message.GetTargetID())

	//searchFlag := user.AddLimitFlag(db.FLAG_SEARCH_COUNT, 1, constant.MAX_SEARCH_LIMIT)
	//if searchFlag != nil {
	//	user.WriteMsg(BuildFlagInfoPush(searchFlag))
	//}

	//user.RandomItem(user)
}

//查找好友
type SearchUserService struct {
}

func (service *SearchUserService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {

	message := request.GetSearchUser()

	if !cache.UserCache.IsUserExist(message.GetId()) {
		exception.GameException(exception.USER_NOTFOUND)
	}

	userInfo := &protocol.UserInfo{
		Id: message.GetId(),
	}

	FillDetailInfo(userInfo)

	userInfo.Data = &protocol.UserDetailInfo{
		Desc: cache.UserCache.GetUserDesc(message.GetId()),
	}

	response.SearchUserRet = &protocol.SearchUserRet{User: userInfo}

	//nickName := message.GetNickname()
	//
	//uidMapping := cache.CommunityCache.GetAllUIDByNickname(nickName)
	//if uidMapping == nil || len(uidMapping) == 0 {
	//	response.SearchFriendRet = &protocol.SearchFriendRet{}
	//} else {
	//	friendInfo := []*protocol.UserInfo{}
	//	for uid, _ :=  range uidMapping {
	//		friendInfo = append(friendInfo, core.BuildFriendInfo1(uid, nickName))
	//	}
	//	response.SearchFriendRet = &protocol.SearchFriendRet{Friend:friendInfo}
	//}

}


type ActivePrivilegeService struct {
}

func (service *ActivePrivilegeService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	active := request.GetActivePrivilege() //是否激活特权

	//特权激活状态有变化，需要更新和推送
	if user.GetSubscribe() != active {
		if active {
			draw := user.GetFlagBooleanValue(db.FLAG_DRAW_PRIVILEGE)
			if !draw {
				user.TakeInPower(conf.DATA.PrivilegeRewardDis, false, constant.OPT_TYPE_PRIVILEGE_GIFT)
				user.UpdateBooleanFlag(db.FLAG_DRAW_PRIVILEGE, true)
			}
		}

		user.SetSubscribe(active)
		user.WriteMsg(BuildRoleSocialPush(user))
	}

}


//重置建筑
type BuildingResetService struct {
}

func (service *BuildingResetService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetBuildingReset()

	powerLimit := message.GetPowerLimit()
	addPower := powerLimit - user.GetPowerLimitBase()
	if addPower > 0 {
		user.TakeInPower(addPower, false, constant.OPT_TYPE_UPGRADE_BUILDING)
		//user.SetDirty()
	}

	starType := user.GetStarType()
	user.SetPowerLimit(powerLimit)
	user.WriteMsg(BuildRoleSocialPush(user))
	//AddRoleSocialPush(user, response)
	user.TakeInItems(message.GetItemID(), constant.OPT_TYPE_BUILDING_RESET, 0)
	user.WriteMsg(util.BuildBuildingInfoPush(user.GetID(), user.GetStarId(), starType, message.GetBuilding()))
	//user.SetDirty()
	//推送消息
	for _, building := range message.GetBuilding() {
		newsFeed := BuildNewsFeed(user.GetID(), constant.NEWSFEED_TYPE_BUILDING_DESTORY, starType, building.GetType(), building.GetLevel())
		user.AddNewsFeed(newsFeed)
		if request.GetOffline() {
			user.AddOfflineMessage(newsFeed)
		} else {
			user.WriteMsg(BuildNewsFeedPush(newsFeed))
			//AddNewsFeedPush(newsFeed, response)
		}
		//user.SetDirty()
	}
}

type AddNewsFeedService struct {
}

func (service *AddNewsFeedService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetAddNewsFeed()

	//订单过期不需要处理通知了
	if message.GetType() == constant.NEWSFEED_TYPE_BE_ACCEPT || message.GetType() == constant.NEWSFEED_TYPE_BE_REJECT {
		deal := user.GetDeal(message.GetId())
		//交易不存在 不需要处理
		if deal == nil {
			return
		}
		//拒绝或通过后需要删除订单
		user.RemoveDeal(message.GetId())
	}

	//user.AddMutualEvent(message.GetRelateID(), message.GetType(), time.Unix(message.GetTime(), 0))

	overdue := user.AddNewsFeed(message) == nil

	//通过索取，添加物品
	//if message.GetType() == constant.NEWSFEED_TYPE_BE_ACCEPT {
	//	//新增交易统计
	//	statistic := user.AddStatisticsValue(constant.STATISTIC_TYPE_SALE, 1)
	//	user.WriteMsg(BuildScorePush(statistic))
	//
	//	item := user.TakeInItem(message.GetParam1(), 1, constant.OPT_TYPE_SALE, message.GetRelateID())
	//	user.WriteMsg(util.BuildItemPush(item))
	//}

	//util.BuildItemPush()

	//if constant.IsBeAttackNewsFeed(message.GetType()) {
	//	rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_TYPE_BE_ATTACK, 1, 0)
	//}
	if message.GetType() == constant.NEWSFEED_TYPE_BE_LOOT_FAITH {
		rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_BE_LOOT_FAITH, 1 , 0)
	}
	if message.GetType() == constant.NEWSFEED_TYPE_BE_LOOT_BELIEVER {
		rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_BE_LOOT_BELIEVER, 1 , 0)
	}
	if message.GetType() == constant.NEWSFEED_TYPE_BE_ATK_BUILD {
		rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_BE_ATK_BUILDING, 1 , 0)
	}

	//收到求助，清除自动推送
	if message.GetType() == constant.NEWSFEED_TYPE_BE_HELP_ITEMHELP {
		user.CleanHelpPublicTime()
	}

	if message.GetType() == constant.NEWSFEED_TYPE_BE_BUY_GOODS {
		//新增交易统计
		statistic := user.AddStatisticsValue(constant.STATISTIC_TYPE_SALE, 1)
		user.WriteMsg(BuildScorePush(statistic))
		//AddScorePush(statistic, response)
		user.TakeInGayPoint(message.GetParam3(), constant.OPT_TYPE_GOODS_BUY, 0)
		user.WriteMsg(BuildRoleSocialPush(user))
		//AddRoleSocialPush(user, response)
	}

	if message.GetType() == constant.NEWSFEED_TYPE_BE_LOOT_FAITH {
		//decFaith := user.ForceTakeOutFaith(message.GetParam1(), constant.OPT_TYPE_EVENT_LOOT, message.GetRelateID())
		decFaith := int32(float64(user.GetFaith())*0.3)
		if decFaith > 0 {
			user.TakeOutFaith(decFaith, constant.OPT_TYPE_EVENT_LOOT, message.GetRelateID())
		}
		message.Param1 = decFaith
		user.WriteMsg(BuildRoleSocialPush(user))
		//AddRoleSocialPush(user, response)
	}

	if message.GetType() == constant.NEWSFEED_TYPE_BE_BUY_SALE {
		//新增交易统计
		statistic := user.AddStatisticsValue(constant.STATISTIC_TYPE_SALE, 1)
		user.WriteMsg(BuildScorePush(statistic))
		//AddScorePush(statistic, response)
		if message.GetParam2() > 0 {
			user.TakeInGayPoint(message.GetParam2(), constant.OPT_TYPE_SALE, 0)
			user.WriteMsg(BuildRoleSocialPush(user))
			//AddRoleSocialPush(user, response)
		}
		return
	}

	//没有关注。需要添加偶遇记录
	//if !cache.UserCache.ExistFollower(character.Int32ToString(user.GetID()), character.Int32ToString(message.GetRelateID())) {
	//	//掠夺，需要添加偶遇记录
	//	if message.GetType() == constant.NEWSFEED_TYPE_BE_LOOT_ITEM {
	//		stranger := BuildStranger(message.GetRelateID(), message.GetRelateNickname(), constant.STRANGER_TYPE_LOOT)
	//		stranger.Param = message.Param1
	//		user.AddStranger(stranger)
	//		user.WriteMsg(BuildStrangerPush(stranger))
	//	} else if message.GetType() == constant.NEWSFEED_TYPE_BE_BUY_SALE {
	//		stranger := BuildStranger(message.GetRelateID(), message.GetRelateNickname(), constant.STRANGER_TYPE_BUY_SELL)
	//		stranger.Param = message.Param1
	//		user.AddStranger(stranger)
	//		user.WriteMsg(BuildStrangerPush(stranger))
	//	}
	//}

	if overdue {
		return
	}

	if request.GetOffline() {
		user.AddOfflineMessage(message)
	} else {
		user.WriteMsg(BuildNewsFeedPush(message))
		//AddNewsFeedPush(message, response)
	}

}

type DrawCorrectCivilRewardService struct {

}

func (service *DrawCorrectCivilRewardService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetCorrectCivilReward()
	power := message.GetPower()
	if power > 0 {
		user.TakeInPower(power, false, constant.OPT_TYPE_DRAW_CIVIL)
		//同步数据给客户端
		user.WriteMsg(user.BuildRoleSocialPush())
	}
}

//-------------------------------------------社交消息----------------------
type FollowService struct {
}

func (service *FollowService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	request.Follow.Id = user.GetID()
	result := rpc.CommunityServiceProxy.HandleMessage(request)
	FillDetailInfo(result.FollowRet.Follower)
	//followId := result.FollowRet.Follower.GetId()
	//user.RemoveStranger(followId)

	//followMessage := BuildStrangerMessage(user.GetID(), user.GetNickName(), constant.STRANGER_TYPE_FOLLOW)
	//推送关注消息
	//rpc.UserServiceProxy.PersistCall(followId, cache.UserCache.GetUserNode(followId), followMessage)

	response.FollowRet = result.GetFollowRet()
}

type UnFollowService struct {
}

func (service *UnFollowService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	request.Unfollow.Id = user.GetID()
	result := rpc.CommunityServiceProxy.HandleMessage(request)
	response.UnfollowRet = result.GetUnfollowRet()
}

type FollowListService struct {
}

func (service *FollowListService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	request.GetFollowerList.Id = user.GetID()
	result := rpc.CommunityServiceProxy.HandleMessage(request)

	for _, follower := range result.GetFollowerListRet.Followers {
		FillDetailInfo(follower)
	}
	response.GetFollowerListRet = result.GetGetFollowerListRet()
}

type GetFollowerDetailService struct {
}

func (service *GetFollowerDetailService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	followerID := request.GetGetFollowerDetail()

	//var saleItem int32 = 0
	//sale := rpc.TradeServiceProxy.GetSale(followerID)
	//if sale != nil {
	//	saleItem = sale.GetItemID()
	//}
	response.GetFollowerDetailRet = &protocol.UserDetailInfo{
		Id:     followerID,
		Desc:   cache.UserCache.GetUserDesc(followerID),
		//ItemID: saleItem,
	}
}

type UserDetailService struct {
}

func (service *UserDetailService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetGetUserDetail()
	results := make([]*protocol.UserDetail, len(message.Uid))
	for index, uid := range message.Uid {
		if uid < 0{
			nickname := conf.DATA.ROBOT_MAPPPING[uid]
			if nickname == "" {
				nickname = words.RandomName()
			}
			avatar := /* conf.DATA.RobotAvatarPrefix +*/ character.Int32ToString(-uid) + ".jpg"
				results[index] = &protocol.UserDetail{
				Uid:      uid,
				Nickname: nickname,
				Avatar:   avatar,
			}
		} else {
			results[index] = &protocol.UserDetail{
				Uid:      uid,
				Nickname: cache.UserCache.GetUserNickname(uid),
				Avatar:   cache.UserCache.GetUserAvatar(uid),
			}
		}
	}
	response.GetUserDetailRet = &protocol.GetUserDetailRet{UserDetail: results}
}

//---------------商品挂售------------------
type PublicSaleService struct {
}

func (service *PublicSaleService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetPublicSale()
	itemID := message.GetItemID()

	if !user.CanTakeoutItem(itemID, 1) {
		exception.GameException(exception.ITEM_NOT_ENOUGH)
	}

	if int64(time.Now().Sub(user.GetPublicTime()).Seconds()) < conf.DATA.CountdownApplying {
		exception.GameException(exception.ITEM_CAN_NOUT_PUBLIC)
	}
	sale := rpc.TradeServiceProxy.AddSale(user.GetID(), itemID).GetSale()
	rpc.CommunityServiceProxy.PublicMoments(user.GetID(), constant.MOMENTS_TYPE_SALE, itemID)

	user.TakeOutItem(itemID, 1, constant.OPT_TYPE_SALE, 0)
	//user.SetDirty()
	//sale.Nickname = user.GetNickName())
	//sale.Avatar = user.GetAvatar())
	saleTime := user.RefreshPublicTime()
	sale.PublicTime = saleTime.Unix()

	//需要广播所有
	global.BroadcastMessage(&protocol.GlobalMessage{Sale: sale})
	//BroadcastGlobalMessage()

	//if cache.UserCache.ExistFollower(character.Int32ToString(user.GetID()), character.Int32ToString(message.GetId())) {
	//	user.WriteMsg(BuildFollowPush(message.GetId()))
	//	return
	//}

	//推送给所有粉丝
	//pushMessage := util.BuildSalePush(sale)
	//user.PushFollowings(pushMessage)

	response.PublicSaleRet = &protocol.PublicSaleRet{Result: true, PublicTime: saleTime.Unix()}
}

type CancelSaleService struct {
}

func (service *CancelSaleService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetCancelSale()
	itemID := message.GetItemID()

	rpc.TradeServiceProxy.RemoveSale(user.GetID(), itemID)
	rpc.CommunityServiceProxy.RemoveMoments(user.GetID())

	user.TakeInItem(itemID, 1, constant.OPT_TYPE_SALE, 0)
	//user.SetDirty()
	response.CancelSaleRet = &protocol.CancelSaleRet{Result: true}
}

type BuySaleService struct {
}

func (service *BuySaleService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetBuySale()
	publicID := message.GetId()
	itemID := message.GetItemID()

	if publicID == user.GetID() {
		exception.GameException(exception.CAN_NOT_BUY_SELF)
	}
	user.EnsureGayPoint(conf.DATA.FriendPointPrices)
	rpc.TradeServiceProxy.RemoveSale(publicID, itemID)
	rpc.CommunityServiceProxy.RemoveMoments(publicID)


	user.TakeOutGayPoint(conf.DATA.FriendPointPrices, constant.OPT_TYPE_SALE, publicID)
	user.TakeInItem(itemID, 1, constant.OPT_TYPE_SALE, publicID)
	statistic := user.AddStatisticsValue(constant.STATISTIC_TYPE_SALE, 1)
	//user.SetDirty()
	//user.WriteMsg(BuildScorePush(statistic))
	AddScorePush(statistic, response)
	//user.WriteMsg(BuildRoleSocialPush(user))
	AddRoleSocialPush(user, response)

	sendMessage := BuildNewsFeedMessage(user.GetID(), constant.NEWSFEED_TYPE_BE_BUY_SALE, itemID, conf.DATA.FriendPointPrices, 0)
	//推送友情点获取消息
	rpc.UserServiceProxy.PersistCall(publicID, sendMessage)
	//AddUserNewsFeed(user, BuildNewsFeed1(publicID, constant.NEWSFEED_TYPE_BUY_SALE, itemID, conf.DATA.FriendPointPrices, nil))

	//物品被卖出
	rpc.CommunityServiceProxy.PublicMoments(publicID, constant.MOMENTS_TYPE_BUY, itemID)

	response.BuySaleRet = &protocol.BuySaleRet{Result: true}

}

type GetSaleInfoService struct {
}

func (service *GetSaleInfoService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	//followers := cache.UserCache.GetFollowers(character.Int32ToString(user.GetID()))

	response.SaleInfoRet = &protocol.GetSaleInfoRet{}
	//if (len(followers) != 0) {
	//	requestIDs := []int32{}
	//	for id, _ := range followers {
	//		requestIDs = append(requestIDs, character.StringToInt32(id))
	//	}
	//
	//	//sales := rpc.CommunityServiceProxy.GetSales(requestIDs)
	//	//if (sales != nil) {
	//	//	for _, sale := range sales {
	//	//		sale.Nickname = cache.UserCache.GetUserNickname(sale.GetId()))
	//	//		sale.Avatar = cache.UserCache.GetUserAvatar(sale.GetId()))
	//	//	}
	//	//}
	//	//response.SaleInfoRet.Sales = sales
	//}
	mySale := rpc.TradeServiceProxy.GetSale(user.GetID())
	if mySale != nil {
		response.SaleInfoRet.MySale = mySale
	}
	response.SaleInfoRet.PublicTime = user.GetPublicTime().Unix()
}

//type GetStrangerListService struct {
//}
//
//func (service *GetStrangerListService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	response.GetStrangerlistRet = &protocol.GetStrangerListRet{
//		Strangers: user.GetProtocolStrangers(),
//	}
//}

type ReadNewsFeedService struct {

}

func (service *ReadNewsFeedService) Request (request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetReadNewsFeed()
	ret := user.ReadNewsFeed(message.GetId())
	response.ReadNewsFeedRet = &protocol.ReadNewsfeedRet{Result:ret}
}

type NewsfeedDetailService struct {
}

func (service *NewsfeedDetailService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetGetNewsfeedDetail()
	newsFeed := user.GetNewsFeed(message.GetId())
	if newsFeed == nil {
		exception.GameException(exception.NEWSFEED_NOT_FOUND)
	}

	response.GetNewsfeedDetailRet = &protocol.GetNewsfeedDetailRet{
		Self:  newsFeed.Self,
		Other: newsFeed.Other,
	}
}

type OfflineMessageService struct {
}

func (service *OfflineMessageService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.GetOfflineNewsFeedListRet = &protocol.GetOfflineNewsfeedListRet{
		NewsFeeds: user.GetOfflineMessage(),
	}
}

type SearchItemService struct {
}

func (service *SearchItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetSearchItem()

	request.SearchItem.Id = user.GetID()
	filter := message.GetExistItems()
	searchFlag := user.GetFlag(db.FLAG_SEARCH_COUNT)

	var cost int32 = 0
	if message.GetCost() {
		//花钱购买
		//interval := time.Now().Sub(searchFlag.UpdateTime).Seconds()
		//retio := int32(interval / conf.DATA.SearchCostPerTime) + 1
		//cost = retio * conf.DATA.SearchCost
		cost = conf.DATA.SearchCost
		user.EnsureGayPoint(cost)
	} else {
		//使用免费次数
		if searchFlag.Value <= 0 {
			exception.GameException(exception.SEARCH_COUNT_LIMIT)
		}
	}

	starType := user.GetStarType()

	request.SearchItem.StarType = starType
	groupBases := conf.DATA.STAR_ITEMGOURP_MAPPING[starType]

	unHaveGroups := []*conf.ItemGroupBase{}
	//haveGroups := []*conf.ItemGroupBase{}
	for _, groupBase := range groupBases {
		//group := user.GetItemGroup(groupBase.ID)
		//if group == nil {
		//	unHaveGroups = append(unHaveGroups, groupBase)
		//} else {
		//	haveGroups = append(haveGroups, groupBase)
		//}
		unHaveGroups = append(unHaveGroups, groupBase)
	}

	searchItems := []int32{}

	randomItem := RandomSearchItem(unHaveGroups, filter)
	if randomItem != 0 {
		searchItems = append(searchItems, randomItem)
	}
	randomItem = RandomSearchItem(unHaveGroups, filter)
	if randomItem != 0 {
		searchItems = append(searchItems, randomItem)
	}
	randomItem = RandomSearchItem(unHaveGroups, filter)
	if randomItem != 0 {
		searchItems = append(searchItems, randomItem)
	}
	randomItem = RandomSearchItem(unHaveGroups, filter)
	if randomItem != 0 {
		searchItems = append(searchItems, randomItem)
	}

	request.SearchItem.ItemIDs = searchItems

	response.SearchItemRet = rpc.SearchServiceProxy.HandleMessage(request).SearchItemRet
	//SearchResult

	user.UpdateSearch(response.SearchItemRet.Strangers)

	//TODO 圣物需要目标需要添加昵称

	if message.GetCost() {
		user.TakeOutGayPoint(cost, constant.OPT_TYPE_SEARCH_ITEM, 0)
		//user.WriteMsg(BuildRoleSocialPush(user))
		AddRoleSocialPush(user, response)
	} else {
		flag := user.UpdateFlag(db.FLAG_SEARCH_COUNT, searchFlag.Value-1)
		//user.WriteMsg(BuildFlagInfoPush(flag))
		AddFlagInfoPush(flag, response)
	}

	//补充机器人
	robotCount := 4 - len(response.SearchItemRet.Strangers)
	if robotCount > 0 {
		for i := 0; i < robotCount; i++ {
			randomIndex := rand.Intn(len(searchItems))
			response.SearchItemRet.Strangers = append(response.SearchItemRet.Strangers, &protocol.SearchResult{
				Id:       int32(-i - 1),
				Nickname: words.RandomName(),
				ItemID:   searchItems[randomIndex],
			})
		}
	}

}

func RandomSearchItem(randomScope []*conf.ItemGroupBase, filter []int32) int32 {
	length := len(randomScope)
	if length == 0 {
		return 0
	}
	groupBase := randomScope[rand.Intn(length)]
	randomData := []int32{}
	for _, item := range groupBase.Content {
		if !character.ContainsInt32(item, filter) {
			randomData = append(randomData, item)
		}
	}
	if len(randomData) == 0 {
		randomData = groupBase.Content
	}
	return randomData[rand.Intn(len(randomData))]
}

type NewsFeedListService struct {
}

func (service *NewsFeedListService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.GetNewsFeedListRet = &protocol.GetNewsfeedListRet{
		NewsFeeds: user.GetProtocolNewsFeeds(),
	}
}

type DealListService struct {
}

func (service *DealListService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.GetDealListRet = &protocol.GetDealListRet{
		Deals: user.GetProtocolDeals(),
	}
}

//type ItemRequestOverdueService struct {
//}
//
//func (service *ItemRequestOverdueService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetItemRequestOverdue()
//	deal := user.GetDeal(message.GetDealID())
//	if deal == nil {
//		response.ItemRequestOverdueRet = &protocol.ItemRequestOverdueRet{
//			Result: false,
//		}
//		return
//	}
//
//	overdue := util.IsOverdue(deal)
//	if overdue {
//		//退还友情点
//		if (deal.GetType() == constant.NEWSFEED_TYPE_REQUEST_ITEM) {
//			user.TakeInGayPoint(conf.DATA.FriendPointPrices)
//			//user.SetDirty()
//		}
//		user.RemoveDeal(message.GetDealID())
//		//user.SetDirty()
//	}
//	response.ItemRequestOverdueRet = &protocol.ItemRequestOverdueRet{
//		Result: overdue,
//	}
//}

type GlobalMessageService struct {
}

func (service *GlobalMessageService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.GetGlobalMessageRet = global.GetGlobalMessage()
}

//type PublicShareService struct {
//}
//
//func (service *PublicShareService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	shareCount := user.GetFlagValue(db.FLAG_SHARE)
//	if shareCount <= 0 {
//		exception.GameException(exception.SHARE_COUNT_NOT_ENOUGH)
//	}
//	message := request.GetPublicShare()
//	shareType := message.GetType()
//
//	var param1 = character.StringToInt32(message.GetRefID())
//	//var param3 int32 = 0
//	var ext []string = nil
//
//	if shareType == constant.NEWSFEED_TYPE_SHARE_ITEM {
//		if !user.CanTakeoutItem(param1, 1) {
//			exception.GameException(exception.ITEM_NOT_FOUND)
//		}
//		rpc.CommunityServiceProxy.PublicMoments(user.GetID(), constant.MOMENTS_TYPE_SHARE_ITEM, param1)
//
//	} else if shareType == constant.NEWSFEED_TYPE_SHARE_ITEM_GROUP {
//		if user.GetActiveGroupID(param1) {
//			//!user.IsDoneItemGroup(param1)
//			exception.GameException(exception.ITEM_GOURP_NOT_FINISH)
//		}
//		rpc.CommunityServiceProxy.PublicMoments(user.GetID(), constant.MOMENTS_TYPE_SHARE_ITEMGROUP, param1)
//	} else {
//		if shareType == constant.NEWSFEED_TYPE_SHARE_BELIEVER {
//			if user.GetLastUpgradeBelieverID() != message.GetRefID() {
//				exception.GameException(exception.STAR_BELIEVER_NOT_ENOUGH)
//			}
//			ext = []string{message.GetRefID()}
//		} else if shareType == constant.NEWSFEED_TYPE_SHARE_TASK {
//			task := user.GetLastTask()
//			if task == nil || task.GetId() != param1 {
//				exception.GameException(exception.TASK_NOT_FOUND)
//			}
//			param1 = task.GetBaseID()
//			data, _ := json.Marshal(task.GetReward())
//			ext = []string{string(data)}
//		}
//		newsFeed := BuildNewsFeed1(user.GetID(), shareType, param1, 0, 0, ext)
//		global.BroadcastMessage(&protocol.GlobalMessage{NewsFeed: newsFeed})
//	}
//
//	flag := user.UpdateFlag(db.FLAG_SHARE, shareCount-1)
//	//user.WriteMsg(BuildFlagInfoPush(flag))
//	AddFlagInfoPush(flag, response)
//	user.TakeInFaith(conf.DATA.ShareReward, constant.OPT_TYPE_SHARE, 0)
//	//user.SetDirty()
//	//user.WriteMsg(BuildRoleSocialPush(user))
//	AddRoleSocialPush(user, response)
//
//	response.PublicShareRet = &protocol.PublicShareRet{
//		Result: true,
//	}
//
//}

type PublicWechatShareService struct {
}

func (service *PublicWechatShareService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetPublicWechatShare()

	resp := &protocol.PublicWechatShareRet{Reward:&protocol.Reward{}}
	if message.GetType() == constant.PUBLIC_WECHAT_TYPE_HELP {
		switch message.GetRefType() {
		case constant.WECHAT_HELP_REF_BELIEVER:
			resp.ShareCount = user.DecFlag(db.FLAG_SHARE_WECHAT_HELP)
			addBeliever := conf.GenHelpBeliever(conf.DATA.HelpBeliever)
			addBeliever = rpc.StarServiceProxy.UpdateBelieverInfo(user.GetID(), addBeliever, int32(constant.OPT_TYPE_WECHAT_SHARE), true).GetUpdateBelieverRet().Believer
			resp.Reward.Believer = addBeliever
			break
		case constant.WECHAT_HELP_REF_FAITH:
			resp.ShareCount = user.DecFlag(db.FLAG_SHARE_WECHAT_HELP)
			resp.Reward.Faith = conf.DATA.HelpFaith
			user.TakeInFaith(resp.Reward.Faith, constant.OPT_TYPE_WECHAT_SHARE, 0)
			break
		case constant.WECHAT_HELP_REF_REPAIR:
			resp.BuildType = rpc.StarServiceProxy.HelpRepairBuildPublic(user.GetID(), message.GetRefNum()).GetBuildingType()
			break
		case constant.WECHAT_HELP_REF_ITEM:
			cache.UserCache.SetHelpItemWechatUid( 0, user.GetID(), message.GetRefNum())
			break
		default:
			break
		}
	}
	if message.GetType() == constant.PUBLIC_WECHAT_TYPE_SHOW {
		resp.ShareCount = user.DecFlag(db.FLAG_SHARE_WECHAT_SHOW)
		resp.Reward.Faith = conf.DATA.ShareReward
		user.TakeInFaith(resp.Reward.Faith, constant.OPT_TYPE_WECHAT_SHARE, 0)
		rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_GAIN_FAITH_SHARE, float64(resp.Reward.Faith), 0)
	}
	response.PublicWechatShareRet = resp
}

type DrawWechatShareRewardService struct {

}

func (service *DrawWechatShareRewardService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetDrawWechatShareReward()
	resp := &protocol.DrawWechatShareRewardRet{}
	var power int32 = 0
	var gayPoint int32 = 0
	switch message.GetNewsfeedType() {
	case constant.NEWSFEED_TYPE_SHARE_SUCC:
		power = conf.DATA.HelperReward
		resp.Reward = &protocol.Reward{Power: power}
		resp.ShareCount = user.DecFlag(db.FLAG_SHARE_WECHAT_SUCC)
		break
	case constant.NEWSFEED_TYPE_MUTUAL_FOLLOW:
		if !user.IsDrawInviteLoginGift() {
			power = conf.DATA.HelpBonusReward
			resp.Reward = &protocol.Reward{Power: power}
		}
		user.UpdateFlag(db.FLAG_DRAW_INVITE_GIFT, 1)
		break
	case constant.NEWSFEED_TYPE_HELP_ITEMHELP:
		gayPoint = conf.DATA.RelicShareReward
		resp.Reward = &protocol.Reward{GayPoint: gayPoint}
		break
	default:
		break
	}
	if power > 0 {
		user.TakeInPower(power, false, constant.OPT_TYPE_WECHAT_SHARE)
	}
	if gayPoint > 0 {
		user.TakeInGayPoint(gayPoint, constant.OPT_TYPE_WECHAT_SHARE, 0)
	}
	response.DrawWechatShareRewardRet = resp
}

type GetWechatShareTimeService struct {

}

func (service *GetWechatShareTimeService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	response.GetWechatShareTimeRet = &protocol.GetWechatShareTimeRet{NextDraw:user.GetNextAdPowerTime().Unix()}
}

type WatchAdSuccessService struct {

}

func (service *WatchAdSuccessService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetWatchAdSuccess()
	resp := &protocol.WatchAdSuccessRet{Reward:&protocol.Reward{}}

	switch message.GetType() {
	case constant.AD_TYPE_GET_POWER:
		if user.GetPower() > conf.DATA.RevengeMana  {
			response.WatchAdSuccessRet = resp
			return
			//exception.GameException(exception.POWER_IS_NOT_ZERO)
		}
		if user.GetStarFlagValue(constant.STAR_FLAG_FIRST_WATACH_AD) != constant.FLAG_VALUE_HAS_FIRST_WATACH_AD {
			flag := rpc.StarServiceProxy.UpdateStarFlag(user.GetID(), constant.STAR_FLAG_FIRST_WATACH_AD, constant.FLAG_VALUE_HAS_FIRST_WATACH_AD).GetFlag()
			if flag != nil {
				user.WriteMsg(util.BuildStarFlagPush(flag))
				user.SetStarFlag(constant.STAR_FLAG_FIRST_WATACH_AD, constant.FLAG_VALUE_UNLOCK)
				user.IsExceedShareTime(conf.DATA.HelpInterval)
				resp.Reward.Power = conf.DATA.FirstWatchAdPower
			}
		} else {
			if user.IsExceedShareTime(conf.DATA.HelpInterval) {
				resp.Reward.Power = conf.DATA.AdBonusReward
			} else {
				resp.Reward.Power = conf.DATA.AdReward
			}
		}
		//if user.IsFirstDrawAdReward(db.FLAG_WATCH_AD_GET_POWER) {
		//	flag := user.UpdateFlag(db.FLAG_WATCH_AD_GET_POWER, 1)
		//	AddFlagInfoPush(flag, response)
		//	resp.Reward.Power = conf.DATA.AdBonusReward
		//} else {
		//	resp.Reward.Power = conf.DATA.AdReward
		//}
	case constant.AD_TYPE_GET_GAY_POINT:
		resp.Reward.GayPoint = conf.DATA.HelpAdReward
	}
	if resp.Reward.Power > 0 {
		user.TakeInPower(resp.Reward.Power, false, constant.OPT_TYPE_WATCH_AD)
		rpc.StarServiceProxy.UpdateStarStatistics(user.GetID(), constant.STAR_STATISTIC_GAIN_POWER_AD, float64(resp.Reward.Power), 0)
	}
	if resp.Reward.GayPoint > 0 {
		succ := user.TakeInGayPoint(resp.Reward.GayPoint, constant.OPT_TYPE_WATCH_AD, 0)
		if !succ {
			exception.GameException(exception.GAYPOINT_NOT_ENOUGH)
		}
	}
	response.WatchAdSuccessRet = resp
}

//type RequestItemService struct {
//}
//
//func (service *RequestItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetRequestItem()
//	searchResult := user.GetSearch(message.GetSearchID())
//	if searchResult == nil {
//		exception.GameException(exception.SEARCH_RESULT_INVALID)
//	}
//
//	user.TakeOutGayPoint(conf.DATA.FriendPointPrices)
//	notifyUID := searchResult.GetId()
//
//	newDealMessage := BuildNewsFeedMessage(user.GetID(), constant.NEWSFEED_TYPE_BE_REQUEST_ITEM, searchResult.GetItemID(), 0, 0)
//	rpc.UserServiceProxy.PersistCall(notifyUID, newDealMessage)
//
//	deal := BuildNewsFeed(notifyUID, constant.NEWSFEED_TYPE_REQUEST_ITEM, searchResult.GetItemID(), 0, 0)
//	deal.Id = newDealMessage.AddNewsFeed.Id
//
//	user.AddNewsFeed(deal)
//	user.WriteMsg(BuildNewsFeedPush(deal))
//	//AddNewsFeedPush(deal, response)
//
//	user.RemoveSearch(message.GetSearchID())
//	//user.SetDirty()
//	response.RequestItemRet = &protocol.RequestItemRet{
//		Result: true,
//		DealID: deal.GetId(),
//	}
//}

func AddUserNewsFeed(user *user.Session, newsFeed *protocol.NewsFeed) {
	user.AddNewsFeed(newsFeed)
	user.WriteMsg(BuildNewsFeedPush(newsFeed))
	//AddNewsFeedPush(newsFeed, response)
}

//type AcceptItemRequestService struct {
//}
//
//func (service *AcceptItemRequestService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
//	message := request.GetAcceptItemRequest()
//	deal := user.GetDeal(message.GetDealID())
//	if (deal == nil || util.IsOverdue(deal)) {
//		exception.GameException(exception.ITEM_REQUEST_OVERDUE)
//	}
//
//	if deal.GetType() != constant.NEWSFEED_TYPE_BE_REQUEST_ITEM {
//		exception.GameException(exception.INVALID_PARAM)
//	}
//
//	itemID := deal.GetParam1()
//
//	user.TakeOutItem(itemID, 1, constant.OPT_TYPE_SALE, deal.GetRelateID())
//	user.TakeInGayPoint(conf.DATA.FriendPointPrices)
//	user.RemoveDeal(message.GetDealID())
//
//	//user.SetDirty()
//
//	acceptMessage := BuildNewsFeedMessage(user.GetID(), constant.NEWSFEED_TYPE_BE_ACCEPT, deal.GetParam1(), 0, 0)
//	acceptMessage.AddNewsFeed.Id = deal.Id
//
//	rpc.UserServiceProxy.PersistCall(deal.GetRelateID(), acceptMessage)
//
//	response.AcceptItemRequestRet = &protocol.AcceptItemRequestRet{
//		Result: true,
//	}
//
//}

type RejectItemRequestService struct {
}

func (service *RejectItemRequestService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetRejectItemRequest()
	deal := user.GetDeal(message.GetDealID())

	if deal == nil || util.IsOverdue(deal) {
		exception.GameException(exception.ITEM_REQUEST_OVERDUE)
	}

	if deal.GetType() != constant.NEWSFEED_TYPE_BE_REQUEST_ITEM {
		exception.GameException(exception.INVALID_PARAM)
	}

	rejectMessage := BuildNewsFeedMessage(user.GetID(), constant.NEWSFEED_TYPE_BE_REJECT, deal.GetParam1(), 0, 0)
	rejectMessage.AddNewsFeed.Id = deal.Id

	rpc.UserServiceProxy.PersistCall(deal.GetRelateID(), rejectMessage)
	user.RemoveDeal(message.GetDealID())
	//user.SetDirty()
	response.RejectItemRequestRet = &protocol.RejectItemRequestRet{
		Result: true,
	}
}

type GetOnNoticesService struct {

}

func (service *GetOnNoticesService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	resp := rpc.CenterServiceProxy.HandleMessage(request).GetGetOnNoticesRet()
	response.GetOnNoticesRet = resp
}

type RefreshMallItemService struct {

}

func (service *RefreshMallItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetRefreshMallItem()
	resp := &protocol.RefreshMallItemRet{}
	cost, costNext := user.CalcReMallItemCost(user.GetFlagValue(db.FLAG_RE_MALL_ITEM_COUNT))
	if message.GetIsWatchAd() {
		user.DecFlag(db.FLAG_AD_RE_MALL_ITEM_COUNT)
		items, _ := user.RefreshMallItems(false, user.GetStarType(),cost, costNext )
		AddFlagInfoPush(user.GetFlag(db.FLAG_AD_RE_MALL_ITEM_COUNT), response)
		resp.Cost = cost
		resp.Items = items
	} else {
		user.EnsureGayPoint(cost)
		items, _ := user.RefreshMallItems(false, user.GetStarType(),cost, costNext )
		user.AddFlag(db.FLAG_RE_MALL_ITEM_COUNT)
		user.TakeOutGayPoint(cost, constant.OPT_YTPE_REFRESH_MALL, 0)
		resp.Cost = costNext
		resp.Items = items
	}
	response.RefreshMallItemRet = resp
}

type BuyMallItemService struct {

}

func (service *BuyMallItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetBuyMallItem()
	item := user.BuyMallItem(message.GetID())
	user.EnsureGayPoint(item.GroupCost)
	user.TakeInItem(item.ItemID, item.Num, constant.OPT_TYPE_BUY_ITEM, 0)
	user.TakeOutGayPoint(item.GroupCost, constant.OPT_TYPE_BUY_ITEM, 0)
	response.BuyMallItemRet = &protocol.BuyMallItemRet{
		Item:item,
	}
}

type AutoRefreshMallItemService struct {

}

func (service *AutoRefreshMallItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	cost, costNext := user.CalcReMallItemCost(user.GetFlagValue(db.FLAG_RE_MALL_ITEM_COUNT))
	items, reItemTime := user.RefreshMallItems(true, user.GetStarType(), cost, costNext)
	response.AutoRefreshMallItemRet = &protocol.AutoRefreshMallItemRet{
		RefreshTime:reItemTime.Unix(),
		Items:items,
		Cost:cost,
	}
}

type GetMallItemService struct {

}

func (service *GetMallItemService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	//message := request.
	cost, _ := user.CalcReMallItemCost(user.GetFlagValue(db.FLAG_RE_MALL_ITEM_COUNT))
	items, reItemTime := user.GetMallItems(user.GetStarType())
	response.GetMallItemRet =  &protocol.GetMallItemRet{
		RefreshTime:reItemTime.Unix(),
		Items:items,
		Cost:cost,
	}
}

type GetActiveGroupService struct {

}

func (service *GetActiveGroupService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetGetCurrentGroup()
	groupID := rpc.StarServiceProxy.GetCurrentGroupItems(message.GetUid()).GetGroupID()
	response.GetCurrentGroupRet = &protocol.GetCurrentGroupRet{
		GroupID:groupID,
	}
}

//--------------------------压测-------------------------
type AddAttachService struct {

}

func (service *AddAttachService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetAddAttach()
	if message.GetPower() > 0 {
		user.TakeInPower(message.GetPower(), false, constant.OPT_TYPE_TEST)
		//user.SetDirty()
	}
	if message.GetPower() < 0 {
		user.TakeOutPower(-message.GetPower(), constant.OPT_TYPE_TEST)
	}

	if message.GetFaith() > 0 {
		user.TakeInFaith(message.GetFaith(), constant.OPT_TYPE_TEST, 0)
		//user.SetDirty()
	}
	if message.GetFaith() < 0 {
		user.TakeOutFaith(-message.GetFaith(), constant.OPT_TYPE_TEST, 0)
		//user.SetDirty()
	}

	if message.GetGayPoint() > 0 {
		user.TakeInGayPoint(message.GetGayPoint(), constant.OPT_TYPE_TEST, 0)
		//user.SetDirty()
	}
	if message.GetGayPoint() < 0 {
		user.TakeOutGayPoint(-message.GetGayPoint(), constant.OPT_TYPE_TEST, 0)
		//user.SetDirty()
	}

	if message.GetDiamond() > 0 {
		user.TakeInDiamond(message.GetDiamond(), constant.OPT_TYPE_TEST, 0)
		//user.SetDirty()
	}
	if message.GetDiamond() < 0 {
		user.TakeOutDiamond(-message.GetDiamond(), constant.OPT_TYPE_TEST, 0)
		//user.SetDirty()
	}

	if message.GetItems() != nil {
		for _, item := range message.GetItems() {
			if item.GetNum() > 0 {
				user.TakeInItem(item.GetId(), item.GetNum(), constant.OPT_TYPE_TEST, 0)
			}
			if item.GetNum() < 0 {
				user.TakeOutItem(item.GetId(), -item.GetNum(), constant.OPT_TYPE_TEST, 0)
			}
		}
		//user.SetDirty()
	}
	if message.GetBelievers() != nil {
		rpc.StarServiceProxy.UpdateBelieverInfo(user.GetID(), message.GetBelievers(), rpc.OP_BELIEVER_ADD, false)
	}
	response.AddAttachRet = &protocol.AddAttachRet{Result:true}
}

type RemoveSaleService struct {

}

func (service * RemoveSaleService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	message := request.GetRemoveSale()
	rpc.TradeServiceProxy.RemoveSale(message.GetId(), message.GetItemID())
	user.TakeInItem(message.GetItemID(), 1, constant.OPT_TYPE_SALE, 0)
	response.RemoveSaleRet = &protocol.RemoveSaleRet{
		Result:true,
	}
}

type SetBuildsService struct {

}

func (service *SetBuildsService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	resp := rpc.StarServiceProxy.Call(user.GetID(),request).GetSetBuildingsRet()
	response.SetBuildingsRet = resp
}

type SetBelieversService struct {

}

func (service *SetBelieversService) Request(request *protocol.C2GS, response *protocol.GS2C, user *user.Session) {
	resp := rpc.StarServiceProxy.Call(user.GetID(),request).GetSetBelieversRet()
	response.SetBelieversRet = resp
}
//---------------------------------------------------------------------------------