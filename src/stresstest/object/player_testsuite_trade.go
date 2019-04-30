package object

import (
	"stresstest/testcase"
	"gok/service/msg/protocol"
	"stresstest/message"
)

func (this *Player) TradeTestSuite() testcase.TestSuite {
	tradeTestCases := []*protocol.C2GS{
		//BuildInternalAddAttach(this.starType, 0, 0, 0, 0 , true, 1,1, false, 0, 0),
		message.BuildAddItems(this.starType, 1,1),
		message.BuildPublicGoods(this.Uid, this.starType, 1),
		message.BuildGetGoodsInfo(this.Uid),
		message.BuildCancelGoods(this.Uid, this.starType, 1),
		message.BuildAddSale(this.Uid, this.starType, 1),
		message.BuildGetSale(this.Uid),
		message.BuildGetSales([]int32{this.Uid}),
		message.BuildRemoveSale(this.Uid, this.starType, 1),
	}
	return testcase.NewTestSuite("trade", true, tradeTestCases)
}


