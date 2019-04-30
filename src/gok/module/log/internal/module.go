package internal

import (
	//"gok/module/game/gamelog"
	"github.com/name5566/leaf/module"
	//"gok/module/game/global"
	"gok/module/base"
	"gok/module/log/db"
	"gok/module/log/conf"
	"gok/module/log/cache"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}


func (m *Module) IsEnable() bool {
	return conf.Config.Enable
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	db.Init()
	cache.Init()

}

func (m *Module) OnDestroy() {
	cache.Close()
	db.Close()
}