package http

import (
	"net/http"
	"aliens/log"
	"gok/module/center/conf"
	"gok/service/rpc"
	"gok/service/msg/protocol"
	"gok/module/center/order"
	"aliens/common/character"
	"aliens/common/helper"
	"gok/service/lpc"
)

const (
	PARAM_PAY_APPID = "appId"
	PARAM_ACCOUNTID = "accountId"
	PARAM_AMOUNT    = "amount"
	PARAM_CP_ORDER  = "cpOrderId"
	PARAM_ORDER  = "orderId"
)

func Pay(responseWriter http.ResponseWriter, request *http.Request) {
	payProxy(request)
	helper.SendToClient(responseWriter, "success")
}

func payProxy(request *http.Request) {
	request.ParseForm()
	//appID := request.FormValue(PARAM_PAY_APPID)
	//amount := request.FormValue(PARAM_AMOUNT)
	cpOrderId := request.FormValue(PARAM_CP_ORDER)
	//sign := request.FormValue(PARAM_SIGN)
	orderId := request.FormValue(PARAM_ORDER)

	//accountID := request.FormValue(PARAM_ACCOUNTID)

	//signText := fmt.Sprintf("accountId=%v&amount=%v&appId=%v&cpOrderId=%v&orderId=%v%v", accountID, amount, appID, cpOrderId, orderId, conf.Server.AppKey)
	//log.Debug("pay:  param : %v", signText)
	//TODO
	//signText := accountID + amount + appID + cpOrderId + orderId + conf.Server.AppKey
	succ := isSignSuccess(request.Form)
	if !succ {
		return
	}
	//signResult := cipher.MD5Hash(signText)
	//
	//if signResult != sign {
	//	log.Debug("signResult %v : sign : %v", signResult, sign)
	//	return
	//}

	//PAY_LOCK.Lock()
	//defer PAY_LOCK.Unlock()
	tempOrder := order.GetTempOrder(cpOrderId)
	if tempOrder == nil {
		log.Debug("order %v : not found", cpOrderId)
		return
	}
	order.RemoveTempOrder(cpOrderId)
	log.Debug("new order %v", tempOrder)
	lpc.LogServiceProxy.AddOrderRecord(tempOrder.ID, tempOrder.UserID, tempOrder.ProductID, tempOrder.Amount)

	shopBase := conf.DATA.ShopBaseData[tempOrder.ProductID]
	if shopBase == nil {
		log.Debug("order %v : shop item %v not found", orderId,tempOrder.ProductID)
		return
	}

	//表示购买钻石
	if shopBase.Type == 0x01 {
		//mail := mail.Manager.CreateMail(tempOrder.UserID, cpOrderId, character.Int32ToString(tempOrder.ProductID), &db.DBMailAttach{
		//	Diamond:shopBase.Amount,
		//})
		mail := rpc.MailServiceProxy.CreateMail(tempOrder.UserID, cpOrderId, character.Int32ToString(tempOrder.ProductID), `"diamond:"` + character.Int32ToString(shopBase.Amount)).GetMail()
		if mail == nil {
			log.Debug("order %v : create mail exception", orderId)
			return
		}
		rpc.UserServiceProxy.Push(tempOrder.UserID, &protocol.GS2C{
			Sequence:[]int32{1059},
			MailPush: mail,
		})
	}

	//发放商品道具
	//rpc.UserServiceProxy.PersistCall(tempOrder.UserID, cache.UserCache.GetUserNode(tempOrder.UserID), &protocol.C2GS{
	//	Sequence: []int32{524},
	//	AddShopItem: tempOrder.ProductID),
	//})
}


