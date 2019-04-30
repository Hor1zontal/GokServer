/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package lpc

import (
	"time"
	"gok/module/log"
	"gok/module/log/db"
	"gok/constant"
)

var LogServiceProxy = &logHandler{}

type logHandler struct {

}

func (this *logHandler) AddLoginRecord(uid int32, ip string, loginTime time.Time, logoutTime time.Time) {
	//fields := logrus.Fields{}
	//fields["uid"] = uid
	//fields["ip"] = ip
	//
	//this.ESDayLog(constant.ES_LOG_USER, "", fields)
	//this.ESLog(constant.ES_LOG_USER, "", fields)
	data := &db.LoginRecord{UserID:uid, Ip:ip, LoginTime:loginTime, LogoutTime:logoutTime}
	log.ChanRPC.Go(constant.LOG_COMMAND, data)
}

func (this *logHandler) AddRegisterRecord(uid int32, channel string) {
	addTime := time.Now()
	data := &db.RegisterRecord{UserID:uid, Channel:channel, Time:addTime}
	log.ChanRPC.Go(constant.LOG_COMMAND, data)
}

func (this *logHandler) AddGuideRecord(uid int32, guideTime time.Time, logoutTime time.Time) {
	data := &db.GuideRecord{ID:uid, Time:guideTime, Duration: logoutTime.Sub(guideTime).Seconds()}
	log.ChanRPC.Go(constant.LOG_COMMAND, data)
}

func (this *logHandler) AddOrderRecord(orderID string, uid int32, productID int32, amount float64) {
	data := &db.OrderRecord{ID: orderID, UserID:uid, ProductID:productID, Amount:amount, Time:time.Now()}
	log.ChanRPC.Go(constant.LOG_COMMAND, data)
}

func (this *logHandler) AddItemRecord(uid int32, itemID int32, refID int32, operation constant.OPT, change int32, total int32) {
	if change == 0 {
		return
	}
	data := &db.ItemRecord{UserID:uid, ItemID:itemID, Operation:uint8(operation), Change:change, Total:total, Time:time.Now()}
	if refID != 0 {
		data.RefID = refID
	}
	log.ChanRPC.Go(constant.LOG_COMMAND, data)
}

func (this *logHandler) AddSocialRecord(uid int32, socialID constant.SOCIAL_ID, refID int32, operation constant.OPT, change int32, total int32) {
	if change == 0 {
		return
	}
	var data interface{} = nil
	switch socialID {
	case constant.SOCIAL_ID_FAITH:
		data = &db.FaithRecord{UserID:uid, Operation:uint8(operation), Change:change, Total:total, RefID: refID, Time:time.Now()}
	case constant.SOCIAL_ID_DIAMOND:
		data = &db.DiamondRecord{UserID:uid, Operation:uint8(operation), Change:change, Total:total, RefID: refID, Time:time.Now()}
	case constant.SOCIAL_ID_POWER:
		data = &db.PowerRecord{UserID:uid, Operation:uint8(operation), Change:change, Total:total, RefID: refID, Time:time.Now()}
	case constant.SOCIAL_ID_GAYPOINT:
		data = &db.GayPointRecord{UserID:uid, Operation:uint8(operation), Change:change, Total:total, RefID: refID, Time:time.Now()}
	}
	log.ChanRPC.Go(constant.LOG_COMMAND, data)
}