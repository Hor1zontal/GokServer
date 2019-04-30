package core

import (
	"gok/module/passport/cache"
	"gok/module/passport/helper"
	"gok/service/exception"
	"gok/service/msg/protocol"
	"gopkg.in/mgo.v2/bson"
	"gok/service/lpc"
	"gok/module/passport/conf"
	"gok/module/passport/db"
)

func ChannelLogin(channel string, channelUID string, openID string, avatar string, nickName string) *protocol.ChannelLoginRet {

	//username := channel + "_" + channelUID
	username := "gok_" + channelUID
	userCache := cache.GetUser(username)
	new := false
	if userCache == nil {
		passwd := helper.PasswordHash(username, conf.Server.DefaultChannelPWD)
		userCache = cache.NewUser(username, passwd, "", channel, channelUID, openID, avatar)
		lpc.LogServiceProxy.AddRegisterRecord(userCache.ID, userCache.Channel)
		new = true
	}

	code := helper.GetCheckState(userCache.ID, new, "")
	if code != exception.NONE {
		return &protocol.ChannelLoginRet{
			Result: int32(code),
		}
	}
	gameServer := helper.AllocGameServer(userCache.ID)
	if gameServer == "" {
		return &protocol.ChannelLoginRet{
			Result: int32(exception.GAMESERVER_NOT_FOUND),
		}
	}

	//头像有变更,需要更新缓存和数据库
	if avatar != "" && avatar != userCache.Avatar {
		cache.UserCache.SetUserAvatar(userCache.ID, avatar)
		qdoc := bson.M{"_id": userCache.ID}
		udoc := bson.M{"$set": bson.M{"avatar": avatar}}
		lpc.DBServiceProxy.UpdateCondition("user", qdoc, udoc, db.DatabaseHandler)
		//db.DatabaseHandler.Opt("user", qdoc, udoc)
	}
	token := helper.NewToken()
	cache.UserCache.SetUserToken(userCache.ID, token, conf.Server.TokenExpireTime)
	cache.UserCache.SetUserNickname(userCache.ID, nickName)

	return &protocol.ChannelLoginRet{
		Result:0,
		Uid: userCache.ID,
		Token: token,
		New:new,
		GameServer:gameServer,
	}
}


