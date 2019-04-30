package internal

import (
	"gok/module/cluster/center"
	"gok/module/cluster/cache"
)

type Module struct {
}

func (m *Module) IsEnable() bool {
	return true
}

func (m *Module) OnInit() {
	cache.Init()
	center.Init()
}

func (m *Module) OnDestroy() {
	center.Close()
	cache.Close()
}

func (s *Module) Run(closeSig chan bool) {

}
