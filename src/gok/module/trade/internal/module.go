package internal

import (
	"gok/module/trade/conf"
	"gok/module/trade/service"
	"gok/module/trade/db"
	"gok/module/base"
	"github.com/name5566/leaf/module"
	"gok/module/trade/cache"
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
	db.Init()
	cache.Init()
	//core.Exchanger.Init()
	service.Init(ChanRPC)
}

func (m *Module) OnDestroy() {
	service.Close()
	cache.Close()
	db.Close()
}
