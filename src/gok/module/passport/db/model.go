package db

import (
	"time"
)

const (
	//collection constants
	//increment id constants
	//ID_USER = "uid"

	USER_STATUS_NONE byte = 0 //正常
	USER_STATUS_NOT_AUTH byte = 1 //被封号
)

//角色
type DBUser struct {
	ID       int32  `bson:"_id" gorm:"AUTO_INCREMENT"`             //用户id
	Username string `bson:"username"  unique:"true" rorm:"uname"` //用户名 渠道信息_渠道用户id存Username
	Password string `bson:"password"  rorm:"pwd"`                  //加密的密码
	Salt     string `bson:"salt"      rorm:"salt"`                 //加密的salt

	ChannelUID string `bson:"cuid"    rorm:"cuid"`    //用户的渠道的渠道用户id
	Channel    string `bson:"channel" rorm:"channel"` //用户的渠道信息 渠道用户id存Username
	Avatar	   string `bson:"avatar"  rorm:"avatar"` //用户的头像地址

	Mobile string `bson:"mobile"` //用户电话
	OpenID string `bson:"openid"` //微信OPENID 绑定微信填写
	IP     string `bson:"ip"`     //最后一次登录的ip
	Status  byte      `bson:"status" rorm:"status"`  //用户状态 0正常  1封号
	RegTime time.Time `bson:"regtime"` //用户注册时间
}