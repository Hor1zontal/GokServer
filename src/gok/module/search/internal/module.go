package internal

import (
	"gok/module/search/conf"
	"gok/module/search/db"
	"gok/module/search/service"
	"github.com/name5566/leaf/module"
	"gok/module/base"
	"gok/module/search/cache"
	"gok/module/search/core"
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
	core.ItemHelpSearcher.Init()
	core.StarSearcher.Init()
	service.Init(ChanRPC)
}

func (m *Module) OnDestroy() {
	service.Close()
	db.Close()
	cache.Close()
	conf.Close()
}

