package helper

import (
	"crypto/md5"
	"encoding/hex"
	"aliens/common"
	"gok/service/exception"
	"time"
	basecache "gok/cache"
	"gok/module/passport/conf"
	"gok/constant"
	baseservice "gok/service"
	clustercache "gok/module/cluster/cache"
	"gok/module/passport/cache"
	"gok/module/passport/db"
	"gok/module/cluster/center"
)

func MD5Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	md5Hash := hex.EncodeToString(h.Sum(nil))
	return md5Hash
}

func PasswordHash(username string, passwd string) string {
	//h.Write([]byte(passwd + userCache.Salt))
	return MD5Hash(username + MD5Hash(passwd))
}

func NewToken() string {
	return util.Rand().Hex()
}

func GetCheckState(uid int32, newUser bool, ip string) exception.GameCode {
	if newUser && conf.GetServerState() == constant.SERVER_STATE_CLOSE_NEW {
		return exception.SERVER_CLOSE_NEW
	}
	if conf.IsMaintain() {
		if !clustercache.Cluster.IsUidWhiteList(uid) &&
			!clustercache.Cluster.IsIpWhiteList(ip) {
			return exception.SERVER_MAINTAIN
		}
	}
	//停服
	if !conf.IsServerOpen(time.Now()) {
		return exception.SERVER_CLOSE
	}

	if !newUser {
		status := cache.UserCache.GetUserAttrInt32(uid, basecache.UPROP_STATUS)
		//用户是否被封号
		if byte(status) == db.USER_STATUS_NOT_AUTH {
			return exception.USER_FORBIDDEN
		}
	}
	return exception.NONE
}

//检查能否登录
func CheckState(uid int32, newUser bool, ip string) {
	gameCode := GetCheckState(uid, newUser, ip)
	if gameCode != exception.NONE {
		exception.GameException(gameCode)
	}
}

//分配游戏服务器
//分配游戏服务器
func AllocGameServer(uid int32) string {
	node := ""
	if uid != 0 {
		node = clustercache.Cluster.GetUserNode(uid)
	}
	if node == "" {
		node = center.GetServerNode()
	}
	service := center.ClusterCenter.GetService(baseservice.SERVICE_USER, node)
	if service == nil {
		service = center.ClusterCenter.AllocService(baseservice.SERVICE_USER)
	}
	//找不到用户节点所在的服务器，随机分配一个用户服务器
	if service == nil {
		return ""
	}
	wbService, ok := service.(*baseservice.WBService)
	if ok {
		return wbService.Address
	}
	return ""
}