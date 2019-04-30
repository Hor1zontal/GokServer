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
	"gok/constant"
	"gok/module/database"
	database2 "aliens/database"
)

var DBServiceProxy = &dbHandler{}

type dbHandler struct {

}

func (this *dbHandler) Insert(data interface{}, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_INSERT, data, dbHandler)
}

func (this *dbHandler) Update(data interface{}, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_UPDATE, data, dbHandler)
}

func (this *dbHandler) ForceUpdate(data interface{}, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_FUPDATE, data, dbHandler)
}


func (this *dbHandler) Delete(data interface{}, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_DELETE, data, dbHandler)
}

func (this *dbHandler) UpdateCondition(collectionName string, selectDoc interface{}, updateDoc interface{}, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_CONDITION_UPDATE, collectionName, selectDoc, updateDoc, dbHandler)
}

func (this *dbHandler) DeleteCondition(data interface{}, selectDoc interface{}, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_CONDITION_DELETE, data, selectDoc, dbHandler)
}


