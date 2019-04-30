package internal

import (
	"github.com/name5566/leaf/gate"
	"gok/module/gate/conf"
	"gok/module/game"
)

type Module struct {
	*gate.Gate
}

func (m *Module) IsEnable() bool {
	return conf.Config.Enable
}


func (m *Module) OnInit() {
	m.Gate = &gate.Gate{
		MaxConnNum:      conf.Config.MaxConnNum,
		PendingWriteNum: conf.PendingWriteNum,
		MaxMsgLen:       conf.MaxMsgLen,
		WSAddr:          conf.Config.WSAddr,
		HTTPTimeout:     conf.HTTPTimeout,
		TCPAddr:         conf.Config.TCPAddr,
		LenMsgLen:       conf.LenMsgLen,
		LittleEndian:    conf.LittleEndian,
		Processor:       Processor,
		AgentChanRPC:    game.ChanRPC,
	}
}
