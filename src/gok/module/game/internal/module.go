package internal

import (
	"gok/module/game/db"
	"gok/module/game/user"
	"gok/module/game/cache"
	"gok/module/game/service"
	"gok/module/game/rank"
	//"gok/module/game/gamelog"
	"github.com/name5566/leaf/module"
	"gok/module/game/conf"
	//"gok/module/game/global"
	"gok/module/base"
	"gok/module/game/global"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}


func (m *Module) IsEnable() bool {
	return conf.Server.Enable
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	conf.Init()
	db.Init()
	cache.Init()
	rank.Init()
	//task.Init()
	user.Init()
	global.Init(skeleton)
	service.Init(ChanRPC)
}

func (m *Module) OnDestroy() {
	service.Close()
	//task.Close()
	//最好和初始化反向
	user.Close()
	cache.Close()
	db.Close()
	conf.Close()
}

//func (s *Module) Run(closeSig chan bool) {
//	go global.Handler.Run(closeSig)
//	go s.Skeleton.Run(closeSig)
//}

