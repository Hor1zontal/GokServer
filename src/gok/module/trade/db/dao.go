package db

import (
	"gok/service/msg/protocol"
)

func (this *UserSale) BuildProtocol() *protocol.Sale {
	return &protocol.Sale{
		Id:this.UserID,
		ItemID:this.ItemID,
		PublicTime:this.PublicTime.Unix(),
		RefID:"",
	}
}

func (this *UserGoods) BuildProtocol() *protocol.Goods {
	return &protocol.Goods{
		Id:this.CompID.SubID2,
		Num:this.Num,
		Price:this.Price,
	}
}

