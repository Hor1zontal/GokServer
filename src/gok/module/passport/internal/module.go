package internal

import (
	"gok/module/passport/cache"
	"gok/module/passport/conf"
	"gok/module/passport/db"
	"gok/module/passport/service"
	"gok/module/base"
	"github.com/name5566/leaf/module"
	"gok/module/passport/service/http"
	"gok/module/passport/vivo"
	"gok/module/passport/wx"
	"gok/module/passport/version"
	"gok/module/passport/notify"
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
	version.Init()
	db.Init()
	cache.Init()
	service.Init(ChanRPC)
	http.Init(ChanRPC)
	wx.Init(conf.Server.WeChat)
	vivo.Init(conf.Server.Vivo)

	notify.Init()
	Init()
}

func (m *Module) OnDestroy() {
	http.Close()
	service.Close()
	db.Close()
	cache.Close()
	conf.Close()

}