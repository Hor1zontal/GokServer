package internal

import (
	"gok/module/community/conf"
	"gok/module/community/cache"
	"gok/module/community/service"
	"gok/module/community/db"
	"gok/module/community/core"
	"gok/module/base"
	"github.com/name5566/leaf/module"
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
	core.Moments.Init()
	service.Init(ChanRPC)
}

func (m *Module) OnDestroy() {
	service.Close()
	core.Moments.Close()
	cache.Close()
	db.Close()
}
