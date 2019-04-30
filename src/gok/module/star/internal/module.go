package internal

import (
	"gok/module/star/conf"
	"gok/module/star/cache"
	"gok/module/star/session"
	"gok/module/star/service"
	"gok/module/star/db"
	"github.com/name5566/leaf/module"
	"gok/module/base"
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
	session.StarManager.Init()
	service.Init(ChanRPC)

}

func (m *Module) OnDestroy() {
	service.Close()
	session.StarManager.Close()
	cache.Close()
	db.Close()
	conf.Close()
}
