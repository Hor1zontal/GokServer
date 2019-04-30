package db

import (
	"gok/service/msg/protocol"
)

func (this *DMoments) BuildProtocol() *protocol.MomentInfo {
	return &protocol.MomentInfo{
		Id:this.ID,
		Uid:this.Uid,
		CreateTime:this.CreateTime.Unix(),
		Data:this.Data,
	}
}