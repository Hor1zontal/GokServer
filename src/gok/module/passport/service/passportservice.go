package service

import (
	"aliens/common"
	"gok/service/msg/protocol"
	baseservice "gok/service"
	basecache "gok/cache"
	"gok/module/passport/cache"
	clustercache "gok/module/cluster/cache"
	"gok/module/passport/helper"
	"gok/service/exception"
	"time"
	"aliens/log"
	"gok/module/passport/conf"
	"gopkg.in/mgo.v2/bson"
	"gok/module/passport/db"
	"github.com/name5566/leaf/chanrpc"
	"gok/service/lpc"
	"gok/service/rpc"
	"gok/constant"
	"gok/module/passport/core"
	"gok/module/passport/version"
	"gok/module/passport/notify"
	"aliens/common/character"
)

var PassportRPCService *baseservice.GRPCService = nil

func Init(chanRpc *chanrpc.Server) {
	var passportService = baseservice.NewLocalService(baseservice.SERVICE_PASSPORT_RPC)
	passportService.RegisterHandler(6, new(PassportRegisterService)) //注册
	passportService.RegisterHandler(7, new(PassportLoginService))    //登录
	passportService.RegisterHandler(9, new(LoginServerService))      //token登录

	passportService.RegisterHandler(650, new(ChannelLoginService))      //渠道登录
	passportService.RegisterHandler(651, new(UserStateService))      //修改用户状态

	passportService.RegisterHandler(652, new(RefreshClientVersionService))

	passportService.RegisterHandler(653, new(WechatEventPushService))
	passportService.RegisterHandler(654, new(CleanTestAccountService)) //清除测试账号
	passportService.RegisterHandler(655, new(QueryByUsernameService))


	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_CENTER_RPC)
	baseservice.ServiceManager.SubscribeRemoteService(baseservice.SERVICE_USER_RPC)
	////配置了RPC，需要发布服务到ZK
	PassportRPCService = baseservice.PublicRPCService1(passportService, conf.Server.RPCAddress, conf.Server.RPCPort, chanRpc)

	baseservice.SubscribeWBService(baseservice.SERVICE_USER)

}

func Close() {
	//LocalPassportRPCService.Close()
	PassportRPCService.Close()
}



type LoginServerService struct {
}

func (service *LoginServerService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	loginRequest := request.GetLoginServer()
	uid := loginRequest.GetUserId()

	helper.CheckState(uid, false, "")
	if uid <= 0 || loginRequest.GetToken() == "" {
		exception.GameException(exception.USER_INVALID_TOKEN)
	}

	if !cache.UserCache.IsUserExist(uid) {
		exception.GameException(exception.USER_NOTFOUND)
	}

	//TODO 验证用户是否封号

	token := cache.UserCache.GetUserToken(uid)
	if token == "" || token != loginRequest.GetToken() {
		log.Debug("login server uid: %v token: %v tokenParam: %v  time: %v", uid, token, loginRequest.GetToken(), time.Now())
		exception.GameException(exception.USER_INVALID_TOKEN)
	}
	if conf.Maintain.IsCheckVersion && version.VersionManager.VersionInfoMapping[loginRequest.GetVersion()] != conf.Server.ServerInfo {
		log.Info("client version: %v||get server version: %v != current server version:%v",loginRequest.GetVersion(),version.VersionManager.VersionInfoMapping[loginRequest.GetVersion()],conf.Server.ServerInfo)
		exception.GameException(exception.CLIENT_NOT_MATAH)
	}
}

//登录账号服务器请求
type PassportLoginService struct {
}

func (service *PassportLoginService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetLoginLogin()
	username := message.GetUsername()
	passwd := message.GetPassword()
	result := &protocol.LoginLoginRet{}
	response.LoginLoginRet = result
	userCache := cache.GetUser(username)
	if userCache == nil {
		//TODO 后续可以做成缓存读不到去数据库并写回到缓存,要考虑数据穿透的情况
		result.Result = protocol.Login_Result_invalidUser
		return
	}
	helper.CheckState(userCache.ID, false, "")
	passwordHash := helper.PasswordHash(username, passwd)
	//密码不对
	if passwordHash != userCache.Password {
		result.Result = protocol.Login_Result_invalidPwd
		return
	}
	gameServer := helper.AllocGameServer(userCache.ID)
	if gameServer == "" {
		result.Result = protocol.Login_Result_invalidGameServer
		return
	}

	result.Uid = userCache.ID
	token := util.Rand().Hex()
	cache.UserCache.SetUserToken(userCache.ID, token, conf.Server.TokenExpireTime)
	result.Token = token
	result.GameServer = gameServer
	result.Result = protocol.Login_Result_loginSuccess
}

//账号服务器请求
type PassportRegisterService struct {
}

func (service *PassportRegisterService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetLoginRegister()
	username := message.GetUsername()
	passwd := message.GetPassword()
	result := &protocol.LoginRegisterRet{}
	response.LoginRegisterRet = result
	helper.CheckState(0, true, "")

	if cache.UserCache.IsUsernameExist(username) {
		result.Msg = "用户名已存在"
		result.Result = protocol.Register_Result_userExists
		return
	}

	passwd = helper.PasswordHash(username, passwd)
	//TODO 有风险最好数据库再加一层判断
	userCache := cache.NewUser(username, passwd, "", "", "", "", "")
	gameServer := helper.AllocGameServer(0)
	if gameServer == "" {
		result.Result = protocol.Register_Result_invalidServer
		return
	}
	result.Result = protocol.Register_Result_registerSuccess
	result.Uid = userCache.ID
	token := util.Rand().Hex()
	cache.UserCache.SetUserToken(userCache.ID, token, conf.Server.TokenExpireTime)
	result.Token = token

	lpc.LogServiceProxy.AddRegisterRecord(userCache.ID, userCache.Channel)
	result.GameServer = gameServer
}


//登录账号服务器请求
type ChannelLoginService struct {
}

func (service *ChannelLoginService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
		message := request.GetChannelLogin()
		response.ChannelLoginRet = core.ChannelLogin(message.GetChannel(), message.GetChannelUID(), message.GetOpenID(), message.GetAvatar(), message.GetNickname())
}

type UserStateService struct {
}

func (service *UserStateService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetChangeUserState()
	uid := message.GetUid()
	status := message.GetState()
	//更新账号状态
	qdoc := bson.M{"_id": uid}
	udoc := bson.M{"$set": bson.M{"status": status}}
	//db.DatabaseHandler.Opt("user", qdoc, udoc)
	lpc.DBServiceProxy.UpdateCondition("user", qdoc, udoc, db.DatabaseHandler)
	cache.UserCache.SetUserAttr(uid, basecache.UPROP_STATUS, status)
	if status != int32(db.USER_STATUS_NONE) {
		node := clustercache.Cluster.GetUserNode(uid)
		if node != "" {
			rpc.UserServiceProxy.KickOut(uid, node, constant.LOGOUT_TYPE_GATE_CLOSE)
		}
	}
	response.ChangeUserStateRet = &protocol.ChangeUserStateRet{Result:true}
}

type RefreshClientVersionService struct {
}

func (service *RefreshClientVersionService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	ret := version.VersionManager.RefreshVersionInfo()
	response.RefreshClientVersionRet = &protocol.RefreshClientVersionRet{Result:ret}
}

type WechatEventPushService struct {
}

func (service *WechatEventPushService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetWechatEventPush()
	notify.PushEventMsg(message.GetEvent(), message.GetUid(), int(message.GetDelay()))
}

type CleanTestAccountService struct {

}

func (service *CleanTestAccountService) Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetCleanTestAccount()
	response.CleanTestAccountRet = &protocol.CleanTestAccountRet{}
	resp := response.CleanTestAccountRet
	//userCache := cache.GetUserByUid(message.GetUid())
	//if userCache == nil {
	//	return
	//}
	userCache := &db.DBUser{}
	err := db.DatabaseHandler.QueryOneCondition(userCache, "_id", message.GetUid())
	if err != nil {
		resp.Message = "err" + err.Error()
		return
	}
	//log.Info("uid: %v ||||| userCache: %v", message.GetUid(), userCache)
	username := userCache.Username

	if !cache.UserCache.ExistTestUserID(message.GetUid()) {
		resp.Message = "test user not exist"
		return
	}
	unionID := cache.UserCache.GetUnionIDByUserID(userCache.ID)
	//log.Info("unionID:%v",unionID)
	if unionID == "" {
		resp.Message = "unionID is null"
		return
	}
	node := clustercache.Cluster.GetUserNode(userCache.ID)
	rpc.UserServiceProxy.KickOut(userCache.ID, node, constant.LOGOUT_TYPE_CLEAN_ACCOUNT)

	//清除unionid与uid的映射
	cache.UserCache.CleanUserIDUnionIDMapping(userCache.ID, unionID)
	//清除uid与username的映射
	cache.UserCache.DelUidByUsername(username)
	//设置测试账号之前的uid与username的映射关系
	cleanusername := username + "_" + character.Int64ToString(time.Now().Unix())
	cache.UserCache.SetUsernameUidMapping(cleanusername, userCache.ID)
	userCache.Username = cleanusername
	err1 := db.DatabaseHandler.UpdateOne(userCache)
	if err != nil {
		resp.Message = "err1" + err1.Error()
		return
	}
	cache.UserCache.DelTestUserID(userCache.ID)
	log.Info("clean test account %v", userCache.ID)
	resp.Result = true

	response.CleanTestAccountRet = resp
}

type QueryByUsernameService struct {

}

func (service *QueryByUsernameService)Request(request *protocol.C2GS, response *protocol.GS2C, network baseservice.IMessageChannel) {
	message := request.GetQueryByUsername()
	var users []*db.DBUser
	username := message.GetUsername()
	db.DatabaseHandler.QueryAllCondition(&db.DBUser{}, "username", username, &users)
	var result []int32
	for _, user := range users {
		result = append(result, user.ID)
	}
	response.QueryByUsernameRet = &protocol.QueryByUsernameRet{Uids:result}
}