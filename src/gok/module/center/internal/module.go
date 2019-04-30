package internal

import (
	"gok/module/center/conf"
	"gok/module/center/service/http"
	"gok/module/center/cache"
	"gok/module/base"
	"github.com/name5566/leaf/module"
	"gok/module/center/service"
	"gok/module/center/db"
	"gok/module/center/notice"
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
	notice.Init()
	service.Init(ChanRPC)
	http.Init(ChanRPC)
}

func (m *Module) OnDestroy() {
	http.Close()
	service.Close()
	cache.Close()
	db.Close()
	conf.Close()

}
