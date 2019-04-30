package main

import (
	"fmt"
	"strconv"
	"stresstest/base"
	"stresstest/object"
	"time"
	test_conf "stresstest/conf"
)

var g_Players map[int32]*object.Player

const (
	passport_server   string = "p29test.jzyx.com:18811" //账号服地址
	game_server       string = "p29test.jzyx.com:18812" //游戏服地址
	SYNC_TIME         int32  = 5                        //同步频率
	REGISTER_NEW_USER bool   = false                    //true表示注册新账号并登录 false表示使用老账号登录
)

var (
	gameserver    string
	accountserver string
	preaccount    string
	password      string
	accountnum    int
	synctime      int
	exist         bool
	secretkey     string
)

//stresstest -accountserver "127.0.0.1:18811" -gameserver "127.0.0.1:18812" -preaccount "test_" -password "11111111" -accountnum 3000 -synctime 5 -exist=true
func main() {
	//形参
	//gameserver := flag.String("gameserver", "", "gameserver address")
	//accountserver := flag.String("accountserver", "", "accountserver address")
	//preaccount := flag.String("preaccount", "", "pre account")
	//password := flag.String("password", "", "password")
	//accountnum := flag.Int("accountnum", "", "account number")
	//synctime := flag.Int("synctime", "", "synctime=5s")
	//exist := flag.Bool("exist", "", "account is exist?")


	//flag.StringVar(&gameserver, "gameserver", "", "gameserver address")
	//flag.StringVar(&accountserver, "accountserver", "", "accountserver address")
	//flag.StringVar(&preaccount, "preaccount", "", "pre account")
	//flag.StringVar(&password, "password", "", "password")
	//flag.StringVar(&secretkey, "secretkey", "", "cipher key")
	//flag.IntVar(&accountnum, "accountnum", 3000, "account number")
	//flag.IntVar(&synctime, "synctime", 5, "synctime=5s")
	//flag.BoolVar(&exist, "exist", false, "account is exist?")
	//
	//flag.Parse()
	gameserver = test_conf.BASE.GameServer
	accountserver = test_conf.BASE.AccountServer
	preaccount = test_conf.BASE.PreAccount
	password = test_conf.BASE.Password
	secretkey = test_conf.BASE.SecretKey
	accountnum = test_conf.BASE.AccountNum
	synctime = test_conf.BASE.SyncTime
	exist = !test_conf.INIT.Enable

	if gameserver == "" || accountserver == "" || preaccount == "" || password == "" {
		println("Please input correct params => stresstest -h")
		println("stresstest -accountserver \"127.0.0.1:18811\" -gameserver \"127.0.0.1:18812\" -preaccount \"test\" -password \"11111111\" -secretkey \"abcd\" -accountnum 3000 -synctime 5 -exist=true")
		return
	}

	println("secript key: " + secretkey)
	g_Players = make(map[int32]*object.Player)
	var accountTotalNum int
	if accountnum == 0 {
		for _, userInfo := range test_conf.INIT.UserInfo {
			accountTotalNum += userInfo.UserCount
			accountnum = accountTotalNum
		}
	}

	fmt.Printf("##Create %v connections\n", accountnum)
	time.Sleep(3 * time.Second)

	if exist {
		fmt.Println("##LoginSession account")
	} else {
		fmt.Println("##Register account")
	}


	if !exist {
		accountnum = 0
		for index, userInfo := range test_conf.INIT.UserInfo {
			PlayerInit(accountnum + 1, accountnum + userInfo.UserCount, int32(index), accountTotalNum)
			accountnum += userInfo.UserCount
		}


		//for index, count := range test_conf.INIT.LevelUserCount {
		//	PlayerInit(accountnum + 1, accountnum + count, int32(index))
		//	accountnum += count
		//}
	} else {
		PlayerInit(1, accountnum, 0, accountnum)
	}


	//for i := 1; i <= accountnum; i++ {
	//	p := &object.Player{}
	//	p.Username = preaccount + "_" + strconv.Itoa(i)
	//	p.Password = password
	//	//p.Init(int32(i), "127.0.0.1:3568")
	//	p.Init(int32(i), accountserver, gameserver, int64(synctime), secretkey, 0)
	//	g_Players[p.GetIdx()] = p
	//	if exist {
	//		p.AcceptOp(base.OP_LOGIN)
	//	} else {
	//		p.AcceptOp(base.OP_REGISTER)
	//	}
	//	//time.Sleep(10 * time.Millisecond)
	//}
	time.Sleep(300 * time.Second)
	//time.Sleep(10000000 * time.Second)
	//for _, p := range g_Players {
	//	if exist {
	//		p.AcceptOp(base.OP_LOGIN)
	//	} else {
	//		p.AcceptOp(base.OP_REGISTER)
	//	}
	//	time.Sleep(10 * time.Millisecond)
	//}

	//fmt.Println("##Sync data")
	//time.Sleep(3 * time.Second)
	//for _, p := range g_Players {
	//	p.AcceptOp(base.OP_SYNC)
	//}


	//for {
	//	time.Sleep(time.Duration(synctime) * time.Second)
	//}
}


func PlayerInit(id int, accountnum int, accountType int32, accountTotalNum int) {
	for i := id; i <= accountnum; i++ {
		p := &object.Player{}
		p.Username = preaccount + "_" + strconv.Itoa(i)
		p.Password = password
		//p.Init(int32(i), "127.0.0.1:3568")
		p.Init(int32(i), accountserver, gameserver, int64(synctime), secretkey, accountType, accountTotalNum)
		g_Players[p.GetIdx()] = p
		if exist {
			p.AcceptOp(base.OP_LOGIN)
		} else {
			p.AcceptOp(base.OP_REGISTER)
		}
		//time.Sleep(10 * time.Millisecond)
	}
}