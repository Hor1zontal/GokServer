package event

import (
	"gok/service/msg/protocol"
	"aliens/common/character"
	"gok/constant"
)

type BelieverInfo struct {
	ID  string `bson:"id"`           //信徒id
	Num int32 `bson:"num"`          //信徒数量
}

func (this *BelieverInfo) getProtocol() *protocol.BelieverInfo {
	return &protocol.BelieverInfo{
		Id:this.ID,
		Num:this.Num,
	}
}


//构建信徒id
func buildBelieverID(level int32) string {
	if level <= 0 || level > constant.MAX_BELIEVER_LEVEL {
		level = 1
	}
	return "b01" + character.Int32ToString(level) + constant.RandSex()
}