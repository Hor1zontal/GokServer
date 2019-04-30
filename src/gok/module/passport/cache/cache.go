package cache

import (
	basecache "gok/cache"
	"gok/module/passport/conf"
	"gok/module/passport/db"
	"time"
	"gok/service/exception"
	"aliens/log"
)

var UserCache = basecache.NewUserCacheManager()

func Init() {
	UserCache.Init1(conf.Server.RedisAddress, conf.Server.RedisPassword,
		conf.Server.RedisMaxActive, conf.Server.RedisMaxIdle, conf.Server.RedisIdleTimeout)

	//if UserCache.SetNX(basecache.FLAG_LOADUSER, 1) {
	//	log.Debug("start load passport data to redis cache...")
	//	var users []*db.DBUser
	//	db.DatabaseHandler.QueryAll(&db.DBUser{}, &users)
	//	UserCache.SetRegisterTotal(len(users))
	//	for _, user := range users {
	//		//log.Debug("update user %v", user)
	//		UserCache.SetUsernameUidMapping(user.Username, user.ID)
	//		UserCache.HSetUser(user.ID, user)
	//	}
	//	log.Debug("end load passport data to redis cache")
	//}
	//其他服务器加载了缓存就不需要再加载到缓存了
	if UserCache.SetNX(basecache.FLAG_LOADUSER, 1) {
		log.Debug("start load passport data to redis cache...")
		count := 0
		var users []*db.DBUser
		err := db.DatabaseHandler.QueryAllLimit(&db.DBUser{}, &users, 10000, func(data interface{}) bool {
			for _, user := range users {
				UserCache.SetUsernameUidMapping(user.Username, user.ID)
				UserCache.HSetUser(user.ID, user)
			}
			currLen := len(users)
			count += currLen
			return currLen == 0
		})
		if err != nil {
			log.Debug("end load passport data to redis cache")
		}
		log.Debug("end load passport data to redis cache count:%v", count)
	}
}

func Close() {

	//清除所有缓存数据
	UserCache.Close()
}

/**
 *  新建用户
 */
func NewUser(username string, password string, ip string, channel string, channelUID string, openID string, avatar string) *db.DBUser {
	user := &db.DBUser{
		Username: username,
		Password: password,
		Salt:     "",
		Channel:  channel,
		ChannelUID: channelUID,
		Mobile:   "",
		IP:       ip,
		OpenID:   openID,
		Status:   0,
		Avatar:   avatar,
		RegTime:  time.Now(),
		//LastLogin:time.Now(),
	}
	err1 := db.DatabaseHandler.Insert(user)
	if err1 != nil {
		exception.GameException(exception.USERNAME_EXISTS)
	}
	//log.Debug("new user %v", user)
	UserCache.SetUsernameUidMapping(user.Username, user.ID)
	UserCache.HSetUser(user.ID, user)
	return user
}

/**
 *  获取用户数据
 */
func GetUser(username string) *db.DBUser {
	uid := UserCache.GetUidByUsername(username)
	if uid == 0 {
		return nil
	}
	user := &db.DBUser{}
	UserCache.HGetUser(uid, user)
	user.ID = uid
	return user
}

func GetUserByUid(uid int32) *db.DBUser {
	if uid == 0 {
		return nil
	}
	user := &db.DBUser{}
	UserCache.HGetUser(uid, user)
	user.ID = uid
	return user
}