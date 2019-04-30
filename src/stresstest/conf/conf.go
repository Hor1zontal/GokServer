package conf

import (
	"encoding/json"
	"io/ioutil"
	"aliens/log"
)

//var Client struct{
//	UserEnable		bool
//	StarEnable		bool
//	DialEnable		bool
//	MailEnable		bool
//	CommunityEnable bool
//	TradeEnable		bool
//	InitEnable		bool
//}
var BASE struct {
	GameServer			string
	AccountServer		string
	PreAccount			string
	Password			string
	SecretKey			string
	SyncTime			int
	AccountNum			int
}


var DIAL struct {
	Enable				bool
	TargetSelectWeight	[3]int32
	Times				int32
}

var INIT struct {
	Enable				bool
	UserInfo			[]*userInfo
}

type userInfo struct {
	BuildRange			[2]int32
	BelieverRange		[2]int32
	UserCount			int
}

var USER struct {
	Enable		bool
}

var STAR struct {
	Enable		bool
}

var MAIL struct {
	Enable		bool
}

var COMMUNITY struct {
	Enable bool
}

var TRADE struct {
	Enable		bool
}

func init() {
	//测试配置
	data, err := ioutil.ReadFile("conf/base.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &BASE)
	if err != nil {
		log.Fatal("%v", err)
	}

	//初始化测试套配置
	data, err = ioutil.ReadFile("conf/init.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &INIT)
	if err != nil {
		log.Fatal("%v", err)
	}

	//转盘测试套配置
	data, err = ioutil.ReadFile("conf/dial.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &DIAL)
	if err != nil {
		log.Fatal("%v", err)
	}

	//用户模块测试套配置
	data, err = ioutil.ReadFile("conf/user.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &USER)
	if err != nil {
		log.Fatal("%v", err)
	}

	//邮件测试套配置
	data, err = ioutil.ReadFile("conf/mail.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &MAIL)
	if err != nil {
		log.Fatal("%v", err)
	}

	//星球模块测试套配置
	data, err = ioutil.ReadFile("conf/star.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &STAR)
	if err != nil {
		log.Fatal("%v", err)
	}

	//社交模块测试套配置
	data, err = ioutil.ReadFile("conf/community.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &COMMUNITY)
	if err != nil {
		log.Fatal("%v", err)
	}

	//交易模块测试套配置
	data, err = ioutil.ReadFile("conf/trade.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &TRADE)
	if err != nil {
		log.Fatal("%v", err)
	}

}
