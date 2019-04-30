package internal

import (
	"github.com/name5566/leaf/module"
	"gok/module/mail/service"
	"gok/module/base"
	"gok/module/mail/conf"
	"gok/module/mail/db"
	"gok/module/mail/cache"
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
	service.Init(ChanRPC)
}

func (m *Module) OnDestroy() {
	//http.Close()
	service.Close()
	cache.Close()
	db.Close()
	conf.Close()

}
